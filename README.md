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

## API Overview

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

# Health check
curl http://localhost:8080/health
```

## Run on Podman and Kubernetes

Isolated session bus (user bus) dedicated to the POD, with no access to the system or host, and without requiring elevated privileges (eliminating related security risks). Only containers within the same POD that share the same user and volume (unix_socket/bus) can access this session bus.

Each container in the POD must implement its own D-Bus interfaces related to its application or service workload. Only these interfaces are exposed through the REST API.

Docker can be used instead of Podman, but Podman is preferred for its POD support and isolation.
In a Docker environment, all containers are able to access the session bus (if configured for it).

## Run on VMs and physical computers (x86-64, ARM-arch32-64 and single-board computers)

The D-Bus (software and libraries from freedesktop.org) works in Unix and Linux systems apart from the init system (Systemd, SysV, OpenRC, BSD rc, etc).
In this scenario, the D-Bus Controller can be run as a service or application that connects to the system and session buses of the host system.
