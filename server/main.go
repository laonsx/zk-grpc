package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zk-grpc/zookeeper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/status"
)

type ecServer struct {
	addr string
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *ecpb.EchoRequest) (*ecpb.EchoResponse, error) {

	return &ecpb.EchoResponse{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func (s *ecServer) ServerStreamingEcho(*ecpb.EchoRequest, ecpb.Echo_ServerStreamingEchoServer) error {

	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *ecServer) ClientStreamingEcho(ecpb.Echo_ClientStreamingEchoServer) error {

	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *ecServer) BidirectionalStreamingEcho(ecpb.Echo_BidirectionalStreamingEchoServer) error {

	return status.Errorf(codes.Unimplemented, "not implemented")
}

func startServer(addr string, sleep time.Duration) {

	time.Sleep(sleep)

	fmt.Println("regitster", zookeeper.Register("localhost:2181", "zk", addr))
	lis, err := net.Listen("tcp", addr)
	if err != nil {

		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	ecpb.RegisterEchoServer(s, &ecServer{addr: addr})

	log.Printf("serving on %s\n", addr)

	if err := s.Serve(lis); err != nil {

		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {

	for i := 0; i < 10; i++ {
		addr := fmt.Sprintf("127.0.0.1:%d", 50000+i)
		go startServer(addr, time.Duration(15*i)*time.Second)
	}

	defer zookeeper.UnRegister()

	handleSignal()
}

func handleSignal() {

	sigstop := syscall.Signal(15)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigstop, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
