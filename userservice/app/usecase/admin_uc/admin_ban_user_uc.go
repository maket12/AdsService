package admin_uc

import (
	"AdsService/userservice/app/dto/admin_dto"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type AdminBanUserUC struct {
	Users port.UserRepository
}

func (uc *AdminBanUserUC) Execute(in admin_dto.BanUserDTO) (admin_dto.BanUserResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return admin_dto.BanUserResponseDTO{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return admin_dto.BanUserResponseDTO{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.BanUser(in.RequestedUserID); err != nil {
		return admin_dto.BanUserResponseDTO{
			UserID: in.RequestedUserID,
			Banned: false,
		}, uc_errors.ErrBanUser
	}
	return admin_dto.BanUserResponseDTO{
		UserID: in.RequestedUserID,
		Banned: true,
	}, nil
}
