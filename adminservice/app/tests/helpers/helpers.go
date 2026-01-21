package helpers

import (
	"ads/adminservice/domain/entity"
)

func MakeTestUser(userID uint64) *entity.User {
	return &entity.User{
		ID:   userID,
		Role: "user",
	}
}

func MakeTestProfile(userID uint64) *entity.Profile {
	return &entity.Profile{
		UserID: userID,
	}
}

func MakeTestProfiles() []entity.Profile {
	return []entity.Profile{
		{
			UserID: 1,
			Name:   "Alex",
		},
		{
			UserID: 2,
			Name:   "Jhon",
		},
		{
			UserID: 3,
			Name:   "Sara",
		},
	}
}
