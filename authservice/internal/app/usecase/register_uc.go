package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/domain/port"
	"ads/authservice/internal/pkg/errs"
	"context"
	"errors"
)

type RegisterUC struct {
	Account        port.AccountRepository
	AccountRole    port.AccountRoleRepository
	PasswordHasher port.PasswordHasher
}

func NewRegisterUC(
	account port.AccountRepository,
	accountRole port.AccountRoleRepository,
	passwordHasher port.PasswordHasher,
) *RegisterUC {
	return &RegisterUC{
		Account:        account,
		AccountRole:    accountRole,
		PasswordHasher: passwordHasher,
	}
}

func (uc *RegisterUC) Execute(ctx context.Context, in dto.Register) (dto.RegisterResponse, error) {
	// Hashing the password
	hashedPassword, err := uc.PasswordHasher.Hash(in.Password)
	if err != nil {
		return dto.RegisterResponse{}, uc_errors.Wrap(
			uc_errors.ErrHashPassword, err,
		)
	}

	// Creating rich-models with validation
	account, err := model.NewAccount(in.Email, hashedPassword)
	if err != nil {
		return dto.RegisterResponse{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}
	accountRole, err := model.NewAccountRole(account.ID())
	if err != nil {
		return dto.RegisterResponse{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}

	// Save all into database
	if err := uc.Account.Create(ctx, account); err != nil {
		if errors.Is(err, errs.ErrObjectAlreadyExists) {
			return dto.RegisterResponse{}, uc_errors.ErrAccountAlreadyExists
		}
		return dto.RegisterResponse{}, uc_errors.Wrap(
			uc_errors.ErrCreateAccountDB, err,
		)
	}
	if err := uc.AccountRole.Create(ctx, accountRole); err != nil {
		return dto.RegisterResponse{}, uc_errors.Wrap(
			uc_errors.ErrCreateAccountRoleDB, err,
		)
	}

	// Response
	return dto.RegisterResponse{AccountID: account.ID()}, nil
}
