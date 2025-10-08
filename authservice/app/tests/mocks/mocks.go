package mocks

import (
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"github.com/stretchr/testify/mock"
)

// ---- Users ----

type MockUsersRepo struct{ mock.Mock }

func (m *MockUsersRepo) CheckUserExist(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUsersRepo) AddUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUsersRepo) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

// ---- Profiles ----

type MockProfilesRepo struct{ mock.Mock }

func (m *MockProfilesRepo) AddProfile(userID uint64, name, phone string) (*entity.Profile, error) {
	args := m.Called(userID, name, phone)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Profile), args.Error(1)
}

// ---- Tokens ----

type MockTokensRepo struct{ mock.Mock }

var _ port.TokenRepository = (*MockTokensRepo)(nil)

func (m *MockTokensRepo) GenerateAccessToken(userID uint64, email, role string) (string, error) {
	args := m.Called(userID, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockTokensRepo) GenerateRefreshToken(userID uint64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokensRepo) ParseAccessToken(tokenStr string) (*entity.AccessClaims, error) {
	args := m.Called(tokenStr)
	if v := args.Get(0); v != nil {
		return v.(*entity.AccessClaims), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTokensRepo) ParseRefreshToken(tokenStr string) (*entity.RefreshClaims, error) {
	args := m.Called(tokenStr)
	if v := args.Get(0); v != nil {
		return v.(*entity.RefreshClaims), args.Error(1)
	}
	return nil, args.Error(1)
}

// ---- Sessions ----

type MockSessionsRepo struct{ mock.Mock }

var _ port.SessionRepository = (*MockSessionsRepo)(nil)

func (m *MockSessionsRepo) CreateSession(s *entity.Session) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockSessionsRepo) GetSessionByJTI(jti string) (*entity.Session, error) {
	args := m.Called(jti)
	if v := args.Get(0); v != nil {
		return v.(*entity.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionsRepo) RevokeByJTI(jti string) error {
	args := m.Called(jti)
	return args.Error(0)
}

func (m *MockSessionsRepo) RevokeAllByUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockSessionsRepo) CleanupExpired() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSessionsRepo) RotateSession(oldJTI string, newS *entity.Session) error {
	args := m.Called(oldJTI, newS)
	return args.Error(0)
}

func (m *MockTokensRepo) CreateSession(s *entity.Session) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockTokensRepo) GetSessionByJTI(jti string) (*entity.Session, error) {
	args := m.Called(jti)
	if v := args.Get(0); v != nil {
		return (v).(*entity.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTokensRepo) RevokeByJTI(jti string) error {
	args := m.Called(jti)
	return args.Error(0)
}

func (m *MockTokensRepo) RevokeAllByUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockTokensRepo) CleanupExpired() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTokensRepo) RotateSession(oldJTI string, newS *entity.Session) error {
	args := m.Called(oldJTI, newS)
	return args.Error(0)
}
