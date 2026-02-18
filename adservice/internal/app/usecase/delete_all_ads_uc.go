package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/app/uc_errors"
	"ads/adservice/internal/domain/port"
	"context"
)

type DeleteAllAdsUC struct {
	ad port.AdRepository
}

func NewDeleteAllAdsUC(ad port.AdRepository) *DeleteAllAdsUC {
	return &DeleteAllAdsUC{ad: ad}
}

func (uc *DeleteAllAdsUC) Execute(ctx context.Context, in dto.DeleteAllAdsInput) (dto.DeleteAllAdsOutput, error) {
	// Delete all ads
	if err := uc.ad.DeleteAll(ctx, in.SellerID); err != nil {
		return dto.DeleteAllAdsOutput{Success: false}, uc_errors.Wrap(
			uc_errors.ErrDeleteAllAdsDB, err,
		)
	}

	// Response
	return dto.DeleteAllAdsOutput{Success: true}, nil
}
