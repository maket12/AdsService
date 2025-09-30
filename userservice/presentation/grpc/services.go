package grpc

import (
	"AdsService/userservice/app/dto/admin_dto"
	"AdsService/userservice/app/dto/profile_dto"
	authmw "AdsService/userservice/infrastructure/grpc"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"AdsService/userservice/app/dto/user_dto"
	"AdsService/userservice/app/usecase/admin_uc"
	"AdsService/userservice/app/usecase/profile_uc"
	"AdsService/userservice/app/usecase/user_uc"
	pb "AdsService/userservice/presentation/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	AssignRoleUC          *user_uc.AssignRoleUC
	GetUserUC             *user_uc.GetUserUC
	AddProfileUC          *profile_uc.AddProfileUC
	UpdateProfileUC       *profile_uc.UpdateProfileUC
	UploadPhotoUC         *profile_uc.UploadPhotoUC
	ChangeSettingUC       *profile_uc.ChangeSettingsUC
	ChangeSubscriptionsUC *profile_uc.ChangeSubscriptionsUC
	GetProfileUC          *profile_uc.GetProfileUC
	AdminBanUserUC        *admin_uc.AdminBanUserUC
	AdminUnbanUserUC      *admin_uc.AdminUnbanUserUC
	AdminGetProfileUC     *admin_uc.AdminGetProfileUC
	AdminGetProfilesUC    *admin_uc.AdminGetProfilesUC
}

// User

func (s *UserService) AssignRole(ctx context.Context, req *pb.AssignRolRequest) (*pb.AssignRoleResponse, error) {
	out, err := s.AssignRoleUC.Execute(user_dto.AssignRoleDTO{UserID: req.UserId})
	if err != nil {
		return MapAssignRoleResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapAssignRoleResponseDTOToPB(out), nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	out, err := s.GetUserUC.Execute(user_dto.GetUserDTO{UserID: req.UserId})
	if err != nil {
		return MapGetUserResponseDTOToPB(out), status.Errorf(codes.NotFound, err.Error())
	}
	return MapGetUserResponseDTOToPB(out), nil
}

// Profile

func (s *UserService) AddProfile(ctx context.Context, req *pb.AddProfileRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.AddProfileUC.Execute(profile_dto.AddProfileDTO{UserID: userID, Name: req.Name, Phone: req.Phone})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) GetProfile(ctx context.Context, req *emptypb.Empty) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.GetProfileUC.Execute(profile_dto.GetProfileDTO{UserID: userID})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.UpdateProfileUC.Execute(profile_dto.UpdateProfileDTO{
		UserID: userID,
		Name:   req.Name,
		Phone:  req.Phone,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) UploadPhoto(ctx context.Context, req *pb.UploadPhotoRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.UploadPhotoUC.Execute(profile_dto.UploadPhotoDTO{
		UserID:      userID,
		Data:        req.Data,
		Filename:    req.Filename,
		ContentType: req.ContentType,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) ChangeSettings(ctx context.Context, req *pb.ChangeSettingsRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.ChangeSettingUC.Execute(profile_dto.ChangeSettingsDTO{
		UserID:               userID,
		NotificationsEnabled: req.NotificationsEnabled,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) ChangeSubscriptions(ctx context.Context, req *pb.ChangeSubscriptionsRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.ChangeSubscriptionsUC.Execute(profile_dto.ChangeSubscriptionsDTO{
		UserID:        userID,
		Subscriptions: req.Subscriptions,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

// Admin

func (s *UserService) AdminBanUser(ctx context.Context, req *pb.AdminBanUserRequest) (*pb.AdminBanUserResponse, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.AdminBanUserUC.Execute(admin_dto.BanUserDTO{
		UserID:          adminID,
		RequestedUserID: req.UserId,
	})
	if err != nil {
		return MapBanUserResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapBanUserResponseDTOToPB(out), nil
}

func (s *UserService) AdminUnbanUser(ctx context.Context, req *pb.AdminUnbanUserRequest) (*pb.AdminUnbanUserResponse, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.AdminUnbanUserUC.Execute(admin_dto.UnbanUserDTO{
		UserID:          adminID,
		RequestedUserID: req.UserId,
	})
	if err != nil {
		return MapUnbanUserResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapUnbanUserResponseDTOToPB(out), nil
}

func (s *UserService) AdminGetProfile(ctx context.Context, req *pb.AdminGetProfileRequest) (*pb.Profile, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.AdminGetProfileUC.Execute(admin_dto.AdminGetProfileDTO{
		UserID:          adminID,
		RequestedUserID: req.UserId,
	})
	if err != nil {
		return MapGetProfileResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapGetProfileResponseDTOToPB(out), nil
}

func (s *UserService) AdminGetProfiles(ctx context.Context, req *pb.AdminGetProfilesListRequest) (*pb.AdminGetProfilesListResponse, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	out, err := s.AdminGetProfilesUC.Execute(admin_dto.AdminGetProfilesListDTO{
		UserID: adminID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return MapGetProfilesResponseDTOToPB(out), status.Errorf(codes.Internal, err.Error())
	}
	return MapGetProfilesResponseDTOToPB(out), nil
}
