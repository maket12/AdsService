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

func TestGetProfileUC_Success(t *testing.T) {
	for c, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)

			uc := &usecase.GetProfileUC{
				Users:    users,
				Profiles: profiles,
			}

			users.
				On("GetUserRole", testCase.AdminUserID).
				Return("admin", nil)
			profiles.
				On("GetProfile", testCase.RequestedUserID).
				Return(helpers.MakeTestProfile(testCase.RequestedUserID), nil)

			out, err := uc.Execute(dto.GetProfileDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.RequestedUserID, out.UserID)

			users.AssertExpectations(t)
			profiles.AssertExpectations(t)
		})
	}
}

func TestGetProfileUC_GetUserError(t *testing.T) {
	for c, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetProfileUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("", errors.New("get user error"))

			_, err := uc.Execute(dto.GetProfileDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetUserRole, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetProfileUC_NotAdminError(t *testing.T) {
	for c, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)

			uc := &usecase.GetProfileUC{
				Users: users,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("user", nil)

			_, err := uc.Execute(dto.GetProfileDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrNotAdmin, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetProfileUC_GetProfileError(t *testing.T) {
	for c, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)

			uc := &usecase.GetProfileUC{
				Users:    users,
				Profiles: profiles,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			profiles.
				On("GetProfile", mock.Anything).
				Return(nil, errors.New("get profile error"))

			_, err := uc.Execute(dto.GetProfileDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetProfile, err)

			users.AssertExpectations(t)
		})
	}
}

func TestGetProfileUC_ProfileNotFoundError(t *testing.T) {
	for c, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(c+1), func(t *testing.T) {
			users := new(mocks.MockUsersRepo)
			profiles := new(mocks.MockProfilesRepo)

			uc := &usecase.GetProfileUC{
				Users:    users,
				Profiles: profiles,
			}

			users.
				On("GetUserRole", mock.Anything).
				Return("admin", nil)
			profiles.
				On("GetProfile", mock.Anything).
				Return(nil, nil)

			_, err := uc.Execute(dto.GetProfileDTO{
				UserID:          testCase.AdminUserID,
				RequestedUserID: testCase.RequestedUserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrProfileNotFound, err)

			users.AssertExpectations(t)
		})
	}
}
