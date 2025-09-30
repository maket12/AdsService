package port

import "AdsService/userservice/domain/entity"

type ProfileRepository interface {
	AddProfile(userID uint64, name, phone string) (*entity.Profile, error)
	UpdateProfileName(userID uint64, name string) (*entity.Profile, error)
	UpdateProfilePhone(userID uint64, phone string) (*entity.Profile, error)
	UpdateProfilePhoto(userID uint64, photoID string) (*entity.Profile, error)
	EnableNotifications(userID uint64) (*entity.Profile, error)
	DisableNotifications(userID uint64) (*entity.Profile, error)
	UpdateProfileSubscriptions(userID uint64, subscriptions []string) (*entity.Profile, error)
	UpdateProfileTime(userID uint64) (*entity.Profile, error)
	GetProfile(userID uint64) (*entity.Profile, error)
	GetAllProfiles() ([]entity.Profile, error)
}
