package grpc

import (
	"ads/adservice/internal/app/usecase"
	"ads/pkg/generated/ad_v1"
	"context"
	"log/slog"

	"google.golang.org/grpc/status"
)

type AdHandler struct {
	ad_v1.UnimplementedAdServiceServer
	log            *slog.Logger
	createAdUC     *usecase.CreateAdUC
	getAdUC        *usecase.GetAdUC
	updateAdUC     *usecase.UpdateAdUC
	publishAdUC    *usecase.PublishAdUC
	rejectAdUC     *usecase.RejectAdUC
	deleteAdUC     *usecase.DeleteAdUC
	deleteAllAdsUC *usecase.DeleteAllAdsUC
}

func NewAdHandler(
	log *slog.Logger,
	createAdUC *usecase.CreateAdUC,
	getAdUC *usecase.GetAdUC,
	updateAdUC *usecase.UpdateAdUC,
	publishAdUC *usecase.PublishAdUC,
	rejectAdUC *usecase.RejectAdUC,
	deleteAdUC *usecase.DeleteAdUC,
	deleteAllAdsUC *usecase.DeleteAllAdsUC,
) *AdHandler {
	return &AdHandler{
		log:            log,
		createAdUC:     createAdUC,
		getAdUC:        getAdUC,
		updateAdUC:     updateAdUC,
		publishAdUC:    publishAdUC,
		rejectAdUC:     rejectAdUC,
		deleteAdUC:     deleteAdUC,
		deleteAllAdsUC: deleteAllAdsUC,
	}
}

func (h *AdHandler) CreateAd(ctx context.Context, req *ad_v1.CreateAdRequest) (*ad_v1.CreateAdResponse, error) {
	ucResp, err := h.createAdUC.Execute(ctx, MapCreateAdPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to create ad",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapCreateAdDTOToPb(ucResp), nil
}

func (h *AdHandler) GetAd(ctx context.Context, req *ad_v1.GetAdRequest) (*ad_v1.GetAdResponse, error) {
	ucResp, err := h.getAdUC.Execute(ctx, MapGetAdPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to get ad",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapGetAdDTOToPb(ucResp), nil
}

func (h *AdHandler) UpdateAd(ctx context.Context, req *ad_v1.UpdateAdRequest) (*ad_v1.UpdateAdResponse, error) {
	ucResp, err := h.updateAdUC.Execute(ctx, MapUpdateAdPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to update ad",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapUpdateAdDTOToPb(ucResp), nil
}

func (h *AdHandler) PublishAd(ctx context.Context, req *ad_v1.PublishAdRequest) (*ad_v1.PublishAdResponse, error) {
	ucResp, err := h.publishAdUC.Execute(ctx, MapPublishAdPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to publish ad",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapPublishAdDTOToPb(ucResp), nil
}

func (h *AdHandler) RejectAd(ctx context.Context, req *ad_v1.RejectAdRequest) (*ad_v1.RejectAdResponse, error) {
	ucResp, err := h.rejectAdUC.Execute(ctx, MapRejectAdPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to reject ad",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapRejectAdDTOToPb(ucResp), nil
}

func (h *AdHandler) DeleteAd(ctx context.Context, req *ad_v1.DeleteAdRequest) (*ad_v1.DeleteAdResponse, error) {
	ucResp, err := h.deleteAdUC.Execute(ctx, MapDeleteAdPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to delete ad",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapDeleteAdDTOToPb(ucResp), nil
}

func (h *AdHandler) DeleteAllAds(ctx context.Context, req *ad_v1.DeleteAllAdsRequest) (*ad_v1.DeleteAllAdsResponse, error) {
	ucResp, err := h.deleteAllAdsUC.Execute(ctx, MapDeleteAllAdsPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to delete all ads",
			slog.String("seller_id", req.GetSellerId()),
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapDeleteAllAdsDTOToPb(ucResp), nil
}
