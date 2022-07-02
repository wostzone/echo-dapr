# echo-dapr

Echo demonstration service in golang using dapr with gRPC and HTTP invocation.

This service is intended as an exploration in how to write and use microservice in a [dapr](https://docs.dapr.io/) environment, using gRPC and HTTP.

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
| internal/ - source code of the echo service
| pkg/ - source code of client and service for various invocation methods  
| proto/ - gRPC api definition in protobuffer and generated go API source
| go.mod - go dependencies
| Makefile - build and run the demo

## Prerequisites

* Self hosted installation of dapr
* Golang 1.18 or newer

## Demos

See each of the demo README's for details

Completed:

* invoke-grpc/README.md) - invoke service with grpc with/without dapr
* invoke-http/README.md - invoke service with http with/without dapr

In progress:

* invoke-dapr-grpc/README.md - invoke service with grpc using dapr golang SDK
* invoke-http-grpc/README.md - invoke service with http using dapr golang SDK

Planned:

* pubsub-grpc/README.md - invoke service with pubsub over grpc
* pubsub-http/README.md - invoke service with pubsub over http
* invoke-grpc-http/README.md - invoke the grpc service via http
