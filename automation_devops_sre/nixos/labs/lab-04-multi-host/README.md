# Lab 4: Multi-Host Configuration Management

## What You'll Learn

- Declaring multiple NixOS systems in one place
- Sharing configurations between systems
- System-specific overrides
- Deploying identical configs to multiple machines
- Managing infrastructure as code

## Estimated Time

- Understanding concepts: 25 minutes
- Creating multi-host setup: 35 minutes
- Deploying to multiple systems: 40 minutes
- Complete lab: 100 minutes

## Prerequisites

- Completed Lab 1 (NixOS installation)
- Comfortable with configuration.nix
- Basic understanding of shared vs specific configs

## Why This Matters

```
Traditional Multi-Server Management:
  Server 1: Manual configuration
  Server 2: Follow wiki but forget a step
  Server 3: Different version of something
  
  Result: 3 different systems
  Problem: Debugging 3 different issues

NixOS Multi-Host Management:
  Define 3 servers in ONE file
  $ deploy to server-1
  $ deploy to server-2
  $ deploy to server-3
  
  Result: 3 identical systems
  Benefit: Same config = same behavior
```

---

## Concept: Shared vs Specific

### The Pattern

```nix
# Common configuration (all servers use this)
let
  commonConfig = {
    hostname = "base";  # Will be overridden
    networking.firewall.enable = true;
    services.openssh.enable = true;
  };
  
  # Server-specific config overrides common
  server1 = commonConfig // {
    networking.hostname = "webserver-01";
  };
  
  # Server 2 overrides different things
  server2 = commonConfig // {
    networking.hostname = "dbserver-01";
    services.postgresql.enable = true;
  };
```

### Real Workflow

```
You define:
  ├─ Base configuration (everything in common)
  ├─ Server A (base + nginx)
  ├─ Server B (base + postgresql)
  └─ Server C (base + redis)

Each server gets deployed with its specific config
All 3 stay in sync when you update base config
```

---

## Part 1: Simple Two-Server Setup

### Step 1: Create Project Structure

```bash
mkdir nixos-fleet
cd nixos-fleet

# Create directory for configurations
mkdir machines

# Create files
touch machines/base.nix          # Shared config
touch machines/web-server.nix    # Server 1
touch machines/db-server.nix     # Server 2
```

### Step 2: Create Base Configuration

```bash
cat > machines/base.nix << 'EOF'
# Base configuration shared by all servers
#
# This contains everything common to all machines in this fleet.
# Server-specific configs import this and add their own customizations.

{ config, pkgs, ... }:

{
  ############################################################################
  # SYSTEM BASICS (same for all servers)
  ############################################################################
  
  system.stateVersion = "23.11";
  
  # Time synchronization (all servers should have correct time)
  services.chrony.enable = true;
  
  ############################################################################
  # NETWORKING (common firewall rules)
  ############################################################################
  
  networking.firewall.enable = true;
  networking.firewall.allowedTCPPorts = [
    22    # SSH (always enabled)
    80    # HTTP
    443   # HTTPS
  ];
  
  ############################################################################
  # SSH (all servers need SSH for administration)
  ############################################################################
  
  services.openssh.enable = true;
  services.openssh.permitRootLogin = "no";  # More secure
  services.openssh.passwordAuthentication = false;
  # SSH keys will be added by specific configs
  
  ############################################################################
  # USERS (base users, all servers)
  ############################################################################
  
  # Don't allow arbitrary user changes
  users.mutableUsers = true;
  
  # Note: specific servers will add their own users/keys
  
  ############################################################################
  # SYSTEM PACKAGES (useful on all servers)
  ############################################################################
  
  environment.systemPackages = with pkgs; [
    git
    curl
    htop
    tmux
    
    # Essential monitoring
    vim
    jq
  ];
  
  ############################################################################
  # LOGGING (centralized setup)
  ############################################################################
  
  # Later: could add centralized logging
  # For now, just journalctl (built-in)
  
  ############################################################################
  # SECURITY
  ############################################################################
  
  security.sudo.enable = true;
  
  # Fail2ban protects against brute force
  # (optional, more advanced)
  # services.fail2ban.enable = true;
  
}
EOF
```

### Step 3: Create Web Server Configuration

```bash
cat > machines/web-server.nix << 'EOF'
# Web Server Configuration
#
# This server runs:
# - Nginx (reverse proxy)
# - Runs on port 80/443
#
# Imports common configuration and adds web-specific parts

{ config, pkgs, ... }:

{
  # Import base configuration first
  imports = [ ./base.nix ];
  
  ############################################################################
  # SYSTEM IDENTIFICATION
  ############################################################################
  
  networking.hostname = "webserver-01";
  
  # IPv4 address (example, adjust to your network)
  networking.interfaces.eth0.ipv4.addresses = [
    {
      address = "192.168.1.10";
      prefixLength = 24;
    }
  ];
  
  ############################################################################
  # SSH KEYS for this server
  ############################################################################
  
  users.users.root.openssh.authorizedKeys.keys = [
    # Add your public SSH key here
    # "ssh-rsa AAAAB3NzaC1... your-key@yourmachine"
  ];
  
  ############################################################################
  # NGINX (Web Server)
  ############################################################################
  
  services.nginx.enable = true;
  
  # Configure virtual host
  services.nginx.virtualHosts."example.com" = {
    forceSSL = false;  # For testing, disable SSL
    enableACME = false;
    
    # Simple reverse proxy to backend
    locations."/" = {
      proxyPass = "http://127.0.0.1:3000";
      proxyWebsockets = true;
    };
  };
  
  ############################################################################
  # FIREWALL - Web server specific ports
  ############################################################################
  
  # Base config allows 80/443, that's what we need
  # If you need more ports, add them here:
  # networking.firewall.allowedTCPPorts = [ 3000 ];  # Backend port
  
}
EOF
```

### Step 4: Create Database Server Configuration

```bash
cat > machines/db-server.nix << 'EOF'
# Database Server Configuration
#
# This server runs:
# - PostgreSQL (database)
# - Redis (cache)
#
# No web server, no public ports

{ config, pkgs, ... }:

{
  # Import base configuration
  imports = [ ./base.nix ];
  
  ############################################################################
  # SYSTEM IDENTIFICATION
  ############################################################################
  
  networking.hostname = "dbserver-01";
  
  # Different IP address
  networking.interfaces.eth0.ipv4.addresses = [
    {
      address = "192.168.1.20";
      prefixLength = 24;
    }
  ];
  
  ############################################################################
  # SSH KEYS
  ############################################################################
  
  users.users.root.openssh.authorizedKeys.keys = [
    # Add your public SSH key here
  ];
  
  ############################################################################
  # POSTGRESQL (Database)
  ############################################################################
  
  services.postgresql.enable = true;
  services.postgresql.package = pkgs.postgresql_15;
  
  # Listen on localhost only (connect from web server via network)
  services.postgresql.settings.listen_addresses = "localhost";
  
  # Create databases/users that backend needs
  # (You'd normally do this with SQL, but can define here too)
  
  # Example backup configuration
  services.postgresql.backups.enable = true;
  services.postgresql.backups.location = "/var/backups/postgresql";
  
  ############################################################################
  # REDIS (Cache/Session Store)
  ############################################################################
  
  services.redis.servers.default = {
    enable = true;
    port = 6379;
    bind = "127.0.0.1";  # Localhost only
  };
  
  ############################################################################
  # FIREWALL - Database server specific
  ############################################################################
  
  # Database server doesn't need HTTP/HTTPS
  # Override firewall to be more restrictive
  networking.firewall.allowedTCPPorts = [
    22    # SSH only
  ];
  
  # Allow connections from web server (if on same network)
  # This is simplified - in production, more specific rules
  # networking.firewall.allowedTCPPorts = [ 5432 6379 ];  # PostgreSQL + Redis
  
  # Better: restrict by source IP
  # networking.firewall.extraCommands = ''
  #   iptables -A INPUT -s 192.168.1.10 -p tcp --dport 5432 -j ACCEPT
  #   iptables -A INPUT -s 192.168.1.10 -p tcp --dport 6379 -j ACCEPT
  # '';
  
}
EOF
```

---

## Part 2: Deploying to Real Systems

### Step 1: SSH Access to Servers

You need SSH access to your servers:

```bash
# Test SSH connection
ssh root@192.168.1.10  # Web server
ssh root@192.168.1.20  # Database server

# If you can connect, good! If not, fix SSH first
```

### Step 2: Building and Copying Config

```bash
# Build configuration for web server
nix-build -A system -I nixpkgs=~/nixpkgs machines/web-server.nix

# This produces a result in ./result that contains the entire system

# Copy to server (method 1: simple)
scp -r result root@192.168.1.10:/tmp/nixos-system

# Then on the server:
ssh root@192.168.1.10
# Inside server:
/tmp/nixos-system/bin/switch-to-configuration switch
```

### Step 3: Better Way - Using nixos-rebuild Over SSH

```bash
# Even better: nixos-rebuild can push to remote servers

# Copy config to server
scp machines/web-server.nix root@192.168.1.10:/etc/nixos/configuration.nix

# SSH into server
ssh root@192.168.1.10

# Inside the server:
sudo nixos-rebuild switch

# Done! System is updated
```

### Step 4: One Command Deploy

Create a deployment script:

```bash
cat > deploy.sh << 'EOF'
#!/bin/bash

# Deploy NixOS configurations to multiple servers

set -e  # Exit on error

echo "NixOS Multi-Host Deployer"
echo "========================="

# Function to deploy to one server
deploy_server() {
  local name=$1
  local ip=$2
  local config=$3
  
  echo ""
  echo "Deploying to $name ($ip)..."
  echo "Config: $config"
  
  # Copy configuration
  scp machines/$config root@$ip:/etc/nixos/configuration.nix
  
  # Rebuild on remote system
  ssh root@$ip 'echo "Building configuration..." && \
    sudo nixos-rebuild dry-build && \
    echo "Switching to new configuration..." && \
    sudo nixos-rebuild switch'
  
  if [ $? -eq 0 ]; then
    echo "✓ $name deployed successfully"
  else
    echo "✗ $name deployment failed"
    return 1
  fi
}

# Deploy all servers
deploy_server "webserver-01" "192.168.1.10" "web-server.nix"
deploy_server "dbserver-01" "192.168.1.20" "db-server.nix"

echo ""
echo "========================="
echo "Deployment complete!"

EOF

chmod +x deploy.sh

# Use it:
./deploy.sh
```

---

## Part 3: Scaling to Many Servers

### Create a Servers Registry

```bash
cat > machines/registry.nix << 'EOF'
# Server Registry
#
# Define all your servers in one place
# Reference from deployment script

{
  webserver-01 = {
    ip = "192.168.1.10";
    config = ./web-server.nix;
    description = "Web server (Nginx)";
  };
  
  webserver-02 = {
    ip = "192.168.1.11";
    config = ./web-server.nix;  # Same config as webserver-01
    description = "Web server (Nginx) - Backup";
  };
  
  dbserver-01 = {
    ip = "192.168.1.20";
    config = ./db-server.nix;
    description = "Database server (PostgreSQL)";
  };
  
  cacheserver-01 = {
    ip = "192.168.1.30";
    config = ./cache-server.nix;
    description = "Cache server (Redis)";
  };
}
EOF
```

### Advanced Deployment Script

```bash
cat > deploy-all.sh << 'EOF'
#!/bin/bash

# Deploy all servers from registry

REGISTRY="machines/registry.nix"

# Parse registry and deploy (advanced shell scripting)
# For production, consider using Terraform or NixOps

echo "Available servers:"
grep -E '^\s+[a-z-]+\s*=\s*{' $REGISTRY | sed 's/[{}=]//g'

read -p "Enter server names to deploy (space-separated) or 'all': " SERVERS

if [ "$SERVERS" = "all" ]; then
  SERVERS="webserver-01 webserver-02 dbserver-01 cacheserver-01"
fi

for SERVER in $SERVERS; do
  echo "Deploying $SERVER..."
  # Implementation: read registry, deploy
done

EOF

chmod +x deploy-all.sh
```

---

## Part 4: Real-World Multi-Server Example

### Scenario: Web Application Architecture

```
┌─────────────────┐
│  Web Server     │
│  (nginx)        │
│  Port 80/443    │
└────────┬────────┘
         │
         └─────────────┬─────────────┐
                       │             │
                   ┌───▼─────┐   ┌──▼───────┐
                   │Database │   │  Cache   │
                   │(PG 15)  │   │ (Redis)  │
                   └─────────┘   └──────────┘
```

### Web Server Config

```bash
cat > machines/web-app.nix << 'EOF'
{ config, pkgs, ... }:

{
  imports = [ ./base.nix ];
  
  networking.hostname = "webapp-01";
  
  # Nginx reverse proxy to backend app
  services.nginx.enable = true;
  services.nginx.virtualHosts."myapp.example.com" = {
    locations."/" = {
      proxyPass = "http://127.0.0.1:3000";
    };
  };
  
  # Node.js application
  # (Could also use systemd service to run Node app)
  
}
EOF
```

### Database Server Config

```bash
cat > machines/database.nix << 'EOF'
{ config, pkgs, ... }:

{
  imports = [ ./base.nix ];
  
  networking.hostname = "database-01";
  
  services.postgresql.enable = true;
  services.postgresql.package = pkgs.postgresql_15;
  
  # Backup strategy
  services.postgresql.backups.enable = true;
  
  # Create initial database
  systemd.services.postgresql.postStart = ''
    ''${pkgs.postgresql_15}/bin/psql -U postgres -c \
      "CREATE DATABASE myapp;" 2>/dev/null || true
  '';
}
EOF
```

### Cache Server Config

```bash
cat > machines/cache.nix << 'EOF'
{ config, pkgs, ... }:

{
  imports = [ ./base.nix ];
  
  networking.hostname = "cache-01";
  
  services.redis.servers.default = {
    enable = true;
    port = 6379;
  };
  
  # Monitor redis
  services.redis.servers.default.bind = "0.0.0.0";  # Accessible to other servers
  
}
EOF
```

---

## Troubleshooting

### Issue 1: Configuration won't build

```bash
# Error: attribute missing, syntax error, etc.

# Solution: Check syntax
nix-env -f machines/web-server.nix -i

# Check for errors before deploying
sudo nixos-rebuild dry-build
```

### Issue 2: SSH key issues

```bash
# Can't SSH to server

# Solution: Check SSH keys
# Add public key to server manually first:
ssh root@server "echo 'ssh-rsa AAAA...' >> ~/.ssh/authorized_keys"

# Or disable password auth only after key is added
```

### Issue 3: Firewall blocking

```bash
# Can't connect to service

# Solution: Check firewall rules
systemctl status firewall
sudo iptables -L

# Temporarily disable for testing:
sudo systemctl stop firewall
```

---

## Verification Checklist

- [ ] Created base.nix with common config
- [ ] Created web-server.nix importing base
- [ ] Created db-server.nix importing base
- [ ] Can build configurations locally (dry-run)
- [ ] Can SSH to at least one test server
- [ ] Successfully deployed to test server
- [ ] Verified deployment worked (services running)
- [ ] Made a change in base.nix, deployed to all servers
- [ ] All servers updated identically
- [ ] Understand how to add new servers (copy pattern)

---

## What You've Learned

✅ Shared vs server-specific configuration
✅ Using `imports` to share configuration
✅ Managing multiple servers from one place
✅ Deployment scripts and automation
✅ Server registry organization
✅ Infrastructure as code patterns

---

## Production Considerations

### What You Should Add

```bash
# 1. Centralized logging
services.syslog.enable = true;

# 2. Monitoring
services.prometheus.enable = true;
services.grafana.enable = true;

# 3. Backup strategy
services.backup.enable = true;

# 4. VPN/private network
# (for inter-server communication)

# 5. Secrets management
# (Nix can manage this, but requires careful setup)
```

### Tools to Graduate To

- **NixOps**: Declarative deployment tool for NixOS
- **Terraform**: For hybrid clouds + NixOS
- **Colmena**: Newer, simpler deployment tool
- **Flakes**: For reproducible, pinned versions

---

## Next Lab

You've mastered managing multiple systems identically.

**Lab 5** goes deeper into production hardening:
- Security best practices
- Service configuration
- Database setup
- Monitoring and logging

[Next: Lab 5 - Production Server Setup](../lab-05-server-setup/README.md)

---

## Reference

- [NixOS Multi-Machine Setup](https://nixos.wiki/wiki/Network)
- [NixOps Documentation](https://nixops.readthedocs.io/)
- [Colmena Deployment Tool](https://colmena.cli.rs/)
- [Managing Secrets in NixOS](https://nixos.wiki/wiki/Secrets)
