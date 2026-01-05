package tests

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/tests/data"
	"ads/authservice/internal/app/tests/helpers"
	"ads/authservice/internal/app/tests/mocks"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/usecase"
	entity2 "ads/authservice/internal/domain/entity"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUC_Success(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)
			tokens := new(mocks.MockTokensRepo)
			sessions := new(mocks.MockSessionsRepo)

			uc := &usecase.RegisterUC{
				Users:    users,
				Sessions: sessions,
				Tokens:   tokens,
				Profiles: profiles,
			}

			users.
				On("AddUser", mock.MatchedBy(func(user *entity2.User) bool {
					return user.Email == testCase.Email &&
						user.Role == testCase.Role &&
						strings.HasPrefix(user.Password, "$2a$")
				})).Return(nil)

			profiles.
				On("AddProfile", mock.AnythingOfType("uint64"), "undefined", "undefined").
				Return(&entity2.Profile{}, nil)

			tokens.
				On("GenerateAccessToken", mock.AnythingOfType("uint64"), testCase.Email, testCase.Role).
				Return(testCase.ExpectedAccessToken, nil)

			tokens.
				On("GenerateRefreshToken", mock.AnythingOfType("uint64")).
				Return(testCase.ExpectedRefreshToken, nil)

			tokens.
				On("ParseRefreshToken", testCase.ExpectedRefreshToken).
				Return(helpers.MakeRefreshClaims("jti-x", time.Now().Add(24*time.Hour), time.Now()), nil)

			sessions.
				On("CreateSession", mock.MatchedBy(func(s *entity2.Session) bool {
					return s.JTI == "jti-x" && !s.IssuedAt.IsZero() && s.ExpiresAt.After(s.IssuedAt)
				})).
				Return(nil)

			out, err := uc.Execute(dto.Register{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedAccessToken, out.AccessToken)
			assert.Equal(t, testCase.ExpectedRefreshToken, out.RefreshToken)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
			tokens.AssertExpectations(t)
			sessions.AssertExpectations(t)
		})
	}
}

func TestRegisterUC_AddUserError(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			uc := &usecase.RegisterUC{
				Users:    users,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   new(mocks.MockTokensRepo),
				Profiles: new(mocks.MockProfilesRepo),
			}

			users.
				On("AddUser", mock.AnythingOfType("*entity.User")).
				Return(errors.New("db insert failed"))

			res, err := uc.Execute(dto.Register{Email: testCase.Email, Password: testCase.Password})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrAddUser, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
		})
	}
}

func TestRegisterUC_AddProfileError(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)
			uc := &usecase.RegisterUC{
				Users:    users,
				Profiles: profiles,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   new(mocks.MockTokensRepo),
			}

			users.
				On("AddUser", mock.AnythingOfType("*entity.User")).Return(nil)
			profiles.
				On("AddProfile", mock.Anything, mock.Anything, mock.Anything).
				Return(nil, errors.New("profile error"))

			res, err := uc.Execute(dto.Register{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrAddProfile, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
		})
	}
}

func TestRegisterUC_GenerateAccessTokenError(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)
			tokens := new(mocks.MockTokensRepo)

			uc := &usecase.RegisterUC{
				Users:    users,
				Profiles: profiles,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   tokens,
			}

			users.
				On("AddUser", mock.AnythingOfType("*entity.User")).Return(nil)
			profiles.
				On("AddProfile", mock.Anything, mock.Anything, mock.Anything).
				Return(&entity2.Profile{}, nil)
			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("", errors.New("token error"))

			res, err := uc.Execute(dto.Register{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenIssue, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestRegisterUC_GenerateRefreshTokenError(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)
			tokens := new(mocks.MockTokensRepo)

			uc := &usecase.RegisterUC{
				Users:    users,
				Profiles: profiles,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   tokens,
			}

			users.
				On("AddUser", mock.AnythingOfType("*entity.User")).Return(nil)
			profiles.
				On("AddProfile", mock.Anything, mock.Anything, mock.Anything).
				Return(&entity2.Profile{}, nil)
			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("access-token", nil)
			tokens.
				On("GenerateRefreshToken", mock.Anything).
				Return("", errors.New("token error"))

			res, err := uc.Execute(dto.Register{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenIssue, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestRegisterUC_ParseRefreshTokenError(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)
			tokens := new(mocks.MockTokensRepo)

			uc := &usecase.RegisterUC{
				Users:    users,
				Profiles: profiles,
				Sessions: new(mocks.MockSessionsRepo),
				Tokens:   tokens,
			}

			users.
				On("AddUser", mock.AnythingOfType("*entity.User")).Return(nil)
			profiles.
				On("AddProfile", mock.Anything, mock.Anything, mock.Anything).
				Return(&entity2.Profile{}, nil)
			tokens.
				On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).
				Return("access-token", nil)
			tokens.
				On("GenerateRefreshToken", mock.Anything).
				Return("refresh-token", nil)
			tokens.
				On("ParseRefreshToken", "refresh-token").
				Return(nil, errors.New("parse error"))

			res, err := uc.Execute(dto.Register{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrTokenParse, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestRegisterUC_CreateSessionError(t *testing.T) {
	for _, testCase := range data.RegisterTestCases {
		t.Run(testCase.Email, func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)
			sessions := new(mocks.MockSessionsRepo)
			tokens := new(mocks.MockTokensRepo)

			uc := &usecase.RegisterUC{
				Users:    users,
				Profiles: profiles,
				Sessions: sessions,
				Tokens:   tokens,
			}

			users.
				On("AddUser", mock.AnythingOfType("*entity.User")).Return(nil)
			profiles.
				On("AddProfile", mock.Anything, mock.Anything, mock.Anything).
				Return(&entity2.Profile{}, nil)
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

			res, err := uc.Execute(dto.Register{
				Email:    testCase.Email,
				Password: testCase.Password,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrSessionSave, err)
			assert.Empty(t, res.AccessToken)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
			sessions.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}
