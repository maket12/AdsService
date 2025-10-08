package entity

import (
	"errors"
	"time"
)

type Profile struct {
	UserID               uint64
	Name                 string
	Phone                string
	NotificationsEnabled bool
	Subscriptions        []string
	UpdatedAt            time.Time
}

func (p *Profile) Validate() error {
	if p.UserID == 0 {
		return errors.New("user ID is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	if len(p.Name) > 100 {
		return errors.New("name is too long")
	}
	if p.Phone == "" {
		return errors.New("phone is required")
	}
	return nil
}

func (p *Profile) EnableNotifications() {
	p.NotificationsEnabled = true
	p.UpdatedAt = time.Now()
}

func (p *Profile) AddSubscription(subscription string) {
	for _, sub := range p.Subscriptions {
		if sub == subscription {
			return
		}
	}
	p.Subscriptions = append(p.Subscriptions, subscription)
	p.UpdatedAt = time.Now()
}
