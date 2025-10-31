

1. Open cmd prompt 

   go get -u google.golang.org/grpc   ==>  GRPC 

2. Open cmd prompt  

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
3. Download and unzip (for Windows) :


   https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip  :old 

   In download u will get a zip folder

4. Gotopath in your PC  => C:\Program Files\protoc
   create protoc folder in cprogram files
   extract all files in this folder             ===> protoc compiler

5. Go To EnvironmentalVariable variable value in system path add it 
   Add C:\Program Files\protoc\bin


cd C:\Users\Administrator\Desktop
mkdir grpc_hello
cd grpc_hello
go mod init github.com/administrator/grpc_hello

  C:\Users\Administrator\Desktop\grpc_hello\hellopb\greet.proto  

syntax = "proto3";

option go_package = "github.com/administrator/grpc_hello/hellopb;hellopb";

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}

service Greeter {
  rpc SayHello(HelloRequest) returns (HelloResponse);
}


cd C:\Users\Administrator\Desktop\grpc_hello\server
go run main.go

package main

import (
"context"
"fmt"
"log"
"net"

"google.golang.org/grpc"
"github.com/administrator/grpc_hello/hellopb"
)

// server implements the Greeter service defined in greet.proto
type server struct {
hellopb.UnimplementedGreeterServer
}

// Implement the SayHello RPC
func (s *server) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
name := req.GetName()
message := fmt.Sprintf("Hello, %s! 👋 Welcome to gRPC in Go.", name)
return &hellopb.HelloResponse{Message: message}, nil
}

func main() {
// Start listening on port 50051
lis, err := net.Listen("tcp", ":50051")
if err != nil {
log.Fatalf("Failed to listen: %v", err)
}

grpcServer := grpc.NewServer()
hellopb.RegisterGreeterServer(grpcServer, &server{})

log.Println("🚀 gRPC server started on port 50051...")
if err := grpcServer.Serve(lis); err != nil {
log.Fatalf("Failed to serve: %v", err)
}
}


cd C:\Users\Administrator\Desktop\grpc_hello\client
go run main.go

C:\Users\Administrator\Desktop\grpc_hello\client\main.go

package main

import (
"context"
"log"
"time"

"google.golang.org/grpc"
"github.com/administrator/grpc_hello/hellopb"
)

func main() {
// Connect to server
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

client := hellopb.NewGreeterClient(conn)

// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

// Make RPC call
resp, err := client.SayHello(ctx, &hellopb.HelloRequest{Name: "CoderRange"})
if err != nil {
log.Fatalf("Error calling SayHello: %v", err)
}

log.Printf("✅ Response from server: %s", resp.GetMessage())
}
