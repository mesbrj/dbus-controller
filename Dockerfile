# Multi-stage Dockerfile for D-Bus Controller Application
# Stage 1: Build the Go application
FROM golang:1.23-bookworm AS builder

WORKDIR /build

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dbus-controller ./cmd/server

# Stage 2: Runtime image with D-Bus and our application
FROM ubuntu:22.04

# Install D-Bus and required tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    dbus \
    dbus-user-session \
    ca-certificates \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user for security
RUN groupadd -g 1000 dbus-user && \
    useradd -u 1000 -g 1000 -m -s /bin/bash dbus-user

# Copy the built application from builder stage
COPY --from=builder /build/dbus-controller /usr/local/bin/dbus-controller

# Make the binary executable
RUN chmod +x /usr/local/bin/dbus-controller

# Create necessary directories with proper permissions
RUN mkdir -p /tmp/dbus-home/.config/dbus-1/session.d && \
    mkdir -p /tmp/docs && \
    chown -R dbus-user:dbus-user /tmp/dbus-home /tmp/docs

# Switch to non-root user
USER dbus-user

# Set environment variables
ENV HOME=/tmp/dbus-home
ENV CGO_ENABLED=0
ENV GOOS=linux

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/buses || exit 1

# Default command (can be overridden)
CMD ["/usr/local/bin/dbus-controller"]