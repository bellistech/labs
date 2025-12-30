# Labs

Production-focused labs and reference material for SRE/DevOps, systems engineering, networking, and software tooling. Content is organized into domain areas with consistent course scaffolding (fundamentals → core skills → advanced patterns → labs → reference).

Maintained by Stevie Bellis. Some material was generated with AI assistance and then curated and edited.

## Repository layout

Top-level directories:

- `automation_devops_sre/` – Infrastructure automation and SRE platform topics (Terraform, Ansible, Kubernetes, monitoring, logging, eBPF, NixOS, Jenkins, Puppet)
- `coding/` – Programming tracks for automation and systems tooling (Bash, Go, Python, Rust)
- `networking/` – Networking course plus operational references
- `databases/` – Database fundamentals, operations, and reliability patterns
- `linux/` – Linux-focused notes and labs
- `ml_llm_ai/` – ML/LLM infrastructure topics and labs

Supporting files:

- `SETUP.md` – environment notes and prerequisites
- `courses.json` – machine-readable metadata
- `build-structure.sh` – scaffolding/structure generator
- `.github/workflows/` – CI workflow definitions (if/when enabled)
- `CONTRIBUTING.md`, `LICENSE.md`, `CHANGELOG.md`

## Course structure conventions

Most course modules follow a common layout:

```
<course>/
├── README.md
├── QUICKSTART.md
├── 01-fundamentals/
├── 02-core-skills/
├── 03-advanced-patterns/
├── 04-labs/
├── 05-reference/
└── scripts/
```

Some tracks use a numbered-document format (e.g., `1-`, `2-`, `3-`) with `examples/` folders, and some include additional `projects/` or `docs/` directories.

## Index

### Automation / DevOps / SRE

Located under `automation_devops_sre/`:

- Terraform: `automation_devops_sre/terraform/`
- Ansible: `automation_devops_sre/ansible/`
- Kubernetes: `automation_devops_sre/kubernetes/` (includes `projects/` and `docs/`)
- Monitoring: `automation_devops_sre/monitoring/`
- Logging: `automation_devops_sre/logging/`
- eBPF: `automation_devops_sre/ebpf/`
- NixOS: `automation_devops_sre/nixos/` (includes `labs/` and `reference/`)
- Jenkins: `automation_devops_sre/jenkins/`
- Puppet: `automation_devops_sre/puppet/`

### Coding tracks

Located under `coding/`:

- Bash: `coding/bash/`
- Go: `coding/go/` (includes `docs/`, `examples/`, `scripts/`)
- Python: `coding/python/` (includes `docs/`, `exercises/`, `capstone/`)
- Rust: `coding/rust/` (includes `ebpf-sidecar/` workspace and docs)

### Networking

- Core course: `networking/`
  - `01-fundamentals/`, `02-core-skills/`, `03-advanced-patterns/`, `04-labs/`, `05-reference/`, `scripts/`
- Additional operational references live alongside the course material (for example: IPv6, BGP convergence, Linux networking stack notes).

### Databases

- Core course: `databases/`
  - `01-fundamentals/`, `02-core-skills/`, `03-advanced-patterns/`, `04-labs/`, `05-reference/`, `scripts/`

### Linux

- Linux notes and labs: `linux/`

### ML / LLM / AI

- Core course: `ml_llm_ai/`
  - `01-fundamentals/`, `02-core-skills/`, `03-advanced-patterns/`, `04-labs/`, `05-reference/`, `scripts/`

## Getting started

1. Read `SETUP.md` for prerequisites and recommended tooling.
2. Pick an entry point:
   - Platform automation: `automation_devops_sre/terraform/` or `automation_devops_sre/ansible/`
   - Container operations: `automation_devops_sre/kubernetes/`
   - Observability: `automation_devops_sre/monitoring/` and `automation_devops_sre/logging/`
   - Tool-building: `coding/python/` or `coding/go/`
   - Systems depth: `coding/rust/` and `automation_devops_sre/ebpf/`
   - Fundamentals refresh: `networking/`, `linux/`, `databases/`
3. Follow the module `README.md`, then `QUICKSTART.md`, then proceed through numbered lessons and labs.

## Suggested learning paths

Weekend refresher:
- `coding/bash/`
- `networking/`
- One automation module: `automation_devops_sre/terraform/` or `automation_devops_sre/ansible/`

Platform foundations (2–4 weeks):
- Terraform or Ansible
- Kubernetes
- Monitoring + Logging
- Networking fundamentals

Operations depth (6–8 weeks):
- Add databases and Linux operations
- Add one programming track for tooling (Python or Go)
- Add eBPF/Rust for kernel-level observability work

## Contributing

See `CONTRIBUTING.md`. Small, incremental improvements are preferred: clarity, correctness, runnable examples, and repeatable lab environments.

Last updated: December 2025
