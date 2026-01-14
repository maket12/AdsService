package model

import "github.com/google/uuid"

type AccountCreatedEvent struct {
	AccountID uuid.UUID
}

func NewAccountCreatedEvent(accountID uuid.UUID) *AccountCreatedEvent {
	return &AccountCreatedEvent{AccountID: accountID}
}
