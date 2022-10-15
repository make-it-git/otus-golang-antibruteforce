DIRECTORY=$(shell pwd)

build:
	CGO_ENABLED=0 go build -o ./bin/antibruteforce ./cmd/service

run:
	docker-compose -f deployment/docker-compose.yaml up --build -d

stop:
	docker-compose -f deployment/docker-compose.yaml stop

install-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

generate:
	protoc --go-grpc_out=pkg/antibruteforce --go_out=pkg/antibruteforce api/*.proto

test:
	go test -race -v -count=100 ./...

integration-test:
	docker compose -f deployment/docker-compose.test.yaml up --exit-code-from antibruteforce-test

lint:
	docker run --rm -v $(DIRECTORY):/app -w /app golangci/golangci-lint:v1.50.0 golangci-lint run -v --timeout 2m

.PHONY: build run install-protoc generate test integration-test lint
