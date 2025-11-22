# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies including sqlc
RUN apk add --no-cache git make && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Run tests
RUN go test -v ./...

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/db/schema.sql /app/db/schema.sql

# Expose ports
EXPOSE 8080 50051

# Run the application
CMD ["/app/api"]
