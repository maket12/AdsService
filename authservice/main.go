package main

import (
	"AdsService/authservice/auth"
	pb "AdsService/authservice/proto"
	"AdsService/infra/database"
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := database.User{
		Email:    req.Email,
		Password: string(hashed),
		Role:     "user",
	}

	if err := database.AddUser(&user); err != nil {
		return nil, status.Errorf(codes.Internal, "error creating user: %v", err)
	}

	if _, err := database.AddProfile(
		user.ID, "undefined", "undefined",
	); err != nil {
		return nil, status.Errorf(codes.Internal, "error creating profile: %v", err)
	}

	access, err := auth.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue access: %v", err)
	}

	refresh, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue refresh: %v", err)
	}

	refreshClaims, err := auth.ParseRefreshToken(refresh)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse refresh: %v", err)
	}

	if err := database.InsertSession(&database.Session{
		UserID:    user.ID,
		JTI:       refreshClaims.ID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: refreshClaims.ExpiresAt.Time.UTC(),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "save session: %v", err)
	}

	return &pb.AuthResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, err := database.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	access, err := auth.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue access: %v", err)
	}

	refresh, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue refresh: %v", err)
	}

	rClaims, err := auth.ParseRefreshToken(refresh)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse refresh: %v", err)
	}

	if err := database.InsertSession(&database.Session{
		UserID:    user.ID,
		JTI:       rClaims.ID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: rClaims.ExpiresAt.Time.UTC(),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "save session: %v", err)
	}

	return &pb.AuthResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := auth.ParseAccessToken(req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:  false,
			UserId: 0,
		}, err
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
	}, nil
}

func (s *AuthService) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	if err := database.EnhanceUser(req.UserId); err != nil {
		return &pb.AssignRoleResponse{
			Ok: false,
		}, err
	}

	return &pb.AssignRoleResponse{
		Ok: true,
	}, nil
}

func (s *AuthService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := database.GetUserByID(req.UserId)
	if err != nil {
		return &pb.GetUserResponse{
			UserId: req.UserId,
			Email:  "",
			Role:   "unknown",
		}, err
	}

	return &pb.GetUserResponse{
		UserId: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

//func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {
//	refresh_token := req.RefreshToken
//
//	tok, err := jwt.ParseWithClaims(refresh_token, &auth.AccessClaims{}, func(t *jwt.Token) (interface{}, error) { return s.RefreshToken() })
//}

func main() {
	database.InitDB()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &AuthService{})
	reflection.Register(s)

	log.Println("Starting server on port 50051...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server!")
	}
}
