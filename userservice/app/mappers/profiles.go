package mappers

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/domain/entity"
)

func MapIntoProfileDTO(profile *entity.Profile) profile_dto.ProfileResponseDTO {
	return profile_dto.ProfileResponseDTO{
		UserID:               profile.UserID,
		Name:                 profile.Name,
		Phone:                profile.Phone,
		PhotoID:              profile.PhotoID,
		NotificationsEnabled: profile.NotificationsEnabled,
		Subscriptions:        profile.Subscriptions,
		UpdatedAt:            profile.UpdatedAt,
	}
}

func MapIntoProfilesDTO(profiles []entity.Profile) profile_dto.ProfilesResponseDTO {
	out := make([]profile_dto.ProfileResponseDTO, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, MapIntoProfileDTO(&p))
	}
	return profile_dto.ProfilesResponseDTO{Profiles: out}
}
