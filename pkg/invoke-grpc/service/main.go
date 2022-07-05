package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"

	"github.com/wostzone/echo/internal"
	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

// Entry point of the grpc echo service
func main() {
	var port int
	flag.IntVar(&port, "port", pkg.EchoDaprServiceGrpcPort, "Service grpc listening port")
	flag.Parse()

	StartGrpcService(port)
}

// StartGrpcService is how a grpc service would handle requests minus auth, tracing, etc
func StartGrpcService(port int) {
	fmt.Println("Starting echo-service on grpc port ", port)

	// Note: you can't use your own gRPC service with this approach.
	//var opts []grpc.ServerOption
	//
	//grpcServer := grpc.NewServer(opts...)
	//echoService := internal.NewEchoService(func() {
	//	grpcServer.Stop()
	//})
	//lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	//if err != nil {
	//	logrus.Fatalf("Failed open listener: %v", err)
	//}
	//
	//// register a gRPC service that can be invoked by any gRPC client including dapr
	//pb.RegisterEchoServiceServer(grpcServer, echoService)
	//
	//go func() {
	//	if err := grpcServer.Serve(lis); err != nil {
	//		logrus.Fatalf("failed to serve: %v", err)
	//	}
	//}()
	//fmt.Println("gRPC service listing on port ", port)

	//
	appcallbackServer, err := daprd.NewService(":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to create app callback server: %v", err)
	}

	// the invocation handler is http/grpc agnostic - otherwise, it is just a wrapper
	if err := appcallbackServer.AddServiceInvocationHandler("stop", stopInvocationHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}
	// the invocation handler is http/grpc agnostic - otherwise, it is just a wrapper
	if err := appcallbackServer.AddServiceInvocationHandler("echo", echoInvocationHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}
	if err := appcallbackServer.AddServiceInvocationHandler("upper", upperInvocationHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	if err := appcallbackServer.Start(); err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("dapr appCallbackService listing on port ", port)
}

// stopInvocationHandler uses the dapr invocation API. It is just another way to handle request/response
func stopInvocationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	go os.Exit(0)
	out = &common.Content{
		Data:        []byte("bye bye"),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}

// echoInvocationHandler uses the dapr invocation API. It is just another way to handle request/response
func echoInvocationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	var args *pb.TextParam
	var response []byte
	log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	// shouldn't this be protobuf encoded?
	err = json.Unmarshal(in.Data, &args)
	if err != nil {
		err := fmt.Errorf("Error unmarshalling payload for handleEcho: %s", err)
		return nil, err
	}
	echoService := internal.NewEchoService(nil)
	result, err := echoService.Echo(nil, args)
	if err != nil {
		err = fmt.Errorf("Error handling echo request: %s", err)
		return nil, err
	} else {
		response, _ = json.Marshal(result)
	}

	// do something with the invocation here
	out = &common.Content{
		Data:        response,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}

func upperInvocationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	var args *pb.TextParam
	var response []byte
	log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	// shouldn't this be protobuf encoded?
	err = json.Unmarshal(in.Data, &args)
	if err != nil {
		err := fmt.Errorf("Error unmarshalling payload for handleEcho: %s", err)
		return nil, err
	}
	echoService := internal.NewEchoService(nil)
	result, err := echoService.UpperCase(nil, args)
	if err != nil {
		err = fmt.Errorf("Error handling echo request: %s", err)
		return nil, err
	} else {
		response, _ = json.Marshal(result)
	}

	// do something with the invocation here
	out = &common.Content{
		Data:        response,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
