package main

import (
	"context"
	"github.com/weilanjin/go-grpc/api/v1/middleware"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// grpc.UnaryClientInterceptor
// 1. set token
// 2. logging
func interceptorClient(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(map[string]string{
		"token": "xxxxxx",
	})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	slog.Info("[client]", "method", method, "request", req, "reply", reply)
	return err
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptorClient),
	}
	conn, err := grpc.Dial("localhost:10007", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := middleware.NewGreeterClient(conn)
	rsp, err := client.SayHello(context.Background(), &middleware.HelloRequest{Greeting: "lanjin.wei"})
	if err != nil {
		panic(err)
	}
	slog.Info("[client]", "reply", rsp.Reply)
}
