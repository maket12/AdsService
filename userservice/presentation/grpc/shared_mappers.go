package grpc

import (
	"ads/userservice/app/dto"
	"ads/userservice/presentation/grpc/pb"
)

func MapProfileResponseDTOToPB(resp dto.ProfileResponse) *pb.Profile {
	return &pb.Profile{
		UserId: resp.UserID,
		Name:   resp.Name,
		Phone:  resp.Phone,
	}
}
