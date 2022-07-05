package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

var subEcho = &common.Subscription{
	PubsubName: "pubsub",
	Topic:      "echo",
	Route:      "/echo",
}
var subUpper = &common.Subscription{
	PubsubName: "pubsub",
	Topic:      "upper",
	Route:      "/upper",
}
var subReverse = &common.Subscription{
	PubsubName: "pubsub",
	Topic:      "reverse",
	Route:      "/reverse",
}

// StartPubSub on the given port
// See also https://docs.dapr.io/developing-applications/sdks/go/go-service/grpc-service/
//func StartPubSub(listener net.Listener) error {
func StartPubSub(port int) error {
	//fmt.Println("StartPubSub. Listening...", listener.Addr().String())
	//s := daprd.NewServiceWithListener(listener)
	s, _ := daprd.NewService(":" + strconv.Itoa(port))

	//if err := s.AddServiceInvocationHandler("echo", serviceInvocationHandler); err != nil {
	//	log.Fatalf("error adding invocation handler: %v", err)
	//}

	if err := s.AddTopicEventHandler(subEcho, topicEventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}
	if err := s.AddTopicEventHandler(subUpper, topicEventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}
	if err := s.AddTopicEventHandler(subReverse, topicEventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting listening service: %v", err)
	}

	return nil
}

func topicEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Println("topicEventHandler. input: ", e.Data)
	return false, nil
}

func serviceInvocationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	fmt.Println("serviceInvocationHandler. input: %v", in.Data)

	out = &common.Content{
		Data:        in.Data,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
