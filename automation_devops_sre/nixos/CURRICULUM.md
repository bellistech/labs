# NixOS ELI5 Course - Complete Curriculum

## What's Included (Core Foundation)

### Foundation Materials ✅

```
00-foundation/
├── 01-what-is-nix.md (Complete)
│   ├── The Problem Nix Solves
│   ├── How Nix is Different
│   ├── Key Nix Concepts
│   ├── The Nix Store
│   ├── Nix Language Introduction
│   └── Real-World Impact
│
├── 02-what-is-nixos.md (Complete)
│   ├── From Package Manager to Entire System
│   ├── The Core Idea
│   ├── What Happens When You Deploy
│   ├── The Nix Store Scales
│   ├── Key NixOS Principles
│   ├── Real-World Examples
│   └── Paradigm Shift
│
├── 03-nix-language-basics.md (Complete)
│   ├── Core Concept 1: Everything is an Expression
│   ├── Core Concept 2: Attribute Sets
│   ├── Core Concept 3: Let-In
│   ├── Core Concept 4: String Interpolation
│   ├── Core Concept 5: Lists
│   ├── The Pattern You'll See Everywhere
│   ├── Understanding 'with'
│   ├── Functions
│   ├── Conditionals
│   ├── Common Pitfalls
│   ├── The Structure of Every Configuration
│   └── Practice Exercises (with answers)
│
└── ascii-diagrams.md (Complete)
    ├── Diagram 1: Traditional vs NixOS Package Management
    ├── Diagram 2: NixOS System Composition
    ├── Diagram 3: Nix Store Structure
    ├── Diagram 4: System Generation Timeline
    ├── Diagram 5: Declarative vs Imperative
    ├── Diagram 6: Reproducibility Promise
    ├── Diagram 7: NixOS Abstractions Over Time
    └── Diagram 8: Common Misconceptions
```

### Beginner Labs ✅

```
labs/lab-01-first-install/ (Complete)
├── README.md (Full installation walkthrough)
│   ├── Lab Overview
│   ├── Installation Steps (Detailed)
│   ├── Explore Your System
│   ├── Make Your First Change
│   ├── Try a Rollback
│   ├── Troubleshooting Common Issues
│   ├── Verification Checklist
│   └── What You've Learned
│
└── configuration.nix (Starter template)
    ├── Heavily commented
    ├── Explains each section
    ├── Good starting point
    └── Extensible for customization
```

### Quick Start & Navigation ✅

```
├── README.md (Complete course overview)
│   └── How to use everything
│
├── QUICKSTART.md (Path selection guide)
│   ├── Learning paths by role
│   ├── Time-based recommendations
│   ├── Decision trees
│   └── Success metrics
```

---

## What's Planned (Future Expansion)

### Intermediate Labs (Planned)

```
labs/lab-02-dev-environment/ (TODO)
├── README.md
│   ├── Creating development environments
│   ├── Virtual package availability
│   ├── Nix-shell basics
│   ├── Shell.nix patterns
│   └── Multiple development contexts
├── shell.nix (Development template)
└── TROUBLESHOOTING.md

labs/lab-03-dev-shells/ (TODO)
├── README.md
│   ├── Declarative dev environments
│   ├── Language-specific shells
│   ├── Direnv integration
│   ├── Team consistency
│   └── Per-project isolation
├── node-shell.nix
├── rust-shell.nix
├── python-shell.nix
└── examples/

labs/lab-04-multi-host/ (TODO)
├── README.md
│   ├── Managing multiple systems
│   ├── Configuration inheritance
│   ├── NixOS modules
│   ├── Shared configurations
│   └── Fleet management patterns
├── base-config.nix
├── server-01.nix
├── server-02.nix
├── server-03.nix
└── networking/

labs/lab-05-server-setup/ (TODO)
├── README.md
│   ├── Production server configuration
│   ├── Services configuration
│   ├── Networking & firewall
│   ├── User management
│   └── Security hardening
├── configuration.nix (Server template)
└── security-checklist.md

labs/lab-06-deployments/ (TODO)
├── README.md
│   ├── Automated deployments
│   ├── Nixops patterns
│   ├── CI/CD integration
│   ├── Container deployments
│   └── Atomic rollouts
├── deployment.nix
└── ci-cd-examples/

labs/lab-07-home-manager/ (TODO)
├── README.md
│   ├── Managing user environments
│   ├── Home configuration
│   ├── Dotfile management
│   ├── Shell configuration
│   └── Application settings
├── home.nix
└── examples/

labs/lab-08-flakes/ (TODO)
├── README.md
│   ├── Modern Nix (Flakes)
│   ├── Pinned dependencies
│   ├── Reproducibility at scale
│   ├── Flake structure
│   └── Lock files (flake.lock)
├── flake.nix
├── flake.lock
└── templates/

labs/lab-09-hybrid-infra/ (TODO)
├── README.md
│   ├── NixOS + Cloud (AWS/GCP)
│   ├── Terraform + NixOS
│   ├── Managed services integration
│   ├── Hybrid deployments
│   └── Cost optimization
├── aws-deployment.nix
└── terraform-examples/
```

### Reference Materials (Planned)

```
reference/
├── common-packages.md (TODO)
│   ├── Essential system packages
│   ├── Development tools
│   ├── Server applications
│   └── Finding packages
│
├── nixpkgs-structure.md (TODO)
│   ├── Understanding package organization
│   ├── Package overlays
│   ├── Custom package creation
│   └── Contributing packages
│
├── debugging-nix.md (TODO)
│   ├── Build failures
│   ├── Runtime errors
│   ├── Configuration errors
│   ├── Debugging tools
│   └── Inspection commands
│
├── performance-tips.md (TODO)
│   ├── Caching strategies
│   ├── Build optimization
│   ├── Dependency minimization
│   └── Binary cache setup
│
├── migration-guide.md (TODO)
│   ├── From traditional Linux
│   ├── From Docker
│   ├── From Ansible
│   └── Common patterns
│
├── security-hardening.md (TODO)
│   ├── Network security
│   ├── User permissions
│   ├── Service isolation
│   ├── Secret management
│   └── Audit logging
│
└── troubleshooting-index.md (TODO)
    ├── Common problems
    ├── Error message index
    ├── Solution strategies
    └── When to ask for help
```

### Advanced Topics (Planned)

```
advanced/
├── custom-packages.md (TODO)
│   ├── Creating your own packages
│   ├── Packaging patterns
│   ├── Build inputs
│   └── Nixpkgs integration
│
├── overlays.md (TODO)
│   ├── Package overriding
│   ├── Custom modifications
│   ├── Conditional packages
│   └── Composition patterns
│
├── flakes-advanced.md (TODO)
│   ├── Complex flake structures
│   ├── Multiple outputs
│   ├── Flake metadata
│   └── Dependency management
│
├── module-system.md (TODO)
│   ├── Writing modules
│   ├── Module options
│   ├── Config merging
│   └── Composition patterns
│
├── nixops.md (TODO)
│   ├── Declarative deployments
│   ├── Multi-system coordination
│   ├── State management
│   └── Production patterns
│
└── contributing.md (TODO)
    ├── Contributing to nixpkgs
    ├── PR process
    ├── Packaging standards
    └── Maintenance guidelines
```

### Examples (Planned)

```
examples/
├── minimal-system/ (TODO)
│   └── Bare minimum viable NixOS
│
├── dev-workstation/ (TODO)
│   └── Full development environment
│
├── server-setup/ (TODO)
│   ├── Web server (nginx + PostgreSQL)
│   ├── API server
│   ├── Database server
│   └── Monitoring stack
│
├── docker-integration/ (TODO)
│   └── NixOS containers
│
├── flake-templates/ (TODO)
│   ├── project-template/
│   ├── python-app/
│   ├── node-app/
│   ├── go-service/
│   └── rust-service/
│
├── terraform-examples/ (TODO)
│   ├── AWS deployment
│   ├── GCP deployment
│   └── Multi-cloud
│
└── ci-cd-examples/ (TODO)
    ├── GitHub Actions
    ├── GitLab CI
    ├── Jenkins
    └── Custom automation
```

---

## Learning Progression

### Phase 1: Foundation (1-2 hours) ✅
**Status**: Complete and ready
- What is Nix?
- What is NixOS?
- Nix language basics (5 concepts)
- ASCII diagrams for visualization

**Outcome**: Understand the philosophy and approach

### Phase 2: First Deployment (1-2 hours) ✅
**Status**: Complete and ready
- Lab 1: First installation
- Hands-on experience with rebuild
- Understanding system generations
- Safe rollback practice

**Outcome**: Working NixOS system, practical knowledge

### Phase 3: Development Environments (2-3 hours) 🔄
**Status**: Planned (shells, dev-shell, multi-language)
- Lab 2: Development environment
- Lab 3: Dev shells (Nix way of project setup)
- Language-specific examples

**Outcome**: Using NixOS for development teams

### Phase 4: Multi-System Management (3-4 hours) 🔄
**Status**: Planned
- Lab 4: Multi-host configuration
- Configuration composition
- NixOS modules
- Fleet management

**Outcome**: Managing infrastructure at scale

### Phase 5: Production Deployment (4-6 hours) 🔄
**Status**: Planned
- Lab 5: Server setup and hardening
- Lab 6: Deployment automation
- Lab 8: Flakes for reproducibility
- Production patterns

**Outcome**: Production-ready NixOS deployments

### Phase 6: Advanced Topics (6-8 hours) 🔄
**Status**: Planned
- Custom packages
- Overlays and composition
- Nixpkgs contribution
- Complex architectures

**Outcome**: Deep NixOS expertise and community contribution

---

## Content Quality Standards

Every piece of content in this course follows these principles:

### 1. Heavily Commented Code ✅
- Every line explained
- Why decisions were made
- Common alternatives
- Pitfalls and gotchas

### 2. Multiple Explanations ✅
- Simple first (ELI5)
- Detailed explanation
- Technical deep dive
- Visual diagrams

### 3. Hands-On Learning ✅
- Read/understand
- Implement/practice
- Troubleshoot/debug
- Verify/confirm

### 4. Real-World Examples 🔄
- Not toy examples
- Production-ready patterns
- Extensible templates
- Clear evolution path

### 5. Comprehensive Troubleshooting ✅
- Common errors
- Why they happen
- How to fix them
- Prevention strategies

### 6. Progressive Difficulty 🔄
- Start simple
- Build complexity
- Reference earlier concepts
- Spiral learning model

---

## Current Status Summary

```
Foundation Materials: ✅ COMPLETE
├── What is Nix: Complete
├── What is NixOS: Complete
├── Nix Language Basics: Complete
└── ASCII Diagrams: Complete

Beginner Labs: ✅ COMPLETE
├── Lab 1: First Install: Complete
└── Lab 1 Configuration Template: Complete

Navigation & Quick Start: ✅ COMPLETE
├── Main README: Complete
├── QUICKSTART Guide: Complete
└── Learning Paths: Complete

Intermediate Labs: 🔄 PLANNED (Ready to build)
├── Lab 2: Dev Environment
├── Lab 3: Dev Shells
└── Lab 4: Multi-Host

Advanced Labs: 🔄 PLANNED
├── Lab 5: Server Setup
├── Lab 6: Deployments
├── Lab 7: Home Manager
├── Lab 8: Flakes
└── Lab 9: Hybrid Infrastructure

Reference Materials: 🔄 PLANNED
├── Common Packages
├── Debugging Guide
├── Security Hardening
└── Troubleshooting Index

Examples: 🔄 PLANNED
├── Minimal System
├── Development Workstation
├── Server Configurations
└── Template Projects
```

---

## How to Expand This Course

### To Add a New Lab

1. Create directory: `labs/lab-XX-topic/`
2. Create `README.md` with:
   - What you'll learn
   - Estimated time
   - Prerequisites
   - Step-by-step walkthrough
   - Troubleshooting section
   - Verification checklist
3. Create templates/examples
4. Create WALKTHROUGH.md for detailed explanation
5. Update course README

### To Add Reference Material

1. Create file in `reference/` directory
2. Include:
   - Table of contents
   - Examples with explanations
   - Links to relevant labs
   - External resources
3. Update reference index

### To Add Examples

1. Create directory in `examples/`
2. Create complete, working configuration
3. Include README explaining:
   - What this example shows
   - Key components
   - How to customize it
   - Common modifications
4. Make it copy-paste ready

---

## Using This Course

### As an Individual Learner

1. Start with QUICKSTART guide
2. Choose your path based on role/experience
3. Follow recommended labs
4. Reference materials as needed
5. Build your own systems

### As a Team Teaching Material

1. Foundation materials in reading group
2. Lab 1 as group exercise
3. Advanced labs as team projects
4. Reference materials on wiki
5. Examples as templates for team infrastructure

### As an Organization Building Infrastructure

1. Foundation for all new team members
2. Lab 1 + Lab 4 for standard deployments
3. Lab 8 for production systems
4. Examples as internal templates
5. Customization guide for your needs

### As Contributing to Open Source

1. Complete all foundation material
2. Create new packages (future advanced lab)
3. Submit to nixpkgs
4. Reference this course in PRs
5. Help improve NixOS ecosystem

---

## Maintenance & Updates

### When to Update Content

- **Bug fixes**: Always, include clarification
- **Nix version changes**: Update examples to current
- **New features**: Add after stabilization
- **Community feedback**: Incorporate learning
- **New patterns emerge**: Document and teach

### Version Strategy

```
Course version: Tied to NixOS release
- Course 23.11: Targets NixOS 23.11 (stable)
- Course 24.05: Targets NixOS 24.05 (stable)

Each version includes:
- Updated package names
- New features explanation
- Backward compatibility notes
- Migration guides (if needed)
```

---

## Getting Help Contributing

Want to expand this course? 

### For Lab Writers

- Pick a planned lab
- Follow the structure of Lab 1
- Make it detailed, commented, walkthrough-friendly
- Include troubleshooting
- Include verification checklist

### For Example Creators

- Pick a use case (dev workstation, server, etc.)
- Create production-ready configuration
- Include heavy comments
- Make it extensible
- Include README explaining

### For Documentation

- Pick a topic from "Planned" sections
- Research current best practices
- Include examples
- Include links to related materials
- Keep tone consistent

---

## Success Metrics

### For Students

After completing this course, they should:
- ✅ Understand NixOS philosophy
- ✅ Deploy working NixOS systems
- ✅ Make safe configuration changes
- ✅ Troubleshoot common issues
- ✅ Manage multiple systems
- ✅ Understand reproducibility benefits

### For Instructors

Course is successful when:
- ✅ Reduces NixOS learning time significantly
- ✅ Provides hands-on practical experience
- ✅ Explains "why" not just "how"
- ✅ Builds real confidence in learners
- ✅ Enables independent system building

### For the NixOS Community

Course impact:
- ✅ More confident NixOS users
- ✅ Better infrastructure practices
- ✅ Fewer "it doesn't work" posts
- ✅ More contributions to nixpkgs
- ✅ Stronger ecosystem

---

## The Ultimate Goal

This course exists to transform NixOS from "cool but confusing" to "my default choice for infrastructure."

By understanding:
1. The philosophy (declarative > imperative)
2. The tools (Nix language, nixos-rebuild)
3. The patterns (modules, compositions, flakes)
4. Real examples (labs, configurations)

Anyone can become proficient with NixOS and understand why it's powerful.

**The course is complete when:** Every learner finishes and thinks "I'll use NixOS for my next project."

---

## Next Steps

1. **Review**: Current complete foundation and Lab 1
2. **Feedback**: What's missing? What's unclear?
3. **Expand**: Labs 2-4 should come next (most practical)
4. **Polish**: Examples and reference materials
5. **Iterate**: Based on learner feedback

---

**Created for**: Infrastructure engineers, SREs, DevOps professionals
**Purpose**: Make NixOS accessible, practical, and powerful
**Philosophy**: Explain like I'm 5, deploy like I'm a professional

Happy learning! 🚀
