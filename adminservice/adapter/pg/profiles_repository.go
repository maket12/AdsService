package pg

import (
	"AdsService/adminservice/domain/entity"
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

type ProfilesRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProfilesRepo(db *gorm.DB, logger *slog.Logger) *ProfilesRepo {
	return &ProfilesRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ProfilesRepo) GetProfile(userID uint64) (*entity.Profile, error) {
	var p entity.Profile
	result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).First(&p)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &p, nil
}

func (r *ProfilesRepo) GetAllProfiles(limit, offset uint32) ([]entity.Profile, error) {
	var profiles []entity.Profile
	result := r.db.Limit(int(limit)).Offset(int(offset)).Find(&profiles)
	if result.Error != nil {
		return nil, result.Error
	}
	return profiles, nil
}
