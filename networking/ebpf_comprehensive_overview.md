# eBPF Deep Dive for Network Engineers

## Table of Contents
1. [eBPF Fundamentals](#fundamentals)
2. [eBPF Architecture & Components](#architecture)
3. [XDP (eXpress Data Path)](#xdp)
4. [Network Programming with eBPF](#network-programming)
5. [Cilium & Kubernetes Networking](#cilium)
6. [Practical Labs & Examples](#labs)
7. [Performance & Observability](#performance)
8. [eBPF for Network Automation](#automation)
9. [Interview Questions & Answers](#interview)

---

## 1. eBPF Fundamentals {#fundamentals}

### What is eBPF?

eBPF (extended Berkeley Packet Filter) is a revolutionary technology that allows running sandboxed programs in the Linux kernel without changing kernel source code or loading kernel modules.

```
┌─────────────────────────────────────────────────┐
│                   User Space                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │   bpftool │  │    tc    │  │    ip    │    │
│  │  bpftrace │  │  bpf()   │  │  libbpf  │    │
│  └──────────┘  └──────────┘  └──────────┘    │
├─────────────────────────────────────────────────┤
│                  eBPF Subsystem                 │
│  ┌──────────────────────────────────────┐     │
│  │         eBPF Verifier                │     │
│  │  • Safety checks                     │     │
│  │  • Bounds checking                   │     │
│  │  • Type checking                     │     │
│  └──────────────────────────────────────┘     │
│  ┌──────────────────────────────────────┐     │
│  │         JIT Compiler                 │     │
│  │  • Compiles to native code          │     │
│  └──────────────────────────────────────┘     │
│  ┌──────────────────────────────────────┐     │
│  │         eBPF Programs                │     │
│  │  • XDP    • TC                       │     │
│  │  • Socket • Tracepoint               │     │
│  │  • Kprobe • Perf Event              │     │
│  └──────────────────────────────────────┘     │
│  ┌──────────────────────────────────────┐     │
│  │         eBPF Maps                    │     │
│  │  • Hash    • Array                   │     │
│  │  • LRU     • Stack                   │     │
│  │  • Queue   • Bloom Filter            │     │
│  └──────────────────────────────────────┘     │
├─────────────────────────────────────────────────┤
│                  Linux Kernel                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │   XDP    │  │    TC    │  │  Socket  │    │
│  │   Hook   │  │   Hook   │  │   Hook   │    │
│  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────┘
```

### Key Concepts

**1. Programs**: Small programs that run in kernel space
**2. Maps**: Data structures for storing state between program runs
**3. Helpers**: Kernel functions callable from eBPF programs
**4. Verifier**: Ensures programs are safe to run
**5. JIT**: Compiles eBPF bytecode to native machine code

### eBPF vs Classic BPF

| Feature | Classic BPF | eBPF |
|---------|------------|------|
| Registers | 2 (A, X) | 11 (R0-R10) |
| Instruction Set | 32-bit | 64-bit |
| Stack Size | Limited | 512 bytes |
| Maps | No | Yes |
| Helper Functions | Few | 200+ |
| Use Cases | Packet Filtering | Everything |

---

## 2. eBPF Architecture & Components {#architecture}

### eBPF Program Types for Networking

```c
// Program types relevant to networking
enum bpf_prog_type {
    BPF_PROG_TYPE_XDP,           // Earliest packet processing
    BPF_PROG_TYPE_SCHED_CLS,     // TC classifier
    BPF_PROG_TYPE_SCHED_ACT,     // TC action
    BPF_PROG_TYPE_SOCKET_FILTER, // Socket filtering
    BPF_PROG_TYPE_SK_SKB,        // Socket buffer processing
    BPF_PROG_TYPE_CGROUP_SKB,    // Cgroup socket buffer
    BPF_PROG_TYPE_LWT_IN,        // Lightweight tunnel in
    BPF_PROG_TYPE_LWT_OUT,       // Lightweight tunnel out
    BPF_PROG_TYPE_LWT_XMIT,      // Lightweight tunnel transmit
    BPF_PROG_TYPE_SOCK_OPS,      // Socket operations
    BPF_PROG_TYPE_SK_REUSEPORT,  // Port reuse logic
};
```

### eBPF Maps for Network State

```c
// Map types commonly used in networking
enum bpf_map_type {
    BPF_MAP_TYPE_HASH,           // Key-value store
    BPF_MAP_TYPE_ARRAY,          // Fixed-size array
    BPF_MAP_TYPE_PERCPU_HASH,    // Per-CPU hash
    BPF_MAP_TYPE_PERCPU_ARRAY,   // Per-CPU array
    BPF_MAP_TYPE_LRU_HASH,       // LRU eviction
    BPF_MAP_TYPE_LPM_TRIE,       // Longest prefix match
    BPF_MAP_TYPE_DEVMAP,         // Device map for XDP redirect
    BPF_MAP_TYPE_CPUMAP,         // CPU map for XDP redirect
    BPF_MAP_TYPE_XSKMAP,         // XDP socket map
    BPF_MAP_TYPE_SOCKHASH,       // Socket hash
    BPF_MAP_TYPE_SOCKMAP,        // Socket map
};
```

### eBPF Workflow

```bash
# 1. Write eBPF program (C)
cat > xdp_drop.c <<'EOF'
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("xdp")
int xdp_drop_prog(struct xdp_md *ctx) {
    return XDP_DROP;
}

char _license[] SEC("license") = "GPL";
EOF

# 2. Compile to eBPF bytecode
clang -O2 -target bpf -c xdp_drop.c -o xdp_drop.o

# 3. Load and attach
ip link set dev eth0 xdp obj xdp_drop.o sec xdp

# 4. Verify it's loaded
bpftool prog show
bpftool map show
```

---

## 3. XDP (eXpress Data Path) {#xdp}

### XDP Overview

XDP provides a programmable, high-performance network data path in the Linux kernel. It runs eBPF programs at the earliest possible point in the network stack - right after packet reception from the NIC driver.

```
┌──────────────────────────────────────────────┐
│                    NIC                       │
│                     ↓                         │
│              [DMA to Memory]                  │
│                     ↓                         │
│           ┌──────────────────┐               │
│           │   NIC Driver RX   │               │
│           └──────────────────┘               │
│                     ↓                         │
│           ┌──────────────────┐               │
│           │   XDP Hook        │ ← eBPF Here  │
│           └──────────────────┘               │
│            ↙    ↓      ↓    ↘                │
│      XDP_DROP  PASS  TX   REDIRECT           │
│                 ↓                             │
│           ┌──────────────────┐               │
│           │   Network Stack   │               │
│           └──────────────────┘               │
└──────────────────────────────────────────────┘
```

### XDP Actions

```c
enum xdp_action {
    XDP_ABORTED = 0,  // Error, packet dropped
    XDP_DROP,         // Drop packet immediately
    XDP_PASS,         // Pass to normal network stack
    XDP_TX,           // Transmit back on same interface
    XDP_REDIRECT,     // Redirect to another interface
};
```

### XDP Program Example: DDoS Mitigation

```c
// xdp_ddos_mitigate.c
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define MAX_RULES 1024
#define RATE_LIMIT 1000000  // 1M pps

// Map for blacklisted IPs
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_RULES);
    __type(key, __u32);   // IPv4 address
    __type(value, __u64); // Packet count
} blacklist SEC(".maps");

// Per-CPU rate limiting
struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u64);
} rate_limit SEC(".maps");

// Connection tracking
struct conn_key {
    __u32 src_ip;
    __u32 dst_ip;
    __u16 src_port;
    __u16 dst_port;
    __u8  protocol;
};

struct conn_state {
    __u64 packets;
    __u64 bytes;
    __u64 last_seen;
};

struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 100000);
    __type(key, struct conn_key);
    __type(value, struct conn_state);
} connections SEC(".maps");

SEC("xdp")
int xdp_ddos_filter(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    // Parse Ethernet header
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;
    
    // Only process IPv4
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return XDP_PASS;
    
    // Parse IP header
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_PASS;
    
    // Check blacklist
    __u32 src_ip = ip->saddr;
    __u64 *count = bpf_map_lookup_elem(&blacklist, &src_ip);
    if (count) {
        __sync_fetch_and_add(count, 1);
        return XDP_DROP;  // Drop blacklisted IPs
    }
    
    // Rate limiting
    __u32 key = 0;
    __u64 *rate_counter = bpf_map_lookup_elem(&rate_limit, &key);
    if (rate_counter) {
        if (*rate_counter > RATE_LIMIT) {
            return XDP_DROP;  // Rate limit exceeded
        }
        __sync_fetch_and_add(rate_counter, 1);
    }
    
    // SYN flood protection for TCP
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)ip + (ip->ihl * 4);
        if ((void *)(tcp + 1) > data_end)
            return XDP_PASS;
        
        // Drop suspicious SYN packets
        if (tcp->syn && !tcp->ack) {
            // Track SYN packets per source
            struct conn_key conn = {
                .src_ip = ip->saddr,
                .dst_ip = ip->daddr,
                .src_port = tcp->source,
                .dst_port = tcp->dest,
                .protocol = IPPROTO_TCP
            };
            
            struct conn_state *state = bpf_map_lookup_elem(&connections, &conn);
            if (state) {
                state->packets++;
                // If too many SYNs from same source, drop
                if (state->packets > 10) {
                    // Add to blacklist
                    __u64 initial_count = 1;
                    bpf_map_update_elem(&blacklist, &src_ip, &initial_count, BPF_ANY);
                    return XDP_DROP;
                }
            } else {
                // New connection
                struct conn_state new_state = {
                    .packets = 1,
                    .bytes = ctx->data_end - ctx->data,
                    .last_seen = bpf_ktime_get_ns()
                };
                bpf_map_update_elem(&connections, &conn, &new_state, BPF_ANY);
            }
        }
    }
    
    // UDP flood protection
    if (ip->protocol == IPPROTO_UDP) {
        struct udphdr *udp = (void *)ip + (ip->ihl * 4);
        if ((void *)(udp + 1) > data_end)
            return XDP_PASS;
        
        // Drop UDP packets to specific ports during attack
        if (bpf_ntohs(udp->dest) == 53 || bpf_ntohs(udp->dest) == 123) {
            // DNS and NTP amplification protection
            if (bpf_ntohs(udp->source) < 1024) {
                return XDP_DROP;  // Suspicious source port
            }
        }
    }
    
    return XDP_PASS;  // Pass legitimate traffic
}

char _license[] SEC("license") = "GPL";
```

### XDP Load Balancer

```c
// xdp_load_balancer.c
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define MAX_BACKENDS 10

struct backend {
    __u32 ip;
    __u8  mac[ETH_ALEN];
    __u32 weight;
    __u32 connections;
};

// Backend servers configuration
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, MAX_BACKENDS);
    __type(key, __u32);
    __type(value, struct backend);
} backends SEC(".maps");

// Active backend count
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u32);
} backend_count SEC(".maps");

// Connection persistence (session affinity)
struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 100000);
    __type(key, __u32);   // Client IP
    __type(value, __u32); // Backend index
} sessions SEC(".maps");

static __always_inline __u32 jhash(const void *key, __u32 length, __u32 initval) {
    // Jenkins hash implementation
    __u32 a, b, c;
    const __u8 *k = key;
    
    a = b = c = 0xdeadbeef + length + initval;
    
    while (length > 12) {
        a += *(__u32 *)k;
        b += *(__u32 *)(k + 4);
        c += *(__u32 *)(k + 8);
        
        // Mix
        a -= c; a ^= ((c << 4) | (c >> 28)); c += b;
        b -= a; b ^= ((a << 6) | (a >> 26)); a += c;
        c -= b; c ^= ((b << 8) | (b >> 24)); b += a;
        
        length -= 12;
        k += 12;
    }
    
    return c;
}

SEC("xdp")
int xdp_lb(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return XDP_PASS;
    
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_PASS;
    
    // Check for VIP destination
    if (ip->daddr != 0x0a000001)  // 10.0.0.1 - VIP
        return XDP_PASS;
    
    __u32 client_ip = ip->saddr;
    __u32 *backend_idx;
    __u32 idx;
    
    // Check session affinity
    backend_idx = bpf_map_lookup_elem(&sessions, &client_ip);
    if (backend_idx) {
        idx = *backend_idx;
    } else {
        // New connection - select backend
        __u32 key = 0;
        __u32 *count = bpf_map_lookup_elem(&backend_count, &key);
        if (!count || *count == 0)
            return XDP_DROP;
        
        // Hash-based selection
        __u32 hash = jhash(&client_ip, sizeof(client_ip), 0);
        idx = hash % *count;
        
        // Store session
        bpf_map_update_elem(&sessions, &client_ip, &idx, BPF_ANY);
    }
    
    // Get backend
    struct backend *backend = bpf_map_lookup_elem(&backends, &idx);
    if (!backend)
        return XDP_DROP;
    
    // Rewrite packet
    ip->daddr = backend->ip;
    
    // Update Ethernet header
    __builtin_memcpy(eth->h_dest, backend->mac, ETH_ALEN);
    // Set source MAC (assuming we know it)
    __u8 src_mac[ETH_ALEN] = {0x00, 0x11, 0x22, 0x33, 0x44, 0x55};
    __builtin_memcpy(eth->h_source, src_mac, ETH_ALEN);
    
    // Recalculate IP checksum
    ip->check = 0;
    __u32 csum = 0;
    __u16 *p = (__u16 *)ip;
    for (int i = 0; i < sizeof(*ip) / 2; i++) {
        csum += *p++;
    }
    while (csum >> 16)
        csum = (csum & 0xffff) + (csum >> 16);
    ip->check = ~csum;
    
    // Redirect to backend interface
    return XDP_TX;  // Or XDP_REDIRECT for different interface
}

char _license[] SEC("license") = "GPL";
```

### XDP Performance Tuning

```bash
#!/bin/bash
# XDP Performance Optimization Script

# Enable XDP native mode (best performance)
ethtool -K eth0 xdp on

# Check XDP support
ethtool -i eth0 | grep driver

# Set ring buffer sizes
ethtool -G eth0 rx 4096 tx 4096

# Enable hardware offload if supported
ethtool -K eth0 hw-tc-offload on

# CPU affinity for NIC interrupts
for irq in $(grep eth0 /proc/interrupts | awk '{print $1}' | sed 's/://'); do
    echo 2 > /proc/irq/$irq/smp_affinity
done

# Load XDP program with native mode
ip link set dev eth0 xdpdrv obj xdp_prog.o sec xdp

# Monitor XDP statistics
watch -n1 'bpftool prog show id $(ip link show eth0 | grep -o "prog/xdp id [0-9]*" | awk "{print \$3}")'
```

---

## 4. Network Programming with eBPF {#network-programming}

### TC (Traffic Control) eBPF

```c
// tc_bandwidth_limiter.c
#include <linux/bpf.h>
#include <linux/pkt_cls.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <bpf/bpf_helpers.h>

#define RATE_LIMIT_BPS 10000000  // 10 Mbps

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, __u32);   // IP address
    __type(value, __u64); // Bytes transferred
} bandwidth_map SEC(".maps");

SEC("tc")
int tc_bandwidth_limit(struct __sk_buff *skb) {
    void *data_end = (void *)(long)skb->data_end;
    void *data = (void *)(long)skb->data;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return TC_ACT_OK;
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return TC_ACT_OK;
    
    struct iphdr *ip = data + sizeof(*eth);
    if ((void *)(ip + 1) > data_end)
        return TC_ACT_OK;
    
    __u32 src_ip = ip->saddr;
    __u64 *bytes = bpf_map_lookup_elem(&bandwidth_map, &src_ip);
    
    if (bytes) {
        *bytes += skb->len;
        if (*bytes > RATE_LIMIT_BPS) {
            return TC_ACT_SHOT;  // Drop packet
        }
    } else {
        __u64 initial_bytes = skb->len;
        bpf_map_update_elem(&bandwidth_map, &src_ip, &initial_bytes, BPF_ANY);
    }
    
    return TC_ACT_OK;
}

char _license[] SEC("license") = "GPL";
```

### Socket eBPF Programs

```c
// socket_filter.c
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>

SEC("socket")
int socket_filter_prog(struct __sk_buff *skb) {
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return 0;  // Drop
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return 0;
    
    struct iphdr *ip = data + sizeof(*eth);
    if ((void *)(ip + 1) > data_end)
        return 0;
    
    // Only accept TCP traffic
    if (ip->protocol != IPPROTO_TCP)
        return 0;
    
    struct tcphdr *tcp = (void *)ip + ip->ihl * 4;
    if ((void *)(tcp + 1) > data_end)
        return 0;
    
    // Only accept HTTP/HTTPS
    __u16 dest_port = bpf_ntohs(tcp->dest);
    if (dest_port != 80 && dest_port != 443)
        return 0;
    
    return skb->len;  // Accept
}

char _license[] SEC("license") = "GPL";
```

### Sockmap for Service Mesh

```c
// sockmap_service_mesh.c
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct {
    __uint(type, BPF_MAP_TYPE_SOCKMAP);
    __uint(max_entries, 1024);
    __type(key, __u32);
    __type(value, __u64);
} sock_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, __u32);   // Service ID
    __type(value, __u32); // Backend sock index
} service_backend_map SEC(".maps");

SEC("sk_msg")
int sk_msg_prog(struct sk_msg_md *msg) {
    __u32 service_id = msg->local_port;  // Use port as service ID
    
    // Lookup backend for service
    __u32 *backend_idx = bpf_map_lookup_elem(&service_backend_map, &service_id);
    if (!backend_idx)
        return SK_PASS;
    
    // Redirect to backend socket
    return bpf_msg_redirect_map(msg, &sock_map, *backend_idx, 0);
}

SEC("sockops")
int sockops_prog(struct bpf_sock_ops *ops) {
    __u32 op = ops->op;
    
    switch (op) {
        case BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB:
        case BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB: {
            // Add socket to map
            __u32 key = ops->local_port;
            __u64 value = 0;
            bpf_sock_map_update(ops, &sock_map, &key, BPF_ANY);
            break;
        }
    }
    
    return 0;
}

char _license[] SEC("license") = "GPL";
```

---

## 5. Cilium & Kubernetes Networking {#cilium}

### Cilium Architecture

```
┌────────────────────────────────────────────────┐
│              Kubernetes API Server              │
└────────────────────────────────────────────────┘
                        ↑
┌────────────────────────────────────────────────┐
│              Cilium Operator                    │
│  • CRD Management                               │
│  • IP Address Management (IPAM)                 │
│  • Node Discovery                               │
└────────────────────────────────────────────────┘
                        ↑
┌────────────────────────────────────────────────┐
│              Cilium Agent (per node)            │
│  ┌──────────────────────────────────────────┐ │
│  │         Policy Repository                 │ │
│  │  • NetworkPolicy                         │ │
│  │  • CiliumNetworkPolicy                   │ │
│  └──────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────┐ │
│  │         eBPF Compiler                    │ │
│  │  • Policy → eBPF                         │ │
│  │  • Load balancing rules                  │ │
│  └──────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────┐ │
│  │         Datapath (eBPF)                  │ │
│  │  • XDP Programs                          │ │
│  │  • TC Programs                           │ │
│  │  • Socket Programs                       │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────┐
│              Linux Kernel                       │
│  • eBPF Maps (state)                           │
│  • Network Stack                               │
└────────────────────────────────────────────────┘
```

### Cilium Network Policy Implementation

```yaml
# cilium_network_policy.yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: api-server-policy
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: api-server
      
  # Ingress Rules
  ingress:
    # Allow from frontend pods
    - fromEndpoints:
        - matchLabels:
            app: frontend
      toPorts:
        - ports:
            - port: "8080"
              protocol: TCP
          rules:
            http:
              - method: GET
                path: "/api/v1/.*"
              - method: POST
                path: "/api/v1/users"
                
    # Allow from specific external IPs
    - fromCIDR:
        - 203.0.113.0/24
      toPorts:
        - ports:
            - port: "443"
              protocol: TCP
              
    # Allow from pods with specific identity
    - fromIdentities:
        - cluster:production
        - k8s:io.kubernetes.pod.namespace=monitoring
        
  # Egress Rules
  egress:
    # Allow DNS
    - toEndpoints:
        - matchLabels:
            k8s:io.kubernetes.pod.namespace: kube-system
            k8s-app: kube-dns
      toPorts:
        - ports:
            - port: "53"
              protocol: UDP
              
    # Allow to database
    - toEndpoints:
        - matchLabels:
            app: postgres
      toPorts:
        - ports:
            - port: "5432"
              protocol: TCP
              
    # Allow to external services with L7 filtering
    - toFQDNs:
        - matchName: "api.external-service.com"
      toPorts:
        - ports:
            - port: "443"
              protocol: TCP
          rules:
            http:
              - method: GET
                path: "/oauth/token"
                headers:
                  - 'Authorization: Bearer .*'
```

### Cilium eBPF Program Example

```c
// cilium_policy_enforcement.c
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>

#define CILIUM_MAP_POLICY 1
#define CILIUM_MAP_IPCACHE 2
#define CILIUM_MAP_ENDPOINTS 3

// Identity for security context
struct identity {
    __u32 id;
    __u32 labels[16];
};

// Policy entry
struct policy_entry {
    __u32 action;  // ALLOW, DENY, REDIRECT
    __u32 proxy_port;  // For L7 proxy redirect
};

// Maps
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 65536);
    __type(key, __u64);  // src_identity << 32 | dst_identity
    __type(value, struct policy_entry);
} policy_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_LPM_TRIE);
    __uint(max_entries, 10000);
    __uint(map_flags, BPF_F_NO_PREALLOC);
    __type(key, struct {
        __u32 prefixlen;
        __u32 ip;
    });
    __type(value, struct identity);
} ipcache_map SEC(".maps");

static __always_inline int policy_can_access(struct __sk_buff *skb,
                                            __u32 src_identity,
                                            __u32 dst_identity) {
    __u64 key = ((__u64)src_identity << 32) | dst_identity;
    struct policy_entry *entry = bpf_map_lookup_elem(&policy_map, &key);
    
    if (!entry) {
        // No explicit policy, check default
        return 0;  // Default deny
    }
    
    if (entry->proxy_port > 0) {
        // Redirect to L7 proxy
        bpf_skb_change_type(skb, PACKET_HOST);
        return bpf_redirect(entry->proxy_port, 0);
    }
    
    return entry->action;  // ALLOW or DENY
}

SEC("tc/ingress")
int cilium_ingress(struct __sk_buff *skb) {
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return TC_ACT_OK;
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return TC_ACT_OK;
    
    struct iphdr *ip = data + sizeof(*eth);
    if ((void *)(ip + 1) > data_end)
        return TC_ACT_OK;
    
    // Lookup source identity
    struct {
        __u32 prefixlen;
        __u32 ip;
    } src_key = {
        .prefixlen = 32,
        .ip = ip->saddr
    };
    
    struct identity *src_id = bpf_map_lookup_elem(&ipcache_map, &src_key);
    if (!src_id)
        return TC_ACT_SHOT;  // Unknown source
    
    // Lookup destination identity
    struct {
        __u32 prefixlen;
        __u32 ip;
    } dst_key = {
        .prefixlen = 32,
        .ip = ip->daddr
    };
    
    struct identity *dst_id = bpf_map_lookup_elem(&ipcache_map, &dst_key);
    if (!dst_id)
        return TC_ACT_SHOT;  // Unknown destination
    
    // Check policy
    int decision = policy_can_access(skb, src_id->id, dst_id->id);
    
    if (decision == 0)
        return TC_ACT_SHOT;  // Deny
    
    // Add identity to packet metadata for egress processing
    skb->mark = src_id->id;
    
    return TC_ACT_OK;  // Allow
}

char _license[] SEC("license") = "GPL";
```

### Cilium BGP Integration

```yaml
# cilium_bgp_config.yaml
apiVersion: cilium.io/v2alpha1
kind: CiliumBGPPeeringPolicy
metadata:
  name: datacenter-bgp-peering
spec:
  # Apply to nodes with BGP label
  nodeSelector:
    matchLabels:
      bgp-enabled: "true"
      
  virtualRouters:
    - localASN: 65001
      exportPodCIDR: true
      serviceSelector:
        matchLabels:
          advertise: "bgp"
          
      neighbors:
        # Upstream router 1
        - peerAddress: "10.0.0.1/32"
          peerASN: 65000
          peerPort: 179
          connectRetryTimeSeconds: 30
          holdTimeSeconds: 90
          keepAliveTimeSeconds: 30
          gracefulRestart:
            enabled: true
            restartTimeSeconds: 120
            
        # Upstream router 2
        - peerAddress: "10.0.0.2/32"
          peerASN: 65000
          peerPort: 179
          connectRetryTimeSeconds: 30
          holdTimeSeconds: 90
          keepAliveTimeSeconds: 30
          
      # BGP Communities
      communities:
        standard:
          - "65001:100"  # Primary path
          - "65001:200"  # Backup path
        large:
          - "65001:1:100"  # Geographic region
          
      # Route policies
      routePolicies:
        - name: "prefer-local"
          type: "export"
          match:
            prefixLength:
              min: 24
              max: 32
          actions:
            - setLocalPreference: 200
            - addCommunity: "65001:100"
```

---

## 6. Practical Labs & Examples {#labs}

### Lab 1: Build a Complete XDP Firewall

```bash
#!/bin/bash
# Complete XDP Firewall Setup

# Create the XDP program
cat > xdp_firewall.c <<'EOF'
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <linux/icmp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

// Rule structure
struct fw_rule {
    __u32 src_ip;
    __u32 src_mask;
    __u32 dst_ip;
    __u32 dst_mask;
    __u16 src_port_min;
    __u16 src_port_max;
    __u16 dst_port_min;
    __u16 dst_port_max;
    __u8  protocol;
    __u8  action;  // 0=DROP, 1=ACCEPT
};

// Firewall rules map
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1000);
    __type(key, __u32);
    __type(value, struct fw_rule);
} firewall_rules SEC(".maps");

// Statistics map
struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 256);
    __type(key, __u32);
    __type(value, __u64);
} stats SEC(".maps");

#define STATS_TOTAL 0
#define STATS_DROPPED 1
#define STATS_ACCEPTED 2

static __always_inline void update_stats(__u32 key) {
    __u64 *value = bpf_map_lookup_elem(&stats, &key);
    if (value)
        __sync_fetch_and_add(value, 1);
}

SEC("xdp")
int xdp_firewall_prog(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    update_stats(STATS_TOTAL);
    
    // Parse Ethernet
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return XDP_PASS;
    
    // Parse IP
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_DROP;
    
    __u16 src_port = 0, dst_port = 0;
    
    // Extract ports for TCP/UDP
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)ip + (ip->ihl * 4);
        if ((void *)(tcp + 1) > data_end)
            return XDP_DROP;
        src_port = bpf_ntohs(tcp->source);
        dst_port = bpf_ntohs(tcp->dest);
    } else if (ip->protocol == IPPROTO_UDP) {
        struct udphdr *udp = (void *)ip + (ip->ihl * 4);
        if ((void *)(udp + 1) > data_end)
            return XDP_DROP;
        src_port = bpf_ntohs(udp->source);
        dst_port = bpf_ntohs(udp->dest);
    }
    
    // Check firewall rules
    for (__u32 i = 0; i < 1000; i++) {
        struct fw_rule *rule = bpf_map_lookup_elem(&firewall_rules, &i);
        if (!rule)
            break;  // No more rules
        
        // Check protocol
        if (rule->protocol != 0 && rule->protocol != ip->protocol)
            continue;
        
        // Check source IP
        if (rule->src_mask != 0) {
            if ((ip->saddr & rule->src_mask) != (rule->src_ip & rule->src_mask))
                continue;
        }
        
        // Check destination IP
        if (rule->dst_mask != 0) {
            if ((ip->daddr & rule->dst_mask) != (rule->dst_ip & rule->dst_mask))
                continue;
        }
        
        // Check ports
        if (src_port < rule->src_port_min || src_port > rule->src_port_max)
            continue;
        if (dst_port < rule->dst_port_min || dst_port > rule->dst_port_max)
            continue;
        
        // Rule matches - take action
        if (rule->action == 0) {
            update_stats(STATS_DROPPED);
            return XDP_DROP;
        } else {
            update_stats(STATS_ACCEPTED);
            return XDP_PASS;
        }
    }
    
    // Default action
    update_stats(STATS_DROPPED);
    return XDP_DROP;  // Default deny
}

char _license[] SEC("license") = "GPL";
EOF

# Compile
clang -O2 -g -target bpf -c xdp_firewall.c -o xdp_firewall.o

# Create rule management tool
cat > manage_firewall.py <<'EOF'
#!/usr/bin/env python3
import sys
import struct
import socket
from bcc import BPF

def ip_to_int(ip):
    return struct.unpack("!I", socket.inet_aton(ip))[0]

def add_rule(bpf, idx, src_ip="0.0.0.0", src_mask="0.0.0.0",
             dst_ip="0.0.0.0", dst_mask="0.0.0.0",
             src_port_min=0, src_port_max=65535,
             dst_port_min=0, dst_port_max=65535,
             protocol=0, action=1):
    
    rule = (ip_to_int(src_ip), ip_to_int(src_mask),
            ip_to_int(dst_ip), ip_to_int(dst_mask),
            src_port_min, src_port_max,
            dst_port_min, dst_port_max,
            protocol, action)
    
    rules_map = bpf.get_table("firewall_rules")
    rules_map[idx] = rule

# Load BPF program
bpf = BPF(src_file="xdp_firewall.c")
fn = bpf.load_func("xdp_firewall_prog", BPF.XDP)

# Add rules
add_rule(bpf, 0, dst_port_min=22, dst_port_max=22, protocol=6, action=1)  # Allow SSH
add_rule(bpf, 1, dst_port_min=80, dst_port_max=80, protocol=6, action=1)  # Allow HTTP
add_rule(bpf, 2, dst_port_min=443, dst_port_max=443, protocol=6, action=1) # Allow HTTPS
add_rule(bpf, 3, protocol=1, action=1)  # Allow ICMP

# Attach to interface
from pyroute2 import IPRoute
ipr = IPRoute()
idx = ipr.link_lookup(ifname="eth0")[0]
ipr.link("set", index=idx, xdp_fd=fn.fd)

print("XDP Firewall loaded. Press Ctrl+C to remove.")
try:
    while True:
        stats = bpf.get_table("stats")
        print(f"Total: {stats[0].value}, Dropped: {stats[1].value}, Accepted: {stats[2].value}")
        time.sleep(1)
except KeyboardInterrupt:
    ipr.link("set", index=idx, xdp_fd=0)
    print("\nXDP Firewall removed.")
EOF

chmod +x manage_firewall.py
```

### Lab 2: Service Mesh with eBPF

```bash
#!/bin/bash
# eBPF-based Service Mesh Implementation

cat > service_mesh.c <<'EOF'
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

// Service registry
struct service {
    __u32 vip;        // Virtual IP
    __u16 vport;      // Virtual Port
    __u32 backends[10]; // Backend IPs
    __u16 backend_ports[10];
    __u8  backend_count;
    __u8  lb_method;  // 0=RR, 1=Hash, 2=LC
};

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 100);
    __type(key, __u64);  // VIP:VPORT
    __type(value, struct service);
} services SEC(".maps");

// Connection tracking for persistence
struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 100000);
    __type(key, __u64);   // Client IP:Port -> VIP:VPort
    __type(value, __u32); // Backend selection
} connections SEC(".maps");

// Round-robin counter
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 100);
    __type(key, __u32);
    __type(value, __u32);
} rr_counter SEC(".maps");

SEC("sk_msg")
int service_mesh_msg(struct sk_msg_md *msg) {
    __u64 service_key = ((__u64)msg->remote_ip4 << 32) | msg->remote_port;
    
    struct service *svc = bpf_map_lookup_elem(&services, &service_key);
    if (!svc)
        return SK_PASS;
    
    // Check existing connection
    __u64 conn_key = ((__u64)msg->local_ip4 << 32) | msg->local_port;
    __u32 *backend_idx = bpf_map_lookup_elem(&connections, &conn_key);
    
    if (!backend_idx) {
        // New connection - select backend
        __u32 idx = 0;
        
        if (svc->lb_method == 0) {
            // Round-robin
            __u32 svc_idx = service_key & 0xFF;
            __u32 *counter = bpf_map_lookup_elem(&rr_counter, &svc_idx);
            if (counter) {
                idx = *counter % svc->backend_count;
                *counter = (*counter + 1) % svc->backend_count;
            }
        } else if (svc->lb_method == 1) {
            // Hash-based
            __u32 hash = jhash_2words(msg->local_ip4, msg->local_port, 0);
            idx = hash % svc->backend_count;
        }
        
        // Store connection
        bpf_map_update_elem(&connections, &conn_key, &idx, BPF_ANY);
        backend_idx = &idx;
    }
    
    // Redirect to backend
    msg->remote_ip4 = svc->backends[*backend_idx];
    msg->remote_port = svc->backend_ports[*backend_idx];
    
    return SK_PASS;
}

SEC("cgroup/sock_ops")
int service_mesh_sockops(struct bpf_sock_ops *ops) {
    switch (ops->op) {
        case BPF_SOCK_OPS_TCP_CONNECT_CB: {
            // Intercept outgoing connections
            __u64 service_key = ((__u64)ops->remote_ip4 << 32) | ops->remote_port;
            struct service *svc = bpf_map_lookup_elem(&services, &service_key);
            if (svc) {
                // This is a service connection - enable msg program
                bpf_sock_ops_cb_flags_set(ops, BPF_SOCK_OPS_STATE_CB_FLAG);
            }
            break;
        }
    }
    return 0;
}

char _license[] SEC("license") = "GPL";
EOF

# Service mesh control plane
cat > service_mesh_control.py <<'EOF'
#!/usr/bin/env python3
import json
import socket
import struct
from flask import Flask, request, jsonify

app = Flask(__name__)

def load_bpf():
    # Load eBPF programs
    # This would use libbpf or bcc in production
    pass

@app.route('/services', methods=['POST'])
def register_service():
    data = request.json
    vip = socket.inet_aton(data['vip'])
    vport = data['vport']
    backends = [socket.inet_aton(ip) for ip in data['backends']]
    
    # Update eBPF map
    # Implementation would update the services map
    
    return jsonify({"status": "registered"})

@app.route('/metrics')
def get_metrics():
    # Read from eBPF maps
    metrics = {
        "connections": 0,
        "requests_per_second": 0,
        "latency_p99": 0
    }
    return jsonify(metrics)

if __name__ == '__main__':
    load_bpf()
    app.run(host='0.0.0.0', port=8080)
EOF
```

### Lab 3: Network Observability with eBPF

```c
// network_observability.c
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

// Flow record
struct flow_record {
    __u32 src_ip;
    __u32 dst_ip;
    __u16 src_port;
    __u16 dst_port;
    __u8  protocol;
    __u8  tcp_flags;
    __u64 packets;
    __u64 bytes;
    __u64 start_time;
    __u64 last_time;
    __u32 latency_sum;
    __u32 latency_count;
};

// Flow tracking
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 100000);
    __type(key, struct flow_record);
    __type(value, struct flow_record);
} flows SEC(".maps");

// TCP RTT tracking
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __type(key, __u64);  // Socket cookie
    __type(value, __u32); // RTT in microseconds
} tcp_rtt SEC(".maps");

// Per-CPU metrics
struct metrics {
    __u64 total_packets;
    __u64 total_bytes;
    __u64 tcp_retransmits;
    __u64 tcp_resets;
    __u64 udp_errors;
};

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, struct metrics);
} metrics SEC(".maps");

SEC("tracepoint/tcp/tcp_retransmit_skb")
int trace_tcp_retransmit(struct trace_event_raw_tcp_event_sk_skb *ctx) {
    __u32 key = 0;
    struct metrics *m = bpf_map_lookup_elem(&metrics, &key);
    if (m)
        __sync_fetch_and_add(&m->tcp_retransmits, 1);
    return 0;
}

SEC("kprobe/tcp_rcv_established")
int kprobe_tcp_rcv(struct pt_regs *ctx) {
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    
    // Read socket cookie for unique identification
    __u64 cookie = bpf_get_socket_cookie(sk);
    
    // Get TCP info
    struct tcp_sock *tp = (struct tcp_sock *)sk;
    __u32 srtt = BPF_CORE_READ(tp, srtt_us) >> 3;  // Smoothed RTT
    
    bpf_map_update_elem(&tcp_rtt, &cookie, &srtt, BPF_ANY);
    return 0;
}

SEC("xdp")
int xdp_flow_monitor(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return XDP_PASS;
    
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_PASS;
    
    struct flow_record flow = {
        .src_ip = ip->saddr,
        .dst_ip = ip->daddr,
        .protocol = ip->protocol,
    };
    
    // Extract ports and flags
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)ip + (ip->ihl * 4);
        if ((void *)(tcp + 1) > data_end)
            return XDP_PASS;
        
        flow.src_port = bpf_ntohs(tcp->source);
        flow.dst_port = bpf_ntohs(tcp->dest);
        flow.tcp_flags = ((tcp->fin) | (tcp->syn << 1) | (tcp->rst << 2) |
                         (tcp->psh << 3) | (tcp->ack << 4) | (tcp->urg << 5));
        
        // Track RST packets
        if (tcp->rst) {
            __u32 key = 0;
            struct metrics *m = bpf_map_lookup_elem(&metrics, &key);
            if (m)
                __sync_fetch_and_add(&m->tcp_resets, 1);
        }
    } else if (ip->protocol == IPPROTO_UDP) {
        struct udphdr *udp = (void *)ip + (ip->ihl * 4);
        if ((void *)(udp + 1) > data_end)
            return XDP_PASS;
        
        flow.src_port = bpf_ntohs(udp->source);
        flow.dst_port = bpf_ntohs(udp->dest);
    }
    
    // Update flow record
    struct flow_record *existing = bpf_map_lookup_elem(&flows, &flow);
    if (existing) {
        __sync_fetch_and_add(&existing->packets, 1);
        __sync_fetch_and_add(&existing->bytes, ctx->data_end - ctx->data);
        existing->last_time = bpf_ktime_get_ns();
    } else {
        flow.packets = 1;
        flow.bytes = ctx->data_end - ctx->data;
        flow.start_time = bpf_ktime_get_ns();
        flow.last_time = flow.start_time;
        bpf_map_update_elem(&flows, &flow, &flow, BPF_ANY);
    }
    
    // Update global metrics
    __u32 key = 0;
    struct metrics *m = bpf_map_lookup_elem(&metrics, &key);
    if (m) {
        __sync_fetch_and_add(&m->total_packets, 1);
        __sync_fetch_and_add(&m->total_bytes, ctx->data_end - ctx->data);
    }
    
    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
```

---

## 7. Performance & Observability {#performance}

### eBPF Performance Metrics

```python
#!/usr/bin/env python3
# ebpf_performance_monitor.py
import time
import psutil
from bcc import BPF
from prometheus_client import Counter, Gauge, Histogram, start_http_server

# Prometheus metrics
packet_counter = Counter('ebpf_packets_total', 'Total packets processed')
byte_counter = Counter('ebpf_bytes_total', 'Total bytes processed')
drop_counter = Counter('ebpf_drops_total', 'Total packets dropped')
cpu_gauge = Gauge('ebpf_cpu_usage', 'CPU usage by eBPF programs')
latency_histogram = Histogram('ebpf_processing_latency_us', 'Processing latency')

# BPF program for monitoring
bpf_text = """
#include <linux/bpf.h>
#include <linux/ptrace.h>

BPF_PERF_OUTPUT(events);
BPF_HISTOGRAM(latency, u64);

struct data_t {
    u64 timestamp;
    u32 cpu;
    u64 instruction_count;
    u64 cycle_count;
};

TRACEPOINT_PROBE(xdp, xdp_cpumap_kthread) {
    struct data_t data = {};
    data.timestamp = bpf_ktime_get_ns();
    data.cpu = bpf_get_smp_processor_id();
    
    // Get performance counters
    data.instruction_count = bpf_perf_prog_read_value(ctx, BPF_PERF_EVENT_INSTRUCTION);
    data.cycle_count = bpf_perf_prog_read_value(ctx, BPF_PERF_EVENT_CYCLES);
    
    events.perf_submit(ctx, &data, sizeof(data));
    return 0;
}
"""

def monitor_ebpf_performance():
    b = BPF(text=bpf_text)
    
    def process_event(cpu, data, size):
        event = b["events"].event(data)
        
        # Update Prometheus metrics
        cpu_gauge.set(psutil.cpu_percent(interval=0.1))
        
        # Calculate latency
        latency = (event.cycle_count / 2400000)  # Assuming 2.4GHz CPU
        latency_histogram.observe(latency)
    
    b["events"].open_perf_buffer(process_event)
    
    # Start Prometheus metrics server
    start_http_server(8000)
    
    print("Monitoring eBPF performance. Metrics available at http://localhost:8000")
    while True:
        try:
            b.perf_buffer_poll()
        except KeyboardInterrupt:
            break

if __name__ == "__main__":
    monitor_ebpf_performance()
```

### eBPF Optimization Techniques

```c
// optimized_ebpf.c - Performance optimization examples

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

// 1. Use per-CPU maps for better performance
struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 256);
    __type(key, __u32);
    __type(value, __u64);
} percpu_stats SEC(".maps");

// 2. Batch operations
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __uint(map_flags, BPF_F_NO_PREALLOC);  // Dynamic allocation
    __type(key, __u32);
    __type(value, __u64);
} batch_map SEC(".maps");

// 3. Use static inline for better performance
static __always_inline int process_packet_fast(void *data, void *data_end) {
    // Unroll loops when possible
    #pragma unroll
    for (int i = 0; i < 4; i++) {
        // Process 4 items at once
    }
    return 0;
}

// 4. Minimize map lookups
SEC("xdp")
int optimized_xdp(struct xdp_md *ctx) {
    // Cache map lookups
    __u32 key = 0;
    __u64 *value = bpf_map_lookup_elem(&percpu_stats, &key);
    if (!value)
        return XDP_PASS;
    
    // Use value multiple times without additional lookups
    *value += 1;
    
    // Use tail calls for complex logic
    bpf_tail_call(ctx, &prog_array, 0);
    
    return XDP_PASS;
}

// 5. Use BPF_MAP_TYPE_ARRAY_OF_MAPS for dynamic programming
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY_OF_MAPS);
    __uint(max_entries, 10);
    __type(key, __u32);
    __uint(value_id, 1);  // ID of inner map
} map_of_maps SEC(".maps");

char _license[] SEC("license") = "GPL";
```

---

## 8. eBPF for Network Automation {#automation}

### Automated Network Policy with eBPF

```python
#!/usr/bin/env python3
# automated_network_policy.py

import yaml
import ipaddress
from bcc import BPF
from kubernetes import client, config

class eBPFNetworkPolicyController:
    def __init__(self):
        self.bpf = None
        self.load_bpf_program()
        config.load_incluster_config()
        self.k8s_v1 = client.CoreV1Api()
        self.k8s_networking = client.NetworkingV1Api()
    
    def load_bpf_program(self):
        with open('network_policy.c', 'r') as f:
            bpf_text = f.read()
        self.bpf = BPF(text=bpf_text)
    
    def watch_network_policies(self):
        w = watch.Watch()
        for event in w.stream(self.k8s_networking.list_network_policy_for_all_namespaces):
            policy = event['object']
            event_type = event['type']
            
            if event_type == 'ADDED' or event_type == 'MODIFIED':
                self.apply_policy(policy)
            elif event_type == 'DELETED':
                self.remove_policy(policy)
    
    def apply_policy(self, policy):
        """Convert Kubernetes NetworkPolicy to eBPF rules"""
        
        # Get pod selector
        selector = policy.spec.pod_selector
        pods = self.get_pods_by_selector(selector, policy.metadata.namespace)
        
        # Process ingress rules
        if policy.spec.ingress:
            for rule in policy.spec.ingress:
                self.process_ingress_rule(pods, rule)
        
        # Process egress rules
        if policy.spec.egress:
            for rule in policy.spec.egress:
                self.process_egress_rule(pods, rule)
    
    def process_ingress_rule(self, pods, rule):
        """Convert ingress rule to eBPF map entry"""
        
        for pod in pods:
            pod_ip = self.get_pod_ip(pod)
            
            # Process 'from' selectors
            if rule.from_:
                for from_rule in rule.from_:
                    if from_rule.pod_selector:
                        source_pods = self.get_pods_by_selector(
                            from_rule.pod_selector,
                            pod.metadata.namespace
                        )
                        for source_pod in source_pods:
                            source_ip = self.get_pod_ip(source_pod)
                            self.add_ebpf_rule(source_ip, pod_ip, 'ALLOW')
            
            # Process ports
            if rule.ports:
                for port in rule.ports:
                    self.add_port_rule(pod_ip, port.port, port.protocol)
    
    def add_ebpf_rule(self, src_ip, dst_ip, action):
        """Add rule to eBPF map"""
        
        src_int = int(ipaddress.ip_address(src_ip))
        dst_int = int(ipaddress.ip_address(dst_ip))
        
        key = self.bpf["policy_map"].Key()
        key.src_ip = src_int
        key.dst_ip = dst_int
        
        value = self.bpf["policy_map"].Value()
        value.action = 1 if action == 'ALLOW' else 0
        
        self.bpf["policy_map"][key] = value
    
    def attach_to_interfaces(self):
        """Attach eBPF programs to network interfaces"""
        
        # Get all node interfaces
        interfaces = self.get_node_interfaces()
        
        for iface in interfaces:
            # Attach XDP program
            self.bpf.attach_xdp(iface, self.bpf.load_func("xdp_policy", BPF.XDP))
            
            # Attach TC program
            self.bpf.attach_tc(iface, self.bpf.load_func("tc_policy", BPF.SCHED_CLS))

def main():
    controller = eBPFNetworkPolicyController()
    controller.attach_to_interfaces()
    controller.watch_network_policies()

if __name__ == "__main__":
    main()
```

### BGP Route Injection with eBPF

```go
// bgp_ebpf_controller.go
package main

import (
    "context"
    "encoding/binary"
    "net"
    "github.com/cilium/ebpf"
    "github.com/cilium/ebpf/link"
    "github.com/cilium/ebpf/rlimit"
    "github.com/osrg/gobgp/v3/pkg/server"
    "github.com/osrg/gobgp/v3/pkg/packet/bgp"
)

type BGPeBPFController struct {
    bgpServer *server.BgpServer
    routeMap  *ebpf.Map
    xdpProg   *ebpf.Program
}

func NewBGPeBPFController() (*BGPeBPFController, error) {
    // Remove memory limit for eBPF
    if err := rlimit.RemoveMemlock(); err != nil {
        return nil, err
    }
    
    // Load eBPF program
    spec, err := ebpf.LoadCollectionSpec("bgp_routes.o")
    if err != nil {
        return nil, err
    }
    
    coll, err := ebpf.NewCollection(spec)
    if err != nil {
        return nil, err
    }
    
    // Initialize BGP server
    bgpServer := server.NewBgpServer()
    go bgpServer.Serve()
    
    return &BGPeBPFController{
        bgpServer: bgpServer,
        routeMap:  coll.Maps["bgp_routes"],
        xdpProg:   coll.Programs["xdp_bgp_filter"],
    }, nil
}

func (c *BGPeBPFController) WatchBGPRoutes(ctx context.Context) {
    // Monitor BGP RIB changes
    w := c.bgpServer.Watch(server.WatchUpdate(true))
    
    for {
        select {
        case ev := <-w.Event():
            switch msg := ev.Message.(type) {
            case *server.WatchEventUpdate:
                for _, path := range msg.PathList {
                    c.updateeBPFRoute(path)
                }
            }
        case <-ctx.Done():
            return
        }
    }
}

func (c *BGPeBPFController) updateeBPFRoute(path *server.Path) error {
    nlri := path.GetNlri()
    
    switch nlri.GetType() {
    case bgp.RF_IPv4_UC:
        prefix := nlri.(*bgp.IPAddrPrefix)
        ip := net.ParseIP(prefix.Prefix.String())
        
        // Convert to eBPF map key
        key := binary.BigEndian.Uint32(ip.To4())
        
        // Prepare value based on path attributes
        value := struct {
            NextHop   uint32
            LocalPref uint32
            ASPath    [10]uint32
        }{}
        
        // Extract next hop
        if nh := path.GetNexthop(); nh != nil {
            value.NextHop = binary.BigEndian.Uint32(nh.To4())
        }
        
        // Update eBPF map
        return c.routeMap.Put(key, value)
    }
    
    return nil
}

func (c *BGPeBPFController) AttachXDP(ifname string) error {
    iface, err := net.InterfaceByName(ifname)
    if err != nil {
        return err
    }
    
    // Attach XDP program
    l, err := link.AttachXDP(link.XDPOptions{
        Program:   c.xdpProg,
        Interface: iface.Index,
        Flags:     link.XDPDriverMode,
    })
    if err != nil {
        return err
    }
    defer l.Close()
    
    return nil
}

func main() {
    controller, err := NewBGPeBPFController()
    if err != nil {
        panic(err)
    }
    
    // Attach to interface
    if err := controller.AttachXDP("eth0"); err != nil {
        panic(err)
    }
    
    // Start monitoring BGP routes
    ctx := context.Background()
    controller.WatchBGPRoutes(ctx)
}
```

---

## 9. Interview Questions & Answers {#interview}

### Q1: What is eBPF and why is it revolutionary for networking?

**Answer:**
eBPF (extended Berkeley Packet Filter) is a technology that allows running sandboxed programs in the Linux kernel without modifying kernel source code or loading modules. It's revolutionary because:

1. **Performance**: Runs at kernel speed without context switches
2. **Safety**: Verified before execution, cannot crash the kernel
3. **Flexibility**: Can be updated without rebooting
4. **Observability**: Provides deep insights without overhead
5. **Programmability**: Allows custom logic at any point in the stack

In networking, eBPF enables:
- XDP for line-rate packet processing
- Advanced load balancing without proxies
- Dynamic security policies
- Zero-overhead observability
- Service mesh without sidecars

### Q2: Explain XDP and its advantages over traditional packet processing

**Answer:**
XDP (eXpress Data Path) processes packets at the earliest possible point - in the driver, before SKB allocation:

**Traditional Path**:
```
NIC → Driver → SKB allocation → Netfilter → TCP/IP stack → Application
```

**XDP Path**:
```
NIC → Driver → XDP → Decision (DROP/PASS/TX/REDIRECT)
```

**Advantages**:
- **Performance**: 20M+ pps on commodity hardware
- **Low Latency**: No SKB allocation overhead
- **CPU Efficiency**: Processes packets with minimal instructions
- **DDoS Mitigation**: Drop malicious traffic before it consumes resources
- **Flexibility**: Programmable with eBPF

### Q3: How does Cilium use eBPF for Kubernetes networking?

**Answer:**
Cilium replaces kube-proxy and traditional CNI plugins with eBPF programs:

1. **Datapath**: eBPF programs handle all packet forwarding
2. **Load Balancing**: BPF maps store service endpoints
3. **Network Policy**: Compiled to eBPF for line-rate enforcement
4. **Identity-Based Security**: Uses eBPF maps for pod identity
5. **Observability**: eBPF provides flow visibility without overhead

Key components:
- **XDP**: Early packet filtering and DDoS protection
- **TC eBPF**: Pod-to-pod communication
- **Socket eBPF**: Service mesh functionality
- **Sockmap**: Kernel-level socket redirection

### Q4: Design a DDoS mitigation system using XDP

**Answer:**
```c
// Architecture:
// 1. XDP for early filtering
// 2. eBPF maps for state
// 3. Userspace controller for updates

// Key components:

// Rate limiting per IP
struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 1000000);
    __type(key, __u32);  // IP address
    __type(value, struct {
        __u64 packets;
        __u64 last_seen;
    });
} rate_limit SEC(".maps");

// Known good IPs (whitelist)
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10000);
    __type(key, __u32);
    __type(value, __u8);
} whitelist SEC(".maps");

// SYN cookie validation
struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 100000);
    __type(key, __u64);  // Cookie
    __type(value, __u64);  // Timestamp
} syn_cookies SEC(".maps");

// The system would:
// 1. Check whitelist (O(1) lookup)
// 2. Rate limit per IP
// 3. Validate SYN cookies
// 4. Pattern matching for known attacks
// 5. Redirect suspicious traffic to scrubbing
```

### Q5: How would you troubleshoot eBPF program issues?

**Answer:**

1. **Verification Errors**:
```bash
# Check verifier output
bpftool prog load program.o /sys/fs/bpf/prog
# Use -d for debug output
```

2. **Runtime Issues**:
```bash
# List loaded programs
bpftool prog list

# Show program details
bpftool prog show id 42

# Dump program instructions
bpftool prog dump xlated id 42

# Check maps
bpftool map list
bpftool map dump id 10
```

3. **Performance Analysis**:
```bash
# Profile eBPF program
perf record -e bpf:* -a
perf script

# Check statistics
bpftool prog stat
```

4. **Debugging Tools**:
- `bpf_trace_printk()` for debugging (not production)
- `bpf_perf_event_output()` for production logging
- Prometheus metrics from maps

### Q6: Explain eBPF map types and when to use each

**Answer:**

| Map Type | Use Case | Example |
|----------|----------|---------|
| **HASH** | Key-value store | Connection tracking |
| **ARRAY** | Fixed-size indexed | Per-CPU stats |
| **LRU_HASH** | Cache with eviction | Recent connections |
| **PERCPU_HASH** | Per-CPU hash table | High-frequency counters |
| **LPM_TRIE** | Longest prefix match | IP routing tables |
| **PROG_ARRAY** | Tail calls | Program chaining |
| **DEVMAP** | XDP redirect | Load balancing |
| **CPUMAP** | CPU redirect | RSS replacement |
| **SOCKMAP** | Socket redirect | Service mesh |
| **QUEUE/STACK** | FIFO/LIFO | Event buffering |

### Q7: How does eBPF ensure safety?

**Answer:**
The eBPF verifier ensures safety through:

1. **Static Analysis**: Analyzes all possible execution paths
2. **Bounds Checking**: Ensures no out-of-bounds access
3. **Loop Detection**: Prevents infinite loops
4. **Stack Limits**: Maximum 512 bytes stack
5. **Instruction Limit**: Max 1M instructions (configurable)
6. **Helper Validation**: Only approved kernel functions
7. **Type Safety**: Ensures correct data types
8. **Privilege Checks**: CAP_BPF or CAP_SYS_ADMIN required

### Q8: Implement connection tracking in XDP

**Answer:**
```c
struct conn_tuple {
    __u32 src_ip;
    __u32 dst_ip;
    __u16 src_port;
    __u16 dst_port;
    __u8  protocol;
};

struct conn_state {
    __u64 packets_orig;
    __u64 packets_reply;
    __u64 bytes_orig;
    __u64 bytes_reply;
    __u64 last_seen;
    __u8  state;  // NEW, ESTABLISHED, etc.
};

struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, 1000000);
    __type(key, struct conn_tuple);
    __type(value, struct conn_state);
} conntrack SEC(".maps");

static __always_inline int track_connection(struct xdp_md *ctx,
                                           struct conn_tuple *tuple) {
    struct conn_state *state = bpf_map_lookup_elem(&conntrack, tuple);
    
    if (state) {
        state->packets_orig++;
        state->bytes_orig += (ctx->data_end - ctx->data);
        state->last_seen = bpf_ktime_get_ns();
        
        if (state->state == CT_NEW)
            state->state = CT_ESTABLISHED;
    } else {
        struct conn_state new_state = {
            .packets_orig = 1,
            .bytes_orig = ctx->data_end - ctx->data,
            .last_seen = bpf_ktime_get_ns(),
            .state = CT_NEW
        };
        bpf_map_update_elem(&conntrack, tuple, &new_state, BPF_ANY);
    }
    
    return XDP_PASS;
}
```

## Quick Reference

### Essential Commands
```bash
# Load XDP program
ip link set dev eth0 xdp obj prog.o sec xdp

# List eBPF programs
bpftool prog list

# List eBPF maps  
bpftool map list

# Monitor map contents
bpftool map dump id 10

# Attach TC eBPF
tc filter add dev eth0 ingress bpf da obj prog.o sec tc

# Trace eBPF events
bpftool prog trace

# Profile eBPF overhead
perf stat -e bpf:* -a
```

This comprehensive eBPF overview should prepare you well for the networking aspects of your interview, especially regarding Cilium and modern packet processing. The practical examples demonstrate real-world applications that Zoom likely uses in their infrastructure.
