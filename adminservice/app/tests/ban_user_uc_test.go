package tests

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/tests/data"
	"AdsService/adminservice/app/tests/mocks"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestBanUserUC_Success(t *testing.T) {
	for c, testCase := range data.BanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.BanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", testCase.AdminUserID).
				Return("admin", nil)
			users.
				On("BanUser", testCase.RequestedUserID).
				Return(nil)

			out, err := uc.Execute(dto.BanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedAnswer, out.Banned)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestBanUserUC_GetUserError(t *testing.T) {
	for c, testCase := range data.BanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.BanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("", errors.New("get user error"))

			out, err := uc.Execute(dto.BanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUserRole, err)
			assert.Equal(t, out.Banned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestBanUserUC_NotAdminError(t *testing.T) {
	for c, testCase := range data.BanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.BanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("user", nil)

			out, err := uc.Execute(dto.BanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrNotAdmin, err)
			assert.Equal(t, out.Banned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestBanUserUC_BanUserError(t *testing.T) {
	for c, testCase := range data.BanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.BanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			users.
				On("BanUser", mock.Anything).
				Return(errors.New("ban user error"))

			out, err := uc.Execute(dto.BanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrBanUser, err)
			assert.Equal(t, out.Banned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}
