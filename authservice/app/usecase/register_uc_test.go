package usecase

import (
	"AdsService/authservice/app/dto"
	"AdsService/authservice/app/uc_errors"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"AdsService/authservice/domain/entity"
)

// --- Mocks ---

// UsersRepo mock
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
	if u := args.Get(0); u != nil {
		return u.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// ProfilesRepo mock
type MockProfilesRepo struct{ mock.Mock }

func (m *MockProfilesRepo) AddProfile(userID uint64, name string, phone string) (*entity.Profile, error) {
	args := m.Called(userID, name, phone)
	if p := args.Get(0); p != nil {
		return p.(*entity.Profile), args.Error(1)
	}
	return nil, args.Error(1)
}

// TokensRepo mock
//type MockTokensRepo struct{ mock.Mock }
//
//func (m *MockTokensRepo) GenerateAccessToken(userID uint64, email string, role string) (string, error) {
//	args := m.Called(userID, email, role)
//	return args.String(0), args.Error(1)
//}
//func (m *MockTokensRepo) GenerateRefreshToken(userID uint64) (string, error) {
//	args := m.Called(userID)
//	return args.String(0), args.Error(1)
//}
//func (m *MockTokensRepo) ParseAccessToken(tokenStr string) (*entity.AccessClaims, error) {
//	args := m.Called(tokenStr)
//	if c := args.Get(0); c != nil {
//		return c.(*entity.AccessClaims), args.Error(1)
//	}
//	return nil, args.Error(1)
//}
//
//func (m *MockTokensRepo) ParseRefreshToken(tokenStr string) (*entity.RefreshClaims, error) {
//	args := m.Called(tokenStr)
//	if c := args.Get(0); c != nil {
//		return c.(*entity.RefreshClaims), args.Error(1)
//	}
//	return nil, args.Error(1)
//}

// SessionsRepo mock
type MockSessionsRepo struct{ mock.Mock }

func (m *MockSessionsRepo) InsertSession(s *entity.Session) error {
	args := m.Called(s)
	return args.Error(0)
}
func (m *MockSessionsRepo) GetSessionByJTI(jti string) (*entity.Session, error) {
	args := m.Called(jti)
	if s := args.Get(0); s != nil {
		return s.(*entity.Session), args.Error(1)
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
func (m *MockSessionsRepo) Save(userID uint64, refreshToken string) error {
	args := m.Called(userID, refreshToken)
	return args.Error(0)
}

// --- Tests ---

// Mocks
type MockTokensRepo struct{ mock.Mock }

func (m *MockTokensRepo) GenerateAccessToken(userID uint64, email string, role string) (string, error) {
	args := m.Called(userID, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockTokensRepo) GenerateRefreshToken(userID uint64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokensRepo) ParseAccessToken(tokenStr string) (*entity.AccessClaims, error) {
	args := m.Called(tokenStr)
	if c := args.Get(0); c != nil {
		return c.(*entity.AccessClaims), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTokensRepo) ParseRefreshToken(tokenStr string) (*entity.RefreshClaims, error) {
	args := m.Called(tokenStr)
	if c := args.Get(0); c != nil {
		return c.(*entity.RefreshClaims), args.Error(1)
	}
	return nil, args.Error(1)
}

// Test for RegisterUC
func TestRegisterUC_Success(t *testing.T) {
	// Моки
	users := new(MockUsersRepo)
	profiles := new(MockProfilesRepo)
	tokens := new(MockTokensRepo)
	sessions := new(MockSessionsRepo)

	// Создаем экземпляр RegisterUC
	uc := &RegisterUC{
		Users:    users,
		Profiles: profiles,
		Tokens:   tokens,
		Sessions: sessions,
	}

	users.On("CheckUserExist", "test@test.com").Return(false, nil).Once() // проверка, что email не существует
	users.On("AddUser", mock.AnythingOfType("*entity.User")).Run(func(args mock.Arguments) {
		u := args.Get(0).(*entity.User)
		u.ID = 1
	}).Return(nil)

	profiles.On("AddProfile", uint64(1), "undefined", "undefined").Return(&entity.Profile{UserID: 1}, nil)

	tokens.On("GenerateAccessToken", uint64(1), "test@test.com", "user").
		Return("access123", nil)
	tokens.On("GenerateRefreshToken", uint64(1)).
		Return("refresh123", nil)

	tokens.On("ParseRefreshToken", "refresh123").
		Return(&entity.RefreshClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID: "token123",
				ExpiresAt: &jwt.NumericDate{
					Time: time.Now().Add(time.Hour),
				},
			},
		}, nil)

	sessions.On("InsertSession", mock.AnythingOfType("*entity.Session")).Return(nil)

	// Создание RegisterDTO
	registerDTO := &dto.RegisterDTO{
		Email:    "test@test.com",
		Password: "pass",
	}

	res, err := uc.Execute(*registerDTO)
	
	assert.NoError(t, err)
	assert.Equal(t, "access123", res.AccessToken)
	assert.Equal(t, "refresh123", res.RefreshToken)

	// Проверка ожиданий для моков
	users.AssertExpectations(t)
	profiles.AssertExpectations(t)
	tokens.AssertExpectations(t)
	sessions.AssertExpectations(t)
}

func TestRegisterUC_UserAlreadyExists(t *testing.T) {
	users := new(MockUsersRepo)
	uc := &RegisterUC{
		Users:    users,
		Profiles: new(MockProfilesRepo),
		Tokens:   new(MockTokensRepo),
		Sessions: new(MockSessionsRepo),
	}

	users.On("CheckUserExist", "exists@test.com").Return(true, nil)

	// Create RegisterDTO
	registerDTO := dto.RegisterDTO{
		Email:    "exists@test.com",
		Password: "pass",
	}
	_, err := uc.Execute(registerDTO)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, uc_errors.ErrAddUser))
}
