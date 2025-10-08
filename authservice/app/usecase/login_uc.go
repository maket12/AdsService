package usecase

import (
	"ads/authservice/app/dto"
	"ads/authservice/app/uc_errors"
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginUC struct {
	Users    port.UserRepository
	Sessions port.SessionRepository
	Tokens   port.TokenRepository
}

func (uc *LoginUC) Execute(ctx context.Context, in dto.Login) (dto.AuthResponse, error) {
	user, err := uc.Users.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrGetUser
	}
	if user == nil {
		return dto.AuthResponse{}, uc_errors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrInvalidCredential
	}

	access, err := uc.Tokens.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	refresh, err := uc.Tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	rClaims, err := uc.Tokens.ParseRefreshToken(refresh)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenParse
	}

	if err := uc.Sessions.CreateSession(ctx, &entity.Session{
		UserID:    user.ID,
		JTI:       rClaims.ID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: rClaims.ExpiresAt.Time.UTC(),
	}); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrSessionSave
	}

	return dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
