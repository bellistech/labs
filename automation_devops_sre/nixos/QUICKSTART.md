# NixOS ELI5 Course - Quick Start Guide

A rapid path through the course for people who want to understand AND deploy quickly.

## For Absolute Beginners (0-2 hours)

**Goal**: Understand what NixOS is and get it running

```
Start here → Read foundation → Deploy first lab → Success
```

### Track: Understanding (30 mins)

1. Read: [What is Nix?](./00-foundation/01-what-is-nix.md)
   - Focus: "The Problem Nix Solves" and "Key Nix Concepts"
   - Skip: Detailed mechanics (you'll learn by doing)

2. Read: [What is NixOS?](./00-foundation/02-what-is-nixos.md)
   - Focus: "The Core Idea", "Declarative System Management"
   - Remember: One config file describes entire system

3. Skim: [Nix Language Basics](./00-foundation/03-nix-language-basics.md)
   - Learn: 5 concepts you'll use
   - Don't memorize: You'll reference this while coding

### Track: Practical Deployment (90 minutes)

1. Prepare:
   - Have VM or spare machine ready
   - Download NixOS ISO
   - Allocate 20GB disk, 4GB RAM

2. Follow: [Lab 1: First Installation](./labs/lab-01-first-install/README.md)
   - This is the "learning by doing" part
   - Every step is explained
   - Troubleshooting guide included

3. Verify:
   - System boots successfully
   - Can login
   - Can run `nixos-rebuild switch`

**Result**: You have working NixOS and understand how changes work.

---

## For Infrastructure People (1-4 hours)

**Goal**: Understand NixOS advantages and deploy multi-system configuration

```
Quick concepts → Deep Lab 1 → Advanced Labs → Production patterns
```

### Track: Rapid Concepts (45 mins)

1. [What is NixOS?](./00-foundation/02-what-is-nixos.md)
   - Read all sections
   - Pay attention to: "Key NixOS Principles" and "Real-World Examples"

2. [ASCII Diagrams](./00-foundation/ascii-diagrams.md)
   - Study: Diagrams 1, 4, 5
   - These show system management advantages

### Track: Deploying Systems (Varies)

**Lab 1**: [First Installation](./labs/lab-01-first-install/README.md) (45 mins)
- Get baseline understanding

**Lab 4**: Multi-Host Configuration (when available)
- Configure 10 identical servers from one config
- Understand composition and modules

**Lab 8**: Flakes in Production (when available)
- Reproducible deployments with pinned versions
- Infrastructure-as-code patterns

### Track: Decision Making

Can you answer these?

1. "Why is NixOS safer for updates than traditional Linux?"
   - Answer should include: Atomic switching, rollback, generations

2. "How would I ensure 5 servers have identical configuration?"
   - Answer should include: Same configuration.nix, nixos-rebuild on each

3. "What happens if I make a bad change?"
   - Answer should include: Rollback to previous generation

If yes to all three: You're ready for production planning.

---

## For Package Maintainers (2-6 hours)

**Goal**: Understand how to create and maintain packages in Nix

Coming soon: Dedicated track with examples of common packages.

For now:
1. Complete [Nix Language Basics](./00-foundation/03-nix-language-basics.md)
2. Lab 2: Declaring Development Environment (coming)
3. Explore: https://search.nixos.org/packages

---

## For DevOps/SRE (3-8 hours)

**Goal**: Integrate NixOS into deployment pipeline, understand configuration management

### Must-Know Concepts

1. **Declarative configuration**: What makes it safer than imperative
2. **Atomic deployments**: How to safely roll out changes
3. **Reproducibility**: Why this matters for operations
4. **Rollback strategy**: How to recover from failures

### Essential Labs

1. [Lab 1](./labs/lab-01-first-install/): First deployment
2. [Lab 4](./labs/lab-04-multi-host/): Multi-system management (coming)
3. [Lab 6](./labs/lab-06-deployments/): Deployment automation (coming)

### Next Steps

- Integrate with Terraform for infrastructure provisioning
- Use `nixos-rebuild` in CI/CD pipelines
- Create system templates for common patterns
- Build deployment automation

---

## For Security/Compliance (2-4 hours)

**Goal**: Understand how NixOS helps with auditing, compliance, reproducibility

### Why NixOS Matters

- **Audit trail**: Every change in version control
- **Reproducibility**: Verify exact system state anytime
- **Rollback**: Fast recovery from security incidents
- **Immutability**: Configuration drift impossible
- **Binary cache**: Verify package sources

### Essential Reading

1. [What is NixOS?](./00-foundation/02-what-is-nixos.md) - "Key NixOS Principles"
2. [ASCII Diagram 5](./00-foundation/ascii-diagrams.md) - Declarative vs Imperative
3. Lab 1 - Understanding system reproducibility

### Key Questions

- "Show me the exact system configuration deployed on 2024-01-15"
  - Answer: `git log configuration.nix --until=2024-01-16`
  
- "What packages were in production last week?"
  - Answer: `git show HEAD~3:configuration.nix | grep systemPackages`
  
- "Did anything unauthorized get installed?"
  - Answer: Impossible - only changes in config happen

---

## Learning Paths by Time Commitment

### 2-Hour Path (Understanding + First Deploy)

```
1. Read: "What is NixOS?" (15 min)
2. Skim: "Nix Language Basics" (10 min)
3. Lab 1: Installation (60 min)
4. Lab 1: First rebuild (15 min)
Result: Working NixOS, basic understanding
```

### 4-Hour Path (Deep Understanding + Multi-System Ready)

```
1. Read: All foundation materials (45 min)
2. Study: ASCII diagrams (20 min)
3. Lab 1: Complete (60 min)
4. Lab 1: Troubleshooting & deep dive (30 min)
5. Understand: Multi-system strategy (15 min)
Result: Ready for multi-system deployment
```

### 8-Hour Path (Production Ready)

```
1. Complete: All foundation materials (90 min)
2. Lab 1: Full deep dive (90 min)
3. Lab 4: Multi-host configuration (120 min)
4. Lab 6: Deployment patterns (120 min)
5. Reference: Debugging guide (30 min)
Result: Production-ready NixOS knowledge
```

---

## Deciding Where to Start

### Question 1: Do you have Nix experience?

**No** → Start with [What is Nix?](./00-foundation/01-what-is-nix.md)
**Yes** → Start with [What is NixOS?](./00-foundation/02-what-is-nixos.md)

### Question 2: Why are you learning NixOS?

**"I want reproducible development environments"**
→ Focus: Lab 3 (Dev Shells) + Lab 2 (Development Environment)

**"I want safer system administration"**
→ Focus: Lab 1 + Lab 4 (Multi-Host)

**"I want atomic deployments for production"**
→ Focus: Lab 1 + Lab 6 + Lab 8 (Flakes)

**"I want to package software"**
→ Focus: Nix Language + Lab 6 (Package creation, coming)

**"I'm curious about how it works"**
→ Focus: ASCII diagrams + Lab 1

### Question 3: How much time do you have?

**30 minutes**: Read "What is NixOS?" only

**2 hours**: Read foundation + do Lab 1

**4 hours**: All foundation + full Lab 1 + understand multi-system patterns

**8+ hours**: Everything + advanced labs

---

## Red Flags If You're Struggling

### "I don't understand declarative vs imperative"

→ Re-read: [ASCII Diagram 5](./00-foundation/ascii-diagrams.md)
→ Key insight: State is explicit, not implicit

### "My system won't build"

→ Check: [Lab 1 Troubleshooting](./labs/lab-01-first-install/README.md)
→ Common: Syntax errors in configuration.nix

### "I don't know what configuration options exist"

→ Resource: https://search.nixos.org/options
→ Search for what you want to configure, copy examples

### "Why is this so complicated?"

→ It's not - you're learning a different paradigm
→ Traditional: "Do these steps in order"
→ NixOS: "Describe what you want"
→ Give yourself time to switch mental models

### "I keep forgetting the Nix language"

→ That's normal - reference [Nix Language Basics](./00-foundation/03-nix-language-basics.md)
→ You'll internalize as you use it
→ 90% of configurations use same 5 patterns

---

## Success Metrics

### After Foundation (1-2 hours)

✅ Can explain what Nix does in 1 sentence
✅ Can explain what NixOS does in 1 sentence
✅ Understand 5 Nix language concepts
✅ Know what configuration.nix contains

### After Lab 1 (3-4 hours)

✅ Successfully installed NixOS
✅ Made a configuration change and rebuilt
✅ Rolled back a change successfully
✅ Understand system generations

### After Labs 1-4 (6-8 hours)

✅ Can deploy identical systems on multiple machines
✅ Understand composition and modularity
✅ Can troubleshoot common issues
✅ Ready for small production deployment

### After Complete Course (16+ hours)

✅ NixOS is your default choice for infrastructure
✅ Can build custom configurations
✅ Comfortable with multi-system deployments
✅ Understand advanced patterns (Flakes, etc.)
✅ Contributing to nixpkgs feels achievable

---

## Alternative: Jump to Labs

If you're impatient, you can learn by doing:

1. **Skip** foundation reading
2. **Start** Lab 1: Follow installation steps exactly
3. **As you go**: Each step explains concepts briefly
4. **Then read** foundation materials when confused
5. **Reference** ASCII diagrams for conceptual understanding

This works if you learn by doing, but takes longer.

**Recommendation**: Spend 30 minutes on foundation, then labs.
The foundation makes labs make sense.

---

## Community & Resources

### Official Resources

- https://nixos.org/manual/nixos/stable/ - Official manual
- https://search.nixos.org/options - Configuration options
- https://search.nixos.org/packages - Available packages
- https://github.com/NixOS/nixpkgs - Package repository
- https://github.com/NixOS/nixos-hardware - Hardware configurations

### Community Resources

- r/NixOS - Reddit community
- NixOS Discourse - Forums
- NixOS Discord - Real-time chat
- Nix YouTube - Videos and talks

### When You Get Stuck

1. Search: https://search.nixos.org/options
2. Check: https://github.com/NixOS/nixpkgs/tree/nixos-*/nixos/modules
3. Search: Existing GitHub issues for similar problem
4. Ask: NixOS Discourse with error messages

---

## Next Steps After This Course

1. **Deploy to production** - Use what you learned for real systems
2. **Write custom packages** - Package tools you use
3. **Contribute to nixpkgs** - Help maintain ecosystem
4. **Learn Flakes** - Modern Nix features
5. **Explore ecosystem** - Home Manager, nix-shells, etc.

---

## Your Roadmap

```
TODAY (2-4 hours)
  ├─ Read foundation materials
  ├─ Complete Lab 1
  └─ Have working NixOS running

WEEK 1 (4-6 hours)
  ├─ Complete Labs 2-3
  ├─ Write custom configurations
  └─ Deploy second system

WEEK 2-4 (8-12 hours)
  ├─ Complete Labs 4-6
  ├─ Multi-system deployments
  └─ Production planning

MONTH 2+ (Advanced)
  ├─ Flakes and advanced patterns
  ├─ Custom packages
  └─ nixpkgs contributions
```

**Remember**: You don't need to learn everything at once.
Start with what you need, expand from there.

---

## One More Thing

The hardest part of learning NixOS isn't the syntax.
It's changing how you think about system configuration.

**Traditional thinking**: "How do I apply these changes?"
**NixOS thinking**: "What should the system be?"

Once you make that mental shift, NixOS becomes natural and powerful.

The foundation materials and labs are designed to help you make that shift.

**You've got this. 🚀**

---

Ready to start?

- [Foundation Materials](./00-foundation/)
- [Lab 1: First Installation](./labs/lab-01-first-install/README.md)
- [ASCII Diagrams](./00-foundation/ascii-diagrams.md)

Pick your path above and get started!
