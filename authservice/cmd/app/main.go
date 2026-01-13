package main

import (
	"ads/authservice/cmd/app/config"
	adaptergrpc "ads/authservice/internal/adapter/in/grpc"
	adapterdb "ads/authservice/internal/adapter/out/db"
	adapterph "ads/authservice/internal/adapter/out/hasher"
	adaptertg "ads/authservice/internal/adapter/out/jwt"
	"ads/authservice/internal/app/usecase"
	"ads/authservice/internal/generated/auth_v1"
	"ads/pkg/pg"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newLogger(level string) *slog.Logger {
	logLevel := parseLogLevel(level)
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

func newPostgresClient(cfg *config.Config) (*pg.PostgresClient, error) {
	pgConfig := pg.NewPostgresConfig(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.DBSSLMode, cfg.DBOpenConn,
		cfg.DBIdleConn, cfg.DBConnLifeTime,
	)

	pgClient, err := pg.NewPostgresClient(pgConfig)
	if err != nil {
		return nil, err
	}

	return pgClient, nil
}

func runServer(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	// Postgres client
	pgClient, err := newPostgresClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres client: %w", err)
	}

	// Close Postgres
	defer func() {
		logger.InfoContext(ctx, "closing postgres connection...")
		if err := pgClient.Close(); err != nil {
			logger.ErrorContext(ctx, "failed to close postgres",
				slog.Any("error", err),
			)
		}
	}()

	// Repositories
	accountRepo := adapterdb.NewAccountsRepository(pgClient)
	accountRoleRepo := adapterdb.NewAccountRolesRepository(pgClient)
	refreshSessionRepo := adapterdb.NewRefreshSessionsRepository(pgClient)
	passwordHasher := adapterph.NewBcryptHasher(cfg.PasswordCost)
	tokenGenerator := adaptertg.NewTokenGenerator(
		cfg.AccessSecret, cfg.RefreshSecret,
		cfg.AccessTTL, cfg.RefreshTTL,
	)

	// Use-cases
	registerUC := usecase.NewRegisterUC(accountRepo, accountRoleRepo, passwordHasher)
	loginUC := usecase.NewLoginUC(
		accountRepo, accountRoleRepo, refreshSessionRepo,
		passwordHasher, tokenGenerator, cfg.RefreshTTL,
	)
	logoutUC := usecase.NewLogoutUC(refreshSessionRepo, tokenGenerator)
	refreshSessionUC := usecase.NewRefreshSessionUC(
		accountRoleRepo, refreshSessionRepo,
		tokenGenerator, cfg.RefreshTTL,
	)
	validateAccessUC := usecase.NewValidateAccessTokenUC(
		accountRepo, tokenGenerator,
	)

	// Handler
	authHandler := adaptergrpc.NewAuthHandler(
		logger,
		registerUC,
		loginUC,
		logoutUC,
		refreshSessionUC,
		validateAccessUC,
	)

	// gRPC server
	gRPCServer := grpc.NewServer()
	auth_v1.RegisterAuthServiceServer(gRPCServer, authHandler)
	reflection.Register(gRPCServer)

	address := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen port %d: %w",
			cfg.GRPCPort, err,
		)
	}

	// Launch gRPC server
	errChan := make(chan error, 1)
	go func() {
		logger.InfoContext(
			ctx, "starting grpc server", slog.String("address", address))
		if err := gRPCServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	// Graceful shutdown
	select {
	case <-ctx.Done():
		logger.InfoContext(
			ctx, "received shutdown signal, stopping grpc server...",
		)
		gRPCServer.GracefulStop()
		return nil
	case err := <-errChan:
		return fmt.Errorf("grpc server failed: %w", err)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := newLogger(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx, cfg, logger); err != nil {
		logger.ErrorContext(
			ctx, "authservice failed", slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.Info("authservice stopped successfully")
}
