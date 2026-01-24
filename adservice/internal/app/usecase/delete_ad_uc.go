package usecase

import (
	"ads/adservice/internal/app/dto"
	"ads/adservice/internal/domain/port"
	"context"
)

type DeleteAdUC struct {
	ad port.AdRepository
}

func NewDeleteAdUC(ad port.AdRepository) *DeleteAdUC {
	return &DeleteAdUC{ad: ad}
}

func (uc *DeleteAdUC) Execute(ctx context.Context, in dto.DeleteAdRequest) (dto.DeleteAdResponse, error) {

}
