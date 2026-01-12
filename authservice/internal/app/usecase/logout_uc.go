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
	RefreshSession port.RefreshSessionRepository
	TokenGenerator port.TokenGenerator
}

func NewLogoutUC(
	refreshSession port.RefreshSessionRepository,
	tokenGenerator port.TokenGenerator,
) *LogoutUC {
	return &LogoutUC{
		RefreshSession: refreshSession,
		TokenGenerator: tokenGenerator,
	}
}

func (uc *LogoutUC) Execute(ctx context.Context, in dto.Logout) (dto.LogoutResponse, error) {
	// Find session
	_, oldSessionID, err := uc.TokenGenerator.ValidateRefreshToken(
		ctx, in.RefreshToken,
	)
	if err != nil {
		return dto.LogoutResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	session, err := uc.RefreshSession.GetByID(ctx, oldSessionID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.LogoutResponse{}, uc_errors.ErrInvalidRefreshToken
		}
		return dto.LogoutResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetRefreshSessionByIDDB, err,
		)
	}

	// Validate and revoke
	if !session.IsActive() {
		return dto.LogoutResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	if utils.HashToken(in.RefreshToken) != session.RefreshTokenHash() {
		return dto.LogoutResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	var reason = "logout"
	if err := session.Revoke(&reason); err != nil {
		return dto.LogoutResponse{}, uc_errors.ErrCannotRevoke
	}

	if err := uc.RefreshSession.Revoke(ctx, session); err != nil {
		return dto.LogoutResponse{}, uc_errors.Wrap(
			uc_errors.ErrRevokeRefreshSessionDB, err,
		)
	}

	return dto.LogoutResponse{Logout: true}, nil
}
