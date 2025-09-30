package entity

import "time"

type Session struct {
	ID          uint64
	UserID      uint64
	JTI         string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	RotatedFrom *string
	UserAgent   *string
	IP          *string
	RevokedAt   time.Time
	ReusedAt    time.Time
}
