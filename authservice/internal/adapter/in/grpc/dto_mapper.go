package grpc

import (
	"ads/authservice/internal/app/dto"
	"ads/authservice/internal/generated/auth_v1"

	"github.com/google/uuid"
)

func MapRegisterPbToDTO(req *auth_v1.RegisterRequest) dto.Register {
	return dto.Register{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func MapRegisterDTOToPb(out dto.RegisterResponse) *auth_v1.RegisterResponse {
	return &auth_v1.RegisterResponse{AccountId: out.AccountID.String()}
}

func MapLoginPbToDTO(req *auth_v1.LoginRequest) dto.Login {
	var ip, userAgent = req.GetIp(), req.GetUserAgent()
	return dto.Login{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		IP:        &ip,
		UserAgent: &userAgent,
	}
}

func MapLoginDTOToPb(out dto.LoginResponse) *auth_v1.LoginResponse {
	return &auth_v1.LoginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}
}

func MapLogoutPbToDTO(req *auth_v1.LogoutRequest) dto.Logout {
	return dto.Logout{RefreshToken: req.GetRefreshToken()}
}

func MapLogoutDTOToPb(out dto.LogoutResponse) *auth_v1.LogoutResponse {
	return &auth_v1.LogoutResponse{Logout: out.Logout}
}

func MapRefreshSessionPbToDTO(req *auth_v1.RefreshSessionRequest) dto.RefreshSession {
	var ip, userAgent = req.GetIp(), req.GetUserAgent()
	return dto.RefreshSession{
		OldRefreshToken: req.GetOldRefreshToken(),
		IP:              &ip,
		UserAgent:       &userAgent,
	}
}

func MapRefreshSessionDTOToPb(out dto.RefreshSessionResponse) *auth_v1.RefreshSessionResponse {
	return &auth_v1.RefreshSessionResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}
}

func MapValidateAccessTokenPbToDTO(req *auth_v1.ValidateAccessTokenRequest) dto.ValidateAccessToken {
	return dto.ValidateAccessToken{AccessToken: req.GetAccessToken()}
}

func MapValidateAccessTokenDTOToPb(out dto.ValidateAccessTokenResponse) *auth_v1.ValidateAccessTokenResponse {
	return &auth_v1.ValidateAccessTokenResponse{
		AccountId: out.AccountID.String(),
		Role:      out.Role,
	}
}

func MapAssignRolePbToDTO(req *auth_v1.AssignRoleRequest) dto.AssignRole {
	accID, _ := uuid.Parse(req.GetAccountId())
	return dto.AssignRole{
		AccountID: accID,
		Role:      req.GetRole(),
	}
}

func MapAssignRoleDTOToPb(out dto.AssignRoleResponse) *auth_v1.AssignRoleResponse {
	return &auth_v1.AssignRoleResponse{Assign: out.Assign}
}
