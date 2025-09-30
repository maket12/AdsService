package entity

import "time"

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"type:text;default:user;not null"`
	Banned    bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
