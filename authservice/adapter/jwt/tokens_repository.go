package jwt

import (
	"ads/internal/core/ports"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"ads/authservice/domain/entity"
)

type TokenRepository struct {
	accessSecret  []byte
	refreshSecret []byte
	logger        *slog.Logger
}

func NewTokenRepository(accessSecret, refreshSecret string, log *slog.Logger) ports.TokenRepository {
	return &TokenRepository{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		logger:        log,
	}
}

func (s *TokenRepository) GenerateAccessToken(ctx context.Context, userID uint64, email, role string) (string, error) {
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

	s.logger.InfoContext(ctx, "Generated access token for user[id=%v]", userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessSecret)
}

func (s *TokenRepository) GenerateRefreshToken(ctx context.Context, userID uint64) (string, error) {
	claims := entity.RefreshClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	s.logger.InfoContext(ctx, "Generated refresh token for user[id=%v]", userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecret)
}

func (s *TokenRepository) ParseAccessToken(ctx context.Context, token string) (*entity.AccessClaims, error) {
	var c entity.AccessClaims
	parsedToken, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
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

func (s *TokenRepository) ParseRefreshToken(ctx context.Context, token string) (*entity.RefreshClaims, error) {
	var c entity.RefreshClaims
	parsedToken, err := jwt.ParseWithClaims(
		token, &c,
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
