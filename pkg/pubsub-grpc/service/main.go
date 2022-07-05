package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"

	"github.com/wostzone/echo/pkg"
)

// Entry point of the http echo service
func main() {
	var port int
	flag.IntVar(&port, "port", pkg.EchoServiceHttpPort, "Service http listening port")
	flag.Parse()

	StartPubSubGrpcService(port)
}

// StartPubSubGrpcService subscribes and publishes echo events over grpc using dapr pub/sub feature
func StartPubSubGrpcService(port int) {
	fmt.Println("Starting echo-service on grpc port ", port)

	service, err := daprd.NewService(":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// subscribe to events
	echoSub := &common.Subscription{PubsubName: "pubsub", Topic: "echo", Route: "/echo"}
	if err := service.AddTopicEventHandler(echoSub, echoEventHandler); err != nil {
		log.Fatalf("error adding topic event handler: %v", err)
	}

	if err := service.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
	}
}

func echoEventHandler(ctx context.Context, ev *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - topic:%s, ID:%s, Data:%v", ev.Topic, ev.ID, ev.Data)
	// do something with the event here
	return false, nil
}
