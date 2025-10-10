package pg

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ProfilesRepo struct {
	db *gorm.DB
}

func NewProfilesRepo(db *gorm.DB) port.ProfileRepository {
	return &ProfilesRepo{
		db: db,
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

	fmt.Printf("Created new profile: %v", userID)
	return p, r.db.WithContext(ctx).Create(&p).Error
}
