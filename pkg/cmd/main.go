package main

import (
	"net"

	"google.golang.org/grpc"
	grpcs "nononsensecode.com/grpc-tutorial/pkg/interfaces/grpc"
)

func main() {
	println("gRPC server tutorial")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpcs.WithServerInterceptor())
	grpcs.RegisterTutorialServer(s, &grpcs.Server{})
	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
