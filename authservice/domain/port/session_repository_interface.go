package port

import "AdsService/authservice/domain/entity"

type SessionRepository interface {
	InsertSession(s *entity.Session) error
	GetSessionByJTI(jti string) (*entity.Session, error)
	RevokeByJTI(jti string) error
	RevokeAllByUser(userID uint64) error
	CleanupExpired() error
	RotateSession(oldJTI string, newS *entity.Session) error
}
