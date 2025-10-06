package resolvers

import (
	adminpb "AdsService/adminservice/presentation/grpc/pb"
	authpb "AdsService/authservice/presentation/grpc/pb"
	userpb "AdsService/userservice/presentation/grpc/pb"
)

type Resolver struct {
	AuthClient   authpb.AuthServiceClient
	UserClient   userpb.UsersServiceClient
	AdminService adminpb.AdminServiceClient
}
