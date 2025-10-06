package usecase

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/mappers"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/domain/port"
)

type GetProfileUC struct {
	Users    port.UserRepository
	Profiles port.ProfileRepository
}

func (uc *GetProfileUC) Execute(in dto.GetProfileDTO) (dto.ProfileResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.ProfileResponseDTO{}, uc_errors.ErrNotAdmin
	}

	profile, err := uc.Profiles.GetProfile(in.RequestedUserID)
	if err != nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrGetProfile
	}
	if profile == nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrProfileNotFound
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
