package port

import (
	"context"

	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateAccessToken(ctx context.Context, accountID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, accountID, sessionID uuid.UUID) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (accountID uuid.UUID, role string, err error)
	ValidateRefreshToken(ctx context.Context, token string) (accountID uuid.UUID, sessionID uuid.UUID, err error)
}
