package grpc

import (
	context "context"
	"errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	UnimplementedTutorialServer
}

func (s *Server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	clientID, ok := ctx.Value(CtxClientIDKey).(string)
	if ok {
		println("ClientID:", clientID)
	} else {
		println("No clientID")
	}
	CloudVendor, ok := ctx.Value(CtxCloudVendorKey).(string)
	if ok {
		println("CloudVendor:", CloudVendor)
	} else {
		println("No CloudVendor")
	}
	return &HelloReply{Message: "Hello " + in.GetName()}, nil
}

func WithServerInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(clientIDServerInterceptor, cloudVendorServerInterceptor))
}

func cloudVendorServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not present in context")
	}
	cloudVendor := md.Get("cloud-vendor")
	if len(cloudVendor) == 0 {
		return nil, errors.New("cloud vendor not found in metadata")
	}
	ctx = context.WithValue(ctx, CtxCloudVendorKey, cloudVendor[0])
	return handler(ctx, req)
}

func clientIDServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not present in context")
	}
	clientID := md.Get("client-id")
	if len(clientID) == 0 {
		return nil, errors.New("client id not found in metadata")
	}
	ctx = context.WithValue(ctx, CtxClientIDKey, clientID[0])
	return handler(ctx, req)
}

type ClientID string

var CtxClientIDKey ClientID = "clientID"

type CloudVendor string

var CtxCloudVendorKey CloudVendor = "cloudVendor"
