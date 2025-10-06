package main

import (
	"AdsService/adminservice/adapter/jwt"
	"AdsService/adminservice/adapter/pg"
	"AdsService/adminservice/app/usecase"
	grpcinfra "AdsService/adminservice/infrastructure/grpc"
	"AdsService/adminservice/infrastructure/postgres"
	adminsvc "AdsService/adminservice/presentation/grpc"
	pb "AdsService/adminservice/presentation/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	postgres.InitDB()
	db := postgres.DB

	usersRepo := pg.NewUsersRepo(db)
	profilesRepo := pg.NewProfilesRepo(db)
	tokensRepo := jwt.NewTokenRepository()

	assignRoleUC := &usecase.AssignRoleUC{Users: usersRepo}
	getUserUC := &usecase.GetUserUC{Users: usersRepo}
	adminBanUserUC := &usecase.BanUserUC{Users: usersRepo}
	adminUnbanUserUC := &usecase.UnbanUserUC{Users: usersRepo}
	adminGetProfileUC := &usecase.GetProfileUC{Profiles: profilesRepo}
	adminGetProfilesUC := &usecase.GetProfilesUC{Profiles: profilesRepo}

	userService := &adminsvc.AdminService{
		AssignRoleUC:  assignRoleUC,
		GetUserUC:     getUserUC,
		BanUserUC:     adminBanUserUC,
		UnbanUserUC:   adminUnbanUserUC,
		GetProfileUC:  adminGetProfileUC,
		GetProfilesUC: adminGetProfilesUC,
	}

	authInterceptor := &grpcinfra.AuthInterceptor{Tokens: tokensRepo}

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryAuth()))
	pb.RegisterAdminServiceServer(s, userService)
	reflection.Register(s)

	log.Println("UserService gRPC running on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
