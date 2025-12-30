# Lab 1: Your First NixOS Installation & Configuration

## What You'll Learn

- How to install NixOS (and yes, you'll actually do it)
- Basic system configuration
- How to rebuild and apply changes
- How to troubleshoot when things go wrong
- The NixOS boot flow

## Estimated Time

- Clean install: 20-30 minutes
- Configuration exploration: 15-20 minutes
- Troubleshooting: varies (hopefully 0!)

## Prerequisites

- A spare machine, VM, or cloud instance
- Basic Linux knowledge (user accounts, packages, systemd)
- 2GB RAM minimum (4GB+ recommended)
- 5GB disk space minimum

## Lab Overview

```
Step 1: Boot NixOS installation media
   ↓
Step 2: Partition disk & install base system
   ↓
Step 3: Generate initial configuration
   ↓
Step 4: Customize configuration.nix
   ↓
Step 5: Rebuild and test
   ↓
Step 6: Make changes safely (with rollback)
   ↓
Step 7: Verify reproducibility
```

## Your First System

By the end of this lab, you'll have:
- ✅ Working NixOS installation
- ✅ Custom configuration.nix (you wrote it)
- ✅ Understanding of system rebuild process
- ✅ Knowledge of how to rollback
- ✅ A reproducible system you can redeploy

## Getting Started

### Option A: In a VM (Recommended for First Time)

```bash
# Using VirtualBox
1. Download NixOS ISO from https://nixos.org/download/
2. Create new VM
   - Name: nixos-lab-1
   - Memory: 4096 MB
   - Disk: 20 GB
3. Attach ISO and boot
4. Follow installation below
```

### Option B: Bare Metal

```bash
1. Download NixOS ISO
2. Write to USB: $ dd if=nixos-latest.iso of=/dev/sdX bs=4M
3. Boot from USB
4. Follow installation below
```

### Option C: Cloud (AWS/GCP/Digital Ocean)

Most cloud providers support bringing custom images. For learning, VirtualBox is easiest.

---

## Installation Steps

### Step 1: Boot into Live Environment

After booting NixOS ISO, you're in a live environment (nothing is installed yet).

```bash
# Verify you can see network
$ ping 8.8.8.8

# Verify disk is detected
$ lsblk

# You should see something like:
# NAME   MAJ:MIN RM SIZE RO TYPE MOUNTPOINT
# sda      8:0    0  20G  0 disk
```

### Step 2: Partition Your Disk

For this lab, simple partitioning:

```bash
# CHOOSE YOUR DISK carefully (be absolutely sure)
DISK=/dev/sda

# Option A: Interactive partitioning (safest)
$ parted $DISK
(parted) mklabel gpt
(parted) mkpart boot fat32 1MiB 512MiB
(parted) mkpart root ext4 512MiB 100%
(parted) set 1 esp on
(parted) quit

# Option B: Scripted (if confident)
$ sudo parted $DISK --script \
    mklabel gpt \
    mkpart boot fat32 1MiB 512MiB \
    mkpart root ext4 512MiB 100% \
    set 1 esp on

# Verify partitions were created
$ lsblk
# Should show: sda1 (512M) and sda2 (rest)
```

### Step 3: Format Partitions

```bash
# Format boot partition (EFI)
$ sudo mkfs.fat -F 32 /dev/sda1

# Format root partition
$ sudo mkfs.ext4 /dev/sda2

# Verify
$ lsblk -f | grep sda
```

### Step 4: Mount Filesystems

```bash
# Create mount points
$ sudo mkdir -p /mnt
$ sudo mkdir -p /mnt/boot

# Mount root first
$ sudo mount /dev/sda2 /mnt

# Mount EFI boot
$ sudo mount /dev/sda1 /mnt/boot

# Verify
$ df -h | grep mnt
# Should show both mounted
```

### Step 5: Generate Initial Configuration

NixOS has a tool that detects hardware and generates starter config:

```bash
# Generate configuration
$ sudo nixos-generate-config --root /mnt

# This creates:
# /mnt/etc/nixos/configuration.nix (hardware auto-detected)
# /mnt/etc/nixos/hardware-configuration.nix (keep this!)

# Verify files exist
$ sudo cat /mnt/etc/nixos/configuration.nix
$ sudo cat /mnt/etc/nixos/hardware-configuration.nix
```

### Step 6: Review Generated Configuration

Let's look at what was auto-generated:

```bash
$ sudo head -100 /mnt/etc/nixos/configuration.nix
```

You'll see comments explaining each section. Don't modify yet - just familiarize.

### Step 7: Customize Configuration

Now we'll modify it. Copy the template configuration from this lab:

```bash
# View the starter template provided in this lab
$ cat /path/to/lab-01/configuration.nix

# Copy to your system
$ sudo cp /path/to/lab-01/configuration.nix /mnt/etc/nixos/configuration.nix

# Important: Keep hardware-configuration.nix untouched!
```

### Step 8: Install System

```bash
# This builds and installs entire system to /mnt
# Takes 10-20 minutes (downloads/compiles packages)
$ sudo nixos-install

# If it completes successfully:
# "Done. Now you can reboot and boot into the new system"

# Set password for root
$ sudo passwd  # when prompted

# Verify (should show 1 generation)
$ ls -la /mnt/nix/var/nix/profiles/system-*
```

### Step 9: Reboot

```bash
$ sudo reboot

# Remove installation media when prompted
# System should boot into your new NixOS!
```

### Step 10: First Boot - Verify Installation

```bash
# Log in as root (with password you set)

# Check system is working
$ uname -a
$ nixos-version

# Check boot generation
$ sudo nixos-rebuild list-generations

# Verify packages from configuration are installed
$ which vim
$ which git
$ which htop
```

---

## Understanding What Just Happened

```
Traditional Linux Installation:
  1. Boot installer
  2. Manually partition
  3. Manually format
  4. Install core packages
  5. Install services separately
  6. Configure services separately
  7. Hope everything works together

NixOS Installation:
  1. Boot installer
  2. Partition & format
  3. $ nixos-generate-config (detects hardware)
  4. Edit configuration.nix (one file!)
  5. $ nixos-install (builds ENTIRE system from config)
  6. Reboot (everything configured from that one file)
```

The `nixos-install` command:
1. Reads your configuration.nix
2. Calculates all dependencies
3. Builds/fetches packages
4. Generates all config files
5. Creates systemd units
6. Installs grub bootloader
7. Creates a system closure (entire system as one unit)

Your system is now a "system closure" - everything needed to boot is recorded.

---

## Explore Your System

### See What's Installed

```bash
# View all installed packages
$ sudo nixos-rebuild dry-build

# See actual configuration being used
$ cat /etc/nixos/configuration.nix

# View generated systemd services
$ systemctl list-unit-files | grep enabled

# Check boot menu (your generations)
$ sudo grub-reboot 0  # Switch between generations
```

### The Secret Weapon: System Generations

```bash
# See all system generations
$ sudo nixos-rebuild list-generations
# Output:
#  1   2024-01-15 12:30:15
#  2   2024-01-15 12:45:32   <- current

# At any point, you can boot to a previous generation
# Go to GRUB menu (hold shift during boot)
# Select "NixOS - All configurations"
# Choose any previous generation
```

---

## Lab 2: Make Your First Change

Now let's modify the configuration and see the rebuild process.

### Change 1: Add a New Package

```bash
# Edit configuration
$ sudo vim /etc/nixos/configuration.nix

# Find this section:
# environment.systemPackages = with pkgs; [
#   git
#   vim
#   htop
# ];

# Add `curl`:
# environment.systemPackages = with pkgs; [
#   git
#   vim
#   htop
#   curl
# ];

# Save and exit
```

### Change 2: Rebuild System

```bash
# Dry run (see what would change, don't apply)
$ sudo nixos-rebuild dry-build

# Actually apply changes
$ sudo nixos-rebuild switch

# This:
# 1. Evaluates configuration.nix
# 2. Calculates what changed
# 3. Builds only new packages
# 4. Atomically switches to new system
# 5. Doesn't require reboot (unless kernel changed)

# Verify change
$ which curl
$ curl --version
```

### Change 3: Try a Rollback

```bash
# See your generations
$ sudo nixos-rebuild list-generations

# Switch back to previous
$ sudo nixos-rebuild switch --rollback

# Verify curl is gone
$ which curl
# Should say "not found"

# Switch forward again
$ sudo nixos-rebuild switch
$ which curl
# Should be back
```

This is the power of NixOS: **changes are experiments, rollback is instant**.

---

## Troubleshooting Common Issues

### Issue 1: "No space left on device"

```bash
# Problem: Package build used all disk space

# Solution: Delete old generations
$ sudo nix-collect-garbage

# Check disk
$ df -h /

# If still needed:
$ sudo nix-collect-garbage -d  # Delete ALL old generations
```

### Issue 2: "Error: attribute 'X' missing in 'configuration'"

```bash
# Problem: Typo in configuration.nix

# Solution:
$ cat /etc/nixos/configuration.nix | grep -n "error_line"
# Check syntax around that line
# Likely missing: { } : ; =
```

### Issue 3: "Build failed"

```bash
# Problem: Package can't be built

# Solution:
# 1. Check if package name is correct
# 2. Try dry-build to see full error
$ sudo nixos-rebuild dry-build -I nixpkgs=...

# 3. Search for package
$ nix-env -qaP "package_name"
```

### Issue 4: System Won't Boot

```bash
# This shouldn't happen with simple changes, but:

# At GRUB menu (hold shift during boot):
# Select NixOS boot submenu
# Boot previous generation

# Once booted:
$ sudo nixos-rebuild list-generations
$ sudo nixos-rebuild switch --rollback
```

---

## Verification Checklist

Before declaring victory, verify:

- [ ] System boots successfully
- [ ] Can login
- [ ] Packages from configuration are available (git, vim, htop)
- [ ] `sudo nixos-rebuild list-generations` shows generations
- [ ] Can add a package, rebuild, and verify it's installed
- [ ] Rollback works (previous generation still boots)
- [ ] Bootloader shows multiple options

---

## What You've Learned

✅ How NixOS installation differs from traditional Linux
✅ Partitioning and formatting
✅ The role of hardware-configuration.nix
✅ How configuration.nix describes entire system
✅ System rebuild process (fast, safe)
✅ System generations and atomic switching
✅ Rollback capability
✅ Safe experimentation with configurations

---

## Next Lab

Lab 2 will dive deeper into what you can configure and why NixOS's approach is powerful.

[Next: Lab 2 - Declaring Your Perfect Development Environment](../lab-02-dev-environment/README.md)

---

## Reference

- [NixOS Manual Installation](https://nixos.org/manual/nixos/stable/#sec-installation)
- [NixOS Configuration Reference](https://search.nixos.org/options/)
- [Hardware Configuration Deep Dive](./hardware-config.md)
- [Debugging Failed Rebuilds](../reference/debugging-nix.md)
