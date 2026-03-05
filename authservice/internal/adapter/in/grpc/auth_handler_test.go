package grpc_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/maket12/ads-service/authservice/internal/adapter/in/grpc"
	"github.com/maket12/ads-service/authservice/internal/app/dto"
	ucerrs "github.com/maket12/ads-service/authservice/internal/app/errs"
	"github.com/maket12/ads-service/authservice/internal/app/usecase/mocks"
	"github.com/maket12/ads-service/pkg/generated/auth_v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

func TestAH_Register(t *testing.T) {
	testUID := uuid.New()

	type testCase struct {
		name      string
		request   *auth_v1.RegisterRequest
		setupMock func(m *mocks.RegisterUseCase)
		wantCode  codes.Code
		wantResp  *auth_v1.RegisterResponse
	}

	testCases := []testCase{
		{
			name: "Success registration",
			request: &auth_v1.RegisterRequest{
				Email:    "shishi12377@weixin.cn",
				Password: "liushi07.12.2006",
			},
			setupMock: func(m *mocks.RegisterUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.RegisterOutput{AccountID: testUID}, nil)
			},
			wantCode: codes.OK,
			wantResp: &auth_v1.RegisterResponse{AccountId: testUID.String()},
		},
		{
			name: "Failure - invalid input",
			request: &auth_v1.RegisterRequest{
				Email:    "",
				Password: "",
			},
			setupMock: func(m *mocks.RegisterUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.RegisterOutput{AccountID: uuid.Nil}, ucerrs.ErrInvalidInput)
			},
			wantCode: codes.InvalidArgument,
			wantResp: &auth_v1.RegisterResponse{AccountId: uuid.Nil.String()},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockReg := mocks.NewRegisterUseCase(t)
			if tt.setupMock != nil {
				tt.setupMock(mockReg)
			}

			handler := grpc.NewAuthHandler(slog.Default(), mockReg, nil,
				nil, nil, nil, nil,
			)

			resp, err := handler.Register(context.Background(), tt.request)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestAH_Login(t *testing.T) {
	testAccess := "access-token"
	testRefresh := "refresh_token"

	type testCase struct {
		name      string
		request   *auth_v1.LoginRequest
		setupMock func(m *mocks.LoginUseCase)
		wantCode  codes.Code
		wantResp  *auth_v1.LoginResponse
	}
	testCases := []testCase{
		{
			name: "Success login",
			request: &auth_v1.LoginRequest{
				Email:    "zaizai@yummy.com",
				Password: "i bother ShiShi",
			},
			setupMock: func(m *mocks.LoginUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.LoginOutput{
						AccessToken:  testAccess,
						RefreshToken: testRefresh,
					}, nil)
			},
			wantCode: codes.OK,
			wantResp: &auth_v1.LoginResponse{
				AccessToken:  testAccess,
				RefreshToken: testRefresh,
			},
		},
		{
			name: "Failure - internal error",
			request: &auth_v1.LoginRequest{
				Email:    "zaizai@yummy.com",
				Password: "i bother ShiShi",
			},
			setupMock: func(m *mocks.LoginUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.LoginOutput{}, ucerrs.ErrGenerateAccessToken)
			},
			wantCode: codes.Internal,
			wantResp: &auth_v1.LoginResponse{
				AccessToken:  "",
				RefreshToken: "",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockLogin := mocks.NewLoginUseCase(t)
			if tt.setupMock != nil {
				tt.setupMock(mockLogin)
			}

			handler := grpc.NewAuthHandler(
				slog.Default(), nil, mockLogin,
				nil, nil, nil,
				nil,
			)

			resp, err := handler.Login(context.Background(), tt.request)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestAH_Logout(t *testing.T) {
	type testCase struct {
		name      string
		request   *auth_v1.LogoutRequest
		setupMock func(m *mocks.LogoutUseCase)
		wantCode  codes.Code
		wantResp  *auth_v1.LogoutResponse
	}
	testCases := []testCase{
		{
			name: "Success logout",
			request: &auth_v1.LogoutRequest{
				RefreshToken: "refresh-token",
			},
			setupMock: func(m *mocks.LogoutUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.LogoutOutput{
						Logout: true,
					}, nil)
			},
			wantCode: codes.OK,
			wantResp: &auth_v1.LogoutResponse{
				Logout: true,
			},
		},
		{
			name: "Failure - unauthenticated",
			request: &auth_v1.LogoutRequest{
				RefreshToken: "no-valid",
			},
			setupMock: func(m *mocks.LogoutUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.LogoutOutput{
						Logout: false,
					}, ucerrs.ErrInvalidRefreshToken)
			},
			wantCode: codes.Unauthenticated,
			wantResp: &auth_v1.LogoutResponse{
				Logout: false,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockLogout := mocks.NewLogoutUseCase(t)
			if tt.setupMock != nil {
				tt.setupMock(mockLogout)
			}

			handler := grpc.NewAuthHandler(
				slog.Default(), nil, nil,
				mockLogout, nil, nil,
				nil,
			)

			resp, err := handler.Logout(context.Background(), tt.request)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestAH_RefreshSession(t *testing.T) {
	type testCase struct {
		name      string
		request   *auth_v1.RefreshSessionRequest
		setupMock func(m *mocks.RefreshSessionUseCase)
		wantCode  codes.Code
		wantResp  *auth_v1.RefreshSessionResponse
	}
	testCases := []testCase{
		{
			name: "Success refresh",
			request: &auth_v1.RefreshSessionRequest{
				OldRefreshToken: "refresh-token",
			},
			setupMock: func(m *mocks.RefreshSessionUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.RefreshSessionOutput{
						AccessToken:  "new-access-token",
						RefreshToken: "new-refresh-token",
					}, nil)
			},
			wantCode: codes.OK,
			wantResp: &auth_v1.RefreshSessionResponse{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			},
		},
		{
			name: "Failure - internal error",
			request: &auth_v1.RefreshSessionRequest{
				OldRefreshToken: "refresh-token",
			},
			setupMock: func(m *mocks.RefreshSessionUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.RefreshSessionOutput{}, ucerrs.ErrRevokeRefreshSessionDB)
			},
			wantCode: codes.Internal,
			wantResp: &auth_v1.RefreshSessionResponse{
				AccessToken:  "",
				RefreshToken: "",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRefresh := mocks.NewRefreshSessionUseCase(t)
			if tt.setupMock != nil {
				tt.setupMock(mockRefresh)
			}

			handler := grpc.NewAuthHandler(
				slog.Default(), nil, nil,
				nil, mockRefresh, nil,
				nil,
			)

			resp, err := handler.RefreshSession(context.Background(), tt.request)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestAH_ValidateAccessToken(t *testing.T) {
	testUID := uuid.New()

	type testCase struct {
		name      string
		request   *auth_v1.ValidateAccessTokenRequest
		setupMock func(m *mocks.ValidateAccessTokenUseCase)
		wantCode  codes.Code
		wantResp  *auth_v1.ValidateAccessTokenResponse
	}
	testCases := []testCase{
		{
			name: "Success validation",
			request: &auth_v1.ValidateAccessTokenRequest{
				AccessToken: "access-token",
			},
			setupMock: func(m *mocks.ValidateAccessTokenUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.ValidateAccessTokenOutput{
						AccountID: testUID,
						Role:      "user",
					}, nil)
			},
			wantCode: codes.OK,
			wantResp: &auth_v1.ValidateAccessTokenResponse{
				AccountId: testUID.String(),
				Role:      "user",
			},
		},
		{
			name: "Failure - internal error",
			request: &auth_v1.ValidateAccessTokenRequest{
				AccessToken: "no-valid",
			},
			setupMock: func(m *mocks.ValidateAccessTokenUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.ValidateAccessTokenOutput{}, ucerrs.ErrInvalidAccessToken)
			},
			wantCode: codes.Unauthenticated,
			wantResp: &auth_v1.ValidateAccessTokenResponse{
				AccountId: "",
				Role:      "",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockValidate := mocks.NewValidateAccessTokenUseCase(t)
			if tt.setupMock != nil {
				tt.setupMock(mockValidate)
			}

			handler := grpc.NewAuthHandler(
				slog.Default(), nil, nil,
				nil, nil, mockValidate,
				nil,
			)

			resp, err := handler.ValidateAccessToken(context.Background(), tt.request)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestAH_AssignRole(t *testing.T) {
	type testCase struct {
		name      string
		request   *auth_v1.AssignRoleRequest
		setupMock func(m *mocks.AssignRoleUseCase)
		wantCode  codes.Code
		wantResp  *auth_v1.AssignRoleResponse
	}
	testCases := []testCase{
		{
			name: "Success logout",
			request: &auth_v1.AssignRoleRequest{
				AccountId: uuid.New().String(),
			},
			setupMock: func(m *mocks.AssignRoleUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.AssignRoleOutput{
						Assign: true,
					}, nil)
			},
			wantCode: codes.OK,
			wantResp: &auth_v1.AssignRoleResponse{
				Assign: true,
			},
		},
		{
			name: "Failure - precondition",
			request: &auth_v1.AssignRoleRequest{
				AccountId: uuid.New().String(),
			},
			setupMock: func(m *mocks.AssignRoleUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).
					Return(dto.AssignRoleOutput{
						Assign: false,
					}, ucerrs.ErrCannotAssign)
			},
			wantCode: codes.FailedPrecondition,
			wantResp: &auth_v1.AssignRoleResponse{
				Assign: false,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockAssign := mocks.NewAssignRoleUseCase(t)
			if tt.setupMock != nil {
				tt.setupMock(mockAssign)
			}

			handler := grpc.NewAuthHandler(
				slog.Default(), nil, nil,
				nil, nil, nil,
				mockAssign,
			)

			resp, err := handler.AssignRole(context.Background(), tt.request)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
