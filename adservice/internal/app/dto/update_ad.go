package dto

import "github.com/google/uuid"

type UpdateAdRequest struct {
	AdID        uuid.UUID
	Title       *string
	Description *string
	Price       *int64
	Images      []string
}

type UpdateAdResponse struct {
	Success bool
}
