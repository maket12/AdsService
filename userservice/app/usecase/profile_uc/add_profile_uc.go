package profile_uc

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type AddProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *AddProfileUC) Execute(in profile_dto.AddProfileDTO) (profile_dto.ProfileResponseDTO, error) {
	existing, err := uc.Profiles.GetProfile(in.UserID)
	if err == nil {
		return mappers.MapIntoProfileDTO(existing), nil
	}

	profile, err := uc.Profiles.AddProfile(in.UserID, in.Name, in.Phone)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrAddProfile
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
