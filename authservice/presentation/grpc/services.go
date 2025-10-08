package grpc

import (
	pb "ads/authservice/presentation/grpc/pb"
	"context"

	"ads/authservice/app/dto"
	"ads/authservice/app/usecase"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer

	RegisterUC      *usecase.RegisterUC
	LoginUC         *usecase.LoginUC
	ValidateTokenUC *usecase.ValidateTokenUC
}

func NewAuthService(
	register *usecase.RegisterUC,
	login *usecase.LoginUC,
	validate *usecase.ValidateTokenUC,
) *AuthService {
	return &AuthService{
		RegisterUC:      register,
		LoginUC:         login,
		ValidateTokenUC: validate,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	in := dto.Register{Email: req.Email, Password: req.Password}
	tokens, err := s.RegisterUC.Execute(ctx, in)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	in := dto.Login{Email: req.Email, Password: req.Password}
	tokens, err := s.LoginUC.Execute(ctx, in)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	in := dto.ValidateToken{AccessToken: req.AccessToken}
	out, err := s.ValidateTokenUC.Execute(in)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, err
	}
	return &pb.ValidateTokenResponse{
		Valid:  out.Valid,
		UserId: out.UserID,
	}, nil
}
