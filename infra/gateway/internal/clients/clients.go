package clients

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "AdsService/authservice/proto"
	userpb "AdsService/userservice/proto"
)

type Clients struct {
	Auth authpb.AuthServiceClient
	User userpb.UserServiceClient
}

func New(ctx context.Context) *Clients {
	dial := func(addr string) *grpc.ClientConn {
		conn, err := grpc.DialContext(ctx, addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("couldn't connect to %s: %v", addr, err)
			return nil
		}
		log.Printf("succesfully connected to %s:", addr)
		return conn
	}

	clients := &Clients{}

	a := dial("authservice:50051")
	if a != nil {
		clients.Auth = authpb.NewAuthServiceClient(a)
	}

	u := dial("userservice:50052")
	if u != nil {
		clients.User = userpb.NewUserServiceClient(u)
	}

	return clients
}
