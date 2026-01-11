package grpc

import (
	"ads/userservice/internal/app/usecase"
	"ads/userservice/internal/generated/user_v1"
	"context"
	"log/slog"
)

type UserHandler struct {
	user_v1.UnimplementedUserServiceServer
	log             *slog.Logger
	getProfileUC    *usecase.GetProfileUC
	updateProfileUC *usecase.UpdateProfileUC
}

func NewUserHandler(
	log *slog.Logger,
	getProfileUC *usecase.GetProfileUC,
	updateProfileUC *usecase.UpdateProfileUC,
) *UserHandler {
	return &UserHandler{
		log:             log,
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
	}
}

func (h *UserHandler) GetProfile(ctx context.Context, req *user_v1.GetProfileRequest) (*user_v1.GetProfileResponse, error) {
	ucResp, err := h.getProfileUC.Execute(ctx)
}
