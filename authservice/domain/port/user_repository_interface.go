package port

import "AdsService/authservice/domain/entity"

type UserRepository interface {
	CheckUserExist(email string) (bool, error)
	AddUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
}
