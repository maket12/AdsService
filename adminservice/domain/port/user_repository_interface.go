package port

import "AdsService/adminservice/domain/entity"

type UserRepository interface {
	GetUserByID(userID uint64) (*entity.User, error)
	GetUserRole(userID uint64) (string, error)
	EnhanceUser(userID uint64) error
	BanUser(userID uint64) error
	UnbanUser(userID uint64) error
}
