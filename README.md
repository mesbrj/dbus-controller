# D-Bus Controller

A REST API server for introspecting and controlling the Linux D-Bus system.

## Features

- **REST API**: Clean, RESTful interface for D-Bus operations
- **System & Session Bus Support**: Access both system and session D-Bus buses
- **Real-time Introspection**: Dynamic discovery of services, interfaces, methods, properties, and signals
- **Method Execution**: Call D-Bus methods via HTTP POST requests
- **Property Management**: Get and set D-Bus properties via REST endpoints
- **Signal Monitoring**: Subscribe to D-Bus signals (work in progress)
- **No Persistence**: All data is introspected at runtime for real-time accuracy
- **OpenAPI Documentation**: Auto-generated API documentation via Fuego
>
- **Framework**: [Go Fuego](https://github.com/go-fuego/fuego) - Modern Go web framework
- **D-Bus Library**: [godbus](https://github.com/godbus/dbus) - Pure Go D-Bus library

## Quick Start

### Prerequisites

- Go 1.21 or later
- Linux with D-Bus installed
- Access to system and/or session D-Bus (for testing)

### Installation

```bash
# Clone the repository
git clone https://github.com/mesbrj/dbus-controller.git
cd dbus-controller

# Install dependencies
make deps

# Build and run
make run
```

The server will start on `http://localhost:8080`.

### Docker

```bash
# Build and run with Docker
make docker-run

# Or use docker-compose
make docker-compose-up
```

## API Overview

### Base URL
```
http://localhost:8080
```

### Main Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/buses` | GET | List available buses (system, session) |
| `/buses/{busType}/services` | GET | List services on a bus |
| `/buses/{busType}/services/{service}` | GET | Get service information |
| `/buses/{busType}/services/{service}/interfaces` | GET | List service interfaces |
| `/buses/{busType}/services/{service}/interfaces/{interface}/methods` | GET | List interface methods |
| `/buses/{busType}/services/{service}/interfaces/{interface}/methods/{method}/call` | POST | Call a method |
| `/buses/{busType}/services/{service}/interfaces/{interface}/properties` | GET | List interface properties |
| `/buses/{busType}/services/{service}/interfaces/{interface}/properties/{property}` | GET/PUT | Get/Set property value |
| `/buses/{busType}/services/{service}/interfaces/{interface}/signals` | GET | List interface signals |
| `/buses/{busType}/services/{service}/introspect` | GET | Get service introspection XML |

### Example Usage

```bash
# List all buses
curl http://localhost:8080/buses

# List services on system bus
curl http://localhost:8080/buses/system/services

# Get service information
curl http://localhost:8080/buses/system/services/org.freedesktop.DBus

# Call a method
curl -X POST http://localhost:8080/buses/system/services/org.freedesktop.DBus/interfaces/org.freedesktop.DBus/methods/Hello/call \
  -H "Content-Type: application/json" \
  -d '{"args": []}'

# Get a property
curl http://localhost:8080/buses/system/services/org.freedesktop.DBus/interfaces/org.freedesktop.DBus/properties/Features
```

## Development

### Docker

```bash
docker run -d \
  --name dbus-controller \
  --privileged \
  -p 8080:8080 \
  -v /var/run/dbus:/var/run/dbus:ro \
  mesbrj/dbus-controller:latest
```

### Systemd Service

Create `/etc/systemd/system/dbus-controller.service`:

```ini
[Unit]
Description=D-Bus Controller API
After=network.target

[Service]
Type=simple
User=dbus-controller
ExecStart=/usr/local/bin/dbus-controller
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Health Checks

The service provides health check endpoints:
```bash
curl http://localhost:8080/health
```

