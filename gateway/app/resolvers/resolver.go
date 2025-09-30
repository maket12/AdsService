package resolvers

import (
	authpb "AdsService/authservice/presentation/grpc/pb"
	userpb "AdsService/userservice/presentation/grpc/pb"
)

type Resolver struct {
	AuthClient authpb.AuthServiceClient
	UserClient userpb.UserServiceClient
}
