package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/mappers"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
)

type GetUserUC struct {
	Users port.UserRepository
}

func (uc *GetUserUC) Execute(in dto.GetUser) (dto.GetUserResponse, error) {
	role, err := uc.Users.GetUserRole(in.AdminUserID)
	if err != nil {
		return dto.GetUserResponse{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.GetUserResponse{}, uc_errors.ErrNotAdmin
	}

	user, err := uc.Users.GetUserByID(in.RequestedUserID)
	if err != nil {
		return dto.GetUserResponse{}, uc_errors.ErrGetUser
	}
	if user == nil {
		return dto.GetUserResponse{}, uc_errors.ErrUserNotFound
	}

	return mappers.MapIntoGetUserDTO(user), nil
}
