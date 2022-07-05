package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"

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

	service := daprd.NewService(":" + strconv.Itoa(port))

	// the invocation handler is http/grpc agnostic - otherwise, it is just a wrapper
	if err := service.AddServiceInvocationHandler("echo", echoInvocationHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	if err := service.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
	}
}

// echoInvocationHandler uses the dapr invocation API. It is just another way to handle request/response
func echoInvocationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	var args *pb.TextParam
	var response []byte
	log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
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
