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

func TestUnbanUserUC_Success(t *testing.T) {
	for c, testCase := range data.UnbanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.UnbanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", testCase.AdminUserID).
				Return("admin", nil)
			users.
				On("UnbanUser", testCase.RequestedUserID).
				Return(nil)

			out, err := uc.Execute(dto.UnbanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedAnswer, out.Unbanned)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestUnbanUserUC_GetUserError(t *testing.T) {
	for c, testCase := range data.UnbanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.UnbanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("", errors.New("get user error"))

			out, err := uc.Execute(dto.UnbanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUserRole, err)
			assert.Equal(t, out.Unbanned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestUnbanUserUC_NotAdminError(t *testing.T) {
	for c, testCase := range data.UnbanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.UnbanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("user", nil)

			out, err := uc.Execute(dto.UnbanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrNotAdmin, err)
			assert.Equal(t, out.Unbanned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestUnbanUserUC_UnbanUserError(t *testing.T) {
	for c, testCase := range data.UnbanUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.UnbanUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			users.
				On("UnbanUser", mock.Anything).
				Return(errors.New("unban user error"))

			out, err := uc.Execute(dto.UnbanUserDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrUnbanUser, err)
			assert.Equal(t, out.Unbanned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}
