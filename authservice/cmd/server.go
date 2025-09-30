package main

import (
	"AdsService/authservice/adapter/jwt"
	"AdsService/authservice/adapter/pg"
	"AdsService/authservice/app/usecase"
	"AdsService/authservice/infrastructure/postgres"
	authgrpc "AdsService/authservice/presentation/grpc"
	pb "AdsService/authservice/presentation/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	// Init DB
	postgres.InitDB()
	db := postgres.DB

	// Repositories
	usersRepo := pg.NewUsersRepo(db)
	sessionsRepo := pg.NewSessionsRepo(db)
	profilesRepo := pg.NewProfilesRepo(db)
	tokensRepo := jwt.NewTokenRepository()

	// Usecases
	registerUC := &usecase.RegisterUC{Users: usersRepo, Sessions: sessionsRepo, Tokens: tokensRepo, Profiles: profilesRepo}
	loginUC := &usecase.LoginUC{Users: usersRepo, Sessions: sessionsRepo, Tokens: tokensRepo}
	validateUC := &usecase.ValidateTokenUC{Tokens: tokensRepo}

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	authService := authgrpc.NewAuthService(registerUC, loginUC, validateUC)
	pb.RegisterAuthServiceServer(s, authService)
	reflection.Register(s)

	log.Println("AuthService gRPC running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
