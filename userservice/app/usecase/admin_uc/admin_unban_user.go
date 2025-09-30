package admin_uc

import (
	"AdsService/userservice/app/dto/admin_dto"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type AdminUnbanUserUC struct {
	Users port.UserRepository
}

func (uc *AdminUnbanUserUC) Execute(in admin_dto.UnbanUserDTO) (admin_dto.UnbanUserResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.UserID)
	if err != nil {
		return admin_dto.UnbanUserResponseDTO{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return admin_dto.UnbanUserResponseDTO{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.BanUser(in.RequestedUserID); err != nil {
		return admin_dto.UnbanUserResponseDTO{
			UserID:   in.RequestedUserID,
			Unbanned: false,
		}, uc_errors.ErrUnbanUser
	}
	return admin_dto.UnbanUserResponseDTO{
		UserID:   in.RequestedUserID,
		Unbanned: true,
	}, nil
}
