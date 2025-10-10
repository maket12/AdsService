package usecase

import (
	"ads/authservice/app/dto"
	"ads/authservice/app/uc_errors"
	"ads/authservice/domain/port"
	"context"
)

type ValidateTokenUC struct {
	Tokens port.TokenRepository
}

func (uc *ValidateTokenUC) Execute(ctx context.Context, in dto.ValidateToken) (dto.ValidateTokenResponse, error) {
	claims, err := uc.Tokens.ParseAccessToken(ctx, in.AccessToken)
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
