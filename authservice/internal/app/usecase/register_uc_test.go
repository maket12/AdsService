package usecase_test

import (
	"context"
	"errors"
	"testing"

	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/usecase"
	"ads/authservice/internal/domain/port/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUC_Execute(t *testing.T) {
	type fields struct {
		account        *mocks.AccountRepository
		accountRole    *mocks.AccountRoleRepository
		passwordHasher *mocks.PasswordHasher
	}

	tests := []struct {
		name    string
		input   dto.Register
		prepare func(f fields)
		wantErr error
	}{
		{
			name: "Success",
			input: dto.Register{
				Email:    "test@example.com",
				Password: "securePassword123",
			},
			prepare: func(f fields) {
				f.passwordHasher.On("Hash", "securePassword123").
					Return("hashed_password", nil)

				f.account.On("Create", mock.Anything, mock.MatchedBy(func(a interface{}) bool {
					acc := a.(interface{ Email() string })
					return acc.Email() == "test@example.com"
				})).Return(nil)

				f.accountRole.On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:  "Error - hashing password",
			input: dto.Register{Email: "test@example.com", Password: "123"},
			prepare: func(f fields) {
				f.passwordHasher.On("Hash", "123").
					Return("", errors.New("salt error"))
			},
			wantErr: uc_errors.ErrHashPassword,
		},
		{
			name:  "Error - create account",
			input: dto.Register{Email: "exists@example.com", Password: "password"},
			prepare: func(f fields) {
				f.passwordHasher.On("Hash", "password").
					Return("hashed", nil)
				f.account.On("Create", mock.Anything, mock.Anything).
					Return(errors.New("db error"))
			},
			wantErr: uc_errors.ErrCreateAccountDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				account:        mocks.NewAccountRepository(t),
				accountRole:    mocks.NewAccountRoleRepository(t),
				passwordHasher: mocks.NewPasswordHasher(t),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			uc := usecase.NewRegisterUC(f.account, f.accountRole, f.passwordHasher)

			res, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res.AccountID)
			}
		})
	}
}
