package dto

import "github.com/google/uuid"

type RejectAdRequest struct {
	AdID uuid.UUID
}

type RejectAdResponse struct {
	Success bool
}
