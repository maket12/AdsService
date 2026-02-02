package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type DeleteAdUC struct {
	ad port.AdRepository
}

func NewDeleteAdUC(ad port.AdRepository) *DeleteAdUC {
	return &DeleteAdUC{ad: ad}
}

func (uc *DeleteAdUC) Execute(ctx context.Context, in dto.DeleteAdRequest) (dto.DeleteAdResponse, error) {
	// Get from db
	ad, err := uc.ad.Get(ctx, in.AdID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.DeleteAdResponse{Success: false}, uc_errors.ErrInvalidAdID
		}
		return dto.DeleteAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrGetAdDB, err,
		)
	}

	// Update status (deleted)
	err = ad.Delete()
	if err != nil {
		return dto.DeleteAdResponse{Success: false}, uc_errors.ErrCannotDelete
	}

	// Update status in db
	err = uc.ad.UpdateStatus(ctx, ad)
	if err != nil {
		return dto.DeleteAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrUpdateAdStatusDB, err,
		)
	}

	// Response
	return dto.DeleteAdResponse{Success: true}, nil
}
