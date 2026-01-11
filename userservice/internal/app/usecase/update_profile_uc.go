package usecase

import (
	"ads/userservice/internal/app/dto"
	"ads/userservice/internal/app/uc_errors"
	"ads/userservice/internal/domain/port"
	"ads/userservice/pkg/errs"
	"context"
	"errors"
)

type UpdateProfileUC struct {
	Profile        port.ProfileRepository
	PhoneValidator port.PhoneValidator
}

func NewUpdateProfileUC(
	profile port.ProfileRepository,
	phoneValidator port.PhoneValidator,
) *UpdateProfileUC {
	return &UpdateProfileUC{
		Profile:        profile,
		PhoneValidator: phoneValidator,
	}
}

func (uc *UpdateProfileUC) Execute(ctx context.Context, in dto.UpdateProfile) (dto.UpdateProfileOutput, error) {
	// Get from db
	profile, err := uc.Profile.Get(ctx, in.AccountID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.UpdateProfileOutput{Success: false},
				uc_errors.ErrInvalidAccountID
		}
		return dto.UpdateProfileOutput{Success: false},
			uc_errors.Wrap(uc_errors.ErrGetProfileDB, err)
	}

	// Phone number validation
	var validatedPhone *string
	if in.Phone != nil {
		normPhone, err := uc.PhoneValidator.Validate(ctx, *in.Phone)
		if err != nil {
			return dto.UpdateProfileOutput{Success: false},
				uc_errors.ErrInvalidPhoneNumber
		}
		validatedPhone = &normPhone
	}

	// Update
	err = profile.Update(
		in.FirstName,
		in.LastName,
		validatedPhone,
		in.AvatarURl,
		in.Bio,
	)
	if err != nil {
		return dto.UpdateProfileOutput{Success: false}, uc_errors.ErrInvalidProfileData
	}

	if err := uc.Profile.Update(ctx, profile); err != nil {
		return dto.UpdateProfileOutput{Success: false},
			uc_errors.Wrap(uc_errors.ErrUpdateProfileDB, err)
	}

	return dto.UpdateProfileOutput{Success: true}, nil
}
