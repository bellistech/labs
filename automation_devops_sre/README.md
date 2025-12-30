# Open-Source Technology Stacks for SRE and DevOps

*Comprehensive Reference Guide*

---

## 1. Monitoring & Observability

Monitoring, metrics collection, logging, and distributed tracing tools essential for maintaining system reliability and performance visibility.

### 1.1 Metrics & Alerting

| Technology | Description |
|------------|-------------|
| **Prometheus** | Leading monitoring and alerting toolkit with time-series database; pull-based metrics collection with PromQL query language |
| **Prometheus Alertmanager** | Manages, routes, deduplicates, groups, and silences alerts from Prometheus |
| **VictoriaMetrics** | Prometheus-compatible TSDB optimized for high-cardinality workloads and long-term storage |
| **Grafana** | Analytics and visualization platform for creating dashboards from Prometheus, Elasticsearch, and other sources |
| **Grafana OnCall** | Open-source on-call management and incident response scheduling |

### 1.2 Logging

| Technology | Description |
|------------|-------------|
| **ELK Stack** (Elasticsearch, Logstash, Kibana) | Powerful suite for centralized logging, searching, and real-time log visualization |
| **Loki** | Grafana's log aggregation system; lightweight, cost-effective alternative to ELK designed for cloud-native |
| **Fluentd / Fluent Bit** | Data collectors and log forwarders optimized for cloud-native environments; CNCF graduated |

### 1.3 Distributed Tracing

| Technology | Description |
|------------|-------------|
| **OpenTelemetry** | Standardized APIs and SDKs for instrumenting applications to generate traces, metrics, and logs; CNCF project |
| **Jaeger** | Open-source distributed tracing system for monitoring microservices; CNCF graduated |
| **Zipkin** | Distributed tracing system for gathering timing data and troubleshooting latency issues |

### 1.4 eBPF-Based Observability

| Technology | Description |
|------------|-------------|
| **Pixie** | eBPF-powered observability for Kubernetes; auto-instrumentation with no code changes required |
| **Cilium Hubble** | Network and security observability built on eBPF; deep visibility into service communication |

---

## 2. Container & Kubernetes Ecosystem

Container runtime, orchestration, and Kubernetes-native tooling for modern application deployment and management.

### 2.1 Container Runtime & Orchestration

| Technology | Description |
|------------|-------------|
| **Docker** | Containerization platform packaging applications and dependencies for portability |
| **Kubernetes (K8s)** | De facto standard for automating deployment, scaling, and management of containerized applications |
| **containerd** | Industry-standard container runtime; CNCF graduated |
| **Podman** | Daemonless container engine; rootless containers and Docker CLI compatibility |

### 2.2 Kubernetes Networking

| Technology | Description |
|------------|-------------|
| **Cilium** | eBPF-based networking, security, and observability for Kubernetes; CNCF graduated |
| **Calico** | Popular CNI for network policy enforcement and pod networking |
| **MetalLB** | Load balancer implementation for bare-metal Kubernetes clusters |
| **CoreDNS** | Flexible DNS server for Kubernetes service discovery; CNCF graduated |

### 2.3 Service Mesh

| Technology | Description |
|------------|-------------|
| **Istio** | Full-featured service mesh providing traffic management, security (mTLS), and observability |
| **Linkerd** | Lightweight, security-focused service mesh; CNCF graduated; simpler than Istio |
| **Envoy** | High-performance proxy used as data plane in Istio and other service meshes; CNCF graduated |

### 2.4 Kubernetes Management & Package Management

| Technology | Description |
|------------|-------------|
| **Helm** | Package manager for Kubernetes; simplifies complex deployments using templated charts |
| **Kustomize** | Configuration customization without templating; native kubectl integration |
| **kubectl** | Essential CLI for interacting with Kubernetes clusters |
| **Rancher** | Multi-cluster management platform with GitOps via Fleet feature |
| **k9s** | Terminal-based UI for managing Kubernetes clusters |

### 2.5 Kubernetes Storage

| Technology | Description |
|------------|-------------|
| **Rook/Ceph** | Cloud-native distributed storage orchestrator for Kubernetes |
| **Longhorn** | Lightweight, reliable distributed block storage for Kubernetes |
| **OpenEBS** | Container-attached storage for Kubernetes workloads |

---

## 3. Infrastructure as Code & Configuration Management

Declarative infrastructure provisioning and configuration management automation tools.

### 3.1 Infrastructure Provisioning (IaC)

| Technology | Description |
|------------|-------------|
| **Terraform** | Industry-standard IaC tool for declaratively provisioning multi-cloud infrastructure |
| **OpenTofu** | Open-source Terraform fork; community-driven alternative after HashiCorp license change |
| **Pulumi** | IaC using standard programming languages (Python, Go, TypeScript, etc.) |
| **Crossplane** | Kubernetes-native control plane for building internal cloud platforms |

### 3.2 Configuration Management

| Technology | Description |
|------------|-------------|
| **Ansible** | Agentless automation using YAML playbooks; SSH-based push model |
| **Puppet** | Model-driven, declarative approach with agent-based pull model |
| **Chef** | Configuration automation using Ruby DSL recipes and cookbooks |
| **Salt (SaltStack)** | High-speed execution with push/pull hybrid; ZeroMQ-based communication |

---

## 4. CI/CD & GitOps

Continuous integration, continuous delivery, and GitOps tools for automated software delivery pipelines.

### 4.1 CI/CD Platforms

| Technology | Description |
|------------|-------------|
| **Jenkins** | Highly extensible automation server; thousands of plugins; self-hosted |
| **Tekton** | Cloud-native CI/CD framework running as Kubernetes CRDs |
| **Drone CI** | Container-native CI/CD platform with simple YAML configuration |
| **Concourse CI** | Pipeline-based CI/CD with containerized builds and declarative configuration |

### 4.2 GitOps & Continuous Delivery

| Technology | Description |
|------------|-------------|
| **Argo CD** | Declarative GitOps CD tool; syncs cluster state from Git with intuitive UI |
| **Flux CD** | CNCF graduated GitOps toolkit; Helm and Kustomize integration |
| **Spinnaker** | Multi-cloud CD platform; advanced deployment strategies (canary, blue/green) |

### 4.3 Artifact Management

| Technology | Description |
|------------|-------------|
| **Harbor** | Cloud-native container registry with vulnerability scanning; CNCF graduated |
| **Nexus Repository** | Universal artifact repository for Docker, npm, Maven, and more |
| **JFrog Artifactory (OSS)** | Binary repository manager with broad format support |

---

## 5. Security & Policy

Container security, secrets management, policy enforcement, and vulnerability scanning tools.

### 5.1 Container Security & Scanning

| Technology | Description |
|------------|-------------|
| **Trivy** | Comprehensive vulnerability scanner for containers, IaC, and more |
| **Falco** | Runtime security and threat detection for containers and Kubernetes; CNCF graduated |
| **Grype** | Vulnerability scanner for container images and filesystems |
| **Clair** | Static vulnerability analysis for container images |

### 5.2 Policy & Admission Control

| Technology | Description |
|------------|-------------|
| **OPA / Gatekeeper** | Policy-as-code for Kubernetes admission control; CNCF graduated |
| **Kyverno** | Kubernetes-native policy engine; simpler YAML-based policies |
| **Kubewarden** | Policy engine using WebAssembly for flexible, portable policies |

### 5.3 Secrets Management

| Technology | Description |
|------------|-------------|
| **HashiCorp Vault** | Industry-standard secrets management; dynamic secrets, encryption as a service |
| **External Secrets Operator** | Syncs secrets from external providers (AWS, GCP, Vault) into Kubernetes |
| **Sealed Secrets** | Encrypts secrets for safe storage in Git; decrypted only in-cluster |
| **SOPS** | Mozilla's encrypted file editor; integrates with KMS providers |

---

## 6. Reliability Engineering

Chaos engineering, SLO management, and incident response tools for building resilient systems.

### 6.1 Chaos Engineering

| Technology | Description |
|------------|-------------|
| **Chaos Mesh** | Cloud-native chaos engineering platform for Kubernetes; CNCF incubating |
| **LitmusChaos** | Kubernetes-native chaos engineering with experiment library; CNCF incubating |
| **Gremlin (Free Tier)** | Chaos-as-a-service with controlled failure injection |

### 6.2 SLO Management

| Technology | Description |
|------------|-------------|
| **Sloth** | SLO generator for Prometheus; creates recording rules and alerts from SLO specs |
| **OpenSLO** | Open standard for defining SLOs in a vendor-neutral YAML format |
| **Pyrra** | SLO dashboard and alerting built on Prometheus metrics |

### 6.3 Incident Response & Runbook Automation

| Technology | Description |
|------------|-------------|
| **Rundeck** | Runbook automation and self-service operations; job scheduling and orchestration |
| **StackStorm** | Event-driven automation platform for ChatOps and incident response |
| **PagerDuty (Free Tier)** | Incident management and on-call scheduling |

---

## 7. Developer Experience & Platform Engineering

Internal developer portals, service catalogs, and platform tooling for improving developer productivity.

| Technology | Description |
|------------|-------------|
| **Backstage** | CNCF incubating developer portal; service catalog, docs, and plugin ecosystem |
| **Port** | Internal developer portal for self-service and visibility |
| **Crossplane** | Build custom cloud platforms on Kubernetes; infrastructure composition |

---

## 8. Build & Test Automation

Build tools, testing frameworks, and code quality analysis for automated software delivery.

### 8.1 Build Tools

| Technology | Description |
|------------|-------------|
| **Maven** | Build automation for Java; dependency management and lifecycle |
| **Gradle** | Flexible build tool supporting Java, Kotlin, and polyglot projects |
| **Bazel** | Google's build system for large-scale, multi-language monorepos |

### 8.2 Testing & Quality

| Technology | Description |
|------------|-------------|
| **Selenium** | Browser automation for functional and E2E testing |
| **Apache JMeter** | Load and performance testing for applications and APIs |
| **SonarQube** | Continuous code quality inspection; static analysis for bugs and vulnerabilities |
| **k6** | Modern load testing tool with JavaScript scripting; Grafana project |

---

## 9. Version Control & Collaboration

Source control, project management, and collaboration platforms.

| Technology | Description |
|------------|-------------|
| **Git** | Distributed version control system; industry standard |
| **GitLab (CE)** | Complete DevOps platform with built-in CI/CD, registry, and issue tracking |
| **Gitea** | Lightweight self-hosted Git service; simple and resource-efficient |
| **Forgejo** | Community fork of Gitea with strong governance |

---

## 10. Example CI/CD Pipeline: Jenkins, Docker, Kubernetes

Reference pipeline architecture demonstrating integration of common tools in a production workflow.

| Stage | Actions | Tools |
|-------|---------|-------|
| **1. Source** | Pipeline triggered on Git push; code checkout | Git, Webhooks |
| **2. Build** | Compile code, run unit tests, static analysis | Maven/Gradle, SonarQube |
| **3. Package** | Build container image from Dockerfile | Docker, Kaniko |
| **4. Scan** | Vulnerability scan of container image | Trivy, Grype |
| **5. Push** | Push image to container registry | Harbor, Docker Registry |
| **6. Deploy (Staging)** | Deploy to staging via Helm/Kustomize | Kubernetes, Argo CD |
| **7. Test** | Run integration and E2E tests | Selenium, k6 |
| **8. Deploy (Prod)** | Promote to production with approval gate | Argo CD, Spinnaker |
| **9. Monitor** | Collect metrics, logs, traces; alert on anomalies | Prometheus, Grafana, Loki |

---

## 11. Example Kubernetes Deployment Manifest

Standard Kubernetes Deployment resource demonstrating pod management with replicas.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

| Field | Description |
|-------|-------------|
| `apiVersion` | Kubernetes API version (apps/v1 for Deployments) |
| `kind` | Resource type being created (Deployment) |
| `metadata` | Object identification: name, labels, annotations |
| `spec.replicas` | Number of pod instances to maintain |
| `spec.selector` | Label selector to identify managed pods |
| `spec.template` | Pod template defining container specs |

---

## 12. Helm Usage Patterns

Common use cases for Helm as the Kubernetes package manager:

- **Managing Complexity** — Bundles multiple YAML files into a single versioned chart package
- **Environment Standardization** — Templating with values.yaml enables consistent deploys across dev/staging/prod
- **Third-Party Software** — Artifact Hub provides thousands of pre-configured charts (Prometheus, PostgreSQL, NGINX)
- **Version Control & Rollbacks** — Every deployment is a versioned release; instant rollback with `helm rollback`
- **Dependency Management** — Charts can declare dependencies on subcharts for ordered deployment
- **CI/CD Integration** — Clean CLI commands enable automated, repeatable deployments in pipelines

---

## 13. Quick Reference: Tool Selection Guide

Decision matrix for common tooling choices based on use case and environment.

| Use Case | Primary Choice | Alternative |
|----------|----------------|-------------|
| Container Orchestration | Kubernetes | Docker Swarm, Nomad |
| IaC (Multi-cloud) | Terraform / OpenTofu | Pulumi, Crossplane |
| Config Management | Ansible | Puppet, Salt |
| CI/CD (Self-hosted) | Jenkins, Tekton | Drone, Concourse |
| GitOps | Argo CD | Flux CD |
| Metrics | Prometheus + VictoriaMetrics | Datadog (commercial) |
| Logging | Loki + Grafana | ELK Stack |
| Tracing | Jaeger + OpenTelemetry | Zipkin |
| Service Mesh | Istio | Linkerd (simpler) |
| Secrets | Vault + External Secrets | Sealed Secrets, SOPS |
| Policy | OPA/Gatekeeper | Kyverno (simpler) |
| Vulnerability Scanning | Trivy | Grype, Clair |
| Chaos Engineering | Chaos Mesh | LitmusChaos |
| K8s Networking (eBPF) | Cilium | Calico |

---

*SRE/DevOps Reference • bellis.tech*
