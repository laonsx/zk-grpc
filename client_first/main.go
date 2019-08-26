package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"zk-grpc/zookeeper"

	"google.golang.org/grpc"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
)

func callUnaryEcho(c ecpb.EchoClient, message string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {

		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println(r.Message)
}

func makeRPCs(hwc ecpb.EchoClient, n int) {

	for i := 0; i < n; i++ {

		time.Sleep(time.Second)
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}

func main() {

	pickfirstConn, err := grpc.Dial(zookeeper.InitGrpcDialUrl("localhost:2181", "zk"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {

		log.Fatalf("did not connect: %v", err)
	}

	defer pickfirstConn.Close()

	hwc := ecpb.NewEchoClient(pickfirstConn)

	makeRPCs(hwc, 100)
}
