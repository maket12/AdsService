package dto

import "github.com/google/uuid"

type DeleteAdRequest struct {
	AdID uuid.UUID
}

type DeleteAdResponse struct {
	Success bool
}
