package tests

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/tests/data"
	"ads/userservice/app/tests/helpers"
	"ads/userservice/app/tests/mocks"
	"ads/userservice/app/uc_errors"
	"ads/userservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateProfileUC_Success(t *testing.T) {
	for _, testCase := range data.UpdateProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.UpdateProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("UpdateProfileName", testCase.UserID, testCase.Name).
				Return(nil, nil)

			profiles.
				On("UpdateProfilePhone", testCase.UserID, testCase.Phone).
				Return(helpers.MakeTestProfile(testCase.UserID, testCase.Name, testCase.Phone, true, nil, ""), nil)

			out, err := uc.Execute(dto.UpdateProfile{
				UserID: testCase.UserID,
				Name:   testCase.Name,
				Phone:  testCase.Phone,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.UserID, out.UserID)

			profiles.AssertExpectations(t)
		})
	}
}

func TestUpdateProfileUC_UpdateNameError(t *testing.T) {
	for _, testCase := range data.UpdateProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.UpdateProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("UpdateProfileName", mock.Anything, mock.Anything).
				Return(nil, errors.New("profile error"))

			_, err := uc.Execute(dto.UpdateProfile{
				UserID: testCase.UserID,
				Name:   testCase.Name,
				Phone:  testCase.Phone,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrUpdateProfile, err)

			profiles.AssertExpectations(t)
		})
	}
}

func TestUpdateProfileUC_UpdatePhoneError(t *testing.T) {
	for _, testCase := range data.UpdateProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.UpdateProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("UpdateProfileName", mock.Anything, mock.Anything).
				Return(nil, nil)
			profiles.
				On("UpdateProfilePhone", mock.Anything, mock.Anything).
				Return(nil, errors.New("profile error"))

			_, err := uc.Execute(dto.UpdateProfile{
				UserID: testCase.UserID,
				Name:   testCase.Name,
				Phone:  testCase.Phone,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrUpdateProfile, err)

			profiles.AssertExpectations(t)
		})
	}
}
