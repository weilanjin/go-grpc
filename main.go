package main

import (
	"fmt"
	"github.com/weilanjin/go-grpc/api/v1/helloworld"
	"google.golang.org/protobuf/proto"
)

func main() {
	req := helloworld.HelloRequest{Name: "lanjin.wei"}
	bytes, _ := proto.Marshal(&req)
	fmt.Println(string(bytes))
}
