# Invoke with HTTP using dapr SDK

This example shows how to invoke a HTTP service that is developed using the dapr SDK. This is based on the [dapr SDK docs for Go](https://docs.dapr.io/developing-applications/sdks/go/go-client/).

Use-case: Newly developed clients that invoke http services that run with a dapr sidecar. The sidecar supports dapr features such as logging, tracing, resiliency and other middleware. The use of pub/sub and bindings is described in the pubsub-http package.

Pros:

* Client does not need to know the addresses and ports of service or sidecar. Instead, the service app-id is used.
* Use of dapr middleware services

Cons:

* Overhead of using HTTP. Prefer the use of gRPC when possible.

# Service Invocation Using The SDK

Both client and service can use the SDK to implement dapr integration.

> flow: client/sdk -[grpc?]-> service-sidecar -[http]-> service

### Client Integration

Clients that use the dapr SDK InvokeMethod only need the name of the service. dapr manages the ports and resolves the service sidecar. The client must be invoked through the sidecart using 'dapr run'. Clients that use the SDK are therefore easier use than clients that don't as these need to be configured to talk to the sidecar port.

When to use:

* When communicating with a http service that has a dapr sidecar, a client can use the sdk to easily invoke a method over http using the app-id and method name.

Client code snippet ([full example code](client/main.go)):

```go
import dapr "github.com/dapr/go-sdk/client"
client, err := dapr.NewClient()
resp, err := client.InvokeMethod(ctx, "echo", "upper", "post")
```

How to run:

```bash
dapr run -- go run client/main.go upper "Hello world"
```

### Service Integration

There is no clear advantage in using the SDK for just service invocation handling. The SDK's 'service.AddServiceInvocationHandler' method does not have benefits over using a simple http.Server instance with a mux.HandleFunc will do the same.

However, when services also use pub/sub event invocation and binding invocation, the same dapr service instance can be used to add other the invocation handlers. Route names must be unique amongst all invocation methods, only the first name registered is used.

Service code snippet ([full example code](service/main.go)):

```go
import daprd "github.com/dapr/go-sdk/service/http"

service := daprd.NewServiceWithMux(":"+strconv.Itoa(port), r)
if err := service.AddServiceInvocationHandler("/echo", echoHandler); err != nil {
log.Fatalf("error adding invocation handler: %v", err)
}

func echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
//...
}
```

How to run:

```bash
dapr run --app-port 40002 --app-id echo --app-protocol http -- go run service/main.go -port 40002
```

On a local network there is no need to assign a dapr http port as a dapr client sidecar will locate the server sidecar using dns. Dapr will auto-assign a port itself.

## 2. Using curl As Client

Since this is a http service, curl can be used as the client to invoke the service using its HTTP port. As no client sidecar is used, the service sidecar must have a HTTP port assigned. dapr uses a url format that includes the version, app-id, and method:

```
curl localhost:{port}/v1.0/invoke/{app-id}/method/{method-name} -d "{"text":"Hello world"}" -X POST
```

Where

* {port} is the service sidecar port. For example "--dapr-http-port 9002"
* {app-id} is the service application name, eg echo
* {method-name} is the registered invocation handler, eg "echo", "upper", "reverse"
