//! eBPF Sidecar - Userspace Loader and Metrics Exporter
//!
//! This program loads the eBPF programs into the kernel, attaches them to
//! the appropriate hooks, and exports collected metrics via Prometheus.
//!
//! # Usage
//! ```bash
//! # Monitor all processes
//! sudo ./sidecar
//!
//! # Monitor specific PID
//! sudo ./sidecar --pid 1234
//!
//! # Monitor specific ports
//! sudo ./sidecar --ports 80,443,8080
//!
//! # Custom Prometheus port
//! sudo ./sidecar --metrics-port 9091
//! ```

use anyhow::{Context, Result};
use aya::{
    include_bytes_aligned,
    maps::{Array, HashMap},
    programs::{KProbe, TracePoint},
    Bpf,
};
use aya_log::BpfLogger;
use clap::Parser;
use log::{debug, error, info, warn};
use prometheus::{
    register_counter_vec, register_gauge_vec, register_histogram_vec,
    CounterVec, Encoder, GaugeVec, HistogramVec, TextEncoder,
};
use sidecar_common::{ConnKey, ConnMetrics, SidecarConfig};
use std::convert::Infallible;
use std::net::{Ipv4Addr, SocketAddr};
use std::sync::Arc;
use std::time::Duration;
use tokio::signal;
use tokio::sync::RwLock;
use tokio::time;

mod config;
mod metrics;

use config::Config;

// ============================================================================
// CLI Arguments
// ============================================================================

#[derive(Debug, Parser)]
#[command(name = "sidecar")]
#[command(about = "eBPF-based service mesh sidecar for network observability")]
#[command(version)]
struct Args {
    /// Target PID to monitor (0 = all processes)
    #[arg(short, long, default_value = "0")]
    pid: u32,

    /// Ports to monitor (comma-separated, empty = all)
    #[arg(long, value_delimiter = ',')]
    ports: Option<Vec<u16>>,

    /// Prometheus metrics port
    #[arg(short, long, default_value = "9090")]
    metrics_port: u16,

    /// Metrics collection interval in seconds
    #[arg(short, long, default_value = "5")]
    interval: u64,

    /// Enable debug logging from eBPF programs
    #[arg(short, long)]
    debug: bool,

    /// Config file path (optional)
    #[arg(short, long)]
    config: Option<String>,
}

// ============================================================================
// Prometheus Metrics
// ============================================================================

lazy_static::lazy_static! {
    static ref CONN_BYTES_SENT: CounterVec = register_counter_vec!(
        "sidecar_connection_bytes_sent_total",
        "Total bytes sent per connection",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();

    static ref CONN_BYTES_RECV: CounterVec = register_counter_vec!(
        "sidecar_connection_bytes_received_total",
        "Total bytes received per connection",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();

    static ref CONN_PACKETS_SENT: CounterVec = register_counter_vec!(
        "sidecar_connection_packets_sent_total",
        "Total packets sent per connection",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();

    static ref CONN_PACKETS_RECV: CounterVec = register_counter_vec!(
        "sidecar_connection_packets_received_total",
        "Total packets received per connection",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();

    static ref CONN_RETRANSMITS: CounterVec = register_counter_vec!(
        "sidecar_connection_retransmits_total",
        "Total TCP retransmissions per connection",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();

    static ref CONN_DURATION: GaugeVec = register_gauge_vec!(
        "sidecar_connection_duration_seconds",
        "Connection duration in seconds",
        &["src_ip", "dst_ip", "dst_port"]
    ).unwrap();

    static ref ACTIVE_CONNECTIONS: prometheus::IntGauge = prometheus::register_int_gauge!(
        "sidecar_active_connections",
        "Number of active connections being tracked"
    ).unwrap();
}

// ============================================================================
// Main Entry Point
// ============================================================================

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();

    // Initialize logging
    env_logger::Builder::from_env(
        env_logger::Env::default().default_filter_or(if args.debug { "debug" } else { "info" }),
    )
    .init();

    info!("Starting eBPF sidecar...");
    info!("Target PID: {}", if args.pid == 0 { "all".to_string() } else { args.pid.to_string() });
    info!("Metrics port: {}", args.metrics_port);

    // Bump memlock rlimit for eBPF
    bump_memlock_rlimit()?;

    // Load eBPF program
    let mut bpf = load_ebpf_program()?;

    // Initialize eBPF logging
    if let Err(e) = BpfLogger::init(&mut bpf) {
        warn!("Failed to initialize eBPF logger: {}", e);
    }

    // Configure the sidecar
    configure_sidecar(&mut bpf, &args)?;

    // Attach programs
    attach_programs(&mut bpf)?;

    info!("eBPF programs loaded and attached successfully");

    // Start Prometheus HTTP server
    let metrics_addr: SocketAddr = ([0, 0, 0, 0], args.metrics_port).into();
    tokio::spawn(async move {
        if let Err(e) = run_metrics_server(metrics_addr).await {
            error!("Metrics server error: {}", e);
        }
    });
    info!("Prometheus metrics available at http://0.0.0.0:{}/metrics", args.metrics_port);

    // Get reference to connections map
    let connections: HashMap<_, ConnKey, ConnMetrics> =
        HashMap::try_from(bpf.map("CONNECTIONS").context("Failed to get CONNECTIONS map")?)?;

    // Metrics collection loop
    let mut interval = time::interval(Duration::from_secs(args.interval));

    info!("Sidecar running. Press Ctrl+C to stop.");

    loop {
        tokio::select! {
            _ = interval.tick() => {
                if let Err(e) = collect_and_export_metrics(&connections) {
                    error!("Failed to collect metrics: {}", e);
                }
            }
            _ = signal::ctrl_c() => {
                info!("Received shutdown signal");
                break;
            }
        }
    }

    info!("Sidecar stopped");
    Ok(())
}

// ============================================================================
// eBPF Loading and Setup
// ============================================================================

fn bump_memlock_rlimit() -> Result<()> {
    let rlim = libc::rlimit {
        rlim_cur: libc::RLIM_INFINITY,
        rlim_max: libc::RLIM_INFINITY,
    };
    let ret = unsafe { libc::setrlimit(libc::RLIMIT_MEMLOCK, &rlim) };
    if ret != 0 {
        anyhow::bail!("Failed to set RLIMIT_MEMLOCK");
    }
    Ok(())
}

fn load_ebpf_program() -> Result<Bpf> {
    // Include the compiled eBPF bytecode at compile time
    #[cfg(debug_assertions)]
    let bpf = Bpf::load(include_bytes_aligned!(
        "../../target/bpfel-unknown-none/debug/sidecar"
    ))?;

    #[cfg(not(debug_assertions))]
    let bpf = Bpf::load(include_bytes_aligned!(
        "../../target/bpfel-unknown-none/release/sidecar"
    ))?;

    Ok(bpf)
}

fn configure_sidecar(bpf: &mut Bpf, args: &Args) -> Result<()> {
    let mut config = SidecarConfig::default();
    config.target_pid = args.pid;
    config.debug_mode = if args.debug { 1 } else { 0 };

    // Set target ports if specified
    if let Some(ref ports) = args.ports {
        for (i, port) in ports.iter().take(8).enumerate() {
            config.target_ports[i] = *port;
        }
        config.num_target_ports = ports.len().min(8) as u8;
    }

    // Write config to eBPF map
    let mut config_map: Array<_, SidecarConfig> =
        Array::try_from(bpf.map_mut("CONFIG").context("Failed to get CONFIG map")?)?;
    config_map.set(0, config, 0)?;

    debug!("Configuration applied: {:?}", config);
    Ok(())
}

fn attach_programs(bpf: &mut Bpf) -> Result<()> {
    // Attach kprobes
    let programs = [
        ("trace_tcp_connect", "tcp_connect"),
        ("trace_tcp_sendmsg", "tcp_sendmsg"),
        ("trace_tcp_recvmsg", "tcp_recvmsg"),
        ("trace_tcp_close", "tcp_close"),
    ];

    for (prog_name, fn_name) in programs {
        let program: &mut KProbe = bpf
            .program_mut(prog_name)
            .context(format!("Failed to get program {}", prog_name))?
            .try_into()?;
        program.load()?;
        program.attach(fn_name, 0)?;
        info!("Attached {} to {}", prog_name, fn_name);
    }

    // Attach tracepoint for retransmits
    let tp: &mut TracePoint = bpf
        .program_mut("trace_tcp_retransmit")
        .context("Failed to get trace_tcp_retransmit")?
        .try_into()?;
    tp.load()?;
    tp.attach("tcp", "tcp_retransmit_skb")?;
    info!("Attached trace_tcp_retransmit to tcp:tcp_retransmit_skb");

    Ok(())
}

// ============================================================================
// Metrics Collection and Export
// ============================================================================

fn collect_and_export_metrics(
    connections: &HashMap<&aya::maps::MapData, ConnKey, ConnMetrics>,
) -> Result<()> {
    let mut count = 0;

    for result in connections.iter() {
        let (key, metrics) = result?;

        let src_ip = Ipv4Addr::from(key.src_ip.to_be()).to_string();
        let dst_ip = Ipv4Addr::from(key.dst_ip.to_be()).to_string();
        let dst_port = key.dst_port.to_string();

        // Update Prometheus metrics
        CONN_BYTES_SENT
            .with_label_values(&[&src_ip, &dst_ip, &dst_port])
            .inc_by(metrics.bytes_sent as f64);

        CONN_BYTES_RECV
            .with_label_values(&[&src_ip, &dst_ip, &dst_port])
            .inc_by(metrics.bytes_recv as f64);

        CONN_PACKETS_SENT
            .with_label_values(&[&src_ip, &dst_ip, &dst_port])
            .inc_by(metrics.packets_sent as f64);

        CONN_PACKETS_RECV
            .with_label_values(&[&src_ip, &dst_ip, &dst_port])
            .inc_by(metrics.packets_recv as f64);

        CONN_RETRANSMITS
            .with_label_values(&[&src_ip, &dst_ip, &dst_port])
            .inc_by(metrics.retransmits as f64);

        let duration_secs = (metrics.last_seen_ns - metrics.start_ns) as f64 / 1_000_000_000.0;
        CONN_DURATION
            .with_label_values(&[&src_ip, &dst_ip, &dst_port])
            .set(duration_secs);

        count += 1;
    }

    ACTIVE_CONNECTIONS.set(count);
    debug!("Collected metrics for {} connections", count);

    Ok(())
}

// ============================================================================
// Prometheus HTTP Server
// ============================================================================

async fn run_metrics_server(addr: SocketAddr) -> Result<()> {
    use hyper::service::{make_service_fn, service_fn};
    use hyper::{Body, Request, Response, Server};

    let make_svc = make_service_fn(|_conn| async {
        Ok::<_, Infallible>(service_fn(|req: Request<Body>| async move {
            match req.uri().path() {
                "/metrics" => {
                    let encoder = TextEncoder::new();
                    let metric_families = prometheus::gather();
                    let mut buffer = Vec::new();
                    encoder.encode(&metric_families, &mut buffer).unwrap();
                    Ok::<_, Infallible>(Response::new(Body::from(buffer)))
                }
                "/health" => Ok(Response::new(Body::from("OK"))),
                _ => Ok(Response::builder()
                    .status(404)
                    .body(Body::from("Not Found"))
                    .unwrap()),
            }
        }))
    });

    Server::bind(&addr).serve(make_svc).await?;
    Ok(())
}
