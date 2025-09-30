package pg

import (
	"AdsService/userservice/domain/entity"
	"AdsService/userservice/domain/port"
	"errors"
	"gorm.io/gorm"
	"log"
)

type UsersRepo struct{ db *gorm.DB }

func NewUsersRepo(db *gorm.DB) port.UserRepository { return &UsersRepo{db: db} }

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
		log.Printf("User with email %s already exists", user.Email)
		return nil
	}

	if err := r.db.Create(user).Error; err != nil {
		log.Printf("Error while adding user: %v", err)
		return err
	}

	log.Printf("Successfully added user = %v", user.ID)
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
		log.Printf("database error: %v", result)
		return "", result
	}
	return role, result
}

func (r *UsersRepo) EnhanceUser(userID uint64) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("role", "admin").Error
}

func (r *UsersRepo) BanUser(userID uint64) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("banned", true).Error
}

func (r *UsersRepo) UnbanUser(userID uint64) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("banned", false).Error
}
