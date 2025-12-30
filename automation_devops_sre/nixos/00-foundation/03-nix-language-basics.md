# Foundation 3: Nix Language Basics

## You Don't Need to Learn Functional Programming

Many people worry: "Nix is a functional language, I don't know that."

Good news: You only need to understand 5 concepts to write most NixOS configurations.

```
Functional Programming (all the scary stuff):
  - Monads, Functors, Applicatives, Combinators...
  - [95% of theory you don't need]
  
Configuration Writing (what you actually need):
  - Let-in expressions (variables)
  - Attribute sets (organized data)
  - Lists (collections)
  - String interpolation (templates)
  - Function calls (using packages)
  - [5% you'll learn in 30 minutes]
```

## Core Concept 1: Everything is an Expression

In Nix, everything evaluates to a value. Nothing "happens" procedurally.

```nix
# NOT like bash (commands that do things)
# LIKE math (expressions that produce values)

# Expression: produces a number
1 + 2

# Expression: produces a string
"Hello " + "World"

# Expression: produces an attribute set (dict)
{ name = "Alice"; age = 30; }

# Expression: produces a list
[ 1 2 3 ]
```

**Key difference from bash**:
```bash
# Bash (imperative - do things):
name="Alice"
echo $name
```

```nix
# Nix (declarative - produce values):
let
  name = "Alice";
in
  name
# Produces the value: "Alice"
```

## Core Concept 2: Attribute Sets (Dicts/Maps)

Think of attribute sets like JSON objects. They're everywhere in NixOS.

```nix
# Simple attribute set
{
  name = "nginx";
  version = "1.24";
  enabled = true;
}

# Nested attribute sets
{
  user = {
    username = "alice";
    groups = [ "wheel" "docker" ];
  };
  system = {
    hostname = "server-01";
    timezone = "UTC";
  };
}

# Access values with dot notation
let
  config = {
    hostname = "server";
    port = 8080;
  };
in
  config.hostname  # Produces: "server"
  config.port      # Produces: 8080
```

This is exactly how NixOS configuration works:

```nix
# /etc/nixos/configuration.nix is an attribute set
{
  # All values are inside this set
  system.stateVersion = "23.11";
  networking.hostname = "myserver";
  services.nginx.enable = true;
}
```

## Core Concept 3: Let-In (Variables)

Use `let` to define temporary values, `in` to use them:

```nix
let
  # Define variables here
  myName = "Alice";
  myAge = 30;
  basePort = 8000;
  
in
  # Use them here
  {
    username = myName;
    age = myAge;
    port = basePort + 80;  # Arithmetic too
  }
# Produces: { username = "Alice"; age = 30; port = 8080; }
```

**Real-world NixOS example**:

```nix
let
  # Define once, use multiple times
  webDomain = "example.com";
  acmeMail = "admin@example.com";
  
in
{
  # Use in multiple places
  networking.domains = [ webDomain ];
  
  services.nginx.virtualHosts."${webDomain}" = {
    forceSSL = true;
    enableACME = true;
  };
  
  security.acme.certs."${webDomain}" = {
    email = acmeMail;
  };
}
```

## Core Concept 4: String Interpolation

Use `${}` to embed expressions inside strings:

```nix
# Simple interpolation
let
  name = "Alice";
in
  "Hello, ${name}!"
  # Produces: "Hello, Alice!"

# Expressions inside ${}
let
  port = 8080;
in
  "Server running on port ${toString port}"
  # Note: toString converts non-strings to strings
  # Produces: "Server running on port 8080"

# Complex interpolation
let
  user = { name = "bob"; uid = 1000; };
in
  "User ${user.name} has ID ${toString user.uid}"
  # Produces: "User bob has ID 1000"
```

## Core Concept 5: Lists

Simple ordered collections:

```nix
# List of strings
[ "git" "vim" "htop" ]

# List of numbers
[ 1 2 3 4 5 ]

# Mixed (less common)
[ 1 "two" 3 ]

# Nested
[ [ 1 2 ] [ 3 4 ] ]

# In configuration
{
  environment.systemPackages = with pkgs; [
    nginx
    postgresql
    git
  ];
  
  networking.firewall.allowedTCPPorts = [ 80 443 ];
}
```

## The Pattern You'll See Everywhere

```nix
{ config, pkgs, ... }:

{
  # Your configuration here
}
```

Let's break this down:

```nix
# Function declaration
{ config, pkgs, ... }:

# This means: "This file is a function that takes three things:"
#   - config: your current system configuration
#   - pkgs: all available packages
#   - ...: other stuff (you don't need to name it)

# The function returns:
{
  # An attribute set (your configuration)
}
```

**Real example**:

```nix
{ config, pkgs, ... }:

{
  # pkgs contains all packages
  environment.systemPackages = with pkgs; [
    nodejs    # from pkgs
    postgresql # from pkgs
  ];
  
  # config contains your current config
  networking.hostname = "server-01";
}
```

## Understanding `with`

`with` is a convenience shortcut:

```nix
# Without with:
environment.systemPackages = [
  pkgs.nodejs
  pkgs.postgresql
  pkgs.git
];

# With with (cleaner):
environment.systemPackages = with pkgs; [
  nodejs      # pkgs.nodejs implied
  postgresql  # pkgs.postgresql implied
  git         # pkgs.git implied
];
```

It says: "For the next expression, assume all values come from pkgs"

## Functions (You'll Use These)

Nix functions are expressions that produce values based on inputs:

```nix
# Simple function
let
  add = a: b: a + b;
in
  add 2 3
# Produces: 5

# Function that returns an attribute set
let
  makeUser = name: uid: {
    inherit name uid;  # shorthand for name = name; uid = uid;
    isNormalUser = true;
  };
in
  makeUser "alice" 1000
# Produces: { name = "alice"; uid = 1000; isNormalUser = true; }

# Function with default arguments
let
  makeService = { name, port ? 8000 }: {
    inherit name port;
    enabled = true;
  };
in
  [
    (makeService { name = "api"; })          # port defaults to 8000
    (makeService { name = "web"; port = 3000; })  # port is 3000
  ]
```

**In NixOS, many things are just function calls**:

```nix
{ config, pkgs, ... }:

{
  # Calling a function to generate nginx configuration
  services.nginx.virtualHosts."example.com" = {
    forceSSL = true;
    enableACME = true;
  };
  
  # Calling a function to create users
  users.users.alice = {
    isNormalUser = true;
    home = "/home/alice";
    shell = pkgs.bash;
  };
}
```

## Conditionals

Sometimes you need if-then-else:

```nix
# Simple conditional
let
  isProduction = true;
in
  if isProduction
  then { debug = false; logLevel = "error"; }
  else { debug = true; logLevel = "debug"; }

# In configuration
{ config, pkgs, ... }:

{
  services.nginx.enable = true;
  services.nginx.recommendedTLSSettings = true;
  
  # Conditional based on your settings
  services.nginx.virtualHosts."example.com" = {
    forceSSL = if config.services.ssl.enable then true else false;
  };
}
```

## Common Pitfalls

### Pitfall 1: Forgetting Spaces Around =
```nix
# WRONG:
{ name="alice"; }

# CORRECT:
{ name = "alice"; }
```

### Pitfall 2: Mixing Lists and Sets
```nix
# WRONG: These are different things
{
  packages = [ "git" "vim" ]  # List (ordered)
  settings = { ssl = true }   # Set (unordered)
}

# Lists and Sets are NOT interchangeable
# Lists: [ a b c ]
# Sets: { a = 1; b = 2; }
```

### Pitfall 3: Forgetting to Reference Variables
```nix
# WRONG: Port is a string "8080", not a number
let
  port = "8080";
in
  config.port = port + 1;  # Can't add string + number

# CORRECT: Use numbers or toString
let
  port = 8080;
in
  config.port = port + 1;  # Now it works: 8081
```

## The Structure of Every Configuration

```nix
# Every NixOS configuration follows this pattern:

{ config, pkgs, ... }:

let
  # (Optional) Define variables/helpers
  domain = "example.com";
  
in

{
  # System basics
  system.stateVersion = "23.11";
  
  # Services
  services.nginx.enable = true;
  
  # User-level packages
  environment.systemPackages = with pkgs; [
    git
    vim
  ];
  
  # Users
  users.users.alice = { /* config */ };
}
```

This pattern is the foundation for everything you'll write.

## Practice Exercises (Do These)

### Exercise 1: Attribute Set
```nix
# Create an attribute set for a web service
{
  name = "nginx";
  port = 80;
  ssl = false;
}

# Questions:
# 1. How do you access the port?
# 2. How do you change ssl to true?
```

### Exercise 2: String Interpolation
```nix
let
  hostname = "server-01";
  domain = "example.com";
in
  # Create a string: "server-01.example.com"
  # Using string interpolation
```

### Exercise 3: Let-In
```nix
let
  basePort = 8000;
  services = 3;
  portRange = 100;
in
{
  service1Port = basePort;
  service2Port = basePort + portRange;
  service3Port = basePort + (portRange * 2);
}

# Produces what?
```

### Answers (after you try):
```nix
# Ex 1
# 1. { name = "nginx"; port = 80; ssl = false; }.port  # produces: 80
# 2. { name = "nginx"; port = 80; ssl = true; }

# Ex 2
let
  hostname = "server-01";
  domain = "example.com";
in
  "${hostname}.${domain}"
# Produces: "server-01.example.com"

# Ex 3
{
  service1Port = 8000;
  service2Port = 8100;
  service3Port = 8200;
}
```

---

## You Now Know Enough Nix

You've learned the 5 concepts that cover ~95% of NixOS configuration writing:
1. ✅ Attribute sets (organized configuration)
2. ✅ Let-in (define variables)
3. ✅ String interpolation (templating)
4. ✅ Lists (collections)
5. ✅ Functions (using packages)

The rest you'll learn by example.

---

## Next: Your First Lab

Time to apply this knowledge to a real system. In Lab 1, you'll install NixOS and write your first configuration.

[Go to: Lab 1 - First Installation](../labs/lab-01-first-install/README.md)

Or continue learning theory:
- [Nix Language Reference (advanced)](https://nixos.org/manual/nix/stable/language/index.html)
- [Nixpkgs Manual](https://nixos.org/manual/nixpkgs/stable/)
