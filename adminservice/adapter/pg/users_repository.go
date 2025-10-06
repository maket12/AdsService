package pg

import (
	"AdsService/adminservice/domain/entity"
	"AdsService/adminservice/domain/port"
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

type UsersRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUsersRepo(db *gorm.DB, logger *slog.Logger) port.UserRepository {
	return &UsersRepo{db: db, logger: logger}
}

func (r *UsersRepo) GetUserByID(userID uint64) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UsersRepo) GetUserRole(userID uint64) (string, error) {
	var role string
	result := r.db.Model(entity.User{}).Where("id = ?", userID).Pluck("role", role).Error
	if result != nil {
		r.logger.Error("database error: %v", result)
		return "", result
	}
	return role, result
}

func (r *UsersRepo) EnhanceUser(userID uint64) error {
	res := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("role", "admin").Error
	if res != nil {
		r.logger.Error("error to enhance user: %w", res)
		return res
	}
	return nil
}

func (r *UsersRepo) BanUser(userID uint64) error {
	res := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("banned", true).Error
	if res != nil {
		r.logger.Error("error to ban user: %w", res)
		return res
	}
	return nil
}

func (r *UsersRepo) UnbanUser(userID uint64) error {
	res := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("banned", false).Error
	if res != nil {
		r.logger.Error("error to unban user: %w", res)
		return res
	}
	return nil
}
