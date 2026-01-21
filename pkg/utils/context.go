package utils

import (
	"ads/pkg/errs"
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// Context keys
type contextKey string

const (
	AccountIDKey contextKey = "account_id"
	RoleKey      contextKey = "role"
)

// Custom errors
var (
	ErrMetadataIsMissing     = errors.New("metadata is missing")
	ErrAccountIDNotSpecified = errors.New("account id not found in metadata")
	ErrInvalidAccountID      = errors.New("metadata contains invalid account id")
)

// ExtractAccountID Extracts account id from incoming context
func ExtractAccountID(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, errs.NewNotAuthenticatedErrorWithReason(ErrMetadataIsMissing)
	}

	vals := md.Get("x-account-id")
	if len(vals) == 0 {
		return uuid.Nil, errs.NewNotAuthenticatedErrorWithReason(ErrAccountIDNotSpecified)
	}

	id, err := uuid.Parse(vals[0])
	if err != nil {
		return uuid.Nil, errs.NewNotAuthenticatedErrorWithReason(ErrInvalidAccountID)
	}

	return id, nil
}

// PackAccountIDForGRPC Packs account id into outgoing context (metadata)
func PackAccountIDForGRPC(ctx context.Context, accountID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "x-account-id", accountID)
}

// SetAccountIDInCtx Sets account id in context (gateway)
func SetAccountIDInCtx(ctx context.Context, accountID string) context.Context {
	return context.WithValue(ctx, AccountIDKey, accountID)
}

// SetAccountRoleInCtx Sets account role in context (gateway)
func SetAccountRoleInCtx(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}
