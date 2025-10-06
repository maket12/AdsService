package pg

import (
	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

type UsersRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUsersRepo(db *gorm.DB, logger *slog.Logger) port.UserRepository {
	return &UsersRepo{
		db:     db,
		logger: logger,
	}
}

func (r *UsersRepo) CheckUserExist(email string) (bool, error) {
	var u entity.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UsersRepo) AddUser(user *entity.User) error {
	exists, err := r.CheckUserExist(user.Email)
	if err != nil {
		return err
	}
	if exists {
		r.logger.Error("User with email %s already exists", user.Email)
		return nil
	}

	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error("Error while adding user: %v", err)
		return err
	}

	r.logger.Info("Successfully added user = %v", user.ID)
	return nil
}

func (r *UsersRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
