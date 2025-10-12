package resolvers

import (
	adminpb "ads/adminservice/presentation/grpc/pb"
	authpb "ads/authservice/presentation/grpc/pb"
	userpb "ads/userservice/presentation/grpc/pb"
)

type Resolver struct {
	AuthClient  authpb.AuthServiceClient
	UserClient  userpb.UsersServiceClient
	AdminClient adminpb.AdminServiceClient
}
