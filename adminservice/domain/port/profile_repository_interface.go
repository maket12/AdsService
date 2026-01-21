package port

import (
	"ads/adminservice/domain/entity"
	"context"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, userID uint64) (*entity.Profile, error)
	GetAllProfiles(ctx context.Context, limit, offset uint32) ([]entity.Profile, error)
}
