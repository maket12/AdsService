package profile_dto

import "time"

type ProfileResponseDTO struct {
	UserID               uint64
	Name                 string
	Phone                string
	PhotoID              string
	NotificationsEnabled bool
	Subscriptions        []string
	UpdatedAt            time.Time
}
