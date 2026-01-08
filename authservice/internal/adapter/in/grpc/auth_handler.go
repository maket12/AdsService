package grpc

import (
	"ads/authservice/internal/app/usecase"
	"ads/authservice/internal/generated/auth_v1"
	"context"
	"log/slog"

	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	auth_v1.UnimplementedAuthServiceServer
	log                   slog.Logger
	registerUC            *usecase.RegisterUC
	loginUC               *usecase.LoginUC
	logoutUC              *usecase.LogoutUC
	refreshSessionUC      *usecase.RefreshSessionUC
	validateAccessTokenUC *usecase.ValidateAccessTokenUC
}

func NewAuthHandler(
	log slog.Logger,
	registerUC *usecase.RegisterUC,
	loginUC *usecase.LoginUC,
	logoutUC *usecase.LogoutUC,
	refreshSessionUC *usecase.RefreshSessionUC,
	validateAccessTokenUC *usecase.ValidateAccessTokenUC,
) *AuthHandler {
	return &AuthHandler{
		log:                   log,
		registerUC:            registerUC,
		loginUC:               loginUC,
		logoutUC:              logoutUC,
		refreshSessionUC:      refreshSessionUC,
		validateAccessTokenUC: validateAccessTokenUC,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {
	ucResp, err := h.registerUC.Execute(ctx, MapRegisterPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to register",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapRegisterDTOToPb(ucResp), nil
}

func (h *AuthHandler) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	ucResp, err := h.loginUC.Execute(ctx, MapLoginPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to login",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapLoginDTOToPb(ucResp), nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *auth_v1.LogoutRequest) (*auth_v1.LogoutResponse, error) {
	ucResp, err := h.logoutUC.Execute(ctx, MapLogoutPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to logout",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapLogoutDTOToPb(ucResp), nil
}

func (h *AuthHandler) RefreshSession(ctx context.Context, req *auth_v1.RefreshSessionRequest) (*auth_v1.RefreshSessionResponse, error) {
	ucResp, err := h.refreshSessionUC.Execute(ctx, MapRefreshSessionPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to refresh session",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapRefreshSessionDTOToPb(ucResp), nil
}

func (h *AuthHandler) ValidateAccessToken(ctx context.Context, req *auth_v1.ValidateAccessTokenRequest) (*auth_v1.ValidateAccessTokenResponse, error) {
	ucResp, err := h.validateAccessTokenUC.Execute(ctx, MapValidateAccessTokenPbToDTO(req))

	if err != nil {
		code, msg, internalErr := gRPCError(err)
		h.log.ErrorContext(ctx, "failed to validate access token",
			slog.Int("code", int(code)),
			slog.String("public_msg", msg),
			slog.Any("reason", internalErr),
		)
		return nil, status.Error(code, msg)
	}

	return MapValidateAccessTokenDTOToPb(ucResp), nil
}
