package dto

import "github.com/google/uuid"

type AssignRole struct {
	AccountID uuid.UUID
	Role      string
}

type AssignRoleResponse struct {
	Assign bool
}
