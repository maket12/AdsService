package usecase

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/mappers"
	"ads/userservice/app/uc_errors"
	"ads/userservice/domain/port"
	"context"
)

type GetProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *GetProfileUC) Execute(ctx context.Context, in dto.GetProfile) (dto.ProfileResponse, error) {
	profile, err := uc.Profiles.GetProfile(ctx, in.UserID)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrGetProfile
	}
	if profile == nil {
		return dto.ProfileResponse{}, uc_errors.ErrProfileNotFound
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
