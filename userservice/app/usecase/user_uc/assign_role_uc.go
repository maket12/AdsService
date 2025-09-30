package user_uc

import (
	"AdsService/userservice/app/dto/user_dto"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
)

type AssignRoleUC struct {
	Users port.UserRepository
}

func (uc *AssignRoleUC) Execute(in user_dto.AssignRoleDTO) (user_dto.AssignRoleResponseDTO, error) {
	if err := uc.Users.EnhanceUser(in.UserID); err != nil {
		return user_dto.AssignRoleResponseDTO{
			UserID:   in.UserID,
			Assigned: false,
		}, uc_errors.ErrEnhanceUser
	}
	return user_dto.AssignRoleResponseDTO{
		UserID:   in.UserID,
		Assigned: true,
	}, nil
}
