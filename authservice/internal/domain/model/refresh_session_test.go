package model_test

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewRefreshSession(t *testing.T) {
	type testCase struct {
		name        string
		accountID   uuid.UUID
		tokenHash   string
		rotatedFrom *uuid.UUID
	}
}
