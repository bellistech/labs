//! Configuration handling for the sidecar.

use serde::{Deserialize, Serialize};
use std::path::Path;

/// Sidecar configuration loaded from YAML file.
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Config {
    /// Target process configuration
    #[serde(default)]
    pub target: TargetConfig,

    /// Metrics export configuration
    #[serde(default)]
    pub metrics: MetricsConfig,

    /// Logging configuration
    #[serde(default)]
    pub logging: LoggingConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct TargetConfig {
    /// PID to monitor (0 = all)
    #[serde(default)]
    pub pid: u32,

    /// Process name to monitor (alternative to PID)
    #[serde(default)]
    pub process_name: Option<String>,

    /// cgroup path to monitor (for container filtering)
    #[serde(default)]
    pub cgroup: Option<String>,

    /// Ports to monitor (empty = all)
    #[serde(default)]
    pub ports: Vec<u16>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricsConfig {
    /// Prometheus metrics port
    #[serde(default = "default_metrics_port")]
    pub port: u16,

    /// Collection interval in seconds
    #[serde(default = "default_interval")]
    pub interval_secs: u64,

    /// Enable HTTP layer 7 metrics
    #[serde(default)]
    pub enable_http: bool,
}

impl Default for MetricsConfig {
    fn default() -> Self {
        Self {
            port: 9090,
            interval_secs: 5,
            enable_http: false,
        }
    }
}

fn default_metrics_port() -> u16 {
    9090
}

fn default_interval() -> u64 {
    5
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LoggingConfig {
    /// Log level: trace, debug, info, warn, error
    #[serde(default = "default_log_level")]
    pub level: String,

    /// Enable eBPF debug logging
    #[serde(default)]
    pub ebpf_debug: bool,
}

impl Default for LoggingConfig {
    fn default() -> Self {
        Self {
            level: "info".to_string(),
            ebpf_debug: false,
        }
    }
}

fn default_log_level() -> String {
    "info".to_string()
}

impl Config {
    /// Load configuration from a YAML file.
    pub fn load<P: AsRef<Path>>(path: P) -> anyhow::Result<Self> {
        let contents = std::fs::read_to_string(path)?;
        let config: Config = serde_yaml::from_str(&contents)?;
        Ok(config)
    }

    /// Create default configuration.
    pub fn default_config() -> Self {
        Self::default()
    }
}
