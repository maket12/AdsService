package tests

import (
	"AdsService/userservice/app/dto"
	"AdsService/userservice/app/tests/data"
	"AdsService/userservice/app/tests/helpers"
	"AdsService/userservice/app/tests/mocks"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/app/usecase"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddProfileUC_Success(t *testing.T) {
	for _, testCase := range data.AddProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.AddProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("GetProfile", testCase.UserID).
				Return(nil, errors.New("profile not found"))

			profiles.
				On("AddProfile", testCase.UserID, testCase.Name, testCase.Phone).
				Return(helpers.MakeTestProfile(testCase.UserID, testCase.Name, testCase.Phone, true, nil, ""), nil)

			out, err := uc.Execute(dto.AddProfileDTO{
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

func TestAddProfileUC_AddProfileError(t *testing.T) {
	for _, testCase := range data.AddProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.AddProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("GetProfile", mock.Anything).
				Return(nil, errors.New("profile not found"))

			profiles.
				On("AddProfile", mock.Anything, mock.Anything, mock.Anything).
				Return(nil, errors.New("profile error"))

			_, err := uc.Execute(dto.AddProfileDTO{
				UserID: testCase.UserID,
				Name:   testCase.Name,
				Phone:  testCase.Phone,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrAddProfile, err)

			profiles.AssertExpectations(t)
		})
	}
}
