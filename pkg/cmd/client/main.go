package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	grpcc "nononsensecode.com/grpc-tutorial/pkg/interfaces/grpc"
)

func main() {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		withDefInt(),
	}
	conn, err := grpc.Dial("localhost:8080", dialOpts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := grpcc.NewTutorialClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, grpcc.CtxCloudVendorKey, "aws")
	ctx = context.WithValue(ctx, grpcc.CtxClientIDKey, "1234")

	resp, err := client.SayHello(ctx, &grpcc.HelloRequest{Name: "John"})
	if err != nil {
		panic(err)
	}
	println(resp.Message)
}

func withDefInt() grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(clientIDInterceptor, cloudVendorInterceptor))
}

func withDefaultInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(clientIDInterceptor, cloudVendorInterceptor))
}

func withClientIDInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientIDInterceptor)
}

func withCloudVendorInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(cloudVendorInterceptor)
}

func clientIDInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	fmt.Println("clientIDInterceptor")
	clientID, ok := ctx.Value(grpcc.CtxClientIDKey).(string)
	if !ok {
		return errors.New("client id not found in context")
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Append("client-id", clientID)
	} else {
		md = metadata.Pairs("client-id", clientID)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func cloudVendorInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	fmt.Println("cloudVendorInterceptor")
	cloudVendor, ok := ctx.Value(grpcc.CtxCloudVendorKey).(string)
	if !ok {
		return errors.New("cloud vendor not found in context")
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Append("cloud-vendor", cloudVendor)
	} else {
		md = metadata.Pairs("cloud-vendor", cloudVendor)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
