# NixOS ELI5 Course - Complete Index

## 📊 Course Status

```
FOUNDATION MATERIALS: ✅ COMPLETE
├── What is Nix? (30 pages)
├── What is NixOS? (30 pages)
├── Nix Language Basics (25 pages)
├── ASCII Diagrams (35 pages)
└── Total: ~120 pages

BEGINNER-INTERMEDIATE LABS: ✅ COMPLETE
├── Lab 1: First Installation (complete with template)
├── Lab 2: Development Environment (20+ pages)
├── Lab 3: Dev Shells (22+ pages)
├── Lab 4: Multi-Host Configuration (28+ pages)
└── Total: ~90 pages

REFERENCE MATERIALS: ✅ COMPLETE
├── Debugging Guide (comprehensive)
└── Total: ~45 pages

NAVIGATION MATERIALS: ✅ COMPLETE
├── Main README
├── QUICKSTART Guide
├── CURRICULUM Overview
├── This Index
└── Total: ~80 pages

TOTAL CURRENT: ~335 pages of production-ready material

FUTURE (Designed, ready to build):
├── Labs 5-9 (Planned: 200+ pages)
├── Advanced Topics (Planned: 80+ pages)
├── Examples (Planned: 60+ pages)
└── Additional Reference (Planned: 40+ pages)
```

---

## 🚀 Quick Navigation

### For Absolute Beginners (First 2 Hours)
1. Start: `README.md`
2. Then: `QUICKSTART.md` → Select "2-Hour Path"
3. Read: `00-foundation/01-what-is-nix.md`
4. Read: `00-foundation/02-what-is-nixos.md`
5. Skim: `00-foundation/03-nix-language-basics.md`
6. Do: `labs/lab-01-first-install/README.md`
7. Result: Working NixOS, understanding basics

### For Developers (4-6 Hours)
1. Foundation: 90 minutes
2. Lab 1: 60 minutes (installation)
3. Lab 2: 45 minutes (dev environment)
4. Lab 3: 45 minutes (dev shells)
5. Result: Can create reproducible dev environments, understand NixOS philosophy

### For Operations/DevOps (8-12 Hours)
1. Foundation: 90 minutes
2. Lab 1: 60 minutes (installation)
3. Lab 4: 90 minutes (multi-host setup)
4. Reference: Debugging guide + examples
5. Result: Can deploy and manage multiple identical NixOS systems

### For Mastery (40-60 Hours)
1. Complete all foundation materials
2. Complete all current labs (1-4)
3. Study reference materials deeply
4. Build real projects with what you learn
5. When ready: Advanced labs (5-9)
6. Result: Production-ready expertise, can contribute to community

---

## 📚 All Materials by Topic

### Foundation & Concepts
- `00-foundation/01-what-is-nix.md` - Package manager deep dive
- `00-foundation/02-what-is-nixos.md` - Operating system layer
- `00-foundation/03-nix-language-basics.md` - Minimal language for config
- `00-foundation/ascii-diagrams.md` - Visual explanations

### Practical Labs (Hands-On)
- `labs/lab-01-first-install/` - Boot, install, rebuild
  - `README.md` - Step-by-step guide
  - `configuration.nix` - Heavily commented template
- `labs/lab-02-dev-environment/` - Project-based shells
  - Multi-language environments
  - Team consistency
  - Virtual environment automation
- `labs/lab-03-dev-shells/` - Quick ad-hoc environments
  - One-liners with `-p`
  - Reusable templates
  - Command line patterns
- `labs/lab-04-multi-host/` - Managing multiple systems
  - Shared vs specific config
  - Deployment automation
  - Fleet management patterns

### Reference & Troubleshooting
- `reference/debugging-nix.md` - Error messages and solutions
  - Syntax errors (typos, brackets)
  - Build failures (hash mismatches, disk space)
  - Package/service issues
  - Debugging workflow
  - Prevention best practices

### Navigation & Meta
- `README.md` - Course overview
- `QUICKSTART.md` - Path selection guide
- `CURRICULUM.md` - Complete curriculum overview
- `STRUCTURE.txt` - Directory structure
- `INDEX.md` - This file

---

## 🎓 Learning Outcomes by Lab

### Lab 1: First Installation
**Time**: 90 minutes
**Outcome**: 
- ✅ Successfully installed NixOS
- ✅ Understand boot/rebuild cycle
- ✅ Can make configuration changes safely
- ✅ Know how to rollback
- ✅ Comfortable with system generations

**Skills Gained**:
- System installation
- Configuration basics
- Rebuild workflow
- Rollback capability
- Risk management

---

### Lab 2: Development Environment
**Time**: 90 minutes
**Outcome**:
- ✅ Create project-specific shell.nix
- ✅ Manage Python, Node.js, database tools
- ✅ Team consistency (same env everywhere)
- ✅ Clean system (nothing installed globally)
- ✅ Virtual environment automation

**Skills Gained**:
- Development workflow
- Language-specific setups
- Team collaboration patterns
- Environment isolation

---

### Lab 3: Development Shells
**Time**: 70 minutes
**Outcome**:
- ✅ Quick ad-hoc shells with -p flag
- ✅ Reusable templates (~/.nix-shells/)
- ✅ One-liners for experimentation
- ✅ Combining multiple tools
- ✅ Custom environment variables

**Skills Gained**:
- Quick experimentation
- Template creation
- Tool discovery
- Workflow optimization

---

### Lab 4: Multi-Host Configuration
**Time**: 100 minutes
**Outcome**:
- ✅ Configure multiple servers identically
- ✅ Shared vs specific configuration
- ✅ Deployment automation
- ✅ Server registry patterns
- ✅ Infrastructure as code

**Skills Gained**:
- Multi-system management
- Configuration composition
- Automation patterns
- Scalable infrastructure

---

## 🔧 Tools & Commands You'll Learn

### Essential Commands
```bash
# Core NixOS commands
nixos-rebuild switch              # Apply configuration
nixos-rebuild dry-build           # Preview changes
nixos-rebuild dry-activate        # What would change
nixos-rebuild switch --rollback   # Revert to previous
nixos-rebuild list-generations    # See all versions

# Development shells
nix-shell                         # Enter with shell.nix
nix-shell -p package             # Enter with ad-hoc package
nix-shell ~/.nix-shells/template.nix  # Use template

# Searching
nix search nixpkgs pattern        # Find packages
nix search nixos option           # Find options

# Maintenance
nix-collect-garbage              # Clean old packages
nix-collect-garbage -d           # Aggressive cleanup
```

### Useful Shortcuts
```bash
# Add to ~/.bashrc or ~/.zshrc:
alias rebuild='sudo nixos-rebuild switch'
alias dry-run='sudo nixos-rebuild dry-build'
alias rollback='sudo nixos-rebuild switch --rollback'
alias generations='sudo nixos-rebuild list-generations'

# Usage:
rebuild        # Instead of full command
dry-run        # Quick check
rollback       # Revert safely
generations    # See history
```

---

## 📖 How to Read This Course

### Linear Path (Recommended First Time)
```
1. README.md (this course)
   ↓
2. QUICKSTART.md (pick your path)
   ↓
3. Foundation materials (01, 02, 03)
   ↓
4. ASCII diagrams (visual understanding)
   ↓
5. Lab 1 (hands-on with real system)
   ↓
6. Your path:
   - Developer: Labs 2, 3 next
   - Operations: Lab 4 next
   - Both: All labs in order
```

### Reference Path (When Needed)
```
1. Run into problem
   ↓
2. Check QUICKSTART.md (is it in FAQ?)
   ↓
3. Search reference/debugging-nix.md
   ↓
4. Search nix search nixpkgs
   ↓
5. Ask community (discourse, reddit, github)
```

### Mastery Path
```
1. Complete all foundation (2-3 hours)
   ↓
2. Do all 4 current labs (6-8 hours)
   ↓
3. Build your own project with what you learned (varies)
   ↓
4. Study reference materials (2-3 hours)
   ↓
5. When advanced labs added (Labs 5-9): take them
   ↓
6. Contribute to nixpkgs or open source
```

---

## 🎯 Use Cases & Recommended Path

### "I want reproducible development environments"
**Start with**: `00-foundation/01-02.md`, then Lab 2 & 3

### "I want to manage servers in code"
**Start with**: All foundation, then Lab 1 & 4

### "I'm curious about how NixOS works"
**Start with**: `ascii-diagrams.md`, then foundation materials

### "I'm evaluating NixOS for my team"
**Start with**: QUICKSTART.md → 4-Hour Developer or Operations path

### "I'm already using NixOS and stuck"
**Start with**: `reference/debugging-nix.md` for your specific error

### "I want to contribute to nixpkgs"
**When available**: Advanced topics → Custom packages section

---

## 🔍 Finding Answers in This Course

### By Problem
- **System won't boot**: Lab 1 troubleshooting + debugging guide
- **Package not available**: Lab 2/3 + debugging guide
- **Multiple servers differ**: Lab 4
- **Don't know Nix syntax**: Foundation 03
- **Build fails mysteriously**: Debugging guide

### By Question
- **"What is Nix?"**: Foundation 01
- **"How do I install NixOS?"**: Lab 1
- **"How do I develop?"**: Lab 2
- **"Why use NixOS?"**: Foundation 02 + ASCII diagrams
- **"How do I manage multiple systems?"**: Lab 4
- **"What does this error mean?"**: Debugging guide

---

## 📝 Course Features

✅ **Heavily Commented Code**
- Every line explained
- Why decisions were made
- Common mistakes highlighted

✅ **Multiple Learning Styles**
- Reading (explanations)
- Diagrams (visual)
- Hands-on (labs)
- Troubleshooting (reference)

✅ **Production-Ready**
- Not toy examples
- Real deployment patterns
- Safe experimentation
- Rollback built in

✅ **Progressive Difficulty**
- Beginner foundation
- Intermediate labs
- Advanced topics (coming)
- Can skip ahead if you know the material

✅ **Comprehensive**
- What to do (labs)
- Why it works (foundation)
- How to debug (reference)
- How to learn more (links)

---

## 🚀 What You Can Do After This Course

### After Foundation Materials
- ✅ Understand NixOS philosophy
- ✅ Know why declarative > imperative
- ✅ Understand package isolation
- ✅ Know how to read NixOS config

### After Lab 1
- ✅ Install NixOS
- ✅ Modify system configuration
- ✅ Rebuild safely
- ✅ Understand generations/rollback

### After Lab 2 & 3
- ✅ Create reproducible dev environments
- ✅ Manage per-project tools
- ✅ Share environment with team
- ✅ Keep system clean (no global installs)

### After Lab 4
- ✅ Manage multiple servers identically
- ✅ Deploy using configuration
- ✅ Update fleet automatically
- ✅ Keep systems in sync

---

## 📚 External Resources

### Official
- [NixOS Manual](https://nixos.org/manual/nixos/stable/)
- [Option Search](https://search.nixos.org/options)
- [Package Search](https://search.nixos.org/packages)

### Learning
- [Nix Pills](https://nixos.org/guides/nix-pills/)
- [Nix Manual](https://nixos.org/manual/nix/stable/)
- [Nixpkgs Manual](https://nixos.org/manual/nixpkgs/stable/)

### Community
- [Discourse](https://discourse.nixos.org/) - Main forum
- [Reddit](https://reddit.com/r/NixOS/) - Community chat
- [GitHub Issues](https://github.com/NixOS/nixpkgs/) - Problem solving
- [Discord](https://discord.gg/RbvHtGU) - Real-time chat

### Tools
- [NixOps](https://nixops.readthedocs.io/) - Declarative deployments
- [Colmena](https://colmena.cli.rs/) - Simple deployment
- [Home Manager](https://github.com/nix-community/home-manager) - User config
- [Flakes](https://nixos.wiki/wiki/Flakes) - Modern Nix

---

## 🎁 What's Included in This Release

### Files
- 14 complete markdown files
- 2 heavily commented Nix templates
- 1 visual directory structure guide
- Complete navigation and learning paths
- Comprehensive debugging reference

### Content
- ~335 pages of material
- 4 complete, walkthrough labs
- 4 foundation modules
- Extensive reference guides
- Production-ready examples

### Coverage
- Package management ✅
- System configuration ✅
- Development workflows ✅
- Multi-host management ✅
- Debugging & troubleshooting ✅
- Learning paths for different roles ✅

---

## 🔄 What's Coming

### Labs 5-9 (Planned)
- Lab 5: Production Server Setup (security, services, monitoring)
- Lab 6: Deployment Automation (CI/CD integration)
- Lab 7: Home Manager (user dotfile management)
- Lab 8: Flakes (modern Nix, reproducibility)
- Lab 9: Hybrid Infrastructure (NixOS + Cloud + Terraform)

### Advanced Topics
- Custom package creation
- Package overlays and customization
- Module system deep dive
- Contributing to nixpkgs
- NixOps/Colmena deployment tools

### Examples
- Minimal system
- Development workstation
- Full-stack web server
- DevOps toolkit
- Docker integration

---

## 💬 Feedback & Contributions

This course is designed to grow. If you:
- Found something confusing
- Want to add your own lab
- Have a better explanation
- Want to contribute examples

See CURRICULUM.md for how to contribute!

---

## 📊 Progress Tracker

Use this to track your progress through the course:

```
Foundation Materials:
  [ ] What is Nix?
  [ ] What is NixOS?
  [ ] Nix Language Basics
  [ ] ASCII Diagrams (at least diagrams 1, 2, 5)

Lab 1: First Installation
  [ ] Read prerequisites
  [ ] Install NixOS
  [ ] Make first change
  [ ] Rollback successfully
  [ ] Verify everything works

Lab 2: Development Environment
  [ ] Understand shell.nix pattern
  [ ] Create simple project
  [ ] Create multi-language project
  [ ] Test virtualenv automation
  [ ] Teach it to a teammate

Lab 3: Development Shells
  [ ] Try nix-shell -p
  [ ] Create template files
  [ ] Use templates successfully
  [ ] Create custom combination
  [ ] Add shell functions to ~/.bashrc

Lab 4: Multi-Host
  [ ] Create base.nix
  [ ] Create server-specific configs
  [ ] Dry-build each server
  [ ] Deploy to test machine
  [ ] Verify servers are identical

Ready for Production?
  [ ] All labs complete
  [ ] Can debug issues
  [ ] Understand documentation
  [ ] Can teach someone else
  [ ] Ready for Labs 5-9 (when available)
```

---

## 🎓 Graduation Criteria

When can you consider yourself "NixOS fluent"?

✅ **Conceptual Understanding**
- Explain declarative vs imperative in your own words
- Understand why NixOS prevents "works on my machine"
- Know the benefits of functional package management

✅ **Practical Skills**
- Can install and configure NixOS
- Can create development shells
- Can manage multiple systems
- Can troubleshoot common issues

✅ **Production Readiness**
- Use version control for configs
- Understand rollback strategy
- Have deployment process
- Can respond to incidents

When you can answer yes to all three: **You're ready!**

---

## 🚀 Next Steps

1. **Right now**: Pick your learning path from QUICKSTART.md
2. **This week**: Complete foundation materials + Lab 1
3. **Next week**: Labs 2-3 (if developer) or Lab 4 (if ops)
4. **This month**: Build something real with NixOS
5. **Long term**: Watch for Labs 5-9 additions, consider contributing

---

## 📞 Need Help?

1. **Course question**: Check this index
2. **Getting error**: See reference/debugging-nix.md
3. **Learning stuck**: Re-read foundation materials
4. **Feature request**: Check CURRICULUM.md for how to contribute
5. **Community**: discourse.nixos.org or reddit.com/r/nixos

---

**Welcome to NixOS!** 

You're learning one of the most powerful infrastructure tools available. The learning curve is real, but the payoff is huge.

Take it step by step, practice with labs, and don't hesitate to reach out to the community.

Happy learning! 🎉

---

*Last Updated: 2024*  
*Course Version: 1.0 (Foundation + Labs 1-4)*  
*Total Material: ~335 pages*  
*Next Update: When Labs 5-9 added*
