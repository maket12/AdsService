package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetAdRequest struct {
	adID uuid.UUID
}

type GetAdResponse struct {
	adID        uuid.UUID
	sellerID    uuid.UUID
	title       string
	description *string
	price       int64
	status      string
	images      []string
	createdAt   time.Time
	updatedAt   time.Time
}
