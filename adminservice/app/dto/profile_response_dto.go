package dto

import "time"

type ProfileResponse struct {
	UserID               uint64
	Name                 string
	Phone                string
	PhotoID              string
	NotificationsEnabled bool
	Subscriptions        []string
	UpdatedAt            time.Time
}
