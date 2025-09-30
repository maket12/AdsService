package profile_uc

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type UpdateProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *UpdateProfileUC) Execute(in profile_dto.UpdateProfileDTO) (profile_dto.ProfileResponseDTO, error) {
	profile, err := uc.Profiles.UpdateProfileName(in.UserID, in.Name)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrUpdateProfile
	}
	profile, err = uc.Profiles.UpdateProfilePhone(in.UserID, in.Phone)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrUpdateProfile
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
