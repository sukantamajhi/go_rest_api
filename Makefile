.PHONY: build run clean test install dev deps

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
	@rm -rf .air.toml
	@rm -rf tmp/
	@go clean

test:
	@go test ./...

## deps: Download modules
deps:
	@go mod download