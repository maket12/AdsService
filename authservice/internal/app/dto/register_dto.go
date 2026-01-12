package dto

import "github.com/google/uuid"

type Register struct {
	Email    string
	Password string
}

type RegisterResponse struct {
	AccountID uuid.UUID
}
