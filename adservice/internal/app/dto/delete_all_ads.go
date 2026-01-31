package dto

import "github.com/google/uuid"

type DeleteAllAdsRequest struct {
	sellerID uuid.UUID
}

type DeleteAllAdsResponse struct {
	success bool
}
