package internal

import (
	"context"
	"fmt"

	commonv1pb "github.com/dapr/dapr/pkg/proto/common/v1"
	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/anypb"

	echo "github.com/wostzone/echo/proto/go"
)

// DAPR pub/sub handlers for the service

//--- UnimplementedAppCallbackServer methods
//OnInvoke(context.Context, *common.InvokeRequest) (*common.InvokeResponse, error)
//ListTopicSubscriptions(context.Context, *emptypb.Empty) (*runtime.ListTopicSubscriptionsResponse, error)
//OnTopicEvent(context.Context, *runtime.TopicEventRequest) (*runtime.TopicEventResponse, error)
//ListInputBindings(context.Context, *emptypb.Empty) (*runtime.ListInputBindingsResponse, error)
//OnBindingEvent(context.Context, *runtime.BindingEventRequest) (*runtime.BindingEventResponse, error)
//mustEmbedUnimplementedAppCallbackServer()

// OnInvoke to test invocation via dapr HTTP???
func (service *EchoService) OnInvoke(ctx context.Context, in *commonv1pb.InvokeRequest) (*commonv1pb.InvokeResponse, error) {
	fmt.Println("OnInvoke. Method ", in.Method)

	var response string

	// Convert the http parameters to grpc parameter for use by the service
	switch in.Method {
	case "echo":
		params := echo.TextParam{Text: in.Data.String()}
		out, _ := service.Echo(ctx, &params)
		response = out.Text
	case "upper":
		params := echo.TextParam{Text: in.Data.String()}
		out, _ := service.UpperCase(ctx, &params)
		response = out.Text
	case "reverse":
		params := echo.TextParam{Text: in.Data.String()}
		out, _ := service.Reverse(ctx, &params)
		response = out.Text
	case "stop":
		service.Stop(ctx, &empty.Empty{})
		response = "stopped"
	default:
		fmt.Println("ignored unknown method: " + in.Method)
	}
	return &commonv1pb.InvokeResponse{
		ContentType: "text/plain; charset=UTF-8",
		Data:        &anypb.Any{Value: []byte(response)},
	}, nil
}

// ListTopicSubscriptions Dapr will call this method to get the list of topics the app wants to subscribe to.
// In this example, we are telling Dapr to subscribe to a topic named TopicA
func (service *EchoService) ListTopicSubscriptions(ctx context.Context, in *empty.Empty) (*pb.ListTopicSubscriptionsResponse, error) {
	fmt.Println("ListTopicSubscriptions")
	return &pb.ListTopicSubscriptionsResponse{
		Subscriptions: []*pb.TopicSubscription{
			{Topic: "echo"},
			{Topic: "upper"},
			{Topic: "reverse"},
			{Topic: "stop"},
		},
	}, nil
}

// ListInputBindings Dapr will call this method to get the list of bindings the app will get invoked by. In this example, we are telling Dapr
// To invoke our app with a binding named storage
func (service *EchoService) ListInputBindings(ctx context.Context, in *empty.Empty) (*pb.ListInputBindingsResponse, error) {
	fmt.Println("ListInputBindings")
	return &pb.ListInputBindingsResponse{
		Bindings: []string{"storage"},
	}, nil
}

// OnBindingEvent This method gets invoked every time a new event is fired from a registered binding. The message carries the binding name, a payload and optional metadata
func (service *EchoService) OnBindingEvent(ctx context.Context, in *pb.BindingEventRequest) (*pb.BindingEventResponse, error) {
	fmt.Println("Invoked from binding")
	return &pb.BindingEventResponse{}, nil
}

// OnTopicEvent This method is fired whenever a message has been published to a topic that has been subscribed. Dapr sends published messages in a CloudEvents 0.3 envelope.
func (service *EchoService) OnTopicEvent(ctx context.Context, in *pb.TopicEventRequest) (*pb.TopicEventResponse, error) {
	fmt.Println("Topic message arrived")
	return &pb.TopicEventResponse{}, nil
}
