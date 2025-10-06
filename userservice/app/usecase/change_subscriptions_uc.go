package usecase

import (
	"AdsService/userservice/app/dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type ChangeSubscriptionsUC struct {
	Profiles port.ProfileRepository
}

func (uc *ChangeSubscriptionsUC) Execute(in dto.ChangeSubscriptionsDTO) (dto.ProfileResponseDTO, error) {
	profile, err := uc.Profiles.UpdateProfileSubscriptions(in.UserID, in.Subscriptions)
	if err != nil {
		return dto.ProfileResponseDTO{}, uc_errors.ErrChangeSubscriptions
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
