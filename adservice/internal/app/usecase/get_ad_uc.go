package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/port"
	"ads/pkg/errs"
	"context"
	"errors"
)

type GetAdUC struct {
	ad    port.AdRepository
	media port.MediaRepository
}

func NewGetAdUC(
	ad port.AdRepository, media port.MediaRepository,
) *GetAdUC {
	return &GetAdUC{
		ad:    ad,
		media: media,
	}
}

func (uc *GetAdUC) Execute(ctx context.Context, in dto.GetAdRequest) (dto.GetAdResponse, error) {
	// Get from db
	ad, err := uc.ad.Get(ctx, in.AdID)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return dto.GetAdResponse{}, uc_errors.ErrInvalidAdID
		}
		return dto.GetAdResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetAdDB, err,
		)
	}

	// Get images from db
	images, err := uc.media.Get(ctx, ad.ID())
	if err != nil {
		return dto.GetAdResponse{}, uc_errors.Wrap(
			uc_errors.ErrGetImagesDB, err,
		)
	}

	// Add images into rich model
	err = ad.Update(nil, nil, nil, images)
	if err != nil {
		return dto.GetAdResponse{}, uc_errors.ErrInvalidInput
	}

	// Response
	return dto.GetAdResponse{
		AdID:        ad.ID(),
		SellerID:    ad.SellerID(),
		Title:       ad.Title(),
		Description: ad.Description(),
		Price:       ad.Price(),
		Status:      string(ad.Status()),
		Images:      ad.Images(),
		CreatedAt:   ad.CreatedAt(),
		UpdatedAt:   ad.UpdatedAt(),
	}, nil
}
