package mappers

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/domain/entity"
)

func MapIntoProfileDTO(profile *entity.Profile) dto.ProfileResponseDTO {
	return dto.ProfileResponseDTO{
		UserID:               profile.UserID,
		Name:                 profile.Name,
		Phone:                profile.Phone,
		PhotoID:              profile.PhotoID,
		NotificationsEnabled: profile.NotificationsEnabled,
		Subscriptions:        profile.Subscriptions,
		UpdatedAt:            profile.UpdatedAt,
	}
}

func MapIntoProfilesDTO(profiles []entity.Profile) dto.ProfilesResponseDTO {
	out := make([]dto.ProfileResponseDTO, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, MapIntoProfileDTO(&p))
	}
	return dto.ProfilesResponseDTO{Profiles: out}
}
