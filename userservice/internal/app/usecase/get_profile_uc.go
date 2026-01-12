package usecase

import (
	"ads/pkg/errs"
	"ads/userservice/internal/app/dto"
	"ads/userservice/internal/app/uc_errors"
	"ads/userservice/internal/domain/port"
	"context"
	"errors"
)

type GetProfileUC struct {
	Profile port.ProfileRepository
}

func NewGetProfileUC(profile port.ProfileRepository) *GetProfileUC {
	return &GetProfileUC{Profile: profile}
}

func (uc *GetProfileUC) Execute(ctx context.Context, in dto.GetProfile) (dto.GetProfileOutput, error) {
	profile, err := uc.Profile.Get(ctx, in.AccountID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.GetProfileOutput{}, uc_errors.ErrInvalidAccountID
		}
		return dto.GetProfileOutput{},
			uc_errors.Wrap(uc_errors.ErrGetProfileDB, err)
	}
	return dto.GetProfileOutput{
		AccountID: profile.AccountID(),
		FirstName: profile.FirstName(),
		LastName:  profile.LastName(),
		Phone:     profile.Phone(),
		AvatarURl: profile.AvatarURL(),
		Bio:       profile.Bio(),
		UpdatedAt: profile.UpdatedAt(),
	}, nil
}
