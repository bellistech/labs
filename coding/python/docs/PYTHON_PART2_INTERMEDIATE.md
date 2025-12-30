# Python Course Part 2: Intermediate (Days 4-6)

## Explain Like I'm 5: Objects and Classes

Imagine you're playing with LEGO. You have a **blueprint** for building a car - it shows you exactly what pieces to use and how to put them together. You can use that same blueprint to build 10 different cars - each car is its own thing, but they all follow the same blueprint.

In programming:
- **Class** = The blueprint (instructions for making something)
- **Object** = A thing built from the blueprint (an actual car)
- **Instance** = Same as object (the car you're holding)

---

## Day 4: Object-Oriented Programming

### Exercise 4.1: Your First Class

```python
#!/usr/bin/env python3
# =============================================================================
# CLASSES AND OBJECTS
# =============================================================================
#
# ELI5: A class is like a cookie cutter. It defines the shape.
# Each cookie you make is an "object" or "instance" of that class.
#

class Dog:
    """
    A simple Dog class.
    
    ELI5: This is the blueprint for making dogs.
    It says every dog has a name and age.
    """
    
    # __init__ is the "constructor" - it runs when you create a new Dog
    #
    # ELI5: This is like the birth of the dog. You give it a name and age.
    #
    def __init__(self, name: str, age: int):
        # 'self' refers to THIS specific dog
        # self.name = "put the name in THIS dog's name tag"
        self.name = name
        self.age = age
    
    # Methods are functions that belong to the class
    #
    # ELI5: Things this dog can DO
    #
    def bark(self):
        print(f"{self.name} says: Woof!")
    
    def describe(self):
        print(f"{self.name} is {self.age} years old")


# Create objects (instances) from the class
my_dog = Dog("Buddy", 3)
your_dog = Dog("Max", 5)

# Each dog is separate!
my_dog.bark()        # Buddy says: Woof!
your_dog.bark()      # Max says: Woof!

my_dog.describe()    # Buddy is 3 years old
your_dog.describe()  # Max is 5 years old

# Access attributes
print(f"My dog's name is {my_dog.name}")
```

### Exercise 4.2: Class with More Features

```python
#!/usr/bin/env python3
# =============================================================================
# NETWORK DEVICE CLASS
# =============================================================================

class NetworkDevice:
    """
    Represents a network device (router, switch, server).
    
    ELI5: This is a blueprint for any device on your network.
    """
    
    # Class attribute - shared by ALL instances
    device_count = 0
    
    def __init__(self, hostname: str, ip_address: str, device_type: str):
        # Instance attributes - unique to each device
        self.hostname = hostname
        self.ip_address = ip_address
        self.device_type = device_type
        self.is_online = False
        
        # Update class attribute
        NetworkDevice.device_count += 1
    
    def connect(self) -> bool:
        """Simulate connecting to the device."""
        print(f"Connecting to {self.hostname} ({self.ip_address})...")
        self.is_online = True
        return True
    
    def disconnect(self):
        """Disconnect from the device."""
        print(f"Disconnecting from {self.hostname}")
        self.is_online = False
    
    def get_status(self) -> str:
        """Get device status."""
        status = "ONLINE" if self.is_online else "OFFLINE"
        return f"{self.hostname} [{self.device_type}] - {status}"
    
    def __str__(self):
        """Called when you print() the object."""
        return f"NetworkDevice({self.hostname}, {self.ip_address})"
    
    def __repr__(self):
        """Called in debugging - more detailed."""
        return f"NetworkDevice(hostname='{self.hostname}', ip='{self.ip_address}', type='{self.device_type}')"


# Create some devices
router = NetworkDevice("router-01", "192.168.1.1", "router")
switch = NetworkDevice("switch-01", "192.168.1.2", "switch")
server = NetworkDevice("web-01", "192.168.1.100", "server")

print(f"Total devices: {NetworkDevice.device_count}")  # 3

router.connect()
print(router.get_status())  # router-01 [router] - ONLINE

print(router)  # Uses __str__
```

### Exercise 4.3: Inheritance

```python
#!/usr/bin/env python3
# =============================================================================
# INHERITANCE - BUILDING ON EXISTING CLASSES
# =============================================================================
#
# ELI5: Inheritance is like a family tree.
# A child class inherits everything from its parent class,
# then can add its own special features.
#

class Device:
    """Base class for all devices."""
    
    def __init__(self, hostname: str, ip_address: str):
        self.hostname = hostname
        self.ip_address = ip_address
        self.is_online = False
    
    def power_on(self):
        self.is_online = True
        print(f"{self.hostname} powered on")
    
    def power_off(self):
        self.is_online = False
        print(f"{self.hostname} powered off")


class Router(Device):
    """Router inherits from Device and adds routing features."""
    
    def __init__(self, hostname: str, ip_address: str, routing_protocol: str):
        # Call parent's __init__ first
        super().__init__(hostname, ip_address)
        
        # Add router-specific attributes
        self.routing_protocol = routing_protocol
        self.routes = []
    
    def add_route(self, network: str, next_hop: str):
        """Add a routing table entry."""
        self.routes.append({"network": network, "next_hop": next_hop})
        print(f"Added route to {network} via {next_hop}")
    
    def show_routes(self):
        """Display the routing table."""
        print(f"Routing table for {self.hostname}:")
        for route in self.routes:
            print(f"  {route['network']} -> {route['next_hop']}")


class Switch(Device):
    """Switch inherits from Device and adds switching features."""
    
    def __init__(self, hostname: str, ip_address: str, num_ports: int):
        super().__init__(hostname, ip_address)
        self.num_ports = num_ports
        self.vlans = []
    
    def create_vlan(self, vlan_id: int, name: str):
        """Create a VLAN."""
        self.vlans.append({"id": vlan_id, "name": name})
        print(f"Created VLAN {vlan_id}: {name}")


# Use the classes
router = Router("core-rtr-01", "10.0.0.1", "OSPF")
router.power_on()  # Inherited from Device!
router.add_route("192.168.1.0/24", "10.0.0.2")
router.add_route("192.168.2.0/24", "10.0.0.3")
router.show_routes()

switch = Switch("access-sw-01", "10.0.0.10", 48)
switch.power_on()
switch.create_vlan(10, "Users")
switch.create_vlan(20, "Servers")
```

---

## Day 5: Error Handling

### Exercise 5.1: Try/Except Basics

```python
#!/usr/bin/env python3
# =============================================================================
# ERROR HANDLING WITH TRY/EXCEPT
# =============================================================================
#
# ELI5: Try/except is like having a safety net.
# "Try to do this thing. If something goes wrong, do this instead."
#

# Without error handling - program crashes!
# result = 10 / 0  # ZeroDivisionError!

# With error handling - program continues
try:
    result = 10 / 0
except ZeroDivisionError:
    print("Oops! Can't divide by zero!")
    result = 0

print(f"Result: {result}")

# -----------------------------------------------------------------------------
# CATCHING MULTIPLE ERRORS
# -----------------------------------------------------------------------------

def get_list_item(items: list, index: int):
    """Get an item from a list, handling errors gracefully."""
    try:
        return items[index]
    except IndexError:
        print(f"Index {index} is out of range!")
        return None
    except TypeError:
        print("Index must be an integer!")
        return None


my_list = [1, 2, 3]
print(get_list_item(my_list, 0))   # 1
print(get_list_item(my_list, 10))  # None (with error message)
print(get_list_item(my_list, "a")) # None (with error message)


# -----------------------------------------------------------------------------
# FINALLY BLOCK
# -----------------------------------------------------------------------------
#
# 'finally' ALWAYS runs, whether there was an error or not.
# Great for cleanup!
#

def read_file_safely(filename: str) -> str:
    """Read a file with proper error handling."""
    file = None
    try:
        file = open(filename, 'r')
        content = file.read()
        return content
    except FileNotFoundError:
        print(f"File not found: {filename}")
        return ""
    except PermissionError:
        print(f"Permission denied: {filename}")
        return ""
    finally:
        # This runs no matter what!
        if file:
            file.close()
            print("File closed.")


# -----------------------------------------------------------------------------
# RAISING YOUR OWN ERRORS
# -----------------------------------------------------------------------------

def validate_port(port: int) -> int:
    """Validate a port number."""
    if not isinstance(port, int):
        raise TypeError("Port must be an integer")
    
    if port < 0 or port > 65535:
        raise ValueError(f"Port must be 0-65535, got {port}")
    
    return port


try:
    port = validate_port(80)       # OK
    port = validate_port(99999)    # Raises ValueError
except ValueError as e:
    print(f"Invalid port: {e}")
```

### Exercise 5.2: Custom Exceptions

```python
#!/usr/bin/env python3
# =============================================================================
# CUSTOM EXCEPTIONS
# =============================================================================

class NetworkError(Exception):
    """Base exception for network errors."""
    pass


class ConnectionError(NetworkError):
    """Raised when connection fails."""
    def __init__(self, host: str, port: int, message: str = None):
        self.host = host
        self.port = port
        self.message = message or f"Failed to connect to {host}:{port}"
        super().__init__(self.message)


class TimeoutError(NetworkError):
    """Raised when operation times out."""
    def __init__(self, operation: str, timeout: float):
        self.operation = operation
        self.timeout = timeout
        super().__init__(f"Operation '{operation}' timed out after {timeout}s")


def connect_to_server(host: str, port: int, timeout: float = 5.0):
    """Simulate connecting to a server."""
    import random
    
    # Simulate random failures for demo
    outcome = random.choice(["success", "connection_failed", "timeout"])
    
    if outcome == "connection_failed":
        raise ConnectionError(host, port)
    elif outcome == "timeout":
        raise TimeoutError("connect", timeout)
    
    print(f"Connected to {host}:{port}")


# Using the custom exceptions
for _ in range(3):
    try:
        connect_to_server("192.168.1.1", 22)
    except ConnectionError as e:
        print(f"Connection failed: {e}")
    except TimeoutError as e:
        print(f"Timeout: {e}")
    except NetworkError as e:
        print(f"Network error: {e}")
```

---

## Day 6: Decorators and Generators

### Exercise 6.1: Decorators

```python
#!/usr/bin/env python3
# =============================================================================
# DECORATORS - WRAPPING FUNCTIONS
# =============================================================================
#
# ELI5: A decorator is like wrapping a gift. You take a function,
# wrap it with extra behavior, and give back the wrapped version.
#

import time
from functools import wraps


def timer(func):
    """
    Decorator that measures how long a function takes.
    
    ELI5: Before calling the function, start a stopwatch.
    After it finishes, stop the stopwatch and print the time.
    """
    @wraps(func)  # Preserves the function's name and docstring
    def wrapper(*args, **kwargs):
        start_time = time.time()
        result = func(*args, **kwargs)
        end_time = time.time()
        print(f"{func.__name__} took {end_time - start_time:.4f} seconds")
        return result
    return wrapper


@timer  # This is the same as: slow_function = timer(slow_function)
def slow_function():
    """A function that takes a while."""
    time.sleep(1)
    return "Done!"


result = slow_function()  # Prints: slow_function took 1.00XX seconds


# -----------------------------------------------------------------------------
# DECORATOR WITH ARGUMENTS
# -----------------------------------------------------------------------------

def retry(max_attempts: int = 3, delay: float = 1.0):
    """
    Decorator that retries a function if it fails.
    
    ELI5: If at first you don't succeed, try, try again!
    """
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            last_exception = None
            
            for attempt in range(1, max_attempts + 1):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    last_exception = e
                    print(f"Attempt {attempt} failed: {e}")
                    if attempt < max_attempts:
                        time.sleep(delay)
            
            raise last_exception
        return wrapper
    return decorator


@retry(max_attempts=3, delay=0.5)
def unreliable_network_call():
    """A function that sometimes fails."""
    import random
    if random.random() < 0.7:  # 70% chance of failure
        raise ConnectionError("Network error!")
    return "Success!"


try:
    result = unreliable_network_call()
    print(result)
except ConnectionError:
    print("All retries failed!")
```

### Exercise 6.2: Generators

```python
#!/usr/bin/env python3
# =============================================================================
# GENERATORS - LAZY SEQUENCES
# =============================================================================
#
# ELI5: A generator is like a vending machine.
# It doesn't make all the snacks at once - it makes one snack
# each time you push the button. This saves memory!
#

def count_up_to(max_num: int):
    """
    Generator that counts from 1 to max_num.
    
    'yield' is like 'return' but pauses the function
    instead of ending it.
    """
    num = 1
    while num <= max_num:
        yield num  # Return this value, then pause here
        num += 1


# Use the generator
for n in count_up_to(5):
    print(n)  # 1, 2, 3, 4, 5


# -----------------------------------------------------------------------------
# WHY GENERATORS MATTER: MEMORY EFFICIENCY
# -----------------------------------------------------------------------------

def get_big_list():
    """Creates a list of 1 million numbers - uses lots of memory!"""
    return [i for i in range(1_000_000)]


def get_big_generator():
    """Generates 1 million numbers - uses almost no memory!"""
    for i in range(1_000_000):
        yield i


# The generator only holds one number at a time in memory
for num in get_big_generator():
    if num > 10:
        break
    print(num)


# -----------------------------------------------------------------------------
# GENERATOR EXPRESSIONS
# -----------------------------------------------------------------------------

# List comprehension - creates all items immediately
squares_list = [x**2 for x in range(10)]

# Generator expression - creates items on demand
squares_gen = (x**2 for x in range(10))

print(type(squares_list))  # <class 'list'>
print(type(squares_gen))   # <class 'generator'>


# -----------------------------------------------------------------------------
# PRACTICAL EXAMPLE: LOG FILE READER
# -----------------------------------------------------------------------------

def read_log_lines(filename: str):
    """Generator that yields log lines one at a time."""
    with open(filename, 'r') as f:
        for line in f:
            yield line.strip()


def filter_errors(lines):
    """Generator that filters for error lines."""
    for line in lines:
        if 'ERROR' in line:
            yield line


# Chain generators together!
# This can process HUGE files without loading them into memory
# log_lines = read_log_lines('huge_logfile.log')
# error_lines = filter_errors(log_lines)
# for error in error_lines:
#     print(error)
```

---

## Summary: What You Learned

### Day 4
- ✅ Classes and objects (blueprints and instances)
- ✅ `__init__` constructor
- ✅ Instance attributes vs class attributes
- ✅ Methods
- ✅ Inheritance and `super()`
- ✅ Special methods (`__str__`, `__repr__`)

### Day 5
- ✅ Try/except for error handling
- ✅ Multiple except blocks
- ✅ Finally for cleanup
- ✅ Raising exceptions
- ✅ Custom exception classes

### Day 6
- ✅ Decorators for wrapping functions
- ✅ `@wraps` to preserve function metadata
- ✅ Decorators with arguments
- ✅ Generators with `yield`
- ✅ Generator expressions
- ✅ Memory-efficient data processing

---

## Next: Part 3 - Data and Files

Continue to `PYTHON_PART3_DATA.md` for:
- File I/O
- JSON and YAML parsing
- HTTP requests
- Working with system commands
