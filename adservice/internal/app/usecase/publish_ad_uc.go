package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type PublishAdUC struct {
	ad port.AdRepository
}

func NewPublishAdUC(ad port.AdRepository) *PublishAdUC {
	return &PublishAdUC{ad: ad}
}

func (uc *PublishAdUC) Execute(ctx context.Context, in dto.PublishAdRequest) (dto.PublishAdResponse, error) {
	// Get from db
	ad, err := uc.ad.Get(ctx, in.AdID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.PublishAdResponse{Success: false}, uc_errors.ErrInvalidAdID
		}
		return dto.PublishAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrGetAdDB, err,
		)
	}

	// Publish
	err = ad.Publish()
	if err != nil {
		return dto.PublishAdResponse{Success: false}, uc_errors.ErrCannotPublish
	}

	// Update in db
	err = uc.ad.UpdateStatus(ctx, ad)
	if err != nil {
		return dto.PublishAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrUpdateAdStatusDB, err,
		)
	}

	// Response
	return dto.PublishAdResponse{Success: true}, nil
}
