# Python Course Part 1: Fundamentals (Days 1-3)

## Explain Like I'm 5: What is Programming?

Imagine you have a very obedient but very literal robot helper. This robot will do 
EXACTLY what you tell it - nothing more, nothing less. If you say "make a sandwich" 
but forget to say "put the peanut butter on the bread BEFORE closing it," you'll 
get a weird sandwich.

**Programming** is writing very precise instructions for this literal robot (the computer).

**Python** is one of the languages you can use to write these instructions. It's 
designed to be readable and look almost like English.

---

## Day 1: Hello World and Basic Data Types

### Exercise 1.1: Your First Python Program

```python
#!/usr/bin/env python3
# =============================================================================
# HELLO WORLD - YOUR FIRST PYTHON PROGRAM
# =============================================================================
#
# ELI5: This is the simplest program possible.
# It just tells the computer to print some text to the screen.
#
# The #! line (called "shebang") tells Linux which program runs this file.
# It's like labeling a document "READ WITH WORD" vs "READ WITH NOTEPAD".
#

# print() is a function that displays text on the screen
# 
# ELI5: Think of print() as the computer's mouth.
# Whatever you put in the parentheses is what the computer "says".
#
print("Hello, World!")

# You can print numbers too
print(42)

# You can print multiple things separated by commas
# They'll be separated by spaces automatically
print("The answer is", 42)

# You can use different quote styles - they work the same
print('Single quotes work too!')
print("Double quotes are also fine!")

# Triple quotes let you write multiple lines
print("""
This is a
multi-line
string!
""")
```

**Run it:**
```bash
# Method 1: Using Python directly
python3 hello.py

# Method 2: Make it executable (Linux/Mac)
chmod +x hello.py
./hello.py
```

---

### Exercise 1.2: Variables - Labeled Boxes

```python
#!/usr/bin/env python3
# =============================================================================
# VARIABLES - STORING DATA IN LABELED BOXES
# =============================================================================
#
# ELI5: A variable is like a labeled box.
# You give the box a name, and you can put something inside.
# Later, you can look in the box by using its name.
#

# Creating variables is simple: name = value
#
# ELI5: We're saying "make a box called 'name', put 'Alice' inside"
#
name = "Alice"
age = 28
height = 5.7
is_student = True

# Print the contents of our boxes
print("Name:", name)
print("Age:", age)
print("Height:", height)
print("Is student:", is_student)

# You can change what's in a box
#
# ELI5: We're emptying the 'age' box and putting 29 in it
#
age = 29
print("Happy birthday! Now age is:", age)

# Variables can use values from other variables
#
# ELI5: We're making a new box and filling it by looking at other boxes
#
birth_year = 2024 - age
print("Birth year:", birth_year)

# Variable names can use letters, numbers, and underscores
# But they can't START with a number!
#
# Good names:
user_name = "bob"
userName = "bob"      # "camelCase" style
user_name_2 = "bob"   # Numbers OK, just not at the start
_private = "hidden"   # Starting underscore has special meaning (by convention)

# Bad names (these would cause errors):
# 2nd_user = "carol"  # Can't start with number!
# user-name = "dave"  # Hyphens not allowed!
# user name = "eve"   # Spaces not allowed!
```

---

### Exercise 1.3: Data Types - Different Kinds of Boxes

```python
#!/usr/bin/env python3
# =============================================================================
# DATA TYPES - DIFFERENT KINDS OF INFORMATION
# =============================================================================
#
# ELI5: Not all boxes hold the same kind of stuff.
# Some boxes hold numbers, some hold text, some hold yes/no answers.
# Python figures out what kind of box to use based on what you put in it.
#

# INTEGER (int) - Whole numbers
#
# ELI5: Numbers without decimal points. Used for counting things.
#
port = 8080
count = 42
negative = -10

print(f"Port {port} is an integer: {type(port)}")

# FLOAT - Decimal numbers
#
# ELI5: Numbers with decimal points. Used for measurements, percentages.
#
temperature = 98.6
percentage = 0.75
pi = 3.14159

print(f"Temperature {temperature} is a float: {type(temperature)}")

# STRING (str) - Text
#
# ELI5: Words, sentences, any text in quotes.
#
hostname = "server-01"
message = "Hello, World!"
empty = ""  # An empty string is still a string

print(f"Hostname '{hostname}' is a string: {type(hostname)}")

# BOOLEAN (bool) - True or False
#
# ELI5: Yes/No, On/Off, True/False. Only two possible values.
#
is_running = True
is_broken = False

print(f"is_running = {is_running}, type: {type(is_running)}")

# NONE - Nothing
#
# ELI5: An empty box. Not zero, not empty string, just... nothing.
# Useful for "this variable exists but has no value yet".
#
result = None
print(f"result = {result}, type: {type(result)}")

# -----------------------------------------------------------------------------
# TYPE CONVERSION
# -----------------------------------------------------------------------------
#
# ELI5: Sometimes you need to change one type to another.
# Like converting a word to a number.
#

# String to Integer
port_string = "8080"
port_number = int(port_string)
print(f"'{port_string}' (string) → {port_number} (int)")

# Integer to String
count = 42
count_string = str(count)
print(f"{count} (int) → '{count_string}' (string)")

# String to Float
price_string = "19.99"
price_number = float(price_string)
print(f"'{price_string}' (string) → {price_number} (float)")

# Float to Integer (rounds DOWN, doesn't round normally!)
temperature = 98.6
temp_int = int(temperature)  # Results in 98, not 99!
print(f"{temperature} (float) → {temp_int} (int) - decimals are DROPPED!")

# -----------------------------------------------------------------------------
# BOOLEAN TRUTHINESS
# -----------------------------------------------------------------------------
#
# ELI5: In Python, many things can be "truthy" or "falsy".
# Empty things are usually falsy. Things with content are truthy.
#

print("\n--- Truthiness ---")
print(f"bool(0) = {bool(0)}")           # False - zero is falsy
print(f"bool(1) = {bool(1)}")           # True - any non-zero is truthy
print(f"bool('') = {bool('')}")         # False - empty string is falsy
print(f"bool('hi') = {bool('hi')}")     # True - non-empty string is truthy
print(f"bool([]) = {bool([])}")         # False - empty list is falsy
print(f"bool([1]) = {bool([1])}")       # True - non-empty list is truthy
print(f"bool(None) = {bool(None)}")     # False - None is always falsy
```

---

### Exercise 1.4: Basic Math and String Operations

```python
#!/usr/bin/env python3
# =============================================================================
# BASIC OPERATIONS
# =============================================================================
#
# ELI5: You can do math with numbers and combine strings.
#

# -----------------------------------------------------------------------------
# MATH OPERATIONS
# -----------------------------------------------------------------------------

a = 10
b = 3

print("--- Basic Math ---")
print(f"{a} + {b} = {a + b}")    # Addition: 13
print(f"{a} - {b} = {a - b}")    # Subtraction: 7
print(f"{a} * {b} = {a * b}")    # Multiplication: 30
print(f"{a} / {b} = {a / b}")    # Division: 3.333... (always returns float!)
print(f"{a} // {b} = {a // b}")  # Floor division: 3 (rounds down to integer)
print(f"{a} % {b} = {a % b}")    # Modulo (remainder): 1
print(f"{a} ** {b} = {a ** b}")  # Exponent (power): 1000

# Practical example: converting bytes to kilobytes
bytes_used = 1536
kilobytes = bytes_used / 1024
print(f"\n{bytes_used} bytes = {kilobytes} KB")

# Practical example: checking if a number is even or odd
number = 17
remainder = number % 2
if remainder == 0:
    print(f"{number} is even")
else:
    print(f"{number} is odd")

# -----------------------------------------------------------------------------
# STRING OPERATIONS
# -----------------------------------------------------------------------------

print("\n--- String Operations ---")

first_name = "John"
last_name = "Doe"

# Concatenation (combining strings)
#
# ELI5: Gluing strings together like train cars
#
full_name = first_name + " " + last_name
print(f"Full name: {full_name}")

# String repetition
#
# ELI5: Make copies of a string
#
line = "-" * 40
print(line)

# String length
#
# ELI5: How many characters in the string?
#
print(f"Length of '{full_name}': {len(full_name)} characters")

# String methods (built-in operations)
#
# ELI5: Strings come with superpowers! Use a dot to access them.
#
message = "  Hello, World!  "

print(f"Original: '{message}'")
print(f"Upper: '{message.upper()}'")           # ALL CAPS
print(f"Lower: '{message.lower()}'")           # all lowercase
print(f"Strip: '{message.strip()}'")           # Remove extra spaces
print(f"Replace: '{message.replace('World', 'Python')}'")

# Checking string content
hostname = "web-server-01"
print(f"\nHostname: {hostname}")
print(f"Starts with 'web': {hostname.startswith('web')}")
print(f"Ends with '01': {hostname.endswith('01')}")
print(f"Contains 'server': {'server' in hostname}")

# -----------------------------------------------------------------------------
# F-STRINGS (FORMATTED STRINGS)
# -----------------------------------------------------------------------------
#
# ELI5: F-strings let you put variables inside strings easily.
# Just put 'f' before the quotes and use {curly braces} for variables.
#

name = "Alice"
age = 28
salary = 75000.50

# The f before the quote makes it a formatted string
print(f"Name: {name}, Age: {age}")

# You can do math inside the braces!
print(f"In 5 years, {name} will be {age + 5}")

# Format numbers nicely
print(f"Salary: ${salary:,.2f}")  # Adds commas and 2 decimal places
print(f"Percentage: {0.856:.1%}") # Shows as 85.6%

# Padding and alignment
print(f"{'left':<10} | {'center':^10} | {'right':>10}")
```

---

## Day 2: Control Flow and Data Structures

### Exercise 2.1: Conditionals (If/Else)

```python
#!/usr/bin/env python3
# =============================================================================
# CONDITIONALS - MAKING DECISIONS
# =============================================================================
#
# ELI5: Sometimes your code needs to make decisions.
# "If it's raining, bring an umbrella. Otherwise, wear sunglasses."
# That's exactly what if/else does!
#

# Basic if statement
#
# ELI5: "If this is true, do this thing"
#
temperature = 75

if temperature > 80:
    print("It's hot! Turn on the AC.")

# If/else - two options
#
# ELI5: "If this is true, do A. Otherwise, do B."
#
hour = 14

if hour < 12:
    print("Good morning!")
else:
    print("Good afternoon!")

# If/elif/else - multiple options
#
# ELI5: "If first thing is true, do A. Else if second thing is true, do B.
#        Otherwise, do C."
#
score = 85

if score >= 90:
    grade = "A"
elif score >= 80:
    grade = "B"
elif score >= 70:
    grade = "C"
elif score >= 60:
    grade = "D"
else:
    grade = "F"

print(f"Score {score} = Grade {grade}")

# -----------------------------------------------------------------------------
# COMPARISON OPERATORS
# -----------------------------------------------------------------------------
#
# ELI5: These compare two things and return True or False
#

x = 10
y = 5

print(f"\n--- Comparisons (x={x}, y={y}) ---")
print(f"x == y: {x == y}")  # Equal to
print(f"x != y: {x != y}")  # Not equal to
print(f"x > y: {x > y}")    # Greater than
print(f"x < y: {x < y}")    # Less than
print(f"x >= y: {x >= y}")  # Greater than or equal
print(f"x <= y: {x <= y}")  # Less than or equal

# -----------------------------------------------------------------------------
# LOGICAL OPERATORS
# -----------------------------------------------------------------------------
#
# ELI5: Combine multiple conditions
# - and: Both must be true
# - or: At least one must be true
# - not: Flip true to false, false to true
#

age = 25
has_license = True

# AND - both conditions must be true
if age >= 18 and has_license:
    print("You can drive!")

# OR - at least one condition must be true
is_weekend = True
is_holiday = False

if is_weekend or is_holiday:
    print("No work today!")

# NOT - reverse the condition
is_raining = False

if not is_raining:
    print("Leave the umbrella at home.")

# -----------------------------------------------------------------------------
# PRACTICAL EXAMPLE: NETWORK PORT CLASSIFIER
# -----------------------------------------------------------------------------

def classify_port(port: int) -> str:
    """
    Classify a network port number.
    
    ELI5: Network ports are like apartment numbers in a building (the server).
    Different "apartments" are reserved for different services.
    
    Args:
        port: The port number to classify
    
    Returns:
        A string describing the port type
    """
    if port < 0 or port > 65535:
        return "Invalid port number"
    elif port == 0:
        return "Reserved"
    elif port < 1024:
        return "Well-known/privileged (requires root)"
    elif port < 49152:
        return "Registered port"
    else:
        return "Dynamic/ephemeral port"


# Test the function
test_ports = [22, 80, 443, 3000, 8080, 49152, 65535, 70000]

print("\n--- Port Classification ---")
for port in test_ports:
    result = classify_port(port)
    print(f"Port {port}: {result}")
```

---

### Exercise 2.2: Loops - Doing Things Repeatedly

```python
#!/usr/bin/env python3
# =============================================================================
# LOOPS - REPEATING ACTIONS
# =============================================================================
#
# ELI5: Sometimes you need to do the same thing many times.
# Instead of writing the same code 100 times, you use a loop!
#

# -----------------------------------------------------------------------------
# FOR LOOP - Repeat for each item in a collection
# -----------------------------------------------------------------------------
#
# ELI5: "For each cookie in the jar, eat the cookie."
#

# Loop through a list
servers = ["web-01", "web-02", "db-01", "cache-01"]

print("--- Server List ---")
for server in servers:
    print(f"  • {server}")

# Loop through numbers using range()
#
# ELI5: range(5) gives you 0, 1, 2, 3, 4 (starts at 0, stops BEFORE 5)
#
print("\n--- Counting to 5 ---")
for i in range(5):
    print(f"Count: {i}")

# range() with start and end
print("\n--- Counting 5 to 10 ---")
for i in range(5, 11):  # 5 to 10 (stops before 11)
    print(f"Count: {i}")

# range() with step
print("\n--- Even numbers 0-10 ---")
for i in range(0, 11, 2):  # 0, 2, 4, 6, 8, 10
    print(f"Even: {i}")

# enumerate() - get index AND value
#
# ELI5: Sometimes you need to know which item you're on (1st, 2nd, 3rd...)
#
print("\n--- Numbered Server List ---")
for index, server in enumerate(servers, start=1):
    print(f"{index}. {server}")

# -----------------------------------------------------------------------------
# WHILE LOOP - Repeat while a condition is true
# -----------------------------------------------------------------------------
#
# ELI5: "While you're still hungry, eat another cookie."
# Be careful! If you never stop being hungry, you'll eat forever!
#

print("\n--- Countdown ---")
countdown = 5
while countdown > 0:
    print(f"{countdown}...")
    countdown = countdown - 1  # Or: countdown -= 1
print("Liftoff! 🚀")

# -----------------------------------------------------------------------------
# BREAK AND CONTINUE
# -----------------------------------------------------------------------------
#
# break: Stop the loop entirely (exit the loop)
# continue: Skip to the next iteration (skip rest of this round)
#

# Example: Find the first even number
print("\n--- Find First Even ---")
numbers = [1, 3, 7, 8, 9, 10]
for num in numbers:
    if num % 2 == 0:
        print(f"Found even number: {num}")
        break  # Stop looking, we found one!

# Example: Skip certain items
print("\n--- Skip Maintenance Servers ---")
servers = ["web-01", "web-02-maintenance", "db-01", "web-03-maintenance"]
for server in servers:
    if "maintenance" in server:
        continue  # Skip this one, move to next
    print(f"Active: {server}")

# -----------------------------------------------------------------------------
# PRACTICAL EXAMPLE: SIMPLE PORT SCANNER
# -----------------------------------------------------------------------------

def scan_ports(start: int, end: int, open_ports: list) -> list:
    """
    Simulate scanning ports.
    
    ELI5: We're checking doors (ports) to see which ones are open.
    In real life, we'd actually try to connect to each port.
    
    Args:
        start: First port to scan
        end: Last port to scan
        open_ports: List of ports that are "open" (for simulation)
    
    Returns:
        List of open ports found
    """
    found = []
    
    for port in range(start, end + 1):
        if port in open_ports:
            found.append(port)
            print(f"Port {port}: OPEN")
        # In a real scanner, you'd try to connect here
    
    return found


# Simulate: these ports are "open"
simulated_open = [22, 80, 443, 8080]
print("\n--- Port Scan Simulation (1-100) ---")
results = scan_ports(1, 100, simulated_open)
print(f"Open ports found: {results}")
```

---

### Exercise 2.3: Data Structures - Collections of Data

```python
#!/usr/bin/env python3
# =============================================================================
# DATA STRUCTURES - ORGANIZING DATA
# =============================================================================
#
# ELI5: So far we've stored single values in variables (one box).
# But what if you have lots of related things?
# Data structures are like filing cabinets, toolboxes, or shopping lists.
#

# -----------------------------------------------------------------------------
# LISTS - Ordered, changeable collections
# -----------------------------------------------------------------------------
#
# ELI5: A list is like a numbered to-do list or a playlist.
# Items have an order (1st, 2nd, 3rd) and you can change them.
#

# Creating lists
servers = ["web-01", "web-02", "db-01"]
numbers = [1, 2, 3, 4, 5]
mixed = [1, "hello", 3.14, True]  # Lists can hold different types

# Accessing items (zero-indexed!)
#
# ELI5: Python counts from 0, not 1.
# First item is index 0, second is index 1, etc.
#
print("--- List Access ---")
print(f"First server: {servers[0]}")   # web-01
print(f"Second server: {servers[1]}")  # web-02
print(f"Last server: {servers[-1]}")   # db-01 (negative = from end)

# Slicing - get a portion of the list
#
# ELI5: Like cutting a piece of a loaf of bread
#
print(f"First two: {servers[0:2]}")    # ['web-01', 'web-02']
print(f"From index 1: {servers[1:]}")  # ['web-02', 'db-01']

# Modifying lists
servers.append("cache-01")           # Add to end
servers.insert(0, "lb-01")          # Insert at position
removed = servers.pop()              # Remove and return last item
servers.remove("web-02")             # Remove specific item

print(f"Modified list: {servers}")

# List operations
print(f"Length: {len(servers)}")
print(f"Is 'db-01' in list? {'db-01' in servers}")

# List comprehension - create lists in one line
#
# ELI5: A shortcut for making lists. Like saying
# "Give me every number doubled" instead of using a loop.
#
numbers = [1, 2, 3, 4, 5]
doubled = [n * 2 for n in numbers]
print(f"Doubled: {doubled}")

# Filter with comprehension
evens = [n for n in numbers if n % 2 == 0]
print(f"Even numbers: {evens}")


# -----------------------------------------------------------------------------
# DICTIONARIES - Key-value pairs
# -----------------------------------------------------------------------------
#
# ELI5: A dictionary is like a real dictionary or phone book.
# You look something up by name (key) and get information (value).
#

# Creating dictionaries
server = {
    "hostname": "web-01",
    "ip": "10.0.1.10",
    "port": 80,
    "status": "running"
}

# Accessing values
print("\n--- Dictionary Access ---")
print(f"Hostname: {server['hostname']}")
print(f"IP: {server['ip']}")

# Safe access with .get() (returns None if key doesn't exist)
print(f"Location: {server.get('location', 'unknown')}")

# Modifying dictionaries
server['status'] = "stopped"         # Update existing
server['location'] = "us-east-1"    # Add new
del server['port']                   # Remove key

# Looping through dictionaries
print("\n--- Server Info ---")
for key, value in server.items():
    print(f"  {key}: {value}")

# Dictionary comprehension
ports = {"http": 80, "https": 443, "ssh": 22}
# Swap keys and values
ports_reversed = {v: k for k, v in ports.items()}
print(f"Reversed: {ports_reversed}")


# -----------------------------------------------------------------------------
# SETS - Unique, unordered collections
# -----------------------------------------------------------------------------
#
# ELI5: A set is like a bag of unique marbles.
# No duplicates allowed, and there's no specific order.
#

# Creating sets
allowed_ports = {22, 80, 443, 8080}
requested_ports = {22, 80, 3306, 5432}

print("\n--- Set Operations ---")
print(f"Allowed: {allowed_ports}")
print(f"Requested: {requested_ports}")

# Set operations (like math!)
print(f"Intersection (both): {allowed_ports & requested_ports}")
print(f"Union (either): {allowed_ports | requested_ports}")
print(f"Difference (allowed only): {allowed_ports - requested_ports}")
print(f"Difference (requested only): {requested_ports - allowed_ports}")

# Practical: Which requests are denied?
denied = requested_ports - allowed_ports
print(f"Denied ports: {denied}")


# -----------------------------------------------------------------------------
# TUPLES - Immutable (unchangeable) ordered collections
# -----------------------------------------------------------------------------
#
# ELI5: A tuple is like a sealed envelope.
# Once you put items in, you can't change them.
# Good for coordinates, database records, things that shouldn't change.
#

# Creating tuples
coordinates = (10.5, 20.3)
rgb_color = (255, 128, 0)  # Orange

# Accessing (same as lists)
print(f"\n--- Tuple Access ---")
print(f"X coordinate: {coordinates[0]}")

# Unpacking - assign multiple variables at once
#
# ELI5: Open the envelope and put each item in its own box
#
x, y = coordinates
print(f"x={x}, y={y}")

r, g, b = rgb_color
print(f"RGB: R={r}, G={g}, B={b}")

# Tuples as dictionary keys (lists can't do this!)
locations = {
    (40.7128, -74.0060): "New York",
    (34.0522, -118.2437): "Los Angeles",
}
print(f"Location at (40.7128, -74.0060): {locations[(40.7128, -74.0060)]}")
```

---

## Day 3: Functions and Modules

### Exercise 3.1: Functions - Reusable Code Blocks

```python
#!/usr/bin/env python3
# =============================================================================
# FUNCTIONS - REUSABLE CODE BLOCKS
# =============================================================================
#
# ELI5: A function is like a recipe card.
# You write the instructions once, then use them whenever you need.
# "make_cookies" might be called many times, but you only write it once.
#

# -----------------------------------------------------------------------------
# BASIC FUNCTIONS
# -----------------------------------------------------------------------------

# Simple function with no inputs
def say_hello():
    """
    Print a greeting.
    
    ELI5: This function does one thing - print hello.
    No ingredients needed.
    """
    print("Hello!")

say_hello()  # Call the function


# Function with parameters (inputs)
def greet(name):
    """
    Greet someone by name.
    
    Args:
        name: The person's name (the ingredient)
    """
    print(f"Hello, {name}!")

greet("Alice")
greet("Bob")


# Function with return value (output)
def add(a, b):
    """
    Add two numbers.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        The sum of a and b
    
    ELI5: This function takes two ingredients (numbers)
    and gives back a result (their sum).
    """
    return a + b

result = add(5, 3)
print(f"5 + 3 = {result}")


# -----------------------------------------------------------------------------
# DEFAULT PARAMETERS
# -----------------------------------------------------------------------------
#
# ELI5: Sometimes you want optional ingredients.
# If not specified, use a default value.
#

def ping(host, count=4, timeout=1.0):
    """
    Simulate a ping command.
    
    Args:
        host: Target to ping (required)
        count: Number of pings (optional, default 4)
        timeout: Timeout in seconds (optional, default 1.0)
    """
    print(f"Pinging {host}, count={count}, timeout={timeout}")

ping("google.com")              # Uses defaults
ping("google.com", count=10)    # Override count
ping("google.com", timeout=0.5) # Override timeout
ping("google.com", 10, 0.5)     # Override both


# -----------------------------------------------------------------------------
# TYPE HINTS
# -----------------------------------------------------------------------------
#
# ELI5: Type hints are like labels that say what type of ingredient
# a function expects. Python doesn't enforce them, but they help
# humans (and IDEs) understand the code.
#

def calculate_bmi(weight_kg: float, height_m: float) -> float:
    """
    Calculate Body Mass Index.
    
    Args:
        weight_kg: Weight in kilograms
        height_m: Height in meters
    
    Returns:
        BMI value
    """
    return weight_kg / (height_m ** 2)

bmi = calculate_bmi(70.0, 1.75)
print(f"BMI: {bmi:.1f}")


# -----------------------------------------------------------------------------
# *ARGS AND **KWARGS
# -----------------------------------------------------------------------------
#
# ELI5: 
# - *args: Accept any number of regular arguments (like a buffet)
# - **kwargs: Accept any number of named arguments (like custom toppings)
#

def make_sandwich(*ingredients):
    """Accept any number of ingredients."""
    print("Making sandwich with:")
    for item in ingredients:
        print(f"  - {item}")

make_sandwich("bread", "cheese", "tomato", "lettuce")


def configure_server(**options):
    """Accept any named options."""
    print("Server configuration:")
    for key, value in options.items():
        print(f"  {key}: {value}")

configure_server(host="0.0.0.0", port=8080, debug=True)


# -----------------------------------------------------------------------------
# LAMBDA FUNCTIONS (One-liners)
# -----------------------------------------------------------------------------
#
# ELI5: A lambda is a tiny, anonymous function.
# Use it when you need a simple function just once.
#

# Regular function
def double(x):
    return x * 2

# Same thing as a lambda
double_lambda = lambda x: x * 2

print(f"double(5) = {double(5)}")
print(f"double_lambda(5) = {double_lambda(5)}")

# Useful with sort, filter, map
numbers = [3, 1, 4, 1, 5, 9, 2, 6]
sorted_numbers = sorted(numbers, key=lambda x: -x)  # Sort descending
print(f"Sorted descending: {sorted_numbers}")


# -----------------------------------------------------------------------------
# PRACTICAL EXAMPLE: SUBNET CALCULATOR
# -----------------------------------------------------------------------------

from typing import Tuple, Dict


def parse_cidr(cidr: str) -> Tuple[str, int]:
    """
    Parse CIDR notation into IP and prefix length.
    
    ELI5: CIDR is like "123 Main Street /24" where /24 tells you
    how big the neighborhood is. We split it into address and size.
    
    Args:
        cidr: CIDR string like "192.168.1.0/24"
    
    Returns:
        Tuple of (ip_address, prefix_length)
    
    Raises:
        ValueError: If CIDR format is invalid
    """
    if "/" not in cidr:
        raise ValueError(f"Invalid CIDR format: {cidr}")
    
    ip, prefix_str = cidr.split("/")
    prefix = int(prefix_str)
    
    if not 0 <= prefix <= 32:
        raise ValueError(f"Invalid prefix length: {prefix}")
    
    return ip, prefix


def calculate_subnet_info(cidr: str) -> Dict:
    """
    Calculate subnet information from CIDR.
    
    Args:
        cidr: Network in CIDR notation
    
    Returns:
        Dictionary with subnet details
    """
    ip, prefix = parse_cidr(cidr)
    
    # Host bits = 32 - prefix
    host_bits = 32 - prefix
    
    # Total addresses = 2^host_bits
    total_addresses = 2 ** host_bits
    
    # Usable hosts = total - 2 (network address and broadcast)
    usable_hosts = max(0, total_addresses - 2)
    
    return {
        "network": ip,
        "prefix": prefix,
        "total_addresses": total_addresses,
        "usable_hosts": usable_hosts,
        "cidr": cidr
    }


# Test it
print("\n--- Subnet Calculator ---")
for cidr in ["10.0.0.0/8", "172.16.0.0/16", "192.168.1.0/24", "10.0.0.0/30"]:
    info = calculate_subnet_info(cidr)
    print(f"{cidr}: {info['usable_hosts']:,} usable hosts")
```

---

## Git Checkpoint

After completing Days 1-3, commit your progress:

```bash
cd ~/dev/python-learning

git add day_01/ day_02/ day_03/
git commit -m "Days 1-3: Python fundamentals complete

- Day 1: Variables, data types, strings, f-strings
- Day 2: Conditionals, loops, lists, dicts, sets, tuples
- Day 3: Functions, type hints, args/kwargs, lambdas

Built:
- Port classifier
- Subnet calculator
- Network utilities module"

git push origin main
```

---

## Summary: What You Learned

### Day 1
- ✅ Variables and data types (int, float, str, bool, None)
- ✅ Type conversion (int(), str(), float())
- ✅ Basic math operations (+, -, *, /, //, %, **)
- ✅ String operations and methods
- ✅ F-strings for formatting

### Day 2
- ✅ Conditionals (if, elif, else)
- ✅ Comparison operators (==, !=, <, >, <=, >=)
- ✅ Logical operators (and, or, not)
- ✅ For loops with range() and enumerate()
- ✅ While loops with break and continue
- ✅ Lists, dictionaries, sets, tuples
- ✅ List/dict comprehensions

### Day 3
- ✅ Defining functions with def
- ✅ Parameters and return values
- ✅ Default parameters
- ✅ Type hints
- ✅ *args and **kwargs
- ✅ Lambda functions
- ✅ Docstrings

---

## Next: Part 2 - Intermediate Python

Continue to `PYTHON_PART2_INTERMEDIATE.md` for:
- Classes and Object-Oriented Programming
- Error handling with try/except
- Decorators and generators
- Context managers (with statements)
