package mocks

import (
	"ads/adminservice/domain/entity"
	"github.com/stretchr/testify/mock"
)

// ---- Users ----

type MockUsersRepo struct{ mock.Mock }

func (m *MockUsersRepo) GetUserRole(userID uint64) (string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *MockUsersRepo) EnhanceUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUsersRepo) BanUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUsersRepo) UnbanUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUsersRepo) GetUserByID(userID uint64) (*entity.User, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

// ---- Profiles ----

type MockProfilesRepo struct{ mock.Mock }

func (m *MockProfilesRepo) GetProfile(userID uint64) (*entity.Profile, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfilesRepo) GetAllProfiles(limit, offset uint32) ([]entity.Profile, error) {
	args := m.Called(limit, offset)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Profile), args.Error(1)
}
