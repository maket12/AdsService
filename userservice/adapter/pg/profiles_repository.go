package pg

import (
	"ads/userservice/domain/entity"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/lib/pq"
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

func (r *ProfilesRepo) AddProfile(ctx context.Context, userID uint64, name, phone string) (*entity.Profile, error) {
	p, err := entity.NewProfile(entity.UserID(userID), name, phone)
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Create(&p).Error
	if err != nil {
		return nil, err
	}

	r.logger.InfoContext(ctx, "Created new profile: %v", userID)

	return p, nil
}

func (r *ProfilesRepo) UpdateProfileName(ctx context.Context, userID uint64, name string) (*entity.Profile, error) {
	if result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("name", name).Error; result != nil {
		return nil, result
	}
	r.logger.InfoContext(ctx, "Updated name = %v of user = %v", name, userID)
	return r.UpdateProfileTime(ctx, userID)
}

func (r *ProfilesRepo) UpdateProfilePhone(ctx context.Context, userID uint64, phone string) (*entity.Profile, error) {
	if result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("phone", phone).Error; result != nil {
		return nil, result
	}
	r.logger.InfoContext(ctx, "Updated phone = %v of user = %v", phone, userID)
	return r.UpdateProfileTime(ctx, userID)
}

func (r *ProfilesRepo) UpdateProfilePhoto(ctx context.Context, userID uint64, photoID string) (*entity.Profile, error) {
	if result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("photo_id", photoID).Error; result != nil {
		return nil, result
	}
	r.logger.InfoContext(ctx, "Updated photo = %v of user = %v", photoID, userID)
	return r.UpdateProfileTime(ctx, userID)
}

func (r *ProfilesRepo) EnableNotifications(ctx context.Context, userID uint64) (*entity.Profile, error) {
	if result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("notifications_enabled", true).Error; result != nil {
		return nil, result
	}
	r.logger.InfoContext(ctx, "Enable notifications for user = %v", userID)
	return r.UpdateProfileTime(ctx, userID)
}

func (r *ProfilesRepo) DisableNotifications(ctx context.Context, userID uint64) (*entity.Profile, error) {
	if result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("notifications_enabled", false).Error; result != nil {
		return nil, result
	}
	r.logger.InfoContext(ctx, "Disable notifications for user = %v", userID)
	return r.UpdateProfileTime(ctx, userID)
}

func (r *ProfilesRepo) UpdateProfileSubscriptions(ctx context.Context, userID uint64, subscriptions []string) (*entity.Profile, error) {
	if result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("subscriptions", pq.Array(subscriptions)).Error; result != nil {
		return nil, result
	}
	r.logger.InfoContext(ctx, "Updated subscriptions = %v of user = %v", subscriptions, userID)
	return r.UpdateProfileTime(ctx, userID)
}

func (r *ProfilesRepo) UpdateProfileTime(ctx context.Context, userID uint64) (*entity.Profile, error) {
	var p entity.Profile
	result := r.db.WithContext(ctx).Model(entity.Profile{}).Where("user_id = ?", userID).
		Update("updated_at", gorm.Expr("now()")).First(&p)
	if result.Error != nil {
		return nil, result.Error
	}

	r.logger.InfoContext(ctx, "Successfully updated profile at %v", time.Now().UTC())

	return &p, nil
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
