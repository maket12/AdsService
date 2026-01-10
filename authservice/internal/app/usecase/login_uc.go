package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/app/utils"
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/domain/port"
	"ads/authservice/internal/pkg/errs"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type LoginUC struct {
	Account        port.AccountRepository
	AccountRole    port.AccountRoleRepository
	RefreshSession port.RefreshSessionRepository
	PasswordHasher port.PasswordHasher
	TokenGenerator port.TokenGenerator

	refreshSessionTTL time.Duration
}

func NewLoginUC(
	account port.AccountRepository,
	accountRole port.AccountRoleRepository,
	refreshSession port.RefreshSessionRepository,
	passwordHasher port.PasswordHasher,
	tokenGenerator port.TokenGenerator,
	refreshSessionTTL time.Duration,
) *LoginUC {
	return &LoginUC{
		Account:           account,
		AccountRole:       accountRole,
		RefreshSession:    refreshSession,
		PasswordHasher:    passwordHasher,
		TokenGenerator:    tokenGenerator,
		refreshSessionTTL: refreshSessionTTL,
	}
}

func (uc *LoginUC) Execute(ctx context.Context, in dto.Login) (dto.LoginResponse, error) {
	// Find account
	account, err := uc.Account.GetByEmail(ctx, in.Email)

	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.LoginResponse{}, uc_errors.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetAccountByEmailDB, err,
		)
	}

	if !uc.PasswordHasher.Compare(account.PasswordHash(), in.Password) {
		return dto.LoginResponse{}, uc_errors.ErrInvalidCredentials
	}

	// Account validation
	if ok := account.CanLogin(); !ok {
		return dto.LoginResponse{}, uc_errors.ErrCannotLogin
	}

	// Update Account
	account.MarkLogin()
	if err := uc.Account.MarkLogin(ctx, account); err != nil {
		return dto.LoginResponse{}, uc_errors.Wrap(
			uc_errors.ErrUpdateAccountDB, err,
		)
	}

	// Find an account role
	accRole, err := uc.AccountRole.Get(ctx, account.ID())
	if err != nil {
		return dto.LoginResponse{}, uc_errors.Wrap(uc_errors.ErrGetAccountRoleDB, err)
	}

	// Generate tokens
	accessToken, err := uc.TokenGenerator.GenerateAccessToken(
		ctx, account.ID(), accRole.Role().String(),
	)
	if err != nil {
		return dto.LoginResponse{}, uc_errors.Wrap(
			uc_errors.ErrGenerateAccessToken, err,
		)
	}

	var sessionID = uuid.New()
	refreshToken, err := uc.TokenGenerator.GenerateRefreshToken(
		ctx, account.ID(), sessionID,
	)
	if err != nil {
		return dto.LoginResponse{}, uc_errors.Wrap(
			uc_errors.ErrGenerateRefreshToken, err,
		)
	}

	hashedRefreshToken := utils.HashToken(refreshToken)

	// Create refresh session
	refreshSession, err := model.NewRefreshSession(
		sessionID, account.ID(), hashedRefreshToken, nil,
		in.IP, in.UserAgent, uc.refreshSessionTTL,
	)
	if err != nil {
		return dto.LoginResponse{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}

	if err := uc.RefreshSession.Create(ctx, refreshSession); err != nil {
		return dto.LoginResponse{}, uc_errors.Wrap(
			uc_errors.ErrCreateRefreshSessionDB, err,
		)
	}

	// Output
	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
