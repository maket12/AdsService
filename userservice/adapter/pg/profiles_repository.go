package pg

import (
	"AdsService/userservice/domain/entity"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log/slog"
	"time"
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

func (r *ProfilesRepo) UpdateProfileName(userID uint64, name string) (*entity.Profile, error) {
	if result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("name", name).Error; result != nil {
		return nil, result
	}
	r.logger.Info("Updated name = %v of user = %v", name, userID)
	return r.UpdateProfileTime(userID)
}

func (r *ProfilesRepo) UpdateProfilePhone(userID uint64, phone string) (*entity.Profile, error) {
	if result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("phone", phone).Error; result != nil {
		return nil, result
	}
	r.logger.Info("Updated phone = %v of user = %v", phone, userID)
	return r.UpdateProfileTime(userID)
}

func (r *ProfilesRepo) UpdateProfilePhoto(userID uint64, photoID string) (*entity.Profile, error) {
	if result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("photo_id", photoID).Error; result != nil {
		return nil, result
	}
	r.logger.Info("Updated photo = %v of user = %v", photoID, userID)
	return r.UpdateProfileTime(userID)
}

func (r *ProfilesRepo) EnableNotifications(userID uint64) (*entity.Profile, error) {
	if result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("notifications_enabled", true).Error; result != nil {
		return nil, result
	}
	r.logger.Info("Enable notifications for user = %v", userID)
	return r.UpdateProfileTime(userID)
}

func (r *ProfilesRepo) DisableNotifications(userID uint64) (*entity.Profile, error) {
	if result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("notifications_enabled", false).Error; result != nil {
		return nil, result
	}
	r.logger.Info("Disable notifications for user = %v", userID)
	return r.UpdateProfileTime(userID)
}

func (r *ProfilesRepo) UpdateProfileSubscriptions(userID uint64, subscriptions []string) (*entity.Profile, error) {
	if result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("subscriptions", pq.Array(subscriptions)).Error; result != nil {
		return nil, result
	}
	r.logger.Info("Updated subscriptions = %v of user = %v", subscriptions, userID)
	return r.UpdateProfileTime(userID)
}

func (r *ProfilesRepo) UpdateProfileTime(userID uint64) (*entity.Profile, error) {
	var p entity.Profile
	result := r.db.Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("updated_at", gorm.Expr("now()")).First(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
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
