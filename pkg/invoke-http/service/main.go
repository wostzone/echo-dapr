package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/gorilla/mux"

	"github.com/wostzone/echo/internal"
	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

// Entry point of the http echo service
func main() {
	var port int
	flag.IntVar(&port, "port", pkg.EchoServiceHttpPort, "Service http listening port")
	flag.Parse()

	StartHttpService(port)
}

// StartHttpService is how a regular http service would handle requests minus auth, tracing, etc
func StartHttpService(port int) {
	fmt.Println("Starting echo-service on http port ", port)
	r := mux.NewRouter()
	//r.HandleFunc("/echo", handleEcho).Methods("POST")
	r.HandleFunc("/upper", handleUpper).Methods("POST")
	r.HandleFunc("/reverse", handleReverse).Methods("POST")

	service := daprd.NewServiceWithMux(":"+strconv.Itoa(port), r)
	//service := daprd.NewService(":" + strconv.Itoa(port))

	// the invocation handler is http/grpc agnostic
	if err := service.AddServiceInvocationHandler("stop", stopInvocationHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}
	if err := service.AddServiceInvocationHandler("echo", echoInvocationHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	if err := service.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
	}
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
	//log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
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

// http style handler
func handleReverse(w http.ResponseWriter, r *http.Request) {
	//log.Println("handleReverse: Received 'reverse' request over http")
	var args *pb.TextParam
	data, err := ioutil.ReadAll(r.Body)

	json.Unmarshal(data, &args)
	echoService := internal.NewEchoService(nil)
	result, err := echoService.Reverse(nil, args)
	if err != nil {
		log.Println("Error handling reverse request:", err)
		w.WriteHeader(500)
	} else {
		response, _ := json.Marshal(result)
		w.Write(response)
	}
}

func handleUpper(w http.ResponseWriter, r *http.Request) {
	//log.Println("handleUpper: Received 'upper' request over http")

	var args *pb.TextParam
	data, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &args)
	echoService := internal.NewEchoService(nil)
	result, err := echoService.UpperCase(nil, args)
	if err != nil {
		log.Println("Error handling upper request: ", err)
		w.WriteHeader(500)
	} else {
		response, _ := json.Marshal(result)
		w.Write(response)
	}
}
