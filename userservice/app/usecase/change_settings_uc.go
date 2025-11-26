package usecase

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/mappers"
	"ads/userservice/app/uc_errors"
	"ads/userservice/domain/entity"
	"ads/userservice/domain/port"
	"context"
)

type ChangeSettingsUC struct {
	Profiles port.ProfileRepository
}

func (uc *ChangeSettingsUC) Execute(ctx context.Context, in dto.ChangeSettings) (dto.ProfileResponse, error) {
	var profile *entity.Profile
	var err error

	if in.NotificationsEnabled {
		profile, err = uc.Profiles.EnableNotifications(ctx, in.UserID)
	} else {
		profile, err = uc.Profiles.DisableNotifications(ctx, in.UserID)
	}

	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrChangeSettings
	}

	return mappers.MapIntoProfileDTO(profile), nil
}
