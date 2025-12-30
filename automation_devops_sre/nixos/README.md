# NixOS: From "What's a Nix?" to Production Systems
## Complete Educational Course with ELI5 Methodology

### Course Philosophy

NixOS isn't just another Linux distro. It's a fundamentally different way of thinking about systems management. Instead of manually tweaking config files (and breaking things), NixOS declares "this is exactly what my system should look like" and makes it happen. Every. Single. Time.

Think of it like:
- **Traditional Linux**: Building a house by telling workers "add a brick here, paint that wall blue, maybe install plumbing"
- **NixOS**: Writing a blueprint that says "this house has 3 bedrooms, blue walls, copper plumbing" and workers build it perfectly every time

### Learning Path

```
[FOUNDATION]
  ├── What is Nix? (The Package Manager)
  ├── What is NixOS? (The Operating System)
  ├── The Nix Language (Minimal but Powerful)
  │
[BEGINNER LABS]
  ├── Lab 1: Your First NixOS Installation & Configuration
  ├── Lab 2: Declaring Your Perfect Development Environment
  ├── Lab 3: Creating Reproducible Development Shells
  │
[INTERMEDIATE]
  ├── Understanding the Nixpkgs Ecosystem
  ├── Writing Custom Packages
  ├── Building Modular Configurations
  │
[INTERMEDIATE LABS]
  ├── Lab 4: Multi-Host Configuration Management
  ├── Lab 5: Building a Development Server
  ├── Lab 6: Deploying Applications with NixOS
  │
[ADVANCED]
  ├── Building Custom Distributions
  ├── Creating Flakes (Modern Nix)
  ├── Integrating with CI/CD
  │
[ADVANCED LABS]
  ├── Lab 7: Home-Manager for User Configuration
  ├── Lab 8: NixOS Flakes in Production
  ├── Lab 9: Hybrid Infrastructure (NixOS + Cloud)
```

### Directory Structure

```
nixos-course/
├── README.md (this file)
├── 00-foundation/
│   ├── 01-what-is-nix.md
│   ├── 02-what-is-nixos.md
│   ├── 03-nix-language-basics.md
│   └── ascii-diagrams.md
├── labs/
│   ├── lab-01-first-install/
│   │   ├── README.md (instructions)
│   │   ├── configuration.nix (starter template)
│   │   ├── WALKTHROUGH.md (step-by-step)
│   │   └── troubleshooting.md
│   ├── lab-02-dev-environment/
│   ├── lab-03-dev-shells/
│   ├── lab-04-multi-host/
│   ├── lab-05-server/
│   ├── lab-06-deployments/
│   ├── lab-07-home-manager/
│   ├── lab-08-flakes/
│   └── lab-09-hybrid-infra/
├── reference/
│   ├── common-packages.md
│   ├── nixpkgs-structure.md
│   ├── debugging-nix.md
│   └── performance-tips.md
└── examples/
    ├── minimal-system/
    ├── dev-workstation/
    ├── server-setup/
    └── flake-templates/
```

### How to Use This Course

**For Absolute Beginners**: Start with `00-foundation/` and read in order. Don't skip the "why" explanations.

**For Infrastructure People**: Jump to foundational concepts, then go straight to `lab-04-multi-host/`. The analogies will make sense given your background.

**For Package Maintainers**: Start with `02-what-is-nixos.md`, then jump to Lab 6 and beyond.

### Key Principle: Every Lab is Production-Ready

Unlike toy examples, each lab provides:
- Complete, working configuration
- Heavily commented code explaining each decision
- Troubleshooting section with real problems and solutions
- "What if I want to..." extension section
- Docker Compose for testing (no NixOS install required for learning)

### Prerequisites

- Comfort with Linux concepts (users, permissions, filesystems, package managers)
- Ability to read shell scripts
- Patience (Nix has a learning curve, but it's worth it)
- Text editor and terminal

### Core Concepts You'll Master

1. **Declarative Configuration**: Describe state, not steps
2. **Reproducibility**: Same config = same system, always
3. **Immutability**: Changes are atomic or they don't happen
4. **Referential Transparency**: Pure functions and no side effects
5. **Composition**: Build complex systems from simple pieces

### What Makes This Different

Most NixOS resources assume you know Nix. This course assumes you know Linux and want to learn a better way. We explain:
- Not just *what* to write, but *why*
- Not just syntax, but philosophy
- Not just examples, but patterns
- Not just success paths, but how to debug failures

Each lab can be completed in 30-60 minutes and builds toward actual systems you'd run in production.

---

**Ready?** Start with `00-foundation/01-what-is-nix.md`
