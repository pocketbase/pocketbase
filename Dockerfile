# Dockerfile
# Stage 1: Clone the repository and set up the development environment
FROM golang:1.21-alpine AS installer
WORKDIR /app

# Install git and essential build tools
RUN apk add --no-cache git make build-base

# Clone the PocketBase repository
RUN git clone https://github.com/pocketbase/pocketbase.git . && \
    # Create the necessary directory structure
    mkdir -p cmd/pocketbase && \
    # Create main.go in the correct location
    echo 'package main

import (
    "log"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/cmd"
)

func main() {
    app := pocketbase.New()
    serveCmd := cmd.NewServeCommand(app, true)
    app.RootCmd.AddCommand(serveCmd)
    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}' > cmd/pocketbase/main.go && \
    # Install dependencies
    go mod download && \
    go mod verify && \
    go mod tidy

# Stage 2: Build the application
FROM installer AS builder
WORKDIR /app

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o bin/pocketbase \
    cmd/pocketbase/main.go

# Stage 3: Create the minimal runtime image
FROM alpine:latest AS runner
WORKDIR /app

# Install certificates for HTTPS support
RUN apk add --no-cache ca-certificates && \
    # Create data directory
    mkdir -p pb_data && \
    # Add non-root user for security
    adduser -D pocketbase

# Copy only the necessary binary from the builder stage
COPY --from=builder /app/bin/pocketbase .

# Set ownership for security
RUN chown -R pocketbase:pocketbase /app

# Switch to non-root user
USER pocketbase

# Expose the default PocketBase port
EXPOSE 8090

# Set default environment variables
ENV PB_ENCRYPTION_KEY=""

# Start PocketBase
CMD ["./pocketbase", "serve", "--http=0.0.0.0:8090", "--origins=*"]
