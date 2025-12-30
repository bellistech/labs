# Debugging Nix and NixOS - Troubleshooting Guide

## Quick Fix Index

If your error is about:
- **"attribute missing"** → [Syntax & Typo Errors](#syntax--typo-errors)
- **"build failed"** → [Build Failures](#build-failures)
- **"package not found"** → [Package Not Found](#package-not-found)
- **"service won't start"** → [Service Issues](#service-issues)
- **"can't write to disk"** → [Disk Issues](#disk-issues)
- **"can't find module"** → [Module Errors](#module-errors)

---

## Understanding Error Messages

### Anatomy of a Nix Error

```
error: attribute 'services.nging' missing in 'config'
       ↑              ↑                    ↑       ↑
    type        what's wrong          where    context
```

The error tells you:
1. **Type**: What kind of problem (attribute, syntax, evaluation)
2. **What's wrong**: The specific issue
3. **Where**: The section causing the problem
4. **Context**: File and line number (sometimes)

---

## Syntax & Typo Errors

### Error: "attribute 'X' missing"

```
error: attribute 'services.nging' missing in 'config'
```

**Cause**: Typo or wrong section name

```nix
# WRONG: (typo)
services.nging.enable = true;

# RIGHT:
services.nginx.enable = true;
```

**How to fix**:
```bash
# 1. Check for typos carefully
grep -n "nging" /etc/nixos/configuration.nix

# 2. Search for correct option
nix search nixpkgs "nginx" | head -20

# 3. Check documentation
visit search.nixos.org/options and search "nginx"
```

### Error: "expected attribute set"

```
error: expected attribute set
```

**Cause**: Using wrong syntax (list instead of set, or vice versa)

```nix
# WRONG: (attribute set syntax with list)
services.nginx = [ 
  enable = true;  # ✗ Can't do this
];

# RIGHT: (attribute set)
services.nginx = {
  enable = true;
};

# Or for lists:
environment.systemPackages = [
  pkgs.git
  pkgs.vim
];
```

**How to fix**:
```bash
# Remember:
# Sets: { key = value; }
# Lists: [ item1 item2 ]

# Check your syntax carefully
# { } = set (for configuration)
# [ ] = list (for collections)
```

### Error: "unexpected end of file"

```
error: unexpected end of file
```

**Cause**: Missing closing brace or bracket

```nix
# WRONG: (missing closing brace)
{
  networking.hostname = "server";
  services.nginx.enable = true;
# ← Missing }

# RIGHT:
{
  networking.hostname = "server";
  services.nginx.enable = true;
}
```

**How to fix**:
```bash
# Check ending of file
tail -5 /etc/nixos/configuration.nix

# Look for matching { }
grep -o "{" configuration.nix | wc -l
grep -o "}" configuration.nix | wc -l

# Should be equal!
```

### Error: "semicolon expected"

```
error: expected ';' at (line 5, column 42)
```

**Cause**: Missing semicolon

```nix
# WRONG:
services.nginx.enable = true  # ← missing ;

# RIGHT:
services.nginx.enable = true;
```

**How to fix**:
```bash
# Check line 5, column 42
sed -n '5p' /etc/nixos/configuration.nix | cut -c 42-50

# All value assignments need semicolons inside sets
# Exception: last item in a set (but it's good practice anyway)
```

---

## Build Failures

### Error: "hash mismatch"

```
error: hash mismatch in fixed-output derivation '/nix/store/...'
expected sha256:abc123...
got sha256:def456...
```

**Cause**: Package source has changed (upstream update or corruption)

**Solutions**:

```bash
# Solution 1: Use fetchurl with wrong hash first, copy correct hash
let
  src = fetchurl {
    url = "...";
    sha256 = lib.fakeSha256;  # Will fail and show real hash
  };
in

# Solution 2: Update nixpkgs (if it's a known issue)
sudo nixos-rebuild switch -I nixpkgs=<nixpkgs>

# Solution 3: Check if source URL is accessible
curl -I https://example.com/package.tar.gz
```

### Error: "trying to fetch..."

```
error: cannot download [package] from any source
```

**Cause**: Package source unavailable or network issue

**Solutions**:

```bash
# 1. Check network
ping 8.8.8.8

# 2. Check if URL is accessible
curl https://github.com/package/releases/download/...

# 3. Use binary cache instead of building
# (NixOS downloads pre-built binaries)
# Usually automatic, but can verify:
nix-channel --list
nix-channel --update

# 4. Wait and retry (temporary network issue)
sudo nixos-rebuild switch
```

### Error: "insufficient disk space"

```
error: cannot create directory: No space left on device
```

**Cause**: /nix/store has filled disk

**Solutions**:

```bash
# 1. Check disk usage
df -h

# 2. Clean old generations
sudo nix-collect-garbage

# 3. Aggressive cleanup (removes ALL old generations)
sudo nix-collect-garbage -d

# 4. Find large packages
du -sh /nix/store/* | sort -h | tail -20

# 5. Expand disk (if you can)
# Varies by system, would need to resize partitions
```

### Error: "infinite recursion"

```
error: infinite recursion encountered
```

**Cause**: Configuration refers to itself somehow

```nix
# WRONG: (refers to itself)
let
  x = x + 1;
in x

# RIGHT: (break the reference)
let
  x = 5;
  y = x + 1;
in y
```

**How to fix**:
```bash
# Look at line number in error
# Check for self-references or circular imports
grep -r "import.*self" /etc/nixos/

# Check for circular imports between files
# A imports B imports A
```

---

## Package Not Found

### Error: "attribute 'python311Packages.X' missing"

```
error: attribute 'python311Packages.nonexistent' missing in 'nixpkgs'
```

**Cause**: Package doesn't exist or has different name

**Solutions**:

```bash
# 1. Search for the package
nix search nixpkgs "package-name"

# 2. Common renames
# Old: python311Packages.requests
# New: python311Packages.requests  # still exists

# 3. Check if it exists for your Python version
nix search nixpkgs "requests" | grep python

# 4. Use pkgs search instead
nix search nixpkgs "name-substring"
```

### Error: "cannot find package in scope"

```
error: package 'some-tool' not found
```

**Cause**: Using tool name instead of package

```nix
# WRONG: (tool name, not package)
environment.systemPackages = with pkgs; [
  nginx  # ← Do you mean the server? Package name might be different
];

# RIGHT: (check actual package name)
environment.systemPackages = with pkgs; [
  nginx  # this is correct, but verify with search
];

# Find it:
nix search nixpkgs "nginx"
```

**How to fix**:
```bash
# 1. Search for what you want
nix search nixpkgs "web server" | grep -i nginx

# 2. Look at package attributes
nix search nixpkgs "nginx" -l | head -5

# 3. Use exact name from search results
```

---

## Service Issues

### Service Won't Start

```bash
# Check service status
systemctl status nginx

# Shows: Unit nginx.service failed to load

# Find error details
journalctl -u nginx -n 50
```

**Causes & Solutions**:

```bash
# 1. Configuration syntax error
sudo nixos-rebuild dry-build  # Catches config errors

# 2. Service name wrong
# Search for service name:
nix search nixos "nginx"

# 3. Service not available in nixpkgs
# Check what version/nixpkgs you're using
nix-channel --list

# 4. Service conflicts with running process
# Example: nginx is already running from manual install
which nginx
sudo systemctl stop nginx
sudo systemctl disable nginx

# 5. Port already in use
sudo lsof -i :80

# 6. Missing dependencies
# Check configuration for required settings
# Example: nginx needs working filesystem
```

### Debugging Service Startup

```bash
# 1. Check if service is enabled
systemctl is-enabled nginx

# 2. Check service file
cat /run/systemd/system/nginx.service

# 3. Watch service startup
journalctl -u nginx -f

# (In another terminal)
sudo systemctl restart nginx

# 4. Check service definition in configuration
grep -A 20 "services.nginx" /etc/nixos/configuration.nix
```

### Port Already in Use

```bash
# Problem: Service won't bind to port

# Find what's using the port
sudo lsof -i :80
sudo netstat -tuln | grep :80

# Solution 1: Stop the other service
sudo systemctl stop httpd

# Solution 2: Use different port
services.nginx.defaultListen = [
  { addr = "0.0.0.0"; port = 8080; }
];
```

---

## Module Errors

### Error: "module 'X' is not a module"

```
error: module 'hardware-configuration' is not a module
```

**Cause**: Trying to import something that isn't a NixOS module

```nix
# WRONG: (not a module format)
imports = [ ./something.nix ];
# where something.nix doesn't export a config attribute set

# RIGHT: (proper module)
imports = [ ./hardware-configuration.nix ];
# where hardware-configuration.nix is a proper module
```

**Module structure**:
```nix
# Correct module format:
{ config, pkgs, ... }:

{
  # Your config here
}

# Not just attribute set
```

### Error: "conflicting definitions"

```
error: The option 'services.X' has conflicting definitions
```

**Cause**: Same option defined in multiple places

```nix
# If using imports:
# A.nix: services.nginx.enable = true;
# B.nix: services.nginx.enable = true;
# Then importing both causes conflict

# Solution: Use lib.mkForce
services.nginx.enable = lib.mkForce true;

# Or define in only one place
```

---

## Disk & Storage Issues

### No Space Left on Device

```bash
# Check disk usage
df -h

# Clean nix store
sudo nix-collect-garbage
sudo nix-collect-garbage -d

# See what's taking space
du -sh /nix/store/* | sort -h | tail -20

# See what's in your home
du -sh ~/* | sort -h | tail -20
```

### Can't Write to /etc/nixos

```bash
# Problem: Permission denied when editing

# NixOS systems are read-only outside of configuration
# Solution: Edit as root

sudo vim /etc/nixos/configuration.nix

# If you're not in sudoers:
# Need root login to add yourself
su -
```

---

## Network Issues

### Can't Reach Internet

```bash
# Check network
ping 8.8.8.8

# Check DNS
nslookup example.com

# Check routes
ip route

# Check interfaces
ip addr

# If all else fails, check configuration
cat /etc/nixos/configuration.nix | grep networking
```

### SSH Connection Issues

```bash
# Check if SSH is enabled
systemctl status ssh

# Check if key is authorized
cat ~/.ssh/authorized_keys

# If not there, add it:
echo "ssh-rsa AAAA..." >> ~/.ssh/authorized_keys

# Verify SSH configuration
sudo grep -i "PermitRootLogin" /etc/ssh/sshd_config
```

---

## Getting Help

### When You're Really Stuck

```bash
# 1. Get the full error output
sudo nixos-rebuild switch 2>&1 | tee rebuild.log

# 2. Search NixOS manual
man nixos.conf  # if available
# Or: https://search.nixos.org/options

# 3. Search on GitHub issues
# https://github.com/NixOS/nixpkgs/issues

# 4. Ask on forums
# https://discourse.nixos.org/
# https://reddit.com/r/NixOS/

# 5. Check configuration option
nix search nixos "your-option"

# 6. Read example configurations
# https://github.com/NixOS/nixos-hardware/
# https://github.com/nix-community/home-manager/
```

### Creating a Minimal Example

When asking for help, reduce to minimum:

```nix
# Minimal example showing the problem
{ pkgs, ... }:

{
  system.stateVersion = "23.11";
  
  # Only the part that fails:
  services.nginx.enable = true;
  
  # This lets helpers reproduce quickly
}
```

---

## Debugging Workflow

### Step 1: Check for Syntax Errors

```bash
sudo nixos-rebuild dry-build

# Builds without applying
# If this passes, your syntax is correct
```

### Step 2: Check Service Status

```bash
sudo systemctl status <service>
journalctl -u <service> -n 50
```

### Step 3: Check Logs

```bash
# System logs
journalctl -n 100

# Service-specific
journalctl -u nginx -n 50

# Follow in real-time
journalctl -f
```

### Step 4: Verify Configuration

```bash
# What's actually running?
systemctl show <service>

# What's in the files?
cat /etc/nixos/configuration.nix | grep "your.option"
```

### Step 5: Test Changes

```bash
# Test without applying
sudo nixos-rebuild dry-build

# If satisfied:
sudo nixos-rebuild switch

# If something broke:
sudo nixos-rebuild switch --rollback
```

---

## Common Patterns That Fail

### Pattern 1: Forgetting `pkgs.`

```nix
# WRONG:
environment.systemPackages = [
  git  # ← Not found, need pkgs.git
  vim
];

# RIGHT:
environment.systemPackages = with pkgs; [
  git  # Now refers to pkgs.git
  vim
];

# OR:
environment.systemPackages = [
  pkgs.git
  pkgs.vim
];
```

### Pattern 2: Wrong Option Type

```nix
# WRONG: (trying to set a set to a string)
services.nginx.defaultServer = "example.com";

# RIGHT: (understand what the option expects)
services.nginx.defaultListenAddresses = [ "80" ];
```

**How to check**:
```bash
nix search nixos "services.nginx.default"
# Check the type shown
```

### Pattern 3: Circular Imports

```nix
# A.nix imports B.nix imports A.nix
# This creates infinite loop

# Solution: Use lib.mkMerge or other functions
# Better: Restructure to avoid circularity
```

---

## Reference Commands

```bash
# Build without applying
sudo nixos-rebuild dry-build

# Apply changes
sudo nixos-rebuild switch

# Rollback last change
sudo nixos-rebuild switch --rollback

# See all generations
sudo nixos-rebuild list-generations

# See what would change
sudo nixos-rebuild dry-activate

# View current configuration
cat /etc/nixos/configuration.nix

# Search for option
nix search nixos "what-you-want"

# See service status
systemctl status <service>

# View service logs
journalctl -u <service>

# View system logs
journalctl -n <number>

# Reload systemd (after config change)
sudo systemctl daemon-reload
```

---

## Prevention

### Best Practices to Avoid Issues

1. **Always test before deploying**: `dry-build` is your friend
2. **Use version control**: Track all changes in git
3. **Read error messages carefully**: They usually tell you exactly what's wrong
4. **Search first**: Most issues have been solved
5. **Small incremental changes**: Don't change 10 things at once
6. **Document your configuration**: Comments help future you

---

## Quick Reference Card

```
Error Type              First Check         Next Check
─────────────────────────────────────────────────────────
Syntax error            Rebuild dry-build   Search option docs
Package not found       nix search nixpkgs  Check package name
Service won't start     systemctl status    journalctl output
No space left           df -h               nix-collect-garbage
Network issue           ping 8.8.8.8        ip route
SSH problem             systemctl status    ~/.ssh/authorized_keys
Permission denied       sudo command        Check sudoers
Port in use             lsof -i :PORT       Change port number
```

---

For more help:
- [Official Manual](https://nixos.org/manual/nixos/stable/)
- [Option Search](https://search.nixos.org/options)
- [Package Search](https://search.nixos.org/packages)
- [NixOS Forum](https://discourse.nixos.org/)
- [Reddit Community](https://reddit.com/r/NixOS/)
