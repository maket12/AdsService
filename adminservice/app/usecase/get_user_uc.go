package usecase

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/mappers"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/domain/port"
)

type GetUserUC struct {
	Users port.UserRepository
}

func (uc *GetUserUC) Execute(in dto.GetUserDTO) (dto.GetUserResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.AdminUserID)
	if err != nil {
		return dto.GetUserResponseDTO{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.GetUserResponseDTO{}, uc_errors.ErrNotAdmin
	}

	user, err := uc.Users.GetUserByID(in.RequestedUserID)
	if err != nil {
		return dto.GetUserResponseDTO{}, uc_errors.ErrGetUser
	}
	if user == nil {
		return dto.GetUserResponseDTO{}, uc_errors.ErrUserNotFound
	}

	return mappers.MapIntoGetUserDTO(user), nil
}
