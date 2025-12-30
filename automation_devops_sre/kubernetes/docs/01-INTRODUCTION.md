# Kubernetes Crash Course: From Zero to Production

## Part 1: What Even IS Kubernetes? (Explain Like I'm 5)

---

## Chapter 1: The Problem Kubernetes Solves

### 1.1 The Old Days: Pets vs. Cattle

Imagine you're running a pizza restaurant. In the old days, you had ONE amazing pizza oven. You named it "Big Bertha." You knew all its quirks - the left side runs hot, you have to jiggle the door handle just right. If Bertha broke, DISASTER! You'd scramble to fix her because she was irreplaceable.

This is how we used to treat servers - like **pets**. Each server had a name, a personality, and was carefully maintained.

```
The Pet Model (Old Way):

    ┌────────────────────────────────┐
    │     🖥️  "WebServer-Prod-01"    │
    │                                │
    │  - Named and special           │
    │  - Manually configured         │
    │  - Fixed when sick             │
    │  - Irreplaceable!              │
    │                                │
    └────────────────────────────────┘
    
    If it dies: PANIC! 😱 Call everyone! Emergency!
```

But what if instead you had 100 identical pizza ovens? If one breaks, who cares? Throw it away and use another one. No names, no special treatment. These are **cattle** - identical, replaceable units.

```
The Cattle Model (Cloud Native Way):

    ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐
    │ 🖥️  │ │ 🖥️  │ │ 🖥️  │ │ 🖥️  │ │ 🖥️  │
    │ #1  │ │ #2  │ │ #3  │ │ #4  │ │ #5  │
    └─────┘ └─────┘ └─────┘ └─────┘ └─────┘
    
    If one dies: ¯\_(ツ)_/¯ "Whatever, spin up a new one!"
```

**Kubernetes helps you manage cattle, not pets!**

### 1.2 Containers: The Building Blocks

Before Kubernetes, we need to understand containers.

A **container** is like a lunchbox for your application:
- Contains everything your app needs to run (code, libraries, settings)
- Isolated from other containers (your lunchbox doesn't mix with others)
- Portable (same lunchbox works at school, home, or the park)

```
Traditional Deployment:
┌──────────────────────────────────────────────────────────┐
│                      Server                               │
│  ┌──────────────────────────────────────────────────────┐│
│  │                Operating System                       ││
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐              ││
│  │  │  App A  │  │  App B  │  │  App C  │              ││
│  │  │Libraries│  │Libraries│  │Libraries│  ← Conflicts!││
│  │  └─────────┘  └─────────┘  └─────────┘              ││
│  └──────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────┘

Problem: App A needs Python 2, App B needs Python 3... CONFLICT!

Container Deployment:
┌──────────────────────────────────────────────────────────┐
│                      Server                               │
│  ┌──────────────────────────────────────────────────────┐│
│  │                Operating System                       ││
│  │              Container Runtime (Docker)               ││
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐              ││
│  │  │Container│  │Container│  │Container│              ││
│  │  │  App A  │  │  App B  │  │  App C  │  ← Isolated! ││
│  │  │Python 2 │  │Python 3 │  │ Node.js │              ││
│  │  └─────────┘  └─────────┘  └─────────┘              ││
│  └──────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────┘

Solution: Each container brings its own dependencies!
```

### 1.3 Kubernetes: The Container Orchestra

So we have containers. Great! But managing ONE container is easy. What about 100? 1,000? 10,000?

You need something to:
- Start containers on available machines
- Restart them if they crash
- Handle networking between them
- Scale up when traffic increases
- Scale down when it's quiet
- Update to new versions without downtime

This is what **Kubernetes** does! It's like a conductor for an orchestra of containers.

```
Without Kubernetes:
┌───────────────────────────────────────────────────────────────┐
│                                                               │
│   "Okay, let me SSH into server 1 and start the web app..."  │
│   "Now server 2 for the database..."                          │
│   "Wait, server 3 crashed! Let me manually restart..."       │
│   "Traffic spike! Gotta manually add more servers..."        │
│                                                               │
│   😓 (You, at 3 AM, for the 5th time this week)              │
│                                                               │
└───────────────────────────────────────────────────────────────┘

With Kubernetes:
┌───────────────────────────────────────────────────────────────┐
│                                                               │
│   You: "I want 5 copies of the web app running."             │
│   Kubernetes: "Done. I'll keep it that way forever."         │
│                                                               │
│   (Server crashes)                                            │
│   Kubernetes: "I noticed. Already fixed. You're welcome."    │
│                                                               │
│   (Traffic spike)                                             │
│   Kubernetes: "Scaling up automatically. No action needed."  │
│                                                               │
│   😴 (You, sleeping peacefully)                              │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

### 1.4 The Name "Kubernetes"

Kubernetes (κυβερνήτης) is Greek for "helmsman" or "pilot" - the person who steers the ship. The logo is a ship's wheel!

Google created it, based on 15+ years of running containers internally (they called their system "Borg" - like Star Trek!). In 2014, they open-sourced it as Kubernetes.

People often call it "K8s" (pronounced "kates") because there are 8 letters between K and S:

```
K-u-b-e-r-n-e-t-e-s
K-[  8 letters  ]-s
      = K8s!
```

---

## Chapter 2: Kubernetes Core Concepts (The Big Picture)

### 2.1 The Cluster: Your Container Kingdom

A Kubernetes "cluster" is a group of machines working together. Think of it like a beehive:

```
A Kubernetes Cluster
═══════════════════════════════════════════════════════════════

     ┌──────────────────────────────────────────────────────┐
     │                   CONTROL PLANE                       │
     │              (The Queen Bee / The Brain)              │
     │                                                       │
     │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
     │  │   API   │ │  etcd   │ │Scheduler│ │Controller│    │
     │  │ Server  │ │(memory) │ │(assigns)│ │ Manager  │    │
     │  └─────────┘ └─────────┘ └─────────┘ └─────────┘    │
     └──────────────────────────────────────────────────────┘
                              │
                              │ (gives orders)
                              ▼
     ┌──────────────────────────────────────────────────────┐
     │                     WORKER NODES                      │
     │              (The Worker Bees / The Muscle)           │
     │                                                       │
     │  ┌────────────┐  ┌────────────┐  ┌────────────┐      │
     │  │  Node #1   │  │  Node #2   │  │  Node #3   │      │
     │  │ [Pod][Pod] │  │ [Pod][Pod] │  │ [Pod][Pod] │      │
     │  │ [Pod][Pod] │  │ [Pod]      │  │ [Pod][Pod] │      │
     │  └────────────┘  └────────────┘  └────────────┘      │
     └──────────────────────────────────────────────────────┘
```

**Control Plane** = The brains. Makes all the decisions.
**Worker Nodes** = The muscle. Actually runs your containers.

### 2.2 Pods: The Smallest Unit

A "Pod" is the smallest thing Kubernetes manages. It's a wrapper around one or more containers.

**Wait, why not just run containers directly?**

Great question! Sometimes containers need to work SUPER closely together. Like, they need to share files or talk over localhost. A Pod groups those tightly-coupled containers together.

```
A Pod:
┌─────────────────────────────────────────────────────────────┐
│  POD                                                         │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   Container 1   │  │   Container 2   │                   │
│  │   (Web App)     │  │   (Log Shipper) │                   │
│  └─────────────────┘  └─────────────────┘                   │
│                                                              │
│  - Shared network (localhost works between them)            │
│  - Shared storage (same files accessible)                   │
│  - Same IP address                                          │
│  - Scheduled together on the same Node                      │
└─────────────────────────────────────────────────────────────┘
```

**In practice**: Most Pods have just ONE container. Multi-container Pods are for special patterns like sidecars (we'll learn about those later).

### 2.3 Nodes: The Workers

A "Node" is a machine (physical or virtual) that runs Pods. It's a worker bee.

Each Node runs:
- **kubelet**: The agent that talks to the Control Plane
- **kube-proxy**: Handles networking
- **Container runtime**: Usually Docker or containerd

```
A Node:
┌─────────────────────────────────────────────────────────────┐
│  NODE (a server)                                             │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  kubelet          "I report to the Control Plane"     │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  kube-proxy       "I handle networking magic"         │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Container Runtime (Docker/containerd)                │   │
│  │  ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐            │   │
│  │  │ Pod 1 │ │ Pod 2 │ │ Pod 3 │ │ Pod 4 │            │   │
│  │  └───────┘ └───────┘ └───────┘ └───────┘            │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.4 The Control Plane: The Brain

The Control Plane makes all decisions. It consists of:

**API Server** - The "front desk" 
- Everything talks to this (kubectl, nodes, etc.)
- The ONLY way to interact with the cluster

**etcd** - The "memory"
- Stores all cluster state
- "There are 3 Pods of nginx running on Node 2"

**Scheduler** - The "matchmaker"
- Decides which Node should run a new Pod
- "This Pod needs 2GB RAM, Node 3 has space, put it there!"

**Controller Manager** - The "supervisor"
- Runs controllers that maintain desired state
- "We should have 3 Pods but only 2 are running... create another!"

```
Control Plane Components:

                      ┌───────────────┐
                      │   kubectl     │ (your commands)
                      │   (CLI)       │
                      └───────┬───────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│  ┌─────────────────┐                                            │
│  │   API SERVER    │◄──── Everything goes through here          │
│  │  (Front Desk)   │                                            │
│  └────────┬────────┘                                            │
│           │                                                      │
│     ┌─────┴─────┬─────────────┐                                 │
│     ▼           ▼             ▼                                 │
│ ┌────────┐ ┌──────────┐ ┌───────────────┐                      │
│ │  etcd  │ │Scheduler │ │  Controller   │                      │
│ │(Memory)│ │(Assigner)│ │   Manager     │                      │
│ └────────┘ └──────────┘ └───────────────┘                      │
│                                                                  │
│  "Where    "Which node  "Are we at the                          │
│   pods      should run   desired state?                         │
│   are?"     this pod?"   If not, fix it!"                       │
└─────────────────────────────────────────────────────────────────┘
```

---

## Chapter 3: Your First Steps with kubectl

### 3.1 Setting Up Your Environment

For learning, you can use:

**Option 1: Minikube** (Single-node cluster on your laptop)
```bash
# Install minikube
brew install minikube   # macOS
# or: Download from https://minikube.sigs.k8s.io/

# Start a cluster
minikube start

# Check it's running
kubectl cluster-info
```

**Option 2: kind** (Kubernetes IN Docker)
```bash
# Install kind
brew install kind   # macOS

# Create a cluster
kind create cluster

# Check it's running
kubectl get nodes
```

**Option 3: Docker Desktop**
- Settings → Kubernetes → Enable Kubernetes → Apply & Restart

### 3.2 Verifying Your Setup

Once you have a cluster running:

```bash
# Check your cluster is running
kubectl cluster-info

# Example output:
# Kubernetes control plane is running at https://127.0.0.1:49157
# CoreDNS is running at https://127.0.0.1:49157/api/v1/...

# See your nodes
kubectl get nodes

# Example output:
# NAME       STATUS   ROLES           AGE   VERSION
# minikube   Ready    control-plane   10m   v1.28.0

# See what's running in kube-system (built-in stuff)
kubectl get pods -n kube-system
```

### 3.3 Your First Pod (Hello World!)

Create a file called `my-first-pod.yaml`:

```yaml
# my-first-pod.yaml
#
# This is a Kubernetes "manifest" - a file that describes what you want.
# Kubernetes reads this and makes it happen!

apiVersion: v1          # Which version of the Kubernetes API to use
kind: Pod               # What are we creating? A Pod!

metadata:               # Information ABOUT this Pod
  name: my-first-pod    # The name (must be unique in the namespace)
  labels:               # Labels are key-value tags for organizing
    app: nginx          # "This Pod is part of the 'nginx' app"

spec:                   # The SPECIFICATION - what should exist
  containers:           # List of containers in this Pod
  - name: nginx         # Name of this container
    image: nginx:1.25   # Docker image to use
    ports:
    - containerPort: 80 # This container listens on port 80
```

Apply it:

```bash
# Create the Pod from the file
kubectl apply -f my-first-pod.yaml

# Check it's running
kubectl get pods

# Output:
# NAME           READY   STATUS    RESTARTS   AGE
# my-first-pod   1/1     Running   0          5s

# Get more details
kubectl describe pod my-first-pod

# See the logs
kubectl logs my-first-pod

# Delete it when done
kubectl delete pod my-first-pod
```

---

## Chapter 4: The Declarative Philosophy

### 4.1 Imperative vs. Declarative

**Imperative** (Do this, then that):
```bash
kubectl run my-pod --image=nginx
kubectl scale deployment my-app --replicas=3
kubectl set image deployment/my-app nginx=nginx:1.26
```

**Declarative** (Make it look like this):
```bash
kubectl apply -f my-deployment.yaml
# The YAML file describes the DESIRED state
# Kubernetes figures out how to get there
```

**Kubernetes prefers declarative!** Why?

```
Imperative:                     Declarative:
                               
You: "Create a pod"             You: "Here's what I want"
You: "Now scale to 3"           K8s: "Got it. I'll make it happen
You: "Update the image"              and KEEP it that way."
You: "Oh no it crashed! Fix it!"     
                               ... (Kubernetes manages the steps)
```

### 4.2 The Reconciliation Loop

Kubernetes constantly runs a "reconciliation loop":

1. Look at the DESIRED state (what you want, from your YAML files)
2. Look at the ACTUAL state (what's really happening)
3. Take ACTION to make actual match desired
4. Repeat forever!

```
         ┌──────────────────────────────────────┐
         │   DESIRED STATE                      │
         │   "I want 3 nginx Pods"              │
         └─────────────┬────────────────────────┘
                       │ compare
                       ▼
         ┌──────────────────────────────────────┐
         │   ACTUAL STATE                       │
         │   "There are 2 nginx Pods"           │
         └─────────────┬────────────────────────┘
                       │ difference!
                       ▼
         ┌──────────────────────────────────────┐
         │   ACTION                             │
         │   "Start 1 more nginx Pod"           │
         └─────────────┬────────────────────────┘
                       │
                       └───────► (repeat forever)
```

---

## Summary: What We Learned in Part 1

1. **The Problem**: Managing containers at scale is hard. Kubernetes automates it.

2. **Core Concepts**:
   - Cluster = Group of machines (Control Plane + Worker Nodes)
   - Pod = Smallest deployable unit (wraps container(s))
   - Node = A machine that runs Pods
   - Control Plane = The brain (API Server, etcd, Scheduler, Controllers)
   - kubectl = Your command-line tool to talk to Kubernetes

3. **Declarative Philosophy**: Tell Kubernetes WHAT you want, not HOW to do it.

4. **YAML Structure**: apiVersion, kind, metadata, spec

5. **Key Commands**:
   - `kubectl get pods` - List Pods
   - `kubectl describe pod X` - Details about Pod X
   - `kubectl logs X` - Logs from Pod X
   - `kubectl apply -f file.yaml` - Apply a manifest
   - `kubectl delete -f file.yaml` - Delete resources

---

## Exercises

1. Create a Pod running the `redis` image
2. Use `kubectl describe` to find what IP address your Pod got
3. Use `kubectl exec -it <pod-name> -- /bin/bash` to shell into a Pod
4. Delete your Pod and watch it NOT come back (we'll fix this in Part 2!)

---

**Coming up in Part 2**: Deployments, ReplicaSets, and making your Pods reliable!
