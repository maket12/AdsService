package port

import (
	"ads/authservice/domain/entity"
	"context"
)

type TokenRepository interface {
	GenerateAccessToken(ctx context.Context, userID uint64, email, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uint64) (string, error)
	ParseAccessToken(ctx context.Context, tokenStr string) (*entity.AccessClaims, error)
	ParseRefreshToken(ctx context.Context, tokenStr string) (*entity.RefreshClaims, error)
}
