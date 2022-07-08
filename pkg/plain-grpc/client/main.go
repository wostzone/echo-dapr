package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

// Demonstrate how a client sends an echo message to a grpc service.
// Use-case: develop services that are independent of dapr but can be invoked by dapr via dapr
// sidecars using gRPC. The sidecars in B and C are loaded separately from the service
//
//  A: client -[grpc]-> service
//  B: client -[grpc]-> service-sidecar -[grpc]-> service
//  C: client -[grpc]-> client-sidecar -[grpc]-> service-sidecar -[grpc]-> service
func main() {
	var text = "Hello echo"
	var appID = pkg.EchoServiceAppID
	var cmd string
	var port int
	var repeat int = 1
	flag.IntVar(&port, "port", pkg.EchoServiceGrpcPort, "Service gRPC listening port")
	flag.StringVar(&appID, "app-id", pkg.EchoServiceAppID, "Service name when using dapr")
	flag.IntVar(&repeat, "repeat", repeat, "Nr of times to invoke")
	flag.Parse()
	values := flag.Args()
	if len(values) == 1 && values[0] == "stop" {
		cmd = values[0]
	} else if len(values) == 2 {
		cmd = values[0]
		text = values[1]
	} else {
		fmt.Println("Missing text: invoke-grpc-client <command> <text>: ", cmd)
		flag.Usage()
		return
	}

	InvokeGrpcService(port, appID, cmd, text, repeat)
}

// InvokeGrpcService invokes the service using grpc
func InvokeGrpcService(port int, appID string, cmd string, text string, repeat int) {
	var response *pb.TextParam
	listenAddress := fmt.Sprintf(":%d", port)

	// Set up a connection to the server.
	fmt.Println("Connecting to service '"+appID+"' on", listenAddress)
	conn, err := grpc.Dial(listenAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// The service name to connect to when connecting via dapr
	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", appID)

	t1 := time.Now()
	for count := 0; count < repeat; count++ {
		if cmd == "echo" {
			response, err = c.Echo(ctx, &pb.TextParam{Text: text + "-" + strconv.Itoa(count)})
		} else if cmd == "upper" {
			response, err = c.UpperCase(ctx, &pb.TextParam{Text: text + "-" + strconv.Itoa(count)})
		} else if cmd == "reverse" {
			response, err = c.Reverse(ctx, &pb.TextParam{Text: text + "-" + strconv.Itoa(count)})
		} else if cmd == "stop" {
			response, err = c.Stop(ctx, &empty.Empty{})
			os.Exit(0)
		} else {
			response, err = c.Echo(ctx, &pb.TextParam{Text: cmd + " - " + text + "-" + strconv.Itoa(count)})
		}
		if err != nil {
			log.Fatalf("Could not echo text: %v", err)
		} else if response != nil {
			//fmt.Println("Response: ", response.GetText())
		}
	}
	t2 := time.Now()
	duration := t2.Sub(t1)
	fmt.Println(fmt.Sprintf("Time to invoke %d grpc calls: %d msec", repeat, duration.Milliseconds()))
}
