# Foundation 2: What is NixOS? The Operating System

## From Package Manager to Entire System

You've learned that Nix is a package manager that solves dependency hell. Now imagine applying that same thinking to your **entire operating system**.

Traditional Linux:
- Kernel somewhere
- Bootloader configured one way
- System services scattered in various configs
- Users manually tweaking files
- "It works until I update something"

NixOS:
- Define: "I want Linux kernel 6.1, Systemd, these services, this user"
- Build: Atomically construct entire system
- Deploy: Instant rollback if anything breaks
- Update: Safe, testable changes

**One configuration file describes your entire system.**

## The Core Idea: One File to Rule Them All

```
Traditional system management:
/etc/nginx/nginx.conf       <- manage separately
/etc/postgresql/pg.conf     <- manage separately
/etc/systemd/system/*.service <- manage separately
/etc/fstab                  <- manage separately
/etc/hosts                  <- manage separately
$ systemctl enable nginx    <- run manually
$ useradd myuser            <- run manually
# Now try to replicate on another server. Good luck.

NixOS:
/etc/nixos/configuration.nix
├── "kernel = linux_6_1"
├── "services.nginx.enable = true"
├── "services.postgresql.enable = true"
├── "users.users.myuser = { ... }"
├── networking configuration
├── filesystems
└── everything else
$ nixos-rebuild switch
# Done. All 10 servers configured identically.
```

## The Breakthrough: Declarative System Management

```
Imperative (traditional Linux administration):
  1. Install nginx
  2. Copy this config file
  3. Enable the service
  4. Restart systemd
  5. Update firewall rules
  6. Check if it's running
  7. Oh wait, I forgot to set the environment variable...

Declarative (NixOS):
  1. "Here's exactly what I want my system to look like"
  2. Tell NixOS to make it so
  3. Done.
```

NixOS reads your description and:
- Installs exact packages
- Generates config files based on your settings
- Enables services
- Sets environment variables
- Creates users
- Configures filesystems
- **All automatically**

## What Happens When You Deploy

```
You: $ nixos-rebuild switch

NixOS does:
  1. Evaluates /etc/nixos/configuration.nix
  2. Builds entire system derivation
     ├─ Fetches Linux kernel 6.1
     ├─ Compiles it with exact options
     ├─ Generates nginx.conf from your settings
     ├─ Generates systemd units
     ├─ Creates user accounts
     └─ Does thousands of other things
  3. Tests the build
  4. If any part fails: STOP (system unchanged)
  5. If all succeeds: Atomic switch
     └─ Bootloader points to new system
     └─ Next boot uses new configuration
  6. Previous system generation saved
     └─ Can rollback instantly
```

## The Nix Store Scales: System Generations

Just like packages in the Nix store have unique hashes, **entire systems** get versioned:

```
Nix boot menu (after several rebuilds):
┌─────────────────────────────────┐
│ NixOS Boot Menu                 │
├─────────────────────────────────┤
│ [*] Generation 45 (current)     │
│ [ ] Generation 44               │
│ [ ] Generation 43               │
│ [ ] Generation 42               │
│ [ ] Generation 41               │
│ [ ] Generation 40               │
│ [ ] Generation 39               │
└─────────────────────────────────┘

Problem: Broke everything with latest change
Solution: Select Generation 44, reboot, working again
         $ nixos-rebuild switch --rollback
```

Each generation is:
- Immutable (never changes)
- Complete (entire system state)
- Independent (can boot any one)
- Stored forever (you choose when to delete)

## NixOS = Nix Package Manager + SystemD + GRUB + Your Config

Conceptually:
```
         Your configuration.nix
                  |
                  v
         ┌────────────────────┐
         │  Nix Evaluator     │
         │ (reads your config)│
         └────────────────────┘
                  |
                  v
         ┌────────────────────┐
         │ Dependency Graph   │
         │ (what needs what)  │
         └────────────────────┘
                  |
                  v
         ┌────────────────────┐
         │ Build Derivation   │
         │ (compile/generate) │
         └────────────────────┘
                  |
                  v
         ┌────────────────────┐
         │ System Profile     │
         │ (ready to boot)    │
         └────────────────────┘
                  |
                  v
         ┌────────────────────┐
         │ GRUB Bootloader    │
         │ (points to it)     │
         └────────────────────┘
```

## What You Declare (Simple Example)

```nix
# /etc/nixos/configuration.nix

{ config, pkgs, ... }:

{
  # System basics
  boot.loader.grub.device = "/dev/sda";
  networking.hostname = "webserver-01";
  
  # Services
  services.nginx.enable = true;
  services.postgresql.enable = true;
  
  # User
  users.users.alice = {
    isNormalUser = true;
    group = "users";
    extraGroups = [ "wheel" "postgres" ];
  };
  
  # Packages available to all users
  environment.systemPackages = with pkgs; [
    git
    vim
    htop
  ];
  
  # Firewall
  networking.firewall.allowedTCPPorts = [ 80 443 ];
}
```

What NixOS does with this:
- Installs nginx (but only nginx, nothing extra)
- Generates `/etc/nginx/nginx.conf` from sensible defaults + your settings
- Creates systemd service for nginx
- Enables it to start on boot
- Same for postgresql
- Creates user `alice` with exact permissions
- Installs git, vim, htop system-wide
- Configures firewall rules
- **Builds one atomic system that includes all of this**

## Key NixOS Principles

### 1. Entire System is Reproducible
```
Same configuration.nix on Server A and Server B
   ↓
Identical systems
   ↓
No "it works on mine but not yours"
```

### 2. Configuration is Version-Controlled
```
Your production system: one git repository
/etc/nixos/configuration.nix

Want to see what changed? $ git log -p configuration.nix
Want to know what version of nginx was running in 2023? $ git blame
Want to understand a production decision? $ git log
```

### 3. Services Are Modular
```nix
# Want to add PostgreSQL? One line:
services.postgresql.enable = true;

# Want to add Let's Encrypt SSL? One line:
security.acme.enable = true;

# Want to configure both to work together?
# NixOS handles it automatically (they know about each other)
```

### 4. Immutability Prevents Configuration Drift

```
Traditional server after 6 months of manual changes:
  - Someone modified nginx.conf
  - Someone added a cron job
  - Someone updated a library by hand
  - Nobody remembers why service Y is installed
  - Original configuration is lost
  Result: Unmaintainable mess

NixOS server after 6 months:
  - Every change in version control
  - Exact state reproducible from git
  - Configuration drift impossible
  - Can rebuild from scratch anytime
```

### 5. Rollback Means Never Choosing Between Safety and Updates

```
Traditional thinking:
  "Should I update? No, might break everything"
  [system gets stale, security vulnerabilities]

NixOS thinking:
  "Should I update? Yes, if anything breaks I rollback in 10 seconds"
  [always current, always safe]
```

## Real-World Examples

### Example 1: Spin Up Identical Development Environment

```
Your team: 8 developers, 3 different laptops, 2 operating systems
Problem: "Works on my machine"
         "I have different Node version"
         "My dependencies are different"

Solution: One configuration.nix
  $ git clone repo
  $ nix flake update
  $ direnv allow  # automatic
  [enters environment with exact same versions]
  [everyone developing identical setup]
```

### Example 2: Staging Matches Production

```
You want to test nginx config changes safely
Traditional: Edit on staging server, hope it's like production
NixOS: 
  - Update configuration.nix
  - Run on staging ($ nixos-rebuild switch --dry-run to preview)
  - Test thoroughly
  - Deploy to production (same configuration.nix)
  - Identical systems, no surprises
```

### Example 3: Team Onboarding

```
New engineer joins
Traditional:
  - "Read the wiki"
  - Manual steps to set up dev environment
  - "Oh, you need version X of tool Y"
  - "Did you set environment variable Z?"
  - Weeks to get productive

NixOS:
  - $ git clone
  - $ nixos-rebuild switch
  - Same environment as all 7 other engineers
  - Ready to commit in minutes
```

## NixOS vs Traditional Linux Distros

| Aspect | Ubuntu/Fedora/Debian | NixOS |
|--------|---------------------|-------|
| **System config** | Manual files | Declarative code |
| **Reproducibility** | "Usually works" | 100% guaranteed |
| **Package conflicts** | Can happen | Impossible |
| **Update safety** | Risky | Safe (rollback) |
| **Version management** | Single system version | Multiple generations |
| **Documentation** | Per-package | Integrated in config |
| **Learning curve** | Easy | Medium (but worth it) |
| **Configuration drift** | Likely | Impossible |

## Why Companies Use NixOS

1. **Reproducible infrastructure** (Nix used by major companies for CI/CD)
2. **Atomic deployments** (banking, critical systems)
3. **Compliance auditing** (every change tracked in git)
4. **Disaster recovery** (rebuild from git, instantly)
5. **Polyglot environments** (support 50 different package versions simultaneously)

## The Paradigm Shift

NixOS requires thinking differently about systems:

**Traditional**: "How do I configure this system?"
**NixOS**: "What should this system be?"

**Traditional**: "Did that change get applied?"
**NixOS**: "Here's exactly what's deployed (verify from git)"

**Traditional**: "How do I avoid breaking prod?"
**NixOS**: "How do I quickly fix prod if something breaks?"

---

## Next: The Nix Language

Now that you understand what NixOS *does*, let's learn the minimal amount of Nix language needed to write powerful configurations.

[Continue to: Foundation 3 - Nix Language Basics](./03-nix-language-basics.md)
