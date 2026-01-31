package dto

import "github.com/google/uuid"

type PublishAdRequest struct {
	adID uuid.UUID
}

type PublishAdResponse struct {
	success bool
}
