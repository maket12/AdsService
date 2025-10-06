package main

import (
	"AdsService/adminservice/adapter/jwt"
	"AdsService/adminservice/adapter/pg"
	"AdsService/adminservice/app/usecase"
	"AdsService/adminservice/config"
	grpcinfra "AdsService/adminservice/infrastructure/grpc"
	"AdsService/adminservice/infrastructure/postgres"
	"AdsService/adminservice/pkg/logger"
	adminsvc "AdsService/adminservice/presentation/grpc"
	pb "AdsService/adminservice/presentation/grpc/pb"
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
	log.Info("🚀 starting adminservice", "port", cfg.GRPCPort)

	db, err := initDependencies(cfg, log)
	if err != nil {
		log.Error("failed to initialize dependencies", "error", err)
		return
	}
	defer closeDependencies(log)

	adminService := initServices(db, cfg, log)

	server := startGRPCServer(adminService, cfg, log)

	waitForShutdown(server, log)

	log.Info("👋 adminservice stopped")
}

// Инициализация зависимостей
func initDependencies(cfg *config.Config, log *slog.Logger) (*gorm.DB, error) {
	log.Info("initializing database connections...")

	if err := postgres.InitDB(cfg, log); err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	return postgres.DB, nil
}

// Закрытие соединений
func closeDependencies(log *slog.Logger) {
	log.Info("closing database connections...")
}

// Инициализация сервисов
func initServices(db *gorm.DB, cfg *config.Config, log *slog.Logger) *adminsvc.AdminService {
	usersRepo := pg.NewUsersRepo(db, log)
	profilesRepo := pg.NewProfilesRepo(db, log)

	assignRoleUC := &usecase.AssignRoleUC{Users: usersRepo}
	getUserUC := &usecase.GetUserUC{Users: usersRepo}
	adminBanUserUC := &usecase.BanUserUC{Users: usersRepo}
	adminUnbanUserUC := &usecase.UnbanUserUC{Users: usersRepo}
	adminGetProfileUC := &usecase.GetProfileUC{Profiles: profilesRepo}
	adminGetProfilesUC := &usecase.GetProfilesUC{Profiles: profilesRepo}

	return &adminsvc.AdminService{
		AssignRoleUC:  assignRoleUC,
		GetUserUC:     getUserUC,
		BanUserUC:     adminBanUserUC,
		UnbanUserUC:   adminUnbanUserUC,
		GetProfileUC:  adminGetProfileUC,
		GetProfilesUC: adminGetProfilesUC,
	}
}

func startGRPCServer(adminService *adminsvc.AdminService, cfg *config.Config, log *slog.Logger) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	tokensRepo := jwt.NewTokenRepository(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
	authInterceptor := &grpcinfra.AuthInterceptor{Tokens: tokensRepo}

	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryAuth()))
	pb.RegisterAdminServiceServer(server, adminService)
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
