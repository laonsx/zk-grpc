package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"zk-grpc/zookeeper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/resolver"
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

	resolver.Register(zookeeper.NewGrpcResolver("localhost:2181", "zk"))

	roundrobinConn, err := grpc.Dial(zookeeper.GrpcDialUrl, grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure())
	if err != nil {

		log.Fatalf("did not connect: %v", err)
	}

	defer roundrobinConn.Close()

	hwcc := ecpb.NewEchoClient(roundrobinConn)

	makeRPCs(hwcc, 150)
}
