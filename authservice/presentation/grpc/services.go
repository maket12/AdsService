package grpc

import (
	dto2 "ads/authservice/internal/app/dto"
	usecase2 "ads/authservice/internal/app/usecase"
	pb "ads/authservice/presentation/grpc/pb"
	"context"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer

	RegisterUC      *usecase2.RegisterUC
	LoginUC         *usecase2.LoginUC
	ValidateTokenUC *usecase2.ValidateTokenUC
}

func NewAuthService(
	register *usecase2.RegisterUC,
	login *usecase2.LoginUC,
	validate *usecase2.ValidateTokenUC,
) *AuthService {
	return &AuthService{
		RegisterUC:      register,
		LoginUC:         login,
		ValidateTokenUC: validate,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	in := dto2.Register{Email: req.Email, Password: req.Password}
	tokens, err := s.RegisterUC.Execute(ctx, in)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	in := dto2.Login{Email: req.Email, Password: req.Password}
	tokens, err := s.LoginUC.Execute(ctx, in)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	in := dto2.ValidateToken{AccessToken: req.AccessToken}
	out, err := s.ValidateTokenUC.Execute(ctx, in)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, err
	}
	return &pb.ValidateTokenResponse{
		Valid:  out.Valid,
		UserId: out.UserID,
	}, nil
}
