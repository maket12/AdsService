package resolvers

import (
	authpb "AdsService/authservice/presentation/grpc/pb"
	userpb "AdsService/userservice/presentation/grpc/pb"
	"fmt"
	"strconv"
	"time"
)

func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func MapPbProfileToUserProfile(profile *userpb.Profile) *UserProfile {
	if profile == nil {
		return nil
	}

	var updatedAt *time.Time
	if profile.UpdatedAt != nil {
		t := profile.UpdatedAt.AsTime()
		updatedAt = &t
	}

	return &UserProfile{
		UserID:               fmt.Sprint(profile.UserId),
		Name:                 toPtr(profile.Name),
		Phone:                toPtr(profile.Phone),
		PhotoID:              toPtr(profile.PhotoId),
		NotificationsEnabled: profile.NotificationsEnabled,
		Subscriptions:        profile.Subscriptions,
		UpdatedAt:            updatedAt,
	}
}

func MapPbGetUserResponseToUser(user *userpb.GetUserResponse) *User {
	return &User{
		UserID: strconv.FormatUint(user.UserId, 10),
		Email:  user.Email,
		Role:   user.Role,
	}
}

func MapPbProfilesToUserProfiles(resp *userpb.AdminGetProfilesListResponse) *ProfilesList {
	if resp == nil {
		return &ProfilesList{Profiles: []*UserProfile{}}
	}

	profiles := make([]*UserProfile, 0, len(resp.Profiles))
	for _, p := range resp.Profiles {
		profiles = append(profiles, MapPbProfileToUserProfile(p))
	}

	return &ProfilesList{Profiles: profiles}
}

func MapPbBanResponseToBanResult(resp *userpb.AdminBanUserResponse) *AdminBanResult {
	return &AdminBanResult{Banned: resp.Banned}
}

func MapPbUnbanResponseToUnbanResult(resp *userpb.AdminUnbanUserResponse) *AdminUnbanResult {
	return &AdminUnbanResult{Unbanned: resp.Unbanned}
}

func MapPbAssignRoleResponseToAssignRoleResult(resp *userpb.AssignRoleResponse) *AssignRoleResult {
	return &AssignRoleResult{
		UserID:   strconv.FormatUint(resp.UserId, 10),
		Assigned: resp.Assigned,
	}
}

func MapPbAuthResponseToAuthPayload(resp *authpb.AuthResponse) *AuthPayload {
	return &AuthPayload{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
}
