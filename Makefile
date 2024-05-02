
clean:
	go clean

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run

rpc:
	protoc --go_out=proto/messaging/v1 --go_opt=paths=source_relative \
    --go-grpc_out=proto/messaging/v1 --go-grpc_opt=paths=source_relative \
    proto/messaging/messaging_service.proto

test:
	go test -v ./...
