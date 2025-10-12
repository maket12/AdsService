package resolvers

import (
	adminpb "ads/adminservice/presentation/grpc/pb"
	authpb "ads/authservice/presentation/grpc/pb"
	userpb "ads/userservice/presentation/grpc/pb"
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

func MapPbBanResponseToBanOutput(resp *adminpb.BanUserResponse) *AdminBanOutput {
	return &AdminBanOutput{Banned: resp.Banned}
}

func MapPbUnbanResponseToUnbanOutput(resp *adminpb.UnbanUserResponse) *AdminUnbanOutput {
	return &AdminUnbanOutput{Unbanned: resp.Unbanned}
}

func MapPbAssignRoleResponseToAssignRoleOutput(resp *adminpb.AssignRoleResponse) *AssignRoleOutput {
	return &AssignRoleOutput{
		UserID:   strconv.FormatUint(resp.UserId, 10),
		Assigned: resp.Assigned,
	}
}

func MapPbAuthResponseToRegisterOutput(resp *authpb.AuthResponse) *RegisterUserOutput {
	return &RegisterUserOutput{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
}

func MapPbAuthResponseToLoginOutput(resp *authpb.AuthResponse) *LoginUserOutput {
	return &LoginUserOutput{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
}
