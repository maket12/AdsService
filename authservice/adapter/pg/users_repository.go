package pg

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type UsersRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUsersRepo(db *gorm.DB, log *slog.Logger) port.UserRepository {
	return &UsersRepo{
		db:     db,
		logger: log,
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

	if u.IsBanned() {
		return false, entity.ErrUserBanned
	}

	return true, nil
}

func (r *UsersRepo) AddUser(ctx context.Context, user *entity.User) error {
	if entity.UserID(user.GetID()).Valid() {
		return errors.New("user ID must be valid")
	}

	exists, err := r.CheckUserExist(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		r.logger.InfoContext(ctx, "User with email %s already exists", user.Email)
		return nil
	}

	if user.GetRole() == "" {
		user.ChangeRole("user")
	}

	if err = r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.InfoContext(ctx, "Error while adding user: %v", err)
		return err
	}

	r.logger.InfoContext(ctx, "Successfully added user = %v", user.ID)
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

	if user.IsBanned() {
		return nil, entity.ErrUserBanned
	}

	return &user, nil
}
