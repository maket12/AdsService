package grpc

import (
	"ads/adminservice/app/dto"
	"ads/adminservice/app/usecase"
	authmw "ads/adminservice/infrastructure/grpc"
	"context"

	pb "ads/adminservice/presentation/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminService struct {
	pb.UnimplementedAdminServiceServer
	BanUserUC     *usecase.BanUserUC
	UnbanUserUC   *usecase.UnbanUserUC
	AssignRoleUC  *usecase.AssignRoleUC
	GetUserUC     *usecase.GetUserUC
	GetProfileUC  *usecase.GetProfileUC
	GetProfilesUC *usecase.GetProfilesUC
}

func (s *AdminService) AdminBanUser(ctx context.Context, req *pb.BanUserRequest) (*pb.BanUserResponse, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.BanUserUC.Execute(ctx, dto.BanUser{
		UserID:          adminID,
		RequestedUserID: req.UserId,
	})
	if err != nil {
		return MapBanUserResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapBanUserResponseDTOToPB(out), nil
}

func (s *AdminService) AdminUnbanUser(ctx context.Context, req *pb.UnbanUserRequest) (*pb.UnbanUserResponse, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.UnbanUserUC.Execute(ctx, dto.UnbanUser{
		UserID:          adminID,
		RequestedUserID: req.UserId,
	})
	if err != nil {
		return MapUnbanUserResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapUnbanUserResponseDTOToPB(out), nil
}

func (s *AdminService) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	out, err := s.AssignRoleUC.Execute(ctx, dto.AssignRole{RequestedUserID: req.UserId})
	if err != nil {
		return MapAssignRoleResponseDTOToPB(out), status.Error(codes.Internal, err.Error())
	}
	return MapAssignRoleResponseDTOToPB(out), nil
}

func (s *AdminService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	out, err := s.GetUserUC.Execute(ctx, dto.GetUser{AdminUserID: req.UserId})
	if err != nil {
		return MapGetUserResponseDTOToPB(out), status.Error(codes.NotFound, err.Error())
	}
	return MapGetUserResponseDTOToPB(out), nil
}

func (s *AdminService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.Profile, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.GetProfileUC.Execute(ctx, dto.GetProfile{
		UserID:          adminID,
		RequestedUserID: req.UserId,
	})
	if err != nil {
		return MapGetProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapGetProfileResponseDTOToPB(out), nil
}

func (s *AdminService) AdminGetProfiles(ctx context.Context, req *pb.GetProfilesListRequest) (*pb.GetProfilesListResponse, error) {
	adminID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.GetProfilesUC.Execute(ctx, dto.GetProfilesList{
		UserID: adminID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return MapGetProfilesResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapGetProfilesResponseDTOToPB(out), nil
}
