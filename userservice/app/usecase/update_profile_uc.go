package usecase

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/mappers"
	"ads/userservice/app/uc_errors"
	"ads/userservice/domain/port"
	"context"
)

type UpdateProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *UpdateProfileUC) Execute(ctx context.Context, in dto.UpdateProfile) (dto.ProfileResponse, error) {
	_, err := uc.Profiles.UpdateProfileName(ctx, in.UserID, in.Name)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrUpdateProfile
	}
	profile, err := uc.Profiles.UpdateProfilePhone(ctx, in.UserID, in.Phone)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrUpdateProfile
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
