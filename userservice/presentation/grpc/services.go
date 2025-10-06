package grpc

import (
	"AdsService/userservice/app/dto"
	"AdsService/userservice/app/usecase"
	authmw "AdsService/userservice/infrastructure/grpc"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "AdsService/userservice/presentation/grpc/pb"

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

	out, err := s.AddProfileUC.Execute(dto.AddProfileDTO{UserID: userID, Name: req.Name, Phone: req.Phone})
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

	out, err := s.GetProfileUC.Execute(dto.GetProfileDTO{UserID: userID})
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

	out, err := s.UpdateProfileUC.Execute(dto.UpdateProfileDTO{
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

	out, err := s.UploadPhotoUC.Execute(dto.UploadPhotoDTO{
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

	out, err := s.ChangeSettingUC.Execute(dto.ChangeSettingsDTO{
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

	out, err := s.ChangeSubscriptionsUC.Execute(dto.ChangeSubscriptionsDTO{
		UserID:        userID,
		Subscriptions: req.Subscriptions,
	})
	if err != nil {
		return MapProfileResponseDTOToPB(out), status.Errorf(codes.Internal, "%s", err.Error())
	}
	return MapProfileResponseDTOToPB(out), nil
}
