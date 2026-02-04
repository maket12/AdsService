package dto

import (
	"github.com/google/uuid"
)

type UpdateProfileInput struct {
	AccountID uuid.UUID
	FirstName *string
	LastName  *string
	Phone     *string
	AvatarURl *string
	Bio       *string
}

type UpdateProfileOutput struct {
	Success bool
}
