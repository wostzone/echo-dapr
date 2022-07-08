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

This is a very simplistic performance test comparing 'echo' invocation times using best out of 5 runs.

There are 3 types of tests, each with grpc and http:

* plain calls are simple calls from a client directly to a service on localhost. Eg: client -> service
* dapr wrapped calls use dapr sidecars with plain client and server. Eg: client -> sidecar -> sidecar -> service
* in sdk calls the client and service use the dapr sdk to invoke and receive messages. This integration adds a dependency on the dapr sdk in return for perfomance gains and simplified port management.

Durations with small payload: {"Text": "Hello world"} with printing of the result disabled. (adds up to 40msec)
The duration only includes the time to make the call and does not include the time to startup dapr and the test.

| Makefile | test                         | Hello world | 1K text      | 10K text     | 100K text      |
|----------|------------------------------|-------------|--------------|--------------|----------------|
| run1     | 1000 plain grpc calls        | 150 msec    | 150 msec     | 200 msec     | 550 msec       |
| run2     | 1000 plain http calls        | 320 msec    | 355 msec     | 520 msec     | 2050 msec      |
| run3     | 1000 dapr wrapped grpc calls | 970 msec    | 980 msec     | 1110 msec    | fails (*1)     |
| run4     | 1000 dapr wrapped http calls | 810 msec    | 820 msec     | 1110 msec    | 3160 msec      |
| run5     | 1000 grpc sdk calls          | 760 msec    | 760 msec     | 980 msec     | 3230 msec      |
| run6     | 1000 http sdk calls          | 690 msec    | 700/685 msec | 960/800 msec | 3120/1830 msec |

*1  (93K message size succeeds with duration 1980 msec for 1000 calls, anything over 93K sec fails with context deadline exceeded. Looks like a timeout at the 2msec mark)

Findings:

* Durations vary by up to 10%
* For small messages < 1K
  * using dapr is 4.5 times slower than direct calls for grpc and 2.1 times slower for http calls using the sdk.
  * Using the dapr SDK improves performance over the wrapped calls by roughly 20% for grpc and 15% for http.
  * Grpc calls via dapr sidecars are slower than http calls by 10% when using the SDK and 15% when using wrapped calls.
* For larger messages 100K the difference between dapr based calls
  * using dapr is still 4.5 times slower than direct calls for gRPC.
  * the performance gap between plain gRPC and http messages widens. Most likely due to compression in gRPC.
  * dapr performance differences disappear. Each message takes roughly 3 msec.
* run 2, run4, run5, run6 all use json unmarshal to decode the payload and json marshal to encode it. At larger payloads this really starts to slow things down. The second number in run6 is without marshalling the payload.
