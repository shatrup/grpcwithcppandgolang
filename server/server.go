package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "example.com/m/proto/proto" // Import the package containing the generated code

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
	//greetSuffix string
}

func (s *server) Hello(_ context.Context, _ *pb.GreetingRequest) (*pb.GreetingReply, error) {
	// Implement your method logic here
	log.Println("Hello called")
	return &pb.GreetingReply{Text: fmt.Sprintf("Hello ")}, nil
}

func main() {
	fmt.Println("gPRC server has started....")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//suffix := "nothing... "
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Println("Server started on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
