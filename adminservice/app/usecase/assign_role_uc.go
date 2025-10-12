package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
)

type AssignRoleUC struct {
	Users port.UserRepository
}

func (uc *AssignRoleUC) Execute(in dto.AssignRole) (dto.AssignRoleResponse, error) {
	role, err := uc.Users.GetUserRole(in.AdminUserID)
	if err != nil {
		return dto.AssignRoleResponse{}, uc_errors.ErrGetUserRole
	}
	if role != "admin" {
		return dto.AssignRoleResponse{}, uc_errors.ErrNotAdmin
	}

	if err := uc.Users.EnhanceUser(in.RequestedUserID); err != nil {
		return dto.AssignRoleResponse{
			UserID:   in.RequestedUserID,
			Assigned: false,
		}, uc_errors.ErrEnhanceUser
	}
	return dto.AssignRoleResponse{
		UserID:   in.RequestedUserID,
		Assigned: true,
	}, nil
}
