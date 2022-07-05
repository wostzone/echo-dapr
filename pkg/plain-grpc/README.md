# Invoke with gRPC without SDK

This example shows how to invoke a gRPC service that is developed without dependency on dapr, but can interoperate with dapr using gRPC.

Use-case: Integrate dapr with existing gRPC services to use dapr's features.

Pros:

* Services can be developed, tested, and used with or without dapr.
* Dapr features are available when invoking the service through the dapr sidecar proxy:
    * Client can invoke the service using the application name. No need to know its port.
    * Mutual authentication and api token based auth
    * Tracing, logging, resiliency features.

Cons:

* Services can not be invoked via http. Dapr does not convert http/json to gRPC. Developers must implement a separate http gateway service that maps http->gRPC. This is not really a con when the HTTP API for external clients looks different from the micro-service gRPC API.
* Multiplies port usage: 1 port for the dapr client sidecar, 1 port for each service instance, and 1 port for each dapr sidecar. Eg a dapr client sidecar with 3 service instances requires 1+3*2 = 7 unique ports instead of 3 ports without dapr.

Open questions still to answer:

1. Can dapr round robin load balancing be used?

# Invocation Methods

(this is based on the current understanding of using dapr)

Before running the examples, compile the proto files and generate the go API messages and methods. From the project root folder run 'make proto'.

This does something like:

```sh
protoc --proto_path=./proto\
  --go_opt=paths=source_relative \
  --go_out=./proto/go \
  --go-grpc_out=./proto/go \
  --go-grpc_opt=paths=source_relative \
  proto/echo.proto
go mod tidy
```

## 1. Direct Client To Service

> client -[grpc]-> service

This usage is simply a client->service invocation over gRPC. No dapr involved.

When to use:

* when using a single service instance, and
* when the client knows the service port, and
* when there are only a few services to manage, and
* when you don't need dapr features such as logging, tracing, ...

How to run:

```bash
go run service/main.go -port 40001 &
go run client/main.go -port 40001 upper "Hello world"
go run client/main.go -port 40001 stop
```

## 2. Indirect Client To Service Sidecar

> client -[grpc]-> service-sidecar -[grpc]-> service

This usage does not use a client sidecar but simply invokes the service sidecar.

Note: this requires the 'dapr-app-id' metadata field set in the client to the name of the service to invoke:
> ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", pkg.EchoServiceAppID)


When to use:

* when using a single service instance, and
* when you want to use dapr features such as logging, tracing, ...

How to run (ports are arbitrary):

```sh
dapr run --app-protocol grpc \
	--app-port 40001 \
	--app-id echo \
	--dapr-grpc-port 9001 \
	 --  go run service/main.go -port 40001 &
go run client/main.go --port 9001 --app-id echo upper "Hello world"
go run client/main.go --port 9001 --app-id echo "stop"
```

## 3. Indirect Using Client Sidecar

> client -[grpc]-> client-sidecar -[grpc]-> service sidecar -[grpc]-> service

Note: this requires the 'dapr-app-id' metadata field set in the client to the name of the service to invoke:
> ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", pkg.EchoServiceAppID)

When to use:

* when the client doesn't know the server port, only the dapr client sidecar, or
* when using multiple service instances with round-robin load balancing, and
* when you want to use dapr features such as logging, tracing, ...

How to run (ports are arbitrary):

```sh
dapr run --app-protocol grpc \
	--app-port 40001 \
	--app-id echo \
	--dapr-grpc-port 9001 \
	 --  go run service/main.go -port 40001 &
dapr run --dapr-grpc-port 9101 \
    -- go run client/main.go -port 9101 --app-id echo upper "Hello"
dapr run --dapr-grpc-port 9101 \ 
    -- go run client/main.go -port 9101 --app-id echo stop
```
