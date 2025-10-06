package mappers

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/domain/entity"
)

func MapIntoGetUserDTO(user *entity.User) dto.GetUserResponseDTO {
	return dto.GetUserResponseDTO{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
}
