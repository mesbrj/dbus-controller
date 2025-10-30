package model

import "time"

// BusInfo represents information about a D-Bus
type BusInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ServiceInfo represents information about a D-Bus service
type ServiceInfo struct {
	Name          string               `json:"name"`
	Owner         string               `json:"owner,omitempty"`
	Interfaces    []string             `json:"interfaces,omitempty"`
	ObjectPaths   []string             `json:"object_paths,omitempty"`
	Introspection *IntrospectionResult `json:"introspection,omitempty"`
}

// InterfaceInfo represents information about a D-Bus interface
type InterfaceInfo struct {
	Name       string         `json:"name"`
	Methods    []MethodInfo   `json:"methods,omitempty"`
	Properties []PropertyInfo `json:"properties,omitempty"`
	Signals    []SignalInfo   `json:"signals,omitempty"`
}

// MethodInfo represents information about a D-Bus method
type MethodInfo struct {
	Name        string            `json:"name"`
	InArgs      []ArgumentInfo    `json:"in_args,omitempty"`
	OutArgs     []ArgumentInfo    `json:"out_args,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ArgumentInfo represents information about a method argument
type ArgumentInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Direction string `json:"direction"` // "in" or "out"
}

// PropertyInfo represents information about a D-Bus property
type PropertyInfo struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Access      string            `json:"access"` // "read", "write", "readwrite"
	Value       interface{}       `json:"value,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// SignalInfo represents information about a D-Bus signal
type SignalInfo struct {
	Name        string            `json:"name"`
	Args        []ArgumentInfo    `json:"args,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// MethodCallResult represents the result of a D-Bus method call
type MethodCallResult struct {
	Success      bool          `json:"success"`
	ReturnValues []interface{} `json:"return_values,omitempty"`
	Error        string        `json:"error,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

// PropertyValue represents a D-Bus property value
type PropertyValue struct {
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
}

// SignalSubscription represents a D-Bus signal subscription
type SignalSubscription struct {
	ID        string    `json:"id"`
	BusType   string    `json:"bus_type"`
	Service   string    `json:"service"`
	Interface string    `json:"interface"`
	Signal    string    `json:"signal"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// IntrospectionResult represents the result of D-Bus introspection
type IntrospectionResult struct {
	Service    string               `json:"service"`
	ObjectPath string               `json:"object_path"`
	XML        string               `json:"xml"`
	ParsedData *ParsedIntrospection `json:"parsed_data,omitempty"`
	Timestamp  time.Time            `json:"timestamp"`
}

// ParsedIntrospection represents parsed introspection data
type ParsedIntrospection struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
	Nodes      []NodeInfo      `json:"nodes,omitempty"`
}

// NodeInfo represents a D-Bus node
type NodeInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
