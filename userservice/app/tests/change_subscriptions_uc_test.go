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

func TestChangeSubscriptionsUC_Success(t *testing.T) {
	for _, testCase := range data.ChangeSubscriptionsTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.ChangeSubscriptionsUC{
				Profiles: &profiles,
			}

			profiles.
				On("UpdateProfileSubscriptions", testCase.UserID, testCase.Subscriptions).
				Return(helpers.MakeTestProfile(testCase.UserID, "", "", false, testCase.Subscriptions, ""), nil)

			out, err := uc.Execute(dto.ChangeSubscriptionsDTO{
				UserID:        testCase.UserID,
				Subscriptions: testCase.Subscriptions,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.UserID, out.UserID)
			assert.Equal(t, testCase.Subscriptions, out.Subscriptions)

			profiles.AssertExpectations(t)
		})
	}
}

func TestChangeSubscriptionsUC_ChangeSubscriptionsError(t *testing.T) {
	for _, testCase := range data.ChangeSubscriptionsTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.ChangeSubscriptionsUC{
				Profiles: &profiles,
			}

			profiles.
				On("UpdateProfileSubscriptions", mock.Anything, mock.Anything).
				Return(nil, errors.New("subscriptions error"))

			_, err := uc.Execute(dto.ChangeSubscriptionsDTO{
				UserID:        testCase.UserID,
				Subscriptions: testCase.Subscriptions,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrChangeSubscriptions, err)

			profiles.AssertExpectations(t)
		})
	}
}
