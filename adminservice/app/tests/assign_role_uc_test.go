package tests

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/tests/data"
	"ads/adminservice/app/tests/mocks"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAssignRoleUC_Success(t *testing.T) {
	for _, testCase := range data.AssignRoleTestCases {
		t.Run(fmt.Sprint(testCase.RequestedUserID), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.AssignRoleUC{
				Users: users,
			}

			users.
				On("GetUserRole", testCase.AdminUserID).
				Return("admin", nil)
			users.
				On("EnhanceUser", testCase.RequestedUserID).
				Return(nil)

			out, err := uc.Execute(dto.AssignRole{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedAnswer, out.Assigned)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}

func TestAssignRoleUC_GetRoleError(t *testing.T) {
	for c, testCase := range data.AssignRoleTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.AssignRoleUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("", errors.New("get role error"))

			_, err := uc.Execute(dto.AssignRole{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUserRole, err)

			users.AssertExpectations(t)
		})
	}
}

func TestAssignRoleUC_NotAdminError(t *testing.T) {
	for c, testCase := range data.AssignRoleTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.AssignRoleUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("user", nil)

			_, err := uc.Execute(dto.AssignRole{
				AdminUserID:     testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrNotAdmin, err)

			users.AssertExpectations(t)
		})
	}
}

func TestAssignRoleUC_EnhanceError(t *testing.T) {
	for _, testCase := range data.AssignRoleTestCases {
		t.Run(fmt.Sprint(testCase.RequestedUserID), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.AssignRoleUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			users.
				On("EnhanceUser", mock.Anything).
				Return(errors.New("enhance user error"))

			out, err := uc.Execute(dto.AssignRole{
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrEnhanceUser, err)
			assert.Equal(t, out.Assigned, false)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
		})
	}
}
