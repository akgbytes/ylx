.PHONY: build run format clean migrate rollback

build:
	@go build -o bin/api ./cmd/api

run:
	@go run ./cmd/api

format:
	@go fmt ./...
	@go mod tidy

clean:
	@rm -rf ./bin

migrate:
	@go run ./cmd/migrator/main.go up

rollback:
	@go run ./cmd/migrator/main.go down 1
