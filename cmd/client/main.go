package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-chat.com/pkg/chat/proto"
)

const defaultName = "world"

var (
	port    = flag.Int("port", 8080, "port to listen to")
	addr    = flag.String("addr", "localhost", "the address to connect to")
	name    = flag.String("name", defaultName, "Name to greet")
	repeats = flag.Int("repeats", 0, "how may times to repeat the name")
)

func main() {
	flag.Parse()
	fullAddr := fmt.Sprintf("%s:%d", *addr, *port)
	conn, err := grpc.Dial(fullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name, Repeats: int32(*repeats)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
