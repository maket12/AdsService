package usecase

import (
	"ads/pkg/errs"
	"ads/userservice/internal/app/dto"
	"ads/userservice/internal/app/uc_errors"
	"ads/userservice/internal/domain/model"
	"ads/userservice/internal/domain/port"
	"context"
	"errors"
)

type CreateProfileUC struct {
	profile port.ProfileRepository
}

func NewCreateProfileUC(profile port.ProfileRepository) *CreateProfileUC {
	return &CreateProfileUC{profile: profile}
}

func (uc *CreateProfileUC) Execute(ctx context.Context, in dto.CreateProfileInput) error {
	// Create profile
	profile, err := model.NewProfile(in.AccountID)
	if err != nil {
		return uc_errors.ErrInvalidAccountID
	}

	if err := uc.profile.Create(ctx, profile); err != nil {
		if errors.Is(err, errs.ErrObjectAlreadyExists) {
			return nil
		}
		return uc_errors.Wrap(uc_errors.ErrCreateProfileDB, err)
	}

	return nil
}
