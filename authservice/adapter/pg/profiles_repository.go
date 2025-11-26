package pg

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type ProfilesRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProfilesRepo(db *gorm.DB, log *slog.Logger) port.ProfileRepository {
	return &ProfilesRepo{
		db:     db,
		logger: log,
	}
}

func (r *ProfilesRepo) AddProfile(ctx context.Context, userID uint64, name, phone string) (*entity.Profile, error) {
	if userID == 0 {
		return nil, errors.New("user ID must be valid")
	}
	if name == "" {
		return nil, errors.New("name must be not empty")
	}
	if phone == "" {
		return nil, errors.New("phone must be not empty")
	}

	p, err := entity.NewProfile(entity.UserID(userID), name, phone)
	if err != nil {
		return nil, err
	}

	r.logger.InfoContext(ctx, "Created new profile: %v", userID)
	return p, r.db.WithContext(ctx).Create(&p).Error
}
