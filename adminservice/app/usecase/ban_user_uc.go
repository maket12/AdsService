package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
)

type BanUserUC struct {
	Users port.UserRepository
}

func (uc *BanUserUC) Execute(in dto.BanUser) (dto.BanUserResponse, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.BanUserResponse{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.BanUserResponse{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.BanUser(in.RequestedUserID); err != nil {
		return dto.BanUserResponse{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrBanUser
	}
	return dto.BanUserResponse{
		UserID: in.RequestedUserID,
		Banned: true,
	}, nil
}
