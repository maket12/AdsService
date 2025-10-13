package grpc

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/usecase"
	authmw "ads/userservice/infrastructure/grpc"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "ads/userservice/presentation/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUsersServiceServer
	AddProfileUC          *usecase.AddProfileUC
	UpdateProfileUC       *usecase.UpdateProfileUC
	UploadPhotoUC         *usecase.UploadPhotoUC
	ChangeSettingUC       *usecase.ChangeSettingsUC
	ChangeSubscriptionsUC *usecase.ChangeSubscriptionsUC
	GetProfileUC          *usecase.GetProfileUC
}

func (s *UserService) AddProfile(ctx context.Context, req *pb.AddProfileRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	out, err := s.AddProfileUC.Execute(ctx, dto.AddProfile{UserID: userID, Name: req.Name, Phone: req.Phone})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) GetProfile(ctx context.Context, req *emptypb.Empty) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.GetProfileUC.Execute(ctx, dto.GetProfile{UserID: userID})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.UpdateProfileUC.Execute(ctx, dto.UpdateProfile{
		UserID: userID,
		Name:   req.Name,
		Phone:  req.Phone,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) UploadPhoto(ctx context.Context, req *pb.UploadPhotoRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.UploadPhotoUC.Execute(ctx, dto.UploadPhoto{
		UserID:      userID,
		Data:        req.Data,
		Filename:    req.Filename,
		ContentType: req.ContentType,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) ChangeSettings(ctx context.Context, req *pb.ChangeSettingsRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.ChangeSettingUC.Execute(ctx, dto.ChangeSettings{
		UserID:               userID,
		NotificationsEnabled: req.NotificationsEnabled,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}

func (s *UserService) ChangeSubscriptions(ctx context.Context, req *pb.ChangeSubscriptionsRequest) (*pb.Profile, error) {
	userID, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}

	out, err := s.ChangeSubscriptionsUC.Execute(ctx, dto.ChangeSubscriptions{
		UserID:        userID,
		Subscriptions: req.Subscriptions,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}
