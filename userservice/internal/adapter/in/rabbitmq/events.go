package rabbitmq

import "github.com/google/uuid"

type AccountCreateEvent struct {
	AccountID uuid.UUID `json:"account_id"`
}
