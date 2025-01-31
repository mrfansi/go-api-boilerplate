FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

# Copy any additional required files
COPY --from=builder /app/.env.example .env

# Create data directory for SQLite
RUN mkdir -p /app/data && \
    chown -R nobody:nobody /app

# Use non-root user
USER nobody

# Expose port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./api"]