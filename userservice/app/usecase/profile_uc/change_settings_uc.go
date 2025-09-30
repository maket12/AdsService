package profile_uc

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/entity"
	"AdsService/userservice/domain/port"
)

type ChangeSettingsUC struct {
	Profiles port.ProfileRepository
}

func (uc *ChangeSettingsUC) Execute(in profile_dto.ChangeSettingsDTO) (profile_dto.ProfileResponseDTO, error) {
	var profile *entity.Profile
	var err error

	if in.NotificationsEnabled {
		profile, err = uc.Profiles.EnableNotifications(in.UserID)
	} else {
		profile, err = uc.Profiles.DisableNotifications(in.UserID)
	}

	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrChangeSettings
	}

	return mappers.MapIntoProfileDTO(profile), nil
}
