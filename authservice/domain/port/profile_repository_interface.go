package port

import "AdsService/authservice/domain/entity"

type ProfileRepository interface {
	AddProfile(userID uint64, name, phone string) (*entity.Profile, error)
}
