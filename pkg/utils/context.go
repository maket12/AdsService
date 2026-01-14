package utils

import (
	"ads/pkg/errs"
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
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
