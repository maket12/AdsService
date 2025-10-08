package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
)

type UnbanUserUC struct {
	Users port.UserRepository
}

func (uc *UnbanUserUC) Execute(in dto.UnbanUser) (dto.UnbanUserResponse, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.UnbanUserResponse{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.UnbanUserResponse{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.UnbanUser(in.RequestedUserID); err != nil {
		return dto.UnbanUserResponse{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrUnbanUser
	}
	return dto.UnbanUserResponse{
		UserID:   in.RequestedUserID,
		Unbanned: true,
	}, nil
}
