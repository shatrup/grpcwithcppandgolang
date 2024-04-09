package main

import (
	"context"
	"log"

	pb "example.com/m/proto/proto" // Import the package containing the generated code
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response
	r, err := c.Hello(context.Background(), &pb.GreetingRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", r.Text)
}
