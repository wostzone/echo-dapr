package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/wostzone/echo/proto/go"
)

// EchoService demonstrates how to build a microservice for grpc, pub/sub and http invocation.
type EchoService struct {
	pb.UnimplementedEchoServiceServer

	// handle stop request
	stopHandler func()
}

//--- EchoService methods

func (service *EchoService) Echo(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	if args == nil {
		return nil, fmt.Errorf("Missing args")
	}
	fmt.Println("EchoService.Echo: ", args.Text)
	return args, nil
}

func (service *EchoService) UpperCase(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	if args == nil {
		return nil, fmt.Errorf("Missing args")
	}
	upper := strings.ToUpper(args.Text)
	response := pb.TextParam{Text: upper}
	fmt.Println("EchoService.UpperCase: ", response.Text)
	return &response, nil
}

func (service *EchoService) Reverse(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	if args == nil {
		return nil, fmt.Errorf("Missing args")
	}
	rns := []rune(args.Text)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	response := pb.TextParam{Text: string(rns)}
	fmt.Println("EchoService.Reverse: ", response.Text)
	return &response, nil
}

// Stop the service, give it some time for a response
func (service *EchoService) Stop(ctx context.Context, args *empty.Empty) (*pb.TextParam, error) {
	fmt.Println("EchoService.Stop")

	go func() {
		time.Sleep(time.Millisecond * 100)
		if service.stopHandler != nil {
			service.stopHandler()
		}
	}()
	response := pb.TextParam{Text: "Stopped"}
	return &response, nil
}

//// QueryTDs return a collection of TDs that match the query parameter
//func (service *ThingDirServer) QueryTDs(ctx context.Context, args *directory.QueryTDs_Args) (*directory.TDList_Result, error) {
//	res := &directory.TDList_Result{}
//	return res, nil
//}
//
//// ListTDs returns the collection of known TDs
//func (service *ThingDirServer) ListTDs(ctx context.Context, args *directory.LimitOffset_Args) (*directory.TDList_Result, error) {
//	res := &directory.TDList_Result{}
//	res.Tds = make([]*thing.TD, 0)
//	res.Tds = append(res.Tds, &thing.TD{ID: "thing1"})
//	res.Tds = append(res.Tds, &thing.TD{ID: "thing2"})
//	return res, nil
//}
//
//// UpdateTD adds a TD
//func (service *ThingDirServer) UpdateTD(ctx context.Context, args *thing.TD) (*empty.Empty, error) {
//	return &empty.Empty{}, nil
//}

// NewEchoService creates and registers the service with gRPC interface
// onShutDown callback is used to handle stop request
func NewEchoService(stopHandler func()) *EchoService {
	service := &EchoService{stopHandler: stopHandler}
	return service
}
