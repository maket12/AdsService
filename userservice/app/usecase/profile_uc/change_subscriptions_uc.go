package profile_uc

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type ChangeSubscriptionsUC struct {
	Profiles port.ProfileRepository
}

func (uc *ChangeSubscriptionsUC) Execute(in profile_dto.ChangeSubscriptionsDTO) (profile_dto.ProfileResponseDTO, error) {
	profile, err := uc.Profiles.UpdateProfileSubscriptions(in.UserID, in.Subscriptions)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrChangeSubscriptions
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
