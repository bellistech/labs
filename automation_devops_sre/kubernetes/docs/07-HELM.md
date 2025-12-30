# Kubernetes Crash Course: From Zero to Production

## Part 7: Helm (The Package Manager)

---

## Chapter 36: What is Helm?

### 36.1 The YAML Explosion Problem

Deploying a real application requires MANY YAML files:
- Deployment, Service, ConfigMap, Secret
- Ingress, PersistentVolumeClaim
- ServiceAccount, RBAC rules
- HorizontalPodAutoscaler
- NetworkPolicy

And you need different values for dev/staging/prod!

```
The Problem:
═══════════════════════════════════════════════════════════════

   Development:           Staging:            Production:
   deployment.yaml        deployment.yaml     deployment.yaml
   service.yaml           service.yaml        service.yaml
   configmap.yaml         configmap.yaml      configmap.yaml
   secret.yaml            secret.yaml         secret.yaml
   ... (10 files)         ... (10 files)      ... (10 files)

   They're 95% identical with small differences!
   Maintaining all this = Nightmare! 😱
```

### 36.2 Helm to the Rescue!

**Helm** is a package manager for Kubernetes:
- `apt` for Ubuntu
- `brew` for macOS
- **Helm** for Kubernetes

A Helm **Chart** is a package containing all Kubernetes resources for an application.

```
With Helm:
═══════════════════════════════════════════════════════════════

   ┌─────────────────────────────────────────────────────────┐
   │                    Helm Chart                            │
   │  (templates + default values)                            │
   │                                                          │
   │  templates/                                              │
   │  ├── deployment.yaml    ← Has {{ .Values.replicas }}    │
   │  ├── service.yaml       ← Has {{ .Values.port }}        │
   │  ├── configmap.yaml                                      │
   │  └── ingress.yaml                                        │
   │                                                          │
   │  values.yaml            ← Default values                 │
   └─────────────────────────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
   ┌───────────┐      ┌───────────┐      ┌───────────┐
   │   dev     │      │  staging  │      │   prod    │
   │  values   │      │  values   │      │  values   │
   │           │      │           │      │           │
   │replicas: 1│      │replicas: 2│      │replicas: 5│
   │port: 3000 │      │port: 3000 │      │port: 3000 │
   └───────────┘      └───────────┘      └───────────┘
```

---

## Chapter 37: Helm Basics

### 37.1 Installing Helm

```bash
# macOS
brew install helm

# Linux
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Verify
helm version
```

### 37.2 Helm Repositories

Charts are stored in repositories:

```bash
# Add popular repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus https://prometheus-community.github.io/helm-charts
helm repo update

# Search for charts
helm search repo nginx
helm search repo postgresql

# Show chart info
helm show chart bitnami/nginx
helm show values bitnami/nginx
```

### 37.3 Installing a Chart

```bash
# Basic install
helm install my-nginx bitnami/nginx

# Install in specific namespace
helm install my-nginx bitnami/nginx -n web --create-namespace

# Install with custom values
helm install my-nginx bitnami/nginx --set replicaCount=3

# Install with values file
helm install my-nginx bitnami/nginx -f my-values.yaml

# See what would be created (dry run)
helm install my-nginx bitnami/nginx --dry-run
```

### 37.4 Managing Releases

```bash
# List installed releases
helm list
helm list -A  # All namespaces

# Get release status
helm status my-nginx

# Upgrade a release
helm upgrade my-nginx bitnami/nginx --set replicaCount=5

# Rollback
helm rollback my-nginx 1  # Rollback to revision 1

# View history
helm history my-nginx

# Uninstall
helm uninstall my-nginx
```

---

## Chapter 38: Creating Your Own Chart

### 38.1 Chart Structure

```bash
# Create a new chart
helm create mychart

# What gets created:
mychart/
├── Chart.yaml          # Chart metadata
├── values.yaml         # Default values
├── templates/          # Template files
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── hpa.yaml
│   ├── serviceaccount.yaml
│   ├── _helpers.tpl    # Template helpers
│   ├── NOTES.txt       # Post-install message
│   └── tests/
│       └── test-connection.yaml
└── charts/             # Dependencies
```

### 38.2 Chart.yaml

```yaml
apiVersion: v2
name: mychart
description: My awesome application
type: application
version: 1.0.0         # Chart version
appVersion: "2.0.0"    # App version
```

### 38.3 values.yaml

```yaml
# Default values for mychart
replicaCount: 1

image:
  repository: nginx
  pullPolicy: IfNotPresent
  tag: "1.25"

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: nginx
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          pathType: Prefix

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

### 38.4 Template Example

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "mychart.fullname" . }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "mychart.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "mychart.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: {{ .Values.service.port }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

### 38.5 Template Functions

```yaml
# String functions
{{ .Values.name | upper }}           # MYAPP
{{ .Values.name | lower }}           # myapp
{{ .Values.name | title }}           # Myapp
{{ .Values.name | quote }}           # "myapp"

# Default values
{{ .Values.replicas | default 1 }}

# Conditionals
{{- if .Values.ingress.enabled }}
# Ingress resource here
{{- end }}

# Loops
{{- range .Values.hosts }}
  - host: {{ .host }}
{{- end }}

# Include templates
{{- include "mychart.labels" . | nindent 4 }}

# YAML conversion
{{- toYaml .Values.resources | nindent 12 }}
```

---

## Chapter 39: Testing and Debugging

```bash
# Lint your chart
helm lint mychart/

# Render templates locally
helm template my-release mychart/

# Render with custom values
helm template my-release mychart/ -f prod-values.yaml

# Debug with --debug
helm install my-release mychart/ --dry-run --debug

# Get rendered manifests of installed release
helm get manifest my-release
```

---

## Chapter 40: Example - Complete Chart

```yaml
# values.yaml for a web application
replicaCount: 2

image:
  repository: mycompany/webapp
  tag: "1.0.0"

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: webapp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: webapp-tls
      hosts:
        - webapp.example.com

config:
  databaseHost: postgres.default.svc
  logLevel: info

secrets:
  databasePassword: ""  # Override at deploy time!

resources:
  limits:
    cpu: 500m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

Deploy with environment-specific values:

```bash
# Development
helm install webapp ./mychart -f values-dev.yaml

# Staging
helm install webapp ./mychart -f values-staging.yaml

# Production
helm install webapp ./mychart \
  -f values-prod.yaml \
  --set secrets.databasePassword=$DB_PASSWORD
```

---

## Summary: What We Learned in Part 7

1. **Helm**: Package manager for Kubernetes
2. **Charts**: Packages containing templates + values
3. **Releases**: Installed instances of charts
4. **Values**: Configuration that varies per environment
5. **Templates**: Go templates for Kubernetes resources

**Key Commands**:
- `helm repo add/update`
- `helm install/upgrade/rollback/uninstall`
- `helm template` (render locally)
- `helm lint` (validate chart)

---

**Coming up in Part 8**: Capstone Project - Production Metrics Platform!
