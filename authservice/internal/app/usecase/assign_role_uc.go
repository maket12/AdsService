package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type AssignRoleUC struct {
	AccountRole port.AccountRoleRepository
}

func NewAssignRoleUC(accountRole port.AccountRoleRepository) *AssignRoleUC {
	return &AssignRoleUC{AccountRole: accountRole}
}

func (uc *AssignRoleUC) Execute(ctx context.Context, in dto.AssignRole) (dto.AssignRoleResponse, error) {
	// Get role
	accRole, err := uc.AccountRole.Get(ctx, in.AccountID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.AssignRoleResponse{Assign: false},
				uc_errors.ErrInvalidAccountID
		}
		return dto.AssignRoleResponse{Assign: false},
			uc_errors.Wrap(uc_errors.ErrGetAccountRoleDB, err)
	}

	// Assign
	if err := accRole.Assign(in.Role); err != nil {
		return dto.AssignRoleResponse{Assign: false},
			uc_errors.ErrCannotAssign
	}

	// Update db
	if err := uc.AccountRole.Update(ctx, accRole); err != nil {
		return dto.AssignRoleResponse{Assign: false},
			uc_errors.Wrap(uc_errors.ErrUpdateAccountRoleDB, err)
	}

	// Output
	return dto.AssignRoleResponse{Assign: true}, nil
}
