package pg

import (
	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
	"gorm.io/gorm"
	"time"
)

type SessionsRepo struct{ db *gorm.DB }

func NewSessionsRepo(db *gorm.DB) port.SessionRepository {
	return &SessionsRepo{db: db}
}

func (r *SessionsRepo) InsertSession(s *entity.Session) error {
	return r.db.Create(s).Error
}

func (r *SessionsRepo) GetSessionByJTI(jti string) (*entity.Session, error) {
	var s entity.Session
	if err := r.db.Where("jti = ?", jti).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionsRepo) RevokeByJTI(jti string) error {
	return r.db.Model(&entity.Session{}).Where("jti = ? AND revoked_at IS NULL", jti).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func (r *SessionsRepo) RevokeAllByUser(userID uint64) error {
	return r.db.Model(&entity.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func (r *SessionsRepo) CleanupExpired() error {
	return r.db.Where("expires_at < now()").Delete(&entity.Session{}).Error
}

func (r *SessionsRepo) RotateSession(oldJTI string, newS *entity.Session) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Session{}).
			Where("jti = ? AND revoked_at IS NULL", oldJTI).
			Update("revoked_at", gorm.Expr("now()")).Error; err != nil {
			return err
		}
		if newS.IssuedAt.IsZero() {
			newS.IssuedAt = time.Now().UTC()
		}
		return tx.Create(newS).Error
	})
}
