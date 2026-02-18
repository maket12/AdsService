package usecase

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/app/uc_errors"
	"ads/authservice/internal/domain/model"
	"ads/authservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type RegisterUC struct {
	account          port.AccountRepository
	accountRole      port.AccountRoleRepository
	passwordHasher   port.PasswordHasher
	accountPublisher port.AccountPublisher
}

func NewRegisterUC(
	account port.AccountRepository,
	accountRole port.AccountRoleRepository,
	passwordHasher port.PasswordHasher,
	accountPublisher port.AccountPublisher,
) *RegisterUC {
	return &RegisterUC{
		account:          account,
		accountRole:      accountRole,
		passwordHasher:   passwordHasher,
		accountPublisher: accountPublisher,
	}
}

func (uc *RegisterUC) Execute(ctx context.Context, in dto.RegisterInput) (dto.RegisterOutput, error) {
	// Hashing the password
	hashedPassword, err := uc.passwordHasher.Hash(in.Password)
	if err != nil {
		return dto.RegisterOutput{}, uc_errors.Wrap(
			uc_errors.ErrHashPassword, err,
		)
	}

	// Creating rich-models with validation
	account, err := model.NewAccount(in.Email, hashedPassword)
	if err != nil {
		return dto.RegisterOutput{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}
	accountRole, err := model.NewAccountRole(account.ID())
	if err != nil {
		return dto.RegisterOutput{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}

	// Save all into database
	if err := uc.account.Create(ctx, account); err != nil {
		if errors.Is(err, errs.ErrObjectAlreadyExists) {
			return dto.RegisterOutput{}, uc_errors.ErrAccountAlreadyExists
		}
		return dto.RegisterOutput{}, uc_errors.Wrap(
			uc_errors.ErrCreateAccountDB, err,
		)
	}
	if err := uc.accountRole.Create(ctx, accountRole); err != nil {
		return dto.RegisterOutput{}, uc_errors.Wrap(
			uc_errors.ErrCreateAccountRoleDB, err,
		)
	}

	// Send even to rabbitmq (create profile)
	if err := uc.accountPublisher.PublishAccountCreate(ctx, account.ID()); err != nil {
		return dto.RegisterOutput{},
			uc_errors.Wrap(uc_errors.ErrPublishEvent, err)
	}

	// Response
	return dto.RegisterOutput{AccountID: account.ID()}, nil
}
