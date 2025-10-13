package usecase

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/uc_errors"
	"ads/adminservice/domain/port"
	"context"
)

type UnbanUserUC struct {
	Users port.UserRepository
}

func (uc *UnbanUserUC) Execute(ctx context.Context, in dto.UnbanUser) (dto.UnbanUserResponse, error) {
	role, err := uc.Users.GetUserRole(ctx, in.UserID)
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

	if err := uc.Users.UnbanUser(ctx, in.RequestedUserID); err != nil {
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
