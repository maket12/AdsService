package usecase

import (
	"ads/authservice/app/dto"
	"ads/authservice/app/uc_errors"
	"ads/authservice/domain/entity"
	"ads/authservice/domain/port"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
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

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrInvalidCredential
	}

	access, err := uc.Tokens.GenerateAccessToken(ctx, user.GetID(), user.GetEmail(), user.GetRole())
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	refresh, err := uc.Tokens.GenerateRefreshToken(ctx, user.GetID())
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	rClaims, err := uc.Tokens.ParseRefreshToken(ctx, refresh)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenParse
	}

	sess, err := entity.NewSession(0, entity.UserID(user.GetID()), rClaims.ID, time.Now().UTC(), rClaims.ExpiresAt.Time.UTC())
	if err = uc.Sessions.InsertSession(ctx, sess); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrSessionSave
	}

	return dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
