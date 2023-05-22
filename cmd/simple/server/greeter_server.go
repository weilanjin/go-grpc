package main

import (
	"context"
	"fmt"
	"github.com/weilanjin/go-grpc/api/v1/simple"
	"log"
	"sync"
	"time"
)

type GreeterServer struct {
	simple.UnimplementedGreeterServer
}

func (svc *GreeterServer) SayHello(ctx context.Context, req *simple.HelloRequest) (*simple.HelloResponse, error) {
	return &simple.HelloResponse{
		Reply: fmt.Sprintf("[server] %s", req.Greeting),
	}, nil
}

func (svc *GreeterServer) LotsOfReplies(req *simple.HelloRequest, serverStream simple.Greeter_LotsOfRepliesServer) error {
	for i := 0; i < 10; i++ {
		err := serverStream.Send(&simple.HelloResponse{
			Reply: fmt.Sprintf("[server] %s %s 实时价格 %d", time.Now().Format(time.DateTime), req.Greeting, i+10),
		})
		if err != nil {
			log.Printf("[server] ERROR %s", err.Error())
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (svc *GreeterServer) LotsOfGreetings(cliStream simple.Greeter_LotsOfGreetingsServer) error {
	for {
		recv, err := cliStream.Recv()
		if err != nil {
			log.Printf("[server] ERROR %s", err.Error())
			return err
		}
		log.Printf(recv.Greeting)
	}
}

func (svc *GreeterServer) BidiHello(stream simple.Greeter_BidiHelloServer) error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				log.Printf("[server] ERROR %s", err.Error())
				wg.Done()
				return
			}
			log.Printf("[server] %s", recv.Greeting)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			err := stream.Send(&simple.HelloResponse{
				Reply: "helle!",
			})
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
	return nil
}
