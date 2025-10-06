package usecase

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/domain/port"
)

type BanUserUC struct {
	Users port.UserRepository
}

func (uc *BanUserUC) Execute(in dto.BanUserDTO) (dto.BanUserResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return dto.BanUserResponseDTO{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.BanUserResponseDTO{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.BanUser(in.RequestedUserID); err != nil {
		return dto.BanUserResponseDTO{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrBanUser
	}
	return dto.BanUserResponseDTO{
		UserID: in.RequestedUserID,
		Banned: true,
	}, nil
}
