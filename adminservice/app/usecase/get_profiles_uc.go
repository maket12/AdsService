package usecase

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/mappers"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/domain/port"
)

type GetProfilesUC struct {
	Users    port.UserRepository
	Profiles port.ProfileRepository
}

func (uc *GetProfilesUC) Execute(in dto.GetProfilesListDTO) (dto.ProfilesResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.ProfilesResponseDTO{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.ProfilesResponseDTO{}, uc_errors.ErrNotAdmin
	}

	profiles, err := uc.Profiles.GetAllProfiles(in.Limit, in.Offset)
	if err != nil {
		return dto.ProfilesResponseDTO{}, uc_errors.ErrGetProfiles
	}
	if profiles == nil {
		return dto.ProfilesResponseDTO{}, nil
	}
	return mappers.MapIntoProfilesDTO(profiles), nil
}
