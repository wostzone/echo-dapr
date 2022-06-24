# echo

Echo demonstration service using gRPC and dapr

This service is intended as an example on how to write a microservice using gRPC in
a [dapr](https://docs.dapr.io/) environment. The service is written in golang and can be invoked via
dapr gRPC and http endpoints.

This project includes the following demonstrations:

1. How to write and invoke a gRPC service in golang using a protobuffers for API definition
2. How to invoke the service using dapr
3. How to extend the service with the dapr SDK to allow invocation via the dapr HTTP endpoint
4. How to extend the service with the dapr SDK to allow invocation via the dapr pub/sub (todo)

Wishlist:

1. Show logging using dapr middleware
2. Show authentication using dapr middleware

Project structure :
| bin/ - service and client binaries from make all
| cmd/ - source code of commands for the service and the client  
| internal/ - source code of the echo service
| proto/ - gRPC api definition in protobuffer and generated go API source
| go.mod - go dependencies
| Makefile - build and run the demo

## Prerequisites

* Self hosted installation of dapr
* Golang 1.18 or newer

## Demos

### Run the demos

Build and run the echo service using gRPC: echo-client/grpc -> echo-service/grpc
> make run1

Build and run the echo service with dapr and invoke using gRPC:
echo-client/grpc -> dapr/grpc (sidecart) -> echo-service/grpc
> make run2

Build and run the echo service using dapr and invoke via curl:
curl/http -> dapr/http -> dapr/grpc (sidecart) -> echo-service/grpc)
> make run3

### Running the echo service via dapr

Launch the service using dapr. Dapr listens on GRPC port 9000
> dapr run --enable-api-logging
> --app-protocol grpc --app-port 40001 --app-id echo
> --dapr-http-port 9000 --dapr-grpc-port 9001 -- bin/echo-service &

to stop:
> dapr stop --app-id echo

### Invoke the service via gRPC

> bin/echo-cli -port 9001 upper "Hello gRPC via dapr gRPC"

### Invoke the service via curl

> curl localhost:9000/v1.0/invoke/echo/method/upper -d "Hello world" -X PUT
