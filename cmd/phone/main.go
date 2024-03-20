package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc-chat.com/pkg/phone/proto"
)

var port = flag.Int("port", 8080, "server port to listen to")

type server struct {
	pb.UnimplementedPhoneServiceServer
}

func (s *server) VoiceCall(ctx context.Context, phoneNumber *pb.Person_PhoneNumber) (*pb.CallResponse, error) {
	log.Printf("recieved: %v", phoneNumber)
	return &pb.CallResponse{
		Reciever: &pb.Person{
			Name:  "Marian",
			Id:    1,
			Email: "mail@marian.com",
			Phones: []*pb.Person_PhoneNumber{
				{Number: "074112345", Type: pb.PhoneType_PHONE_TYPE_HOME.Enum()}}},
		CallStatus: &pb.ErrorStatus{Message: "Call ok"}}, nil
}

func main() {
	log.Printf("Starting server")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterPhoneServiceServer(s, &server{})
	log.Printf("listening on %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
