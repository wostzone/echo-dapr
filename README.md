# echo

Echo demonstration service using gRPC and dapr

This service is intended as an example on how to write a microservice using gRPC in
a [dapr](https://docs.dapr.io/) environment. The service is written in golang and can be invoked via
dapr gRPC and http endpoints.

This project includes 3 demonstrations:

1. How to write and invoke a gRPC service in golang using a protobuffers for API definition
2. How to invoke the service using dapr
3. How to extend the service with the dapr SDK to allow invocation via the dapr HTTP endpoint
4. How to extend the service with the dapr SDK to allow invocation via the dapr pub/sub

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

Self hosted installation of dapr.

## Writing a gRPC service

## Invoking a gRPC service

## Extend the service to support dapr HTTP endpoint

## Extend the service to respond to pub/sub requests
