package main

import (
	"context"
	"fmt"
	"github.com/weilanjin/go-grpc/api/v1/simple"
	"github.com/weilanjin/go-grpc/pkg/util"
	"log"
	"sync"
	"time"
)

type GreeterClient struct {
	Client simple.GreeterClient
}

func (clt *GreeterClient) SayHello() {
	rsp, err := clt.Client.SayHello(context.Background(), &simple.HelloRequest{Greeting: "hello lanjin.wei"})
	if err != nil {
		panic(err)
	}
	log.Println(rsp.Reply)
}

func (clt *GreeterClient) LotsOfReplies() {
	serverStream, err := clt.Client.LotsOfReplies(context.Background(), &simple.HelloRequest{Greeting: "恒大集团"})
	if err != nil {
		panic(err)
	}
	for {
		rsp, err := serverStream.Recv()
		if err != nil {
			log.Printf("[client] ERROR %s", err.Error())
			break
		}
		log.Println(rsp.Reply)
	}
}

func (clt *GreeterClient) LotsOfGreetings() {
	clientStream, err := clt.Client.LotsOfGreetings(context.Background())
	if err != nil {
		panic(err)
	}
	randNum := util.RandNumber(19, 26)
	for i := 0; i < 10; i++ {
		err := clientStream.Send(&simple.HelloRequest{
			Greeting: fmt.Sprintf("[client] 大棚室内 %s 温度 %d", time.Now().Format(time.DateTime), randNum()),
		})
		if err != nil {
			log.Printf("[client] ERROR %s", err.Error())
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (clt *GreeterClient) BidiHello() {
	stream, err := clt.Client.BidiHello(context.Background())
	if err != nil {
		log.Printf("[client] ERROR %s", err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				log.Printf("[client] ERROR %s", err.Error())
				wg.Done()
				return
			}
			log.Printf("[client] %s", recv.Reply)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			err := stream.Send(&simple.HelloRequest{Greeting: "Hi!"})
			if err != nil {
				log.Printf("[server] ERROR %s", err.Error())
				wg.Done()
				return
			}
			time.Sleep(2 * time.Second)
		}
		wg.Done()
	}()
	wg.Wait()
}
