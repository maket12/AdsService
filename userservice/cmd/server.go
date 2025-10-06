package main

import (
	"AdsService/userservice/adapter/jwt"
	"AdsService/userservice/adapter/mongo"
	"AdsService/userservice/adapter/pg"
	"AdsService/userservice/app/usecase"
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
	profilesRepo := pg.NewProfilesRepo(db)
	photosRepo := mongo.NewPhotoRepo(mongodb.Bucket)
	tokensRepo := jwt.NewTokenRepository()

	// Init usecases
	addProfileUC := &usecase.AddProfileUC{Profiles: profilesRepo}
	updateProfileUC := &usecase.UpdateProfileUC{Profiles: profilesRepo}
	uploadPhotoUC := &usecase.UploadPhotoUC{Profiles: profilesRepo, Photos: photosRepo}
	changeSettingsUC := &usecase.ChangeSettingsUC{Profiles: profilesRepo}
	changeSubscriptionsUC := &usecase.ChangeSubscriptionsUC{Profiles: profilesRepo}
	getProfileUC := &usecase.GetProfileUC{Profiles: profilesRepo}

	// Init service
	userService := &usersvc.UserService{
		AddProfileUC:          addProfileUC,
		UpdateProfileUC:       updateProfileUC,
		UploadPhotoUC:         uploadPhotoUC,
		ChangeSettingUC:       changeSettingsUC,
		ChangeSubscriptionsUC: changeSubscriptionsUC,
		GetProfileUC:          getProfileUC,
	}

	// Interceptor with JWT
	authInterceptor := &grpcinfra.AuthInterceptor{Tokens: tokensRepo}

	// gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryAuth()))
	pb.RegisterUsersServiceServer(s, userService)
	reflection.Register(s)

	log.Println("UserService gRPC running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
