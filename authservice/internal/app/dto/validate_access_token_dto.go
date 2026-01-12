package dto

import "github.com/google/uuid"

type ValidateAccessToken struct {
	AccessToken string
}

type ValidateAccessTokenResponse struct {
	AccountID uuid.UUID
	Role      string
}
