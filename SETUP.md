# Labs Repository - Quick Setup Guide

Welcome to the Labs repository bundle! This guide will get you up and running in minutes.

## 📦 What's Included

```
labs/
├── README.md                 # Main documentation and learning paths
├── CONTRIBUTING.md           # How to contribute and add content
├── LICENSE.md                # MIT License
├── CHANGELOG.md              # Version history and roadmap
├── .gitignore                # Git ignore rules
├── courses.json              # Machine-readable course metadata
├── build-structure.sh        # Automated repository builder
└── SETUP.md                  # This file
```

## 🚀 Quick Start (5 minutes)

### Step 1: Extract the Bundle

```bash
# If you have a zip file
unzip labs-bundle.zip
cd labs

# If you have a tar.gz file
tar xzf labs-bundle.tar.gz
cd labs
```

### Step 2: Build the Directory Structure

The bundle includes a bash script that automatically creates the entire repository structure:

```bash
# Build in current directory
bash build-structure.sh .

# Or build in a specific location
bash build-structure.sh /path/to/labs

# Show detailed output
bash build-structure.sh . --verbose
```

**Expected Output:**
```
[SUCCESS] Starting Labs repository structure builder
[INFO] Building course structure: terraform
[INFO] Building course structure: ansible
...
[SUCCESS] Repository structure created successfully!

Summary:
  Directories created: 155
  Files created: 265
  Total courses: 14
```

### Step 3: Initialize Git (Optional)

```bash
git init
git add .
git commit -m "Initial Labs repository structure"
git remote add origin https://github.com/yourusername/labs.git
git push -u origin main
```

### Step 4: Start Adding Content

Pick a course and start building:

```bash
cd terraform
# Edit the README.md to add course overview
nano README.md

# Add course content to lesson files
nano 01-fundamentals/lesson-01-terraform-basics.md

# Create working examples
mkdir -p 01-fundamentals/examples/my-first-example
cd 01-fundamentals/examples/my-first-example
# Add your working code here
```

## 📚 Understanding the Structure

### Course Organization

Each course follows this consistent pattern:

```
course-name/
├── README.md                    # Course overview & prerequisites
├── QUICKSTART.md               # 30-minute quick intro
├── 01-fundamentals/            # Basic concepts (ELI5 style)
├── 02-core-skills/             # Practical implementation
├── 03-advanced-patterns/       # Production patterns
├── 04-labs/                    # Hands-on exercises
├── 05-reference/               # Lookup materials & FAQs
└── scripts/                    # Setup & utility scripts
```

### Learning Paths

The repository includes 6 different learning paths optimized for different goals:

1. **Weekend Refresher** (12 hours) - Quick skill refresh
2. **Intermediate SRE Foundation** (50 hours) - 2-4 weeks of learning
3. **Advanced SRE Mastery** (150 hours) - Comprehensive depth
4. **Cloud-Native Specialist** (70 hours) - Focus on cloud tech
5. **Systems Engineering Specialist** (90 hours) - Deep systems focus
6. **Automation & IaC Specialist** (70 hours) - Infrastructure focus

See `README.md` for details on each path.

## 📋 Available Courses

### Foundation Courses
- **bash** - Shell scripting essentials
- **networking** - Infrastructure networking fundamentals

### Infrastructure Automation
- **terraform** - Infrastructure as Code
- **ansible** - Configuration management

### Container & Orchestration
- **kubernetes** - 8-part comprehensive course
- **docker** - Container fundamentals (future)

### Programming & Development
- **go** - Metrics collection & backend services

### System Administration
- **nixos** - 24-part declarative system configuration

### Observability
- **monitoring** - Prometheus & Grafana
- **logging** - ELK Stack & log aggregation

### Deployment & Operations
- **gitops** - Flux & ArgoCD for GitOps
- **llm-deployment** - AI/ML infrastructure

### Advanced Topics
- **ebpf** - Kernel-level observability with Rust
- **databases** - Database design & operations

## 🎯 Content Creation Guide

### Adding a Lesson

1. Navigate to the course section:
   ```bash
   cd terraform/01-fundamentals
   ```

2. Edit the lesson file:
   ```bash
   nano lesson-01-terraform-basics.md
   ```

3. Use this template:
   ```markdown
   # Lesson: Terraform Basics

   ## Overview
   [What you'll learn in this lesson]

   ## The Concept Explained Simply
   [ELI5 explanation with analogy]

   ## Technical Deep Dive
   [Detailed technical explanation]

   ## Practical Example
   [Working code example]

   ## Common Mistakes
   [What to watch out for]

   ## Next Steps
   [What comes next]
   ```

### Adding Working Examples

1. Create example directory:
   ```bash
   mkdir -p 02-core-skills/examples/deploy-vpc
   cd 02-core-skills/examples/deploy-vpc
   ```

2. Create essential files:
   ```bash
   # Main implementation
   touch main.tf variables.tf outputs.tf
   
   # Example configuration
   touch terraform.tfvars.example
   
   # Setup script
   touch setup.sh
   
   # Documentation
   touch README.md
   ```

3. In README.md, include:
   - What this example demonstrates
   - Prerequisites
   - How to run it
   - Expected output
   - Common modifications

### Adding Labs & Exercises

1. Create lab file:
   ```bash
   cd 04-labs
   nano lab-01-deploy-app-infrastructure.md
   ```

2. Structure:
   - Scenario description
   - Objectives
   - Success criteria
   - Hints (optional)
   - Reference solution in `solutions/`

## 🛠️ Using courses.json

The `courses.json` file contains machine-readable metadata:

```bash
# View course information
cat courses.json | jq '.courses[] | {name, estimatedHours}'

# Get learning paths
cat courses.json | jq '.learningPaths[] | {name, description}'

# Find all courses on a topic
cat courses.json | jq '.topics.kubernetes'
```

This enables:
- Automated course discovery
- Learning path recommendations
- Integration with learning management systems
- Progress tracking tools

## 📖 Key Principles to Follow

When adding content, remember these core principles:

✅ **Production-Ready First**
- Every example must work in real environments
- Include error handling and edge cases
- Use industry best practices

✅ **ELI5 Methodology**
- Start with analogies and simple explanations
- Progress to technical depth
- Don't sacrifice accuracy for simplicity

✅ **Extensive Documentation**
- Inline comments in all code
- Explain the "why" not just the "what"
- Include multiple learning levels

✅ **Real-World Focus**
- Use production scenarios
- Show common pitfalls
- Include troubleshooting guides

## 🔗 Integration with GitHub

### Initial Setup
```bash
git init
git add .
git commit -m "Initial Labs repository structure"
git branch -M main
git remote add origin https://github.com/yourusername/labs.git
git push -u origin main
```

### Recommended .github/workflows

Create `.github/workflows/validate.yml` to:
- Lint markdown files
- Validate code examples
- Check for broken links
- Verify directory structure

### Create a .github/pull_request_template.md

```markdown
## Description
What content is being added/updated?

## Type of Change
- [ ] New course
- [ ] New lesson
- [ ] New example
- [ ] New lab
- [ ] Bug fix
- [ ] Improvement

## Checklist
- [ ] Extensive comments included
- [ ] Examples have been tested
- [ ] No hardcoded paths or secrets
- [ ] README updated
- [ ] Follows ELI5 methodology
```

## 📊 Building Your Own

### Extend with New Courses

1. Create course directory:
   ```bash
   mkdir -p my-course/{01-fundamentals,02-core-skills,03-advanced-patterns,04-labs,05-reference,scripts}
   ```

2. Add documentation:
   ```bash
   touch my-course/README.md
   touch my-course/QUICKSTART.md
   ```

3. Add content following the patterns in existing courses

### Create Learning Paths

Edit `courses.json` to add custom learning paths:

```json
{
  "id": "my-custom-path",
  "name": "My Custom Learning Path",
  "description": "Description of the learning path",
  "estimatedHours": 50,
  "courses": ["course1", "course2", "course3"],
  "difficulty": "intermediate"
}
```

## ❓ Troubleshooting

### build-structure.sh not executable
```bash
chmod +x build-structure.sh
```

### Permission denied when creating directories
```bash
# Ensure you have write permissions
chmod u+w /path/to/labs
```

### Need help with courses.json
```bash
# Validate JSON syntax
cat courses.json | jq . > /dev/null
```

## 📞 Getting Help

1. **Check existing courses** - See how other courses are structured
2. **Read CONTRIBUTING.md** - Detailed guidelines for contributing
3. **Review examples** - Look at working examples in other courses
4. **Test with build-structure.sh** - The script is self-documenting

## 🎉 Next Steps

1. **Build the structure** → `bash build-structure.sh .`
2. **Read the main README** → `cat README.md`
3. **Choose a course** → Pick one that interests you
4. **Start with QUICKSTART.md** → Get oriented quickly
5. **Add your first content** → Create a lesson or example
6. **Contribute back** → Share improvements via pull request

## 📝 Attribution

This repository was generated with assistance from Claude AI and has been customized, modified, and completed as comprehensive SRE/DevOps and systems/networking/software engineering educational materials.

---

**Happy Learning! 🚀**

For the full README and detailed learning paths, see `README.md`.
For contribution guidelines, see `CONTRIBUTING.md`.
For version history and roadmap, see `CHANGELOG.md`.
