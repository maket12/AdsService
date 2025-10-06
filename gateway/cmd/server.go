package main

import (
	authpb "AdsService/authservice/presentation/grpc/pb"
	"AdsService/gateway/app/resolvers"
	"AdsService/gateway/infrastructure/authctx"
	userpb "AdsService/userservice/presentation/grpc/pb"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	authAddr := getenv("AUTH_GRPC_ADDR", "authservice:50051")
	userAddr := getenv("USER_GRPC_ADDR", "userservice:50052")

	// gRPC connections
	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to authservice: %v", err)
	}
	defer authConn.Close()

	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to userservice: %v", err)
	}
	defer userConn.Close()

	// gRPC clients
	authClient := authpb.NewAuthServiceClient(authConn)
	userClient := userpb.NewUsersServiceClient(userConn)

	// GraphQL resolver
	resolver := &resolvers.Resolver{
		AuthClient: authClient,
		UserClient: userClient,
	}

	// GraphQL server
	srv := handler.NewDefaultServer(resolvers.NewExecutableSchema(resolvers.Config{Resolvers: resolver}))

	// HTTP middleware: кладём Authorization -> gRPC metadata в контекст
	authMD := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := authctx.MetadataFromRequest(r) // берёт Authorization из заголовка и кладёт md "authorization"
			next.ServeHTTP(w, r.WithContext(ctx)) // этот ctx попадёт во все резолверы
		})
	}

	// Routes
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", authMD(srv)) // <--- оборачиваем /query

	log.Println("🚀 Gateway running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
