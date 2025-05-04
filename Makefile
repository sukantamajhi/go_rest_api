.PHONY: build run clean test install dev

build:
	@go build -o bin/go_rest_api

start: build
	@./bin/go_rest_api

run:
	@go run main.go

dev:
	@air

clean:
	@rm -rf bin/
	@go clean

test:
	@go test ./...

install:
	@go install 