package port

import "ads/userservice/domain/entity"

type TokenRepository interface {
	ParseAccessToken(token string) (*entity.AccessClaims, error)
}
