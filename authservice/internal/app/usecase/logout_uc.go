package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/utils"
	"ads/authservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type LogoutUC struct {
	refreshSession port.RefreshSessionRepository
	tokenGenerator port.TokenGenerator
}

func NewLogoutUC(
	refreshSession port.RefreshSessionRepository,
	tokenGenerator port.TokenGenerator,
) *LogoutUC {
	return &LogoutUC{
		refreshSession: refreshSession,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *LogoutUC) Execute(ctx context.Context, in dto.LogoutInput) (dto.LogoutOutput, error) {
	// Find session
	_, oldSessionID, err := uc.tokenGenerator.ValidateRefreshToken(
		ctx, in.RefreshToken,
	)
	if err != nil {
		return dto.LogoutOutput{}, uc_errors.ErrInvalidRefreshToken
	}

	session, err := uc.refreshSession.GetByID(ctx, oldSessionID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.LogoutOutput{}, uc_errors.ErrInvalidRefreshToken
		}
		return dto.LogoutOutput{}, uc_errors.Wrap(
			uc_errors.ErrGetRefreshSessionByIDDB, err,
		)
	}

	// Validate and revoke
	if !session.IsActive() {
		return dto.LogoutOutput{}, uc_errors.ErrInvalidRefreshToken
	}

	if utils.HashToken(in.RefreshToken) != session.RefreshTokenHash() {
		return dto.LogoutOutput{}, uc_errors.ErrInvalidRefreshToken
	}

	var reason = "logout"
	if err := session.Revoke(&reason); err != nil {
		return dto.LogoutOutput{}, uc_errors.ErrCannotRevoke
	}

	if err := uc.refreshSession.Revoke(ctx, session); err != nil {
		return dto.LogoutOutput{}, uc_errors.Wrap(
			uc_errors.ErrRevokeRefreshSessionDB, err,
		)
	}

	return dto.LogoutOutput{Logout: true}, nil
}
