package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetProfile struct {
	AccountID uuid.UUID
}

type GetProfileOutput struct {
	AccountID uuid.UUID
	FirstName *string
	LastName  *string
	Phone     *string
	AvatarURl *string
	Bio       *string
	UpdatedAt time.Time
}
