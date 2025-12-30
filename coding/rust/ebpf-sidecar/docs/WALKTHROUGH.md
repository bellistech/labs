# eBPF Sidecar: Complete Walkthrough

A comprehensive guide to understanding and building eBPF-based network observability tools with Rust.

## Table of Contents

1. [Introduction to eBPF](#1-introduction-to-ebpf)
2. [Why Rust for eBPF?](#2-why-rust-for-ebpf)
3. [Project Architecture](#3-project-architecture)
4. [Understanding the Code](#4-understanding-the-code)
5. [Attaching to Processes](#5-attaching-to-processes)
6. [Exporting Metrics](#6-exporting-metrics)
7. [Production Considerations](#7-production-considerations)
8. [Extending the Sidecar](#8-extending-the-sidecar)

---

## 1. Introduction to eBPF

### What is eBPF?

eBPF (extended Berkeley Packet Filter) is a revolutionary Linux kernel technology that allows you to run sandboxed programs inside the kernel without modifying kernel source code or loading kernel modules.

Think of it as "JavaScript for the kernel" - you write small programs that the kernel verifies for safety, then runs at specific hook points.

### eBPF vs Traditional Approaches

```
Traditional Monitoring (ptrace, strace):
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    App      │ ──► │   Kernel    │ ──► │  Tracer     │
│             │ ◄── │             │ ◄── │  (userspace)│
└─────────────┘     └─────────────┘     └─────────────┘
                         ▲
                    Context switches!
                    High overhead!

eBPF Monitoring:
┌─────────────┐     ┌──────────────────────────────────┐
│    App      │ ──► │   Kernel + eBPF programs         │
│             │ ◄── │   (runs in-kernel, no switches)  │
└─────────────┘     └──────────────────────────────────┘
                         ▲
                    Zero context switches!
                    Minimal overhead!
```

### eBPF Attach Points

| Attach Point | Description | Use Case |
|--------------|-------------|----------|
| **kprobes** | Hook any kernel function | Trace syscalls, TCP functions |
| **kretprobes** | Hook function return | Capture return values |
| **tracepoints** | Predefined stable hooks | tcp_retransmit_skb, sched_switch |
| **XDP** | Packet at NIC level | DDoS mitigation, load balancing |
| **TC** | Traffic control | Packet modification |
| **uprobe** | Hook userspace functions | Application tracing |
| **socket filters** | Socket-level filtering | Per-socket monitoring |

### eBPF Maps

Maps are key-value stores shared between kernel eBPF programs and userspace:

```
┌───────────────────┐          ┌───────────────────┐
│   Kernel Space    │          │   User Space      │
│                   │          │                   │
│  eBPF Program     │          │  Rust Program     │
│       │           │          │       │           │
│       ▼           │          │       ▼           │
│  ┌─────────┐      │          │  ┌─────────┐      │
│  │   Map   │◄─────┼──────────┼──►   Map   │      │
│  │ (shared)│      │          │  │ (same!) │      │
│  └─────────┘      │          │  └─────────┘      │
└───────────────────┘          └───────────────────┘
```

Map types:
- **HashMap** - Key-value lookup (our CONNECTIONS map)
- **Array** - Fixed-size array (our CONFIG map)
- **PerfEventArray** - Send events to userspace
- **RingBuffer** - Efficient event streaming
- **LRU HashMap** - Auto-evicting cache

---

## 2. Why Rust for eBPF?

### The Aya Framework

Aya is a pure-Rust eBPF library. Unlike other solutions (BCC, libbpf), it doesn't require C toolchain.

| Aspect | C + libbpf | Rust + Aya |
|--------|------------|------------|
| Kernel code | C | Rust |
| Userspace code | C/Go/Python | Rust |
| Type sharing | Manual structs | Shared crate |
| Build system | Clang + make | Cargo |
| Memory safety | Manual | Compiler-enforced |

### Project Structure Explained

```
ebpf-sidecar/
├── Cargo.toml              # Workspace - ties everything together
│
├── sidecar-common/         # SHARED TYPES
│   └── src/lib.rs          # ConnKey, ConnMetrics structs
│                           # Used by BOTH kernel and userspace
│                           # #![no_std] - works in kernel
│
├── sidecar-ebpf/           # KERNEL CODE
│   ├── Cargo.toml          # Depends on aya-ebpf
│   └── src/main.rs         # eBPF programs (kprobes)
│                           # Compiled to BPF bytecode
│                           # Runs INSIDE the kernel
│
├── sidecar/                # USERSPACE CODE
│   ├── Cargo.toml          # Depends on aya, prometheus
│   └── src/main.rs         # Loads eBPF, exports metrics
│                           # Normal Rust binary
│                           # Runs as regular process
│
└── xtask/                  # BUILD TOOLS
    └── src/main.rs         # cargo xtask build
```

### The Magic of Shared Types

```rust
// sidecar-common/src/lib.rs
#![no_std]  // Works in kernel (no standard library)

#[repr(C)]  // C-compatible memory layout
pub struct ConnKey {
    pub src_ip: u32,
    pub dst_ip: u32,
    pub src_port: u16,
    pub dst_port: u16,
}

// In kernel (sidecar-ebpf):
use sidecar_common::ConnKey;
CONNECTIONS.insert(&key, &metrics, 0)?;

// In userspace (sidecar):
use sidecar_common::ConnKey;
for (key, value) in connections.iter() { ... }
```

Same struct definition, guaranteed to match!

---

## 3. Project Architecture

### Data Flow

```
1. Application makes TCP connection
         │
         ▼
2. Kernel calls tcp_connect()
         │
         ▼
3. Our eBPF kprobe fires
   ┌─────────────────────────────┐
   │ fn trace_tcp_connect()      │
   │   key = extract_conn_info() │
   │   CONNECTIONS.insert(key)   │
   └─────────────────────────────┘
         │
         ▼
4. Application sends data
         │
         ▼
5. Kernel calls tcp_sendmsg()
         │
         ▼
6. Our eBPF kprobe fires
   ┌─────────────────────────────┐
   │ fn trace_tcp_sendmsg()      │
   │   metrics.bytes_sent += len │
   └─────────────────────────────┘
         │
         ▼
7. Userspace reads CONNECTIONS map every N seconds
   ┌─────────────────────────────┐
   │ for (key, metrics) in map   │
   │   prometheus.update(metrics)│
   └─────────────────────────────┘
         │
         ▼
8. Prometheus scrapes /metrics endpoint
```

### Components

| Component | Location | Runs In | Purpose |
|-----------|----------|---------|---------|
| eBPF Programs | sidecar-ebpf | Kernel | Hook TCP functions |
| CONNECTIONS map | Kernel memory | Kernel | Store per-conn metrics |
| CONFIG map | Kernel memory | Kernel | Runtime configuration |
| Loader | sidecar/main.rs | Userspace | Load eBPF, attach probes |
| Exporter | sidecar/main.rs | Userspace | Read maps, serve Prometheus |

---

## 4. Understanding the Code

### Kernel Side (sidecar-ebpf/src/main.rs)

#### Defining Maps

```rust
// HashMap: key=ConnKey, value=ConnMetrics
#[map]
static CONNECTIONS: HashMap<ConnKey, ConnMetrics> =
    HashMap::with_max_entries(10240, BPF_F_NO_PREALLOC);
```

- `#[map]` - Aya macro to define eBPF map
- `10240` - Maximum entries (tune based on expected connections)
- `BPF_F_NO_PREALLOC` - Allocate on demand (saves memory)

#### Kprobe Handler

```rust
#[kprobe]
pub fn trace_tcp_connect(ctx: ProbeContext) -> u32 {
    match try_trace_tcp_connect(&ctx) {
        Ok(()) => 0,   // Success
        Err(_) => 1,   // Error (logged elsewhere)
    }
}

fn try_trace_tcp_connect(ctx: &ProbeContext) -> Result<(), i64> {
    // Get first argument (struct sock *)
    let sock: *const u8 = ctx.arg(0).ok_or(1i64)?;
    
    // Extract connection 4-tuple from sock struct
    let key = unsafe { read_conn_key_from_sock(sock)? };
    
    // Create initial metrics
    let metrics = ConnMetrics {
        bytes_sent: 0,
        start_ns: unsafe { bpf_ktime_get_ns() },
        ...
    };
    
    // Insert into map
    CONNECTIONS.insert(&key, &metrics, 0)?;
    
    Ok(())
}
```

#### Reading Kernel Structs

```rust
unsafe fn read_conn_key_from_sock(sock: *const u8) -> Result<ConnKey, i64> {
    // struct sock layout (simplified):
    // offset 0:  __sk_common
    //   offset 0:  skc_daddr (dest IP)
    //   offset 4:  skc_rcv_saddr (source IP)
    //   offset 12: skc_dport (dest port, network order)
    //   offset 14: skc_num (source port)
    
    let src_ip = bpf_probe_read_kernel(sock.add(4) as *const u32)?;
    let dst_ip = bpf_probe_read_kernel(sock.add(0) as *const u32)?;
    // ...
}
```

**Note:** These offsets are kernel-version specific. Production code should use CO-RE (Compile Once Run Everywhere) with BTF.

### Userspace Side (sidecar/src/main.rs)

#### Loading eBPF

```rust
// Include compiled eBPF bytecode AT COMPILE TIME
let mut bpf = Bpf::load(include_bytes_aligned!(
    "../../target/bpfel-unknown-none/debug/sidecar"
))?;
```

The eBPF bytecode is embedded in the binary - no separate file needed.

#### Attaching Programs

```rust
// Get program by name (from #[kprobe] function name)
let program: &mut KProbe = bpf
    .program_mut("trace_tcp_connect")?
    .try_into()?;

// Load into kernel (verified by BPF verifier)
program.load()?;

// Attach to kernel function
program.attach("tcp_connect", 0)?;
```

#### Reading Maps

```rust
// Get typed reference to map
let connections: HashMap<_, ConnKey, ConnMetrics> =
    HashMap::try_from(bpf.map("CONNECTIONS")?)?;

// Iterate all entries
for result in connections.iter() {
    let (key, metrics) = result?;
    // key: ConnKey, metrics: ConnMetrics
    // Update Prometheus counters...
}
```

---

## 5. Attaching to Processes

### Understanding Attachment

eBPF programs attach to **kernel functions**, not processes. When `tcp_connect` is called by ANY process, our probe fires.

Filtering happens INSIDE the eBPF program:

```rust
fn should_trace(ctx: &impl EbpfContext) -> bool {
    let config = CONFIG.get(0)?;
    
    // Get current process ID
    let pid = (bpf_get_current_pid_tgid() >> 32) as u32;
    
    // Filter by configured PID
    if config.target_pid != 0 && pid != config.target_pid {
        return false;  // Skip this process
    }
    
    true
}
```

### Method 1: PID Filtering

```bash
# Find your app's PID
$ pgrep nginx
12345

# Run sidecar with filter
$ sudo ./sidecar --pid 12345
```

The eBPF program will ignore all connections not from PID 12345.

### Method 2: Port Filtering

```bash
# Monitor only HTTP/HTTPS
$ sudo ./sidecar --ports 80,443

# Monitor database connections
$ sudo ./sidecar --ports 5432,3306
```

### Method 3: cgroup Filtering (Containers)

Every container runs in a cgroup. You can filter by cgroup ID:

```bash
# Find container's cgroup
$ cat /proc/$(docker inspect -f '{{.State.Pid}}' mycontainer)/cgroup
0::/docker/abc123def456...

# Configure in config.yaml
target:
  cgroup: "/docker/abc123def456"
```

### Method 4: Process Name (requires uprobe)

```rust
// Advanced: attach uprobe to specific binary
let program: &mut UProbe = bpf.program_mut("trace_app_function")?.try_into()?;
program.load()?;
program.attach(Some("my_function"), 0, "/usr/bin/my-app", None)?;
```

### Kubernetes Sidecar Pattern

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  shareProcessNamespace: true  # Share PID namespace
  containers:
  - name: app
    image: my-app:latest
  - name: ebpf-sidecar
    image: ebpf-sidecar:latest
    securityContext:
      privileged: true  # Required for eBPF
    volumeMounts:
    - name: sys
      mountPath: /sys
      readOnly: true
  volumes:
  - name: sys
    hostPath:
      path: /sys
```

---

## 6. Exporting Metrics

### Prometheus Integration

The sidecar exposes a standard Prometheus `/metrics` endpoint:

```rust
// Define Prometheus metrics
lazy_static! {
    static ref BYTES_SENT: CounterVec = register_counter_vec!(
        "sidecar_connection_bytes_sent_total",
        "Total bytes sent",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();
}

// Update from eBPF map data
for (key, metrics) in connections.iter() {
    BYTES_SENT
        .with_label_values(&[&src_ip, &dst_ip, &port])
        .inc_by(metrics.bytes_sent as f64);
}

// Serve HTTP
async fn metrics_handler() -> Response {
    let encoder = TextEncoder::new();
    let families = prometheus::gather();
    encoder.encode(&families, &mut buffer);
    Response::new(Body::from(buffer))
}
```

### Prometheus Configuration

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'ebpf-sidecar'
    static_configs:
      - targets: ['localhost:9090']
    scrape_interval: 15s
```

### Alternative: Push to Pushgateway

```rust
use prometheus::push_metrics;

// Push every N seconds
push_metrics(
    "ebpf_sidecar",
    labels,
    "http://pushgateway:9091",
    prometheus::gather(),
)?;
```

### Alternative: OpenTelemetry

```rust
use opentelemetry::metrics::MeterProvider;
use opentelemetry_otlp::WithExportConfig;

let provider = MeterProvider::builder()
    .with_reader(PeriodicReader::builder(
        opentelemetry_otlp::new_exporter()
            .tonic()
            .with_endpoint("http://otel-collector:4317")
    ).build())
    .build();

let meter = provider.meter("ebpf-sidecar");
let bytes_counter = meter.u64_counter("connection.bytes_sent").init();
```

### Grafana Dashboard

Example queries:

```promql
# Bytes per second by destination
rate(sidecar_connection_bytes_sent_total[5m])

# Top 10 destinations by traffic
topk(10, sum by (dst_ip) (rate(sidecar_connection_bytes_sent_total[5m])))

# Retransmit rate (network quality)
rate(sidecar_connection_retransmits_total[5m]) 
  / rate(sidecar_connection_packets_sent_total[5m])

# Connection duration histogram
histogram_quantile(0.99, sidecar_connection_duration_seconds)
```

---

## 7. Production Considerations

### Kernel Compatibility

eBPF features vary by kernel version:

| Feature | Minimum Kernel |
|---------|----------------|
| Basic kprobes | 4.1 |
| HashMap maps | 4.1 |
| Per-CPU maps | 4.6 |
| BTF (CO-RE) | 5.2 |
| Ring buffer | 5.8 |
| bpf_loop | 5.17 |

Check your kernel:
```bash
uname -r  # Should be 5.8+ for full features
```

### CO-RE (Compile Once Run Everywhere)

Our code uses hardcoded struct offsets. For production:

```rust
// Instead of hardcoded offsets:
let src_ip = bpf_probe_read_kernel(sock.add(4) as *const u32)?;

// Use CO-RE with BTF:
use aya_ebpf::helpers::bpf_probe_read_kernel;

#[repr(C)]
struct sock {
    __sk_common: sock_common,
}

#[repr(C)]
struct sock_common {
    skc_daddr: u32,
    skc_rcv_saddr: u32,
    // ... BTF-derived layout
}
```

### Performance Tuning

```yaml
# config.yaml
metrics:
  # Reduce map reads for high-traffic systems
  interval_secs: 15
  
  # Limit tracked connections
  max_connections: 50000
```

### Security

```bash
# Instead of root, use capabilities:
sudo setcap cap_bpf,cap_perfmon,cap_net_admin+ep ./sidecar

# Or in Kubernetes:
securityContext:
  capabilities:
    add: ["BPF", "PERFMON", "NET_ADMIN"]
```

---

## 8. Extending the Sidecar

### Adding HTTP L7 Metrics

```rust
// In eBPF: Parse HTTP from socket data
#[kprobe]
pub fn trace_tcp_recvmsg(ctx: ProbeContext) -> u32 {
    let data = read_socket_data(&ctx)?;
    
    if data.starts_with(b"HTTP/") {
        let status = parse_http_status(data)?;
        let event = HttpEvent { status, ... };
        EVENTS.output(&ctx, &event, 0);
    }
}

// In userspace: Read perf events
let mut events = AsyncPerfEventArray::try_from(bpf.map_mut("EVENTS")?)?;
let mut buf = events.open(0, None)?;

loop {
    let events = buf.read_events(&mut buffers).await?;
    for event in events {
        let http: HttpEvent = parse(&event);
        HTTP_REQUESTS.inc();
    }
}
```

### Adding TLS Interception (uprobe)

```rust
// Attach to OpenSSL
#[uprobe]
pub fn trace_ssl_write(ctx: ProbeContext) -> u32 {
    // Called when app writes to TLS connection
    let data_ptr: *const u8 = ctx.arg(1)?;
    let len: usize = ctx.arg(2)?;
    // Now we can see plaintext!
}

// Attach to specific library
program.attach(Some("SSL_write"), 0, "/usr/lib/libssl.so", None)?;
```

### Adding DNS Monitoring

```rust
// Hook UDP sendto for DNS queries
#[kprobe]
pub fn trace_udp_sendmsg(ctx: ProbeContext) -> u32 {
    let sock: *const u8 = ctx.arg(0)?;
    let port = get_dest_port(sock)?;
    
    if port == 53 {
        // Parse DNS query
        let query = parse_dns_query(&ctx)?;
        DNS_QUERIES.insert(&query, &1, 0)?;
    }
}
```

---

## Summary

You now understand:

1. **eBPF fundamentals** - Kernel programs, maps, attach points
2. **Rust + Aya** - Type-safe eBPF development
3. **Project structure** - Shared types, kernel code, userspace loader
4. **Process attachment** - PID, port, cgroup filtering
5. **Metrics export** - Prometheus, OTLP, Grafana

### Next Steps

1. Run the sidecar on a test system
2. Add HTTP L7 parsing
3. Integrate with your Prometheus/Grafana stack
4. Deploy as Kubernetes sidecar
5. Explore XDP for packet-level processing

### Resources

- [Aya Book](https://aya-rs.dev/book/)
- [eBPF.io](https://ebpf.io/)
- [Cilium eBPF Guide](https://docs.cilium.io/en/latest/bpf/)
- [Linux Kernel BTF](https://www.kernel.org/doc/html/latest/bpf/btf.html)
