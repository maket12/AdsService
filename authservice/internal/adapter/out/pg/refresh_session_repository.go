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

type RefreshSessionRepository struct {
	q *sqlc.Queries
}

func NewRefreshSessionsRepository(q *sqlc.Queries) *RefreshSessionRepository {
	return &RefreshSessionRepository{q: q}
}

func (r *RefreshSessionRepository) Create(ctx context.Context, session *model.RefreshSession) error {
	params := mapper.MapRefreshSessionToSQLCCreate(session)
	return r.q.CreateRefreshSession(ctx, params)
}

func (r *RefreshSessionRepository) GetByHash(ctx context.Context, tokenHash string) (*model.RefreshSession, error) {
	rawSession, err := r.q.GetRefreshSessionByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("refresh_session", tokenHash)
		}
		return nil, err
	}
	refreshSession := mapper.MapSQLCToRefreshSession(rawSession)
	return refreshSession, nil
}

func (r *RefreshSessionRepository) GetByID(ctx context.Context, tokenID uuid.UUID) (*model.RefreshSession, error) {
	rawSession, err := r.q.GetRefreshSessionByID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("refresh_session", tokenID)
		}
		return nil, err
	}
	refreshSession := mapper.MapSQLCToRefreshSession(rawSession)
	return refreshSession, nil
}

func (r *RefreshSessionRepository) Revoke(ctx context.Context, session *model.RefreshSession) error {
	var (
		revokedAt    sql.NullTime
		revokeReason sql.NullString
	)
	if session.RevokedAt() != nil {
		revokedAt = sql.NullTime{
			Time:  *session.RevokedAt(),
			Valid: true,
		}
	}
	if session.RevokeReason() != nil {
		revokeReason = sql.NullString{
			String: *session.RevokeReason(),
			Valid:  true,
		}
	}

	params := sqlc.RevokeRefreshSessionParams{
		ID:           session.ID(),
		RevokedAt:    revokedAt,
		RevokeReason: revokeReason,
	}

	return r.q.RevokeRefreshSession(ctx, params)
}

func (r *RefreshSessionRepository) RevokeAllForAccount(ctx context.Context, accountID uuid.UUID, reason *string) error {
	var revokeReason sql.NullString
	if reason != nil {
		revokeReason = sql.NullString{
			String: *reason,
			Valid:  true,
		}
	}
	params := sqlc.RevokeAllAccountRefreshSessionsParams{
		AccountID: accountID,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		RevokeReason: revokeReason,
	}
	return r.q.RevokeAllAccountRefreshSessions(ctx, params)
}

func (r *RefreshSessionRepository) RevokeDescendants(ctx context.Context, sessionID uuid.UUID, reason *string) error {
	var revokeReason sql.NullString
	if reason != nil {
		revokeReason = sql.NullString{
			String: *reason,
			Valid:  true,
		}
	}
	params := sqlc.RevokeRefreshSessionDescendantsParams{
		ID: sessionID,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		RevokeReason: revokeReason,
	}
	return r.q.RevokeRefreshSessionDescendants(ctx, params)
}

func (r *RefreshSessionRepository) DeleteExpired(ctx context.Context, expiresAt time.Time) error {
	return r.q.DeleteExpiredRefreshSessions(ctx, expiresAt)
}

func (r *RefreshSessionRepository) ListActiveForAccount(ctx context.Context, accountID uuid.UUID) ([]*model.RefreshSession, error) {
	params := sqlc.ListAccountActiveRefreshSessionsParams{
		AccountID: accountID,
		ExpiresAt: time.Now(),
	}
	rawList, err := r.q.ListAccountActiveRefreshSessions(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*model.RefreshSession, 0, len(rawList))
	for i := range rawList {
		mappedSession := mapper.MapSQLCToRefreshSession(rawList[i])
		result = append(result, mappedSession)
	}

	return result, nil
}
