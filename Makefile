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

TEST_COMMAND = "echo"
#TEST_PAYLOAD = "Hello world"
TEST_PAYLOAD=$(shell cat test/payload-1K.txt)
#TEST_PAYLOAD=$(shell cat test/payload-10K.txt)
#TEST_PAYLOAD=$(shell cat test/payload-100K.txt)

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


run1:  ## Run plain echo gRPC service without dapr
	go run pkg/plain-grpc/service/main.go &
	sleep 1
	go run pkg/plain-grpc/client/main.go --repeat 1000 $(TEST_COMMAND) "$(TEST_PAYLOAD)"
	sleep 1
	go run pkg/plain-grpc/client/main.go "stop"

run2:  ## Run plain echo http service without dapr
	go run pkg/plain-http/service/main.go &
	sleep 1
	go run pkg/plain-http/client/main.go --repeat 1000 $(TEST_COMMAND) "$(TEST_PAYLOAD)"
	sleep 1
	go run pkg/plain-http/client/main.go "stop"

run3: ## Run echo gRPC service with dapr sidecars
	dapr run --app-protocol grpc --app-port 40001 \
		--app-id echo --dapr-grpc-port 9001 \
		-- go run pkg/plain-grpc/service/main.go --port 40001&
	dapr run --dapr-grpc-port=9002  -- \
		go run pkg/plain-grpc/client/main.go --port 9002 --repeat 1000 $(TEST_COMMAND) "$(TEST_PAYLOAD)"
	go run pkg/plain-grpc/client/main.go "stop"

run4: ## Run echo http service with dapr sidecars
	dapr run --app-protocol http --app-port 40002 \
		--app-id echo --dapr-http-port 9003 \
		-- go run pkg/plain-http/service/main.go --port 40002&
	dapr run --dapr-http-port=9004  -- \
		go run pkg/plain-http/client/main.go --port 9004 --repeat 1000 $(TEST_COMMAND) "$(TEST_PAYLOAD)"
	go run pkg/plain-http/client/main.go "stop"

run5: ## Run echo gRPC service with invoke SDK
	dapr run --app-protocol grpc --app-port 40001 \
		--app-id echo --dapr-grpc-port 9001 -- \
		go run pkg/invoke-grpc/service/main.go --port 40001&
	dapr run --dapr-grpc-port=9002  -- \
		go run pkg/invoke-grpc/client/main.go --repeat 1000 $(TEST_COMMAND) "$(TEST_PAYLOAD)"
	dapr run --dapr-grpc-port=9002  -- \
    	go run pkg/invoke-grpc/client/main.go "stop" >/dev/null

run6: ## Run echo http service with invoke SDK
	dapr run --app-protocol http --app-port 40002 \
		--app-id echo --dapr-http-port 9003 \
		-- go run pkg/invoke-http/service/main.go --port 40002&
	dapr run --dapr-http-port=9004  -- \
		go run pkg/invoke-http/client/main.go --repeat 1000 $(TEST_COMMAND) "$(TEST_PAYLOAD)"
	dapr run --dapr-http-port=9004  -- \
		go run pkg/invoke-http/client/main.go "stop" >/dev/null
