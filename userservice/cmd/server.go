package main

import (
	"AdsService/userservice/adapter/jwt"
	"AdsService/userservice/adapter/mongo"
	"AdsService/userservice/adapter/pg"
	"AdsService/userservice/app/usecase/admin_uc"
	"AdsService/userservice/app/usecase/profile_uc"
	"AdsService/userservice/app/usecase/user_uc"
	grpcinfra "AdsService/userservice/infrastructure/grpc"
	"AdsService/userservice/infrastructure/mongodb"
	"AdsService/userservice/infrastructure/postgres"
	usersvc "AdsService/userservice/presentation/grpc"
	pb "AdsService/userservice/presentation/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	// Init DB connections
	postgres.InitDB()
	db := postgres.DB

	mongodb.InitMongoDB()
	defer mongodb.CloseMongoDB()

	// Init repositories
	usersRepo := pg.NewUsersRepo(db)
	profilesRepo := pg.NewProfilesRepo(db)
	photosRepo := mongo.NewPhotoRepo(mongodb.Bucket)
	tokensRepo := jwt.NewTokenRepository()

	// Init usecases
	assignRoleUC := &user_uc.AssignRoleUC{Users: usersRepo}
	getUserUC := &user_uc.GetUserUC{Users: usersRepo}
	addProfileUC := &profile_uc.AddProfileUC{Profiles: profilesRepo}
	updateProfileUC := &profile_uc.UpdateProfileUC{Profiles: profilesRepo}
	uploadPhotoUC := &profile_uc.UploadPhotoUC{Profiles: profilesRepo, Photos: photosRepo}
	changeSettingsUC := &profile_uc.ChangeSettingsUC{Profiles: profilesRepo}
	changeSubscriptionsUC := &profile_uc.ChangeSubscriptionsUC{Profiles: profilesRepo}
	getProfileUC := &profile_uc.GetProfileUC{Profiles: profilesRepo}
	adminBanUserUC := &admin_uc.AdminBanUserUC{Users: usersRepo}
	adminUnbanUserUC := &admin_uc.AdminUnbanUserUC{Users: usersRepo}
	adminGetProfileUC := &admin_uc.AdminGetProfileUC{Profiles: profilesRepo}
	adminGetProfilesUC := &admin_uc.AdminGetProfilesUC{Profiles: profilesRepo}

	// Init service
	userService := &usersvc.UserService{
		AssignRoleUC:          assignRoleUC,
		GetUserUC:             getUserUC,
		AddProfileUC:          addProfileUC,
		UpdateProfileUC:       updateProfileUC,
		UploadPhotoUC:         uploadPhotoUC,
		ChangeSettingUC:       changeSettingsUC,
		ChangeSubscriptionsUC: changeSubscriptionsUC,
		GetProfileUC:          getProfileUC,
		AdminBanUserUC:        adminBanUserUC,
		AdminUnbanUserUC:      adminUnbanUserUC,
		AdminGetProfileUC:     adminGetProfileUC,
		AdminGetProfilesUC:    adminGetProfilesUC,
	}

	// Interceptor with JWT
	authInterceptor := &grpcinfra.AuthInterceptor{Tokens: tokensRepo}

	// gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryAuth()))
	pb.RegisterUserServiceServer(s, userService)
	reflection.Register(s)

	log.Println("UserService gRPC running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
