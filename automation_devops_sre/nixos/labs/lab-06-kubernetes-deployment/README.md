# Lab 6: Deploying Kubernetes Across Clos Network

## What You'll Learn

- Deploy Kubernetes control plane on core/dist layer
- Deploy worker nodes across pods
- Use network fabric for pod-to-pod communication  
- Scale applications across datacenters
- Automatic service discovery via EVPN

## Estimated Time: 90 minutes

## Prerequisites

- Completed Lab 5 (Network fabric running)
- BGP/VXLAN connectivity verified
- All 12 LXD containers networked

---

## Part 1: Kubernetes on NixOS

### Why NixOS + Kubernetes?

```
Traditional:
  Install Kubernetes manually
  Install networking plugin
  Configure routes manually
  Fragile, hard to replicate
  
NixOS + Kubernetes:
  Declare K8s in configuration.nix
  Network already working
  Automatic pod networking
  Reproducible, version-controlled
```

---

## Part 2: Deploy Control Plane

### On Core-1: Kubernetes Master

Edit `/etc/nixos/configuration.nix`:

```nix
{ config, pkgs, ... }:

{
  # Kubernetes master node
  services.kubernetes = {
    roles = ["master"];
    masterAddress = "10.100.1.1";
    clusterCidr = "192.168.0.0/16";  # Pod CIDR
    serviceClusterIpRange = "10.32.0.0/24";
    
    # Networking plugin - Cilium for VXLAN support
    networking.plugin = "cni";
  };
  
  # Cilium CNI for Kubernetes
  services.cilium.enable = true;
  services.cilium.debug = true;
  
  networking.hostname = "core-1-k8s";
  system.stateVersion = "23.11";
}
```

Rebuild:
```bash
sudo nixos-rebuild switch

# Wait for control plane
kubectl get nodes

# Should eventually show:
# NAME      STATUS   ROLES    AGE
# core-1    Ready    master   2m
```

**What this does (ELI5):**
```
"I am a Kubernetes master node"
"My address is 10.100.1.1"
"Pods should use 192.168.0.0/16 addresses"
"Use Cilium to connect pods"
"Cilium knows about VXLAN so pods can communicate"
```

---

## Part 3: Deploy Worker Nodes on Pods

### On Pod-1: Worker Node Configuration

Edit `/etc/nixos/configuration.nix`:

```nix
{ config, pkgs, ... }:

{
  # Kubernetes worker node
  services.kubernetes = {
    roles = ["node"];
    masterAddress = "10.100.1.1";  # Connect to master on core-1
    clusterCidr = "192.168.0.0/16";
  };
  
  # Cilium CNI
  services.cilium.enable = true;
  
  # Container runtime
  virtualisation.docker.enable = true;
  
  networking.hostname = "pod-1-worker";
  system.stateVersion = "23.11";
}
```

Rebuild:
```bash
sudo nixos-rebuild switch

# From core-1, check nodes joined
kubectl get nodes

# Should show:
# NAME         STATUS   ROLES    
# core-1       Ready    master   
# pod-1        Ready    worker   
# pod-2        Ready    worker   
# ... (add all 6 pods as workers)
```

**Repeat for Pod-2 through Pod-6** with appropriate hostnames.

---

## Part 4: Deploy Test Application

### Deploy Nginx Across Pods

Create `/tmp/nginx-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-multi-pod
  namespace: default
spec:
  replicas: 6  # One per pod
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      # Spread across all pods
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - nginx
            topologyKey: kubernetes.io/hostname
      
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
```

Deploy:
```bash
kubectl apply -f /tmp/nginx-deployment.yaml

# Watch deployment
kubectl get pods -w

# Should show:
# NAME                             READY   STATUS    RESTARTS
# nginx-multi-pod-xxx-yyy          1/1     Running   0
# nginx-multi-pod-aaa-bbb          1/1     Running   0
# ... (6 pods total, one per node)
```

### Create Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
  type: ClusterIP
```

Deploy:
```bash
kubectl apply -f /tmp/nginx-service.yaml

# Get service
kubectl get svc nginx-service

# Test from any pod
kubectl exec -it nginx-multi-pod-xxx-yyy -- curl nginx-service

# Should show nginx default page
```

---

## Part 5: Cross-Datacenter Deployment

### Simulate Multi-Datacenter

**Datacenter A: Pods 1-3**
```bash
# Label nodes
kubectl label nodes pod-1 datacenter=dc-a
kubectl label nodes pod-2 datacenter=dc-a
kubectl label nodes pod-3 datacenter=dc-a
```

**Datacenter B: Pods 4-6**
```bash
# Label nodes
kubectl label nodes pod-4 datacenter=dc-b
kubectl label nodes pod-5 datacenter=dc-b
kubectl label nodes pod-6 datacenter=dc-b
```

### Deploy Frontend (DC-A) and Backend (DC-B)

Frontend Deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      # Only on DC-A
      nodeSelector:
        datacenter: dc-a
      
      containers:
      - name: frontend
        image: nginx:latest
        env:
        - name: BACKEND_URL
          value: "http://backend-service:8080"
```

Backend Deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      # Only on DC-B
      nodeSelector:
        datacenter: dc-b
      
      containers:
      - name: backend
        image: python:3.11
        command:
        - python
        - -m
        - http.server
        - "8080"
```

Deploy both:
```bash
kubectl apply -f frontend.yaml
kubectl apply -f backend.yaml

# Check pods deployed correctly
kubectl get pods -L datacenter

# Frontend should be on pod-1, pod-2, pod-3
# Backend should be on pod-4, pod-5, pod-6
```

---

## Part 6: Network Magic Happens

### What's Happening Behind the Scenes

```
Frontend Pod (Pod-1, DC-A):
  IP: 192.168.10.50
  
Backend Pod (Pod-4, DC-B):
  IP: 192.168.20.40
  
Frontend sends request to Backend:
  
  Step 1: Frontend resolves backend-service
    DNS returns: 10.32.0.100 (Kubernetes service IP)
    
  Step 2: Frontend sends packet to 10.32.0.100
    
  Step 3: Cilium intercepts (network magic!)
    "Oh, 10.32.0.100 is backend service"
    "Backend pods are at 192.168.20.40"
    "But they're in DC-B, use VXLAN tunnel"
    
  Step 4: VXLAN wraps packet
    Outer: 10.100.10.1 (Dist-1) → 10.100.12.1 (Dist-3)
    Inner: 192.168.10.50 → 192.168.20.40
    
  Step 5: Packet travels through network fabric
    Pod-1 → Dist-1 → Core-1/2 → Dist-3 → Pod-4
    (All automatic via BGP routing!)
    
  Step 6: VXLAN unwraps on Dist-3
    
  Step 7: Packet reaches Pod-4
    Backend pod receives request
    Sends response back (same tunnel path)
    
Step 8: Frontend receives response
```

**The amazing part: This all works AUTOMATICALLY!**
- No manual routing configuration
- BGP discovers paths
- Cilium learns pod locations
- VXLAN creates virtual network
- Frontend and Backend think they're local!

---

## Part 7: Verify Cross-Datacenter Communication

### Test Frontend → Backend

```bash
# Get frontend pod
kubectl get pods -l app=frontend -o name | head -1
POD=pod/frontend-xxx-yyy

# Exec into frontend
kubectl exec -it $POD -- sh

# From inside frontend pod
curl backend-service:8080

# Should get response from backend in DC-B!
```

### Monitor Network Traffic

From any distribution router:

```bash
lxc shell dist-1

# Watch VXLAN traffic
sudo tcpdump -i vxlan100

# Should show packets traveling through tunnel
```

### Check EVPN Routes

```bash
lxc shell dist-1
sudo vtysh

# Show EVPN routes
show bgp l2vpn evpn route

# Should show all pod locations
```

---

## Part 8: Scale Applications

### Scale Frontend to 6 Replicas

```bash
kubectl scale deployment frontend --replicas=6

# Wait for deployment
kubectl get pods

# Now frontend pods on all 6 nodes!
```

### Scale Backend

```bash
kubectl scale deployment backend --replicas=6

# Backend pods on all 6 nodes!
```

### Observe Network Balancing

```bash
# Each frontend pod can reach any backend pod
# Traffic automatically load-balanced
# Network ensures optimal paths

kubectl get pods -o wide

# Shows which pod on which node
# All communication works seamlessly
```

---

## Part 9: Simulate Failure

### Stop a Router

```bash
lxc stop dist-1

# Observe what happens
kubectl get pods

# Pods should still work!
# BGP finds alternate path
# No application downtime
```

### Restart Router

```bash
lxc start dist-1

# Automatically rejoins network
# Routes reestablish
# All pods reconnect
```

---

## Verification Checklist

- [ ] Kubernetes master on core-1
- [ ] All 6 pods joined as workers
- [ ] Nginx deployment running (6 replicas)
- [ ] Service accessible
- [ ] Frontend pods in DC-A
- [ ] Backend pods in DC-B
- [ ] Cross-DC communication working
- [ ] Frontend can reach Backend service
- [ ] EVPN showing pod locations
- [ ] Scaling works
- [ ] Failure recovery works

---

## What You've Built

✅ **Kubernetes Cluster**
- 1 master (core-1)
- 6 workers (pods 1-6)
- Spanning 2 virtual datacenters

✅ **Multi-Datacenter Application**
- Frontend tier in DC-A
- Backend tier in DC-B
- Automatic service discovery
- Transparent communication

✅ **Network Intelligence**
- BGP routing between datacenters
- VXLAN tunnels for pods
- EVPN service location discovery
- Automatic load balancing

✅ **Self-Healing**
- Survives router failures
- Automatic path recalculation
- No manual intervention needed

---

## Real-World Comparison

```
What You Built:                    Real Cloud:
─────────────────                  ──────────
Pod deployment across DCs    →     Multi-region deployment
BGP routing                  →     BGP in CDNs/cloud
VXLAN tunnels               →     Overlay networks (AWS VPC, GCP VPC)
EVPN discovery              →     Service meshes (Istio, Linkerd)
Automatic failover          →     Active-active datacenters
NixOS configs               →     IaC (Terraform, CloudFormation)
```

---

## Next: Lab 7

Explore advanced scenarios:
- Add observability (monitoring)
- Deploy with Helm charts
- Multi-region Istio service mesh
- Chaos engineering (intentional failures)

[Next: Lab 7 - Observability & Scaling](../lab-07-observability/README.md)

---

## Deep Dive Questions

1. **Why VXLAN instead of direct routing?**
   - Allows flexible subnet design
   - Decouples physical network from logical
   - Enables live migration

2. **How does service discovery work?**
   - Cilium watches Kubernetes API
   - Updates routes when pods start/stop
   - EVPN gossips locations

3. **Why BGP instead of OSPF?**
   - Hierarchical (scales to many routers)
   - EVPN extension for L2
   - What real datacenters use

4. **How does failover work?**
   - BGP detects router down
   - Announces alternate paths
   - All nodes update routes (LFIB)
   - <100ms convergence

---

## Troubleshooting

```
Problem: "Pods can't reach each other"
Debug: 
  kubectl exec -it pod -- bash
  ping other-pod-ip
  If fails: Check VXLAN tunnel active
  
Problem: "Services unreachable"
Debug:
  kubectl get svc
  kubectl get endpoints service-name
  curl service-ip:port from another pod
  
Problem: "High latency across DCs"
Debug:
  Check BGP path selection
  Check VXLAN tunnel MTU
  Monitor link bandwidth

Problem: "Pod stuck pending"
Debug:
  kubectl describe pod pod-name
  Check resource availability
  Check node capacity
```

---

## Reference

- [Kubernetes Networking](https://kubernetes.io/docs/concepts/cluster-administration/networking/)
- [Cilium Documentation](https://docs.cilium.io/)
- [VXLAN in Kubernetes](https://kubernetes.io/docs/tasks/administer-cluster/flannel/)
- [BGP in Kubernetes](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
