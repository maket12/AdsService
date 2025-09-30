package pg

import (
	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
	"gorm.io/gorm"
	"log"
	"time"
)

type ProfilesRepo struct{ db *gorm.DB }

func NewProfilesRepo(db *gorm.DB) port.ProfileRepository {
	return &ProfilesRepo{db: db}
}

func (r *ProfilesRepo) AddProfile(userID uint64, name, phone string) (*entity.Profile, error) {
	var p = entity.Profile{
		UserID:    userID,
		Name:      name,
		Phone:     phone,
		UpdatedAt: time.Now().UTC(),
	}
	log.Printf("Created new profile: %v", userID)
	return &p, r.db.Create(&p).Error
}
