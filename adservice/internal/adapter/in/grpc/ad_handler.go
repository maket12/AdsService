package grpc

import (
	"ads/pkg/generated/ad_v1"
	"log/slog"
)

type AdHandler struct {
	ad_v1.UnimplementedAdServiceServer
	log *slog.Logger
}
