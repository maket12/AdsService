package tests

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/tests/data"
	"ads/adminservice/app/tests/helpers"
	"ads/adminservice/app/tests/mocks"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetUserUC_Success(t *testing.T) {
	for c, testCase := range data.GetUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", testCase.AdminUserID).
				Return("admin", nil)
			users.
				On("GetUserByID", testCase.RequestedUserID).
				Return(helpers.MakeTestUser(testCase.RequestedUserID), nil)

			out, err := uc.Execute(dto.GetUser{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.NoError(t, err)
			assert.Equal(t, out.UserID, testCase.RequestedUserID)

			users.AssertExpectations(t)
		})
	}
}

func TestGetUserUC_GetRoleError(t *testing.T) {
	for c, testCase := range data.GetUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("", errors.New("get role error"))

			_, err := uc.Execute(dto.GetUser{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUserRole, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetUserUC_NotAdminError(t *testing.T) {
	for c, testCase := range data.GetUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("user", nil)

			_, err := uc.Execute(dto.GetUser{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrNotAdmin, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetUserUC_GetUserError(t *testing.T) {
	for c, testCase := range data.GetUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			users.
				On("GetUserByID", mock.Anything).
				Return(nil, errors.New("get user error"))

			_, err := uc.Execute(dto.GetUser{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUser, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetUserUC_UserNotFoundError(t *testing.T) {
	for c, testCase := range data.GetUserTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetUserUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			users.
				On("GetUserByID", mock.Anything).
				Return(nil, nil)

			_, err := uc.Execute(dto.GetUser{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrUserNotFound, err)

			users.AssertExpectations(t)
		})
	}
}
