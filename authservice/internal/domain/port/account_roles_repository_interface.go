package port

import (
	"ads/authservice/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type AccountRolesRepository interface {
	Create(ctx context.Context, accountRole *model.AccountRole) error
	Get(ctx context.Context, accountID uuid.UUID) (*model.AccountRole, error)
	Update(ctx context.Context, accountRole *model.AccountRole) error
	Delete(ctx context.Context, accountID uuid.UUID) error
}
