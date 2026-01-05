# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies needed for build (if any)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ink-flow main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for external API calls
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/ink-flow .

# Expose ports (Application and Prometheus)
EXPOSE 8080 8081

# Command to run the application
CMD ["./ink-flow"]
