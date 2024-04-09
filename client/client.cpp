#include <iostream>
#include <memory>
#include <string>
//#include <grpc/grpc.h>
#include <grpcpp/grpcpp.h>
#include "../proto/myservice.pb.h"
#include "../proto/myservice.grpc.pb.h"

using namespace std;
using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using example::GreetingRequest;
using example::GreetingReply;

int main() {
    shared_ptr<grpc::Channel> channel = grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials());
    //unique_ptr<example::Greeter::Stub> stub = example::Greeter::NewStub(channel);
    auto stub = example::Greeter::NewStub(channel);
    // std::unique_ptr<myservice::MyService::Stub> stub = myservice::MyService::NewStub(channel);
    example::GreetingRequest request;
    request.set_name("World");

    example::GreetingReply response;
    grpc::ClientContext context;

    grpc::Status status = stub->Hello(&context, request, &response);

    if (status.ok()) {
        cout << "Response: " << response.text() << endl;
    } else {
        cerr << "RPC failed: " << status.error_code() << ", " << status.error_message() << endl;
    }

    return 0;
}