package entity

import (
	"github.com/lib/pq"
	"time"
)

type Profile struct {
	UserID               uint64         `gorm:"primaryKey"`
	Name                 string         `gorm:"type:text"`
	Phone                string         `gorm:"type:text"`
	PhotoID              string         `gorm:"type:text"`
	NotificationsEnabled bool           `gorm:"default:true"`
	Subscriptions        pq.StringArray `gorm:"type:jsonb;serializer:json;not null;default:'[]'::jsonb"`
	UpdatedAt            time.Time      `gorm:"not null"`
	Banned               bool           `gorm:"default:false"`
}
