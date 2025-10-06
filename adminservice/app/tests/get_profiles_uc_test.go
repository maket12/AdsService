package tests

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/tests/data"
	"AdsService/adminservice/app/tests/helpers"
	"AdsService/adminservice/app/tests/mocks"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetProfilesUC_Success(t *testing.T) {
	for c, testCase := range data.GetProfilesTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)

			uc := &usecase.GetProfilesUC{
				Users:    users,
				Profiles: profiles,
			}

			users.
				On("GetUserRole", testCase.AdminUserID).
				Return("admin", nil)
			profiles.
				On("GetAllProfiles", testCase.Limit, testCase.Offset).
				Return(helpers.MakeTestProfiles(), nil)

			_, err := uc.Execute(dto.GetProfilesListDTO{
				UserID: testCase.AdminUserID,
				Limit:  testCase.Limit,
				Offset: testCase.Offset,
			})

			assert.NoError(t, err)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
		})
	}
}

func TestGetProfilesUC_GetUserError(t *testing.T) {
	for c, testCase := range data.GetProfilesTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetProfilesUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("", errors.New("get user error"))

			_, err := uc.Execute(dto.GetProfilesListDTO{
				UserID: testCase.AdminUserID,
				Limit:  testCase.Limit,
				Offset: testCase.Offset,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUserRole, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetProfilesUC_NotAdminError(t *testing.T) {
	for c, testCase := range data.GetProfilesTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetProfilesUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("user", nil)

			_, err := uc.Execute(dto.GetProfilesListDTO{
				UserID: testCase.AdminUserID,
				Limit:  testCase.Limit,
				Offset: testCase.Offset,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrNotAdmin, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetProfilesUC_GetProfilesError(t *testing.T) {
	for c, testCase := range data.GetProfilesTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)

			uc := &usecase.GetProfilesUC{
				Users:    users,
				Profiles: profiles,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			profiles.
				On("GetAllProfiles", mock.Anything, mock.Anything).
				Return(nil, errors.New("get profiles error"))

			_, err := uc.Execute(dto.GetProfilesListDTO{
				UserID: testCase.AdminUserID,
				Limit:  testCase.Limit,
				Offset: testCase.Offset,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetProfiles, err)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
		})
	}
}
