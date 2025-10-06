package usecase

import (
	"AdsService/userservice/app/dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type UpdateProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *UpdateProfileUC) Execute(in dto.UpdateProfileDTO) (dto.ProfileResponseDTO, error) {
	_, err := uc.Profiles.UpdateProfileName(in.UserID, in.Name)
	if err != nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrUpdateProfile
	}
	profile, err := uc.Profiles.UpdateProfilePhone(in.UserID, in.Phone)
	if err != nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrUpdateProfile
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
