package helpers

import (
	"AdsService/userservice/domain/entity"
)

func MakeTestProfile(
	userID uint64, name, phone string,
	notificationsEnabled bool, subscriptions []string, photoID string) *entity.Profile {
	return &entity.Profile{
		UserID:               userID,
		Name:                 name,
		Phone:                phone,
		PhotoID:              photoID,
		NotificationsEnabled: notificationsEnabled,
		Subscriptions:        subscriptions,
	}
}
