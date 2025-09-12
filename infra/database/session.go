package database

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID          uint64
	UserID      uint64
	JTI         string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	RotatedFrom *string
	UserAgent   *string
	IP          *string
	RevokedAt   time.Time
	ReusedAt    time.Time
}

func InsertSession(s *Session) error {
	return DB.Create(s).Error
}

func GetSessionByJTI(jti string) (*Session, error) {
	var s Session
	if err := DB.Where("jti = ?", jti).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func RevokeByJTI(jti string) error {
	return DB.Model(&Session{}).Where("jti = ? AND revoked_at IS NULL", jti).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func RevokeAllByUser(userID uint64) error {
	return DB.Model(&Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func CleanupExpired() error {
	return DB.Where("expires_at < now()").Delete(&Session{}).Error
}

func RotateSession(oldJTI string, newS *Session) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Session{}).
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
