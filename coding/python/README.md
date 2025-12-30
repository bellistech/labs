# 🐍 Python Learning Path: Beginner to Production

## Explain Like I'm 5: Why Python?

Imagine you're learning to talk. Some languages are really complicated with lots of 
grammar rules (that's like C++ or Rust). Python is like learning a language that's 
almost like English - you can start having real conversations in just a few days!

Python is perfect for:
- **Scripting** - Automate boring stuff
- **Web Development** - Build APIs and websites
- **Data Analysis** - Crunch numbers and make charts
- **DevOps** - Manage servers and infrastructure
- **AI/ML** - Machine learning and artificial intelligence

---

## 🎯 What You'll Build

By the end of this course, you'll have built:

1. **Network Utilities Module** - IP address validation, subnet calculations
2. **Database Client** - PostgreSQL connection with proper patterns
3. **REST API** - FastAPI web service with full CRUD operations
4. **Production Deployment** - Systemd service, CI/CD pipeline

---

## 📚 Course Structure

| Phase | Days | Topics | Project |
|-------|------|--------|---------|
| **1** | 1-3 | Python Fundamentals | Calculator, Data Explorer |
| **2** | 4-6 | OOP, Errors, Decorators | Custom Classes |
| **3** | 7-9 | Files, HTTP, System | HTTP Client, File Processor |
| **4** | 10-12 | Databases | PostgreSQL Client |
| **5** | 13-16 | FastAPI Web Dev | REST API |
| **6** | 17-20 | Production Deploy | Systemd + CI/CD |

---

## 🚀 Quick Start

```bash
# Create your workspace
mkdir -p ~/dev/python-learning
cd ~/dev/python-learning

# Create virtual environment
# 
# ELI5: A virtual environment is like a sandbox.
# Packages you install here don't affect the rest of your system.
#
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Verify
which python  # Should show: .../venv/bin/python
python --version  # Should be 3.10+

# Start learning!
cd day_01
python hello.py
```

---

## 📁 Directory Structure

```
python-course/
├── README.md                     # You are here!
├── docs/
│   ├── PYTHON_PART1_FUNDAMENTALS.md    # Days 1-3
│   ├── PYTHON_PART2_INTERMEDIATE.md    # Days 4-6
│   ├── PYTHON_PART3_DATAFILES.md       # Days 7-9
│   ├── PYTHON_PART4_DATABASES.md       # Days 10-12
│   ├── PYTHON_PART5_FASTAPI.md         # Days 13-16
│   └── PYTHON_PART6_PRODUCTION.md      # Days 17-20
├── day_01/                       # Hello World, Data Types
│   ├── hello.py
│   ├── calculator.py
│   └── data_types.py
├── day_02/                       # Control Flow, Data Structures
│   ├── conditions.py
│   ├── loops.py
│   └── data_structures.py
├── day_03/                       # Functions, Modules
│   ├── functions.py
│   ├── netutils/                 # Custom module
│   │   ├── __init__.py
│   │   └── ip.py
│   └── test_netutils.py
├── network_api/                  # FastAPI Capstone Project
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py               # FastAPI application
│   │   ├── config.py             # Configuration
│   │   ├── database.py           # Database setup
│   │   ├── models.py             # SQLAlchemy models
│   │   ├── schemas.py            # Pydantic schemas
│   │   ├── crud.py               # CRUD operations
│   │   └── routers/
│   │       ├── devices.py
│   │       └── stats.py
│   └── tests/
│       ├── unit/
│       └── integration/
├── deployment/
│   ├── network-api.service       # Systemd unit file
│   ├── config.env.example        # Example configuration
│   └── deploy.sh                 # Deployment script
└── tests/
    └── conftest.py               # Pytest configuration
```

---

## 🎓 Teaching Philosophy

Every code example in this course follows these principles:

### 1. Heavy Comments

```python
# Every non-obvious line gets explained
def calculate_subnet(ip: str, prefix: int) -> dict:
    """
    Calculate subnet information from IP and prefix.
    
    ELI5: Given an address like "192.168.1.0" and a prefix like 24,
    figure out how many computers can be on this network.
    
    Args:
        ip: The network address (like a street name)
        prefix: The subnet mask length (how big the neighborhood is)
    
    Returns:
        Dictionary with network info
    """
    # The prefix tells us how many bits are for the network
    # 32 - prefix = bits for hosts (individual computers)
    host_bits = 32 - prefix
    
    # 2^host_bits = total addresses, minus 2 for network and broadcast
    usable_hosts = (2 ** host_bits) - 2
    
    return {"network": ip, "prefix": prefix, "usable_hosts": usable_hosts}
```

### 2. Real Analogies

Every concept gets a real-world comparison:
- **Variables** = Labeled boxes
- **Functions** = Factory machines
- **Classes** = Cookie cutters
- **APIs** = Restaurant waiters
- **Databases** = Filing cabinets

### 3. Progressive Complexity

We start simple and add complexity gradually:

```python
# Day 1: Simple
print("Hello!")

# Day 5: More complex
class Device:
    def __init__(self, hostname: str):
        self.hostname = hostname
    
    def ping(self) -> bool:
        return True

# Day 15: Production-ready
@router.post("/devices", response_model=DeviceResponse)
async def create_device(
    device: DeviceCreate,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> DeviceResponse:
    """Create a new device with authentication."""
    return crud.create_device(db, device, current_user.id)
```

---

## 🔧 Prerequisites

### Required Software

```bash
# Python 3.10 or higher
python3 --version

# pip (Python package manager)
pip3 --version

# Git
git --version

# PostgreSQL (for database lessons)
psql --version
```

### Installation

**macOS:**
```bash
brew install python3 git postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install python3 python3-pip python3-venv git postgresql
```

**Fedora/RHEL:**
```bash
sudo dnf install python3 python3-pip git postgresql-server
```

---

## 📖 How to Use This Course

### Daily Routine

1. **Read the theory** - Start with the docs/ file for that phase
2. **Type the code** - Don't copy-paste! Typing helps you learn
3. **Run and experiment** - Change things and see what happens
4. **Break things** - Try to cause errors on purpose, then fix them
5. **Git commit** - Save your progress at the end of each day

### Example Session

```bash
# Morning: Read the theory
cat docs/PYTHON_PART1_FUNDAMENTALS.md | head -200

# Midday: Work through exercises
cd day_01
python hello.py
# Edit, experiment, learn!

# Evening: Commit your work
git add .
git commit -m "Day 1: Hello world and data types complete"
```

---

## 🏆 Capstone Project

The course culminates in a **Network Inventory API** - a real production-ready 
web service that:

- Manages network device inventory (routers, switches, servers)
- Uses FastAPI for the REST API
- PostgreSQL for persistent storage
- Runs as a systemd service
- Has CI/CD with Jenkins
- Can be packaged as .deb or .rpm

```
┌─────────────────────────────────────────────────────────────┐
│                    Network Inventory API                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   ┌──────────┐     ┌───────────┐     ┌────────────────┐     │
│   │  Users   │────▶│  FastAPI  │────▶│  PostgreSQL    │     │
│   │(clients) │     │   (API)   │     │  (database)    │     │
│   └──────────┘     └───────────┘     └────────────────┘     │
│        │                                     │               │
│        │           ┌───────────┐             │               │
│        └──────────▶│  Swagger  │◀────────────┘               │
│                    │  (docs)   │                             │
│                    └───────────┘                             │
│                                                              │
│   Features:                                                  │
│   ✓ CRUD operations for devices                              │
│   ✓ Search and filtering                                     │
│   ✓ Pagination                                               │
│   ✓ Statistics endpoint                                      │
│   ✓ Health checks                                            │
│   ✓ Systemd integration                                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 📝 Git Workflow

### Initial Setup

```bash
cd ~/dev/python-learning
git init
git add .
git commit -m "Initial commit: Python learning path setup"
```

### Daily Commits

```bash
# After each day's work
git add day_XX/
git commit -m "Day X: Brief description of what you learned"

# Example commits:
# "Day 1: Hello world, calculator, data types"
# "Day 3: Functions and custom netutils module"
# "Day 12: PostgreSQL client with connection pooling"
# "Day 16: FastAPI CRUD endpoints complete"
```

### Push to GitHub

```bash
# Create repo on GitHub first, then:
git remote add origin https://github.com/bellistech/python-course.git
git push -u origin main
```

---

## 🔗 Related Projects

After completing this course, continue your learning with:

1. **Go Course** - Learn compiled language with similar ELI5 style
2. **Kubernetes Course** - Deploy your Python apps to containers
3. **Jenkins Setup** - Automate your Python CI/CD
4. **Package Your App** - Create .deb/.rpm from your Python project

---

## 📚 Additional Resources

### Official Documentation
- [Python Docs](https://docs.python.org/3/)
- [FastAPI Docs](https://fastapi.tiangolo.com/)
- [SQLAlchemy Docs](https://docs.sqlalchemy.org/)
- [Pydantic Docs](https://docs.pydantic.dev/)

### Books
- "Automate the Boring Stuff with Python" (free online)
- "Fluent Python" (intermediate/advanced)

### Practice
- [Python Exercises](https://exercism.io/tracks/python)
- [LeetCode](https://leetcode.com/) (algorithm practice)

---

## ✅ Checklist

Use this to track your progress:

- [ ] Day 1: Environment setup, hello world, data types
- [ ] Day 2: Conditionals, loops, data structures
- [ ] Day 3: Functions, modules, netutils package
- [ ] Day 4: Classes and OOP basics
- [ ] Day 5: Error handling, exceptions
- [ ] Day 6: Decorators, generators, context managers
- [ ] Day 7: File I/O, JSON, YAML
- [ ] Day 8: HTTP requests, REST basics
- [ ] Day 9: System commands, subprocess
- [ ] Day 10: SQLite basics
- [ ] Day 11-12: PostgreSQL with psycopg2
- [ ] Day 13-14: FastAPI fundamentals
- [ ] Day 15-16: FastAPI with database
- [ ] Day 17-18: Systemd service deployment
- [ ] Day 19-20: CI/CD pipeline, testing

---

**Happy Learning! 🐍**

*"The best way to learn Python is to build real things."*
