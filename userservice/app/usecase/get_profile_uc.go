package usecase

import (
	"AdsService/userservice/app/dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type GetProfileUC struct {
	Profiles port.ProfileRepository
}

func (uc *GetProfileUC) Execute(in dto.GetProfileDTO) (dto.ProfileResponseDTO, error) {
	profile, err := uc.Profiles.GetProfile(in.UserID)
	if err != nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrGetProfile
	}
	if profile == nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrProfileNotFound
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
