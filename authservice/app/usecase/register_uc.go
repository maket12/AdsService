package usecase

import (
	"AdsService/authservice/app/dto"
	"AdsService/authservice/app/uc_errors"
	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUC struct {
	Users    port.UserRepository
	Sessions port.SessionRepository
	Tokens   port.TokenRepository
	Profiles port.ProfileRepository
}

func (uc *RegisterUC) Execute(in dto.RegisterDTO) (dto.AuthResponseDTO, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)

	user := &entity.User{
		Email:    in.Email,
		Password: string(hashed),
		Role:     "user",
	}
	if err := uc.Users.AddUser(user); err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrAddUser
	}

	if _, err := uc.Profiles.AddProfile(user.ID, "undefined", "undefined"); err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrAddProfile
	}

	access, err := uc.Tokens.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrTokenIssue
	}
	refresh, err := uc.Tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrTokenIssue
	}

	rc, err := uc.Tokens.ParseRefreshToken(refresh)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrTokenParse
	}

	sess := &entity.Session{
		UserID:    user.ID,
		JTI:       rc.ID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: rc.ExpiresAt.Time.UTC(),
	}
	if err := uc.Sessions.InsertSession(sess); err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrSessionSave
	}

	return dto.AuthResponseDTO{AccessToken: access, RefreshToken: refresh}, nil
}
