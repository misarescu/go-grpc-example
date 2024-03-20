package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc-chat.com/pkg/chat/proto"
)

var port = flag.Int("port", 8080, "server port to listen to")

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Recieved: %v", in.GetName())
	result := "Hello " + in.GetName() + " "
	for i := int32(0); i < in.GetRepeats(); i++ {
		result += in.GetName()
		result += " "
	}
	log.Printf("Output: %v", result)
	return &pb.HelloReply{Message: result}, nil
}

func main() {
	log.Printf("Starting server")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("listening on %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
