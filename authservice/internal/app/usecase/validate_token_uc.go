package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/domain/port"
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
