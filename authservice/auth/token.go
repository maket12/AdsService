package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

type AccessClaims struct {
	Type   string `json:"type"`
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func getAccessSecret() []byte {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		panic("JWT_ACCESS_SECRET is not set in environment")
	}
	return []byte(secret)
}

func getRefreshSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		panic("JWT_REFRESH_SECRET is not set in environment")
	}
	return []byte(secret)
}

func GenerateAccessToken(userID uint64, email, role string) (string, error) {
	claims := &AccessClaims{
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getAccessSecret())
}

func GenerateRefreshToken(userID uint64) (string, error) {
	claims := RefreshClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getRefreshSecret())
}

func ParseAccessToken(tokenStr string) (*AccessClaims, error) {
	var c AccessClaims
	parsedToken, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {
		return getAccessSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithLeeway(30*time.Second))

	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	if c.Type != "access" {
		return nil, errors.New("wrong token type")
	}
	return &c, nil
}

func ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	var c RefreshClaims
	parsedToken, err := jwt.ParseWithClaims(
		tokenStr, &c,
		func(t *jwt.Token) (interface{}, error) { return getRefreshSecret(), nil },
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
