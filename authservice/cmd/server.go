package main

import (
	"AdsService/authservice/adapter/jwt"
	"AdsService/authservice/adapter/pg"
	"AdsService/authservice/app/usecase"
	"AdsService/authservice/config"
	"AdsService/authservice/infrastructure/postgres"
	"AdsService/authservice/pkg/logger"
	authgrpc "AdsService/authservice/presentation/grpc"
	pb "AdsService/authservice/presentation/grpc/pb"
	"fmt"
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
	log.Info("🚀 starting authservice", "port", cfg.GRPCPort)

	db, err := initDependencies(cfg, log)
	if err != nil {
		log.Error("failed to initialize dependencies", "error", err)
		return
	}
	defer closeDependencies(log)

	authService := initServices(db, cfg, log)

	server := startGRPCServer(authService, cfg, log)

	waitForShutdown(server, log)

	log.Info("👋 authservice stopped")
}

// Инициализация зависимостей
func initDependencies(cfg *config.Config, log *slog.Logger) (*gorm.DB, error) {
	log.Info("initializing database connections...")

	if err := postgres.InitDB(cfg, log); err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	return postgres.DB, nil
}

func closeDependencies(log *slog.Logger) {
	log.Info("closing database connections...")
}

func initServices(db *gorm.DB, cfg *config.Config, log *slog.Logger) *authgrpc.AuthService {
	usersRepo := pg.NewUsersRepo(db, log)
	sessionsRepo := pg.NewSessionsRepo(db, log)
	profilesRepo := pg.NewProfilesRepo(db, log)
	tokensRepo := jwt.NewTokenRepository(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, log)

	registerUC := &usecase.RegisterUC{
		Users:    usersRepo,
		Sessions: sessionsRepo,
		Tokens:   tokensRepo,
		Profiles: profilesRepo,
	}
	loginUC := &usecase.LoginUC{
		Users:    usersRepo,
		Sessions: sessionsRepo,
		Tokens:   tokensRepo,
	}
	validateUC := &usecase.ValidateTokenUC{Tokens: tokensRepo}

	return authgrpc.NewAuthService(registerUC, loginUC, validateUC)
}

func startGRPCServer(authService *authgrpc.AuthService, cfg *config.Config, log *slog.Logger) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authService)
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
