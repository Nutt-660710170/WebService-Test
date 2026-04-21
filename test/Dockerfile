# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/oil-price-api ./cmd/main.go

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Install runtime dependencies (timezone data, CA certificates)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/oil-price-api .

# Expose port
EXPOSE 3000

# Run the application
CMD ["./oil-price-api"]
