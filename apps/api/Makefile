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
	@go run ./cmd/migrator up

rollback:
	@go run ./cmd/migrator down 1

seed:
	@go run ./cmd/seeder

seed-listings:
	@go run ./cmd/seeder --file=001_listings.sql
