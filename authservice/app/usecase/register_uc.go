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

type RegisterUC struct {
	Users    port.UserRepository
	Sessions port.SessionRepository
	Tokens   port.TokenRepository
	Profiles port.ProfileRepository
}

func (uc *RegisterUC) Execute(ctx context.Context, in dto.Register) (dto.AuthResponse, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)

	user := &entity.User{
		Email:    in.Email,
		Password: string(hashed),
		Role:     "user",
	}
	if err := uc.Users.AddUser(ctx, user); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrAddUser
	}

	if _, err := uc.Profiles.AddProfile(ctx, user.ID, "undefined", "undefined"); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrAddProfile
	}

	access, err := uc.Tokens.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}
	refresh, err := uc.Tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	rc, err := uc.Tokens.ParseRefreshToken(refresh)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenParse
	}

	sess := &entity.Session{
		UserID:    user.ID,
		JTI:       rc.ID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: rc.ExpiresAt.Time.UTC(),
	}
	if err := uc.Sessions.CreateSession(ctx, sess); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrSessionSave
	}

	return dto.AuthResponse{AccessToken: access, RefreshToken: refresh}, nil
}
