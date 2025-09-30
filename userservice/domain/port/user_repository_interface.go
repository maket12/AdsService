package port

import "AdsService/userservice/domain/entity"

type UserRepository interface {
	CheckUserExist(email string) (bool, error)
	AddUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(userID uint64) (*entity.User, error)
	GetUserRole(userID uint64) (string, error)
	EnhanceUser(userID uint64) error
	BanUser(userID uint64) error
	UnbanUser(userID uint64) error
}
