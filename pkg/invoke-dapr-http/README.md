# Invoke with HTTP using dapr SDK

This example shows how to invoke a HTTP service that is developed using the dapr SDK. This is based on the [dapr SDK docs for Go](https://docs.dapr.io/developing-applications/sdks/go/go-client/).

Use-case: New service with direct dapr integration

Pros:

* Client does not need to know the addresses and ports of service or sidecar. Instead, the service app-id is used.
* Use of dapr middleware services

Cons:

* Overhead of using HTTP. Prefer the use of gRPC when possible.

# Invocation Methods

(this is based on the current understanding of using dapr)

## 1. Invoke A Service Via Its Sidecar Using The SDK

Both client and service can use the SDK to implement dapr integration.

> flow: client/sdk -[grpc?]-> service-sidecar -[http]-> service

### Client Integration

Clients that use the dapr SDK InvokeMethod only need the name of the service. dapr manages the ports and resolves the service sidecar. The client must be invoked through the sidecart using 'dapr run'.

When to use:

* When communicating with an existing http service, a new client can use the sdk to easily connect over http.

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

When to use:

* There is no clear advantage in using the SDK for service invocation. service.AddServiceInvocationHandler seems to be similar to mux.HandleFunc.
* When adding event handling to the service
* When adding binding invocation handling to the service

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
