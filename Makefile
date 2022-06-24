# Makefile to build and test the echo service
DISTGO=./proto/go
PROTOPATH=--proto_path=./proto
GOOPT=--go_opt=paths=source_relative
GOOUT=--go_out=$(DISTGO)
GRPC=--go-grpc_out=$(DISTGO) --go-grpc_opt=paths=source_relative
SERVICE_GRPC_PORT = 40001
DAPR_HTTP_PORT = 9000
DAPR_GRPC_PORT = 9001
PROTOC=protoc $(PROTOPATH) $(GOOPT) $(GOOUT) $(GRPC)
.DEFAULT_GOAL := help

.FORCE:

proto: .FORCE  ## Compile hub protobuffer files for go
	$(PROTOC) proto/echo.proto
	go mod tidy

echo-service: proto ## Compile the gRPC service
	go build -o bin/echo-service cmd/echo-service/main.go

echo-cli: proto ## Compile the gRPC CLI
	go build  -o bin/echo-cli cmd/echo-cli/main.go


help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


run1: echo-cli echo-service ## Run echo gRPC service without dapr
	@bin/echo-service&
	@bin/echo-cli "upper" "Hello gRPC"
	@bin/echo-cli "stop"

run2: echo-cli echo-service ## Run echo service via dapr and use the gRPC client to call the service via dapr
	@dapr run --enable-api-logging \
		--app-protocol grpc --app-port $(SERVICE_GRPC_PORT) \
		--app-id echo \
		--dapr-http-port $(DAPR_HTTP_PORT) --dapr-grpc-port $(DAPR_GRPC_PORT)  --  bin/echo-service &
	@bin/echo-cli -port $(DAPR_GRPC_PORT) upper "Hello gRPC via dapr gRPC"
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