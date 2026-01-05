package tests

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/tests/data"
	"ads/authservice/internal/app/tests/helpers"
	"ads/authservice/internal/app/tests/mocks"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/usecase"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateToken_Success(t *testing.T) {
	for c, testCase := range data.ValidateTokenTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			tokens := new(mocks.MockTokensRepo)

			uc := &usecase.ValidateTokenUC{
				Tokens: tokens,
			}

			tokens.
				On("ParseAccessToken", testCase.AccessToken).
				Return(helpers.MakeAccessClaims(), nil)

			out, err := uc.Execute(dto.ValidateToken{
				AccessToken: testCase.AccessToken,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedAnswer, out.Valid)

			tokens.AssertExpectations(t)
		})
	}
}

func TestValidateToken_TokenError(t *testing.T) {
	for c, testCase := range data.ValidateTokenTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			tokens := new(mocks.MockTokensRepo)

			uc := &usecase.ValidateTokenUC{
				Tokens: tokens,
			}

			tokens.
				On("ParseAccessToken", testCase.AccessToken).
				Return(nil, errors.New("token error"))

			out, err := uc.Execute(dto.ValidateToken{
				AccessToken: testCase.AccessToken,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenIssue, err)
			assert.Equal(t, out.Valid, false)

			tokens.AssertExpectations(t)
		})
	}
}
