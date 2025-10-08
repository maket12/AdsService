package pg

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) port.SessionRepository {
	return &SessionsRepo{
		db: db,
	}
}

func (r *SessionsRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	if session.UserID == 0 {
		return errors.New("user ID is required")
	}
	if session.ExpiresAt.Before(time.Now()) {
		return errors.New("expiration date must be in future")
	}

	result := r.db.WithContext(ctx).Create(session).Error
	if result != nil {
		return fmt.Errorf("error while insert session: %w", result)
	}
	fmt.Printf("Inserted session for user[id=%v]", session.UserID)
	return nil
}

func (r *SessionsRepo) GetSessionByJTI(ctx context.Context, jti string) (*entity.Session, error) {
	if jti == "" {
		return nil, errors.New("empty jti")
	}

	var s entity.Session
	if err := r.db.WithContext(ctx).Where("jti = ?", jti).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionsRepo) RevokeByJTI(ctx context.Context, jti string) error {
	if jti == "" {
		return errors.New("empty jti")
	}

	return r.db.WithContext(ctx).Model(&entity.Session{}).Where("jti = ? AND revoked_at IS NULL", jti).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func (r *SessionsRepo) RevokeAllByUser(ctx context.Context, userID uint64) error {
	if userID == 0 {
		return errors.New("user ID is required")
	}

	return r.db.WithContext(ctx).Model(&entity.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func (r *SessionsRepo) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < now()").Delete(&entity.Session{}).Error
}

func (r *SessionsRepo) RotateSession(ctx context.Context, oldJTI string, newSession *entity.Session) error {
	if oldJTI == "" {
		return errors.New("empty old jti")
	}
	if newSession.UserID == 0 {
		return errors.New("user ID is required")
	}
	if newSession.ExpiresAt.Before(time.Now()) {
		return errors.New("expiration date must be in future")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Session{}).
			Where("jti = ? AND revoked_at IS NULL", oldJTI).
			Update("revoked_at", gorm.Expr("now()")).Error; err != nil {
			return err
		}
		if newSession.IssuedAt.IsZero() {
			newSession.IssuedAt = time.Now().UTC()
		}
		return tx.Create(newSession).Error
	})
}
