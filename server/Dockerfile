# syntax=docker/dockerfile:1
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Enable caching of Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

# Copy source code
COPY . .

# Build binary utilizing cache mounts for packages and go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build/ \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/goat-chat ./cmd/api/main.go

# Production Stage - use a minimal, secure image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/goat-chat .

EXPOSE 8888

ENTRYPOINT ["./goat-chat"]
