package main

import (
	"ads/userservice/adapter/jwt"
	mongorepo "ads/userservice/adapter/mongo"
	"ads/userservice/adapter/pg"
	"ads/userservice/app/usecase"
	grpcinfra "ads/userservice/infrastructure/grpc"
	"ads/userservice/infrastructure/mongodb"
	"ads/userservice/infrastructure/postgres"
	"ads/userservice/pkg"
	usersvc "ads/userservice/presentation/grpc"
	pb "ads/userservice/presentation/grpc/pb"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := pkg.New()

	cfg, err := pkg.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		return
	}

	log.Info("🚀 starting userservice", "port", cfg.GRPCPort)

	db, mongoData, err := initDependencies(cfg, log)
	if err != nil {
		log.Error("failed to initialize dependencies", "error", err)
		return
	}
	defer closeDependencies(mongoData)

	userService := initServices(db, mongoData, log)

	server := startGRPCServer(userService, cfg, log)

	waitForShutdown(server, log)

	log.Info("👋 userservice stopped")
}

func initDependencies(cfg *pkg.Config, log *slog.Logger) (*gorm.DB, *mongodb.MongoData, error) {
	log.Info("initializing database connections...")

	db, err := postgres.InitDB(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres: %w", err)
	}

	mongoData, err := mongodb.InitMongoDB(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("mongodb: %w", err)
	}

	return db, mongoData, nil
}

func closeDependencies(mongoData *mongodb.MongoData) {
	fmt.Print("closing database connections...")
	mongodb.CloseMongoDB(mongoData)
}

func initServices(db *gorm.DB, mongoData *mongodb.MongoData, log *slog.Logger) *usersvc.UserService {
	profilesRepo := pg.NewProfilesRepo(db, log)
	photosRepo := mongorepo.NewPhotoRepo(mongoData.Bucket, log)

	addProfileUC := &usecase.AddProfileUC{Profiles: profilesRepo}
	updateProfileUC := &usecase.UpdateProfileUC{Profiles: profilesRepo}
	uploadPhotoUC := &usecase.UploadPhotoUC{Profiles: profilesRepo, Photos: photosRepo}
	changeSettingsUC := &usecase.ChangeSettingsUC{Profiles: profilesRepo}
	changeSubscriptionsUC := &usecase.ChangeSubscriptionsUC{Profiles: profilesRepo}
	getProfileUC := &usecase.GetProfileUC{Profiles: profilesRepo}

	return &usersvc.UserService{
		AddProfileUC:          addProfileUC,
		UpdateProfileUC:       updateProfileUC,
		UploadPhotoUC:         uploadPhotoUC,
		ChangeSettingUC:       changeSettingsUC,
		ChangeSubscriptionsUC: changeSubscriptionsUC,
		GetProfileUC:          getProfileUC,
	}
}

func startGRPCServer(userService *usersvc.UserService, cfg *pkg.Config, log *slog.Logger) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	tokensRepo := jwt.NewTokenRepository(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
	authInterceptor := &grpcinfra.AuthInterceptor{Tokens: tokensRepo}

	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryAuth()))
	pb.RegisterUsersServiceServer(server, userService)
	reflection.Register(server)

	go func() {
		log.Info("📡 gRPC server listening", "port", cfg.GRPCPort)
		if err := server.Serve(lis); err != nil {
			log.Error("gRPC server failed", "error", err)
		}
	}()

	return server
}

func waitForShutdown(server *grpc.Server, log *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("🛑 shutting down server...")

	stopped := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info("✅ server stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Warn("⏰ forcing server shutdown")
		server.Stop()
	}
}
