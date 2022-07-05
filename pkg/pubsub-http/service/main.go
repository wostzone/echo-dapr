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

// Entry point of the http echo service
func main() {
	var port int
	flag.IntVar(&port, "port", pkg.EchoServiceHttpPort, "Service http listening port")
	flag.Parse()

	StartPubSubHttpService(port)
}

// StartPubSubHttpService subscribes and publishes echo events over http using dapr pub/sub feature
func StartPubSubHttpService(port int) {
	fmt.Println("Starting echo-service on http port ", port)

	service := daprd.NewService(":" + strconv.Itoa(port))
	//service := daprd.NewService(":" + strconv.Itoa(port))

	if err := service.AddServiceInvocationHandler("echo", echoServiceHandler); err != nil {
		log.Fatalf("error adding echo invocation handler: %v", err)
	}

	//// Handle binding invocations
	//if err := service.AddBindingInvocationHandler("echo1", echoBindingHandler); err != nil {
	//	log.Fatalf("error adding binding handler: %v", err)
	//}
	//if err := service.AddBindingInvocationHandler("upper1", upperBindingHandler); err != nil {
	//	log.Fatalf("error adding binding handler: %v", err)
	//}
	//if err := service.AddBindingInvocationHandler("reverse1", reverseBindingHandler); err != nil {
	//	log.Fatalf("error adding binding handler: %v", err)
	//}
	// subscribe to events
	echoSub := &common.Subscription{PubsubName: "pubsub", Topic: "echo", Route: "/echo"}
	if err := service.AddTopicEventHandler(echoSub, echoEventHandler); err != nil {
		log.Fatalf("error adding topic event handler: %v", err)
	}
	upperSub := &common.Subscription{PubsubName: "pubsub", Topic: "upper", Route: "/upper"}
	if err := service.AddTopicEventHandler(upperSub, upperEventHandler); err != nil {
		log.Fatalf("error adding topic event handler: %v", err)
	}
	reverseSub := &common.Subscription{PubsubName: "pubsub", Topic: "reverse", Route: "/reverse"}
	if err := service.AddTopicEventHandler(reverseSub, reverseEventHandler); err != nil {
		log.Fatalf("error adding topic event handler: %v", err)
	}

	if err := service.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
	}
}

func echoServiceHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	var args *pb.TextParam
	var response []byte
	log.Printf("echoServiceHandler - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	err = json.Unmarshal(in.Data, &args)
	echoService := internal.NewEchoService(nil)
	result, _ := echoService.Echo(nil, args)
	response, _ = json.Marshal(result)

	// do something with the invocation here
	out = &common.Content{
		Data:        response,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}

func echoEventHandler(ctx context.Context, ev *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - topic:%s, ID:%s, Data:%v", ev.Topic, ev.ID, ev.Data)
	// do something with the event here
	return false, nil
}

func upperEventHandler(ctx context.Context, ev *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - topic:%s, ID:%s, Data:%v", ev.Topic, ev.ID, ev.Data)
	// do something with the event here
	return false, nil
}
func reverseEventHandler(ctx context.Context, ev *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - topic:%s, ID:%s, Data:%v", ev.Topic, ev.ID, ev.Data)
	// do something with the event here
	return false, nil
}

func echoBindingHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Printf("binding - Data:%v, Meta:%v", string(in.Data), in.Metadata)
	// do something with the invocation here
	return in.Data, nil
}

func reverseBindingHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Println("reverseBindingHandler: Received 'reverse' request over http")
	var args *pb.TextParam

	json.Unmarshal(in.Data, &args)
	echoService := internal.NewEchoService(nil)
	result, err := echoService.Reverse(nil, args)
	if err != nil {
		err := fmt.Errorf("error handling reverse request:", err)
		fmt.Println(err)
		return nil, err
	} else {
		response, _ := json.Marshal(result)
		return response, nil
	}
}

func upperBindingHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Println("upperBindingHandler: Received 'upper' request over http")
	var args *pb.TextParam

	json.Unmarshal(in.Data, &args)
	echoService := internal.NewEchoService(nil)
	result, err := echoService.UpperCase(nil, args)
	if err != nil {
		err := fmt.Errorf("error handling upper request:", err)
		fmt.Println(err)
		return nil, err
	} else {
		response, _ := json.Marshal(result)
		return response, nil
	}
}
