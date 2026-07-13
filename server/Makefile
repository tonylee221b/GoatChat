APP_NAME := goat-chat
ENV ?= local
ENV_FILE := .env.$(ENV)
-include $(ENV_FILE)
export

CMD_DIR := ./cmd/api/
BUILD_DIR := ./bin
DATABASE_URL ?= postgres://goatchat:goatchat@localhost:5432/goatchat?sslmode=disable
MIGRATIONS_DIR := ./db/migration
COMPOSE := docker compose --env-file $(ENV_FILE)
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.27.1
SQLC := go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

.PHONY: help run build test test-race fmt vet tidy check clean db-up db-down db-migrate-up sqlc-generate require-env-file

help:
	@echo "Available targets:"
	@echo "  make run ENV=local          - Run the application from $(CMD_DIR)"
	@echo "  make build                  - Build the application binary into $(BUILD_DIR)"
	@echo "  make test                   - Run all Go tests"
	@echo "  make test-race              - Run all Go tests with the race detector"
	@echo "  make fmt                    - Format all Go packages"
	@echo "  make vet                    - Run go vet on all packages"
	@echo "  make tidy                   - Run go mod tidy"
	@echo "  make check                  - Run fmt, vet, and test"
	@echo "  make db-up ENV=local        - Start PostgreSQL and Redis with Docker Compose"
	@echo "  make db-down ENV=local      - Stop PostgreSQL and Redis"
	@echo "  make db-migrate-up ENV=local - Apply PostgreSQL migration"
	@echo "  make sqlc-generate          - Generate typed database code"
	@echo "  make clean                  - Remove build artifacts"
	@echo ""
	@echo "Environment:"
	@echo "  ENV defaults to local. Use ENV=dev or ENV=prod to load .env.dev or .env.prod."

require-env-file:
	@test -f "$(ENV_FILE)" || { echo "Missing $(ENV_FILE). Use ENV=local, ENV=dev, or ENV=prod."; exit 1; }

run: require-env-file
	go run $(CMD_DIR)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

test:
	gotestsum --format testname -- -count=1 ./...

test-race:
	gotestsum --format testname -- -race ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

check: fmt vet test

clean:
	rm -rf $(BUILD_DIR)

db-up: require-env-file
	$(COMPOSE) up -d postgres redis

db-down: require-env-file
	$(COMPOSE) down

db-migrate-up: require-env-file
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

sqlc-generate:
	$(SQLC) generate
