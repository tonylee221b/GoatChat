APP_NAME := goat-chat
CMD_DIR := ./cmd/api/
BUILD_DIR := ./bin
DATABASE_URL ?= postgres://goatchat:goatchat@localhost:5432/goatchat?sslmode=disable
MIGRATIONS_DIR := ./db/migration
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.27.1
SQLC := go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

.PHONY: help run build test test-race fmt vet tidy check clean db-up db-down db-migrate-up sqlc-generate

help:
	@echo "Available targets:"
	@echo "  make run        - Run the application from $(CMD_DIR)"
	@echo "  make build      - Build the application binary into $(BUILD_DIR)"
	@echo "  make test       - Run all Go tests"
	@echo "  make test-race  - Run all Go tests with the race detector"
	@echo "  make fmt        - Format all Go packages"
	@echo "  make vet        - Run go vet on all packages"
	@echo "  make tidy       - Run go mod tidy"
	@echo "  make check      - Run fmt, vet, and test"
	@echo "  make db-up      - Start PostgreSQL and Redis with Docker Compose"
	@echo "  make db-down    - Stop PostgreSQL and Redis"
	@echo "  make db-migrate-up     - Apply PostgreSQL migration"
	@echo "  make sqlc-generate     - Generate typed database code"
	@echo "  make clean      - Remove build artifacts"

run:
	go run $(CMD_DIR)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

check: fmt vet test

clean:
	rm -rf $(BUILD_DIR)

db-up:
	docker compose up -d postgres redis

db-down:
	docker compose down

db-migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

sqlc-generate:
	$(SQLC) generate
