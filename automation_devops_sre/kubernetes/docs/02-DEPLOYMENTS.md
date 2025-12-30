# Kubernetes Crash Course: From Zero to Production

## Part 2: Deployments & ReplicaSets (Making Things Reliable)

---

## Chapter 5: The Problem with Raw Pods

### 5.1 Pods Are Fragile!

Remember our Pod from Part 1? Here's a harsh truth: **Pods are meant to die.**

If a Pod crashes, gets deleted, or the Node it's on fails... that Pod is GONE. Forever. Kubernetes doesn't automatically bring it back.

```
Raw Pod Life Cycle:

  You: kubectl run my-pod --image=nginx
                    │
                    ▼
            ┌──────────────┐
            │   Pod Born   │
            │  my-pod:nginx│
            └──────────────┘
                    │
          (something goes wrong)
                    │
                    ▼
            ┌──────────────┐
            │  Pod Died!   │
            │    ☠️ RIP    │
            └──────────────┘
                    │
                    ▼
            ┌──────────────┐
            │   NOTHING    │
            │  (it's gone) │
            └──────────────┘

  User: "Why is the website down?!"
```

We need Pods to:
1. Come back if they crash
2. Have multiple copies for redundancy
3. Update without downtime

**Enter ReplicaSets and Deployments!**

---

## Chapter 6: ReplicaSets - The Pod Babysitter

A **ReplicaSet** ensures a specified number of Pod copies are always running.

Think of it like a babysitter who counts kids:
- "I'm supposed to have 3 kids"
- "I only see 2 kids..."
- "Better get another one!"

```yaml
# replicaset-example.yaml
#
# ELI5: "Always keep 3 copies of this Pod running!"
#
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-rs
spec:
  replicas: 3                    # <-- "I want 3 Pods!"
  
  selector:                      # <-- "How do I find my Pods?"
    matchLabels:
      app: nginx                 # "Look for Pods with label app=nginx"
  
  template:                      # <-- "Here's how to make new Pods"
    metadata:
      labels:
        app: nginx               # Must match selector above!
    spec:
      containers:
      - name: nginx
        image: nginx:1.25
```

```bash
# Create the ReplicaSet
kubectl apply -f replicaset-example.yaml

# See 3 Pods created!
kubectl get pods
# NAME             READY   STATUS    RESTARTS   AGE
# nginx-rs-abc12   1/1     Running   0          5s
# nginx-rs-def34   1/1     Running   0          5s
# nginx-rs-ghi56   1/1     Running   0          5s

# Try deleting one...
kubectl delete pod nginx-rs-abc12

# Check again - it's back!
kubectl get pods
# NAME             READY   STATUS    RESTARTS   AGE
# nginx-rs-def34   1/1     Running   0          30s
# nginx-rs-ghi56   1/1     Running   0          30s
# nginx-rs-xyz99   1/1     Running   0          2s   # <-- NEW!
```

### 6.1 Why Not Use ReplicaSets Directly?

You CAN, but there's something even better: **Deployments!**

ReplicaSets can't:
- Handle rolling updates (updating without downtime)
- Roll back to a previous version
- Track revision history

---

## Chapter 7: Deployments - The Full Package

A **Deployment** manages ReplicaSets and provides:
- ✅ Self-healing (like ReplicaSet)
- ✅ Scaling (like ReplicaSet)
- ✅ Rolling updates (new!)
- ✅ Rollbacks (new!)
- ✅ Revision history (new!)

```
Hierarchy:

    ┌─────────────────────────────────────────────────────────┐
    │                      DEPLOYMENT                          │
    │  "I manage ReplicaSets and handle updates gracefully"   │
    └───────────────────────────┬─────────────────────────────┘
                                │ manages
                                ▼
    ┌─────────────────────────────────────────────────────────┐
    │                      REPLICASET                          │
    │        "I maintain 3 copies of the Pod template"        │
    └───────────────────────────┬─────────────────────────────┘
                                │ creates
            ┌───────────────────┼───────────────────┐
            ▼                   ▼                   ▼
        ┌───────┐          ┌───────┐          ┌───────┐
        │  Pod  │          │  Pod  │          │  Pod  │
        │  #1   │          │  #2   │          │  #3   │
        └───────┘          └───────┘          └───────┘
```

### 7.1 Creating a Deployment

```yaml
# nginx-deployment.yaml
#
# ELI5: "Run 3 nginx web servers, keep them running,
#        and help me update them safely!"
#
apiVersion: apps/v1
kind: Deployment           # <-- Deployment, not ReplicaSet!
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3              # How many Pods we want
  
  selector:                # How to find our Pods
    matchLabels:
      app: nginx
  
  template:                # Pod template (what each Pod looks like)
    metadata:
      labels:
        app: nginx         # Must match selector!
    spec:
      containers:
      - name: nginx
        image: nginx:1.25
        ports:
        - containerPort: 80
        
        # Resource limits (very important for production!)
        resources:
          requests:        # Minimum resources needed
            memory: "64Mi"
            cpu: "100m"    # 100 millicores = 0.1 CPU
          limits:          # Maximum resources allowed
            memory: "128Mi"
            cpu: "200m"
```

```bash
# Create the Deployment
kubectl apply -f nginx-deployment.yaml

# See the Deployment
kubectl get deployments
# NAME               READY   UP-TO-DATE   AVAILABLE   AGE
# nginx-deployment   3/3     3            3           10s

# See the ReplicaSet it created
kubectl get replicasets
# NAME                          DESIRED   CURRENT   READY   AGE
# nginx-deployment-7fb96c846b   3         3         3       10s

# See the Pods
kubectl get pods
# NAME                                READY   STATUS    RESTARTS   AGE
# nginx-deployment-7fb96c846b-abc12   1/1     Running   0          10s
# nginx-deployment-7fb96c846b-def34   1/1     Running   0          10s
# nginx-deployment-7fb96c846b-ghi56   1/1     Running   0          10s
```

---

## Chapter 8: Rolling Updates (Zero Downtime!)

### 8.1 The Magic of Rolling Updates

When you update a Deployment, Kubernetes replaces Pods gradually:

```
Rolling Update Process:
═══════════════════════════════════════════════════════════════

Step 1: Current state (v1.25)
    [Pod v1.25] [Pod v1.25] [Pod v1.25]
        ✓           ✓           ✓

Step 2: Start new Pod (v1.26), old still running
    [Pod v1.25] [Pod v1.25] [Pod v1.25] [Pod v1.26]
        ✓           ✓           ✓          (starting)

Step 3: New Pod ready, kill one old
    [Pod v1.25] [Pod v1.25] [Pod v1.26]
        ✓           ✓           ✓

Step 4: Repeat...
    [Pod v1.25] [Pod v1.26] [Pod v1.26]
        ✓           ✓           ✓

Step 5: Done!
    [Pod v1.26] [Pod v1.26] [Pod v1.26]
        ✓           ✓           ✓

Users never noticed! Zero downtime! 🎉
```

### 8.2 Triggering an Update

```bash
# Method 1: Edit the YAML and apply
# Change: image: nginx:1.25 → image: nginx:1.26
kubectl apply -f nginx-deployment.yaml

# Method 2: Direct command
kubectl set image deployment/nginx-deployment nginx=nginx:1.26

# Watch the rollout happen
kubectl rollout status deployment/nginx-deployment
# Waiting for deployment "nginx-deployment" rollout to finish:
# 1 out of 3 new replicas have been updated...
# 2 out of 3 new replicas have been updated...
# 3 out of 3 new replicas have been updated...
# deployment "nginx-deployment" successfully rolled out

# See the history
kubectl rollout history deployment/nginx-deployment
# REVISION  CHANGE-CAUSE
# 1         <none>
# 2         <none>
```

### 8.3 Rollbacks (Undo Button!)

Deployed a bad version? No problem!

```bash
# Oh no, v1.26 has a bug! Go back!
kubectl rollout undo deployment/nginx-deployment

# Or go to a specific revision
kubectl rollout undo deployment/nginx-deployment --to-revision=1

# Check it worked
kubectl rollout status deployment/nginx-deployment
```

### 8.4 Update Strategy Options

```yaml
spec:
  strategy:
    type: RollingUpdate      # or "Recreate"
    rollingUpdate:
      maxUnavailable: 25%    # Max Pods that can be down during update
      maxSurge: 25%          # Max extra Pods during update
```

**RollingUpdate** (default): Gradual replacement, zero downtime
**Recreate**: Kill ALL old Pods, then create ALL new Pods (causes downtime!)

---

## Chapter 9: Health Checks (Probes)

### 9.1 Why Health Checks Matter

Kubernetes needs to know:
1. **Is this Pod alive?** (livenessProbe)
2. **Is this Pod ready to receive traffic?** (readinessProbe)
3. **Has this Pod finished starting?** (startupProbe)

```yaml
# deployment-with-probes.yaml
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
        image: myapp:1.0
        ports:
        - containerPort: 8080
        
        # LIVENESS PROBE
        # "Is the application alive?"
        # If this fails, Kubernetes RESTARTS the container
        livenessProbe:
          httpGet:
            path: /healthz       # Hit this endpoint
            port: 8080
          initialDelaySeconds: 15  # Wait 15s before first check
          periodSeconds: 10        # Check every 10s
          timeoutSeconds: 5        # Timeout after 5s
          failureThreshold: 3      # Fail after 3 bad responses
        
        # READINESS PROBE
        # "Is the application ready for traffic?"
        # If this fails, Pod is removed from Service (no traffic)
        # Container is NOT restarted
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          failureThreshold: 3
        
        # STARTUP PROBE (for slow-starting apps)
        # "Has the application finished starting?"
        # Liveness/readiness don't run until this passes
        startupProbe:
          httpGet:
            path: /healthz
            port: 8080
          failureThreshold: 30      # Allow 30 attempts
          periodSeconds: 10         # = 5 minutes to start
```

### 9.2 Probe Types

```yaml
# HTTP GET - Most common for web apps
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080

# TCP Socket - Good for databases, non-HTTP services
livenessProbe:
  tcpSocket:
    port: 3306

# Exec Command - Run a command inside the container
livenessProbe:
  exec:
    command:
    - cat
    - /tmp/healthy
```

---

## Chapter 10: Resource Management

### 10.1 Requests and Limits

```yaml
resources:
  requests:          # "Minimum I need to run"
    memory: "256Mi"  # The scheduler uses this to place Pods
    cpu: "250m"      # 250 millicores = 0.25 CPU
  
  limits:            # "Maximum I'm allowed to use"
    memory: "512Mi"  # If exceeded: Pod is killed (OOMKilled)
    cpu: "500m"      # If exceeded: Pod is throttled (slowed down)
```

```
Memory Units:          CPU Units:
- Ki, Mi, Gi, Ti      - 1 = 1 full CPU core
- K, M, G, T          - 100m = 0.1 CPU (100 millicores)
                      - 500m = 0.5 CPU
                      - 2000m = 2 CPUs
```

### 10.2 Quality of Service (QoS) Classes

Based on your resource settings, Pods get a QoS class:

**Guaranteed** (Best treatment)
```yaml
# requests == limits for both CPU and memory
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "256Mi"
    cpu: "250m"
```

**Burstable** (Middle ground)
```yaml
# requests < limits, or only one set
resources:
  requests:
    memory: "256Mi"
  limits:
    memory: "512Mi"
```

**BestEffort** (First to be killed)
```yaml
# No requests or limits set
resources: {}
```

When the Node runs out of resources, Kubernetes kills Pods in this order:
1. BestEffort (killed first)
2. Burstable  
3. Guaranteed (killed last)

---

## Chapter 11: Scaling

### 11.1 Manual Scaling

```bash
# Scale up
kubectl scale deployment/nginx-deployment --replicas=5

# Scale down  
kubectl scale deployment/nginx-deployment --replicas=2

# Check current replicas
kubectl get deployment nginx-deployment
```

### 11.2 Horizontal Pod Autoscaler (HPA)

Automatically scale based on CPU/memory usage:

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: webapp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: webapp
  
  minReplicas: 2       # Never go below 2
  maxReplicas: 10      # Never go above 10
  
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70   # Target 70% CPU
```

```bash
# Apply HPA
kubectl apply -f hpa.yaml

# Check HPA status
kubectl get hpa
# NAME         REFERENCE           TARGETS   MINPODS   MAXPODS   REPLICAS
# webapp-hpa   Deployment/webapp   45%/70%   2         10        3
```

---

## Summary: What We Learned in Part 2

1. **Raw Pods are fragile** - They die and don't come back

2. **ReplicaSets** maintain a desired number of identical Pods

3. **Deployments** manage ReplicaSets and enable:
   - Rolling updates (zero downtime!)
   - Easy rollbacks
   - Scaling up/down

4. **Health Probes**:
   - `livenessProbe`: Is it alive? (restart if not)
   - `readinessProbe`: Is it ready for traffic?
   - `startupProbe`: Has it finished starting?

5. **Resources**:
   - `requests`: Minimum guaranteed
   - `limits`: Maximum allowed

6. **Key Commands**:
   - `kubectl rollout status deployment/X` - Watch rollout
   - `kubectl rollout undo deployment/X` - Rollback
   - `kubectl scale deployment/X --replicas=N` - Scale

---

## Exercises

1. Create a Deployment with 3 replicas of `httpd:2.4`
2. Update it to `httpd:2.4.58` and watch the rollout
3. Roll back to the previous version
4. Add health probes to your Deployment
5. Create an HPA that scales between 2-5 replicas based on CPU

---

**Coming up in Part 3**: Services and Networking - How Pods find and talk to each other!
