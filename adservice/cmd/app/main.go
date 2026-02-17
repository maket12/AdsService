package app

import (
	"ads/adservice/cmd/app/config"
	adaptergrpc "ads/adservice/internal/adapter/in/grpc"
	adapterpg "ads/adservice/internal/adapter/out/postgres"
	adaptermongo "ads/adservice/internal/adapter/out/mongodb"
	"ads/adservice/internal/app/usecase"
	adapterph "ads/authservice/internal/adapter/out/hasher"
	adaptertg "ads/authservice/internal/adapter/out/jwt"
	adapterdb "ads/authservice/internal/adapter/out/postgres"
	"ads/pkg/generated/ad_v1"
	pkgpostgres "ads/pkg/postgres"
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

func newPostgresClient(cfg *config.Config) (*pkgpostgres.Client, error) {
	pgConfig := pkgpostgres.NewConfig(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.DBSSLMode, cfg.DBOpenConn,
		cfg.DBIdleConn, cfg.DBConnLifeTime,
	)

	pgClient, err := pkgpostgres.NewClient(pgConfig)
	if err != nil {
		return nil, err
	}

	return pgClient, nil
}

func closePostgresClient(
	ctx context.Context,
	logger *slog.Logger,
	pgClient *pkgpostgres.Client,
) {
	logger.InfoContext(ctx, "closing postgres connection...")
	if err := pgClient.Close(); err != nil {
		logger.ErrorContext(ctx, "failed to close postgres",
			slog.Any("error", err),
		)
	}
}

func runServer(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	// Postgres client
	pgClient, err := newPostgresClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres client: %w", err)
	}

	// Close Postgres
	defer closePostgresClient(ctx, logger, pgClient)

	// Repositories
	adRepo := adapterpg.NewAdRepository(pgClient)
	mediaRepo := adaptermongo.NewMediaRepository(pgClient)

	// Use-cases
	createAdUC := usecase.NewCreateAdUC(adRepo, mediaRepo)
	getAdUC := usecase.NewGetAdUC(adRepo, mediaRepo)
	updateAdUC := usecase.NewUpdateAdUC(adRepo, mediaRepo)
	publishAdUC := usecase.NewPublishAdUC(adRepo)
	rejectAdUC := usecase.NewRejectAdUC(adRepo)
	deleteAdUC := usecase.NewDeleteAdUC(adRepo)
	deleteAllAdsUC := usecase.NewDeleteAllAdsUC(adRepo)

	// Handler
	adHandler := adaptergrpc.NewAdHandler(
		logger,
		createAdUC,
		getAdUC,
		updateAdUC,
		publishAdUC,
		rejectAdUC,
		deleteAdUC,
		deleteAllAdsUC,
	)

	// gRPC server
	gRPCServer := grpc.NewServer()
	ad_v1.RegisterAdServiceServer(gRPCServer, adHandler)
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
