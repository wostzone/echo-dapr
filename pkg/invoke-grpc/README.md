# Invoke with gRPC using dapr SDK

This example shows how to invoke a gRPC service that is developed using the dapr SDK. This is based on the [dapr SDK gRPC service docs for Go](https://docs.dapr.io/developing-applications/sdks/go/go-service/grpc-service/).

Use-case: Newly developed services that communicate via gRPC while using dapr for features such as logging, tracing, resiliency and other middleware.

Pros:

* Client is protocol agnostic using the dapr SDK. Just create the client instance and invoke the method on a named service. The sidecar handles the port mapping and routing behind the scenes.
* The service can also be invoked over plain http with and without client sidecar. Use a client sidecar lets you use dapr middleware and auto resolves the server.
* Access to dapr's state store, bindings, pubsub services.
* Fast communication between services if all goes over gRPC. (not really)

Cons:

* This is incompatible with your service gRPC protobuf API. It cannot be invoked with a gRPC client built against your proto file. If you just need to communicate with an existing gRPC service then simply launch a sidecar for both client and server. However this is not compatible with using http, dapr's invoke, pubsub or bindings.
* The http client-server invocation is faster than the grpc client-server invocation. 700msec for 1000 calls vs 890msec. Why?

# gRPC Invocation Using The SDK

Both client and service use the SDK to implement dapr integration using communication over gRPC.

> flow: client/sdk -[grpc]-> client-sidecar -[grpc] -> service-sidecar -[grpc]-> service

## Service

The documentation shows to create a service you run daprd.NewService(port) then attach any number of event, binding and invocation handlers.

Service code snippet ([full example code](service/main.go)):

```go
package main

import daprd "github.com/dapr/go-sdk/service/grpc"

func main() {
	address := ":40001"
	appcallbackServer, err := daprd.NewService(address)
	if err != nil {
		log.Fatalf("failed to create app callback server: %v", err)
	}
	if err := appcallbackServer.AddServiceInvocationHandler("echo", echoHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}
}

func echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	//...
}
```

How to run:

```bash
dapr run --app-port 40001 --app-id echo --app-protocol grpc -- go run service/main.go -port 40001
```

On a local network there is no need to assign a dapr port, as a dapr client sidecar will locate the server sidecar using dns. Dapr will auto-assign ports.

## Client

Clients that use the dapr SDK InvokeMethod only need the name of the service. dapr manages the ports and resolves the service sidecar. The client must be invoked through the sidecart using 'dapr run'. Clients that use the SDK are therefore easier use than clients that don't as these need to be configured to talk to the sidecar port.

Client code snippet ([full example code](client/main.go)):

```go
package main

import dapr "github.com/dapr/go-sdk/client"

func main() {
	client, err := dapr.NewClient()
	ctx := context.Background()
	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        data,
	}
	resp, err := client.InvokeMethodWithContent(ctx, "echo", "upper", "post", content)
}
```

How to run:

```bash
dapr run --app-id echoclient -- go run client/main.go upper "Hello world"
```

Note 1: No need to specify an app protocol. dapr resolves it automagically
Note 2: If the client subscribes to events or bindings, an app-port must be provided to dapr on which the client is listening.
