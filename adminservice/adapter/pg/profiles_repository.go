package pg

import (
	"ads/adminservice/domain/entity"
	"errors"
	"gorm.io/gorm"
)

type ProfilesRepo struct {
	db *gorm.DB
}

func NewProfilesRepo(db *gorm.DB) *ProfilesRepo {
	return &ProfilesRepo{
		db: db,
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
