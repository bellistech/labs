# Python Learning Path: Beginner to Production Web Service

A structured progression from absolute beginner to building a PostgreSQL-backed web service running as a systemd service, with GitHub integration, Jenkins CI/CD, and professional scaffolding.

**Similar structure to the Rust and Go projects we've discussed - progressive exercises with hands-on labs and GitHub portfolio building.**

---

## Table of Contents

1. [Phase 1: Python Fundamentals (Days 1-3)](#phase-1-python-fundamentals)
2. [Phase 2: Intermediate Python (Days 4-6)](#phase-2-intermediate-python)
3. [Phase 3: Working with Data & Files (Days 7-9)](#phase-3-working-with-data--files)
4. [Phase 4: Database Fundamentals (Days 10-12)](#phase-4-database-fundamentals)
5. [Phase 5: Web Development with FastAPI (Days 13-16)](#phase-5-web-development-with-fastapi)
6. [Phase 6: Production Deployment (Days 17-20)](#phase-6-production-deployment)
7. [Git & GitHub Workflow](#git--github-workflow)
8. [Jenkins CI/CD Pipeline](#jenkins-cicd-pipeline)
9. [Project Scaffolding Templates](#project-scaffolding-templates)
10. [Capstone Project: Full-Stack Service](#capstone-project)

---

## Initial Setup

### Prerequisites

```bash
# macOS
brew install python3 git postgresql

# Ubuntu/Debian
sudo apt update
sudo apt install python3 python3-pip python3-venv git postgresql postgresql-contrib

# Verify installations
python3 --version   # Should be 3.10+
git --version
psql --version
```

### Project Directory Structure

```bash
# Create your Python learning workspace
mkdir -p ~/dev/python
cd ~/dev/python

# Initialize git repo for the entire learning path
git init python-learning
cd python-learning

# Create basic structure
mkdir -p {day_01,day_02,day_03,day_04,day_05,day_06}
mkdir -p {day_07,day_08,day_09,day_10,day_11,day_12}
mkdir -p {day_13_16_fastapi,day_17_20_deployment}
mkdir -p {tests,docs}

# Initial .gitignore
cat > .gitignore << 'EOF'
# Python
__pycache__/
*.py[cod]
*$py.class
venv/
.venv/
env/
*.egg-info/

# IDE
.idea/
.vscode/
*.swp

# Testing
.pytest_cache/
.coverage
htmlcov/

# Environment
.env
*.env.local

# Logs
*.log

# Database
*.db
*.sqlite3

# OS
.DS_Store
EOF

git add .gitignore
git commit -m "Initial commit: project structure"
```

---

## Phase 1: Python Fundamentals

### Day 1: Environment Setup & Hello World

#### Exercise 1.1: Environment Setup

```bash
# Create Day 1 directory and virtual environment
cd ~/dev/python/python-learning/day_01
python3 -m venv venv
source venv/bin/activate   # On Windows: venv\Scripts\activate

# Verify
which python    # Should show venv path
python --version
pip --version
```

#### Exercise 1.2: Hello World Variations

Create `day_01/hello.py`:

```python
#!/usr/bin/env python3
"""
Day 1 Exercise: Hello World variations
Learning: print(), variables, f-strings, input()
"""

# Basic print
print("Hello, World!")

# Variables and f-strings
name = "Muck"
print(f"Hello, {name}!")

# User input
user_name = input("What is your name? ")
print(f"Welcome, {user_name}!")

# Multiple variables
greeting = "Hello"
target = "Python"
print(f"{greeting}, {target}!")

# String methods
message = "  hello world  "
print(f"Original: '{message}'")
print(f"Stripped: '{message.strip()}'")
print(f"Upper: '{message.upper()}'")
print(f"Title: '{message.strip().title()}'")
```

Run it:
```bash
chmod +x hello.py
./hello.py
# or
python hello.py
```

#### Exercise 1.3: Basic Calculator

Create `day_01/calculator.py`:

```python
#!/usr/bin/env python3
"""
Day 1 Exercise: Basic calculator
Learning: arithmetic operators, type conversion, error handling basics
"""

def main():
    print("=== Simple Calculator ===")
    
    # Get user input and convert to numbers
    try:
        num1 = float(input("Enter first number: "))
        num2 = float(input("Enter second number: "))
    except ValueError:
        print("Error: Please enter valid numbers")
        return
    
    # Perform operations
    addition = num1 + num2
    subtraction = num1 - num2
    multiplication = num1 * num2
    
    # Handle division by zero
    if num2 != 0:
        division = num1 / num2
        floor_div = num1 // num2
        modulo = num1 % num2
    else:
        division = "undefined (division by zero)"
        floor_div = "undefined"
        modulo = "undefined"
    
    power = num1 ** num2
    
    # Display results
    print(f"\nResults:")
    print(f"{num1} + {num2} = {addition}")
    print(f"{num1} - {num2} = {subtraction}")
    print(f"{num1} * {num2} = {multiplication}")
    print(f"{num1} / {num2} = {division}")
    print(f"{num1} // {num2} = {floor_div}")
    print(f"{num1} % {num2} = {modulo}")
    print(f"{num1} ** {num2} = {power}")

if __name__ == "__main__":
    main()
```

#### Exercise 1.4: Data Types Explorer

Create `day_01/data_types.py`:

```python
#!/usr/bin/env python3
"""
Day 1 Exercise: Exploring Python data types
Learning: int, float, str, bool, type(), isinstance()
"""

def explore_types():
    # Integer
    age = 28
    print(f"age = {age}, type: {type(age).__name__}")
    
    # Float
    temperature = 98.6
    print(f"temperature = {temperature}, type: {type(temperature).__name__}")
    
    # String
    hostname = "router-01"
    print(f"hostname = '{hostname}', type: {type(hostname).__name__}")
    
    # Boolean
    is_active = True
    print(f"is_active = {is_active}, type: {type(is_active).__name__}")
    
    # None
    result = None
    print(f"result = {result}, type: {type(result).__name__}")
    
    # Type conversion
    print("\n=== Type Conversion ===")
    port_str = "8080"
    port_int = int(port_str)
    print(f"'{port_str}' (str) -> {port_int} (int)")
    
    count = 42
    count_str = str(count)
    print(f"{count} (int) -> '{count_str}' (str)")
    
    # Boolean conversion
    print("\n=== Boolean Truthiness ===")
    values = [0, 1, "", "hello", [], [1,2,3], None, {}, {"a": 1}]
    for val in values:
        print(f"bool({val!r}) = {bool(val)}")

def main():
    explore_types()

if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add day_01/
git commit -m "Day 1: Hello world, calculator, and data types"
```

---

### Day 2: Control Flow & Data Structures

#### Exercise 2.1: Conditionals

Create `day_02/conditions.py`:

```python
#!/usr/bin/env python3
"""
Day 2 Exercise: Control flow with conditionals
Learning: if/elif/else, comparison operators, logical operators
"""

def check_network_port(port: int) -> str:
    """Classify a network port number."""
    if port < 0 or port > 65535:
        return "Invalid port number"
    elif port == 0:
        return "Reserved port"
    elif port < 1024:
        return "Well-known/privileged port (requires root)"
    elif port < 49152:
        return "Registered port"
    else:
        return "Dynamic/ephemeral port"

def check_http_status(code: int) -> str:
    """Interpret HTTP status codes."""
    if 100 <= code < 200:
        return "Informational"
    elif 200 <= code < 300:
        return "Success"
    elif 300 <= code < 400:
        return "Redirection"
    elif 400 <= code < 500:
        return "Client Error"
    elif 500 <= code < 600:
        return "Server Error"
    else:
        return "Unknown status code"

def check_ip_class(ip: str) -> str:
    """Determine IP address class (simplified)."""
    try:
        first_octet = int(ip.split(".")[0])
    except (ValueError, IndexError):
        return "Invalid IP"
    
    if 1 <= first_octet <= 126:
        return "Class A"
    elif 128 <= first_octet <= 191:
        return "Class B"
    elif 192 <= first_octet <= 223:
        return "Class C"
    elif 224 <= first_octet <= 239:
        return "Class D (Multicast)"
    elif 240 <= first_octet <= 255:
        return "Class E (Reserved)"
    else:
        return "Special"

def main():
    # Test port classification
    test_ports = [22, 80, 443, 3000, 8080, 49152, 65535, 70000]
    print("=== Port Classification ===")
    for port in test_ports:
        result = check_network_port(port)
        print(f"Port {port}: {result}")
    
    print("\n=== HTTP Status Codes ===")
    test_codes = [100, 200, 301, 404, 500, 503]
    for code in test_codes:
        result = check_http_status(code)
        print(f"HTTP {code}: {result}")
    
    print("\n=== IP Address Classes ===")
    test_ips = ["10.0.0.1", "172.16.0.1", "192.168.1.1", "224.0.0.1"]
    for ip in test_ips:
        result = check_ip_class(ip)
        print(f"{ip}: {result}")

if __name__ == "__main__":
    main()
```

#### Exercise 2.2: Loops

Create `day_02/loops.py`:

```python
#!/usr/bin/env python3
"""
Day 2 Exercise: Loops
Learning: for loops, while loops, range(), enumerate(), break/continue
"""

def count_down(start: int) -> None:
    """Countdown timer using while loop."""
    print(f"Countdown from {start}:")
    while start > 0:
        print(f"  {start}...")
        start -= 1
    print("  Liftoff!")

def scan_ports(start: int, end: int, skip_ports: list) -> list:
    """Simulate port scanning with continue."""
    open_ports = []
    for port in range(start, end + 1):
        if port in skip_ports:
            continue  # Skip certain ports
        # Simulate: even ports are "open"
        if port % 2 == 0:
            open_ports.append(port)
        if len(open_ports) >= 5:
            break  # Stop after finding 5 open ports
    return open_ports

def enumerate_interfaces() -> None:
    """Using enumerate for indexed iteration."""
    interfaces = ["lo", "eth0", "eth1", "docker0", "br-lan"]
    print("\nNetwork interfaces:")
    for idx, iface in enumerate(interfaces, start=1):
        print(f"  {idx}. {iface}")

def nested_loops_example() -> None:
    """Subnet scanning with nested loops."""
    print("\nScanning subnets:")
    subnets = ["10.0.1", "10.0.2"]
    hosts = [1, 2, 3]
    
    for subnet in subnets:
        print(f"  Subnet {subnet}.0/24:")
        for host in hosts:
            ip = f"{subnet}.{host}"
            print(f"    Scanning {ip}...")

def main():
    count_down(5)
    
    enumerate_interfaces()
    
    skip = [21, 23, 25]  # Skip FTP, Telnet, SMTP
    found = scan_ports(20, 100, skip)
    print(f"\nFirst 5 open ports (skipping {skip}): {found}")
    
    nested_loops_example()

if __name__ == "__main__":
    main()
```

#### Exercise 2.3: Data Structures

Create `day_02/data_structures.py`:

```python
#!/usr/bin/env python3
"""
Day 2 Exercise: Lists, Dictionaries, Sets, Tuples
Learning: core data structures and their operations
"""

def list_operations():
    """Working with lists."""
    print("=== Lists ===")
    
    # Create and modify
    servers = ["web01", "web02", "db01"]
    servers.append("cache01")
    servers.insert(0, "lb01")
    
    print(f"Servers: {servers}")
    print(f"First: {servers[0]}, Last: {servers[-1]}")
    print(f"Web servers: {servers[1:3]}")  # Slicing
    
    # List comprehension
    web_servers = [s for s in servers if s.startswith("web")]
    print(f"Filtered web servers: {web_servers}")
    
    # Sorting
    servers_sorted = sorted(servers)
    print(f"Sorted: {servers_sorted}")
    
    # Common operations
    print(f"Length: {len(servers)}")
    print(f"'db01' in servers: {'db01' in servers}")

def dict_operations():
    """Working with dictionaries."""
    print("\n=== Dictionaries ===")
    
    # Server inventory
    inventory = {
        "web01": {"ip": "10.0.1.10", "role": "web", "status": "up"},
        "web02": {"ip": "10.0.1.11", "role": "web", "status": "up"},
        "db01": {"ip": "10.0.2.10", "role": "database", "status": "up"},
    }
    
    # Access and update
    print(f"web01 IP: {inventory['web01']['ip']}")
    inventory["web01"]["status"] = "maintenance"
    
    # Safe access with .get()
    unknown = inventory.get("web99", {"ip": "unknown"})
    print(f"web99 IP: {unknown['ip']}")
    
    # Iterate
    print("\nServer Status:")
    for name, info in inventory.items():
        print(f"  {name}: {info['status']}")
    
    # Dict comprehension
    ips = {name: info["ip"] for name, info in inventory.items()}
    print(f"IP map: {ips}")
    
    # Keys and values
    print(f"All servers: {list(inventory.keys())}")

def set_operations():
    """Working with sets."""
    print("\n=== Sets ===")
    
    # Network ACL example
    allowed_ports = {22, 80, 443, 8080}
    requested_ports = {22, 80, 3306, 5432}
    
    print(f"Allowed: {allowed_ports}")
    print(f"Requested: {requested_ports}")
    print(f"Permitted (intersection): {allowed_ports & requested_ports}")
    print(f"Denied (difference): {requested_ports - allowed_ports}")
    print(f"All mentioned (union): {allowed_ports | requested_ports}")
    
    # Set operations
    allowed_ports.add(3306)
    print(f"After adding 3306: {allowed_ports}")

def tuple_operations():
    """Working with tuples (immutable)."""
    print("\n=== Tuples ===")
    
    # Coordinates, config values (immutable)
    server_location = ("us-east-1", "az-a", "rack-42")
    region, az, rack = server_location  # Unpacking
    
    print(f"Location: {server_location}")
    print(f"Region: {region}, AZ: {az}, Rack: {rack}")
    
    # Named tuple for cleaner code
    from collections import namedtuple
    Server = namedtuple("Server", ["name", "ip", "port"])
    
    web = Server("web01", "10.0.1.10", 80)
    print(f"Server: {web.name} at {web.ip}:{web.port}")
    
    # Tuple as dict key (immutable = hashable)
    connections = {
        ("10.0.1.10", 80): "active",
        ("10.0.1.10", 443): "active",
    }
    print(f"Connection states: {connections}")

def main():
    list_operations()
    dict_operations()
    set_operations()
    tuple_operations()

if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add day_02/
git commit -m "Day 2: Control flow and data structures"
```

---

### Day 3: Functions & Modules

#### Exercise 3.1: Functions Deep Dive

Create `day_03/functions.py`:

```python
#!/usr/bin/env python3
"""
Day 3 Exercise: Functions
Learning: def, args, kwargs, return values, docstrings, type hints
"""

from typing import Optional, List, Dict, Tuple, Union


def ping(host: str, count: int = 4, timeout: float = 1.0) -> dict:
    """
    Simulate a ping command.
    
    Args:
        host: Target hostname or IP address
        count: Number of ping packets (default: 4)
        timeout: Timeout in seconds (default: 1.0)
    
    Returns:
        dict with success status and stats
    """
    import random
    latencies = [random.uniform(1, 100) for _ in range(count)]
    
    return {
        "host": host,
        "packets_sent": count,
        "packets_received": count,
        "packet_loss": 0.0,
        "min_latency": round(min(latencies), 2),
        "max_latency": round(max(latencies), 2),
        "avg_latency": round(sum(latencies) / len(latencies), 2),
    }


def parse_cidr(cidr: str) -> Tuple[str, int]:
    """
    Parse CIDR notation into IP and prefix length.
    
    Args:
        cidr: CIDR string like "192.168.1.0/24"
    
    Returns:
        Tuple of (ip_address, prefix_length)
    
    Raises:
        ValueError: If CIDR format is invalid
    """
    if "/" not in cidr:
        raise ValueError(f"Invalid CIDR format: {cidr}")
    
    ip, prefix = cidr.split("/")
    prefix_int = int(prefix)
    
    if not 0 <= prefix_int <= 32:
        raise ValueError(f"Invalid prefix length: {prefix}")
    
    return ip, prefix_int


def configure_interface(*args, **kwargs) -> None:
    """
    Demonstrate *args and **kwargs.
    
    Args:
        *args: Positional arguments (interface names)
        **kwargs: Keyword arguments (configuration options)
    """
    print(f"Configuring interfaces: {args}")
    print("Options:")
    for key, value in kwargs.items():
        print(f"  {key}: {value}")


def create_vlan(
    vlan_id: int,
    name: str,
    *ports: str,
    tagged: bool = False,
    description: Optional[str] = None
) -> Dict:
    """
    Create a VLAN configuration.
    
    Demonstrates combining regular args, *args, and kwargs.
    """
    return {
        "vlan_id": vlan_id,
        "name": name,
        "ports": list(ports),
        "tagged": tagged,
        "description": description or f"VLAN {vlan_id}",
    }


def calculate_subnet_info(cidr: str) -> Dict[str, Union[str, int]]:
    """Calculate subnet information from CIDR."""
    ip, prefix = parse_cidr(cidr)
    
    # Calculate number of hosts
    host_bits = 32 - prefix
    total_hosts = 2 ** host_bits
    usable_hosts = max(0, total_hosts - 2)  # Minus network and broadcast
    
    return {
        "network": ip,
        "prefix": prefix,
        "total_addresses": total_hosts,
        "usable_hosts": usable_hosts,
        "cidr": cidr,
    }


# Lambda functions
get_ip = lambda device: device.get("ip", "N/A")
is_up = lambda device: device.get("status") == "up"


def main():
    # Basic function calls
    print("=== Basic Function Calls ===")
    result = ping("8.8.8.8")
    print(f"Ping result: {result['avg_latency']}ms avg")
    
    # With different arguments
    result = ping("10.0.0.1", count=10, timeout=0.5)
    print(f"Ping 10.0.0.1: {result['packets_received']}/{result['packets_sent']} received")
    
    # Tuple unpacking from return
    print("\n=== CIDR Parsing ===")
    ip, prefix = parse_cidr("192.168.1.0/24")
    print(f"Network: {ip}, Prefix: /{prefix}")
    
    # Error handling
    try:
        parse_cidr("invalid")
    except ValueError as e:
        print(f"Expected error: {e}")
    
    # *args and **kwargs
    print("\n=== *args and **kwargs ===")
    configure_interface("eth0", "eth1", mtu=9000, speed="1000")
    
    # Mixed arguments
    print("\n=== VLAN Creation ===")
    vlan = create_vlan(100, "Management", "eth0", "eth1", tagged=True)
    print(f"VLAN config: {vlan}")
    
    # Subnet calculation
    print("\n=== Subnet Calculation ===")
    for cidr in ["10.0.0.0/8", "172.16.0.0/16", "192.168.1.0/24", "10.0.0.0/30"]:
        info = calculate_subnet_info(cidr)
        print(f"{cidr}: {info['usable_hosts']} usable hosts")
    
    # Lambda examples
    print("\n=== Lambda Functions ===")
    devices = [
        {"name": "router1", "ip": "10.0.0.1", "status": "up"},
        {"name": "switch1", "ip": "10.0.0.2", "status": "down"},
    ]
    print(f"IPs: {[get_ip(d) for d in devices]}")
    print(f"Up devices: {[d['name'] for d in devices if is_up(d)]}")


if __name__ == "__main__":
    main()
```

#### Exercise 3.2: Creating Modules

Create module structure:

```bash
mkdir -p day_03/netutils
touch day_03/netutils/__init__.py
```

Create `day_03/netutils/ip.py`:

```python
"""IP address utilities module."""

import ipaddress
from typing import List, Optional


def is_valid_ipv4(ip: str) -> bool:
    """Check if string is a valid IPv4 address."""
    try:
        ipaddress.IPv4Address(ip)
        return True
    except ipaddress.AddressValueError:
        return False


def is_valid_ipv6(ip: str) -> bool:
    """Check if string is a valid IPv6 address."""
    try:
        ipaddress.IPv6Address(ip)
        return True
    except ipaddress.AddressValueError:
        return False


def get_network_hosts(cidr: str, max_hosts: int = 10) -> List[str]:
    """
    Get list of host IPs in a network.
    
    Args:
        cidr: Network in CIDR notation
        max_hosts: Maximum hosts to return
    
    Returns:
        List of IP addresses as strings
    """
    network = ipaddress.ip_network(cidr, strict=False)
    hosts = []
    for idx, host in enumerate(network.hosts()):
        if idx >= max_hosts:
            break
        hosts.append(str(host))
    return hosts


def is_private(ip: str) -> bool:
    """Check if IP address is in private range."""
    try:
        addr = ipaddress.ip_address(ip)
        return addr.is_private
    except ValueError:
        return False


def is_in_network(ip: str, network: str) -> bool:
    """Check if IP is in the given network."""
    try:
        addr = ipaddress.ip_address(ip)
        net = ipaddress.ip_network(network, strict=False)
        return addr in net
    except ValueError:
        return False


def summarize_networks(networks: List[str]) -> List[str]:
    """Collapse a list of networks into summary routes."""
    network_objs = [ipaddress.ip_network(n, strict=False) for n in networks]
    collapsed = ipaddress.collapse_addresses(network_objs)
    return [str(n) for n in collapsed]
```

Create `day_03/netutils/__init__.py`:

```python
"""
Network utilities package.

Provides IP address utilities for network automation.
"""

from .ip import (
    is_valid_ipv4,
    is_valid_ipv6,
    get_network_hosts,
    is_private,
    is_in_network,
    summarize_networks,
)

__version__ = "0.1.0"
__all__ = [
    "is_valid_ipv4",
    "is_valid_ipv6",
    "get_network_hosts",
    "is_private",
    "is_in_network",
    "summarize_networks",
]
```

Create `day_03/test_netutils.py`:

```python
#!/usr/bin/env python3
"""Test our netutils module."""

# Import from our package
from netutils import (
    is_valid_ipv4,
    is_valid_ipv6,
    get_network_hosts,
    is_private,
    is_in_network,
    summarize_networks,
)


def main():
    print("=== Testing netutils module ===\n")
    
    # IP validation
    test_ips = ["192.168.1.1", "10.0.0.1", "256.1.1.1", "2001:db8::1", "invalid"]
    print("IP Validation:")
    for addr in test_ips:
        v4 = is_valid_ipv4(addr)
        v6 = is_valid_ipv6(addr)
        print(f"  {addr}: IPv4={v4}, IPv6={v6}")
    
    # Private IP check
    print("\nPrivate IP Check:")
    for addr in ["10.0.0.1", "8.8.8.8", "192.168.1.1", "172.16.0.1", "1.1.1.1"]:
        private = is_private(addr)
        print(f"  {addr}: private={private}")
    
    # Network membership
    print("\nNetwork Membership (10.0.0.0/24):")
    network = "10.0.0.0/24"
    for ip in ["10.0.0.1", "10.0.0.254", "10.0.1.1", "192.168.1.1"]:
        in_net = is_in_network(ip, network)
        print(f"  {ip} in {network}: {in_net}")
    
    # Network hosts
    print("\nNetwork Hosts (192.168.1.0/24):")
    hosts = get_network_hosts("192.168.1.0/24", max_hosts=5)
    for host in hosts:
        print(f"  {host}")
    
    # Network summarization
    print("\nNetwork Summarization:")
    networks = [
        "192.168.0.0/24",
        "192.168.1.0/24",
        "192.168.2.0/24",
        "192.168.3.0/24",
    ]
    summarized = summarize_networks(networks)
    print(f"  Input: {networks}")
    print(f"  Summarized: {summarized}")


if __name__ == "__main__":
    main()
```

Run from day_03 directory:
```bash
cd day_03
python test_netutils.py
```

**Git checkpoint:**
```bash
git add day_03/
git commit -m "Day 3: Functions and custom modules"
```

---

## Continue to Part 2...

The complete learning path continues with:

- **Phase 2 (Days 4-6)**: OOP, Error Handling, Decorators, Generators
- **Phase 3 (Days 7-9)**: File I/O, HTTP Requests, System Commands
- **Phase 4 (Days 10-12)**: SQLite and PostgreSQL
- **Phase 5 (Days 13-16)**: FastAPI Web Development
- **Phase 6 (Days 17-20)**: Production Deployment with Systemd
- **Jenkins CI/CD Pipeline**
- **Complete Project Scaffolding**

See the full document for all exercises and the capstone project.

## Phase 4 Continued: PostgreSQL Integration

### Day 11-12: PostgreSQL with Python

#### Exercise 11.1: PostgreSQL Setup and Connection

Create `postgres_client.py`:

```python
#!/usr/bin/env python3
"""
Day 11-12 Exercise: PostgreSQL Integration
Learning: psycopg2/psycopg3, connection pooling, transactions

Install: pip install psycopg2-binary
"""

import os
from typing import List, Dict, Optional, Any
from contextlib import contextmanager
from dataclasses import dataclass
from datetime import datetime


# Configuration from environment
DB_CONFIG = {
    "host": os.getenv("POSTGRES_HOST", "localhost"),
    "port": int(os.getenv("POSTGRES_PORT", 5432)),
    "database": os.getenv("POSTGRES_DB", "network_inventory"),
    "user": os.getenv("POSTGRES_USER", "netadmin"),
    "password": os.getenv("POSTGRES_PASSWORD", "secretpassword"),
}


@dataclass
class Device:
    """Device data model."""
    id: Optional[int] = None
    hostname: str = ""
    ip_address: str = ""
    device_type: str = "unknown"
    vendor: Optional[str] = None
    status: str = "unknown"
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class PostgresClient:
    """
    PostgreSQL database client with connection pooling.
    """
    
    def __init__(self, config: Dict = None):
        """Initialize with configuration."""
        try:
            import psycopg2
            from psycopg2 import pool
            self.psycopg2 = psycopg2
        except ImportError:
            print("psycopg2 not installed. Run: pip install psycopg2-binary")
            raise
        
        self.config = config or DB_CONFIG
        self._pool = None
    
    def connect(self, min_conn: int = 1, max_conn: int = 10) -> None:
        """Create connection pool."""
        from psycopg2 import pool
        
        self._pool = pool.ThreadedConnectionPool(
            min_conn,
            max_conn,
            **self.config
        )
        print(f"Connected to PostgreSQL at {self.config['host']}:{self.config['port']}")
    
    def disconnect(self) -> None:
        """Close all connections in the pool."""
        if self._pool:
            self._pool.closeall()
            print("Disconnected from PostgreSQL")
    
    @contextmanager
    def get_cursor(self, commit: bool = True):
        """Context manager for database operations."""
        conn = self._pool.getconn()
        try:
            cursor = conn.cursor()
            yield cursor
            if commit:
                conn.commit()
        except Exception as e:
            conn.rollback()
            raise
        finally:
            cursor.close()
            self._pool.putconn(conn)
    
    @contextmanager
    def transaction(self):
        """Context manager for explicit transactions."""
        conn = self._pool.getconn()
        try:
            cursor = conn.cursor()
            yield cursor
            conn.commit()
        except Exception as e:
            conn.rollback()
            raise
        finally:
            cursor.close()
            self._pool.putconn(conn)
    
    def init_schema(self) -> None:
        """Initialize database schema."""
        with self.get_cursor() as cursor:
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS devices (
                    id SERIAL PRIMARY KEY,
                    hostname VARCHAR(255) NOT NULL UNIQUE,
                    ip_address INET NOT NULL,
                    device_type VARCHAR(50) DEFAULT 'unknown',
                    vendor VARCHAR(100),
                    status VARCHAR(50) DEFAULT 'unknown',
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            """)
            
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS interfaces (
                    id SERIAL PRIMARY KEY,
                    device_id INTEGER NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
                    name VARCHAR(100) NOT NULL,
                    ip_address INET,
                    mac_address MACADDR,
                    status VARCHAR(50) DEFAULT 'down',
                    speed_mbps INTEGER,
                    mtu INTEGER DEFAULT 1500,
                    UNIQUE (device_id, name)
                )
            """)
            
            cursor.execute("""
                CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status)
            """)
            
            cursor.execute("""
                CREATE INDEX IF NOT EXISTS idx_devices_vendor ON devices(vendor)
            """)
            
            # Create updated_at trigger
            cursor.execute("""
                CREATE OR REPLACE FUNCTION update_updated_at()
                RETURNS TRIGGER AS $$
                BEGIN
                    NEW.updated_at = CURRENT_TIMESTAMP;
                    RETURN NEW;
                END;
                $$ LANGUAGE plpgsql
            """)
            
            cursor.execute("""
                DROP TRIGGER IF EXISTS devices_updated_at ON devices
            """)
            
            cursor.execute("""
                CREATE TRIGGER devices_updated_at
                BEFORE UPDATE ON devices
                FOR EACH ROW
                EXECUTE FUNCTION update_updated_at()
            """)
        
        print("Schema initialized")
    
    # CRUD Operations
    def create_device(self, device: Device) -> int:
        """Create a new device."""
        with self.get_cursor() as cursor:
            cursor.execute("""
                INSERT INTO devices (hostname, ip_address, device_type, vendor, status)
                VALUES (%s, %s, %s, %s, %s)
                RETURNING id
            """, (device.hostname, device.ip_address, device.device_type,
                  device.vendor, device.status))
            return cursor.fetchone()[0]
    
    def get_device(self, device_id: int) -> Optional[Device]:
        """Get device by ID."""
        with self.get_cursor(commit=False) as cursor:
            cursor.execute("""
                SELECT id, hostname, ip_address, device_type, vendor, status,
                       created_at, updated_at
                FROM devices WHERE id = %s
            """, (device_id,))
            row = cursor.fetchone()
            if row:
                return Device(
                    id=row[0],
                    hostname=row[1],
                    ip_address=str(row[2]),
                    device_type=row[3],
                    vendor=row[4],
                    status=row[5],
                    created_at=row[6],
                    updated_at=row[7],
                )
            return None
    
    def get_all_devices(self, status: str = None, vendor: str = None,
                        limit: int = 100, offset: int = 0) -> List[Device]:
        """Get devices with optional filtering."""
        with self.get_cursor(commit=False) as cursor:
            query = "SELECT * FROM devices WHERE 1=1"
            params = []
            
            if status:
                query += " AND status = %s"
                params.append(status)
            if vendor:
                query += " AND vendor = %s"
                params.append(vendor)
            
            query += " ORDER BY hostname LIMIT %s OFFSET %s"
            params.extend([limit, offset])
            
            cursor.execute(query, params)
            devices = []
            for row in cursor.fetchall():
                devices.append(Device(
                    id=row[0],
                    hostname=row[1],
                    ip_address=str(row[2]),
                    device_type=row[3],
                    vendor=row[4],
                    status=row[5],
                    created_at=row[6],
                    updated_at=row[7],
                ))
            return devices
    
    def update_device(self, device_id: int, **kwargs) -> bool:
        """Update device fields."""
        if not kwargs:
            return False
        
        set_clauses = []
        values = []
        for key, value in kwargs.items():
            set_clauses.append(f"{key} = %s")
            values.append(value)
        values.append(device_id)
        
        with self.get_cursor() as cursor:
            cursor.execute(f"""
                UPDATE devices
                SET {', '.join(set_clauses)}
                WHERE id = %s
            """, values)
            return cursor.rowcount > 0
    
    def delete_device(self, device_id: int) -> bool:
        """Delete a device."""
        with self.get_cursor() as cursor:
            cursor.execute("DELETE FROM devices WHERE id = %s", (device_id,))
            return cursor.rowcount > 0
    
    def bulk_insert_devices(self, devices: List[Device]) -> int:
        """Bulk insert devices efficiently."""
        from psycopg2.extras import execute_values
        
        with self.get_cursor() as cursor:
            data = [
                (d.hostname, d.ip_address, d.device_type, d.vendor, d.status)
                for d in devices
            ]
            execute_values(
                cursor,
                """
                INSERT INTO devices (hostname, ip_address, device_type, vendor, status)
                VALUES %s
                ON CONFLICT (hostname) DO NOTHING
                """,
                data
            )
            return cursor.rowcount
    
    def search_devices(self, query: str) -> List[Device]:
        """Search devices by hostname or IP."""
        with self.get_cursor(commit=False) as cursor:
            cursor.execute("""
                SELECT * FROM devices
                WHERE hostname ILIKE %s
                   OR ip_address::text LIKE %s
                ORDER BY hostname
            """, (f"%{query}%", f"%{query}%"))
            
            return [
                Device(
                    id=row[0], hostname=row[1], ip_address=str(row[2]),
                    device_type=row[3], vendor=row[4], status=row[5],
                    created_at=row[6], updated_at=row[7]
                )
                for row in cursor.fetchall()
            ]


def main():
    """Demo PostgreSQL operations."""
    print("=== PostgreSQL Client Demo ===\n")
    print("Note: This requires a running PostgreSQL server.")
    print("Set environment variables: POSTGRES_HOST, POSTGRES_DB, etc.\n")
    
    # Show what the code would do
    print("Example usage:")
    print("""
    db = PostgresClient()
    db.connect()
    db.init_schema()
    
    # Create device
    device = Device(
        hostname="router-01",
        ip_address="10.0.0.1",
        device_type="router",
        vendor="juniper",
        status="up"
    )
    device_id = db.create_device(device)
    
    # Read device
    device = db.get_device(device_id)
    print(f"Device: {device.hostname}")
    
    # Update device
    db.update_device(device_id, status="maintenance")
    
    # Search
    results = db.search_devices("router")
    
    # Bulk insert
    devices = [
        Device(hostname=f"switch-{i:02d}", ip_address=f"10.0.1.{i}",
               device_type="switch", vendor="arista")
        for i in range(1, 11)
    ]
    count = db.bulk_insert_devices(devices)
    
    db.disconnect()
    """)


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add postgres_client.py
git commit -m "Day 11-12: PostgreSQL client with connection pooling"
```

---

## Phase 5: Web Development with FastAPI

### Day 13-14: FastAPI Fundamentals

#### Exercise 13.1: Basic FastAPI Application

Create the project structure:

```bash
mkdir -p network_api/{app,tests}
cd network_api
python -m venv venv
source venv/bin/activate
pip install fastapi uvicorn pydantic python-dotenv
```

Create `network_api/app/main.py`:

```python
#!/usr/bin/env python3
"""
Day 13-14 Exercise: FastAPI Web Application
Learning: REST APIs, Pydantic models, dependency injection
"""

from fastapi import FastAPI, HTTPException, Depends, Query, status
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field, IPvAnyAddress
from typing import List, Optional
from datetime import datetime
from enum import Enum
import uvicorn


# Pydantic Models (Request/Response schemas)
class DeviceStatus(str, Enum):
    UP = "up"
    DOWN = "down"
    MAINTENANCE = "maintenance"
    UNKNOWN = "unknown"


class DeviceType(str, Enum):
    ROUTER = "router"
    SWITCH = "switch"
    FIREWALL = "firewall"
    SERVER = "server"
    OTHER = "other"


class DeviceBase(BaseModel):
    """Base device model with common fields."""
    hostname: str = Field(..., min_length=1, max_length=255, 
                          description="Device hostname")
    ip_address: str = Field(..., description="Device IP address")
    device_type: DeviceType = DeviceType.OTHER
    vendor: Optional[str] = Field(None, max_length=100)
    status: DeviceStatus = DeviceStatus.UNKNOWN


class DeviceCreate(DeviceBase):
    """Schema for creating a device."""
    pass


class DeviceUpdate(BaseModel):
    """Schema for updating a device (all fields optional)."""
    hostname: Optional[str] = Field(None, min_length=1, max_length=255)
    ip_address: Optional[str] = None
    device_type: Optional[DeviceType] = None
    vendor: Optional[str] = None
    status: Optional[DeviceStatus] = None


class Device(DeviceBase):
    """Full device model with ID and timestamps."""
    id: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class DeviceList(BaseModel):
    """Paginated device list response."""
    items: List[Device]
    total: int
    page: int
    page_size: int
    pages: int


class HealthCheck(BaseModel):
    """Health check response."""
    status: str
    timestamp: datetime
    version: str


# In-memory database (replace with real DB later)
_devices_db: dict = {}
_device_counter = 0


def get_next_id() -> int:
    """Generate next device ID."""
    global _device_counter
    _device_counter += 1
    return _device_counter


# Initialize FastAPI app
app = FastAPI(
    title="Network Inventory API",
    description="REST API for managing network device inventory",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# Health check endpoint
@app.get("/health", response_model=HealthCheck, tags=["Health"])
async def health_check():
    """Check API health status."""
    return HealthCheck(
        status="healthy",
        timestamp=datetime.now(),
        version="1.0.0"
    )


# Device endpoints
@app.post("/devices", response_model=Device, status_code=status.HTTP_201_CREATED,
          tags=["Devices"])
async def create_device(device: DeviceCreate):
    """Create a new network device."""
    # Check for duplicate hostname
    for existing in _devices_db.values():
        if existing["hostname"] == device.hostname:
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail=f"Device with hostname '{device.hostname}' already exists"
            )
    
    device_id = get_next_id()
    now = datetime.now()
    
    db_device = {
        "id": device_id,
        **device.model_dump(),
        "created_at": now,
        "updated_at": now,
    }
    _devices_db[device_id] = db_device
    
    return Device(**db_device)


@app.get("/devices", response_model=DeviceList, tags=["Devices"])
async def list_devices(
    status: Optional[DeviceStatus] = Query(None, description="Filter by status"),
    vendor: Optional[str] = Query(None, description="Filter by vendor"),
    device_type: Optional[DeviceType] = Query(None, description="Filter by type"),
    search: Optional[str] = Query(None, description="Search hostname or IP"),
    page: int = Query(1, ge=1, description="Page number"),
    page_size: int = Query(20, ge=1, le=100, description="Items per page"),
):
    """List all devices with optional filtering and pagination."""
    # Filter devices
    filtered = list(_devices_db.values())
    
    if status:
        filtered = [d for d in filtered if d["status"] == status]
    if vendor:
        filtered = [d for d in filtered if d.get("vendor") == vendor]
    if device_type:
        filtered = [d for d in filtered if d["device_type"] == device_type]
    if search:
        search_lower = search.lower()
        filtered = [
            d for d in filtered
            if search_lower in d["hostname"].lower() 
            or search_lower in d["ip_address"]
        ]
    
    # Pagination
    total = len(filtered)
    pages = (total + page_size - 1) // page_size
    start = (page - 1) * page_size
    end = start + page_size
    
    items = [Device(**d) for d in filtered[start:end]]
    
    return DeviceList(
        items=items,
        total=total,
        page=page,
        page_size=page_size,
        pages=pages
    )


@app.get("/devices/{device_id}", response_model=Device, tags=["Devices"])
async def get_device(device_id: int):
    """Get a specific device by ID."""
    if device_id not in _devices_db:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device with id {device_id} not found"
        )
    return Device(**_devices_db[device_id])


@app.put("/devices/{device_id}", response_model=Device, tags=["Devices"])
async def update_device(device_id: int, device_update: DeviceUpdate):
    """Update a device."""
    if device_id not in _devices_db:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device with id {device_id} not found"
        )
    
    # Update only provided fields
    update_data = device_update.model_dump(exclude_unset=True)
    if update_data:
        _devices_db[device_id].update(update_data)
        _devices_db[device_id]["updated_at"] = datetime.now()
    
    return Device(**_devices_db[device_id])


@app.delete("/devices/{device_id}", status_code=status.HTTP_204_NO_CONTENT,
            tags=["Devices"])
async def delete_device(device_id: int):
    """Delete a device."""
    if device_id not in _devices_db:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device with id {device_id} not found"
        )
    del _devices_db[device_id]


@app.patch("/devices/{device_id}/status", response_model=Device, tags=["Devices"])
async def update_device_status(device_id: int, new_status: DeviceStatus):
    """Update only the device status."""
    if device_id not in _devices_db:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device with id {device_id} not found"
        )
    
    _devices_db[device_id]["status"] = new_status
    _devices_db[device_id]["updated_at"] = datetime.now()
    
    return Device(**_devices_db[device_id])


# Bulk operations
@app.post("/devices/bulk", response_model=List[Device], tags=["Devices"])
async def bulk_create_devices(devices: List[DeviceCreate]):
    """Create multiple devices at once."""
    created = []
    for device in devices:
        device_id = get_next_id()
        now = datetime.now()
        db_device = {
            "id": device_id,
            **device.model_dump(),
            "created_at": now,
            "updated_at": now,
        }
        _devices_db[device_id] = db_device
        created.append(Device(**db_device))
    return created


# Statistics endpoint
@app.get("/stats", tags=["Statistics"])
async def get_statistics():
    """Get inventory statistics."""
    devices = list(_devices_db.values())
    
    status_counts = {}
    vendor_counts = {}
    type_counts = {}
    
    for d in devices:
        # Count by status
        s = d["status"]
        status_counts[s] = status_counts.get(s, 0) + 1
        
        # Count by vendor
        v = d.get("vendor", "unknown")
        vendor_counts[v] = vendor_counts.get(v, 0) + 1
        
        # Count by type
        t = d["device_type"]
        type_counts[t] = type_counts.get(t, 0) + 1
    
    return {
        "total_devices": len(devices),
        "by_status": status_counts,
        "by_vendor": vendor_counts,
        "by_type": type_counts,
    }


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
```

#### Running the API:

```bash
# Development mode with auto-reload
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Or directly
python -m app.main
```

#### Testing the API:

```bash
# Health check
curl http://localhost:8000/health

# Create device
curl -X POST http://localhost:8000/devices \
  -H "Content-Type: application/json" \
  -d '{"hostname": "router-01", "ip_address": "10.0.0.1", "device_type": "router", "vendor": "juniper"}'

# List devices
curl http://localhost:8000/devices

# Get specific device
curl http://localhost:8000/devices/1

# Update device
curl -X PUT http://localhost:8000/devices/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "up"}'

# Delete device
curl -X DELETE http://localhost:8000/devices/1

# Interactive docs
# Open: http://localhost:8000/docs
```

**Git checkpoint:**
```bash
git add network_api/
git commit -m "Day 13-14: FastAPI REST application"
```

---

### Day 15-16: Database Integration & Advanced Features

#### Exercise 15.1: FastAPI with PostgreSQL

Create `network_api/app/database.py`:

```python
"""Database configuration and session management."""

import os
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql://netadmin:secretpassword@localhost:5432/network_inventory"
)

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()


def get_db():
    """Dependency for database sessions."""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
```

Create `network_api/app/models.py`:

```python
"""SQLAlchemy ORM models."""

from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, Enum
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from .database import Base
import enum


class DeviceStatus(str, enum.Enum):
    UP = "up"
    DOWN = "down"
    MAINTENANCE = "maintenance"
    UNKNOWN = "unknown"


class DeviceType(str, enum.Enum):
    ROUTER = "router"
    SWITCH = "switch"
    FIREWALL = "firewall"
    SERVER = "server"
    OTHER = "other"


class Device(Base):
    """Device ORM model."""
    __tablename__ = "devices"
    
    id = Column(Integer, primary_key=True, index=True)
    hostname = Column(String(255), unique=True, nullable=False, index=True)
    ip_address = Column(String(45), nullable=False)
    device_type = Column(String(50), default="other")
    vendor = Column(String(100))
    status = Column(String(50), default="unknown", index=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), 
                       onupdate=func.now())
    
    # Relationship to interfaces
    interfaces = relationship("Interface", back_populates="device",
                             cascade="all, delete-orphan")


class Interface(Base):
    """Network interface ORM model."""
    __tablename__ = "interfaces"
    
    id = Column(Integer, primary_key=True, index=True)
    device_id = Column(Integer, ForeignKey("devices.id"), nullable=False)
    name = Column(String(100), nullable=False)
    ip_address = Column(String(45))
    mac_address = Column(String(17))
    status = Column(String(50), default="down")
    speed_mbps = Column(Integer)
    mtu = Column(Integer, default=1500)
    
    device = relationship("Device", back_populates="interfaces")
```

Create `network_api/app/crud.py`:

```python
"""CRUD operations for database models."""

from sqlalchemy.orm import Session
from sqlalchemy import or_
from typing import List, Optional
from . import models, schemas


def get_device(db: Session, device_id: int) -> Optional[models.Device]:
    """Get device by ID."""
    return db.query(models.Device).filter(models.Device.id == device_id).first()


def get_device_by_hostname(db: Session, hostname: str) -> Optional[models.Device]:
    """Get device by hostname."""
    return db.query(models.Device).filter(models.Device.hostname == hostname).first()


def get_devices(
    db: Session,
    skip: int = 0,
    limit: int = 100,
    status: str = None,
    vendor: str = None,
    device_type: str = None,
    search: str = None,
) -> List[models.Device]:
    """Get devices with filtering."""
    query = db.query(models.Device)
    
    if status:
        query = query.filter(models.Device.status == status)
    if vendor:
        query = query.filter(models.Device.vendor == vendor)
    if device_type:
        query = query.filter(models.Device.device_type == device_type)
    if search:
        query = query.filter(
            or_(
                models.Device.hostname.ilike(f"%{search}%"),
                models.Device.ip_address.ilike(f"%{search}%")
            )
        )
    
    return query.offset(skip).limit(limit).all()


def count_devices(db: Session, **filters) -> int:
    """Count devices with optional filters."""
    query = db.query(models.Device)
    if filters.get("status"):
        query = query.filter(models.Device.status == filters["status"])
    if filters.get("vendor"):
        query = query.filter(models.Device.vendor == filters["vendor"])
    return query.count()


def create_device(db: Session, device: schemas.DeviceCreate) -> models.Device:
    """Create a new device."""
    db_device = models.Device(**device.model_dump())
    db.add(db_device)
    db.commit()
    db.refresh(db_device)
    return db_device


def update_device(db: Session, device_id: int, 
                  device_update: schemas.DeviceUpdate) -> Optional[models.Device]:
    """Update a device."""
    db_device = get_device(db, device_id)
    if not db_device:
        return None
    
    update_data = device_update.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(db_device, key, value)
    
    db.commit()
    db.refresh(db_device)
    return db_device


def delete_device(db: Session, device_id: int) -> bool:
    """Delete a device."""
    db_device = get_device(db, device_id)
    if not db_device:
        return False
    
    db.delete(db_device)
    db.commit()
    return True


def bulk_create_devices(db: Session, 
                        devices: List[schemas.DeviceCreate]) -> List[models.Device]:
    """Bulk create devices."""
    db_devices = [models.Device(**d.model_dump()) for d in devices]
    db.add_all(db_devices)
    db.commit()
    for d in db_devices:
        db.refresh(d)
    return db_devices
```

---

## Phase 6: Production Deployment

### Day 17-18: Systemd Service & Configuration

#### Exercise 17.1: Systemd Unit File

Create `deployment/network-api.service`:

```ini
[Unit]
Description=Network Inventory API Service
Documentation=https://github.com/yourusername/network-api
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=exec
User=netapi
Group=netapi
WorkingDirectory=/opt/network-api
Environment="PATH=/opt/network-api/venv/bin"
EnvironmentFile=/etc/network-api/config.env
ExecStart=/opt/network-api/venv/bin/uvicorn app.main:app --host 0.0.0.0 --port 8000 --workers 4
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=network-api

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ReadWritePaths=/var/log/network-api

[Install]
WantedBy=multi-user.target
```

Create `deployment/config.env.example`:

```bash
# Database configuration
DATABASE_URL=postgresql://netadmin:changeme@localhost:5432/network_inventory
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=network_inventory
POSTGRES_USER=netadmin
POSTGRES_PASSWORD=changeme

# Application settings
API_HOST=0.0.0.0
API_PORT=8000
API_WORKERS=4
LOG_LEVEL=INFO

# Security
SECRET_KEY=your-secret-key-here
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
```

#### Exercise 17.2: Deployment Script

Create `deployment/deploy.sh`:

```bash
#!/bin/bash
# Network API Deployment Script

set -euo pipefail

# Configuration
APP_NAME="network-api"
APP_USER="netapi"
APP_DIR="/opt/${APP_NAME}"
CONFIG_DIR="/etc/${APP_NAME}"
LOG_DIR="/var/log/${APP_NAME}"
VENV_DIR="${APP_DIR}/venv"
REPO_URL="${REPO_URL:-https://github.com/yourusername/network-api.git}"
BRANCH="${BRANCH:-main}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

# Create application user
create_user() {
    log_info "Creating application user..."
    if ! id "${APP_USER}" &>/dev/null; then
        useradd --system --shell /bin/false --home-dir "${APP_DIR}" "${APP_USER}"
        log_info "Created user ${APP_USER}"
    else
        log_info "User ${APP_USER} already exists"
    fi
}

# Create directories
create_directories() {
    log_info "Creating directories..."
    mkdir -p "${APP_DIR}" "${CONFIG_DIR}" "${LOG_DIR}"
    chown "${APP_USER}:${APP_USER}" "${APP_DIR}" "${LOG_DIR}"
}

# Install system dependencies
install_dependencies() {
    log_info "Installing system dependencies..."
    apt-get update
    apt-get install -y python3 python3-venv python3-pip postgresql-client git
}

# Clone or update repository
setup_code() {
    log_info "Setting up application code..."
    
    if [[ -d "${APP_DIR}/.git" ]]; then
        cd "${APP_DIR}"
        git fetch origin
        git checkout "${BRANCH}"
        git pull origin "${BRANCH}"
    else
        git clone --branch "${BRANCH}" "${REPO_URL}" "${APP_DIR}"
    fi
    
    chown -R "${APP_USER}:${APP_USER}" "${APP_DIR}"
}

# Setup Python virtual environment
setup_venv() {
    log_info "Setting up Python virtual environment..."
    
    if [[ ! -d "${VENV_DIR}" ]]; then
        python3 -m venv "${VENV_DIR}"
    fi
    
    "${VENV_DIR}/bin/pip" install --upgrade pip
    "${VENV_DIR}/bin/pip" install -r "${APP_DIR}/requirements.txt"
}

# Setup configuration
setup_config() {
    log_info "Setting up configuration..."
    
    if [[ ! -f "${CONFIG_DIR}/config.env" ]]; then
        cp "${APP_DIR}/deployment/config.env.example" "${CONFIG_DIR}/config.env"
        chmod 600 "${CONFIG_DIR}/config.env"
        log_warn "Please edit ${CONFIG_DIR}/config.env with your settings"
    fi
}

# Install systemd service
install_service() {
    log_info "Installing systemd service..."
    
    cp "${APP_DIR}/deployment/${APP_NAME}.service" "/etc/systemd/system/"
    systemctl daemon-reload
    systemctl enable "${APP_NAME}"
}

# Database setup
setup_database() {
    log_info "Running database migrations..."
    
    # Source configuration
    source "${CONFIG_DIR}/config.env"
    
    # Run migrations (if using alembic)
    cd "${APP_DIR}"
    "${VENV_DIR}/bin/alembic" upgrade head 2>/dev/null || log_warn "No migrations to run"
}

# Start/restart service
start_service() {
    log_info "Starting service..."
    systemctl restart "${APP_NAME}"
    systemctl status "${APP_NAME}" --no-pager
}

# Health check
health_check() {
    log_info "Running health check..."
    sleep 3
    
    if curl -sf http://localhost:8000/health > /dev/null; then
        log_info "Health check passed!"
    else
        log_error "Health check failed!"
        journalctl -u "${APP_NAME}" --no-pager -n 50
        exit 1
    fi
}

# Main deployment function
deploy() {
    log_info "Starting deployment of ${APP_NAME}..."
    
    check_root
    install_dependencies
    create_user
    create_directories
    setup_code
    setup_venv
    setup_config
    install_service
    setup_database
    start_service
    health_check
    
    log_info "Deployment completed successfully!"
}

# Command handling
case "${1:-deploy}" in
    deploy)
        deploy
        ;;
    restart)
        start_service
        ;;
    status)
        systemctl status "${APP_NAME}"
        ;;
    logs)
        journalctl -u "${APP_NAME}" -f
        ;;
    *)
        echo "Usage: $0 {deploy|restart|status|logs}"
        exit 1
        ;;
esac
```

---

## Git & GitHub Workflow

### Repository Setup

```bash
# Initialize repository
git init
git branch -M main

# Create .gitignore
cat > .gitignore << 'EOF'
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
venv/
ENV/
.venv/
env/

# IDEs
.idea/
.vscode/
*.swp
*.swo

# Testing
.pytest_cache/
.coverage
htmlcov/
.tox/

# Build
dist/
build/
*.egg-info/

# Environment
.env
*.env.local
config.env

# Logs
*.log
logs/

# Database
*.db
*.sqlite3

# OS
.DS_Store
Thumbs.db
EOF

# Create README
cat > README.md << 'EOF'
# Network Inventory API

A RESTful API for managing network device inventory.

## Features

- CRUD operations for network devices
- PostgreSQL database with SQLAlchemy ORM
- FastAPI with automatic OpenAPI documentation
- Systemd service for production deployment
- Jenkins CI/CD pipeline

## Quick Start

```bash
# Clone repository
git clone https://github.com/yourusername/network-api.git
cd network-api

# Setup virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Run development server
uvicorn app.main:app --reload
```

## API Documentation

- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## Deployment

See `deployment/` directory for production deployment instructions.

## License

MIT
EOF

# Initial commit
git add .
git commit -m "Initial commit: Network Inventory API"

# Connect to GitHub
git remote add origin https://github.com/yourusername/network-api.git
git push -u origin main
```

### Branching Strategy

```bash
# Create feature branch
git checkout -b feature/add-interface-endpoints

# Make changes and commit
git add .
git commit -m "feat: add interface CRUD endpoints"

# Push feature branch
git push -u origin feature/add-interface-endpoints

# Create pull request on GitHub
# After review and merge:
git checkout main
git pull origin main
git branch -d feature/add-interface-endpoints
```

### Commit Message Convention

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance

Examples:
```bash
git commit -m "feat(api): add bulk device creation endpoint"
git commit -m "fix(db): handle connection pool exhaustion"
git commit -m "docs: update API documentation"
git commit -m "test: add unit tests for device CRUD"
```

---

## Jenkins CI/CD Pipeline

### Jenkinsfile

Create `Jenkinsfile` in repository root:

```groovy
pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = 'network-api'
        DOCKER_TAG = "${env.BUILD_NUMBER}"
        REGISTRY = 'your-registry.com'
        DEPLOY_HOST = 'deploy@your-server.com'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                sh 'git log --oneline -5'
            }
        }
        
        stage('Setup Python') {
            steps {
                sh '''
                    python3 -m venv venv
                    . venv/bin/activate
                    pip install --upgrade pip
                    pip install -r requirements.txt
                    pip install -r requirements-dev.txt
                '''
            }
        }
        
        stage('Lint') {
            steps {
                sh '''
                    . venv/bin/activate
                    echo "Running flake8..."
                    flake8 app/ --max-line-length=100 --ignore=E501
                    echo "Running black check..."
                    black --check app/
                    echo "Running isort check..."
                    isort --check-only app/
                '''
            }
        }
        
        stage('Type Check') {
            steps {
                sh '''
                    . venv/bin/activate
                    mypy app/ --ignore-missing-imports || true
                '''
            }
        }
        
        stage('Unit Tests') {
            steps {
                sh '''
                    . venv/bin/activate
                    pytest tests/unit/ \
                        --junitxml=test-results/unit.xml \
                        --cov=app \
                        --cov-report=xml:coverage.xml \
                        --cov-report=html:htmlcov \
                        -v
                '''
            }
            post {
                always {
                    junit 'test-results/*.xml'
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'htmlcov',
                        reportFiles: 'index.html',
                        reportName: 'Coverage Report'
                    ])
                }
            }
        }
        
        stage('Integration Tests') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                }
            }
            steps {
                sh '''
                    . venv/bin/activate
                    # Start test database
                    docker-compose -f docker-compose.test.yml up -d db
                    sleep 10
                    
                    # Run integration tests
                    pytest tests/integration/ \
                        --junitxml=test-results/integration.xml \
                        -v || true
                    
                    # Cleanup
                    docker-compose -f docker-compose.test.yml down
                '''
            }
        }
        
        stage('Build Docker Image') {
            when {
                branch 'main'
            }
            steps {
                sh '''
                    docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
                    docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
                '''
            }
        }
        
        stage('Push to Registry') {
            when {
                branch 'main'
            }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'docker-registry',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                        echo $DOCKER_PASS | docker login ${REGISTRY} -u $DOCKER_USER --password-stdin
                        docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}
                        docker push ${REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}
                        docker push ${REGISTRY}/${DOCKER_IMAGE}:latest
                    '''
                }
            }
        }
        
        stage('Deploy to Staging') {
            when {
                branch 'main'
            }
            steps {
                sshagent(['deploy-key']) {
                    sh '''
                        ssh ${DEPLOY_HOST} "
                            cd /opt/network-api && 
                            git pull origin main &&
                            ./deployment/deploy.sh restart
                        "
                    '''
                }
            }
        }
        
        stage('Smoke Tests') {
            when {
                branch 'main'
            }
            steps {
                sh '''
                    sleep 10
                    curl -f http://staging.your-server.com:8000/health || exit 1
                    echo "Smoke tests passed!"
                '''
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        success {
            slackSend(
                channel: '#deployments',
                color: 'good',
                message: "✅ ${env.JOB_NAME} #${env.BUILD_NUMBER} succeeded"
            )
        }
        failure {
            slackSend(
                channel: '#deployments',
                color: 'danger',
                message: "❌ ${env.JOB_NAME} #${env.BUILD_NUMBER} failed"
            )
        }
    }
}
```

### requirements-dev.txt

```
# Development dependencies
pytest>=7.0.0
pytest-cov>=4.0.0
pytest-asyncio>=0.21.0
httpx>=0.24.0
black>=23.0.0
flake8>=6.0.0
isort>=5.12.0
mypy>=1.0.0
pre-commit>=3.0.0
```

---

## Project Scaffolding Templates

### Complete Project Structure

```
network-api/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions alternative
├── app/
│   ├── __init__.py
│   ├── main.py                 # FastAPI application
│   ├── config.py               # Configuration management
│   ├── database.py             # Database connection
│   ├── models.py               # SQLAlchemy models
│   ├── schemas.py              # Pydantic schemas
│   ├── crud.py                 # CRUD operations
│   └── routers/
│       ├── __init__.py
│       ├── devices.py          # Device endpoints
│       ├── interfaces.py       # Interface endpoints
│       └── stats.py            # Statistics endpoints
├── tests/
│   ├── __init__.py
│   ├── conftest.py             # Pytest fixtures
│   ├── unit/
│   │   ├── __init__.py
│   │   ├── test_models.py
│   │   └── test_crud.py
│   └── integration/
│       ├── __init__.py
│       └── test_api.py
├── deployment/
│   ├── network-api.service     # Systemd unit
│   ├── config.env.example      # Example configuration
│   ├── deploy.sh               # Deployment script
│   └── nginx.conf              # Nginx reverse proxy
├── migrations/
│   ├── versions/               # Alembic migrations
│   └── env.py
├── .gitignore
├── .pre-commit-config.yaml
├── Dockerfile
├── docker-compose.yml
├── docker-compose.test.yml
├── Jenkinsfile
├── Makefile
├── README.md
├── requirements.txt
├── requirements-dev.txt
└── pyproject.toml
```

### Makefile for Common Tasks

```makefile
.PHONY: help install dev test lint format run deploy clean

VENV := venv
PYTHON := $(VENV)/bin/python
PIP := $(VENV)/bin/pip

help:
	@echo "Available targets:"
	@echo "  install  - Install production dependencies"
	@echo "  dev      - Install development dependencies"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linters"
	@echo "  format   - Format code"
	@echo "  run      - Run development server"
	@echo "  deploy   - Deploy to production"
	@echo "  clean    - Clean build artifacts"

$(VENV)/bin/activate:
	python3 -m venv $(VENV)

install: $(VENV)/bin/activate
	$(PIP) install --upgrade pip
	$(PIP) install -r requirements.txt

dev: install
	$(PIP) install -r requirements-dev.txt
	pre-commit install

test:
	$(PYTHON) -m pytest tests/ -v --cov=app

test-unit:
	$(PYTHON) -m pytest tests/unit/ -v

test-integration:
	$(PYTHON) -m pytest tests/integration/ -v

lint:
	$(PYTHON) -m flake8 app/ tests/
	$(PYTHON) -m black --check app/ tests/
	$(PYTHON) -m isort --check-only app/ tests/
	$(PYTHON) -m mypy app/

format:
	$(PYTHON) -m black app/ tests/
	$(PYTHON) -m isort app/ tests/

run:
	$(PYTHON) -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

run-prod:
	$(PYTHON) -m uvicorn app.main:app --host 0.0.0.0 --port 8000 --workers 4

docker-build:
	docker build -t network-api:latest .

docker-run:
	docker-compose up -d

deploy:
	./deployment/deploy.sh deploy

clean:
	rm -rf $(VENV)
	rm -rf __pycache__
	rm -rf .pytest_cache
	rm -rf htmlcov
	rm -rf .coverage
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
```

---

## Capstone Project

### Complete Network Inventory Service

Your final project combines everything learned:

1. **Python Fundamentals**: Classes, modules, error handling
2. **Data Processing**: JSON/YAML parsing, file operations
3. **Database**: PostgreSQL with SQLAlchemy ORM
4. **Web API**: FastAPI with full CRUD operations
5. **Deployment**: Systemd service, health checks
6. **CI/CD**: Jenkins pipeline with testing and deployment
7. **Git**: Proper branching, commits, and GitHub workflow

### Final Exercise

1. Fork the template repository
2. Implement all endpoints from the exercises
3. Add unit and integration tests
4. Set up Jenkins pipeline
5. Deploy to a server as a systemd service
6. Document the API and deployment process

### Success Criteria

- [ ] API responds to all CRUD operations
- [ ] Tests pass with >80% coverage
- [ ] Jenkins pipeline builds and deploys successfully
- [ ] Service runs reliably with systemd
- [ ] Documentation is complete and accurate

---

## Quick Reference Commands

```bash
# Virtual environment
python3 -m venv venv
source venv/bin/activate
deactivate

# Package management
pip install package_name
pip install -r requirements.txt
pip freeze > requirements.txt

# Running code
python script.py
python -m module_name
uvicorn app.main:app --reload

# Testing
pytest
pytest -v --cov=app
pytest tests/unit/ -k "test_create"

# Git workflow
git checkout -b feature/new-feature
git add .
git commit -m "feat: add new feature"
git push -u origin feature/new-feature

# Systemd
sudo systemctl start network-api
sudo systemctl status network-api
sudo journalctl -u network-api -f

# Docker
docker build -t myapp .
docker run -p 8000:8000 myapp
docker-compose up -d
```

---

**Congratulations!** You now have a complete Python learning path from basics to production deployment. Each exercise builds on the previous, and the final project demonstrates all skills combined.
