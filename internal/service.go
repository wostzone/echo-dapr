package internal

import (
	"context"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "github.com/wostzone/echo/proto/go"
)

// EchoService demonstrates how to build a microservice using dapr for pub/sub and http invocation.
type EchoService struct {
	pb.UnimplementedEchoServiceServer
	dapr.UnimplementedAppCallbackServer

	// keep the grpc server used to stop the service
	grpcServer *grpc.Server
}

//--- EchoService methods

func (service *EchoService) Echo(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	return args, nil
}

func (service *EchoService) UpperCase(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	upper := strings.ToUpper(args.Text)
	response := pb.TextParam{Text: upper}
	return &response, nil
}

func (service *EchoService) Reverse(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	rns := []rune(args.Text)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	response := pb.TextParam{Text: string(rns)}
	return &response, nil
}

// Stop the service, give it some time for a response
func (service *EchoService) Stop(ctx context.Context, args *empty.Empty) (*pb.TextParam, error) {
	go func() {
		time.Sleep(time.Millisecond * 100)
		service.grpcServer.Stop()
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
// grpcServer is used to stop the service
func NewEchoService(grpcServer *grpc.Server) *EchoService {
	service := &EchoService{grpcServer: grpcServer}
	return service
}
