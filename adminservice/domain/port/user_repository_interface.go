package port

import (
	"ads/adminservice/domain/entity"
	"context"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID uint64) (*entity.User, error)
	GetUserRole(ctx context.Context, userID uint64) (string, error)
	EnhanceUser(ctx context.Context, userID uint64) error
	BanUser(ctx context.Context, userID uint64) error
	UnbanUser(ctx context.Context, userID uint64) error
}
