# Kubernetes Crash Course: From Zero to Production

## Part 5: Storage (Persistent Data)

---

## Chapter 24: The Storage Problem

### 24.1 Containers Are Ephemeral

By default, when a container restarts, all its data is GONE.

```
The Problem:
═══════════════════════════════════════════════════════════════

  Day 1:
    [Database Pod]
    └── data/
        ├── users.db      ← 1000 users
        └── orders.db     ← 5000 orders

  Day 2 (Pod crashes and restarts):
    [Database Pod]  
    └── data/
        └── (empty)       ← ALL DATA GONE! 😱
```

We need **persistent storage** that survives Pod restarts!

---

## Chapter 25: Volume Types

### 25.1 emptyDir (Temporary Shared Storage)

Exists only while Pod exists. Deleted when Pod is deleted.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: shared-storage
spec:
  containers:
  - name: writer
    image: busybox
    command: ['sh', '-c', 'echo "Hello!" > /data/hello.txt; sleep 3600']
    volumeMounts:
    - name: shared-data
      mountPath: /data
  
  - name: reader
    image: busybox
    command: ['sh', '-c', 'cat /data/hello.txt; sleep 3600']
    volumeMounts:
    - name: shared-data
      mountPath: /data
  
  volumes:
  - name: shared-data
    emptyDir: {}          # Temporary storage
```

**Use cases**: Scratch space, sharing files between containers in same Pod

### 25.2 hostPath (Node's Filesystem)

Uses a directory on the Node. **Dangerous in production!**

```yaml
volumes:
- name: host-data
  hostPath:
    path: /var/log        # Path on the Node
    type: Directory
```

**Problems**: Pod might run on different Node next time!

### 25.3 Persistent Volumes (The Right Way)

Kubernetes abstracts storage into two resources:

- **PersistentVolume (PV)**: The actual storage (created by admin)
- **PersistentVolumeClaim (PVC)**: A request for storage (created by user)

```
Storage Architecture:
═══════════════════════════════════════════════════════════════

    [ Admin creates storage ]           [ User requests storage ]
              │                                    │
              ▼                                    ▼
    ┌───────────────────┐             ┌───────────────────┐
    │ PersistentVolume  │◄────────────│     PVC           │
    │ (PV)              │   "binding" │ (request)         │
    │                   │             │                   │
    │ 100Gi SSD         │             │ "I need 50Gi"     │
    │ ReadWriteOnce     │             │                   │
    └───────────────────┘             └───────────────────┘
           │                                    │
           │                                    │
           ▼                                    ▼
    ┌───────────────────┐             ┌───────────────────┐
    │  Actual Storage   │             │      Pod          │
    │  (AWS EBS, GCE    │             │  mounts the PVC   │
    │   PD, NFS, etc.)  │             │                   │
    └───────────────────┘             └───────────────────┘
```

---

## Chapter 26: PersistentVolumes and Claims

### 26.1 Creating a PersistentVolume

```yaml
# pv.yaml (usually created by admin or dynamically)
apiVersion: v1
kind: PersistentVolume
metadata:
  name: my-pv
spec:
  capacity:
    storage: 10Gi
  
  accessModes:
    - ReadWriteOnce        # Can be mounted by ONE node
  
  persistentVolumeReclaimPolicy: Retain   # Keep data after PVC deleted
  
  storageClassName: standard
  
  # Backend storage (varies by environment)
  hostPath:                # For testing only!
    path: /mnt/data
```

### 26.2 Creating a PersistentVolumeClaim

```yaml
# pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  accessModes:
    - ReadWriteOnce
  
  resources:
    requests:
      storage: 5Gi
  
  storageClassName: standard
```

### 26.3 Using PVC in a Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: database
spec:
  containers:
  - name: postgres
    image: postgres:15
    volumeMounts:
    - name: db-storage
      mountPath: /var/lib/postgresql/data
  
  volumes:
  - name: db-storage
    persistentVolumeClaim:
      claimName: my-pvc      # Reference the PVC
```

### 26.4 Access Modes

| Mode | Description | Use Case |
|------|-------------|----------|
| ReadWriteOnce (RWO) | Mount by one node, read-write | Databases |
| ReadOnlyMany (ROX) | Mount by many nodes, read-only | Shared configs |
| ReadWriteMany (RWX) | Mount by many nodes, read-write | Shared files (NFS) |

---

## Chapter 27: StorageClasses (Dynamic Provisioning)

Instead of pre-creating PVs, let Kubernetes create them automatically!

```yaml
# storageclass.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs    # Cloud-specific
parameters:
  type: gp3
  fsType: ext4
reclaimPolicy: Delete                  # Delete PV when PVC deleted
volumeBindingMode: WaitForFirstConsumer
```

```yaml
# pvc-dynamic.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: database-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast-ssd          # Use the StorageClass
  resources:
    requests:
      storage: 50Gi
```

Now when you create this PVC, Kubernetes automatically provisions a 50Gi SSD!

---

## Chapter 28: StatefulSets (Stateful Applications)

For stateful apps (databases), use **StatefulSet** instead of Deployment.

StatefulSet provides:
- Stable network identity (pod-0, pod-1, pod-2)
- Stable storage (each Pod gets its own PVC)
- Ordered deployment and scaling

```yaml
# postgres-statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres      # Required for stable DNS
  replicas: 3
  selector:
    matchLabels:
      app: postgres
  
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
        env:
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
  
  # Each Pod gets its own PVC!
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: fast-ssd
      resources:
        requests:
          storage: 10Gi
```

```bash
# Created resources:
# Pods: postgres-0, postgres-1, postgres-2
# PVCs: data-postgres-0, data-postgres-1, data-postgres-2

# Stable DNS names:
# postgres-0.postgres.default.svc.cluster.local
# postgres-1.postgres.default.svc.cluster.local
# postgres-2.postgres.default.svc.cluster.local
```

---

## Chapter 29: Complete Database Example

```yaml
# Complete PostgreSQL deployment with persistent storage
---
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
type: Opaque
stringData:
  password: "super-secret-password"
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  clusterIP: None              # Headless service for StatefulSet
  selector:
    app: postgres
  ports:
  - port: 5432
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
```

---

## Summary: What We Learned in Part 5

1. **emptyDir**: Temporary, dies with Pod
2. **hostPath**: Node filesystem (dangerous in production)
3. **PersistentVolume (PV)**: The actual storage
4. **PersistentVolumeClaim (PVC)**: Request for storage
5. **StorageClass**: Dynamic provisioning
6. **StatefulSet**: For stateful applications

**Key Commands**:
- `kubectl get pv` - List PersistentVolumes
- `kubectl get pvc` - List PersistentVolumeClaims
- `kubectl get storageclass` - List StorageClasses

---

**Coming up in Part 6**: Ingress - HTTP routing and TLS!
