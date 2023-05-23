package main

import (
	"context"
	"github.com/weilanjin/go-grpc/api/v1/middleware"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

type Credential struct{}

func (c *Credential) GetRequestMetadata(ctx context.Context, url ...string) (map[string]string, error) {
	return map[string]string{
		"token": "xxxxx",
	}, nil
}

func (c *Credential) RequireTransportSecurity() bool {
	return false
}

// grpc.UnaryClientInterceptor
// 1. set token
// 2. logging
// use grpc.WithUnaryInterceptor(interceptorClient)
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
		grpc.WithPerRPCCredentials(new(Credential)),
	}
	conn, err := grpc.Dial("localhost:10007", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := middleware.NewGreeterClient(conn)

	// err DeadlineExceeded
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3) // 设置请求超时时间
	rsp, err := client.SayHello(ctx, &middleware.HelloRequest{Greeting: "lanjin.wei"})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			panic(err)
		}
		slog.Info("[client]", "code", st.Code(), "message", st.Message())
		return
	}
	slog.Info("[client]", "reply", rsp.Reply)
}
