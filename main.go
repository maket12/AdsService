package main

import (
	"AdsService/infra/authmw"
	"AdsService/userservice"
	pb "AdsService/userservice/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(authmw.UnaryAuth()))
	pb.RegisterUserServiceServer(s, &userservice.UserService{}) // теперь ок
	reflection.Register(s)

	log.Println("userservice started on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
