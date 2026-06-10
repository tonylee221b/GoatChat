APP_NAME := goat-chat
CMD_DIR := ./cmd/api/
BUILD_DIR := ./bin

.PHONY: help run build test test-race fmt vet tidy check clean

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
