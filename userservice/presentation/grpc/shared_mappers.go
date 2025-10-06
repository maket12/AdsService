package grpc

import (
	"AdsService/userservice/app/dto"
	"AdsService/userservice/presentation/grpc/pb"
)

func MapProfileResponseDTOToPB(resp dto.ProfileResponseDTO) *pb.Profile {
	return &pb.Profile{
		UserId: resp.UserID,
		Name:   resp.Name,
		Phone:  resp.Phone,
	}
}
