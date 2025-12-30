# Lab 2: Declaring Your Perfect Development Environment

## What You'll Learn

- Creating reproducible development environments with NixOS
- Using `nix-shell` for per-project isolation
- Managing language-specific toolchains
- Team consistency in development setups
- Switching between projects without conflicts

## Estimated Time

- Basic understanding: 20 minutes
- Creating your first dev environment: 30 minutes
- Multi-language setup: 30 minutes
- Complete lab: 90 minutes

## Why This Matters

```
Traditional Development Environment:
  Developer A: Uses Node 18.5.0, gets bugs
  Developer B: Has Node 20.0.0 globally, everything works
  Developer C: Uses NVM with 16.0.0 in old shell
  
  Result: "It works on my machine" nightmare
  
NixOS Development Environment:
  All developers: $ nix-shell
  [Automatically enters exact same environment]
  Same Node version, same npm version, same everything
  Leave project: $ exit
  [Node automatically "uninstalled"]
```

## Core Concept: The Development Shell

In NixOS, instead of installing dev tools globally (or using version managers like nvm, rbenv, pyenv), you declare exactly what each project needs:

```nix
# project/shell.nix
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    nodejs_18
    postgresql_15
    redis
  ];
}
```

Then:
```bash
$ cd project
$ nix-shell
# You're now in an environment with exactly those tools
# Nothing is installed globally
# When you exit, everything is "uninstalled"
```

---

## Part 1: Simple Node.js Project

### Step 1: Create a Project Directory

```bash
mkdir my-node-project
cd my-node-project
```

### Step 2: Create the shell.nix File

```bash
cat > shell.nix << 'EOF'
# Development environment for Node.js project
#
# This file declares exactly what tools this project needs.
# When you run `nix-shell`, NixOS reads this file and creates
# an isolated environment with exactly these packages.
#
# Key benefit: Your project is never "installed" on your system.
# It's only available when you nix-shell into the project.

{ pkgs ? import <nixpkgs> {} }:

# mkShell creates a development environment
pkgs.mkShell {
  
  # buildInputs: packages that should be available in the shell
  buildInputs = with pkgs; [
    # Node.js and npm (together in nodejs package)
    nodejs_18
    
    # Useful tools for Node development
    git
    curl
    jq  # JSON query tool
  ];
  
  # Optional: Set environment variables when entering shell
  shellHook = ''
    echo "================================"
    echo "Node.js Development Environment"
    echo "================================"
    echo "Node version: $(node --version)"
    echo "npm version: $(npm --version)"
    echo ""
    echo "Available tools:"
    echo "  - node, npm, npx"
    echo "  - git, curl, jq"
    echo ""
    echo "To exit this environment: exit"
    echo "================================"
  '';
  
}
EOF
```

### Step 3: Enter the Development Shell

```bash
# First time: might need to accept some Nix paths
nix-shell

# You should see:
# ================================
# Node.js Development Environment
# ================================
# Node version: v18.x.x
# npm version: x.x.x
# ...

# Verify you have Node.js
$ node --version
v18.x.x

$ npm --version
x.x.x

# Create a simple package.json
$ npm init -y
```

### Step 4: Create a Sample App

```bash
# Inside nix-shell
cat > index.js << 'EOF'
#!/usr/bin/env node

// Simple Node.js app to verify environment

const http = require('http');

const server = http.createServer((req, res) => {
  res.writeHead(200, {'Content-Type': 'text/plain'});
  res.end('Hello from NixOS dev environment!\n');
});

const PORT = 3000;
server.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}/`);
});
EOF

# Run it
$ node index.js
# Should say: "Server running at http://localhost:3000/"
# Stop with Ctrl+C

# Exit the nix-shell
$ exit
```

### Step 5: Verify Isolation

```bash
# You're now back in your normal shell
# Node.js should NOT be available

$ node --version
command not found: node

# But you can re-enter anytime
$ nix-shell
$ node --version
v18.x.x

$ exit
```

**This is the power**: Node.js is only available inside the nix-shell. Your system stays clean.

---

## Part 2: Multi-Language Project

### Scenario: Full-Stack App

```
Backend: Python with FastAPI
Frontend: Node.js with React
Database: PostgreSQL
Cache: Redis
```

### Step 1: Create Complex shell.nix

```bash
cat > shell.nix << 'EOF'
# Full-stack development environment
#
# This project needs:
# - Python 3.11 with specific packages
# - Node.js 18 for frontend
# - PostgreSQL client tools (for connecting to dev DB)
# - Redis tools (for debugging)
# - Git and other utilities

{ pkgs ? import <nixpkgs> {} }:

let
  # We can define versions at the top for easy maintenance
  pythonVersion = pkgs.python311;
  nodeVersion = pkgs.nodejs_18;
  postgresVersion = pkgs.postgresql_15;
  
in

pkgs.mkShell {
  
  buildInputs = with pkgs; [
    # Python + common packages for FastAPI development
    pythonVersion
    pythonVersion.pkgs.pip        # pip for installing Python packages
    pythonVersion.pkgs.virtualenv # for virtual environments
    
    # Python packages we need (directly, not via pip)
    # (Most should be in requirements.txt, but some we want in shell)
    
    # Node.js for frontend
    nodeVersion
    
    # Database tools
    postgresVersion  # Includes psql, createdb, etc.
    
    # Redis tools
    redis
    
    # General utilities
    git
    curl
    jq
    htop  # System monitor
    
    # For monitoring/debugging
    tmux  # Terminal multiplexer
  ];
  
  # Environment variables for development
  shellHook = ''
    # Set Python to use UTF-8
    export PYTHONIOENCODING=utf-8
    
    # Create/activate Python virtual environment
    export VENV_DIR=".venv"
    if [ ! -d "$VENV_DIR" ]; then
      echo "Creating Python virtual environment..."
      python -m venv $VENV_DIR
    fi
    source $VENV_DIR/bin/activate
    
    # Display helpful information
    echo "========================================="
    echo "Full-Stack Development Environment"
    echo "========================================="
    echo ""
    echo "Languages:"
    echo "  Python: $(python --version)"
    echo "  Node.js: $(node --version)"
    echo "  npm: $(npm --version)"
    echo ""
    echo "Databases:"
    echo "  PostgreSQL client: $(psql --version)"
    echo "  Redis: $(redis-cli --version)"
    echo ""
    echo "Project Setup:"
    echo "  Backend: pip install -r requirements.txt"
    echo "  Frontend: npm install"
    echo ""
    echo "To exit: exit"
    echo "========================================="
  '';
  
}
EOF
```

### Step 2: Set Up Project Files

```bash
# Backend setup
mkdir -p backend
cat > backend/requirements.txt << 'EOF'
fastapi==0.104.0
uvicorn==0.24.0
sqlalchemy==2.0.0
psycopg2-binary==2.9.0
redis==5.0.0
EOF

cat > backend/main.py << 'EOF'
from fastapi import FastAPI
from fastapi.responses import JSONResponse

app = FastAPI()

@app.get("/")
async def root():
    return JSONResponse({"message": "Hello from FastAPI"})

@app.get("/health")
async def health():
    return JSONResponse({"status": "ok"})

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
EOF

# Frontend setup
mkdir -p frontend
cat > frontend/package.json << 'EOF'
{
  "name": "frontend",
  "version": "0.1.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  }
}
EOF
```

### Step 3: Use the Environment

```bash
# Enter the shell (creates virtualenv if needed)
nix-shell

# Install backend dependencies
cd backend
pip install -r requirements.txt

# Start backend
python main.py
# Should see: "Uvicorn running on http://127.0.0.1:8000"

# In another terminal, enter nix-shell again
nix-shell
cd frontend
npm install
npm list

# Test the API
curl http://localhost:8000/
# Should respond: {"message":"Hello from FastAPI"}

# Test database connection (if you have a dev DB)
psql -h localhost -U user -d mydb -c "SELECT 1"

# Exit when done
exit
```

---

## Part 3: Team Consistency

### The Git-Based Workflow

```bash
# Team workflow:
# Developer A clones repo
git clone https://github.com/team/project.git
cd project

# Developer A doesn't need to follow setup wiki
# Just one command:
nix-shell

# [Environment is automatically identical to everyone else]

# Developer B on different machine:
git clone https://github.com/team/project.git
cd project
nix-shell

# [Same environment, guaranteed]
```

### Adding Project-Specific Documentation

```bash
# Create project README
cat > README.md << 'EOF'
# My Project

## Quick Start

```bash
# Enter development environment
nix-shell

# Install dependencies
pip install -r requirements.txt
npm install

# Run backend
python backend/main.py

# Run frontend (in another nix-shell)
nix-shell
npm start
```

## Environment

This project uses:
- Python 3.11 + FastAPI
- Node.js 18 + React
- PostgreSQL 15
- Redis

All managed by NixOS via `shell.nix`. No manual setup needed!

## Switching Projects

```bash
cd other-project
nix-shell
# [Now in OTHER project's environment]
```

No conflicts. Each project has exactly what it needs.
EOF
```

---

## Part 4: Advanced Patterns

### Pattern 1: Language-Specific Shells

Sometimes you want different environments:

```bash
# Create shell-python.nix for Python-only work
cat > shell-python.nix << 'EOF'
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    python311
    python311.pkgs.pip
    python311.pkgs.black     # Code formatter
    python311.pkgs.pytest    # Testing
    python311.pkgs.mypy      # Type checking
  ];
}
EOF

# Create shell-node.nix for Node-only work
cat > shell-node.nix << 'EOF'
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    nodejs_18
    yarn
  ];
}
EOF

# Use specific shell
nix-shell shell-python.nix
# or
nix-shell shell-node.nix
```

### Pattern 2: Direnv Integration (Automatic Entry)

Direnv automatically enters nix-shell when you cd into a directory:

```bash
# Enable direnv (one-time)
# Requires direnv package in your system

# In project directory, create .envrc
cat > .envrc << 'EOF'
use nix
EOF

# Authorize it (security feature)
direnv allow

# Now whenever you cd into project:
cd my-project
# [Automatically enters nix-shell]
# [Name shows as "(nix-shell)" in prompt]

cd ..
# [Automatically exits nix-shell]
```

### Pattern 3: Pinned Nixpkgs Version

For reproducibility across months/years:

```bash
cat > shell.nix << 'EOF'
{ pkgs ? import (fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/23.11.tar.gz";
    sha256 = "0dig043qz7yzmq41qw3m0r8nw5xl73pzcz69l1665yghfqkbk5cx";
  }) {}
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    nodejs_18
    postgresql_15
  ];
}
EOF
```

This guarantees exact versions forever (even if nixpkgs updates).

---

## Troubleshooting

### Issue 1: "nix-shell not found"

```bash
# Problem: Nix isn't installed on system
# Solution: Install Nix first
# https://nixos.org/download.html
```

### Issue 2: "Python virtual environment won't activate"

```bash
# The shellHook in shell.nix should handle it
# But if it doesn't:
source .venv/bin/activate  # Manual activation

# Or recreate it
rm -rf .venv
nix-shell  # Re-enter, it'll recreate
```

### Issue 3: "Package version is wrong"

```bash
# Verify what's in shell
nix-shell
$ node --version

# If wrong, check shell.nix for correct name
# Search for package: nix search nixpkgs nodejs

# Common mistake: using "nodejs" instead of "nodejs_18"
# or "nodejs_18" when you meant "nodejs_20"
```

### Issue 4: "Module not found" errors in Python

```bash
# If you have a .venv from outside NixOS:
rm -rf .venv

# Re-enter nix-shell (will create fresh virtualenv)
nix-shell

# Reinstall: pip install -r requirements.txt
```

---

## Verification Checklist

Before moving to next lab, verify:

- [ ] Created shell.nix file
- [ ] Entered nix-shell successfully
- [ ] Tools were available (node, python, etc.)
- [ ] Created simple project files
- [ ] Exited nix-shell, tools unavailable
- [ ] Re-entered nix-shell, tools available again
- [ ] Multi-language example works
- [ ] Virtual environment creates automatically
- [ ] Understand isolation benefit (clean system)
- [ ] Can explain to teammate how it works

---

## What You've Learned

✅ Creating development shells with `shell.nix`
✅ Per-project environment isolation
✅ Multi-language setups (Python, Node.js, DB tools)
✅ Virtual environment automation
✅ Team consistency (same environment everywhere)
✅ Clean system (nothing installed globally)
✅ Quick project switching

---

## Real-World Usage

### Daily Workflow

```bash
# Morning
cd project-a
nix-shell          # Automatically in project A environment
npm start          # Frontend
python app.py      # Backend (in another terminal's nix-shell)

# Midday - switch projects
exit               # Leave project A
cd project-b
nix-shell          # Now in project B environment (different versions!)

# Afternoon - back to project A
exit
cd project-a
nix-shell          # Same tools as this morning (reproducible)
```

### Onboarding New Developers

Old way:
```
1. Install Node 18.5.0
2. Install Python 3.11
3. Install PostgreSQL 15
4. Configure database
5. Set environment variables
6. npm install
7. pip install -r requirements.txt
8. Hope nothing conflicts
```

New way:
```
1. nix-shell
2. pip install -r requirements.txt
3. npm install
4. npm start
```

Done.

---

## Next Steps

You now understand development environment isolation. 

**Lab 3** will dive deeper with:
- Using nix-shell for quick environments (without creating files)
- Language-specific templates
- Home Manager for personal dotfiles

[Next: Lab 3 - Creating Reproducible Development Shells](../lab-03-dev-shells/README.md)

---

## Reference

- [Nix Shell Documentation](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html)
- [Direnv Homepage](https://direnv.net/)
- [Nixpkgs Package Search](https://search.nixos.org/packages)
- [Creating Development Environments](https://nixos.wiki/wiki/Development_environment_with_nix-shell)
