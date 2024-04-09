# GrpcExample

This repository contains an example of a service defined with gRPC, it's my playground to explore gRPC implementation in different languages.

As of now it contains a server written in Go, and two clients, One for go and second for c++

DISCLAIMER:  it's just a showcase, the project doesn't necessarily follow all the best practices.

## Requirements

Minimum:

* make;
* [Go 1.8.3](https://golang.org/dl/);
* [gRPC-go](https://github.com/golang/protobuf/);
* pkg-config;
* [Protobuf 3.3.2](https://github.com/protocolbuffers/protobuf/releases).

Requirements for the command line client:

* [gRPC-c++](https://grpc.io/docs/quickstart/cpp.html).

## Project structure

```dir
- client
    |- client.go
    |- client.cpp
- proto
    |- myservice.proto
- server
    |- server.go
```


## Build

To build the protobuf, go to the proto folder and run this command:
```sh
protoc -I. --grpc_out=. --plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin myservice.proto
protoc --cpp_out=. myservice.proto
protoc --go_out=. --go-grpc_out=. myservice.proto
```

To build the server, go to the server folder and run this command:

```sh
go run main.go
```
To build the client, go to the client folder and run this command only for go client:

```sh
go run main.go
```

To build the client, go to the client folder and run this command only for c++ client:

```sh
g++ -std=c++17 client.cpp ../proto/myservice.pb.cc ../proto/myservice.grpc.pb.cc -o client `pkg-config --libs --cflags protobuf grpc++`
./client
```