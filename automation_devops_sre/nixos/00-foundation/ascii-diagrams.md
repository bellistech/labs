# NixOS: Visual Diagrams for Key Concepts

## Diagram 1: Traditional vs NixOS Package Management

```
TRADITIONAL PACKAGE MANAGER (apt, yum)
=====================================

User: $ apt install nodejs

    ┌─────────────────────────────┐
    │ Package Index               │
    │ (what's available to        │
    │ install)                    │
    └─────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────┐
    │ Dependency Resolver         │
    │ (figures out what nodejs    │
    │ needs)                      │
    └─────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────┐
    │ Download & Install          │
    │ /usr/bin/node (shared)      │
    │ /usr/lib/libssl.so (shared) │
    │ /usr/lib/libz.so (shared)   │
    └─────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────┐
    │ Problems:                   │
    │ - Conflicts with other apps │
    │ - "Works on my machine"     │
    │ - Hard to rollback          │
    │ - Version hell              │
    └─────────────────────────────┘

NIX PACKAGE MANAGER
===================

User: $ nix-env -i nodejs

    ┌─────────────────────────────┐
    │ Nix Expression              │
    │ (full recipe: source code   │
    │ + dependencies +            │
    │ build flags)                │
    └─────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────┐
    │ Dependency Graph            │
    │ (recursively calculate ALL  │
    │ dependencies)               │
    └─────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────┐
    │ Compute Output Hash         │
    │ (based on everything)       │
    │ → abc123nodejs18.5.0        │
    └─────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────────────┐
    │ Pure Build Environment              │
    │ (isolated, deterministic)           │
    │ Inputs: source + deps + compiler    │
    └─────────────────────────────────────┘
             │
             ▼
    ┌─────────────────────────────────────┐
    │ Output Hash Verified                │
    │ (same inputs = same output always)  │
    └─────────────────────────────────────┘
             │
             ▼
    ┌──────────────────────────────────────┐
    │ Store in Nix Store                   │
    │ /nix/store/abc123-nodejs-18.5.0/bin/ │
    │ /nix/store/def456-openssl-3.0.1/lib/ │
    │ /nix/store/ghi789-icu-72.1/lib/      │
    │                                      │
    │ Benefits:                            │
    │ ✓ No conflicts (separate instances) │
    │ ✓ Reproducible (same everywhere)    │
    │ ✓ Rollback (previous version exists)│
    │ ✓ Coexistence (multiple versions)   │
    └──────────────────────────────────────┘
```

---

## Diagram 2: NixOS System Composition

```
YOUR SYSTEM DECLARATION
=======================

  ┌────────────────────────────────────────────────────────┐
  │   /etc/nixos/configuration.nix                         │
  │                                                        │
  │   This file describes your ENTIRE system              │
  │   It's the single source of truth                      │
  └────────────────────────────────────────────────────────┘
           │
           │  $ nixos-rebuild switch
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Nix Evaluator (reads your config)                    │
  │                                                        │
  │   Interprets all Nix code                              │
  │   Resolves all references                              │
  │   Calculates all dependencies                          │
  └────────────────────────────────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Dependency Graph (What needs what?)                  │
  │                                                        │
  │   services.nginx                                       │
  │     ├─ nginx binary (depends on openssl)              │
  │     ├─ openssl library (depends on libc)              │
  │     ├─ libc (system library)                          │
  │     ├─ systemd unit (service startup)                 │
  │     └─ nginx.conf (generated from your settings)      │
  │                                                        │
  │   users.users.alice                                    │
  │     ├─ Create account                                  │
  │     ├─ Set shell to /nix/store/.../bash              │
  │     └─ Add to groups                                   │
  │                                                        │
  │   environment.systemPackages = [ git vim htop ]       │
  │     ├─ git (depends on perl, openssl, curl...)       │
  │     ├─ vim (depends on ncurses...)                    │
  │     └─ htop (depends on ncurses...)                   │
  └────────────────────────────────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Realize Derivations (build/fetch packages)           │
  │                                                        │
  │   /nix/store/abc123-nginx-1.24/                       │
  │   /nix/store/def456-openssl-3.0.1/                    │
  │   /nix/store/ghi789-libc-2.37/                        │
  │   /nix/store/... (thousands more)                      │
  │                                                        │
  │   System Closure: Complete dependency tree            │
  │   (everything needed to boot this system)             │
  └────────────────────────────────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Generate Configuration Files                         │
  │                                                        │
  │   NixOS creates:                                       │
  │   /etc/passwd (from users config)                     │
  │   /etc/nginx/nginx.conf (from nginx settings)         │
  │   /etc/systemd/system/nginx.service                   │
  │   /etc/fstab (from filesystems config)                │
  │   ... (everything else)                               │
  └────────────────────────────────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Build System Derivation                              │
  │                                                        │
  │   One atomic "system" package                          │
  │   /nix/store/system123-nixos-23.11/                   │
  │   Contains: kernel, bootloader, modules, all configs  │
  │                                                        │
  │   Compute: system-closure-hash                        │
  │   (unique identifier for this exact system)           │
  └────────────────────────────────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Atomic Switch                                        │
  │                                                        │
  │   GRUB bootloader now points to:                       │
  │   /nix/store/system123.../                            │
  │                                                        │
  │   ✓ If build succeeds: switch completes               │
  │   ✓ If build fails: switch never happens              │
  │   ✓ Previous system still exists (can rollback)       │
  │                                                        │
  │   /nix/var/nix/profiles/system -> gen-45              │
  │                                                        │
  │   Older generations:                                   │
  │   /nix/var/nix/profiles/system-44-link                │
  │   /nix/var/nix/profiles/system-43-link                │
  │   ... (all kept, you choose when to delete)           │
  └────────────────────────────────────────────────────────┘
           │
           ▼
  ┌────────────────────────────────────────────────────────┐
  │   Running System                                       │
  │                                                        │
  │   Kernel: Exactly what config specified               │
  │   Services: Exactly what config enabled               │
  │   Packages: Exactly what config listed                │
  │   Users: Exactly what config declared                 │
  │   Network: Exactly what config set                    │
  │                                                        │
  │   Every. Single. Detail. Matches. Your. Config.        │
  └────────────────────────────────────────────────────────┘
```

---

## Diagram 3: Nix Store Structure

```
NIX STORE: Your Software Vault
================================

  /nix/store/
  
  ├─ [hash1]-package-name-version/
  │  ├─ bin/
  │  │  └─ executable
  │  ├─ lib/
  │  │  └─ libfoo.so
  │  └─ share/
  │     └─ docs
  │
  ├─ [hash2]-dependency-version/
  │  ├─ lib/
  │  │  └─ libdep.so
  │  └─ ...
  │
  ├─ [hash3]-nodejs-18.5.0/
  │  ├─ bin/node
  │  ├─ lib/libssl.so -> ../../[hash4]-openssl/lib/
  │  └─ ...
  │
  ├─ [hash4]-openssl-3.0.1/
  │  └─ lib/libssl.so
  │
  └─ [hashN]-package-N-version/
     └─ ...


KEY INSIGHT: Each directory's name includes a HASH

  [hash] = SHA256(package + all dependencies + compiler + flags)

  Same inputs → Same hash → Same directory
  Different inputs → Different hash → Different directory
  No conflicts! Multiple versions coexist peacefully


EXAMPLE PATH STRUCTURE:

  Traditional:
    /usr/bin/node
    /usr/lib/libssl.so.1.1
    
  Problem: What version? Conflicts when upgrading.

  Nix:
    /nix/store/abc123-nodejs-18.5.0-with-openssl-3.0.1/bin/node
    /nix/store/def456-nodejs-16.0.0-with-openssl-1.1.1/bin/node
    /nix/store/ghi789-openssl-3.0.1/lib/libssl.so.1.1
    /nix/store/jkl012-openssl-1.1.1/lib/libssl.so.1.1
    
  Solution: Multiple versions coexist with specific dependencies


PACKAGES ARE LAZY-LOADED (only use what you need):

  User profile (symlinks to actual packages):
  
    ~/.nix-profile/bin/node -> /nix/store/abc123.../bin/node
    ~/.nix-profile/bin/npm  -> /nix/store/abc123.../bin/npm
    ~/.nix-profile/lib/...  -> /nix/store/abc123.../lib/...
    
  If you `nix-env -e nodejs`, just remove symlinks
  Package in store remains (garbage collection later)
```

---

## Diagram 4: System Generation Timeline

```
SYSTEM GENERATIONS: Version Control for Your Entire OS
=========================================================

  Configuration Change Timeline:

  Time 1:
    $ nixos-rebuild switch
    ✓ Generation 1 created
    
  Time 2:
    $ echo "services.ssh.enable = true;" >> configuration.nix
    $ nixos-rebuild switch
    ✓ Generation 2 created
    
  Time 3:
    $ echo "environment.systemPackages += [ docker ];" >> configuration.nix
    $ nixos-rebuild switch
    ✓ Generation 3 created
    
  Time 4:
    $ cat < /dev/zero | dd of=/tmp/file (oops!)
    💥 System getting weird, broke something
    $ sudo nixos-rebuild switch --rollback
    ✓ Generation 2 restored
    ✓ System back to working state
    ✓ Boot sees previous generation in GRUB
    
  Time 5:
    $ nixos-rebuild switch --rollback
    ✓ Generation 1 restored


GRUB BOOT MENU (After Several Changes):
=========================================

  ┌──────────────────────────────────────┐
  │  NixOS Boot Menu                     │
  ├──────────────────────────────────────┤
  │ [*] NixOS-23.11 (Generation 5)       │  ← Current boot (generation 5)
  │ [ ] NixOS-23.11 (Generation 4)       │     Can select to boot into
  │ [ ] NixOS-23.11 (Generation 3)       │     previous version
  │ [ ] NixOS-23.11 (Generation 2)       │
  │ [ ] NixOS-23.11 (Generation 1)       │
  ├──────────────────────────────────────┤
  │ Press 'e' to edit, 'c' for console  │
  └──────────────────────────────────────┘


WHAT GETS STORED IN EACH GENERATION:

  /nix/var/nix/profiles/system-45-link
    ├─ kernel
    ├─ initrd (boot ramdisk)
    ├─ bootloader config (grub.cfg)
    ├─ systemd units (services)
    ├─ /etc configs (passwd, hostname, etc.)
    ├─ udev rules
    ├─ installed packages
    └─ everything else to boot


CLEANUP:

  Old generations stay forever by default
  
  View them:
    $ nix-env --list-generations
    # 1   2024-01-15 12:30:15
    # 2   2024-01-15 12:45:32
    # 3   2024-01-15 13:20:11
    
  Delete all but current:
    $ nix-collect-garbage
    
  Delete all (including current):
    $ nix-collect-garbage -d
    
  Keep last N generations:
    $ nix-env --delete-generations +5  # Keep last 5
```

---

## Diagram 5: Declarative vs Imperative

```
IMPERATIVE (Traditional Linux)
==============================

  What you do:        What the system looks like:
  
  $ apt install nginx      /usr/bin/nginx (somewhere)
  $ nginx -g daemon on     [process running, maybe]
  $ cp nginx.conf /etc/    /etc/nginx/nginx.conf (maybe yours?)
  $ systemctl enable nginx [nginx starts next boot, hopefully?]
  
  Problem: Hard to know exact state, hard to replicate
  
  
DECLARATIVE (NixOS)
===================

  What you do:        What the system becomes:
  
  services.nginx.enable = true;
       │
       └─→ Reads Nix expression
           Resolves dependencies
           Builds derivations
           Generates nginx.conf automatically
           Creates systemd unit automatically
           Enables service automatically
           Switches entire system atomically
           └─→ NGINX RUNNING EXACTLY AS DECLARED


IMPERATIVE PROBLEMS:

  Configuration Drift:
    ┌─────────────────┐
    │ Server A        │   Manual changes over time
    │ manual changes  │   (what was changed? why? when?)
    └─────────────────┘   
                          Lost history
                          Unmaintainable
                          Can't replicate
    
    ┌─────────────────┐
    │ Server B        │   Different history
    │ different drift │   Different state
    └─────────────────┘   
           ↕
      NOT IDENTICAL!


DECLARATIVE SOLUTION:

    ┌─────────────────────────────────────┐
    │ configuration.nix (source control)  │
    └─────────────────────────────────────┘
              │
              ├─→ Server A:  nix-rebuild switch
              │   Result: Identical system
              │
              ├─→ Server B:  nix-rebuild switch
              │   Result: Identical system
              │
              └─→ Server C:  nix-rebuild switch
                  Result: Identical system
    
    Any server, any time, same config = same system
    Configuration drift: IMPOSSIBLE
```

---

## Diagram 6: Reproducibility Promise

```
THE NIX REPRODUCIBILITY GUARANTEE
==================================

  A = Input specification
      (source code + dependencies + compiler version + flags)
  
  f(A) = Build process (Nix build system)
  
  B = Build result
      (compiled binary + all dependencies)
  
  
THEORY:
  f(A) = f(A)  →  B = B (for all time)
  
  Same input forever → Same output forever
  (deterministic builds)


PRACTICE:

  Monday 2024-01-15:
    Developer on MacBook builds package with configuration A
    Result hash: abc123
    Binary stored with configuration A
    
  Wednesday 2024-01-17:
    Same developer, same machine, rebuilds with configuration A
    Result hash: abc123 (SAME!)
    
  Following year 2025-01-15:
    Different developer, Linux server, builds with configuration A
    Result hash: abc123 (STILL SAME!)
    
  Why? Because:
    ✓ Source code is pinned (specific commit)
    ✓ Dependencies are pinned (exact versions)
    ✓ Build environment is pure (no system state pollution)
    ✓ Compiler versions are pinned
    ✓ All flags are identical
    
  
WHAT THIS MEANS:

  ✗ "Works on my machine but not yours" - IMPOSSIBLE with Nix
  ✗ "Let's try rebuilding" (hoping it magically works) - FIXED
  ✗ "Did you upgrade that library?" - TRACKED in expression
  
  ✓ "This exact version works on Linux/Mac/CI/Prod" - GUARANTEED
  ✓ "Production matches dev environment" - IDENTICAL
  ✓ "Can rebuild from 5 years ago" - EXACT SAME RESULT


ANTI-EXAMPLE: npm/pip/cargo without lockfile

  npm install                    ← "Install latest of everything"
  Dependency: express ^4.0       ← "4.0 to 4.999"
  
  Day 1: Installs express 4.17.1 (latest that day)
  Day 100: `npm install` installs express 4.20.5 (new latest)
  
  Different versions → Different behavior → Breaks things
  
  This is WHY every modern language needs lock files now


WITH NIXOS:

  nix-shell -p nodejs=18.5.0    ← Exact version, always
  Result: Every developer, always gets 18.5.0
  (even if 19.0.0 exists)
  
  Configuration.nix = the lock file for your entire system
```

---

## Diagram 7: NixOS Abstractions Over Time

```
LEARNING NixOS: Layers of Understanding
=========================================

LAYER 1: Package Management
  
    "I want Node.js"
    $ nix-env -i nodejs
    [Node.js installed in isolation]
    
    Concepts: packages, profiles, garbage collection


LAYER 2: Declarative Configuration
  
    "I want this system to always have Node.js"
    $ echo "environment.systemPackages = with pkgs; [ nodejs ];" 
    [System remembers, survives reboot]
    
    Concepts: configuration.nix, nixos-rebuild, reproducibility


LAYER 3: Atomic Deployments
  
    "Safe updates with rollback"
    $ nixos-rebuild switch
    [Entire system updated atomically, rollback possible]
    
    Concepts: system closures, atomic switching, generations


LAYER 4: Multi-System Management
  
    "Deploy same config to 10 servers"
    nixos-build-vms configuration.nix
    [Build VMs matching exact configuration]
    
    Concepts: configuration management at scale


LAYER 5: Flakes (Modern Nix)
  
    "Pinned dependencies for system, reproducible globally"
    flake.nix: locks nixpkgs version + overlays
    [System reproduces exactly across time/place]
    
    Concepts: flake.lock, inputs, outputs


LAYER 6: Module System
  
    "Compose configurations from modular pieces"
    imports = [ ./hardware.nix ./services.nix ];
    [Build complex systems from simple, reusable modules]
    
    Concepts: options, config, implementation


YOU ARE HERE (After completing this course):
  ✓ Layers 1-3 solidly understood
  ✓ Layer 4 (multi-system) ready to tackle
  ✓ Layer 5 (Flakes) available for depth dives
```

---

## Diagram 8: Common Misconceptions Clarified

```
MISCONCEPTION 1: "Nix uses more disk space"

  Assumption: Multiple package versions = disk bloat
  
  Reality:
    Shared dependencies ARE deduplicated
    Only unique combinations stored separately
    1000 packages all using libc 2.37 = ONE libc 2.37 on disk
    Only different when: libc 2.37 vs libc 2.36 (different hash)
    
  Typical: Traditional Linux + Nix ≈ 2-3x for same packages
  (not 1000x)


MISCONCEPTION 2: "Nix is slower than traditional"

  Build time: Same or faster (pure builds enable caching)
  Installation time: Fast (packages are pre-built)
  Runtime: Identical (same binaries)
  
  Slowness comes from: First time you build something


MISCONCEPTION 3: "Nix locks you into Nix packages"

  False!
  
  Nix can package ANYTHING:
  - Proprietary software
  - Custom scripts
  - AppImages
  - Docker containers
  - Virtual machines
  
  You're never locked in


MISCONCEPTION 4: "I need to learn Haskell to use Nix"

  False!
  
  Nix language is simple for configuration
  You need: 5 concepts (covered in foundation)
  You don't need: Functors, monads, type theory


MISCONCEPTION 5: "NixOS is unstable"

  False!
  
  NixOS is used in production by companies like:
  - Tweag
  - Determinate Systems
  - Various financial companies
  - Major open source projects
  
  Reputation: Unstable only because of learning curve
  Actual reality: Very stable, predictable


MISCONCEPTION 6: "I can't use NixOS for production"

  False!
  
  NixOS in production:
  ✓ Atomic deployments (safer than traditional)
  ✓ Rollback capability (disaster recovery)
  ✓ Immutability (prevents configuration drift)
  ✓ Reproducible infrastructure (compliance auditing)
  
  It's arguably MORE suitable than traditional Linux
```

---

## Key Takeaway

Each diagram shows a different aspect of WHY NixOS works:

1. **Package management** is deterministic
2. **System composition** is declarative
3. **Store structure** prevents conflicts
4. **Generations** enable safe testing
5. **State management** is explicit not implicit
6. **Reproducibility** is baked in
7. **Abstraction layers** make it learnable
8. **Misconceptions** are just that - misconceptions

Master these concepts, and NixOS becomes your superpower for infrastructure management.
