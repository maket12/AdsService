package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/utils"
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshSessionUC struct {
	accountRole    port.AccountRoleRepository
	refreshSession port.RefreshSessionRepository
	tokenGenerator port.TokenGenerator

	refreshSessionTTL time.Duration
}

func NewRefreshSessionUC(
	accountRole port.AccountRoleRepository,
	refreshSession port.RefreshSessionRepository,
	tokenGenerator port.TokenGenerator,
	refreshSessionTTL time.Duration,
) *RefreshSessionUC {
	return &RefreshSessionUC{
		accountRole:       accountRole,
		refreshSession:    refreshSession,
		tokenGenerator:    tokenGenerator,
		refreshSessionTTL: refreshSessionTTL,
	}
}

func (uc *RefreshSessionUC) Execute(ctx context.Context, in dto.RefreshSession) (dto.RefreshSessionResponse, error) {
	// Find old session
	accountID, oldSessionID, err := uc.tokenGenerator.ValidateRefreshToken(
		ctx, in.OldRefreshToken,
	)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	oldSession, err := uc.refreshSession.GetByID(ctx, oldSessionID)

	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.RefreshSessionResponse{}, uc_errors.ErrInvalidRefreshToken
		}
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetRefreshSessionByIDDB, err,
		)
	}

	// Validate and revoke
	if !oldSession.IsActive() ||
		!utils.ComparePtr(oldSession.IP(), in.IP) ||
		!utils.ComparePtr(oldSession.UserAgent(), in.UserAgent) {
		return dto.RefreshSessionResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	if utils.HashToken(in.OldRefreshToken) != oldSession.RefreshTokenHash() {
		return dto.RefreshSessionResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	var reason = "token rotation"
	if err := oldSession.Revoke(&reason); err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrCannotRevoke, err,
		)
	}

	if err := uc.refreshSession.Revoke(ctx, oldSession); err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrRevokeRefreshSessionDB, err,
		)
	}

	// Get account role
	accRole, err := uc.accountRole.Get(ctx, accountID)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetAccountRoleDB, err,
		)
	}

	// Generate new tokens
	accessToken, err := uc.tokenGenerator.GenerateAccessToken(
		ctx, accountID, accRole.Role().String(),
	)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGenerateAccessToken, err,
		)
	}

	var sessionID = uuid.New()
	refreshToken, err := uc.tokenGenerator.GenerateRefreshToken(
		ctx, accountID, sessionID,
	)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGenerateRefreshToken, err,
		)
	}

	hashedRefreshToken := utils.HashToken(refreshToken)

	// Create new refresh session with rotation
	refreshSession, err := model.NewRefreshSession(
		sessionID, accountID, hashedRefreshToken, &oldSessionID,
		in.IP, in.UserAgent, uc.refreshSessionTTL,
	)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}

	if err := uc.refreshSession.Create(ctx, refreshSession); err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrCreateRefreshSessionDB, err,
		)
	}

	// Output
	return dto.RefreshSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
