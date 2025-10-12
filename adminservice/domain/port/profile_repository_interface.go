package port

import "ads/adminservice/domain/entity"

type ProfileRepository interface {
	GetProfile(userID uint64) (*entity.Profile, error)
	GetAllProfiles(limit, offset uint32) ([]entity.Profile, error)
}
