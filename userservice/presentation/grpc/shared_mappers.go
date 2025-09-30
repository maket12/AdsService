package grpc

import (
	"AdsService/userservice/app/dto/admin_dto"
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/dto/user_dto"
	"AdsService/userservice/presentation/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAssignRoleResponseDTOToPB(resp user_dto.AssignRoleResponseDTO) *pb.AssignRoleResponse {
	return &pb.AssignRoleResponse{
		UserId:   resp.UserID,
		Assigned: resp.Assigned,
	}
}

func MapGetUserResponseDTOToPB(resp user_dto.GetUserResponseDTO) *pb.GetUserResponse {
	return &pb.GetUserResponse{
		UserId: resp.UserID,
		Email:  resp.Email,
		Role:   resp.Role,
	}
}

func MapProfileResponseDTOToPB(resp profile_dto.ProfileResponseDTO) *pb.Profile {
	return &pb.Profile{
		UserId: resp.UserID,
		Name:   resp.Name,
		Phone:  resp.Phone,
	}
}

func MapBanUserResponseDTOToPB(resp admin_dto.BanUserResponseDTO) *pb.AdminBanUserResponse {
	return &pb.AdminBanUserResponse{
		Banned: resp.Banned,
	}
}

func MapUnbanUserResponseDTOToPB(resp admin_dto.UnbanUserResponseDTO) *pb.AdminUnbanUserResponse {
	return &pb.AdminUnbanUserResponse{
		Unbanned: resp.Unbanned,
	}
}

func MapGetProfileResponseDTOToPB(resp profile_dto.ProfileResponseDTO) *pb.Profile {
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

func MapGetProfilesResponseDTOToPB(resp profile_dto.ProfilesResponseDTO) *pb.AdminGetProfilesListResponse {
	profiles := make([]*pb.Profile, 0, len(resp.Profiles))
	for _, p := range resp.Profiles {
		profiles = append(profiles, MapGetProfileResponseDTOToPB(p))
	}
	return &pb.AdminGetProfilesListResponse{Profiles: profiles}
}
