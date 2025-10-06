package jwt

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
)

type TokenRepository struct {
	accessSecret  []byte
	refreshSecret []byte
	logger        *slog.Logger
}

func NewTokenRepository(accessSecret, refreshSecret string, logger *slog.Logger) port.TokenRepository {
	return &TokenRepository{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		logger:        logger,
	}
}

func (s *TokenRepository) GenerateAccessToken(userID uint64, email, role string) (string, error) {
	claims := &entity.AccessClaims{
		Type:   "access",
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	s.logger.Info("Generated access token for user[id=%v]", userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessSecret)
}

func (s *TokenRepository) GenerateRefreshToken(userID uint64) (string, error) {
	claims := entity.RefreshClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	s.logger.Info("Generated refresh token for user[id=%v]", userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecret)
}

func (s *TokenRepository) ParseAccessToken(tokenStr string) (*entity.AccessClaims, error) {
	var c entity.AccessClaims
	parsedToken, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {
		return s.accessSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), jwt.WithLeeway(30*time.Second))

	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	if c.Type != "access" {
		return nil, errors.New("wrong token type")
	}
	return &c, nil
}

func (s *TokenRepository) ParseRefreshToken(tokenStr string) (*entity.RefreshClaims, error) {
	var c entity.RefreshClaims
	parsedToken, err := jwt.ParseWithClaims(
		tokenStr, &c,
		func(t *jwt.Token) (interface{}, error) { return s.refreshSecret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	if c.Type != "refresh" {
		return nil, errors.New("wrong token type")
	}
	return &c, nil
}
