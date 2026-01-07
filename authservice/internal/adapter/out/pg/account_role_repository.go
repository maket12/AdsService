package pg

import (
	"ads/authservice/internal/adapter/out/pg/mapper"
	"ads/authservice/internal/adapter/out/pg/sqlc"
	"ads/authservice/internal/domain/model"
	"ads/authservice/pkg/errs"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type AccountRoleRepository struct {
	q *sqlc.Queries
}

func NewAccountRolesRepository(q *sqlc.Queries) *AccountRoleRepository {
	return &AccountRoleRepository{q: q}
}

func (r *AccountRoleRepository) Create(ctx context.Context, accountRole *model.AccountRole) error {
	params := mapper.MapAccountRoleToSQLCCreate(accountRole)
	return r.q.CreateAccountRole(ctx, params)
}

func (r *AccountRoleRepository) Get(ctx context.Context, accountID uuid.UUID) (*model.AccountRole, error) {
	rawAccRole, err := r.q.GetAccountRole(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("account_role", accountID)
		}
		return nil, err
	}

	accountRole := mapper.MapSQLCToAccountRole(rawAccRole)

	return accountRole, nil
}

func (r *AccountRoleRepository) Update(ctx context.Context, accountRole *model.AccountRole) error {
	var params = sqlc.UpdateAccountRoleParams{
		AccountID: accountRole.AccountID(),
		Role:      sqlc.RoleType(accountRole.Role()),
	}

	if err := r.q.UpdateAccountRole(ctx, params); err != nil {
		return err
	}

	return nil
}

func (r *AccountRoleRepository) Delete(ctx context.Context, accountID uuid.UUID) error {
	return r.q.DeleteAccountRole(ctx, accountID)
}
