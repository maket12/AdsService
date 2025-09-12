package authmw

import (
	"AdsService/authservice/auth"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const (
	ctxUserID ctxKey = "uid"
	ctxRole   ctxKey = "role"
)

func UnaryAuth() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no metadata")
		}

		vals := md.Get("authorization")
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}

		raw := strings.TrimSpace(vals[0])
		if strings.HasPrefix(strings.ToLower(raw), "bearer ") {
			raw = strings.TrimSpace(raw[7:])
		}

		claims, err := auth.ParseAccessToken(raw)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, ctxUserID, claims.UserID)
		ctx = context.WithValue(ctx, ctxRole, claims.Role)

		if strings.HasPrefix(info.FullMethod, "/userservice.UserService/Admin") && claims.Role != "admin" {
			return nil, status.Error(codes.PermissionDenied, "admin only")
		}

		return handler(ctx, req)
	}
}

func GetUserIDFromCtx(ctx context.Context) (uint64, error) {
	v := ctx.Value(ctxUserID)
	id, ok := v.(uint64)
	if !ok || id == 0 {
		return 0, status.Error(codes.Unauthenticated, "no user in context")
	}
	return id, nil
}
