# Pub/Sub with HTTP

This example shows how to use dapr's publish-subscribe in an HTTP service . This is based on the [dapr How to: Publish & subscribe to topics](https://docs.dapr.io/developing-applications/building-blocks/pubsub/howto-publish-subscribe/).

Several issues/questions came up during experimentation. July 2022 the documentation above is confusing:

1. Examples use ports (--dapr-http-port, --dapr-grpc-port) that aren't needed nor explained.
2. The subscription.yaml file doesn't seem to do anything or be needed to run the examples.
3. The documentation example uses both the subscription.yaml file and a AddTopicEventHandler, which contains the same topic, route, and pubsub name. Isn't this duplicate or must the handler be registered with AddTopicEventHandler? (as opposed to running a regular non-dapr http service on the route address)?
4. The documentation refers [more info on pub/sub](https://docs.dapr.io/developing-applications/building-blocks/pubsub/subscription-methods/) which describes a function 'configureSubscribeHandler' on address /dapr/subscribe for programmatic subscriptions. Everything works fine without it however.
5. subscription.yaml includes a scopes section. The documentation claims that it is used to restrict access to topics to the listed application IDs. This doesn't seem to have any effect though.
6. subscription.yaml includes 'apiVersion v1alpha1' which suggests this is an unstable feature. How stable is this? Is there a stable version?
7. On startup historical events are passed to the subscriber. However, these can arrive out of order which turns this feature into a problem for order dependent events. The documentation doesn't explain event order.

Use-case: clients publish messages without knowledge of the consumers. Anyone with an interest, and authorization, can receive the information.
Use-case: dapr clients need to publish messages to a service without being tied to a particular pub/sub messaging system.
Use-case: Add a pub/sub feature to a non-dapr service using its http endpoint. (using declarative subscriptions)
Use-case: Messages have a limited time to live (TTL). Subscribers do not receive those messages after expiration.
Use-case: Subscribers receive queued messages published before connecting - is this possible?

# Prerequisites

The dapr pubsub component must be active. This example uses dapr in standalone mode with the default ~/.dapr/components/pubsub.yaml (Linux) as created during dapr init. This uses pubsub.redis.

Pub/sub topics can be declared in a subscription yaml file, or provided programmatically by the service via a http /dapr/subscribe handler. Both methods are described below.

# 1. Declarative Method

The declarative method defines the topics in the pub/sub component configuration file for the service sidecar.

> flow: client/event -> client-sidecar -[grpc]-> server-sidecar -[http]-> service event handler

Clients that use the dapr SDK InvokeEvent only need the name of the service. dapr manages the ports, resolves the service sidecar and determines the pub/sub connection. The client must be run through the sidecart using 'dapr run'.

When to use:

* When pub/sub uses static routes
* Legacy services that only have HTTP endpoints, declarative routes can be used to route published messages to the endpoint.

## Service Setup

Configure dapr's components/subscription.yaml:

```yaml
apiVersion: dapr.io/v1alpha1
kind: Subscription
metadata:
  name: echo          # the name of the subscription to use
spec:
  topic: echo         # the publication topic
  route: /echo        # the http endpoint of the service to send the publications to
  pubsubname: pubsub  # the name of the pubsub component to use as defined in components.yaml
scopes:
  - echoservice
  - echoclient
```

How to run (using the default components folder, eg ~/.dapr/components):

```bash
cp ~/.dapr/components/pubsub.yaml ./components
dapr run --app-port 40002 --app-id echoservice --app-protocol http \
 --components-path ./components
 -- go run service/main.go -port 40002
```

* Pub/sub subscriptions remember which events an appID has received. On Restart only new events are passed to the handler. This would have been a neat feature to ensure no events are missed, if only they arrived in order. It is also possible to set an auto expiry on events. Can it be turned off?
* --app-port and -port are needed by dapr to invoke the /echo endpoint on the service.

## Client Invocation

Client invocation is simple when the client is launched via dapr.

How to run:

```sh
dapr run --app-id echoclient --dapr-http-port 3601 &

dapr publish --publish-app-id echoclient --pubsub pubsub --topic upper --data '{"Text": "Hello world"}'
# or
curl -X POST http://localhost:3601/v1.0/publish/pubsub/echo -H "Content-Type: application/json" -d '{"Text": "hello world"}'
```

Where:

* --app-id identifies the client for logging and scoped access control
* --pubsub (required) is the name of the pubsub component to publish the event to. The component definition can control the scoping and routing of the events.

# 2. Programmatic Method

The programmatic method lets the service define the topics in code at startup.

> flow: client/event -> client-sidecar -[grpc]-> server-sidecar -[http]-> service event handler

When to use:

* Services determine topics at runtime.
* Services remain agnostic of the messaging system used

## Service Setup

Service code snippet ([full example code](service/main.go)):

```go
package main

import daprd "github.com/dapr/go-sdk/service/http"

func main() {
	serviceAddress := ":40001"
	s := daprd.NewService(serviceAddress)
	echoSub := &common.Subscription{PubsubName: "pubsub", Topic: "echo", Route: "/echo"}
	if err := service.AddTopicEventHandler(echoSub, echoEventHandler); err != nil {
		log.Fatalf("error adding topic event handler: %v", err)
	}
}

func echoEventHandler(ctx context.Context, ev *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - topic:%s, ID:%s, Data:%v", ev.Topic, ev.ID, ev.Data)
	// do something with the event here
	return false, nil
}
```

## Client Invocation

Client invocation is simple when the client is launched via dapr.

Client code snippet ([full example code](client/main.go)):

```go
package main

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
)

func main() {
	// Use the internal mechanism to connect to the sidecar. Do not use NewClientWithPort or it will fail. 
	client, _ := dapr.NewClient()
	defer client.Close()
	ctx := context.Background()
	pubsubName := "pubsub"
	cmd := "upper"
	data := []byte(`{"text":"hello world"}'`)
	client.PublishEvent(ctx, pubsubName, cmd, data)
}
```

How to run:

```bash
dapr run --app-port=9102 --app-id echoclient -- go run client/main.go --pubsub=pubsub upper "Hello world"
```

Where:

* --app-port is (presumably) required to send subscriptions to the client (eg app-channel). This part is not understood as the client creation should not use a port and dapr initializes with its own http and grdp ports. Maybe this isn't used?
* --app-id identifies the client for logging and scoped access control
* --pubsub is the name of the pubsub component to publish the event to. The component definition can control the scoping and routing of the events.

Issues/Questions:

1. Without --app-port dapr issues a warning, although events are still accepted still work. What is it for? 
