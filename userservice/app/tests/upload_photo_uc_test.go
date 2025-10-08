package tests

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/tests/data"
	"ads/userservice/app/tests/helpers"
	"ads/userservice/app/tests/mocks"
	"ads/userservice/app/uc_errors"
	"ads/userservice/app/usecase"
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUploadPhotoUC_Success(t *testing.T) {
	for _, testCase := range data.UploadPhotoTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}
			photos := mocks.MockPhotosRepo{}

			uc := &usecase.UploadPhotoUC{
				Profiles: &profiles,
				Photos:   &photos,
			}

			photos.
				On("UploadPhoto", testCase.UserID, testCase.FileName, testCase.ContentType, bytes.NewReader(testCase.Data), int64(len(testCase.Data))).
				Return(testCase.ExpectedAnswer, nil)
			profiles.
				On("UpdateProfilePhoto", testCase.UserID, testCase.ExpectedAnswer).
				Return(helpers.MakeTestProfile(testCase.UserID, "", "", false, nil, testCase.ExpectedAnswer), nil)

			out, err := uc.Execute(dto.UploadPhoto{
				UserID:      testCase.UserID,
				Data:        testCase.Data,
				Filename:    testCase.FileName,
				ContentType: testCase.ContentType,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.UserID, out.UserID)
			assert.Equal(t, testCase.ExpectedAnswer, out.PhotoID)

			profiles.AssertExpectations(t)
			photos.AssertExpectations(t)
		})
	}
}

func TestUploadPhotoUC_DataError(t *testing.T) {
	for _, testCase := range data.UploadPhotoErrTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}
			photos := mocks.MockPhotosRepo{}

			uc := &usecase.UploadPhotoUC{
				Profiles: &profiles,
				Photos:   &photos,
			}

			_, err := uc.Execute(dto.UploadPhoto{
				UserID:      testCase.UserID,
				Data:        testCase.Data,
				Filename:    testCase.FileName,
				ContentType: testCase.ContentType,
			})

			assert.Error(t, err)

			if len(testCase.Data) == 0 {
				assert.Equal(t, uc_errors.ErrEmptyDataPhoto, err)
			} else if testCase.FileName == "" {
				assert.Equal(t, uc_errors.ErrEmptyFilenamePhoto, err)
			} else {
				assert.Equal(t, uc_errors.ErrEmptyContentTypePhoto, err)
			}

			profiles.AssertExpectations(t)
			photos.AssertExpectations(t)
		})
	}
}

func TestUploadPhotoUC_MongoUploadError(t *testing.T) {
	for _, testCase := range data.UploadPhotoTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			photos := mocks.MockPhotosRepo{}

			uc := &usecase.UploadPhotoUC{
				Photos: &photos,
			}

			photos.
				On("UploadPhoto", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("", errors.New("mongo error"))

			_, err := uc.Execute(dto.UploadPhoto{
				UserID:      testCase.UserID,
				Data:        testCase.Data,
				Filename:    testCase.FileName,
				ContentType: testCase.ContentType,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrMongoUploadPhoto, err)

			photos.AssertExpectations(t)
		})
	}
}

func TestUploadPhotoUC_UpdatePhotoError(t *testing.T) {
	for _, testCase := range data.UploadPhotoTestCases {
		t.Run(fmt.Sprint(testCase.UserID), func(t *testing.T) {
			profiles := mocks.MockProfilesRepo{}
			photos := mocks.MockPhotosRepo{}

			uc := &usecase.UploadPhotoUC{
				Profiles: &profiles,
				Photos:   &photos,
			}

			photos.
				On("UploadPhoto", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("hex-string", nil)
			profiles.
				On("UpdateProfilePhoto", mock.Anything, mock.Anything).
				Return(nil, errors.New("update photo error"))

			_, err := uc.Execute(dto.UploadPhoto{
				UserID:      testCase.UserID,
				Data:        testCase.Data,
				Filename:    testCase.FileName,
				ContentType: testCase.ContentType,
			})

			assert.Error(t, err)
			assert.Equal(t, uc_errors.ErrUpdatePhoto, err)

			photos.AssertExpectations(t)
			profiles.AssertExpectations(t)
		})
	}
}
