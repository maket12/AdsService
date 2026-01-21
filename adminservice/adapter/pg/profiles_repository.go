package pg

import (
	"ads/adminservice/domain/entity"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type ProfilesRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProfilesRepo(db *gorm.DB, log *slog.Logger) *ProfilesRepo {
	return &ProfilesRepo{
		db:     db,
		logger: log,
	}
}

func (r *ProfilesRepo) GetProfile(ctx context.Context, userID uint64) (*entity.Profile, error) {
	var p entity.Profile
	result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).First(&p)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &p, nil
}

func (r *ProfilesRepo) GetAllProfiles(ctx context.Context, limit, offset uint32) ([]entity.Profile, error) {
	var profiles []entity.Profile
	result := r.db.WithContext(ctx).Limit(int(limit)).Offset(int(offset)).Find(&profiles)
	if result.Error != nil {
		return nil, result.Error
	}
	return profiles, nil
}
