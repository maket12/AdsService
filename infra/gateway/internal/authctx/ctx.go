package authctx

import (
	"context"
	"google.golang.org/grpc/metadata"
	"net/http"
)

func MetadataFromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	if auth := r.Header.Get("Authorization"); auth != "" {
		md := metadata.Pairs("authorization", auth)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}
