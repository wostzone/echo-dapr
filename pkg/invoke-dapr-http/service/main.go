package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	// what can invocation handler do that mux doesn't
	if err := service.AddServiceInvocationHandler("/echo", echoHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	if err := service.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
	}
}

func echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
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

//func handleEcho(w http.ResponseWriter, r *http.Request) {
//	var args *pb.TextParam
//	data, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Println("Error reading request body handling echo request: ", err)
//		w.WriteHeader(500)
//		return
//	}
//	err = json.Unmarshal(data, &args)
//	if err != nil {
//		log.Println("Error unmarshalling payload for handleEcho: ", err)
//		w.WriteHeader(500)
//		return
//	}
//	log.Println("handleEcho: Received 'echo' request over http")
//	echoService := internal.NewEchoService(nil)
//	result, err := echoService.Echo(nil, args)
//	if err != nil {
//		log.Println("Error handling echo request: ", err)
//		w.WriteHeader(500)
//	} else {
//		response, _ := json.Marshal(result)
//		w.Write(response)
//	}
//}

func handleReverse(w http.ResponseWriter, r *http.Request) {
	log.Println("handleReverse: Received 'reverse' request over http")
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
	log.Println("handleUpper: Received 'upper' request over http")

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
