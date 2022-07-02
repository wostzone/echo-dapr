# Invoke HTTP

This example shows how to invoke a HTTP service that is developed without dependency on dapr, but can interoperate with dapr using http.

Use-case: Integrate dapr with existing HTTP services to use dapr's features.

Pros:

* Services can be developed, tested, and used with or without dapr.
* Dapr features are available when invoking the service through the dapr sidecar proxy:
    * Client can invoke the service using the application name. No need to know its port.
    * Mutual authentication and api token based auth
    * Tracing, logging, resiliency features.

Cons:

* Services can not be invoked via gRPC. Dapr does not convert gRPC to http/json. Developers must implement a separate gateway service that maps gRPC->http. This is not really a con when the HTTP API is for external clients only or looks different from the micro-service API.
* Multiplies port usage: 1 port for the dapr client sidecar, 1 port for each service instance, and 1 port for each dapr sidecar. Eg a dapr client sidecar with 3 service instances requires 1+3*2 = 7 unique ports instead of 3 ports without dapr.

# Invocation Methods

(this is based on the current understanding of using dapr)

## 1. Direct Client To Service

> client -[http]-> service

This usage is simply a client->service invocation over HTTP. No dapr involved.

When to use:

* when using a single service instance, and
* when the client knows the service port, and
* when there are only a few services to manage, and
* when you don't need dapr features such as logging, tracing, ...

How to run:

```bash
go run service/main.go -port 40002 &
go run client/main.go -port 40002 upper "Hello world"
go run client/main.go -port 40002 stop
```

## 2. Indirect Client To Service Sidecar

> client -[http]-> service-sidecar -[http]-> service

This usage does not use a client sidecar but simply invokes the service sidecar, which in turn invokes the service.

Note: this requires the 'dapr-app-id' metadata field set in the client to the name of the service to invoke:
> req.Header.Add("dapr-app-id", "order-processor")


When to use:

* when using a single service instance, and
* when you want to use dapr features such as logging, tracing, ...

How to run (ports are arbitrary):

```sh
dapr run --app-protocol http \
	--app-port 40002 \
	--app-id echo \
	--dapr-http-port 9002 \
	 --  go run service/main.go -port 40002 &
go run client/main.go --port 9002 --app-id echo upper "Hello world"
go run client/main.go --port 9002 --app-id echo "stop"
```

## 3. Indirect Using Client Sidecar

> client -[http]-> client-sidecar -[grpc]-> service sidecar -[http]-> service

Note: this requires the 'dapr-app-id' metadata field set in the client to the name of the service to invoke:
> req.Header.Add("dapr-app-id", "order-processor")

When to use:

* when the client doesn't know the server port, only the dapr client sidecar, or
* when using multiple service instances with round-robin load balancing, and
* when you want to use dapr features such as logging, tracing, ...

How to run (ports are arbitrary):

```sh
dapr run --app-protocol http \
	--app-port 40002 \
	--app-id echo \
	--dapr-http-port 9002 \
	 --  go run service/main.go -port 40002 &
dapr run --dapr-http-port 9003 \
    -- go run client/main.go -port 9003 --app-id echo upper "Hello"
dapr run --dapr-http-port 9003 \ 
    -- go run client/main.go -port 9003 --app-id echo stop
```
