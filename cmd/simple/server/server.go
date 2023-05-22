package main

import (
	"github.com/weilanjin/go-grpc/api/v1/simple"
	"google.golang.org/grpc"
	"net"
)

const addr = ":10006"

func main() {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	simple.RegisterGreeterServer(s, new(GreeterServer))
	if err = s.Serve(listen); err != nil {
		panic(err)
	}
}
