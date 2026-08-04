.PHONY: build run format clean

build:
	@ go build -o bin/api ./cmd/api

run:
	@ go run ./cmd/api

format:
	@ go fmt ./...
	@ go mod tidy

clean:
	@ rm -rf ./bin
