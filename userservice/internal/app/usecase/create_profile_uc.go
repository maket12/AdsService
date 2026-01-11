package usecase

import (
	"ads/userservice/internal/app/dto"
	"ads/userservice/internal/app/uc_errors"
	"ads/userservice/internal/domain/model"
	"ads/userservice/internal/domain/port"
	"ads/userservice/pkg/errs"
	"context"
	"errors"
)

type CreateProfileUC struct {
	Profile port.ProfileRepository
}

func NewCreateProfileUC(profile port.ProfileRepository) *CreateProfileUC {
	return &CreateProfileUC{Profile: profile}
}

func (uc *CreateProfileUC) Execute(ctx context.Context, in dto.CreateProfile) error {
	// Create profile
	profile, err := model.NewProfile(in.AccountID)
	if err != nil {
		return uc_errors.ErrInvalidAccountID
	}

	if err := uc.Profile.Create(ctx, profile); err != nil {
		if errors.Is(err, errs.ErrObjectAlreadyExists) {
			return nil
		}
		return uc_errors.ErrCreateProfileDB
	}

	return nil
}
