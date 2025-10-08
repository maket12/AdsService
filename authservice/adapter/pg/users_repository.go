package pg

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) port.UserRepository {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) CheckUserExist(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email must be not empty")
	}

	var u entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UsersRepo) AddUser(ctx context.Context, user *entity.User) error {
	if user.ID == 0 {
		return errors.New("user ID must be valid")
	}

	exists, err := r.CheckUserExist(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("User with email %s already exists", user.Email)
		return nil
	}

	if err = r.db.WithContext(ctx).Create(user).Error; err != nil {
		fmt.Printf("Error while adding user: %v", err)
		return err
	}

	fmt.Printf("Successfully added user = %v", user.ID)
	return nil
}

func (r *UsersRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email must be not empty")
	}

	var user entity.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
