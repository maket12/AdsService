package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/model"
	"ads/adservice/internal/domain/port"
	"context"
)

type CreateAdUC struct {
	ad    port.AdRepository
	media port.MediaRepository
}

func NewCreateAdUC(
	ad port.AdRepository, media port.MediaRepository,
) *CreateAdUC {
	return &CreateAdUC{
		ad:    ad,
		media: media,
	}
}

func (uc *CreateAdUC) Execute(ctx context.Context, in dto.CreateAdInput) (dto.CreateAdOutput, error) {
	// Create ad
	ad, err := model.NewAd(
		in.SellerID, in.Title,
		in.Description, in.Price, in.Images,
	)
	if err != nil {
		return dto.CreateAdOutput{}, uc_errors.Wrap(
			uc_errors.ErrInvalidInput, err,
		)
	}

	// Save into database
	if err := uc.ad.Create(ctx, ad); err != nil {
		return dto.CreateAdOutput{}, uc_errors.Wrap(
			uc_errors.ErrCreateAdDB, err,
		)
	}

	// Save images into database
	if err := uc.media.Save(ctx, ad.ID(), ad.Images()); err != nil {
		return dto.CreateAdOutput{}, uc_errors.Wrap(uc_errors.ErrSaveImagesDB, err)
	}

	// Response
	return dto.CreateAdOutput{AdID: ad.ID()}, nil
}
