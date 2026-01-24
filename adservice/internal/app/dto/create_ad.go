package dto

import "github.com/google/uuid"

type CreateAdRequest struct {
	SellerID    uuid.UUID
	Title       string
	Description *string
	Price       int64
	Images      []string
}

type CreateAdResponse struct {
	AdID uuid.UUID
}
