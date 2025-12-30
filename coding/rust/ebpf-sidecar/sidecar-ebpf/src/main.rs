//! eBPF Sidecar - Kernel Space Program
//!
//! This eBPF program attaches to kernel functions to observe network traffic
//! without modifying the application or adding latency through userspace proxying.
//!
//! # Attach Points
//! - `tcp_connect`: Track new outbound connections
//! - `tcp_sendmsg`: Track bytes sent
//! - `tcp_recvmsg`: Track bytes received  
//! - `tcp_close`: Clean up connection tracking
//! - `tcp_retransmit_skb`: Track retransmissions
//!
//! # Maps
//! - `CONNECTIONS`: Per-connection metrics (HashMap)
//! - `CONFIG`: Runtime configuration (Array)
//! - `EVENTS`: HTTP events perf buffer

#![no_std]
#![no_main]

use aya_ebpf::{
    bindings::BPF_F_NO_PREALLOC,
    helpers::{bpf_get_current_pid_tgid, bpf_ktime_get_ns, bpf_probe_read_kernel},
    macros::{kprobe, kretprobe, map, tracepoint},
    maps::{Array, HashMap, PerfEventArray},
    programs::{ProbeContext, RetProbeContext, TracePointContext},
    EbpfContext,
};
use aya_log_ebpf::{debug, info, warn};
use sidecar_common::{ConnKey, ConnMetrics, HttpEvent, SidecarConfig};

// ============================================================================
// eBPF Maps - Shared data structures between kernel and userspace
// ============================================================================

/// Per-connection metrics storage
/// Key: ConnKey (4-tuple), Value: ConnMetrics
#[map]
static CONNECTIONS: HashMap<ConnKey, ConnMetrics> =
    HashMap::with_max_entries(10240, BPF_F_NO_PREALLOC);

/// Runtime configuration from userspace
/// Index 0 contains the current SidecarConfig
#[map]
static CONFIG: Array<SidecarConfig> = Array::with_max_entries(1, 0);

/// HTTP events sent to userspace via perf buffer
#[map]
static EVENTS: PerfEventArray<HttpEvent> = PerfEventArray::new(0);

// ============================================================================
// Helper Functions
// ============================================================================

/// Check if we should trace this process based on config
#[inline(always)]
fn should_trace(ctx: &impl EbpfContext) -> bool {
    let config = match unsafe { CONFIG.get(0) } {
        Some(c) => c,
        None => return true, // No config = trace everything
    };

    // If target_pid is set, only trace that PID
    if config.target_pid != 0 {
        let pid = (bpf_get_current_pid_tgid() >> 32) as u32;
        if pid != config.target_pid {
            return false;
        }
    }

    true
}

/// Extract connection key from sock struct pointer
/// 
/// # Safety
/// Caller must ensure sock pointer is valid
#[inline(always)]
unsafe fn read_conn_key_from_sock(sock: *const u8) -> Result<ConnKey, i64> {
    // Offsets into struct sock -> __sk_common
    // These are for Linux 5.x+ kernels - may need adjustment
    // In production, use CO-RE (Compile Once Run Everywhere) for portability
    const SK_COMMON_OFFSET: usize = 0;
    const SKADDR_OFFSET: usize = 4;   // __sk_common.skc_rcv_saddr
    const DADDR_OFFSET: usize = 0;    // __sk_common.skc_daddr
    const SPORT_OFFSET: usize = 14;   // __sk_common.skc_num (source port)
    const DPORT_OFFSET: usize = 12;   // __sk_common.skc_dport (dest port, network order)

    let common = sock.add(SK_COMMON_OFFSET);

    let src_ip = bpf_probe_read_kernel(common.add(SKADDR_OFFSET) as *const u32)
        .map_err(|_| 1i64)?;
    let dst_ip = bpf_probe_read_kernel(common.add(DADDR_OFFSET) as *const u32)
        .map_err(|_| 2i64)?;
    let src_port = bpf_probe_read_kernel(common.add(SPORT_OFFSET) as *const u16)
        .map_err(|_| 3i64)?;
    let dst_port_be = bpf_probe_read_kernel(common.add(DPORT_OFFSET) as *const u16)
        .map_err(|_| 4i64)?;

    Ok(ConnKey {
        src_ip,
        dst_ip,
        src_port,
        dst_port: u16::from_be(dst_port_be),
    })
}

// ============================================================================
// Kprobe Programs - Attach to kernel functions
// ============================================================================

/// Track new TCP connections (outbound connect)
#[kprobe]
pub fn trace_tcp_connect(ctx: ProbeContext) -> u32 {
    match try_trace_tcp_connect(&ctx) {
        Ok(()) => 0,
        Err(e) => {
            warn!(&ctx, "tcp_connect error: {}", e);
            1
        }
    }
}

fn try_trace_tcp_connect(ctx: &ProbeContext) -> Result<(), i64> {
    if !should_trace(ctx) {
        return Ok(());
    }

    // First argument is struct sock *
    let sock: *const u8 = ctx.arg(0).ok_or(1i64)?;
    let key = unsafe { read_conn_key_from_sock(sock)? };

    let now = unsafe { bpf_ktime_get_ns() };
    let metrics = ConnMetrics {
        bytes_sent: 0,
        bytes_recv: 0,
        packets_sent: 0,
        packets_recv: 0,
        start_ns: now,
        last_seen_ns: now,
        retransmits: 0,
        _padding: 0,
    };

    CONNECTIONS.insert(&key, &metrics, 0)?;

    debug!(
        ctx,
        "NEW CONN: {}:{} -> {}:{}",
        key.src_ip,
        key.src_port,
        key.dst_ip,
        key.dst_port
    );

    Ok(())
}

/// Track TCP send operations
#[kprobe]
pub fn trace_tcp_sendmsg(ctx: ProbeContext) -> u32 {
    match try_trace_tcp_sendmsg(&ctx) {
        Ok(()) => 0,
        Err(_) => 1,
    }
}

fn try_trace_tcp_sendmsg(ctx: &ProbeContext) -> Result<(), i64> {
    if !should_trace(ctx) {
        return Ok(());
    }

    let sock: *const u8 = ctx.arg(0).ok_or(1i64)?;
    let size: usize = ctx.arg(2).ok_or(2i64)?;

    let key = unsafe { read_conn_key_from_sock(sock)? };

    if let Some(metrics) = unsafe { CONNECTIONS.get_ptr_mut(&key) } {
        let m = unsafe { &mut *metrics };
        m.bytes_sent += size as u64;
        m.packets_sent += 1;
        m.last_seen_ns = unsafe { bpf_ktime_get_ns() };
    }

    Ok(())
}

/// Track TCP receive operations
#[kprobe]
pub fn trace_tcp_recvmsg(ctx: ProbeContext) -> u32 {
    match try_trace_tcp_recvmsg(&ctx) {
        Ok(()) => 0,
        Err(_) => 1,
    }
}

fn try_trace_tcp_recvmsg(ctx: &ProbeContext) -> Result<(), i64> {
    if !should_trace(ctx) {
        return Ok(());
    }

    let sock: *const u8 = ctx.arg(0).ok_or(1i64)?;
    let key = unsafe { read_conn_key_from_sock(sock)? };

    // Note: We increment packet count here, but can't easily get size
    // For accurate byte counts, use kretprobe to capture return value
    if let Some(metrics) = unsafe { CONNECTIONS.get_ptr_mut(&key) } {
        let m = unsafe { &mut *metrics };
        m.packets_recv += 1;
        m.last_seen_ns = unsafe { bpf_ktime_get_ns() };
    }

    Ok(())
}

/// Track TCP receive return to get actual bytes received
#[kretprobe]
pub fn trace_tcp_recvmsg_ret(ctx: RetProbeContext) -> u32 {
    match try_trace_tcp_recvmsg_ret(&ctx) {
        Ok(()) => 0,
        Err(_) => 1,
    }
}

fn try_trace_tcp_recvmsg_ret(ctx: &RetProbeContext) -> Result<(), i64> {
    // Return value is bytes received (or negative error)
    let ret: i64 = ctx.ret().ok_or(1i64)?;
    if ret <= 0 {
        return Ok(()); // Error or no data
    }

    // We need to track which connection this belongs to
    // This is tricky without the sock pointer - in production use
    // a temporary map to store sock -> key mapping
    
    Ok(())
}

/// Track TCP connection close for cleanup
#[kprobe]
pub fn trace_tcp_close(ctx: ProbeContext) -> u32 {
    match try_trace_tcp_close(&ctx) {
        Ok(()) => 0,
        Err(_) => 1,
    }
}

fn try_trace_tcp_close(ctx: &ProbeContext) -> Result<(), i64> {
    let sock: *const u8 = ctx.arg(0).ok_or(1i64)?;
    let key = unsafe { read_conn_key_from_sock(sock)? };

    // Log final stats before removing
    if let Some(metrics) = unsafe { CONNECTIONS.get(&key) } {
        let duration_ns = unsafe { bpf_ktime_get_ns() } - metrics.start_ns;
        info!(
            ctx,
            "CLOSE: {}:{} -> {}:{} | TX:{} RX:{} RTX:{} dur:{}ms",
            key.src_ip,
            key.src_port,
            key.dst_ip,
            key.dst_port,
            metrics.bytes_sent,
            metrics.bytes_recv,
            metrics.retransmits,
            duration_ns / 1_000_000
        );
    }

    // Remove from map (cleanup)
    let _ = CONNECTIONS.remove(&key);

    Ok(())
}

// ============================================================================
// Tracepoint Programs - Attach to kernel tracepoints
// ============================================================================

/// Track TCP retransmissions via tracepoint
/// This is more reliable than kprobing tcp_retransmit_skb
#[tracepoint]
pub fn trace_tcp_retransmit(ctx: TracePointContext) -> u32 {
    match try_trace_tcp_retransmit(&ctx) {
        Ok(()) => 0,
        Err(_) => 1,
    }
}

fn try_trace_tcp_retransmit(ctx: &TracePointContext) -> Result<(), i64> {
    // Tracepoint format: tcp:tcp_retransmit_skb
    // Fields at specific offsets (check /sys/kernel/debug/tracing/events/tcp/tcp_retransmit_skb/format)
    // This is kernel-version specific
    
    let saddr: u32 = unsafe { ctx.read_at(16)? };
    let daddr: u32 = unsafe { ctx.read_at(20)? };
    let sport: u16 = unsafe { ctx.read_at(24)? };
    let dport: u16 = unsafe { ctx.read_at(26)? };

    let key = ConnKey {
        src_ip: saddr,
        dst_ip: daddr,
        src_port: sport,
        dst_port: dport,
    };

    if let Some(metrics) = unsafe { CONNECTIONS.get_ptr_mut(&key) } {
        let m = unsafe { &mut *metrics };
        m.retransmits += 1;
        
        debug!(ctx, "RETRANSMIT: {}:{} -> {}:{} (count: {})", 
            saddr, sport, daddr, dport, m.retransmits);
    }

    Ok(())
}

// ============================================================================
// Panic Handler (required for no_std)
// ============================================================================

#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}
