package dto

import "github.com/google/uuid"

type DeleteAdRequest struct {
	adID uuid.UUID
}

type DeleteAdResponse struct {
	success bool
}
