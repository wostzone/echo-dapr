package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	r.HandleFunc("/echo", handleEcho).Methods("POST")
	r.HandleFunc("/upper", handleUpper).Methods("POST")
	r.HandleFunc("/reverse", handleReverse).Methods("POST")
	_ = http.ListenAndServe(":"+strconv.Itoa(port), r)
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	var args *pb.TextParam
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body handling echo request: ", err)
		w.WriteHeader(500)
		return
	}
	err = json.Unmarshal(data, &args)
	if err != nil {
		log.Println("Error unmarshalling payload for handleEcho: ", err)
		w.WriteHeader(500)
		return
	}
	log.Println("handleEcho: Received 'echo' request over http")
	echoService := internal.NewEchoService(nil)
	result, err := echoService.Echo(nil, args)
	if err != nil {
		log.Println("Error handling echo request: ", err)
		w.WriteHeader(500)
	} else {
		response, _ := json.Marshal(result)
		w.Write(response)
	}
}

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
