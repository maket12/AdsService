package profile_uc

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type GetProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *GetProfileUC) Execute(in profile_dto.GetProfileDTO) (profile_dto.ProfileResponseDTO, error) {
	profile, err := uc.Profiles.GetProfile(in.UserID)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrGetProfile
	}
	if profile == nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrProfileNotFound
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
