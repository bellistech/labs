# Kubernetes Crash Course: From Zero to Production

## Part 4: ConfigMaps & Secrets (External Configuration)

---

## Chapter 19: Why External Configuration?

### 19.1 The Problem with Hardcoding

```go
// BAD: Configuration hardcoded in your application
func main() {
    db, _ := sql.Open("postgres", 
        "host=prod-db.example.com user=admin password=secret123")
}
```

Problems:
- Different values for dev/staging/prod
- Need to rebuild for each environment
- Secrets in source code = security nightmare!

**Solution**: Put configuration OUTSIDE the container image.

```
Same Image, Different Config:
═══════════════════════════════════════════════════════════════

    ┌─────────────────┐
    │   myapp:1.0.0   │  ← Same Docker image everywhere!
    └────────┬────────┘
             │
    ┌────────┼────────┬────────────────┐
    ▼        ▼        ▼                ▼
┌───────┐┌───────┐┌───────┐      ┌───────┐
│  Dev  ││ Stage ││ Prod  │      │ Test  │
│Config ││Config ││Config │      │Config │
└───────┘└───────┘└───────┘      └───────┘
DB: local DB: stage DB: prod     DB: test
LOG: debug LOG:info LOG:warn     LOG:debug
```

---

## Chapter 20: ConfigMaps (Non-Sensitive Config)

A **ConfigMap** stores non-sensitive configuration data.

### 20.1 Creating ConfigMaps

```yaml
# configmap.yaml
#
# ELI5: "A bag of settings my app can read"
#
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  # Key-value pairs
  DATABASE_HOST: "postgres.default.svc"
  DATABASE_PORT: "5432"
  LOG_LEVEL: "info"
  
  # Or entire files
  app.properties: |
    database.pool.size=20
    cache.ttl=3600
    feature.newui=true
```

```bash
# Create from YAML
kubectl apply -f configmap.yaml

# Create from literal values
kubectl create configmap my-config \
  --from-literal=KEY1=value1 \
  --from-literal=KEY2=value2

# Create from file
kubectl create configmap my-config --from-file=config.properties

# View ConfigMap
kubectl get configmap app-config -o yaml
```

### 20.2 Using ConfigMaps in Pods

**Method 1: Environment Variables**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    
    # Load specific keys as env vars
    env:
    - name: DB_HOST
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: DATABASE_HOST
    
    # OR load ALL keys as env vars
    envFrom:
    - configMapRef:
        name: app-config
```

**Method 2: Volume Mounts (Files)**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    volumeMounts:
    - name: config-volume
      mountPath: /etc/config      # Files appear here
      readOnly: true
  
  volumes:
  - name: config-volume
    configMap:
      name: app-config
```

```bash
# Inside the Pod:
ls /etc/config/
# DATABASE_HOST  DATABASE_PORT  LOG_LEVEL  app.properties

cat /etc/config/DATABASE_HOST
# postgres.default.svc
```

---

## Chapter 21: Secrets (Sensitive Data)

**Secrets** are like ConfigMaps but for sensitive data (passwords, tokens, keys).

### 21.1 Creating Secrets

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
type: Opaque
stringData:              # Use stringData for plain text (easier)
  username: admin
  password: super-secret-password-123
```

```bash
# Create from YAML
kubectl apply -f secret.yaml

# Create from literals
kubectl create secret generic db-credentials \
  --from-literal=username=admin \
  --from-literal=password=secret123

# Create for Docker registry auth
kubectl create secret docker-registry my-registry \
  --docker-server=registry.example.com \
  --docker-username=user \
  --docker-password=pass

# View Secret (values are base64 encoded)
kubectl get secret db-credentials -o yaml

# Decode a secret value
kubectl get secret db-credentials -o jsonpath='{.data.password}' | base64 -d
```

### 21.2 Using Secrets in Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  containers:
  - name: myapp
    image: myapp:1.0
    
    env:
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: password
    
    # Or mount as files
    volumeMounts:
    - name: secret-volume
      mountPath: /etc/secrets
      readOnly: true
  
  volumes:
  - name: secret-volume
    secret:
      secretName: db-credentials
```

### 21.3 Secret Types

```yaml
type: Opaque                    # Generic key-value (default)
type: kubernetes.io/tls         # TLS certificates
type: kubernetes.io/dockerconfigjson  # Docker registry auth
type: kubernetes.io/basic-auth  # Basic authentication
```

---

## Chapter 22: Best Practices

### 22.1 ConfigMaps

✅ **Do**:
- Use for non-sensitive configuration
- Version your ConfigMaps (app-config-v1, app-config-v2)
- Use descriptive key names

❌ **Don't**:
- Store secrets in ConfigMaps
- Make ConfigMaps too large (1MB limit)

### 22.2 Secrets

✅ **Do**:
- Use for all sensitive data
- Enable encryption at rest (in production!)
- Use RBAC to limit secret access
- Consider external secret managers (Vault, AWS Secrets Manager)

❌ **Don't**:
- Commit secrets to git
- Log secret values
- Base64 encode = encryption (it's NOT!)

### 22.3 Updating Configuration

ConfigMap/Secret updates:
- **Volume mounts**: Update automatically (may take up to a minute)
- **Environment variables**: Require Pod restart

```bash
# Force Pod restart after env var config change
kubectl rollout restart deployment/myapp
```

---

## Chapter 23: Complete Example

```yaml
# Full application with ConfigMap and Secret
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp-config
data:
  DATABASE_HOST: "postgres.default.svc"
  DATABASE_PORT: "5432"
  LOG_LEVEL: "info"
---
apiVersion: v1
kind: Secret
metadata:
  name: webapp-secrets
type: Opaque
stringData:
  DATABASE_USER: "webapp"
  DATABASE_PASSWORD: "super-secret-password"
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
        image: mywebapp:1.0
        ports:
        - containerPort: 8080
        
        # Load all config as env vars
        envFrom:
        - configMapRef:
            name: webapp-config
        - secretRef:
            name: webapp-secrets
        
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
```

---

## Summary: What We Learned in Part 4

1. **ConfigMaps**: Non-sensitive configuration
2. **Secrets**: Sensitive data (passwords, keys, tokens)
3. **Usage**: Environment variables or volume mounts
4. **Best Practices**: Version configs, limit secret access, encrypt at rest

**Key Commands**:
- `kubectl create configmap X --from-literal=key=value`
- `kubectl create secret generic X --from-literal=key=value`
- `kubectl get secret X -o jsonpath='{.data.key}' | base64 -d`

---

**Coming up in Part 5**: Storage - Persistent Volumes and StatefulSets!
