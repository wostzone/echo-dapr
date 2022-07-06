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

* [plain-grpc/README.md](pkg/plain-grpc/README.md) - client-service echo with grpc
* [plain-http/README.md](pkg/plain-http/README.md) - client-service echo with http
* [invoke-grpc/README.md](pkg/invoke-grpc/README.md) - invoke service with grpc using dapr SDK
* [invoke-http/README.md](pkg/invoke-http/README.md) - invoke service with http using dapr SDK

Planned:

* pubsub-grpc/README.md - invoke service with pubsub over grpc
* pubsub-http/README.md - invoke service with pubsub over http

## Performance

This is a very simplistic performance test comparing invocation times using best out of 3 runs.

| Makefile | test                         | duration  |
|----------|------------------------------|-----------|
| run1     | 1000 plain grpc calls        | 164 msec  |
| run2     | 1000 plain http calls        | 320 msec  |
| run3     | 1000 dapr wrapped grpc calls | 986 msec  |
| run4     | 1000 dapr wrapped http calls | 824 msec  |
| run5     | 1000 grpc sdk calls          | 772 msec  |
| run6     | 1000 http sdk calls          | 700 msec  |

'wrapped calls' invoke the plain non-dapr service via dapr sidecars.
Findings:

* Using dapr is 4.7 times slower than direct calls for grpc and 2.1 times slower for http calls using the sdk.
* Using the dapr SDK improves performance over the wrapped calls by roughly 20% for grpc and 15% for http.
* Grpc calls via dapr sidecars are slower than http calls by 10% when using the SDK and 15% when using wrapped calls.
