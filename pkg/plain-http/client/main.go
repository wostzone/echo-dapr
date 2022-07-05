package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wostzone/echo/pkg"
	pb "github.com/wostzone/echo/proto/go"
)

// Client used to invoke the echo service over http using dapr
func main() {
	var text = "Hello echo"
	var appID = pkg.EchoServiceAppID
	var cmd string
	var port int
	var repeat int = 1
	flag.IntVar(&port, "port", pkg.EchoServiceHttpPort, "Service http listening port")
	flag.StringVar(&appID, "app-id", pkg.EchoServiceAppID, "Service name when using dapr")
	flag.IntVar(&repeat, "repeat", repeat, "Nr of times to invoke")
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

	InvokeHttpService(port, appID, cmd, text, repeat)
}

// InvokeHttpService invokes the service with http using a dapr sidecar
func InvokeHttpService(port int, appID string, cmd string, text string, repeat int) {
	fmt.Println("Invoking echo service over http on :"+strconv.Itoa(port), "command: ", cmd)
	message := pb.TextParam{Text: text}
	data, _ := json.Marshal(message)

	client := &http.Client{}
	url := "http://localhost:" + strconv.Itoa(port) + "/" + cmd
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Invoking the service
	// Set the service name to connect to when connecting via dapr
	req.Header.Add("dapr-app-id", appID)

	t1 := time.Now()
	for count := 0; count < repeat; count++ {

		response, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		result, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Result: ", string(result))
	}
	t2 := time.Now()
	duration := t2.Sub(t1)
	fmt.Println(fmt.Sprintf("Time to invoke %d http calls: %d msec", repeat, duration.Milliseconds()))
}
