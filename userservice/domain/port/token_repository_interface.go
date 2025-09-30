package port

import "AdsService/userservice/domain/entity"

type TokenRepository interface {
	ParseAccessToken(tokenStr string) (*entity.AccessClaims, error)
}
