# Kubernetes Crash Course: From Zero to Production

## Part 3: Services & Networking (How Pods Talk to Each Other)

---

## Chapter 12: The Networking Problem

### 12.1 Pods Are Ephemeral

Pods can die and be replaced at any moment. Each new Pod gets a NEW IP address. How do other Pods find them?

```
The Problem:

Day 1:
  Frontend Pod needs Backend Pod
  Backend Pod IP: 10.0.1.5
  Frontend: "Calling 10.0.1.5... works!"

Day 2:
  Backend Pod crashed and was replaced!
  New Backend Pod IP: 10.0.2.17
  Frontend: "Calling 10.0.1.5... ERROR! Nothing there!"
  User: "The website is broken!" 😡
```

We need a **stable way** to find Pods. Enter **Services**!

---

## Chapter 13: Services - The Phone Book

A **Service** gives a stable IP and DNS name to a group of Pods.

Think of it like a phone book:
- You don't memorize phone numbers (Pod IPs)
- You look up "Pizza Place" (Service name)
- The phone book gives you the current number

```
Service as Phone Book:
═══════════════════════════════════════════════════════════════

    Frontend Pod                  "backend" Service
         │                              │
         │ "I need to call              │
         │  the backend!"               │
         │                              │
         └──────────────────────────────┤
                                        │ 10.96.0.50
                                        │ (stable IP!)
                                        │
                    ┌───────────────────┼───────────────────┐
                    │                   │                   │
                    ▼                   ▼                   ▼
               ┌─────────┐        ┌─────────┐        ┌─────────┐
               │Backend 1│        │Backend 2│        │Backend 3│
               │10.0.1.5 │        │10.0.2.17│        │10.0.3.42│
               └─────────┘        └─────────┘        └─────────┘

    The Service IP (10.96.0.50) NEVER changes!
    It load-balances to healthy Backend Pods.
```

### 13.1 Creating a Service

```yaml
# backend-service.yaml
#
# ELI5: "Create a stable phone number for my backend Pods"
#
apiVersion: v1
kind: Service
metadata:
  name: backend          # This becomes the DNS name!
spec:
  selector:              # "Which Pods should this Service route to?"
    app: backend         # Pods with label "app: backend"
  
  ports:
  - port: 80             # Port the Service listens on
    targetPort: 8080     # Port the Pods are listening on
    protocol: TCP
```

```bash
# Create the Service
kubectl apply -f backend-service.yaml

# See the Service
kubectl get services
# NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
# backend      ClusterIP   10.96.0.50     <none>        80/TCP    5s
# kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP   1d

# Now any Pod can reach backend at:
# - http://backend (within same namespace)
# - http://backend.default.svc.cluster.local (full DNS name)
# - http://10.96.0.50 (Cluster IP)
```

---

## Chapter 14: Service Types

### 14.1 ClusterIP (Default)

Only accessible from INSIDE the cluster.

```yaml
spec:
  type: ClusterIP      # Default, often omitted
  selector:
    app: backend
  ports:
  - port: 80
    targetPort: 8080
```

```
ClusterIP:
┌──────────────────────────────────────────────────────────────┐
│                       CLUSTER                                 │
│                                                               │
│   ┌─────────────┐         ┌─────────────┐                    │
│   │ Frontend    │────────▶│  Service    │                    │
│   │ Pod         │         │  10.96.0.50 │                    │
│   └─────────────┘         └──────┬──────┘                    │
│                                  │                            │
│                    ┌─────────────┼─────────────┐             │
│                    ▼             ▼             ▼             │
│               [Backend]     [Backend]     [Backend]          │
│                                                               │
└──────────────────────────────────────────────────────────────┘
         │
         │ ✗ Not accessible from outside!
         ▼
    [ OUTSIDE WORLD ]
```

### 14.2 NodePort

Exposes Service on each Node's IP at a static port.

```yaml
spec:
  type: NodePort
  selector:
    app: webapp
  ports:
  - port: 80           # Service port (inside cluster)
    targetPort: 8080   # Pod port
    nodePort: 30080    # External port (30000-32767)
```

```
NodePort:
┌──────────────────────────────────────────────────────────────┐
│                       CLUSTER                                 │
│                                                               │
│   Node 1 (192.168.1.10)    Node 2 (192.168.1.11)            │
│   Port: 30080              Port: 30080                       │
│         │                        │                           │
│         └────────────┬───────────┘                           │
│                      │                                        │
│                      ▼                                        │
│               ┌─────────────┐                                │
│               │  Service    │                                │
│               │  ClusterIP  │                                │
│               └──────┬──────┘                                │
│                      │                                        │
│         ┌────────────┼────────────┐                          │
│         ▼            ▼            ▼                          │
│    [Pod]        [Pod]        [Pod]                           │
└──────────────────────────────────────────────────────────────┘
         ▲
         │ ✓ Access via http://192.168.1.10:30080
         │ ✓ Or http://192.168.1.11:30080
         │
    [ OUTSIDE WORLD ]
```

### 14.3 LoadBalancer

Creates an external load balancer (in cloud environments).

```yaml
spec:
  type: LoadBalancer
  selector:
    app: webapp
  ports:
  - port: 80
    targetPort: 8080
```

```
LoadBalancer:
                    [ INTERNET ]
                         │
                         ▼
              ┌─────────────────────┐
              │  Cloud Load         │ ← Created by cloud provider
              │  Balancer           │   (AWS ELB, GCP LB, etc.)
              │  (public IP)        │
              └──────────┬──────────┘
                         │
┌────────────────────────┼─────────────────────────────────────┐
│                        │             CLUSTER                  │
│                        ▼                                      │
│               ┌─────────────┐                                │
│               │  Service    │                                │
│               │  (NodePort) │                                │
│               └──────┬──────┘                                │
│                      │                                        │
│         ┌────────────┼────────────┐                          │
│         ▼            ▼            ▼                          │
│    [Pod]        [Pod]        [Pod]                           │
└──────────────────────────────────────────────────────────────┘
```

---

## Chapter 15: DNS in Kubernetes

Kubernetes has built-in DNS! Every Service gets a DNS name.

```
DNS Name Formats:
═══════════════════════════════════════════════════════════════

Service in SAME namespace:
    backend
    
Service in DIFFERENT namespace:
    backend.other-namespace
    
Fully Qualified Domain Name (FQDN):
    backend.default.svc.cluster.local
    └──┬──┘ └──┬──┘ └─┬─┘ └────┬────┘
       │      │      │         │
       │      │      │         └── cluster domain
       │      │      └── "svc" = it's a Service
       │      └── namespace
       └── service name
```

```bash
# From inside a Pod, you can use DNS:

# Same namespace
curl http://backend/api

# Different namespace  
curl http://backend.production/api

# Full FQDN (works from anywhere)
curl http://backend.production.svc.cluster.local/api
```

---

## Chapter 16: Endpoints and Labels

### 16.1 How Services Find Pods

Services use **label selectors** to find Pods. The matching Pods become **Endpoints**.

```yaml
# Service selector
spec:
  selector:
    app: backend      # "Find Pods with label app=backend"
    tier: api         # AND tier=api

# Pod labels (must match!)
metadata:
  labels:
    app: backend
    tier: api
    version: v2       # Extra labels are okay
```

```bash
# See Endpoints
kubectl get endpoints backend
# NAME      ENDPOINTS                               AGE
# backend   10.0.1.5:8080,10.0.2.17:8080           5m

# Detailed view
kubectl describe endpoints backend
```

---

## Chapter 17: Network Policies (Firewall Rules)

By default, all Pods can talk to all other Pods. **Network Policies** restrict this.

```yaml
# network-policy.yaml
#
# ELI5: "Only allow frontend Pods to talk to backend Pods"
#
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-policy
spec:
  # Apply this policy to Pods with label app=backend
  podSelector:
    matchLabels:
      app: backend
  
  policyTypes:
  - Ingress      # Control incoming traffic
  
  ingress:
  # Allow traffic FROM Pods with label app=frontend
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
```

```
With Network Policy:

    ┌───────────────────────────────────────────────────────────┐
    │                                                           │
    │   [Frontend]                     [Backend]                │
    │   app=frontend ──────────────▶  app=backend               │
    │                     ✓ ALLOWED                             │
    │                                                           │
    │   [Attacker]                     [Backend]                │
    │   app=hacker ────────────────X  app=backend               │
    │                     ✗ BLOCKED                             │
    │                                                           │
    └───────────────────────────────────────────────────────────┘
```

---

## Chapter 18: Putting It Together

Here's a complete example with Deployment + Service:

```yaml
# webapp.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: webapp
  template:
    metadata:
      labels:
        app: webapp
    spec:
      containers:
      - name: webapp
        image: nginx:1.25
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: webapp
spec:
  type: ClusterIP
  selector:
    app: webapp
  ports:
  - port: 80
    targetPort: 80
```

```bash
# Deploy everything
kubectl apply -f webapp.yaml

# Check Deployment
kubectl get deployment webapp

# Check Service
kubectl get service webapp

# Check Endpoints
kubectl get endpoints webapp

# Test from another Pod
kubectl run test --rm -it --image=busybox -- wget -qO- http://webapp
```

---

## Summary: What We Learned in Part 3

1. **The Problem**: Pod IPs change, we need stable addressing

2. **Services** provide:
   - Stable IP address (ClusterIP)
   - DNS name (service-name.namespace)
   - Load balancing across Pods

3. **Service Types**:
   - ClusterIP: Internal only (default)
   - NodePort: Expose on each node's IP
   - LoadBalancer: Cloud load balancer

4. **DNS**: `service.namespace.svc.cluster.local`

5. **Network Policies**: Firewall rules for Pod-to-Pod traffic

6. **Key Commands**:
   - `kubectl get services`
   - `kubectl get endpoints`
   - `kubectl describe service X`

---

## Exercises

1. Create a Deployment with 2 replicas and a ClusterIP Service
2. Create another Pod and use `wget` to access your Service
3. Change the Service to NodePort and access it from your browser
4. Create a Network Policy that only allows specific Pods to access your Service

---

**Coming up in Part 4**: ConfigMaps and Secrets - External configuration!
