# 🚀 Kubernetes Crash Course: From Zero to Production

A comprehensive, hands-on Kubernetes course using the "Explain Like I'm 5" teaching methodology. Learn Kubernetes from absolute basics to deploying production-ready systems.

## 🎯 What You'll Learn

By the end of this course, you'll be able to:

- ✅ Understand Kubernetes architecture and core concepts
- ✅ Deploy and manage containerized applications
- ✅ Handle configuration with ConfigMaps and Secrets
- ✅ Manage persistent storage for stateful applications
- ✅ Expose applications with Services and Ingress
- ✅ Package applications with Helm charts
- ✅ Build production-ready systems with best practices

## 📚 Course Structure

| Part | Topic | What You'll Build |
|------|-------|-------------------|
| 1 | [Introduction](docs/01-INTRODUCTION.md) | First Pod, understanding YAML |
| 2 | [Deployments](docs/02-DEPLOYMENTS.md) | Reliable, self-healing applications |
| 3 | [Services & Networking](docs/03-SERVICES.md) | Pod-to-Pod communication |
| 4 | [ConfigMaps & Secrets](docs/04-CONFIGMAPS-SECRETS.md) | External configuration |
| 5 | [Storage](docs/05-STORAGE.md) | Persistent data with PVs/PVCs |
| 6 | [Ingress](docs/06-INGRESS.md) | HTTP routing and TLS |
| 7 | [Helm](docs/07-HELM.md) | Package management |
| 8 | **Capstone Project** | Production metrics platform |

## 🚀 Prerequisites

- Basic command line knowledge
- Familiarity with containers (Docker)
- A working Kubernetes cluster

### Setting Up Your Environment

```bash
# Option 1: Minikube (recommended for learning)
brew install minikube  # macOS
minikube start

# Option 2: kind (Kubernetes in Docker)
brew install kind
kind create cluster

# Option 3: Docker Desktop
# Enable Kubernetes in Settings → Kubernetes

# Verify your cluster
kubectl cluster-info
kubectl get nodes
```

## 📖 How to Use This Course

1. **Read the documentation** in `docs/` in order (01 through 07)
2. **Follow along** with the examples - type them out, don't just copy!
3. **Complete the exercises** at the end of each section
4. **Build the capstone project** to solidify your knowledge

## 🏗️ Capstone Project

Deploy a complete **metrics collection platform** with:

- 📊 PostgreSQL StatefulSet with persistent storage
- 🚀 API Server Deployment with HPA autoscaling
- 📈 Metrics Agent DaemonSet (runs on every node)
- 📉 Grafana Deployment for visualization
- 🌐 Ingress for external access
- 🔒 Network Policies for security
- ⚙️ ConfigMaps/Secrets for configuration

See `projects/capstone-metrics-platform/` for the complete manifests.

## 🎓 Teaching Philosophy

Every concept is explained with:
- **Simple analogies** ("Pods are like shipping containers")
- **ASCII diagrams** showing architecture
- **Heavily commented YAML** explaining every field
- **Common mistakes** and how to avoid them
- **Debugging commands** for troubleshooting

## 📁 Directory Structure

```
k8s-course/
├── README.md                 # This file
├── docs/                     # Course documentation
│   ├── 01-INTRODUCTION.md
│   ├── 02-DEPLOYMENTS.md
│   ├── 03-SERVICES.md
│   ├── 04-CONFIGMAPS-SECRETS.md
│   ├── 05-STORAGE.md
│   ├── 06-INGRESS.md
│   └── 07-HELM.md
└── projects/                 # Hands-on projects
    ├── 00-hello-pod/
    ├── 01-deployments/
    ├── 02-services/
    ├── 03-configmaps-secrets/
    ├── 04-storage/
    ├── 05-ingress/
    ├── 06-helm/
    └── capstone-metrics-platform/
        ├── manifests/
        └── helm/
```

## Quick Reference

### Essential kubectl Commands

```bash
# Viewing Resources
kubectl get pods/deployments/services/ingress
kubectl describe pod <name>
kubectl logs <pod-name>

# Creating/Updating
kubectl apply -f manifest.yaml
kubectl create -f manifest.yaml

# Scaling
kubectl scale deployment <name> --replicas=5

# Rollouts
kubectl rollout status deployment/<name>
kubectl rollout undo deployment/<name>

# Debugging
kubectl exec -it <pod> -- /bin/bash
kubectl port-forward <pod> 8080:80
```

## 🤝 Contributing

Found an error or have a suggestion? Open an issue or submit a PR!

## 📜 License

MIT License - Feel free to use this for learning and teaching.
