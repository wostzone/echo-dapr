// Package main CLI to communicate with the directory service via dapr
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/wostzone/echo/proto/go"
)

// ServiceName used by the service when connecting via dapr
const ServiceName = "echo"

const (
	// the default dapr instance
	//address = "localhost:40001"
	//daprHttpAddress    = "localhost:9000"
	defaultGrpcPort = "40001"
)

// Client used to invoke the grpc echo service
func main() {
	var text = "Hello echo"
	var grpcServicePort = "9001"
	var cmd string
	flag.StringVar(&grpcServicePort, "port", defaultGrpcPort, "gRPC port to connect to the echo service")
	flag.Parse()
	values := flag.Args()
	if len(values) == 0 {
		fmt.Println("Missing text: echo-cli <command> <text>: ", cmd)
		flag.Usage()
		return
	} else if len(values) == 1 && values[0] == "stop" {
		cmd = values[0]
	} else {
		cmd = values[0]
		text = values[1]
	}

	directGrpc(grpcServicePort, cmd, text)
	//usingInvoke()
}

//func usingInvoke() {
//	client, err := dapr.NewClientWithAddress(daprGrpcAddress)
//	if err != nil {
//		panic(err)
//	}
//	defer client.Close()
//
//	ctx := context.Background()
//	//Using Dapr SDK to invoke a method
//	result, err := client.InvokeMethod(ctx, "directory", "ListTDs", "get")
//	if err != nil {
//		log.Fatalf("could not list the directory: %v", err)
//	}
//	log.Println("Result: %v", result)
//}

func directGrpc(port string, cmd string, text string) {
	var response *pb.TextParam
	// Set up a connection to the server.
	fmt.Println("Connecting to echo service on :" + port)
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// The service name to connect to when connecting via dapr
	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", ServiceName)

	if cmd == "echo" {
		response, err = c.Echo(ctx, &pb.TextParam{Text: text})
	} else if cmd == "upper" {
		response, err = c.UpperCase(ctx, &pb.TextParam{Text: text})
	} else if cmd == "reverse" {
		response, err = c.Reverse(ctx, &pb.TextParam{Text: text})
	} else if cmd == "stop" {
		response, err = c.Stop(ctx, &empty.Empty{})
	} else {
		response, err = c.Echo(ctx, &pb.TextParam{Text: cmd + " - " + text})
	}
	if err != nil {
		log.Fatalf("Could not echo text: %v", err)
	} else if response != nil {
		log.Printf(response.GetText())
	}
}
