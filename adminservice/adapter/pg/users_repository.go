package pg

import (
	"ads/adminservice/domain/entity"
	"ads/adminservice/domain/port"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) port.UserRepository {
	return &UsersRepo{db: db}
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
		return "", fmt.Errorf("database error: %v", result)
	}
	return role, result
}

func (r *UsersRepo) EnhanceUser(userID uint64) error {
	res := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("role", "admin").Error
	if res != nil {
		return fmt.Errorf("error to enhance user: %w", res)
	}
	return nil
}

func (r *UsersRepo) BanUser(userID uint64) error {
	res := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("banned", true).Error
	if res != nil {
		return fmt.Errorf("error to ban user: %w", res)
	}
	return nil
}

func (r *UsersRepo) UnbanUser(userID uint64) error {
	res := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("banned", false).Error
	if res != nil {
		return fmt.Errorf("error to unban user: %w", res)
	}
	return nil
}
