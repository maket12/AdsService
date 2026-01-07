package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/utils"
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/domain/port"
	"ads/authservice/pkg/errs"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshSessionUC struct {
	AccountRole    port.AccountRoleRepository
	RefreshSession port.RefreshSessionRepository
	TokenGenerator port.TokenGenerator

	refreshSessionTTL time.Duration
}

func NewRefreshSessionUC(
	accountRole port.AccountRoleRepository,
	refreshSession port.RefreshSessionRepository,
	tokenGenerator port.TokenGenerator,
	refreshSessionTTL time.Duration,
) *RefreshSessionUC {
	return &RefreshSessionUC{
		AccountRole:       accountRole,
		RefreshSession:    refreshSession,
		TokenGenerator:    tokenGenerator,
		refreshSessionTTL: refreshSessionTTL,
	}
}

func (uc *RefreshSessionUC) Execute(ctx context.Context, in dto.RefreshSession) (dto.RefreshSessionResponse, error) {
	// Find old session
	accountID, oldSessionID, err := uc.TokenGenerator.ValidateRefreshToken(
		ctx, in.OldRefreshToken,
	)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.ErrInvalidRefreshToken
	}

	oldSession, err := uc.RefreshSession.GetByID(ctx, oldSessionID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.RefreshSessionResponse{}, uc_errors.ErrInvalidRefreshToken
		}
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetRefreshSessionByIDDB, err,
		)
	}

	// Validate and revoke
	if !oldSession.IsActive() || *oldSession.IP() != *in.IP || *oldSession.UserAgent() != *in.UserAgent {
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

	if err := uc.RefreshSession.Revoke(ctx, oldSession); err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrRevokeRefreshSession, err,
		)
	}

	// Get account role
	accRole, err := uc.AccountRole.Get(ctx, accountID)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetAccountRoleDB, err,
		)
	}

	// Generate new tokens
	accessToken, err := uc.TokenGenerator.GenerateAccessToken(
		ctx, accountID, accRole.Role().String(),
	)
	if err != nil {
		return dto.RefreshSessionResponse{}, uc_errors.Wrap(
			uc_errors.ErrGenerateAccessToken, err,
		)
	}

	var sessionID = uuid.New()
	refreshToken, err := uc.TokenGenerator.GenerateRefreshToken(
		ctx, accountID, sessionID,
	)

	hashedRefreshToken := utils.HashToken(refreshToken)

	// Create new refresh session with rotation
	refreshSession, err := model.NewRefreshSession(
		sessionID, accountID, hashedRefreshToken, &oldSessionID,
		in.IP, in.UserAgent, uc.refreshSessionTTL,
	)

	if err := uc.RefreshSession.Create(ctx, refreshSession); err != nil {
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
