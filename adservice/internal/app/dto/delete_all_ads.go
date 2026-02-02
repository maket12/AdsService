package dto

import "github.com/google/uuid"

type DeleteAllAdsRequest struct {
	SellerID uuid.UUID
}

type DeleteAllAdsResponse struct {
	Success bool
}
