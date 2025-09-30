package user_uc

import (
	"AdsService/userservice/app/dto/user_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type GetUserUC struct {
	Users port.UserRepository
}

func (uc *GetUserUC) Execute(in user_dto.GetUserDTO) (user_dto.GetUserResponseDTO, error) {
	user, err := uc.Users.GetUserByID(in.UserID)
	if err != nil {
		return user_dto.GetUserResponseDTO{
			UserID: in.UserID,
			Email:  "",
			Role:   "",
		}, uc_errors.ErrGetUser
	}
	if user == nil {
		return user_dto.GetUserResponseDTO{}, uc_errors.ErrUserNotFound
	}

	return mappers.MapIntoGetUserDTO(user), nil
}
