package dto

import "github.com/google/uuid"

type RejectAdRequest struct {
	adID uuid.UUID
}

type RejectAdResponse struct {
	success bool
}
