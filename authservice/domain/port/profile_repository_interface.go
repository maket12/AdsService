package port

import (
	"ads/authservice/domain/entity"
	"context"
)

type ProfileRepository interface {
	AddProfile(ctx context.Context, userID uint64, name, phone string) (*entity.Profile, error)
}
