package port

import (
	"ads/authservice/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type AccountsRepository interface {
	Create(ctx context.Context, account *model.Account) error
	GetByEmail(ctx context.Context, email string) (*model.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	MarkLogin(ctx context.Context, account *model.Account) error
	VerifyEmail(ctx context.Context, account *model.Account) error
}
