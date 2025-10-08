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

func TestChangeSettingsUC_Success(t *testing.T) {
	for _, testCase := range data.ChangeSettingsTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.ChangeSettingsUC{
				Profiles: &profiles,
			}

			if testCase.NotificationsEnabled {
				profiles.
					On("EnableNotifications", testCase.UserID).
					Return(helpers.MakeTestProfile(testCase.UserID, "", "", true, nil, ""), nil)
			} else {
				profiles.
					On("DisableNotifications", testCase.UserID).
					Return(helpers.MakeTestProfile(testCase.UserID, "", "", false, nil, ""), nil)
			}

			out, err := uc.Execute(dto.ChangeSettings{
				UserID:               testCase.UserID,
				NotificationsEnabled: testCase.NotificationsEnabled,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.UserID, out.UserID)
			assert.Equal(t, testCase.NotificationsEnabled, out.NotificationsEnabled)

			profiles.AssertExpectations(t)
		})
	}
}

func TestAddProfileUC_ChangeSettingsError(t *testing.T) {
	for _, testCase := range data.ChangeSettingsTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}

			uc := &usecase.ChangeSettingsUC{
				Profiles: &profiles,
			}

			if testCase.NotificationsEnabled {
				profiles.
					On("EnableNotifications", mock.Anything).
					Return(nil, errors.New("notifications error"))
			} else {
				profiles.
					On("DisableNotifications", mock.Anything).
					Return(nil, errors.New("notifications error"))
			}

			_, err := uc.Execute(dto.ChangeSettings{
				UserID:               testCase.UserID,
				NotificationsEnabled: testCase.NotificationsEnabled,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrChangeSettings, err)

			profiles.AssertExpectations(t)
		})
	}
}
