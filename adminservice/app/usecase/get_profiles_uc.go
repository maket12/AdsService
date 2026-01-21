package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/mappers"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
	"context"
)

type GetProfilesUC struct {
	Users    port.UserRepository
	Profiles port.ProfileRepository
}

func (uc *GetProfilesUC) Execute(ctx context.Context, in dto.GetProfilesList) (dto.ProfilesResponse, error) {
	role, err := uc.Users.GetUserRole(ctx, in.UserID)
	if err != nil {
		return dto.ProfilesResponse{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.ProfilesResponse{}, uc_errors.ErrNotAdmin
	}

	profiles, err := uc.Profiles.GetAllProfiles(ctx, in.Limit, in.Offset)
	if err != nil {
		return dto.ProfilesResponse{}, uc_errors.ErrGetProfiles
	}
	if profiles == nil {
		return dto.ProfilesResponse{}, nil
	}
	return mappers.MapIntoProfilesDTO(profiles), nil
}
