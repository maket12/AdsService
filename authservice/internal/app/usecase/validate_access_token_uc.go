package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type ValidateAccessTokenUC struct {
	account        port.AccountRepository
	tokenGenerator port.TokenGenerator
}

func NewValidateAccessTokenUC(
	account port.AccountRepository,
	tokenGenerator port.TokenGenerator,
) *ValidateAccessTokenUC {
	return &ValidateAccessTokenUC{
		account:        account,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *ValidateAccessTokenUC) Execute(ctx context.Context, in dto.ValidateAccessToken) (dto.ValidateAccessTokenResponse, error) {
	// Parse access token
	accountID, role, err := uc.tokenGenerator.ValidateAccessToken(
		ctx, in.AccessToken,
	)
	if err != nil {
		return dto.ValidateAccessTokenResponse{}, uc_errors.ErrInvalidAccessToken
	}

	// Get account and check if it is not active
	account, err := uc.account.GetByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.ValidateAccessTokenResponse{}, uc_errors.ErrInvalidAccessToken
		}
		return dto.ValidateAccessTokenResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetAccountByIDDB, err,
		)
	}
	if !account.CanLogin() {
		return dto.ValidateAccessTokenResponse{}, uc_errors.ErrCannotLogin
	}

	// Output
	return dto.ValidateAccessTokenResponse{
		AccountID: accountID,
		Role:      role,
	}, nil
}
