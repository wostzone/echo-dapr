// Package main that launches the directory service on localhost
package main

import (
	"flag"
	"net"

	dapr "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/wostzone/echo/internal"
	pb "github.com/wostzone/echo/proto/go"
)

// Entry point of the gRPC echo service
func main() {
	grpcServicePort := *flag.String("port", "40001", "Service gRPC listening port")
	flag.Parse()

	lis, err := net.Listen("tcp", "localhost:"+grpcServicePort)
	if err != nil {
		logrus.Fatalf("Failed open listener: %v", err)
	}

	// listen on gRPC
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	echoService := internal.NewEchoService(grpcServer)
	// the service implements both the Echo API and the dapr callback api
	pb.RegisterEchoServiceServer(grpcServer, echoService)
	dapr.RegisterAppCallbackServer(grpcServer, echoService)

	logrus.Infof("Directory service listening on gRPC port: %s", grpcServicePort)
	//go func() {
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
	//}()

}
