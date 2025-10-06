package main

import (
	"AdsService/userservice/adapter/jwt"
	mongorepo "AdsService/userservice/adapter/mongo"
	"AdsService/userservice/adapter/pg"
	"AdsService/userservice/app/usecase"
	"AdsService/userservice/config"
	grpcinfra "AdsService/userservice/infrastructure/grpc"
	"AdsService/userservice/infrastructure/mongodb"
	"AdsService/userservice/infrastructure/postgres"
	"AdsService/userservice/pkg/logger"
	usersvc "AdsService/userservice/presentation/grpc"
	pb "AdsService/userservice/presentation/grpc/pb"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New()
	log.Info("🚀 starting userservice", "port", cfg.GRPCPort)

	db, mongoBucket, err := initDependencies(cfg, log)
	if err != nil {
		log.Error("failed to initialize dependencies", "error", err)
		return
	}
	defer closeDependencies(log)

	userService := initServices(db, mongoBucket, log)

	server := startGRPCServer(userService, cfg, log)

	waitForShutdown(server, log)

	log.Info("👋 userservice stopped")
}

func initDependencies(cfg *config.Config, log *slog.Logger) (*gorm.DB, *mongo.GridFSBucket, error) {
	log.Info("initializing database connections...")

	if err := postgres.InitDB(cfg, log); err != nil {
		return nil, nil, fmt.Errorf("postgres: %w", err)
	}

	if err := mongodb.InitMongoDB(cfg, log); err != nil {
		return nil, nil, fmt.Errorf("mongodb: %w", err)
	}

	return postgres.DB, mongodb.Bucket, nil
}

// Закрытие всех соединений
func closeDependencies(log *slog.Logger) {
	log.Info("closing database connections...")
	mongodb.CloseMongoDB(log)
}

func initServices(db *gorm.DB, mongoBucket *mongo.GridFSBucket, log *slog.Logger) *usersvc.UserService {
	profilesRepo := pg.NewProfilesRepo(db, log)
	photosRepo := mongorepo.NewPhotoRepo(mongoBucket, log)

	// Use cases
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

func startGRPCServer(userService *usersvc.UserService, cfg *config.Config, log *slog.Logger) *grpc.Server {
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
