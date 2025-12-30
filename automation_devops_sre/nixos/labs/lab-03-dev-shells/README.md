# Lab 3: Creating Reproducible Development Shells

## What You'll Learn

- Quick ad-hoc development shells (no files needed)
- Nix-shell one-liners for temporary environments
- Creating reusable shell templates
- Language-specific development environments
- Mixing and matching tools on-the-fly

## Estimated Time

- Quick shells: 20 minutes
- Templates: 25 minutes
- Custom combinations: 25 minutes
- Complete lab: 70 minutes

## Prerequisites

- Completed Lab 1 (NixOS installed)
- Completed Lab 2 (Understand shell.nix files)

## Why This Matters

Lab 2 was about permanent project environments. This lab is about:
- **Quick experimentation** ("Let me test this tool")
- **One-off tasks** ("I need a PostgreSQL client right now")
- **Template reuse** ("Every Python project needs these")
- **No file creation** ("I just want to try something")

---

## Part 1: Ad-Hoc Shells (Command Line)

### The Pattern

```bash
# Basic pattern: nix-shell -p package1 package2 package3

# Example: One-time Python environment
nix-shell -p python311

# You're now in Python environment (no files created)
$ python --version
Python 3.11.0

$ exit
```

### One-Liners for Common Tasks

#### Python Experimentation

```bash
# Quick Python with common packages
nix-shell -p python311 python311Packages.numpy python311Packages.matplotlib

# Python with all scientific packages
nix-shell -p python311Full

# Python with data science stack
nix-shell -p python311 \
  python311Packages.numpy \
  python311Packages.pandas \
  python311Packages.matplotlib \
  python311Packages.jupyter
```

#### Node.js Quick Test

```bash
# Latest Node.js with npm
nix-shell -p nodejs

# Specific version
nix-shell -p nodejs_18

# Node with package managers
nix-shell -p nodejs yarn
```

#### Database Tools

```bash
# PostgreSQL client tools only (no server)
nix-shell -p postgresql_15

# Connect to database
nix-shell -p postgresql_15
$ psql -h localhost -U user -d dbname
```

#### Rust Development

```bash
# Full Rust environment
nix-shell -p rustup

# Or specific Rust version
nix-shell -p rust cargo rustfmt clippy
```

#### Go Development

```bash
# Go compiler and tools
nix-shell -p go

# With common tools
nix-shell -p go gopls
```

### Real-World Usage

```bash
# Scenario: You want to try a tool for 10 minutes
# Don't want to install it permanently

# Option 1: Install globally (pollutes system)
# Option 2: Use Docker (heavy)
# Option 3: nix-shell (clean!)

nix-shell -p ffmpeg
$ ffmpeg -version
$ ffmpeg -i video.mp4 -c:v libx264 output.mp4
$ exit
# ffmpeg is gone, system clean

# Try another tool
nix-shell -p imagemagick
$ convert image.jpg -resize 50% thumb.jpg
$ exit
# Both tools work perfectly, nothing left behind
```

---

## Part 2: Shell with Custom Script

### Pattern: Bash Script That Enters Shell

Useful for development workflows that need multiple steps:

```bash
# Create dev-shell.sh
cat > dev-shell.sh << 'EOF'
#!/bin/bash

# Development shell with automatic setup

# Enter nix-shell with specific packages
exec nix-shell -p \
  nodejs_18 \
  postgresql_15 \
  redis \
  git \
  --run "bash"

EOF

chmod +x dev-shell.sh

# Use it
./dev-shell.sh
# Now in environment with all those tools
```

### Pattern: Shell with initialization

```bash
cat > dev-init.sh << 'EOF'
#!/bin/bash

# Complex shell with auto-setup

nix-shell -p nodejs_18 postgresql_15 --run bash -c '
  echo "=== Development Environment ==="
  echo "Node: $(node --version)"
  echo "npm: $(npm --version)"
  echo "PostgreSQL: $(psql --version)"
  echo ""
  echo "Checking local development database..."
  
  if pg_isready -h localhost -p 5432; then
    echo "✓ PostgreSQL is running"
  else
    echo "✗ PostgreSQL not running (start it if needed)"
  fi
  
  echo ""
  echo "Ready to develop!"
  
  # Drop into bash
  bash
'
EOF

chmod +x dev-init.sh
./dev-init.sh
```

---

## Part 3: Reusable Shell Templates

### Create a Templates Directory

```bash
# Organize shell templates by language
mkdir -p ~/.nix-shells

# Now create templates for common stacks
```

### Template 1: Python Data Science

```bash
cat > ~/.nix-shells/python-datascience.nix << 'EOF'
# Data science environment
#
# Use: nix-shell ~/.nix-shells/python-datascience.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    python311
    python311Packages.numpy
    python311Packages.pandas
    python311Packages.matplotlib
    python311Packages.scikit-learn
    python311Packages.jupyter
    python311Packages.ipython
    
    # Tools
    git
    curl
    jq
  ];
  
  shellHook = ''
    echo "Data Science Environment"
    echo "Available: Python, NumPy, Pandas, Matplotlib, scikit-learn, Jupyter"
    echo "Start Jupyter: jupyter notebook"
  '';
}
EOF

# Use it
nix-shell ~/.nix-shells/python-datascience.nix
```

### Template 2: Full Stack (Node + Python + DB)

```bash
cat > ~/.nix-shells/fullstack.nix << 'EOF'
# Full-stack development environment
#
# Use: nix-shell ~/.nix-shells/fullstack.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # Frontend
    nodejs_18
    yarn
    
    # Backend
    python311
    python311.pkgs.pip
    python311.pkgs.virtualenv
    
    # Databases
    postgresql_15
    redis
    
    # Tools
    git
    curl
    docker
    
    # Utilities
    tmux
    htop
    jq
  ];
  
  shellHook = ''
    export VENV_DIR=".venv"
    if [ ! -d "$VENV_DIR" ]; then
      python -m venv $VENV_DIR
    fi
    source $VENV_DIR/bin/activate
    
    echo "Full-Stack Environment Ready"
    echo "Frontend: Node $(node --v), npm, yarn"
    echo "Backend: Python $(python --version 2>&1)"
    echo "DB: PostgreSQL, Redis"
  '';
}
EOF

# Use it
nix-shell ~/.nix-shells/fullstack.nix
```

### Template 3: DevOps Tools

```bash
cat > ~/.nix-shells/devops.nix << 'EOF'
# DevOps environment
#
# Use: nix-shell ~/.nix-shells/devops.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # Infrastructure as Code
    terraform
    ansible
    
    # Container tools
    docker
    docker-compose
    
    # Kubernetes
    kubectl
    helm
    
    # Cloud CLI
    awscli2
    google-cloud-sdk
    
    # Configuration management
    git
    
    # Utilities
    tmux
    htop
    jq
    yq
    curl
    wget
    openssh
    
    # Debugging
    netcat
    tcpdump
    strace
  ];
  
  shellHook = ''
    echo "DevOps Environment Ready"
    echo "Available:"
    echo "  IaC: terraform, ansible"
    echo "  Container: docker, docker-compose"
    echo "  Kubernetes: kubectl, helm"
    echo "  Cloud: aws, gcloud"
  '';
}
EOF

# Use it
nix-shell ~/.nix-shells/devops.nix
```

### Template 4: Rust Development

```bash
cat > ~/.nix-shells/rust.nix << 'EOF'
# Rust development environment
#
# Use: nix-shell ~/.nix-shells/rust.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    rustup
    rust-analyzer
    cargo
    rustfmt
    clippy
    
    # Dependencies
    openssl
    pkg-config
    
    # Tools
    git
    curl
    
    # Optional: for WebAssembly
    # wasm-pack
  ];
  
  shellHook = ''
    echo "Rust Environment Ready"
    echo "rustup: $(rustup --version)"
    echo "cargo: $(cargo --version)"
    echo ""
    echo "Create new project: cargo new project_name"
  '';
}
EOF

# Use it
nix-shell ~/.nix-shells/rust.nix
```

### How to Use Templates

```bash
# List your templates
ls ~/.nix-shells/

# Enter data science environment
nix-shell ~/.nix-shells/python-datascience.nix

# Enter devops environment
nix-shell ~/.nix-shells/devops.nix

# Create alias for quick access (add to ~/.bashrc or ~/.zshrc)
alias ds='nix-shell ~/.nix-shells/python-datascience.nix'
alias fs='nix-shell ~/.nix-shells/fullstack.nix'
alias devops='nix-shell ~/.nix-shells/devops.nix'
alias rust='nix-shell ~/.nix-shells/rust.nix'

# Now just type:
ds
# [Python data science environment]

devops
# [DevOps environment]
```

---

## Part 4: Combining Templates

### Pattern: Combining Multiple Templates

```bash
# Create a super environment combining multiple
cat > ~/.nix-shells/super.nix << 'EOF'
# Super environment combining everything
#
# Use: nix-shell ~/.nix-shells/super.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # Development
    nodejs_18 python311 rustup go
    
    # Databases
    postgresql_15 redis mongodb
    
    # DevOps
    terraform docker kubectl
    
    # Tools
    git curl jq tmux
  ];
  
  shellHook = ''
    echo "SuperEnvironment - Everything available!"
  '';
}
EOF

nix-shell ~/.nix-shells/super.nix
```

### Pattern: Environment Variables

```bash
cat > ~/.nix-shells/custom.nix << 'EOF'
# Custom environment with variables
#
# Use: nix-shell ~/.nix-shells/custom.nix

{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    nodejs_18
    python311
  ];
  
  # Set custom environment variables
  shellHook = ''
    export MY_PROJECT_ROOT="$PWD"
    export API_URL="http://localhost:8000"
    export DB_HOST="localhost"
    export DEBUG=1
    
    echo "Environment variables set:"
    echo "  MY_PROJECT_ROOT=$MY_PROJECT_ROOT"
    echo "  API_URL=$API_URL"
    echo "  DB_HOST=$DB_HOST"
    echo "  DEBUG=$DEBUG"
  '';
}
EOF
```

---

## Part 5: Practical Workflows

### Workflow 1: Quick Python Script Testing

```bash
# You have a quick Python script to test
cat > script.py << 'EOF'
import json
import requests

data = {"name": "test"}
print(json.dumps(data))
EOF

# Enter Python environment with required packages
nix-shell -p python311 python311Packages.requests

# Run the script
python script.py

# Exit when done
exit
```

### Workflow 2: Multiple Temporary Projects

```bash
# Project 1 - Node testing
mkdir project-node
cd project-node
nix-shell -p nodejs
npm init
npm install express
npm list
exit

# Project 2 - Python testing (same machine, different tools!)
mkdir project-python
cd project-python
nix-shell -p python311
python -m pip list
exit

# Project 1 again
cd project-node
nix-shell -p nodejs
npm list
exit

# No conflicts! Each shell has exactly what it needs
```

### Workflow 3: Trying Tools Before Committing

```bash
# Option A: Try a tool for 5 minutes
nix-shell -p postgresql_15
$ psql --version
$ createdb testdb
$ exit
# Database tools gone, system clean

# Option B: If I like it, commit it to a project
cat >> shell.nix << 'EOF'
postgresql_15
EOF

# Now it's permanent for this project
```

---

## Advanced: Creating Custom Shell Functions

### Add to ~/.bashrc or ~/.zshrc

```bash
# Instant development environment by language

# Quick Python
py() {
  nix-shell -p python311 python311Packages.pip --run bash
}

# Quick Node
node-dev() {
  nix-shell -p nodejs yarn --run bash
}

# Quick Go
go-dev() {
  nix-shell -p go --run bash
}

# Quick Rust
rust-dev() {
  nix-shell ~/.nix-shells/rust.nix
}

# Usage:
# $ py
# $ node-dev
# $ go-dev
# $ rust-dev
```

---

## Troubleshooting

### Issue 1: "Package not found"

```bash
# If: nix-shell -p pythonwrong
# Error: attribute 'pythonwrong' missing

# Solution: Search for correct name
nix search nixpkgs python | grep -i data

# Common corrections:
# python → python311
# node → nodejs_18
# pg → postgresql_15
```

### Issue 2: "Different packages have conflicts"

```bash
# This shouldn't happen in nix-shell
# Each package is isolated
# But if issues:

# Solution 1: Use only compatible versions
# Solution 2: Use separate shells for incompatible tools
nix-shell -p tool1
# later
nix-shell -p tool2

# Not: nix-shell -p tool1 tool2 (if incompatible)
```

### Issue 3: "Environment variable not set"

```bash
# Problem: Expected variable not in shell
# Solution: Add to shellHook

cat > shell.nix << 'EOF'
shellHook = ''
  export MY_VAR="value"
'';
EOF
```

---

## Verification Checklist

- [ ] Created ad-hoc shell with `nix-shell -p`
- [ ] Multiple packages in one shell
- [ ] Exited shell and verified tools gone
- [ ] Created shell template file
- [ ] Used template with `nix-shell path/to/template.nix`
- [ ] Understand difference: temp shell vs project shell
- [ ] Can quickly create environment for any language
- [ ] Know where to search for package names
- [ ] Tested combining multiple templates
- [ ] Set environment variables in shellHook

---

## What You've Learned

✅ Ad-hoc shells with `-p` flag
✅ One-liners for quick environments
✅ Reusable shell templates
✅ Language-specific environments
✅ Combining templates
✅ No polluting system with global installs
✅ Quick switching between environments

---

## Real-World Pattern

```bash
# Morning: I need to work with multiple languages today

# Python task
nix-shell ~/.nix-shells/python-datascience.nix
[work with data]
exit

# DevOps task
nix-shell ~/.nix-shells/devops.nix
[configure infrastructure]
exit

# Node project
cd node-project
nix-shell  # uses project's shell.nix
[develop]
exit

# Back to devops
nix-shell ~/.nix-shells/devops.nix
[more infrastructure]
exit

# System state: CLEAN. Nothing installed globally.
# All work is reproducible and isolated.
```

---

## Next: Multi-Host Configuration

You now understand how developers stay in isolated environments.

**Lab 4** will scale this to managing multiple complete systems:
- Configuring 10 servers identically
- Sharing configuration between systems
- NixOS modules for composability

[Next: Lab 4 - Multi-Host Configuration Management](../lab-04-multi-host/README.md)

---

## Reference

- [Nix-shell Manual](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html)
- [Package Search](https://search.nixos.org/packages)
- [Direnv](https://direnv.net/) - Automatic shell on directory change
- [Nix Pills - Working with local directories](https://nixos.org/guides/nix-pills/developing-with-nix-shell.html)
