package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type RejectAdUC struct {
	ad port.AdRepository
}

func NewRejectAdUC(ad port.AdRepository) *RejectAdUC {
	return &RejectAdUC{ad: ad}
}

func (uc *RejectAdUC) Execute(ctx context.Context, in dto.RejectAdRequest) (dto.RejectAdResponse, error) {
	// Get from db
	ad, err := uc.ad.Get(ctx, in.AdID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.RejectAdResponse{Success: false}, uc_errors.ErrInvalidAdID
		}
		return dto.RejectAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrGetAdDB, err,
		)
	}

	// Reject
	err = ad.Reject()
	if err != nil {
		return dto.RejectAdResponse{Success: false}, uc_errors.ErrCannotReject
	}

	// Update in db
	err = uc.ad.UpdateStatus(ctx, ad)
	if err != nil {
		return dto.RejectAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrUpdateAdStatusDB, err,
		)
	}

	// Response
	return dto.RejectAdResponse{Success: true}, nil
}
