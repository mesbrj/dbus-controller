package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	serviceName   = "com.example.HelloWorld"
	objectPath    = "/com/example/HelloWorld"
	interfaceName = "com.example.HelloWorld"
)

// HelloWorldService represents D-Bus service
type HelloWorldService struct {
	data string
	pid  int32
}

// NewHelloWorldService creates a new instance of the service
func NewHelloWorldService() *HelloWorldService {
	return &HelloWorldService{
		data: "Hello World",
		pid:  int32(os.Getpid()),
	}
}

// Hello method - responds with greeting message
func (h *HelloWorldService) Hello(name string) (string, *dbus.Error) {
	response := fmt.Sprintf("Hello, %s!", name)
	fmt.Printf("Hello method called with name: %s, responding: %s\n", name, response)
	return response, nil
}

// GetData returns the data property
func (h *HelloWorldService) GetData() (string, *dbus.Error) {
	fmt.Printf("GetData property accessed: %s\n", h.data)
	return h.data, nil
}

// GetPID returns the PID property
func (h *HelloWorldService) GetPID() (int32, *dbus.Error) {
	fmt.Printf("GetPID property accessed: %d\n", h.pid)
	return h.pid, nil
}

// D-Bus introspection XML
const introspectXML = `
<node>
	<interface name="com.example.HelloWorld">
		<method name="Hello">
			<arg direction="in" name="name" type="s"/>
			<arg direction="out" name="greeting" type="s"/>
		</method>
		<property name="Data" type="s" access="read"/>
		<property name="PID" type="i" access="read"/>
	</interface>` + introspect.IntrospectDataString + `</node> `

func main() {
	fmt.Println("=== D-Bus Hello World Service ===")
	fmt.Printf("Starting D-Bus service: %s\n", serviceName)
	fmt.Printf("Process PID: %d\n", os.Getpid())

	// Connect to session bus
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Printf("Failed to connect to session bus: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create service instance
	service := NewHelloWorldService()

	// Request service name
	reply, err := conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Printf("Failed to request name: %v\n", err)
		os.Exit(1)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Printf("Name already taken: %s\n", serviceName)
		os.Exit(1)
	}

	fmt.Printf("Successfully acquired D-Bus name: %s\n", serviceName)

	// Create property store
	propsSpec := map[string]map[string]*prop.Prop{
		interfaceName: {
			"Data": {
				Value:    service.data,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"PID": {
				Value:    service.pid,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
		},
	}

	// Export properties
	_, err = prop.Export(conn, objectPath, propsSpec)
	if err != nil {
		fmt.Printf("Failed to export properties: %v\n", err)
		os.Exit(1)
	}

	// Export the service object
	err = conn.Export(service, objectPath, interfaceName)
	if err != nil {
		fmt.Printf("Failed to export object: %v\n", err)
		os.Exit(1)
	}

	// Export introspection interface
	err = conn.Export(introspect.Introspectable(introspectXML), objectPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		fmt.Printf("Failed to export introspectable: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Service exported on object path: %s\n", objectPath)
	fmt.Printf("Interface: %s\n", interfaceName)
	fmt.Printf("Available methods: Hello(name string) -> string\n")
	fmt.Printf("Available properties: Data (string), PID (int32)\n")
	fmt.Println()
	fmt.Println("Service is ready! Testing instructions:")
	fmt.Println("1. List services:")
	fmt.Println("   dbus-send --session --print-reply --dest=org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.ListNames")
	fmt.Println()
	fmt.Println("2. Test Hello method:")
	fmt.Printf("   dbus-send --session --print-reply --dest=%s %s %s.Hello string:\"World\"\n", serviceName, objectPath, interfaceName)
	fmt.Println()
	fmt.Println("3. Get Data property:")
	fmt.Printf("   dbus-send --session --print-reply --dest=%s %s org.freedesktop.DBus.Properties.Get string:\"%s\" string:\"Data\"\n", serviceName, objectPath, interfaceName)
	fmt.Println()
	fmt.Println("4. Get PID property:")
	fmt.Printf("   dbus-send --session --print-reply --dest=%s %s org.freedesktop.DBus.Properties.Get string:\"%s\" string:\"PID\"\n", serviceName, objectPath, interfaceName)
	fmt.Println()
	fmt.Println("5. Introspect service:")
	fmt.Printf("   dbus-send --session --print-reply --dest=%s %s org.freedesktop.DBus.Introspectable.Introspect\n", serviceName, objectPath)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the service...")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Keep the service running
	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v\n", sig)
	fmt.Println("Shutting down D-Bus service...")

	// Release the service
	_, err = conn.ReleaseName(serviceName)
	if err != nil {
		fmt.Printf("Failed to release service: %v\n", err)
	}

	fmt.Println("D-Bus Hello World service stopped")
}
