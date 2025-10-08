package mappers

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/domain/entity"
)

func MapIntoGetUserDTO(user *entity.User) dto.GetUserResponse {
	return dto.GetUserResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
}
