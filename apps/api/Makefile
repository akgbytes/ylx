.PHONY: build run format format-check lint clean migrate rollback

build:
	@go build -o bin/api ./cmd/api

run:
	@go run ./cmd/api

format:
	@go mod tidy
	@golangci-lint fmt

format-check:
	@golangci-lint fmt --diff
	@go mod tidy -diff

lint:
	@golangci-lint run

clean:
	@rm -rf ./bin

migrate:
	@go run ./cmd/migrator/main.go up

rollback:
	@go run ./cmd/migrator/main.go down 1
