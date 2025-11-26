package usecase

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/mappers"
	"ads/userservice/app/uc_errors"
	"ads/userservice/domain/port"
	"context"
)

type AddProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *AddProfileUC) Execute(ctx context.Context, in dto.AddProfile) (dto.ProfileResponse, error) {
	existing, err := uc.Profiles.GetProfile(ctx, in.UserID)
	if err == nil {
		return mappers.MapIntoProfileDTO(existing), nil
	}

	profile, err := uc.Profiles.AddProfile(ctx, in.UserID, in.Name, in.Phone)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrAddProfile
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
