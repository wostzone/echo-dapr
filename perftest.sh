#!/bin/bash

echo make sure the service is running first:
# dapr run --app-protocol grpc --app-port 40001 --enable-profiling --profile-port 7777 \
# --app-id echo --dapr-grpc-port 9001 -- pkg/plain-grpc/grpc-service --port 40001

curl http://localhost:7777/debug/pprof/profile?seconds=120 > echo-service.pprof &

dapr run --app-id echoclient --dapr-grpc-port=9002 \
--enable-profiling --profile-port 7778 -- \
pkg/plain-grpc/grpc-client --port 9002 --repeat 55000 echo "Hello world" &

sleep 1
curl http://localhost:7778/debug/pprof/profile?seconds=120 > echo-client.pprof


# generate profiling PDF
go tool pprof --pdf pkg/plain-grpc/grpc-client ./echo-client.pprof > echo-client.pdf
go tool pprof --pdf pkg/plain-grpc/grpc-service ./echo-service.pprof > echo-service.pdf
