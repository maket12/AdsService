package admin_uc

import (
	"AdsService/userservice/app/dto/admin_dto"
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type AdminGetProfilesUC struct {
	Users    port.UserRepository
	Profiles port.ProfileRepository
}

func (uc *AdminGetProfilesUC) Execute(in admin_dto.AdminGetProfilesListDTO) (profile_dto.ProfilesResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return profile_dto.ProfilesResponseDTO{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return profile_dto.ProfilesResponseDTO{}, uc_errors.ErrNotAdmin
	}

	profiles, err := uc.Profiles.GetAllProfiles()
	if err != nil {
		return profile_dto.ProfilesResponseDTO{}, uc_errors.ErrGetProfiles
	}
	if profiles == nil {
		return profile_dto.ProfilesResponseDTO{}, nil
	}
	return mappers.MapIntoProfilesDTO(profiles), nil
}
