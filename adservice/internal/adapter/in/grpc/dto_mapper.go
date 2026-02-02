package grpc

import (
	"ads/adservice/internal/app/dto"
	"ads/pkg/generated/ad_v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapCreateAdPbToDTO(req *ad_v1.CreateAdRequest) dto.CreateAdRequest {
	sellerID, _ := uuid.Parse(req.GetSellerId())
	return dto.CreateAdRequest{
		SellerID:    sellerID,
		Title:       req.GetTitle(),
		Description: req.Description,
		Price:       req.GetPrice(),
		Images:      req.GetImages(),
	}
}

func MapCreateAdDTOToPb(out dto.CreateAdResponse) *ad_v1.CreateAdResponse {
	return &ad_v1.CreateAdResponse{AdId: out.AdID.String()}
}

func MapGetAdPbToDTO(req *ad_v1.GetAdRequest) dto.GetAdRequest {
	adID, _ := uuid.Parse(req.GetAdId())
	return dto.GetAdRequest{AdID: adID}
}

func MapGetAdDTOToPb(out dto.GetAdResponse) *ad_v1.GetAdResponse {
	return &ad_v1.GetAdResponse{
		AdId:        out.AdID.String(),
		SellerId:    out.SellerID.String(),
		Title:       out.Title,
		Description: out.Description,
		Price:       out.Price,
		Status:      out.Status,
		Images:      out.Images,
		CreatedAt:   timestamppb.New(out.CreatedAt),
		UpdatedAt:   timestamppb.New(out.UpdatedAt),
	}
}

func MapUpdateAdPbToDTO(req *ad_v1.UpdateAdRequest) dto.UpdateAdRequest {
	adID, _ := uuid.Parse(req.GetAdId())
	return dto.UpdateAdRequest{
		AdID:        adID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Images:      req.Images,
	}
}

func MapUpdateAdDTOToPb(out dto.UpdateAdResponse) *ad_v1.UpdateAdResponse {
	return &ad_v1.UpdateAdResponse{Success: out.Success}
}

func MapPublishAdPbToDTO(req *ad_v1.PublishAdRequest) dto.PublishAdRequest {
	adID, _ := uuid.Parse(req.GetAdId())
	return dto.PublishAdRequest{AdID: adID}
}

func MapPublishAdDTOToPb(out dto.PublishAdResponse) *ad_v1.PublishAdResponse {
	return &ad_v1.PublishAdResponse{Success: out.Success}
}

func MapRejectAdPbToDTO(req *ad_v1.RejectAdRequest) dto.RejectAdRequest {
	adID, _ := uuid.Parse(req.GetAdId())
	return dto.RejectAdRequest{AdID: adID}
}

func MapRejectAdDTOToPb(out dto.RejectAdResponse) *ad_v1.RejectAdResponse {
	return &ad_v1.RejectAdResponse{Success: out.Success}
}

func MapDeleteAdPbToDTO(req *ad_v1.DeleteAdRequest) dto.DeleteAdRequest {
	adID, _ := uuid.Parse(req.GetAdId())
	return dto.DeleteAdRequest{AdID: adID}
}

func MapDeleteAdDTOToPb(out dto.DeleteAdResponse) *ad_v1.DeleteAdResponse {
	return &ad_v1.DeleteAdResponse{Success: out.Success}
}

func MapDeleteAllAdsPbToDTO(req *ad_v1.DeleteAllAdsRequest) dto.DeleteAllAdsRequest {
	sellerID, _ := uuid.Parse(req.GetSellerId())
	return dto.DeleteAllAdsRequest{SellerID: sellerID}
}

func MapDeleteAllAdsDTOToPb(out dto.DeleteAllAdsResponse) *ad_v1.DeleteAllAdsResponse {
	return &ad_v1.DeleteAllAdsResponse{Success: out.Success}
}
