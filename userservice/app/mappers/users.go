package mappers

import (
	"AdsService/userservice/app/dto/user_dto"
	"AdsService/userservice/domain/entity"
)

func MapIntoGetUserDTO(user *entity.User) user_dto.GetUserResponseDTO {
	return user_dto.GetUserResponseDTO{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
}
