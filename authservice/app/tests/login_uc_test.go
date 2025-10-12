package tests

import (
	"ads/authservice/app/dto"
	"ads/authservice/app/tests/data"
	"ads/authservice/app/tests/helpers"
	"ads/authservice/app/tests/mocks"
	"ads/authservice/app/uc_errors"
	"ads/authservice/app/usecase"
	"ads/authservice/domain/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestLoginUC_Success(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			tokens := new(mocks.MockTokensRepo)
			sessions := new(mocks.MockSessionsRepo)

			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: sessions,
				Tokens:   tokens,
			}

			existingUser := helpers.MakeTestUser(testCase.Email, testCase.Password, testCase.Role)

			users.
				On("GetUserByEmail", testCase.Email).
				Return(existingUser, nil)

			tokens.
				On("GenerateAccessToken", existingUser.ID, testCase.Email, testCase.Role).
				Return(testCase.ExpectedAccessToken, nil)

			tokens.
				On("GenerateRefreshToken", existingUser.ID).
				Return(testCase.ExpectedRefreshToken, nil)

			tokens.
				On("ParseRefreshToken", testCase.ExpectedRefreshToken).
				Return(helpers.MakeRefreshClaims("jti-x", time.Now().Add(24*time.Hour), time.Now()), nil)

			sessions.
				On("CreateSession", mock.MatchedBy(func(s *entity.Session) bool {
					return s.UserID == existingUser.ID && s.JTI == "jti-x" && !s.IssuedAt.IsZero() && s.ExpiresAt.After(s.IssuedAt)
				})).
				Return(nil)

			out, err := uc.Execute(dto.Login{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedAccessToken, out.AccessToken)
			assert.Equal(t, testCase.ExpectedRefreshToken, out.RefreshToken)

			users.AssertExpectations(t)
			tokens.AssertExpectations(t)
			sessions.AssertExpectations(t)
		})
	}
}

func TestLoginUC_GetUserError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   new(mocks.MockTokensRepo),
			}

			users.
				On("GetUserByEmail", testCase.Email).
				Return(&entity.User{}, errors.New("db error"))

			res, err := uc.Execute(dto.Login{Email: testCase.Email, Password: testCase.Password})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUser, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestLoginUC_UserNotFoundError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   new(mocks.MockTokensRepo),
			}

			users.
				On("GetUserByEmail", testCase.Email).
				Return(nil, nil)

			res, err := uc.Execute(dto.Login{Email: testCase.Email, Password: testCase.Password})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrUserNotFound, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestLoginUC_InvalidPasswordError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   new(mocks.MockTokensRepo),
			}

			existingUser := helpers.MakeTestUser(testCase.Email, testCase.Password, testCase.Role)

			users.
				On("GetUserByEmail", testCase.Email).
				Return(existingUser, nil)

			res, err := uc.Execute(dto.Login{
				Email:    testCase.Email,
				Password: "wrong_pass",
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrInvalidCredential, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestLoginUC_AccessTokenError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			tokens := new(mocks.MockTokensRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   tokens,
			}

			existingUser := helpers.MakeTestUser(testCase.Email, testCase.Password, testCase.Role)

			users.
				On("GetUserByEmail", mock.Anything).
				Return(existingUser, nil)

			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("", errors.New("token error"))

			res, err := uc.Execute(dto.Login{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenIssue, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestLoginUC_RefreshTokenError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			tokens := new(mocks.MockTokensRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   tokens,
			}

			existingUser := helpers.MakeTestUser(testCase.Email, testCase.Password, testCase.Role)

			users.
				On("GetUserByEmail", mock.Anything).
				Return(existingUser, nil)
			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("access-token", nil)
			tokens.
				On("GenerateRefreshToken", mock.Anything).
				Return("", errors.New("token error"))

			res, err := uc.Execute(dto.Login{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenIssue, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestLoginUC_ParseTokenError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			tokens := new(mocks.MockTokensRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   tokens,
			}

			existingUser := helpers.MakeTestUser(testCase.Email, testCase.Password, testCase.Role)

			users.
				On("GetUserByEmail", mock.Anything).
				Return(existingUser, nil)
			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("access-token", nil)
			tokens.
				On("GenerateRefreshToken", mock.Anything).
				Return("refresh-token", nil)
			tokens.
				On("ParseRefreshToken", mock.Anything).
				Return(&entity.RefreshClaims{}, errors.New("parse error"))

			res, err := uc.Execute(dto.Login{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenParse, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestLoginUC_CreateSessionError(t *testing.T) {
	for _, testCase := range data.LoginTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			sessions := new(mocks.MockSessionsRepo)
			tokens := new(mocks.MockTokensRepo)
			uc := &usecase.LoginUC{
				Users:    users,
				Sessions: sessions,
				Tokens:   tokens,
			}

			existingUser := helpers.MakeTestUser(testCase.Email, testCase.Password, testCase.Role)

			users.
				On("GetUserByEmail", mock.Anything).
				Return(existingUser, nil)
			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("access-token", nil)
			tokens.
				On("GenerateRefreshToken", mock.Anything).
				Return("refresh-token", nil)
			tokens.
				On("ParseRefreshToken", mock.Anything).
				Return(helpers.MakeRefreshClaims("jti-x", time.Now().Add(24*time.Hour), time.Now()), nil)
			sessions.
				On("CreateSession", mock.Anything).
				Return(errors.New("session error"))

			res, err := uc.Execute(dto.Login{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrSessionSave, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}
