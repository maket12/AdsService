package pg

import (
	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type ProfilesRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProfilesRepo(db *gorm.DB, logger *slog.Logger) port.ProfileRepository {
	return &ProfilesRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ProfilesRepo) AddProfile(userID uint64, name, phone string) (*entity.Profile, error) {
	var p = entity.Profile{
		UserID:    userID,
		Name:      name,
		Phone:     phone,
		UpdatedAt: time.Now().UTC(),
	}
	r.logger.Info("Created new profile: %v", userID)
	return &p, r.db.Create(&p).Error
}
