package usecase

import (
	"AdsService/authservice/app/dto"
	"AdsService/authservice/app/uc_errors"
	"AdsService/authservice/domain/entity"
	"AdsService/authservice/domain/port"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginUC struct {
	Users    port.UserRepository
	Sessions port.SessionRepository
	Tokens   port.TokenRepository
}

func (uc *LoginUC) Execute(in dto.LoginDTO) (dto.AuthResponseDTO, error) {
	user, err := uc.Users.GetUserByEmail(in.Email)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrGetUser
	}
	if user == nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrInvalidCredential
	}

	access, err := uc.Tokens.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrTokenIssue
	}

	refresh, err := uc.Tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrTokenIssue
	}

	rClaims, err := uc.Tokens.ParseRefreshToken(refresh)
	if err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrTokenParse
	}

	if err := uc.Sessions.InsertSession(&entity.Session{
		UserID:    user.ID,
		JTI:       rClaims.ID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: rClaims.ExpiresAt.Time.UTC(),
	}); err != nil {
		return dto.AuthResponseDTO{}, uc_errors.ErrSessionSave
	}

	return dto.AuthResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
