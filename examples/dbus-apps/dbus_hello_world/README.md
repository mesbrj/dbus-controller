# D-Bus Hello World Service

A simple Go application that exposes a D-Bus service with properties and methods.

## Service Details

- **Service Name**: `com.example.HelloWorld`
- **Object Path**: `/com/example/HelloWorld`
- **Interface**: `com.example.HelloWorld`

## Features

### Properties (Read-only)
- `Data` (string): Returns "Hello World"
- `PID` (int32): Returns the process ID of the HelloWorld running service

### Methods
- `Hello(name string)`: Returns "Hello, {name}!"

## Example HelloWorld Service Pod Deployment (podman)

```bash
# On the root of the repository
podman build -f Dockerfile -t dbus-controller:latest
cd examples/dbus-apps/dbus_hello_world
podman build -f Dockerfile -t dbus-hello-world:latest
# Publishes all containerPort definitions
podman play kube --publish-all dbus-hello-world-pod.yaml
```

## Example HelloWorld Service Pod Deployment (kubernetes)

 ...

## Expected results and behavior

- **controller-service** Container:
    - D-Bus session daemon
    - D-Bus Controller API listening on port 8080

![](/docs/controller-service.png)
>

- **dbus-hello-world-service** Container:
    - Registers the `com.example.HelloWorld` service on the D-Bus session bus
    - Responds via D-Bus to property and method calls as defined

![](/docs/dbus-hello-world-service.png)
>

- **dbus-hello-world-client** Container:
    - Connects to the D-Bus session bus
    - Enter a loop of "testing" the HelloWorld service every 30 seconds:
        - Calls (via D-Bus) and prints the `Data` and `PID` properties from the `com.example.HelloWorld` service
        - Calls (via D-Bus) and prints the result of the `Hello` method with a sample name: Client-<test_round_number>

![](/docs/dbus-hello-world-client.png)
>

- **D-Bus Controller API**:
    - Accessible via HTTP on port 8080
    - Listenning only in the `controller-service` container
    - Will perform D-Bus operations on `com.example.HelloWorld` service running in the `dbus-hello-world-service` container

![](/docs/d-bus-controller-api.png)
>

- **D-Bus Controller API**: Bugs starting here:

![](/docs/bugs-d-bus-controller-api.png)
