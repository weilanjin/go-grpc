package main

import (
	"github.com/weilanjin/go-grpc/api/v1/simple"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:10006", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := simple.NewGreeterClient(conn)
	gClient := GreeterClient{Client: client}

	gClient.SayHello() // 普通模式

	// 流模式

	gClient.LotsOfReplies()   // 服务端流模式
	gClient.LotsOfGreetings() // 客户端流模式
	gClient.BidiHello()       // 双向流模式
}
