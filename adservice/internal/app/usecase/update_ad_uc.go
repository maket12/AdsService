package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type UpdateAdUC struct {
	ad    port.AdRepository
	media port.MediaRepository
}

func NewUpdateAdUC(
	ad port.AdRepository, media port.MediaRepository,
) *UpdateAdUC {
	return &UpdateAdUC{
		ad:    ad,
		media: media,
	}
}

func (uc *UpdateAdUC) Execute(ctx context.Context, in dto.UpdateAdRequest) (dto.UpdateAdResponse, error) {
	// Get from db
	ad, err := uc.ad.Get(ctx, in.AdID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.UpdateAdResponse{Success: false}, uc_errors.ErrInvalidAdID
		}
		return dto.UpdateAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrGetAdDB, err,
		)
	}

	// Update
	err = ad.Update(in.Title, in.Description, in.Price, in.Images)
	if err != nil {
		return dto.UpdateAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}

	// Update in db
	err = uc.ad.Update(ctx, ad)
	if err != nil {
		return dto.UpdateAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrUpdateAdDB, err,
		)
	}

	// Update images in db
	err = uc.media.Save(ctx, ad.ID(), ad.Images())
	if err != nil {
		return dto.UpdateAdResponse{Success: false}, uc_errors.Wrap(
			uc_errors.ErrSaveImagesDB, err,
		)
	}

	// Response
	return dto.UpdateAdResponse{Success: true}, nil
}
