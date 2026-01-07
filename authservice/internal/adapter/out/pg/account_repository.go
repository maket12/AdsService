package pg

import (
	"ads/authservice/internal/adapter/out/pg/mapper"
	"ads/authservice/internal/adapter/out/pg/sqlc"
	"ads/authservice/internal/domain/model"
	"ads/authservice/pkg/errs"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountRepository struct {
	q *sqlc.Queries
}

func NewAccountsRepository(q *sqlc.Queries) *AccountRepository {
	return &AccountRepository{q: q}
}

func (r *AccountRepository) Create(ctx context.Context, account *model.Account) error {
	params := mapper.MapAccountToSQLCCreate(account)
	return r.q.CreateAccount(ctx, params)
}

func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*model.Account, error) {
	rawAcc, err := r.q.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("account", email)
		}
		return nil, err
	}

	account := mapper.MapSQLCToAccount(rawAcc)

	return account, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	rawAcc, err := r.q.GetAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("account", id)
		}
		return nil, err
	}

	account := mapper.MapSQLCToAccount(rawAcc)

	return account, nil
}

func (r *AccountRepository) MarkLogin(ctx context.Context, account *model.Account) error {
	var lastLoginTime time.Time
	if account.LastLoginAt() != nil {
		lastLoginTime = *account.LastLoginAt()
	}

	var params = sqlc.MarkAccountLoginParams{
		ID: account.ID(),
		LastLoginAt: sql.NullTime{
			Time:  lastLoginTime,
			Valid: true,
		},
		UpdatedAt: account.UpdatedAt(),
	}

	if err := r.q.MarkAccountLogin(ctx, params); err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) VerifyEmail(ctx context.Context, account *model.Account) error {
	params := sqlc.VerifyAccountEmailParams{
		ID:        account.ID(),
		UpdatedAt: account.UpdatedAt(),
	}
	return r.q.VerifyAccountEmail(ctx, params)
}
