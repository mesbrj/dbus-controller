# D-Bus Hello World Service

A simple Go application that exposes a D-Bus service with properties and methods.

## Service Details

- **Service Name**: `com.example.HelloWorld`
- **Object Path**: `/com/example/HelloWorld`
- **Interface**: `com.example.HelloWorld`

## Features

### Properties (Read-only)
- `Data` (string): Returns "Hello World"
- `PID` (int32): Returns the process ID of the service

### Methods
- `Hello(name string)`: Returns "Hello, {name}!"

## Building and Running

### Option 1: Docker/Podman (Recommended)

1. **Build the Docker image:**
   ```bash
   cd examples/dbus-apps/dbus_hello_world
   podman build -t dbus-hello-world:latest .
   ```
   
   **Note**: The Dockerfile uses `dbus-user-session` (not `dbus-x11`) for proper containerized D-Bus session management.

2. **Run the service in a container:**
   ```bash
   podman run --rm -d --name dbus-hello-world dbus-hello-world:latest
   ```

3. **Test the service:**
   ```bash
   # Test Hello method
   podman exec dbus-hello-world dbus-send --session --print-reply --dest=com.example.HelloWorld /com/example/HelloWorld com.example.HelloWorld.Hello string:"World"
   
   # Test Data property
   podman exec dbus-hello-world dbus-send --session --print-reply --dest=com.example.HelloWorld /com/example/HelloWorld org.freedesktop.DBus.Properties.Get string:"com.example.HelloWorld" string:"Data"
   ```

4. **View logs:**
   ```bash
   podman logs dbus-hello-world
   ```

5. **Stop the service:**
   ```bash
   podman stop dbus-hello-world
   ```

### Option 2: Local Build

1. **Download dependencies:**
   ```bash
   go mod tidy
   ```

2. **Build the application:**
   ```bash
   go build -o dbus-hello-world .
   ```

3. **Run the service:**
   ```bash
   ./dbus-hello-world
   ```

## Testing with D-Bus Commands

### 1. List Available Services
```bash
dbus-send --session --print-reply --dest=org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.ListNames
```

### 2. Call Hello Method
```bash
dbus-send --session --print-reply --dest=com.example.HelloWorld /com/example/HelloWorld com.example.HelloWorld.Hello string:"World"
```

### 3. Get Data Property
```bash
dbus-send --session --print-reply --dest=com.example.HelloWorld /com/example/HelloWorld org.freedesktop.DBus.Properties.Get string:"com.example.HelloWorld" string:"Data"
```

### 4. Get PID Property
```bash
dbus-send --session --print-reply --dest=com.example.HelloWorld /com/example/HelloWorld org.freedesktop.DBus.Properties.Get string:"com.example.HelloWorld" string:"PID"
```

### 5. Introspect Service
```bash
dbus-send --session --print-reply --dest=com.example.HelloWorld /com/example/HelloWorld org.freedesktop.DBus.Introspectable.Introspect
```

## Testing with D-Bus Controller API

Once the service is running, you can also test it using the D-Bus Controller REST API:

1. **List session services** (should include com.example.HelloWorld):
   ```bash
   curl http://localhost:8080/buses/session/services
   ```

2. **Get service info**:
   ```bash
   curl http://localhost:8080/buses/session/services/com.example.HelloWorld
   ```

3. **Call Hello method**:
   ```bash
   curl -X POST http://localhost:8080/buses/session/services/com.example.HelloWorld/interfaces/com.example.HelloWorld/methods/Hello/call \
     -H "Content-Type: application/json" \
     -d '{"arguments": ["World"]}'
   ```

4. **Get properties**:
   ```bash
   curl http://localhost:8080/buses/session/services/com.example.HelloWorld/interfaces/com.example.HelloWorld/properties/Data
   curl http://localhost:8080/buses/session/services/com.example.HelloWorld/interfaces/com.example.HelloWorld/properties/PID
   ```

## Expected Output

When you call the Hello method with "World", you should get:
```
"Hello, World!"
```

The Data property should return:
```
"Hello World"
```

The PID property should return the process ID of the running service.