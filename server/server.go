package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	pb "example.com/m/proto/proto" // Import the package containing the generated code
	"github.com/gorilla/mux"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// type GreetingRequest struct {
// 	Text string `json:"text"`
// }

// Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
type server struct {
	pb.UnimplementedGreeterServer
	//greetSuffix string
}

// func (s *server) Hello(ctx context.Context, msg *pb.GreetingRequest) (*Message, error) {
// 	log.Printf("Received message: %s", msg.Text)
// 	return &Message{Text: "Message received successfully"}, nil
// }

func (s *server) Hello(ctx context.Context, in *pb.GreetingRequest) (*pb.GreetingReply, error) {
	// Implement your method logic here
	log.Println("Hello called")
	reply, err := getReply(ctx, in.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.GreetingReply{Text: reply}, nil
}

func main() {
	fmt.Println("gPRC server has started....")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//suffix := "nothing... "
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})
	log.Println("Server started on port 50051...")
	// if err := grpcServer.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	// }

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start REST server
	r := mux.NewRouter()
	r.HandleFunc("/sayhello", sayHelloREST).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func sayHelloREST(w http.ResponseWriter, r *http.Request) {
	// First, grab the arguments from the request body

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("reading request json body: %s", err.Error()), http.StatusInternalServerError)
	}
	defer r.Body.Close()

	var req pb.GreetingRequest
	err = protojson.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Next, get the reply from the shared helper function

	reply, err := getReply(r.Context(), req.GetName())
	if err != nil {
		code, message := convertGRPCErrorToHTTPStatus(err)
		http.Error(w, message, code)
		return
	}

	// Finally, prepare the response and send it

	resp := &pb.GreetingReply{Text: reply}
	jsonBytes, err := protojson.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// getReply is a helper function that we can reuse in both the gRPC and REST APIs
// so that we don't have to duplicate the implementation logic.
func getReply(_ context.Context, name string) (message string, err error) {
	if name == "" {
		return "", status.Error(codes.InvalidArgument, "name was not provided")
	}

	return fmt.Sprintf("Hello, %s!", name), nil
}

// convertGRPCErrorToHTTPStatus translates gRPC error codes to HTTP status codes. See
// https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md
// for more information.
func convertGRPCErrorToHTTPStatus(err error) (httpCode int, errorText string) {
	s, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, err.Error()
	}

	switch s.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, s.Message()
	default:
		return http.StatusInternalServerError, s.Message()
	}
}
