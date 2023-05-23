package main

import (
	"context"
	"fmt"
	"github.com/weilanjin/go-grpc/api/v1/middleware"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

type Server struct {
	middleware.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *middleware.HelloRequest) (*middleware.HelloResponse, error) {
	// timeout test
	time.Sleep(5 * time.Second)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "get metadata error")
	}
	userInfo, ok := md["userInfo"]
	if !ok {
		slog.Error("[server] get user info fail", "request", request.Greeting)
	}
	slog.Info("[server] %s", "userInfo", userInfo, "request", request.Greeting)
	return &middleware.HelloResponse{Reply: "Hi"}, nil
}

// grpc.UnaryServerInterceptor
// 1.耗时统计
// 2.验证token
func interceptorServer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	slog.Info("[server]", "serverInfo", *info)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "get metadata error")
	}
	token, ok := md["token"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "token not found")
	}
	slog.Info("[server]", "token", token)

	startTime := time.Now()
	res, err := handler(ctx, req)
	time.Sleep(300 * time.Millisecond)
	d := time.Since(startTime)
	slog.Info("[server]", "delay", d.String())
	return res, err
}

var addr = ":10007"

func main() {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	} else {
		slog.Info(fmt.Sprintf("[server] listen: %s", addr))
	}
	defer listen.Close()
	opts := []grpc.ServerOption{grpc.UnaryInterceptor(interceptorServer)}
	s := grpc.NewServer(opts...)
	defer s.GracefulStop() // 优雅停止

	middleware.RegisterGreeterServer(s, &Server{})
	if err = s.Serve(listen); err != nil {
		panic(err)
	}
}
