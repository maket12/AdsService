package dto

import "github.com/google/uuid"

type CreateAdRequest struct {
	sellerID    uuid.UUID
	title       string
	description *string
	price       int64
	images      []string
}

type CreateAdResponse struct {
	adID uuid.UUID
}
