package pg

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) port.SessionRepository {
	return &SessionsRepo{
		db: db,
	}
}

func (r *SessionsRepo) InsertSession(ctx context.Context, session *entity.Session) error {
	if session.UserID().Valid() {
		return errors.New("user ID is required")
	}
	if session.ExpiresAt().Before(time.Now()) {
		return errors.New("expiration date must be in future")
	}

	if session.IsExpired() {
		return entity.ErrSessionExpired
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

	session, err := r.GetSessionByJTI(ctx, jti)
	if err != nil {
		return err
	}

	if err := session.Revoke(); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(session).Error
}

func (r *SessionsRepo) RevokeAllByUser(ctx context.Context, userID uint64) error {
	if userID == 0 {
		return errors.New("user ID is required")
	}

	var sessions []entity.Session
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Find(&sessions).Error; err != nil {
		return err
	}

	for _, session := range sessions {
		if err := session.Revoke(); err != nil {
			return err
		}

		if err := r.db.WithContext(ctx).Save(&session).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SessionsRepo) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < now()").Delete(&entity.Session{}).Error
}

func (r *SessionsRepo) RotateSession(ctx context.Context, oldJTI string, newSession *entity.Session) error {
	if oldJTI == "" {
		return errors.New("empty old jti")
	}
	if newSession.UserID() == 0 {
		return errors.New("user ID is required")
	}
	if newSession.ExpiresAt().Before(time.Now()) {
		return errors.New("expiration date must be in future")
	}

	oldSession, err := r.GetSessionByJTI(ctx, oldJTI)
	if err != nil {
		return err
	}

	if err := oldSession.Rotate(string(newSession.JTI())); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(oldSession).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(newSession).Error; err != nil {
		return err
	}

	return nil
}
