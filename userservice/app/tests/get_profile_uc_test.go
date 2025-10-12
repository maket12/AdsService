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

func TestGetProfileUC_Success(t *testing.T) {
	for _, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.GetProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("GetProfile", testCase.UserID).
				Return(helpers.MakeTestProfile(testCase.UserID, "", "", false, nil, ""), nil)

			out, err := uc.Execute(dto.GetProfile{
				UserID: testCase.UserID,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.UserID, out.UserID)

			profiles.AssertExpectations(t)
		})
	}
}

func TestGetProfileUC_GetProfileError(t *testing.T) {
	for _, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.GetProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("GetProfile", mock.Anything).
				Return(nil, errors.New("get profile error"))

			_, err := uc.Execute(dto.GetProfile{
				UserID: testCase.UserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrGetProfile, err)

			profiles.AssertExpectations(t)
		})
	}
}

func TestGetProfileUC_NotFoundError(t *testing.T) {
	for _, testCase := range data.GetProfileTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.GetProfileUC{
				Profiles: &profiles,
			}

			profiles.
				On("GetProfile", mock.Anything).
				Return(nil, nil)

			_, err := uc.Execute(dto.GetProfile{
				UserID: testCase.UserID,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrProfileNotFound, err)

			profiles.AssertExpectations(t)
		})
	}
}
