package usecase

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/mappers"
	"ads/userservice/app/uc_errors"
	"ads/userservice/domain/port"
)

type ChangeSubscriptionsUC struct {
	Profiles port.ProfileRepository
}

func (uc *ChangeSubscriptionsUC) Execute(in dto.ChangeSubscriptions) (dto.ProfileResponse, error) {
	profile, err := uc.Profiles.UpdateProfileSubscriptions(in.UserID, in.Subscriptions)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrChangeSubscriptions
	}
	return mappers.MapIntoProfileDTO(profile), nil
}
