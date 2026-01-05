package usecase

import (
	dto2 "ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/domain/entity"
	port2 "ads/authservice/internal/domain/port"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginUC struct {
	Users    port2.UserRepository
	Sessions port2.SessionRepository
	Tokens   port2.TokenRepository
}

func (uc *LoginUC) Execute(ctx context.Context, in dto2.Login) (dto2.AuthResponse, error) {
	user, err := uc.Users.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrGetUser
	}
	if user == nil {
		return dto2.AuthResponse{}, uc_errors.ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrInvalidCredential
	}

	access, err := uc.Tokens.GenerateAccessToken(ctx, user.GetID(), user.GetEmail(), user.GetRole())
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	refresh, err := uc.Tokens.GenerateRefreshToken(ctx, user.GetID())
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	rClaims, err := uc.Tokens.ParseRefreshToken(ctx, refresh)
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrTokenParse
	}

	sess, err := model.NewSession(0, model.UserID(user.GetID()), rClaims.ID, time.Now().UTC(), rClaims.ExpiresAt.Time.UTC())
	if err = uc.Sessions.InsertSession(ctx, sess); err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrSessionSave
	}

	return dto2.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
