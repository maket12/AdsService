package ports

import (
	"ads/authservice/domain/entity"
	"context"
)

type UserRepository interface {
	CheckUserExist(ctx context.Context, email string) (bool, error)
	AddUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
