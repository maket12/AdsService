package port

import (
	"ads/userservice/domain/entity"
	"context"
)

type ProfileRepository interface {
	AddProfile(ctx context.Context, userID uint64, name, phone string) (*entity.Profile, error)
	UpdateProfileName(ctx context.Context, userID uint64, name string) (*entity.Profile, error)
	UpdateProfilePhone(ctx context.Context, userID uint64, phone string) (*entity.Profile, error)
	UpdateProfilePhoto(ctx context.Context, userID uint64, photoID string) (*entity.Profile, error)
	EnableNotifications(ctx context.Context, userID uint64) (*entity.Profile, error)
	DisableNotifications(ctx context.Context, userID uint64) (*entity.Profile, error)
	UpdateProfileSubscriptions(ctx context.Context, userID uint64, subscriptions []string) (*entity.Profile, error)
	UpdateProfileTime(ctx context.Context, userID uint64) (*entity.Profile, error)
	GetProfile(ctx context.Context, userID uint64) (*entity.Profile, error)
}
