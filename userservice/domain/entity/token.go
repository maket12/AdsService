package entity

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	Type   string `json:"type"`
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
