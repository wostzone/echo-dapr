// Package main that launches the directory service on localhost
package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/wostzone/echo/internal"
	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

// Entry point of the gRPC echo service
func main() {
	var port int
	flag.IntVar(&port, "port", pkg.EchoServiceGrpcPort, "Service gRPC listening port")
	flag.Parse()

	StartGrpcService(port)
}

// StartGrpcService starts the service listening with grpc
func StartGrpcService(port int) {
	fmt.Println("StartGrpcService starting echo-service on grpc port ", port)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	echoService := internal.NewEchoService(func() {
		grpcServer.Stop()
	})

	// register a gRPC service that can be invoked by any gRPC client including dapr
	pb.RegisterEchoServiceServer(grpcServer, echoService)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		logrus.Fatalf("Failed open listener: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
