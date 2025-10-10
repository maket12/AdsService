package ports

import "ads/authservice/domain/entity"

type TokenRepository interface {
	GenerateAccessToken(userID uint64, email, role string) (string, error)
	GenerateRefreshToken(userID uint64) (string, error)
	ParseAccessToken(tokenStr string) (*entity.AccessClaims, error)
	ParseRefreshToken(tokenStr string) (*entity.RefreshClaims, error)
}
