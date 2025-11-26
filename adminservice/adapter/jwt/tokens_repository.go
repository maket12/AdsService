package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"ads/userservice/domain/entity"
	"ads/userservice/domain/port"
)

type TokenRepository struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewTokenRepository(accessToken, refreshToken string) port.TokenRepository {
	return &TokenRepository{
		accessSecret:  []byte(accessToken),
		refreshSecret: []byte(refreshToken),
	}
}

func (s *TokenRepository) ParseAccessToken(token string) (*entity.AccessClaims, error) {
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
