package port

import "ads/userservice/domain/entity"

type TokenRepository interface {
	ParseAccessToken(tokenStr string) (*entity.AccessClaims, error)
}
