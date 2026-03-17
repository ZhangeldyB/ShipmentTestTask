PROTOC          := protoc
PROTOC_GEN_GO   := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc
MODULE          := github.com/ZhangeldyB/ShipmentTestTask

.PHONY: proto build test test-race up down logs lint

## proto: regenerate gRPC stubs from proto/shipment.proto
proto:
	$(PROTOC) \
		--go_out=. \
		--go_opt=module=$(MODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(MODULE) \
		proto/shipment.proto

## build: compile the server binary to bin/server
build:
	go build -o bin/server ./cmd/server

## test: run all unit tests
test:
	go test ./...

## test-race: run all tests with the race detector enabled
test-race:
	go test -race ./...

## up: build Docker images and start all services (MongoDB + shipment-service)
up:
	docker-compose up --build

## down: stop and remove containers
down:
	docker-compose down

## logs: tail logs for the shipment-service container
logs:
	docker-compose logs -f shipment-service

## lint: run golangci-lint (install: https://golangci-lint.run/usage/install)
lint:
	golangci-lint run ./...
