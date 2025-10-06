package usecase

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/domain/port"
)

type UnbanUserUC struct {
	Users port.UserRepository
}

func (uc *UnbanUserUC) Execute(in dto.UnbanUserDTO) (dto.UnbanUserResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.UnbanUserResponseDTO{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.UnbanUserResponseDTO{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.UnbanUser(in.RequestedUserID); err != nil {
		return dto.UnbanUserResponseDTO{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrUnbanUser
	}
	return dto.UnbanUserResponseDTO{
		UserID:   in.RequestedUserID,
		Unbanned: true,
	}, nil
}
