package database

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"type:text;default:user;not null"`
	Banned    bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CheckUserExist(email string) (bool, error) {
	var u User
	err := DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func AddUser(user *User) error {
	exists, err := CheckUserExist(user.Email)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("User with email %s already exists", user.Email)
		return nil
	}

	if err := DB.Create(user).Error; err != nil {
		log.Printf("Error while adding user: %v", err)
		return err
	}

	log.Printf("Successfully added user = %v", user.ID)
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByID(userID uint64) (*User, error) {
	var user User
	result := DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserRole(userID uint64) (string, error) {
	var role string
	result := DB.Model(User{}).Where("id = ?", userID).Pluck("role", role).Error
	if result != nil {
		log.Printf("database error: %v", result)
		return "", result
	}
	return role, result
}

func EnhanceUser(userID uint64) error {
	return DB.Model(&User{}).Where("id = ?", userID).Update("role", "admin").Error
}

func BanUser(userID uint64) error {
	return DB.Model(&User{}).Where("id = ?", userID).Update("banned", true).Error
}

func UnbanUser(userID uint64) error {
	return DB.Model(&User{}).Where("id = ?", userID).Update("banned", false).Error
}
