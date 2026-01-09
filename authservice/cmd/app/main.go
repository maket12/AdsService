package main

import (
	"ads/authservice/cmd/app/config"
	"ads/authservice/internal/adapter/in/grpc"
	adapterph "ads/authservice/internal/adapter/out/hasher"
	adaptertg "ads/authservice/internal/adapter/out/jwt"
	adapterpg "ads/authservice/internal/adapter/out/pg"
	"ads/authservice/internal/app/usecase"
	"ads/authservice/internal/generated/auth_v1"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	googlegrpc "google.golang.org/grpc"
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

func newPostgresClient(cfg *config.Config) (*adapterpg.PostgresClient, error) {
	pgConfig := adapterpg.NewPostgresConfig(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.SSLMode, cfg.OpenConn,
		cfg.IdleConn, cfg.ConnLifeTime,
	)

	pgClient, err := adapterpg.NewPostgresClient(pgConfig)
	if err != nil {
		return nil, err
	}

	return pgClient, nil
}

func runServer(ctx context.Context, config *config.Config, logger *slog.Logger) error {
	// Postgres client
	pgClient, err := newPostgresClient(config)
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
	accountRepo := adapterpg.NewAccountsRepository(pgClient)
	accountRoleRepo := adapterpg.NewAccountRolesRepository(pgClient)
	refreshSessionRepo := adapterpg.NewRefreshSessionsRepository(pgClient)
	passwordHasher := adapterph.NewBcryptHasher(config.PasswordCost)
	tokenGenerator := adaptertg.NewTokenGenerator(
		config.AccessSecret, config.RefreshSecret,
		config.AccessTTL, config.RefreshTTL,
	)

	// Use-cases
	registerUC := usecase.NewRegisterUC(accountRepo, accountRoleRepo, passwordHasher)
	loginUC := usecase.NewLoginUC(
		accountRepo, accountRoleRepo, refreshSessionRepo,
		passwordHasher, tokenGenerator, config.RefreshTTL,
	)
	logoutUC := usecase.NewLogoutUC(refreshSessionRepo, tokenGenerator)
	refreshSessionUC := usecase.NewRefreshSessionUC(
		accountRoleRepo, refreshSessionRepo,
		tokenGenerator, config.RefreshTTL,
	)
	validateAccessUC := usecase.NewValidateAccessTokenUC(
		accountRepo, tokenGenerator,
	)

	// Handler
	authHandler := grpc.NewAuthHandler(
		logger,
		registerUC,
		loginUC,
		logoutUC,
		refreshSessionUC,
		validateAccessUC,
	)

	// gRPC server
	gRPCServer := googlegrpc.NewServer()
	auth_v1.RegisterAuthServiceServer(gRPCServer, authHandler)

	address := fmt.Sprintf(":%d", config.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen port %d: %w",
			config.GRPCPort, err,
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
