package grpc

import (
	"ads/userservice/internal/app/dto"
	"ads/userservice/internal/generated/user_v1"
)

func MapGetProfilePbToDTO(req *user_v1.GetProfileRequest) dto.GetProfile {
	//return dto.GetProfileOutput{
	//	AccountID: req.,
	//	FirstName: nil,
	//	LastName:  nil,
	//	Phone:     nil,
	//	AvatarURl: nil,
	//	Bio:       nil,
	//	UpdatedAt: time.Time{},
	//}
	return dto.GetProfile{}
}
