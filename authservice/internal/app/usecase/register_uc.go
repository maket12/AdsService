package usecase

import (
	dto2 "ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	entity2 "ads/authservice/internal/domain/entity"
	port2 "ads/authservice/internal/domain/port"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUC struct {
	Users    port2.UserRepository
	Sessions port2.SessionRepository
	Tokens   port2.TokenRepository
	Profiles port2.ProfileRepository
}

func (uc *RegisterUC) Execute(ctx context.Context, in dto2.Register) (dto2.AuthResponse, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)

	u, err := entity2.NewUser(0, in.Email, string(hashed), "user")
	if err != nil {
		return dto2.AuthResponse{}, err
	}
	if err = uc.Users.AddUser(ctx, u); err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrAddUser
	}

	if _, err = uc.Profiles.AddProfile(ctx, u.GetID(), "undefined", "undefined"); err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrAddProfile
	}

	access, err := uc.Tokens.GenerateAccessToken(ctx, u.GetID(), u.GetEmail(), u.GetRole())
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrTokenIssue
	}
	refresh, err := uc.Tokens.GenerateRefreshToken(ctx, u.GetID())
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrTokenIssue
	}

	rc, err := uc.Tokens.ParseRefreshToken(ctx, refresh)
	if err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrTokenParse
	}

	sess, err := entity2.NewSession(0, entity2.UserID(u.GetID()), rc.ID, time.Now().UTC(), rc.ExpiresAt.Time.UTC())
	if err != nil {
		return dto2.AuthResponse{}, err
	}
	if err := uc.Sessions.InsertSession(ctx, sess); err != nil {
		return dto2.AuthResponse{}, uc_errors.ErrSessionSave
	}

	return dto2.AuthResponse{AccessToken: access, RefreshToken: refresh}, nil
}
