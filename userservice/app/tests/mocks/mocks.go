package mocks

import (
	"AdsService/userservice/domain/entity"
	"github.com/stretchr/testify/mock"
	"io"
)

type MockProfilesRepo struct{ mock.Mock }

func (m *MockProfilesRepo) AddProfile(userID uint64, name, phone string) (*entity.Profile, error) {
	args := m.Called(userID, name, phone)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) GetProfile(userID uint64) (*entity.Profile, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) EnableNotifications(userID uint64) (*entity.Profile, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) DisableNotifications(userID uint64) (*entity.Profile, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) UpdateProfileSubscriptions(userID uint64, subscriptions []string) (*entity.Profile, error) {
	args := m.Called(userID, subscriptions)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) UpdateProfileName(userID uint64, name string) (*entity.Profile, error) {
	args := m.Called(userID, name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) UpdateProfilePhone(userID uint64, phone string) (*entity.Profile, error) {
	args := m.Called(userID, phone)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) UpdateProfilePhoto(userID uint64, photoID string) (*entity.Profile, error) {
	args := m.Called(userID, photoID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) UpdateProfileTime(userID uint64) (*entity.Profile, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

type MockPhotosRepo struct{ mock.Mock }

func (m *MockPhotosRepo) UploadPhoto(userID uint64, title, contentType string, r io.Reader, size int64) (string, error) {
	args := m.Called(userID, title, contentType, r, size)
	return args.String(0), args.Error(1)
}
