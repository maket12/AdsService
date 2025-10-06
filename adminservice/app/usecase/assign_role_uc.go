package usecase

import (
	"AdsService/adminservice/app/dto"
	"AdsService/adminservice/app/uc_errors"
	"AdsService/adminservice/domain/port"
)

type AssignRoleUC struct {
	Users port.UserRepository
}

func (uc *AssignRoleUC) Execute(in dto.AssignRoleDTO) (dto.AssignRoleResponseDTO, error) {
	role, err := uc.Users.GetUserRole(in.AdminUserID)
	if err != nil {
		return dto.AssignRoleResponseDTO{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.AssignRoleResponseDTO{}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.EnhanceUser(in.RequestedUserID); err != nil {
		return dto.AssignRoleResponseDTO{
			UserID:   in.RequestedUserID,
			Assigned: false,
		}, uc_errors.ErrEnhanceUser
	}
	return dto.AssignRoleResponseDTO{
		UserID:   in.RequestedUserID,
		Assigned: true,
	}, nil
}
