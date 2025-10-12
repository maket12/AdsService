package mappers

import (
	"ads/userservice/app/dto"
	"ads/userservice/domain/entity"
)

func MapIntoProfileDTO(profile *entity.Profile) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserID:               profile.UserID,
		Name:                 profile.Name,
		Phone:                profile.Phone,
		PhotoID:              profile.PhotoID,
		NotificationsEnabled: profile.NotificationsEnabled,
		Subscriptions:        profile.Subscriptions,
		UpdatedAt:            profile.UpdatedAt,
	}
}

func MapIntoProfilesDTO(profiles []entity.Profile) dto.ProfilesResponse {
	out := make([]dto.ProfileResponse, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, MapIntoProfileDTO(&p))
	}
	return dto.ProfilesResponse{Profiles: out}
}
