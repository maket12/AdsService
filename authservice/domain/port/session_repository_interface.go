package port

import (
	"ads/authservice/domain/entity"
	"context"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSessionByJTI(ctx context.Context, jti string) (*entity.Session, error)
	RevokeByJTI(ctx context.Context, jti string) error
	RevokeAllByUser(ctx context.Context, userID uint64) error
	CleanupExpired(ctx context.Context) error
	RotateSession(ctx context.Context, oldJTI string, newSession *entity.Session) error
}
