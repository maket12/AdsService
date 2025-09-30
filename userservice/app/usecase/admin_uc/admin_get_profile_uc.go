package admin_uc

import (
	"AdsService/userservice/app/dto/admin_dto"
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type AdminGetProfileUC struct {
	Users    port.UserRepository
	Profiles port.ProfileRepository
}

func (uc *AdminGetProfileUC) Execute(in admin_dto.AdminGetProfileDTO) (profile_dto.ProfileResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrNotAdmin
	}

	profile, err := uc.Profiles.GetProfile(in.RequestedUserID)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrGetProfile
	}
	if profile == nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrProfileNotFound
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
