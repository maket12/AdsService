package resolvers

import (
	adminpb "AdsService/adminservice/presentation/grpc/pb"
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

func MapUserPbProfileToUserProfile(profile *userpb.Profile) *UserProfile {
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

func MapAdminPbProfileToUserProfile(profile *adminpb.Profile) *UserProfile {
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

func MapPbGetUserResponseToUser(user *adminpb.GetUserResponse) *User {
	return &User{
		UserID: strconv.FormatUint(user.UserId, 10),
		Email:  user.Email,
		Role:   user.Role,
	}
}

func MapAdminPbProfilesToUserProfiles(resp *adminpb.GetProfilesListResponse) *ProfilesList {
	if resp == nil {
		return &ProfilesList{Profiles: []*UserProfile{}}
	}

	profiles := make([]*UserProfile, 0, len(resp.Profiles))
	for _, p := range resp.Profiles {
		profiles = append(profiles, MapAdminPbProfileToUserProfile(p))
	}

	return &ProfilesList{Profiles: profiles}
}

func MapPbBanResponseToBanResult(resp *adminpb.BanUserResponse) *AdminBanResult {
	return &AdminBanResult{Banned: resp.Banned}
}

func MapPbUnbanResponseToUnbanResult(resp *adminpb.UnbanUserResponse) *AdminUnbanResult {
	return &AdminUnbanResult{Unbanned: resp.Unbanned}
}

func MapPbAssignRoleResponseToAssignRoleResult(resp *adminpb.AssignRoleResponse) *AssignRoleResult {
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
