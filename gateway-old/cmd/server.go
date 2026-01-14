package main

import (
	adminpb "ads/adminservice/presentation/grpc/pb"
	authpb "ads/authservice/presentation/grpc/pb"
	"ads/gateway-old/app/resolvers"
	"ads/gateway-old/config"
	"ads/gateway-old/infrastructure/authctx"
	"ads/gateway-old/pkg/logger"
	userpb "ads/userservice/presentation/grpc/pb"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthConn    *grpc.ClientConn
	UserConn    *grpc.ClientConn
	AdminConn   *grpc.ClientConn
	AuthClient  authpb.AuthServiceClient
	UserClient  userpb.UsersServiceClient
	AdminClient adminpb.AdminServiceClient
}

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		return
	}

	clients, err := initGrpcClients(cfg, log)
	if err != nil {
		log.Error("failed to initialize gRPC clients", "error", err)
		return
	}
	defer closeGrpcClients(clients, log)

	resolver := &resolvers.Resolver{
		AuthClient:  clients.AuthClient,
		UserClient:  clients.UserClient,
		AdminClient: clients.AdminClient,
	}

	server := startHTTPServer(resolver, log)

	waitForShutdown(server, clients, log)

	log.Info("👋 gateway stopped")
}

func initGrpcClients(cfg *config.Config, log *slog.Logger) (*GRPCClients, error) {
	log.Info("initializing gRPC clients...")

	clients := &GRPCClients{}

	authConn, err := grpc.NewClient(cfg.AuthGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", err)
	}
	clients.AuthConn = authConn
	clients.AuthClient = authpb.NewAuthServiceClient(authConn)

	userConn, err := grpc.NewClient(cfg.UserGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		authConn.Close()
		return nil, fmt.Errorf("user service: %w", err)
	}
	clients.UserConn = userConn
	clients.UserClient = userpb.NewUsersServiceClient(userConn)

	adminConn, err := grpc.NewClient(cfg.AdminGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		authConn.Close()
		userConn.Close()
		return nil, fmt.Errorf("admin service: %w", err)
	}
	clients.AdminConn = adminConn
	clients.AdminClient = adminpb.NewAdminServiceClient(adminConn)

	log.Info("✅ gRPC clients connected successfully",
		slog.String("auth", cfg.AuthGrpcAddr),
		slog.String("user", cfg.UserGrpcAddr),
		slog.String("admin", cfg.AdminGrpcAddr),
	)

	return clients, nil
}

func closeGrpcClients(clients *GRPCClients, log *slog.Logger) {
	log.Info("closing gRPC connections...")

	if clients.AuthConn != nil {
		if err := clients.AuthConn.Close(); err != nil {
			log.Error("failed to close auth connection", "error", err)
		} else {
			log.Info("✅ auth connection closed")
		}
	}

	if clients.UserConn != nil {
		if err := clients.UserConn.Close(); err != nil {
			log.Error("failed to close user connection", "error", err)
		} else {
			log.Info("✅ user connection closed")
		}
	}

	if clients.AdminConn != nil {
		if err := clients.AdminConn.Close(); err != nil {
			log.Error("failed to close admin connection", "error", err)
		} else {
			log.Info("✅ admin connection closed")
		}
	}
}

func startHTTPServer(resolver *resolvers.Resolver, log *slog.Logger) *http.Server {
	srv := handler.NewDefaultServer(resolvers.NewExecutableSchema(resolvers.Config{Resolvers: resolver}))

	authMD := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := authctx.MetadataFromRequest(r)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", authMD(srv))

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		log.Info("🚀 gateway running at http://localhost:8080/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed", "error", err)
		}
	}()

	return server
}

func waitForShutdown(server *http.Server, clients *GRPCClients, log *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("🛑 shutting down gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server gracefully", "error", err)
	} else {
		log.Info("✅ HTTP server stopped gracefully")
	}

	closeGrpcClients(clients, log)
}
