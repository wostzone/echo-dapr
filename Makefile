# Makefile to build and test the echo service
DISTGO=./proto/go
PROTOPATH=--proto_path=./proto
GOOPT=--go_opt=paths=source_relative
GOOUT=--go_out=$(DISTGO)
GRPC=--go-grpc_out=$(DISTGO) --go-grpc_opt=paths=source_relative
SERVICE_GRPC_PORT = 40001
SERVICE_HTTP_PORT = 40002
DAPR_GRPC_PORT = 9001
DAPR_HTTP_PORT = 9002
PROTOC=protoc $(PROTOPATH) $(GOOPT) $(GOOUT) $(GRPC)
.DEFAULT_GOAL := help

.FORCE:

proto: .FORCE  ## Compile hub protobuffer files for go
	$(PROTOC) proto/echo.proto
	go mod tidy


invoke-grpc: proto ## Compile the echo invoke grpc client and service
	go build -o bin/invoke-grpc-service pkg/invoke-grpc/service/main.go
	go build -o bin/invoke-grpc-client pkg/invoke-grpc/client/main.go

invoke-http: proto ## Compile the echo invoke http client and service
	go build -o bin/invoke-http-service pkg/invoke-http/service/main.go
	go build -o bin/invoke-http-client pkg/invoke-http/client/main.go

pubsub-grpc: proto  ## Compile the echo pub/sub grpc client and service
	go build -o bin/pubsub-grpc-service pkg/pubsub-grpc/service/main.go
	go build -o bin/pubsub-grpc-client pkg/pubsub-grpc/client/main.go

pubsub-http: proto  ## Compile the echo pub/sub http client and service
	go build -o bin/pubsub-http-service pkg/pubsub-http/service/main.go
	go build -o bin/pubsub-http-client pkg/pubsub-http/client/main.go




help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


run1: echo-cli echo-service ## Run echo gRPC service without dapr
	@bin/echo-service&
	@bin/echo-cli "upper" "Hello gRPC"
	@bin/echo-cli "stop"

run2: echo-cli echo-service ## Run echo service via dapr: client->dapr[40002]->dapr[echo]->service[40001]
	dapr run --enable-api-logging \
		--app-protocol grpc --app-port $(SERVICE_GRPC_PORT) \
		--app-id echo \
		--dapr-http-port $(DAPR_HTTP_PORT) --dapr-grpc-port $(DAPR_GRPC_PORT)  --  bin/echo-service &
	@#bin/echo-cli -port $(DAPR_GRPC_PORT) upper "Hello gRPC via dapr client gRPC"
	@dapr run --dapr-grpc-port=40002  -- bin/echo-cli -port 40002 upper "hello via client->dapr->dapr->service"
	@bin/echo-cli -port $(DAPR_GRPC_PORT) "stop"

run3: echo-cli echo-service ## Run echo service via dapr and use curl to invoke the http API
	@dapr run --enable-api-logging \
		--app-protocol grpc --app-port $(SERVICE_GRPC_PORT) \
		--app-id echo \
		--dapr-http-port $(DAPR_HTTP_PORT) --dapr-grpc-port $(DAPR_GRPC_PORT)  --  bin/echo-service &
	@sleep 2
	@curl localhost:$(DAPR_HTTP_PORT)/v1.0/invoke/echo/method/upper -d "Hello world" -X PUT
	@echo
	@curl localhost:$(DAPR_HTTP_PORT)/v1.0/invoke/echo/method/stop -X PUT

run4: echo-cli echo-service ## Run echo service via dapr pub/sub
	@dapr run --enable-api-logging \
		--app-protocol grpc --app-port $(SERVICE_GRPC_PORT) \
		--app-id echo \
		--dapr-http-port $(DAPR_HTTP_PORT) --dapr-grpc-port $(DAPR_GRPC_PORT)  --  bin/echo-service &
	@sleep 2
	#// todo@curl localhost:$(DAPR_HTTP_PORT)/v1.0/invoke/echo/method/upper -d "Hello world" -X PUT
	@echo
	@curl localhost:$(DAPR_HTTP_PORT)/v1.0/invoke/echo/method/stop -X PUT
