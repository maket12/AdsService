package grpc

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/presentation/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAssignRoleResponseDTOToPB(resp dto.AssignRoleResponseDTO) *pb.AssignRoleResponse {
	return &pb.AssignRoleResponse{
		UserId:   resp.UserID,
		Assigned: resp.Assigned,
	}
}

func MapGetUserResponseDTOToPB(resp dto.GetUserResponseDTO) *pb.GetUserResponse {
	return &pb.GetUserResponse{
		UserId: resp.UserID,
		Email:  resp.Email,
		Role:   resp.Role,
	}
}

func MapBanUserResponseDTOToPB(resp dto.BanUserResponseDTO) *pb.BanUserResponse {
	return &pb.BanUserResponse{
		Banned: resp.Banned,
	}
}

func MapUnbanUserResponseDTOToPB(resp dto.UnbanUserResponseDTO) *pb.UnbanUserResponse {
	return &pb.UnbanUserResponse{
		Unbanned: resp.Unbanned,
	}
}

func MapGetProfileResponseDTOToPB(resp dto.ProfileResponseDTO) *pb.Profile {
	return &pb.Profile{
		UserId:               resp.UserID,
		Name:                 resp.Name,
		Phone:                resp.Phone,
		PhotoId:              resp.PhotoID,
		NotificationsEnabled: resp.NotificationsEnabled,
		Subscriptions:        resp.Subscriptions,
		UpdatedAt:            timestamppb.New(resp.UpdatedAt),
	}
}

func MapGetProfilesResponseDTOToPB(resp dto.ProfilesResponseDTO) *pb.GetProfilesListResponse {
	profiles := make([]*pb.Profile, 0, len(resp.Profiles))
	for _, p := range resp.Profiles {
		profiles = append(profiles, MapGetProfileResponseDTOToPB(p))
	}
	return &pb.GetProfilesListResponse{Profiles: profiles}
}
