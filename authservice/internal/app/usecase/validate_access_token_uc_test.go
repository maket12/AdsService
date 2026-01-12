package usecase_test

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/usecase"
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/domain/port/mocks"
	"ads/authservice/internal/pkg/errs"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateAccessTokenUC_Execute(t *testing.T) {
	type adapter struct {
		account        *mocks.AccountRepository
		tokenGenerator *mocks.TokenGenerator
	}

	type testCase struct {
		name    string
		input   dto.ValidateAccessToken
		prepare func(a adapter)
		wantErr error
	}

	accountID := uuid.New()
	role := "user"
	accessToken := "valid-access-token"

	activeAcc, _ := model.NewAccount("test@test.com", "hash")

	bannedAcc, _ := model.NewAccount("banned@test.com", "hash")
	bannedAcc.Block()

	var tests = []testCase{
		{
			name: "Success",
			input: dto.ValidateAccessToken{
				AccessToken: accessToken,
			},
			prepare: func(a adapter) {
				a.tokenGenerator.On("ValidateAccessToken", mock.Anything, accessToken).
					Return(accountID, role, nil)
				a.account.On("GetByID", mock.Anything, accountID).
					Return(activeAcc, nil)
			},
			wantErr: nil,
		},
		{
			name: "Fail - Invalid Token",
			input: dto.ValidateAccessToken{
				AccessToken: "expired-or-fake-token",
			},
			prepare: func(a adapter) {
				a.tokenGenerator.On("ValidateAccessToken", mock.Anything, "expired-or-fake-token").
					Return(uuid.Nil, "", assert.AnError)
			},
			wantErr: uc_errors.ErrInvalidAccessToken,
		},
		{
			name: "Fail - Account Not Found In DB",
			input: dto.ValidateAccessToken{
				AccessToken: accessToken,
			},
			prepare: func(a adapter) {
				a.tokenGenerator.On("ValidateAccessToken", mock.Anything, accessToken).
					Return(accountID, role, nil)
				a.account.On("GetByID", mock.Anything, accountID).
					Return(nil, errs.ErrObjectNotFound)
			},
			wantErr: uc_errors.ErrInvalidAccessToken,
		},
		{
			name: "Fail - Account Is Banned",
			input: dto.ValidateAccessToken{
				AccessToken: accessToken,
			},
			prepare: func(a adapter) {
				a.tokenGenerator.On("ValidateAccessToken", mock.Anything, accessToken).
					Return(accountID, role, nil)
				a.account.On("GetByID", mock.Anything, accountID).
					Return(bannedAcc, nil)
			},
			wantErr: uc_errors.ErrCannotLogin,
		},
		{
			name: "Fail - Database Error",
			input: dto.ValidateAccessToken{
				AccessToken: accessToken,
			},
			prepare: func(a adapter) {
				a.tokenGenerator.On("ValidateAccessToken", mock.Anything, accessToken).
					Return(accountID, role, nil)
				a.account.On("GetByID", mock.Anything, accountID).
					Return(nil, assert.AnError)
			},
			wantErr: uc_errors.ErrGetAccountByIDDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := adapter{
				account:        mocks.NewAccountRepository(t),
				tokenGenerator: mocks.NewTokenGenerator(t),
			}

			tt.prepare(a)

			uc := usecase.NewValidateAccessTokenUC(a.account, a.tokenGenerator)

			res, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, res.Role)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, role, res.Role)
				assert.Equal(t, accountID, res.AccountID)
			}
		})
	}
}
