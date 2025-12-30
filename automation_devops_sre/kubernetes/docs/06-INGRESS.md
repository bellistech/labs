# Kubernetes Crash Course: From Zero to Production

## Part 6: Ingress (HTTP Routing and TLS)

---

## Chapter 30: Why Ingress?

### 30.1 The Problem with LoadBalancer Services

Each LoadBalancer Service creates a separate cloud load balancer:

```
Without Ingress:
═══════════════════════════════════════════════════════════════

    api.example.com     app.example.com     admin.example.com
          │                   │                    │
          ▼                   ▼                    ▼
    ┌──────────┐        ┌──────────┐        ┌──────────┐
    │  LB #1   │        │  LB #2   │        │  LB #3   │
    │  $$$     │        │  $$$     │        │  $$$     │
    └────┬─────┘        └────┬─────┘        └────┬─────┘
         │                   │                    │
         ▼                   ▼                    ▼
    ┌──────────┐        ┌──────────┐        ┌──────────┐
    │  API     │        │  App     │        │  Admin   │
    │  Service │        │  Service │        │  Service │
    └──────────┘        └──────────┘        └──────────┘

    Problem: 3 load balancers = 3x the cost! 💸
```

### 30.2 Ingress: One Entry Point

**Ingress** gives you ONE load balancer that routes based on host/path:

```
With Ingress:
═══════════════════════════════════════════════════════════════

    api.example.com     app.example.com     admin.example.com
          │                   │                    │
          └─────────────┬─────┴─────────────┬─────┘
                        │                   │
                        ▼                   ▼
              ┌───────────────────────────────────┐
              │         INGRESS CONTROLLER        │
              │         (One Load Balancer)       │
              │              $ (cheaper!)         │
              └─────────────────┬─────────────────┘
                    ┌───────────┼───────────┐
                    │           │           │
                    ▼           ▼           ▼
              ┌─────────┐ ┌─────────┐ ┌─────────┐
              │   API   │ │   App   │ │  Admin  │
              │ Service │ │ Service │ │ Service │
              └─────────┘ └─────────┘ └─────────┘
```

---

## Chapter 31: Ingress Controller

Ingress resources need an **Ingress Controller** to work. Popular options:

- **nginx-ingress** (Most common)
- **Traefik**
- **HAProxy**
- Cloud-specific: AWS ALB, GCP Ingress

### 31.1 Installing nginx-ingress

```bash
# Using Helm
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx

# Or using manifests
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

# Check it's running
kubectl get pods -n ingress-nginx
kubectl get service ingress-nginx-controller -n ingress-nginx
```

---

## Chapter 32: Creating Ingress Resources

### 32.1 Basic Host-Based Routing

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx      # Which controller to use
  
  rules:
  - host: api.example.com      # Route based on hostname
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
  
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
```

### 32.2 Path-Based Routing

```yaml
# Route different paths to different services
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: path-ingress
spec:
  ingressClassName: nginx
  rules:
  - host: example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
      
      - path: /app
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
      
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

### 32.3 Path Types

| Type | Behavior |
|------|----------|
| `Exact` | Must match exactly |
| `Prefix` | Matches path prefix |
| `ImplementationSpecific` | Controller decides |

---

## Chapter 33: TLS / HTTPS

### 33.1 Manual TLS Certificate

```yaml
# Create a Secret with your certificate
apiVersion: v1
kind: Secret
metadata:
  name: tls-secret
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-cert>
  tls.key: <base64-encoded-key>
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tls-ingress
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - app.example.com
    secretName: tls-secret    # Reference the Secret
  
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
```

### 33.2 Automatic TLS with cert-manager

**cert-manager** automatically gets Let's Encrypt certificates!

```bash
# Install cert-manager
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true
```

```yaml
# ClusterIssuer for Let's Encrypt
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
---
# Ingress with automatic certificate
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: auto-tls-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - app.example.com
    secretName: app-tls        # cert-manager creates this!
  
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
```

---

## Chapter 34: Useful Annotations

```yaml
metadata:
  annotations:
    # Redirect HTTP to HTTPS
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    
    # Custom timeout
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    
    # Rate limiting
    nginx.ingress.kubernetes.io/limit-rps: "100"
    
    # Enable CORS
    nginx.ingress.kubernetes.io/enable-cors: "true"
    
    # Basic auth
    nginx.ingress.kubernetes.io/auth-type: basic
    nginx.ingress.kubernetes.io/auth-secret: basic-auth
    
    # Custom backend protocol
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
```

---

## Chapter 35: Debugging Ingress

```bash
# Check Ingress status
kubectl get ingress
kubectl describe ingress my-ingress

# Check Ingress Controller logs
kubectl logs -n ingress-nginx deploy/ingress-nginx-controller

# Test DNS resolution
nslookup app.example.com

# Test endpoint directly
curl -H "Host: app.example.com" http://<ingress-controller-ip>/

# Check certificate
openssl s_client -connect app.example.com:443 -servername app.example.com
```

---

## Summary: What We Learned in Part 6

1. **Ingress**: Single entry point for HTTP/HTTPS traffic
2. **Ingress Controller**: nginx, Traefik, or cloud-native
3. **Routing**: Host-based and path-based
4. **TLS**: Manual or automatic with cert-manager
5. **Annotations**: Customize behavior per-ingress

**Key Commands**:
- `kubectl get ingress`
- `kubectl describe ingress X`
- `kubectl logs -n ingress-nginx deploy/ingress-nginx-controller`

---

**Coming up in Part 7**: Helm - Package management for Kubernetes!
