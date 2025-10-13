package pg

import (
	"ads/adminservice/domain/entity"
	"ads/adminservice/domain/port"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type UsersRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUsersRepo(db *gorm.DB, log *slog.Logger) port.UserRepository {
	return &UsersRepo{db: db, logger: log}
}

func (r *UsersRepo) GetUserByID(ctx context.Context, userID uint64) (*entity.User, error) {
	var user entity.User
	result := r.db.WithContext(ctx).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UsersRepo) GetUserRole(ctx context.Context, userID uint64) (string, error) {
	var role string
	result := r.db.WithContext(ctx).Model(entity.User{}).Where("id = ?", userID).Pluck("role", role).Error
	if result != nil {
		return "", fmt.Errorf("database error: %v", result)
	}
	return role, result
}

func (r *UsersRepo) EnhanceUser(ctx context.Context, userID uint64) error {
	res := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("role", "admin").Error
	if res != nil {
		return fmt.Errorf("error to enhance user: %w", res)
	}

	r.logger.InfoContext(ctx, "Successfully enhanced user[%v].", userID)

	return nil
}

func (r *UsersRepo) BanUser(ctx context.Context, userID uint64) error {
	res := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("banned", true).Error
	if res != nil {
		return fmt.Errorf("error to ban user: %w", res)
	}
	return nil
}

func (r *UsersRepo) UnbanUser(ctx context.Context, userID uint64) error {
	res := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("banned", false).Error
	if res != nil {
		return fmt.Errorf("error to unban user: %w", res)
	}
	return nil
}
