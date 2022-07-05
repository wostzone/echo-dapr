package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	dapr "github.com/dapr/go-sdk/client"

	pb "github.com/wostzone/echo/proto/go"
)

func main() {
	var text = "Hello echo"
	var pubsub = "pubsub"
	var cmd string
	flag.StringVar(&pubsub, "pubsub", pubsub, "Pubsub component to use")
	flag.Parse()
	values := flag.Args()
	if len(values) == 1 && values[0] == "stop" {
		cmd = values[0]
	} else if len(values) == 2 {
		cmd = values[0]
		text = values[1]
	} else {
		fmt.Println("Missing text: pubsub-http-client <command> <text>: ", cmd)
		flag.Usage()
		os.Exit(1)
	}

	PublishEchoOverHttp(pubsub, cmd, text)
}

// PublishEchoOverHttp publishes events and bindings using the dapr sdk.
func PublishEchoOverHttp(pubsub string, cmd string, text string) {
	fmt.Println(fmt.Sprintf("Publishing event '%s' over http on publisher '%s'", cmd, pubsub))
	message := pb.TextParam{Text: text}
	data, _ := json.Marshal(message)

	// This creates a dapr runtime able to connect to sidecars and access the state stores
	// FYI, if you get context deadline exceeded error then the sidecar isnt running
	//client, err := dapr.NewClientWithAddress("localhost:" + strconv.Itoa(clientPort))
	//client, err := dapr.NewClientWithPort(strconv.Itoa(clientPort))
	client, err := dapr.NewClient()
	if err != nil {
		err2 := fmt.Errorf("error initializing client. Make sure this runs with a sidecart.: %s", err)
		log.Println(err2)
		os.Exit(1)
	}
	defer client.Close()
	ctx := context.Background()
	err = client.PublishEvent(ctx, pubsub, cmd, data)
	if err != nil {
		msg := fmt.Sprintf("Error publishing event '%s' on app '%s': %s", cmd, pubsub, err)
		log.Println(msg)
		os.Exit(3)
	}
	fmt.Println(fmt.Sprintf("Event '%s' published succesfully", cmd))
}
