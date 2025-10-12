package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/mappers"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
)

type GetProfileUC struct {
	Users    port.UserRepository
	Profiles port.ProfileRepository
}

func (uc *GetProfileUC) Execute(in dto.GetProfile) (dto.ProfileResponse, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.ProfileResponse{}, uc_errors.ErrNotAdmin
	}

	profile, err := uc.Profiles.GetProfile(in.RequestedUserID)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrGetProfile
	}
	if profile == nil {
		return dto.ProfileResponse{}, uc_errors.ErrProfileNotFound
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
