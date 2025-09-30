package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	
	"AdsService/userservice/domain/entity"
	"AdsService/userservice/domain/port"
)

type TokenRepository struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewTokenRepository() port.TokenRepository {
	access := os.Getenv("JWT_ACCESS_SECRET")
	refresh := os.Getenv("JWT_REFRESH_SECRET")
	if access == "" || refresh == "" {
		panic("JWT secrets are not set in env")
	}
	return &TokenRepository{
		accessSecret:  []byte(access),
		refreshSecret: []byte(refresh),
	}
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
