#!/usr/bin/env python3
"""
D-Bus Controller Python Client Example

This script demonstrates how to interact with the D-Bus Controller REST API
using Python requests library.

Requirements:
    pip install requests

Usage:
    python examples/client.py
"""

import json
import requests
from typing import Dict, List, Any

class DBusControllerClient:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
    
    def list_buses(self) -> List[Dict[str, str]]:
        """List available D-Bus buses"""
        response = self.session.get(f"{self.base_url}/buses")
        response.raise_for_status()
        return response.json()
    
    def list_services(self, bus_type: str) -> List[str]:
        """List services on a specific bus"""
        response = self.session.get(f"{self.base_url}/buses/{bus_type}/services")
        response.raise_for_status()
        return response.json()
    
    def get_service_info(self, bus_type: str, service_name: str) -> Dict[str, Any]:
        """Get detailed information about a service"""
        response = self.session.get(f"{self.base_url}/buses/{bus_type}/services/{service_name}")
        response.raise_for_status()
        return response.json()
    
    def list_interfaces(self, bus_type: str, service_name: str) -> List[str]:
        """List interfaces for a service"""
        response = self.session.get(f"{self.base_url}/buses/{bus_type}/services/{service_name}/interfaces")
        response.raise_for_status()
        return response.json()
    
    def get_interface_info(self, bus_type: str, service_name: str, interface_name: str) -> Dict[str, Any]:
        """Get detailed information about an interface"""
        response = self.session.get(f"{self.base_url}/buses/{bus_type}/services/{service_name}/interfaces/{interface_name}")
        response.raise_for_status()
        return response.json()
    
    def call_method(self, bus_type: str, service_name: str, interface_name: str, method_name: str, args: List[Any] = None) -> Dict[str, Any]:
        """Call a D-Bus method"""
        if args is None:
            args = []
        
        payload = {"args": args}
        response = self.session.post(
            f"{self.base_url}/buses/{bus_type}/services/{service_name}/interfaces/{interface_name}/methods/{method_name}/call",
            json=payload
        )
        response.raise_for_status()
        return response.json()
    
    def get_property(self, bus_type: str, service_name: str, interface_name: str, property_name: str) -> Dict[str, Any]:
        """Get a D-Bus property value"""
        response = self.session.get(f"{self.base_url}/buses/{bus_type}/services/{service_name}/interfaces/{interface_name}/properties/{property_name}")
        response.raise_for_status()
        return response.json()
    
    def set_property(self, bus_type: str, service_name: str, interface_name: str, property_name: str, value: Any) -> Dict[str, Any]:
        """Set a D-Bus property value"""
        payload = {"value": value}
        response = self.session.put(
            f"{self.base_url}/buses/{bus_type}/services/{service_name}/interfaces/{interface_name}/properties/{property_name}",
            json=payload
        )
        response.raise_for_status()
        return response.json()
    
    def introspect_service(self, bus_type: str, service_name: str) -> Dict[str, Any]:
        """Get introspection data for a service"""
        response = self.session.get(f"{self.base_url}/buses/{bus_type}/services/{service_name}/introspect")
        response.raise_for_status()
        return response.json()

def main():
    """Example usage of the D-Bus Controller client"""
    client = DBusControllerClient()
    
    try:
        print("=== D-Bus Controller Python Client Example ===")
        print()
        
        # List available buses
        print("1. Available buses:")
        buses = client.list_buses()
        for bus in buses:
            print(f"   - {bus['type']}: {bus['description']}")
        print()
        
        # List services on system bus
        print("2. Services on system bus:")
        services = client.list_services("system")
        print(f"   Found {len(services)} services")
        for service in services[:5]:  # Show first 5 services
            print(f"   - {service}")
        if len(services) > 5:
            print(f"   ... and {len(services) - 5} more")
        print()
        
        # Get info about org.freedesktop.DBus
        print("3. Information about org.freedesktop.DBus:")
        service_info = client.get_service_info("system", "org.freedesktop.DBus")
        print(f"   Name: {service_info['name']}")
        print(f"   Owner: {service_info.get('owner', 'unknown')}")
        print(f"   Interfaces: {len(service_info.get('interfaces', []))}")
        print()
        
        # List interfaces
        print("4. Interfaces for org.freedesktop.DBus:")
        interfaces = client.list_interfaces("system", "org.freedesktop.DBus")
        for interface in interfaces:
            print(f"   - {interface}")
        print()
        
        # Get interface details
        if "org.freedesktop.DBus" in interfaces:
            print("5. Details for org.freedesktop.DBus interface:")
            interface_info = client.get_interface_info("system", "org.freedesktop.DBus", "org.freedesktop.DBus")
            print(f"   Methods: {len(interface_info.get('methods', []))}")
            print(f"   Properties: {len(interface_info.get('properties', []))}")
            print(f"   Signals: {len(interface_info.get('signals', []))}")
            
            # Show method names
            methods = interface_info.get('methods', [])
            if methods:
                print("   Method names:")
                for method in methods[:3]:  # Show first 3 methods
                    print(f"     - {method['name']}")
            print()
        
        # Call Hello method
        print("6. Calling Hello method:")
        try:
            result = client.call_method("system", "org.freedesktop.DBus", "org.freedesktop.DBus", "Hello")
            if result['success']:
                print(f"   ✅ Success: {result.get('return_values', [])}")
            else:
                print(f"   ❌ Error: {result.get('error', 'Unknown error')}")
        except Exception as e:
            print(f"   ❌ Failed to call method: {e}")
        print()
        
        # Get property
        print("7. Getting Features property:")
        try:
            features = client.get_property("system", "org.freedesktop.DBus", "org.freedesktop.DBus", "Features")
            print(f"   Type: {features['type']}")
            print(f"   Value: {features['value']}")
        except Exception as e:
            print(f"   ❌ Failed to get property: {e}")
        print()
        
        # Introspect service
        print("8. Introspecting org.freedesktop.DBus:")
        try:
            introspection = client.introspect_service("system", "org.freedesktop.DBus")
            print(f"   Service: {introspection['service']}")
            print(f"   Object Path: {introspection['object_path']}")
            if introspection.get('parsed_data'):
                parsed = introspection['parsed_data']
                print(f"   Parsed Interfaces: {len(parsed.get('interfaces', []))}")
                print(f"   Parsed Nodes: {len(parsed.get('nodes', []))}")
        except Exception as e:
            print(f"   ❌ Failed to introspect: {e}")
        print()
        
        print("=== Example completed successfully ===")
        
    except requests.exceptions.ConnectionError:
        print("❌ Could not connect to D-Bus Controller API")
        print("Make sure the server is running on http://localhost:8080")
        print("Start it with: make run")
    except Exception as e:
        print(f"❌ An error occurred: {e}")

if __name__ == "__main__":
    main()