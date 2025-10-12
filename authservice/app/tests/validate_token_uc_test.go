package tests

import (
	"ads/authservice/app/dto"
	"ads/authservice/app/tests/data"
	"ads/authservice/app/tests/helpers"
	"ads/authservice/app/tests/mocks"
	"ads/authservice/app/uc_errors"
	"ads/authservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
