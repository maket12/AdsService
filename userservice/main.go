package main

import (
	"AdsService/infra/authmw"
	"AdsService/infra/database"
	"AdsService/infra/mongodb"
	pb "AdsService/userservice/proto"
	"bytes"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

func toPbProfile(p *database.Profile) *pb.Profile {
	if p == nil {
		return nil
	}
	subs := p.Subscriptions
	if subs == nil {
		subs = []string{}
	}
	return &pb.Profile{
		UserId:               p.UserID,
		Name:                 p.Name,
		Phone:                p.Phone,
		PhotoId:              p.PhotoID,
		NotificationsEnabled: p.NotificationsEnabled,
		Subscriptions:        subs,
		UpdatedAt:            timestamppb.New(p.UpdatedAt),
	}
}

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserService) AddProfile(ctx context.Context, req *pb.AddProfileRequest) (*pb.Profile, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := database.AddProfile(uid, req.Name, req.Phone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return toPbProfile(profile), nil
}

func (s *UserService) GetProfile(ctx context.Context, _ *emptypb.Empty) (*pb.Profile, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := database.GetProfile(uid)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return toPbProfile(profile), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Profile, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := database.UpdateProfileName(uid, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}
	profile, err = database.UpdateProfilePhone(uid, req.Phone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return toPbProfile(profile), nil
}

func (s *UserService) UploadPhoto(ctx context.Context, req *pb.UploadPhotoRequest) (*pb.Profile, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty file data")
	}
	if req.Filename == "" {
		return nil, status.Error(codes.InvalidArgument, "filename required")
	}
	if req.ContentType == "" {
		return nil, status.Error(codes.InvalidArgument, "content_type required")
	}

	reader := bytes.NewReader(req.Data)
	objectHexID, err := mongodb.UploadPhoto(uid, req.Filename, req.ContentType, reader, int64(len(req.Data)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "mongo upload failed: %v", err)
	}

	profile, err := database.UpdateProfilePhoto(uid, objectHexID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return toPbProfile(profile), nil
}

func (s *UserService) ChangeSettings(ctx context.Context, req *pb.ChangeSettingsRequest) (*pb.Profile, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var profile *database.Profile
	if req.NotificationsEnabled {
		profile, err = database.EnableNotifications(uid)
	} else {
		profile, err = database.DisableNotifications(uid)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return toPbProfile(profile), nil
}

func (s *UserService) ChangeSubscriptions(ctx context.Context, req *pb.ChangeSubscriptionsRequest) (*pb.ChangeSubscriptionsResponse, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := database.UpdateProfileSubscriptions(uid, req.Subscriptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.ChangeSubscriptionsResponse{
		Profile: toPbProfile(profile),
	}, nil
}

func (s *UserService) AdminGetProfile(ctx context.Context, req *pb.AdminGetProfileRequest) (*pb.Profile, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	role, err := database.GetUserRole(uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}
	if role != "admin" {
		return nil, status.Errorf(codes.PermissionDenied, "admin only")
	}

	profile, err := database.GetProfile(req.UserId)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return toPbProfile(profile), nil
}

func (s *UserService) AdminGetProfilesList(ctx context.Context, req *pb.AdminGetProfilesListRequest) (*pb.AdminGetProfilesListResponse, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	role, err := database.GetUserRole(uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}
	if role != "admin" {
		return nil, status.Errorf(codes.PermissionDenied, "admin only")
	}

	profiles, err := database.GetAllProfiles()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	res := make([]*pb.Profile, 0, len(profiles))
	for _, p := range profiles {
		pp := p
		res = append(res, toPbProfile(&pp))
	}

	return &pb.AdminGetProfilesListResponse{Profiles: res}, nil
}

func (s *UserService) AdminBanUser(ctx context.Context, req *pb.AdminBanUserRequest) (*pb.AdminBanUserResponse, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return &pb.AdminBanUserResponse{Banned: false}, err
	}

	role, err := database.GetUserRole(uid)
	if err != nil {
		return &pb.AdminBanUserResponse{Banned: false}, status.Errorf(codes.Internal, "db error: %v", err)
	}
	if role != "admin" {
		return &pb.AdminBanUserResponse{Banned: false}, status.Errorf(codes.PermissionDenied, "admin only")
	}

	if res := database.BanUser(req.UserId); res != nil {
		return &pb.AdminBanUserResponse{Banned: false}, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.AdminBanUserResponse{Banned: true}, nil
}

func (s *UserService) AdminUnbanUser(ctx context.Context, req *pb.AdminUnbanUserRequest) (*pb.AdminUnbanUserResponse, error) {
	uid, err := authmw.GetUserIDFromCtx(ctx)
	if err != nil {
		return &pb.AdminUnbanUserResponse{Unbanned: false}, err
	}

	role, err := database.GetUserRole(uid)
	if err != nil {
		return &pb.AdminUnbanUserResponse{Unbanned: false}, status.Errorf(codes.Internal, "db error: %v", err)
	}
	if role != "admin" {
		return &pb.AdminUnbanUserResponse{Unbanned: false}, status.Errorf(codes.PermissionDenied, "admin only")
	}

	if res := database.UnbanUser(req.UserId); res != nil {
		return &pb.AdminUnbanUserResponse{Unbanned: false}, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.AdminUnbanUserResponse{Unbanned: true}, nil
}

func main() {
	database.InitDB()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(authmw.UnaryAuth()))
	pb.RegisterUserServiceServer(s, &UserService{})
	reflection.Register(s)

	log.Println("UserService started on port 50052...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
