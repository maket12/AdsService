package dto

import "github.com/google/uuid"

type PublishAdRequest struct {
	AdID uuid.UUID
}

type PublishAdResponse struct {
	Success bool
}
