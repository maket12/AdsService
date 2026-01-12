package model_test

import (
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccountRole(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		accountID uuid.UUID
		expect    error
	}

	var tests = []testCase{
		{
			name:      "success",
			accountID: uuid.New(),
			expect:    nil,
		},
		{
			name:      "nullable account id",
			accountID: uuid.Nil,
			expect:    errs.ErrValueIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accRole, err := model.NewAccountRole(tt.accountID)
			if tt.expect == nil {
				require.NoError(t, err)
				require.NotNil(t, accRole)
				assert.Equal(t, tt.accountID, accRole.AccountID())
				assert.Equal(t, model.RoleUser, accRole.Role())
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrValueIsInvalid)
				assert.Nil(t, accRole)
			}
		})
	}
}

func TestAccountRole_Assign(t *testing.T) {
	t.Parallel()
	var accRole = model.RestoreAccountRole(uuid.New(), model.RoleUser)
	accRole.Assign()
	assert.Equal(t, accRole.Role(), model.RoleAdmin)
}
