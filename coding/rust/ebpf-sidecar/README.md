# eBPF Sidecar

A service mesh sidecar built with Rust and eBPF for transparent network observability. Unlike traditional proxy-based sidecars (Envoy, Linkerd), this operates at the kernel level with zero application changes and minimal latency overhead.

## Features

- **Zero-latency observability** - eBPF hooks run in-kernel, no userspace proxy hop
- **Per-connection metrics** - Bytes, packets, retransmits, duration
- **Prometheus export** - Native `/metrics` endpoint for scraping
- **Process filtering** - Monitor specific PIDs, processes, or containers
- **Port filtering** - Focus on specific services/ports
- **Pure Rust** - Both kernel (eBPF) and userspace components

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    User Space                           │
│  ┌─────────────┐                    ┌───────────────┐   │
│  │ Your App    │                    │ Sidecar       │   │
│  │ (unchanged) │                    │ (metrics      │   │
│  └──────┬──────┘                    │  exporter)    │   │
│         │                           └───────┬───────┘   │
│         │                                   │           │
│         │ tcp_connect/send/recv             │ read maps │
│         ▼                                   ▼           │
├─────────────────────────────────────────────────────────┤
│                    Kernel Space                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │              eBPF Programs                        │   │
│  │  • trace_tcp_connect  → CONNECTIONS map          │   │
│  │  • trace_tcp_sendmsg  → update bytes_sent        │   │
│  │  • trace_tcp_recvmsg  → update bytes_recv        │   │
│  │  • trace_tcp_close    → cleanup                  │   │
│  │  • trace_retransmit   → network quality          │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y \
    build-essential \
    pkg-config \
    libssl-dev \
    linux-headers-$(uname -r) \
    clang \
    llvm

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# Install nightly toolchain (required for eBPF)
rustup install nightly
rustup component add rust-src --toolchain nightly

# Install bpf-linker
cargo install bpf-linker
```

### Build

```bash
# Build everything
cargo xtask build

# Build release version
cargo xtask build --release

# Build only eBPF programs
cargo xtask build-ebpf
```

### Run

```bash
# Monitor all processes (requires root)
sudo ./target/debug/sidecar

# Monitor specific PID
sudo ./target/debug/sidecar --pid 1234

# Monitor specific ports
sudo ./target/debug/sidecar --ports 80,443,8080

# Custom metrics port
sudo ./target/debug/sidecar --metrics-port 9091

# Enable debug logging
sudo ./target/debug/sidecar --debug
```

### Scrape Metrics

```bash
# Prometheus metrics endpoint
curl http://localhost:9090/metrics

# Health check
curl http://localhost:9090/health
```

## Attaching to a Process

### Method 1: PID Filtering

```bash
# Find your application's PID
pgrep -f "my-application"
# Output: 12345

# Run sidecar with PID filter
sudo ./target/debug/sidecar --pid 12345
```

### Method 2: Port Filtering

```bash
# Monitor only web traffic
sudo ./target/debug/sidecar --ports 80,443

# Monitor database connections
sudo ./target/debug/sidecar --ports 5432,3306,6379
```

### Method 3: Container/cgroup Filtering (Advanced)

For Kubernetes or Docker, you can filter by cgroup:

```bash
# Find container's cgroup
cat /proc/<container-pid>/cgroup

# Use in config.yaml
target:
  cgroup: "/sys/fs/cgroup/system.slice/docker-abc123.scope"
```

### Method 4: Kubernetes Sidecar

```yaml
# Add as a sidecar container in your Pod
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: my-app
    image: my-app:latest
  - name: ebpf-sidecar
    image: bellistech/ebpf-sidecar:latest
    securityContext:
      privileged: true  # Required for eBPF
    ports:
    - containerPort: 9090
      name: metrics
```

## Prometheus Integration

### prometheus.yml

```yaml
scrape_configs:
  - job_name: 'ebpf-sidecar'
    static_configs:
      - targets: ['localhost:9090']
    scrape_interval: 15s
```

### Available Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `sidecar_connection_bytes_sent_total` | Counter | Total bytes sent per connection |
| `sidecar_connection_bytes_received_total` | Counter | Total bytes received per connection |
| `sidecar_connection_packets_sent_total` | Counter | Total packets sent |
| `sidecar_connection_packets_received_total` | Counter | Total packets received |
| `sidecar_connection_retransmits_total` | Counter | TCP retransmissions (network quality) |
| `sidecar_connection_duration_seconds` | Gauge | Connection duration |
| `sidecar_active_connections` | Gauge | Currently tracked connections |

### Example Queries

```promql
# Bytes per second to a destination
rate(sidecar_connection_bytes_sent_total{dst_ip="10.0.0.5"}[5m])

# High retransmit rate (network issues)
rate(sidecar_connection_retransmits_total[5m]) > 10

# Connection count by destination port
count by (dst_port) (sidecar_connection_duration_seconds)
```

## Grafana Dashboard

Import the included dashboard or create panels:

```json
{
  "panels": [
    {
      "title": "Bytes Sent/Received",
      "type": "graph",
      "targets": [
        {"expr": "rate(sidecar_connection_bytes_sent_total[5m])"},
        {"expr": "rate(sidecar_connection_bytes_received_total[5m])"}
      ]
    }
  ]
}
```

## Project Structure

```
ebpf-sidecar/
├── Cargo.toml              # Workspace definition
├── config.yaml             # Example configuration
├── sidecar-common/         # Shared types (kernel & userspace)
│   └── src/lib.rs          # ConnKey, ConnMetrics, etc.
├── sidecar-ebpf/           # eBPF programs (runs in kernel)
│   └── src/main.rs         # Kprobes, tracepoints
├── sidecar/                # Userspace loader & exporter
│   └── src/
│       ├── main.rs         # CLI, eBPF loading, Prometheus
│       ├── config.rs       # YAML config parsing
│       └── metrics.rs      # Metrics aggregation
└── xtask/                  # Build tooling
    └── src/main.rs         # cargo xtask commands
```

## How It Works

### 1. eBPF Programs (Kernel)

Attach to kernel functions via kprobes:

- **tcp_connect** - New outbound connection → create entry in CONNECTIONS map
- **tcp_sendmsg** - Data sent → increment bytes_sent
- **tcp_recvmsg** - Data received → increment bytes_recv  
- **tcp_close** - Connection closed → log and cleanup
- **tcp_retransmit_skb** - Retransmit → increment counter

### 2. Shared Maps

eBPF maps are shared between kernel and userspace:

```rust
// In kernel (sidecar-ebpf)
CONNECTIONS.insert(&key, &metrics, 0)?;

// In userspace (sidecar)
for (key, metrics) in connections.iter() {
    // Export to Prometheus
}
```

### 3. Prometheus Export

Userspace periodically reads maps and updates Prometheus counters/gauges.

## Comparison with Traditional Sidecars

| Feature | eBPF Sidecar | Envoy/Linkerd |
|---------|--------------|---------------|
| Latency overhead | ~0 (in-kernel) | 1-5ms per hop |
| App changes | None | Config/inject |
| L7 parsing | Limited | Full HTTP/gRPC |
| mTLS | No | Yes |
| Traffic routing | No | Yes |
| Resource usage | Very low | Moderate |

**Use eBPF sidecar for:** Pure observability, latency-sensitive apps, legacy apps

**Use Envoy/Linkerd for:** Traffic management, mTLS, advanced routing

## Troubleshooting

### Permission denied
```bash
# eBPF requires root or CAP_BPF+CAP_PERFMON
sudo ./target/debug/sidecar
# Or with capabilities
sudo setcap cap_bpf,cap_perfmon+ep ./target/debug/sidecar
```

### BPF program rejected
```bash
# Check kernel version (need 5.8+)
uname -r

# Check BTF support
ls /sys/kernel/btf/vmlinux
```

### No metrics appearing
```bash
# Verify eBPF programs are loaded
sudo bpftool prog list | grep sidecar

# Check maps have data
sudo bpftool map dump name CONNECTIONS
```

## License

MIT License - see LICENSE file.
