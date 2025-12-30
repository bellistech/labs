# Labs - Production-Ready SRE/DevOps Educational Materials

This repository contains a comprehensive collection of educational materials, tutorials, and walkthroughs designed to bridge the gap between basic scripting knowledge and enterprise-level infrastructure management. All materials emphasize production-ready implementations with extensive documentation following an "Explain Like I'm 5" (ELI5) methodology.

**Generated with assistance from Claude AI - customized, modified, and completed as SRE/DevOps and systems/networking/software engineering refresher**

---

## 📚 Course Collections

### Core Infrastructure Automation

**[Terraform for Infrastructure as Code](./terraform/)**
- Comprehensive guide to Terraform fundamentals and advanced patterns
- State management best practices
- Module design and reusability
- Real-world examples with self-registering, configuration-driven systems
- Production deployment patterns

**[Ansible Configuration Management](./ansible/)**
- Playbook design and best practices
- Dynamic inventory management
- Role-based architecture
- Integration with infrastructure as code workflows
- Hands-on labs and examples

**[Bash Scripting Essentials](./bash/)**
- From basics to production-grade scripts
- Error handling and logging patterns
- Shell script architecture
- Integration with infrastructure automation
- Practical examples and utilities

### Container Orchestration

**[Kubernetes Comprehensive Course](./kubernetes/)**
An 8-part deep dive into Kubernetes:
- Kubernetes fundamentals (comparing to beehives for intuitive understanding)
- Deployments, StatefulSets, and DaemonSets
- Services and networking
- Storage and persistent volumes
- ConfigMaps, Secrets, and configuration management
- Resource management and scheduling
- Advanced patterns and production readiness
- Multi-cluster and GitOps integration

### Language & Application Frameworks

**[Go Programming for Metrics Collection](./go/)**
- Building production-grade metrics collection systems
- gRPC and Protocol Buffers fundamentals
- Service architecture and patterns
- Self-registering metric collectors
- Integration with Prometheus and observability stacks
- From basics to distributed system patterns

### System Administration & Configuration

**[NixOS Comprehensive Curriculum](./nixos/)**
A 24-file modular learning path:
- NixOS fundamentals (comparing to restaurants for understanding)
- Package management and expressions
- System configuration with Nix language
- Reproducible environments and home-manager
- Modular system design
- Hands-on labs and practical deployments
- Development environments and containers
- Advanced topics and production patterns

### Monitoring & Observability

**[Prometheus & Grafana Integration](./monitoring/)**
- Prometheus scraping and storage
- Time series data fundamentals
- Building effective dashboards with Grafana
- Alert rules and notification workflows
- Integration with application metrics
- Production monitoring patterns

**[Logging Architecture & ELK Stack](./logging/)**
- Centralized logging fundamentals
- Elasticsearch, Logstash, and Kibana setup
- Log aggregation patterns
- Query optimization and troubleshooting
- Integration with container platforms
- Production-scale implementations

### Modern Deployment Patterns

**[GitOps with Flux & ArgoCD](./gitops/)**
- GitOps fundamentals and philosophy
- Flux CD configuration and deployment
- ArgoCD for application delivery
- Git workflows for infrastructure
- Multi-environment management
- Production rollout patterns

**[Networking Fundamentals for Infrastructure](./networking/)**
- TCP/IP essentials and troubleshooting
- DNS architecture and management
- Network design patterns
- Service mesh concepts
- Load balancing and traffic management
- Production network design

### Advanced & Emerging Topics

**[LLM Deployment Walkthroughs](./llm-deployment/)**
- Cloud-hosted LLM deployments
- Local deployment on consumer hardware
- Integration with infrastructure platforms
- Model serving patterns
- Cost optimization and scaling
- Practical examples for production and development

**[eBPF & Kernel-Level Observability with Rust](./ebpf/)**
- eBPF fundamentals and capabilities
- Rust for systems programming
- Kernel-level tracing and monitoring
- Performance analysis at scale
- Custom observability tools

**[Databases for Backend Engineers](./databases/)**
- Relational database design and optimization
- NoSQL patterns and tradeoffs
- Backup and recovery strategies
- Replication and high availability
- Query optimization and monitoring
- Production database patterns

---

## 🎯 Learning Paths

### Quick Start (Weekend Refresher)
Ideal for experienced engineers returning to the field or looking for specific skills:

1. Start with **Bash Scripting Essentials** for core Unix/Linux skills
2. Explore **Terraform for Infrastructure as Code** for modern infrastructure patterns
3. Reference **Networking Fundamentals** as needed

**Estimated time:** 8-16 hours

### Intermediate Path (2-4 Weeks)
Comprehensive foundational knowledge:

1. **Bash Scripting Essentials** - master shell environments
2. **Terraform for Infrastructure as Code** - infrastructure foundations
3. **Kubernetes Comprehensive Course** - container orchestration
4. **Monitoring & Observability** - observability tooling
5. **Networking Fundamentals** - infrastructure networking

**Estimated time:** 40-60 hours

### Advanced Path (6-8 Weeks)
Deep expertise across the full stack:

Complete the Intermediate Path, then add:
1. **Ansible Configuration Management** - advanced automation
2. **NixOS Comprehensive Curriculum** - declarative system design
3. **Go Programming for Metrics** - application-level instrumentation
4. **GitOps with Flux & ArgoCD** - modern deployment patterns
5. **Databases for Backend Engineers** - data layer mastery

**Estimated time:** 120-160 hours

### Specialized Tracks

**Cloud-Native Specialist:**
- Kubernetes Comprehensive Course
- GitOps with Flux & ArgoCD
- LLM Deployment Walkthroughs
- eBPF & Kernel-Level Observability

**Systems Engineering Specialist:**
- NixOS Comprehensive Curriculum
- eBPF & Kernel-Level Observability
- Networking Fundamentals
- Databases for Backend Engineers

**Automation & IaC Specialist:**
- Terraform for Infrastructure as Code
- Ansible Configuration Management
- Bash Scripting Essentials
- GitOps with Flux & ArgoCD

---

## ✨ Philosophy & Approach

All materials in this repository follow core principles:

**Production-Ready from Day One**
- Every example is designed to work in real environments
- Complete project scaffolding with industry-standard structures
- No toy implementations or oversimplified demonstrations
- Self-registering, configuration-driven architectures that avoid hardcoding

**Accessible Without Sacrificing Depth**
- ELI5 (Explain Like I'm 5) methodology with practical analogies
- Extensive inline documentation explaining both "what" and "why"
- Progressive complexity from fundamentals to advanced patterns
- Conceptual understanding paired with hands-on implementation

**Comprehensive Learning Resources**
- Multiple learning paths accommodating different time commitments
- Detailed architecture explanations with working code
- Real-world tradeoffs and production considerations
- Complete file structures ready for deployment in training scenarios

**Iterative & Feedback-Driven**
- Materials evolve based on practical experience
- Architectural patterns improve as better approaches emerge
- Quality verification through content metrics and practical testing
- Community-focused improvements and refinements

---

## 🚀 Getting Started

1. **Identify your learning path** above based on your experience level and goals
2. **Clone or download** the specific course you want to explore
3. **Follow the course README** for setup instructions and prerequisites
4. **Work through materials progressively** - each section builds on previous knowledge
5. **Implement the examples** in your own lab environment
6. **Adapt and modify** for your specific use cases and infrastructure

All materials are structured to be immediately usable in training, onboarding, and professional development scenarios.

---

## 📋 Directory Structure

```
labs/
├── README.md (this file)
├── terraform/
├── ansible/
├── bash/
├── kubernetes/
├── go/
├── nixos/
├── monitoring/
├── logging/
├── gitops/
├── networking/
├── llm-deployment/
├── ebpf/
└── databases/
```

Each directory contains a complete course with:
- Comprehensive `README.md` with setup instructions
- Progressive lesson modules
- Working code examples and configurations
- Hands-on labs and exercises
- Additional resources and references

---

## 🛠️ Tools & Technologies

The materials cover modern infrastructure and backend engineering tools:

- **Infrastructure as Code:** Terraform, Ansible
- **Container Platforms:** Docker, Kubernetes
- **Programming Languages:** Go, Bash, Rust
- **System Configuration:** NixOS
- **Monitoring & Observability:** Prometheus, Grafana
- **Logging:** ELK Stack (Elasticsearch, Logstash, Kibana)
- **Deployment Patterns:** GitOps (Flux, ArgoCD)
- **Advanced Topics:** eBPF, LLM Infrastructure, Kernel Observability

---

## 🏗️ Complete Repository Scaffolding Map

The full repository structure is automatically created using the included `build-structure.sh` script:

```
labs/
├── README.md                              # Main repository overview
├── CONTRIBUTING.md                        # Contribution guidelines
├── LICENSE.md                             # MIT License
├── CHANGELOG.md                           # Version history and updates
├── courses.json                           # Machine-readable course metadata
├── build-structure.sh                     # Automated structure builder
├── .gitignore                             # Git ignore rules
├── .github/
│   └── workflows/                         # GitHub Actions (future)
│
├── terraform/                             # Infrastructure as Code
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ (lesson files + examples/)
│   ├── 02-core-skills/ (lesson files + examples/)
│   ├── 03-advanced-patterns/ (lesson files + examples/)
│   ├── 04-labs/ (lab files + solutions/)
│   ├── 05-reference/ (reference docs)
│   └── scripts/ (setup.sh, validate-examples.sh)
│
├── kubernetes/                            # Container Orchestration (8-part)
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 07-reference/ (with examples/)
│   ├── 06-labs/ (with solutions/)
│   └── scripts/ (setup-minikube.sh, cleanup.sh)
│
├── bash/                                  # Core Scripting Skills
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (setup.sh, test-runner.sh)
│
├── go/                                    # Metrics & Backend Services
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-systems-integration/ (with examples/)
│   ├── 04-labs/ (with solutions/)
│   ├── 05-reference/
│   └── scripts/ (setup.sh, run-tests.sh)
│
├── nixos/                                 # System Configuration (24-part)
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 07-reference/ (with examples/)
│   ├── 06-labs/ (with solutions/)
│   └── scripts/ (setup.sh, validate.sh)
│
├── ansible/                               # Configuration Management
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (setup.sh)
│
├── monitoring/                            # Prometheus & Grafana
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (docker-compose.yml, setup.sh)
│
├── logging/                               # ELK Stack & Log Aggregation
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (docker-compose.yml, setup.sh)
│
├── gitops/                                # Flux & ArgoCD
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (setup.sh, validate.sh)
│
├── networking/                            # Infrastructure Networking
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (diagnostic-tools.sh)
│
├── llm-deployment/                        # AI Infrastructure
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 05-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (setup-linux.sh, setup-windows.sh, benchmark.sh)
│
├── ebpf/                                  # eBPF & Kernel Observability
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── 01-fundamentals/ through 04-labs/ (with examples/)
│   ├── 05-reference/
│   └── scripts/ (setup.sh, build-tools.sh)
│
└── databases/                             # Data Infrastructure
    ├── README.md
    ├── QUICKSTART.md
    ├── 01-fundamentals/ through 05-labs/ (with examples/)
    ├── 05-reference/
    └── scripts/ (setup.sh, backup-automation.sh)
```

### Building the Structure

To create this entire directory structure automatically:

```bash
# Build in current directory
bash build-structure.sh

# Build in specific directory
bash build-structure.sh /path/to/labs

# Show detailed output
bash build-structure.sh . --verbose

# Get help
bash build-structure.sh --help
```

**Output:**
- 14 complete course directories
- 70+ lesson templates
- 50+ example directories (ready for working code)
- 14+ lab directories with solution templates
- 5+ reference document templates per course
- Utility scripts for each course

---

## 📝 Attribution

These educational materials were generated with assistance from Claude AI and have been customized, modified, and completed as comprehensive SRE/DevOps and systems/networking/software engineering refresher materials. The content represents a synthesis of industry best practices, production experience, and modern infrastructure patterns.

---

## 🤝 Using These Materials

These resources are designed for:
- **Individual learning:** Self-paced skill development
- **Team training:** Onboarding engineers to infrastructure practices
- **Organizational reference:** Implementation patterns and best practices
- **Interview preparation:** Deep technical preparation for backend/SRE roles

All examples are working implementations suitable for:
- Immediate deployment in training environments
- Adaptation for your specific infrastructure
- Reference implementations for your own projects
- Foundation for advanced customization

---

## 📚 Additional Notes

- All materials assume comfort with command-line interfaces and basic Linux/Unix concepts
- Examples use real-world configurations and patterns
- Heavily commented code and configurations aid learning
- Each course includes multiple time commitment options
- Materials emphasize understanding over memorization

---

**Last Updated:** December 2024

**Maintainer:** Built with Claude AI assistance - customized, modified, and completed by the author

---

*Start your infrastructure mastery journey today. Choose your learning path and dive deep into production-ready SRE/DevOps knowledge.*
