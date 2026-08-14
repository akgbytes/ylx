.PHONY: build clean format lint migrate rollback run typecheck seed

build:
	@mkdir -p bin
	@go build -o bin/api ./cmd/api

run:
	@go run ./cmd/api

format:
	@go mod tidy
	@golangci-lint fmt

lint:
	@golangci-lint run

clean:
	@rm -rf ./bin

typecheck:
	@go vet ./...

migrate:
	@go run ./cmd/migrator up

rollback:
	@go run ./cmd/migrator down 1

seed:
	@go run ./cmd/seeder $(SEED_ARGS)
