package ports

import (
	"ads/internal/core/domain/model/profile"
	"context"
)

type ProfileRepository interface {
	AddProfile(ctx context.Context, aggregate *profile.Profile) (*profile.Profile, error)
}
