//! Metrics collection and aggregation utilities.

use sidecar_common::{ConnKey, ConnMetrics};
use std::collections::HashMap;
use std::net::Ipv4Addr;

/// Aggregated metrics for a destination endpoint.
#[derive(Debug, Default, Clone)]
pub struct EndpointMetrics {
    pub total_bytes_sent: u64,
    pub total_bytes_recv: u64,
    pub total_packets_sent: u64,
    pub total_packets_recv: u64,
    pub total_retransmits: u64,
    pub connection_count: u64,
    pub avg_duration_ms: f64,
}

/// Aggregate per-connection metrics by destination.
pub fn aggregate_by_destination(
    connections: impl Iterator<Item = (ConnKey, ConnMetrics)>,
) -> HashMap<(Ipv4Addr, u16), EndpointMetrics> {
    let mut aggregated: HashMap<(Ipv4Addr, u16), EndpointMetrics> = HashMap::new();

    for (key, metrics) in connections {
        let dst_ip = Ipv4Addr::from(key.dst_ip.to_be());
        let endpoint = (dst_ip, key.dst_port);

        let entry = aggregated.entry(endpoint).or_default();
        entry.total_bytes_sent += metrics.bytes_sent;
        entry.total_bytes_recv += metrics.bytes_recv;
        entry.total_packets_sent += metrics.packets_sent;
        entry.total_packets_recv += metrics.packets_recv;
        entry.total_retransmits += metrics.retransmits as u64;
        entry.connection_count += 1;

        let duration_ms = (metrics.last_seen_ns - metrics.start_ns) as f64 / 1_000_000.0;
        // Running average
        let n = entry.connection_count as f64;
        entry.avg_duration_ms = entry.avg_duration_ms * ((n - 1.0) / n) + duration_ms / n;
    }

    aggregated
}

/// Format bytes as human-readable string.
pub fn format_bytes(bytes: u64) -> String {
    if bytes >= 1_073_741_824 {
        format!("{:.2} GB", bytes as f64 / 1_073_741_824.0)
    } else if bytes >= 1_048_576 {
        format!("{:.2} MB", bytes as f64 / 1_048_576.0)
    } else if bytes >= 1024 {
        format!("{:.2} KB", bytes as f64 / 1024.0)
    } else {
        format!("{} B", bytes)
    }
}

/// Format duration in milliseconds as human-readable string.
pub fn format_duration(ms: f64) -> String {
    if ms >= 60_000.0 {
        format!("{:.1} min", ms / 60_000.0)
    } else if ms >= 1000.0 {
        format!("{:.2} s", ms / 1000.0)
    } else {
        format!("{:.0} ms", ms)
    }
}
