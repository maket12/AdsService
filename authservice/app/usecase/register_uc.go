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

	u, err := entity.NewUser(0, in.Email, string(hashed), "user")
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if err = uc.Users.AddUser(ctx, u); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrAddUser
	}

	if _, err = uc.Profiles.AddProfile(ctx, u.GetID(), "undefined", "undefined"); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrAddProfile
	}

	access, err := uc.Tokens.GenerateAccessToken(ctx, u.GetID(), u.GetEmail(), u.GetRole())
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}
	refresh, err := uc.Tokens.GenerateRefreshToken(ctx, u.GetID())
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	rc, err := uc.Tokens.ParseRefreshToken(ctx, refresh)
	if err != nil {
		return dto.AuthResponse{}, uc_errors.ErrTokenParse
	}

	sess, err := entity.NewSession(0, entity.UserID(u.GetID()), rc.ID, time.Now().UTC(), rc.ExpiresAt.Time.UTC())
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if err := uc.Sessions.InsertSession(ctx, sess); err != nil {
		return dto.AuthResponse{}, uc_errors.ErrSessionSave
	}

	return dto.AuthResponse{AccessToken: access, RefreshToken: refresh}, nil
}
