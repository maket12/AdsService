package main

import (
	"AdsService/infra/gateway/internal/clients"
	"context"
	"log"
	"net/http"

	"AdsService/infra/gateway/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	ctx := context.Background()
	cls := clients.New(ctx)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{C: cls}}))

	http.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "httpRequest", r)
		srv.ServeHTTP(w, r.WithContext(ctx))
	}))
	http.Handle("/", playground.Handler("GraphQL", "/graphql"))

	log.Println("gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
