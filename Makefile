.PHONY: build start run clean test install dev deps add

build:
	@go build -o bin/go_rest_api

start: build
	@./bin/go_rest_api

run:
	@go run main.go

dev:
	@air -c air.toml

clean:
	@echo Cleaning up...
ifeq ($(OS),Windows_NT)
	@if exist bin rmdir /s /q bin
	@if exist .air.toml del .air.toml
	@if exist tmp rmdir /s /q tmp
else
	@rm -rf bin/
	@rm -rf .air.toml
	@rm -rf tmp/
endif
	@go clean
	@go mod tidy
	@echo Clean complete!

test:
	@go test ./...

install-air:
	@go install github.com/cosmtrek/air@latest

add:
	@go get -u ./...

## deps: Download modules
deps:
	@go mod download