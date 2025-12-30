//! Shared types between eBPF (kernel) and userspace programs.
//!
//! These types are used in eBPF maps and must have matching memory layouts
//! in both the kernel eBPF program and the userspace loader.
//!
//! # Important
//! - All types must be `#[repr(C)]` for consistent memory layout
//! - Types must implement `Clone`, `Copy`, `Default` for eBPF map operations
//! - The `user` feature enables `aya::Pod` implementations for userspace

#![no_std]

/// Connection identifier - used as a key in the connections map.
///
/// Uniquely identifies a TCP connection by its 4-tuple:
/// source IP, destination IP, source port, destination port.
#[repr(C)]
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq, Hash)]
pub struct ConnKey {
    /// Source IP address (network byte order)
    pub src_ip: u32,
    /// Destination IP address (network byte order)
    pub dst_ip: u32,
    /// Source port (host byte order)
    pub src_port: u16,
    /// Destination port (host byte order)
    pub dst_port: u16,
}

#[cfg(feature = "user")]
unsafe impl aya::Pod for ConnKey {}

/// Per-connection metrics stored in eBPF map.
///
/// Updated by kernel eBPF programs on every packet send/receive.
/// Read by userspace for metrics export.
#[repr(C)]
#[derive(Clone, Copy, Debug, Default)]
pub struct ConnMetrics {
    /// Total bytes sent on this connection
    pub bytes_sent: u64,
    /// Total bytes received on this connection
    pub bytes_recv: u64,
    /// Total packets sent
    pub packets_sent: u64,
    /// Total packets received
    pub packets_recv: u64,
    /// Connection start time (nanoseconds since boot)
    pub start_ns: u64,
    /// Last activity time (nanoseconds since boot)
    pub last_seen_ns: u64,
    /// Number of TCP retransmissions (indicates network quality)
    pub retransmits: u32,
    /// Padding for 8-byte alignment
    pub _padding: u32,
}

#[cfg(feature = "user")]
unsafe impl aya::Pod for ConnMetrics {}

/// HTTP request/response event sent via perf buffer.
///
/// Captures HTTP-level metrics for L7 observability.
#[repr(C)]
#[derive(Clone, Copy, Debug, Default)]
pub struct HttpEvent {
    /// Connection this event belongs to
    pub conn: ConnKey,
    /// Request/response latency in nanoseconds
    pub latency_ns: u64,
    /// HTTP status code (e.g., 200, 404, 500)
    pub status_code: u16,
    /// HTTP method: 0=GET, 1=POST, 2=PUT, 3=DELETE, 4=PATCH, 5=HEAD, 6=OPTIONS
    pub method: u8,
    /// Padding for alignment
    pub _padding: u8,
    /// Request path hash (for grouping similar requests)
    pub path_hash: u32,
}

#[cfg(feature = "user")]
unsafe impl aya::Pod for HttpEvent {}

/// Process information for filtering by PID/cgroup.
#[repr(C)]
#[derive(Clone, Copy, Debug, Default)]
pub struct ProcessInfo {
    /// Process ID
    pub pid: u32,
    /// Thread group ID (usually same as PID for main thread)
    pub tgid: u32,
    /// User ID
    pub uid: u32,
    /// Group ID
    pub gid: u32,
    /// cgroup ID (for container-aware filtering)
    pub cgroup_id: u64,
}

#[cfg(feature = "user")]
unsafe impl aya::Pod for ProcessInfo {}

/// Configuration passed from userspace to eBPF.
#[repr(C)]
#[derive(Clone, Copy, Debug, Default)]
pub struct SidecarConfig {
    /// Target PID to monitor (0 = all processes)
    pub target_pid: u32,
    /// Target cgroup ID to monitor (0 = all cgroups)
    pub target_cgroup: u64,
    /// Ports to monitor (0 = all ports, otherwise filter)
    pub target_ports: [u16; 8],
    /// Number of ports in target_ports array
    pub num_target_ports: u8,
    /// Enable HTTP parsing (L7 inspection)
    pub enable_http: u8,
    /// Enable detailed per-packet logging (debug mode)
    pub debug_mode: u8,
    /// Padding
    pub _padding: u8,
}

#[cfg(feature = "user")]
unsafe impl aya::Pod for SidecarConfig {}

/// HTTP method constants
pub mod http_method {
    pub const GET: u8 = 0;
    pub const POST: u8 = 1;
    pub const PUT: u8 = 2;
    pub const DELETE: u8 = 3;
    pub const PATCH: u8 = 4;
    pub const HEAD: u8 = 5;
    pub const OPTIONS: u8 = 6;
    pub const UNKNOWN: u8 = 255;
}
