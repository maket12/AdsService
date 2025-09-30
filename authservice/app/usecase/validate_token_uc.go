package usecase

import (
	"AdsService/authservice/app/dto"
	"AdsService/authservice/app/uc_errors"
	"AdsService/authservice/domain/port"
)

type ValidateTokenUC struct {
	Tokens port.TokenRepository
}

func (uc *ValidateTokenUC) Execute(in dto.ValidateTokenDTO) (dto.ValidateTokenResponse, error) {
	claims, err := uc.Tokens.ParseAccessToken(in.AccessToken)
	if err != nil {
		return dto.ValidateTokenResponse{
			Valid:  false,
			UserID: 0,
		}, uc_errors.ErrTokenIssue
	}

	return dto.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
	}, nil
}
