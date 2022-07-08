package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

func main() {
	var text = "Hello echo"
	var appID = pkg.EchoServiceAppID
	var cmd string
	var repeat int = 1
	flag.IntVar(&repeat, "repeat", repeat, "Nr of times to invoke")
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
	InvokeHttpServiceWithSDK(appID, cmd, text, repeat)
}

// InvokeHttpServiceWithSDK invokes a service using the dapr sdk. See also:
//  https://docs.dapr.io/developing-applications/building-blocks/service-invocation/howto-invoke-discover-services/
func InvokeHttpServiceWithSDK(appID string, cmd string, text string, repeat int) {
	fmt.Println(fmt.Sprintf("Invoking command '%s' on service '%s'", cmd, appID))
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
	t1 := time.Now()
	for count := 0; count < repeat; count++ {
		resp, err := client.InvokeMethodWithContent(ctx, appID, cmd, "post", content)
		if err != nil {
			msg := fmt.Sprintf("Error invoking method '%s' on app '%s': %s", cmd, appID, err)
			log.Println(msg)
		}
		_ = resp
		//fmt.Println("Response:", string(resp))
	}
	t2 := time.Now()
	duration := t2.Sub(t1)
	fmt.Println("Time to invoke: ", duration.Milliseconds(), "msec")
}
