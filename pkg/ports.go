package pkg

// ports used in the echo demos

const (
	// EchoServiceAppID app ID of the echo service
	EchoServiceAppID = "echo"
	// EchoServiceGrpcPort with service gRPC listening port
	EchoServiceGrpcPort = 40001
	// EchoServiceHttpPort with service http listening port
	EchoServiceHttpPort = 40002

	// PubSubServiceGrpcPort listening for echo requests with gRPC handler
	PubSubServiceGrpcPort = 40003
	// PubSubServiceHttpPort listening for echo requests with http handler
	PubSubServiceHttpPort = 40004

	// EchoDaprServiceGrpcPort with service sidecar GRPC port
	EchoDaprServiceGrpcPort = 9001
	// EchoDaprServiceHttpPort with service sidecar HTTP port
	EchoDaprServiceHttpPort = 9002
	// EchoDaprClientGrpcPort with client sidecar GRPC port
	EchoDaprClientGrpcPort = 9101
	// EchoDaprClientHttpPort with client sidecar HTTP port
	EchoDaprClientHttpPort = 9102

	// PubSubDaprGrpcPort with pubsub sidecar gRPC port
	//PubSubDaprGrpcPort = 9003
	// PubSubDaprHttpPort with pubsub sidecar HTTP port
	//PubSubDaprHttpPort = 9004
)
