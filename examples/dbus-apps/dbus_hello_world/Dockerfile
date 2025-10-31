# Multi-stage build for D-Bus Hello World Service
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o dbus-hello-world main.go

# Runtime stage
FROM ubuntu:22.04

# Install D-Bus and dependencies
RUN apt-get update && \
    apt-get install -y \
    dbus \
    dbus-user-session \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -u 1000 -s /bin/bash dbususer

# Copy binary from builder
COPY --from=builder /app/dbus-hello-world /usr/local/bin/dbus-hello-world
RUN chmod +x /usr/local/bin/dbus-hello-world

# Create shared directory for D-Bus socket
RUN mkdir -p /shared/dbus && chown dbususer:dbususer /shared/dbus

# Switch to non-root user
USER dbususer

# Set environment for D-Bus session bus
ENV DBUS_SESSION_BUS_ADDRESS=unix:path=/shared/dbus/session_bus_socket

# Run the D-Bus Hello World service
CMD ["/usr/local/bin/dbus-hello-world"]
