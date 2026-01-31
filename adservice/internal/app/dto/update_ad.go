package dto

import "github.com/google/uuid"

type UpdateAdRequest struct {
	adID        uuid.UUID
	title       *string
	description *string
	price       *int64
	images      []string
}

type UpdateAdResponse struct {
	success bool
}
