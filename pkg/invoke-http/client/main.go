package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

func main() {
	var text = "Hello echo"
	var appID = pkg.EchoServiceAppID
	var cmd string
	var port int
	flag.IntVar(&port, "port", pkg.EchoDaprClientHttpPort, "client sidecar http listening port")
	flag.StringVar(&appID, "app-id", pkg.EchoServiceAppID, "Service app-id to invoke")
	flag.Parse()
	values := flag.Args()
	if len(values) == 1 && values[0] == "stop" {
		cmd = values[0]
	} else if len(values) == 2 {
		cmd = values[0]
		text = values[1]
	} else {
		fmt.Println("Missing text: invoke-http-client <command> <text>: ", cmd)
		flag.Usage()
		return
	}

	InvokeHttpServiceWithSDK(port, appID, cmd, text)
}

// InvokeHttpServiceWithSDK invokes a service using the dapr sdk. See also:
//  https://docs.dapr.io/developing-applications/building-blocks/service-invocation/howto-invoke-discover-services/
func InvokeHttpServiceWithSDK(clientPort int, appID string, cmd string, text string) {
	fmt.Println("Invoking echo service over http on :"+strconv.Itoa(clientPort), "command: ", cmd)
	message := pb.TextParam{Text: text}
	data, _ := json.Marshal(message)

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        data,
	}
	// This creates a dapr runtime able to connect to sidecars and access the state stores
	// FYI, if you get context deadline exceeded error then the sidecar isnt running
	//client, err := dapr.NewClientWithAddress("localhost:" + strconv.Itoa(clientPort))
	client, err := dapr.NewClient()
	if err != nil {
		err2 := fmt.Errorf("error initializing client. Make sure this runs with a sidecart.: %s", err)
		log.Println(err2)
		return
	}
	defer client.Close()
	ctx := context.Background()
	resp, err := client.InvokeMethodWithContent(ctx, appID, cmd, "post", content)
	if err != nil {
		msg := fmt.Sprintf("Error invoking method '%s' on app '%s': %s", cmd, appID, err)
		log.Println(msg)
	}
	fmt.Println("Response:", string(resp))
}
