//! Build tooling for eBPF sidecar.
//!
//! Usage:
//!   cargo xtask build           # Build debug
//!   cargo xtask build --release # Build release
//!   cargo xtask build-ebpf      # Build only eBPF programs

use anyhow::Result;
use clap::{Parser, Subcommand};
use std::process::Command;

#[derive(Debug, Parser)]
#[command(name = "xtask")]
#[command(about = "Build tooling for eBPF sidecar")]
struct Args {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Debug, Subcommand)]
enum Commands {
    /// Build eBPF programs only
    BuildEbpf {
        /// Build in release mode
        #[arg(long)]
        release: bool,
    },
    /// Build everything (eBPF + userspace)
    Build {
        /// Build in release mode
        #[arg(long)]
        release: bool,
    },
    /// Run the sidecar (builds first if needed)
    Run {
        /// Arguments to pass to sidecar
        #[arg(trailing_var_arg = true)]
        args: Vec<String>,
    },
}

fn main() -> Result<()> {
    let args = Args::parse();

    match args.command {
        Commands::BuildEbpf { release } => {
            build_ebpf(release)?;
        }
        Commands::Build { release } => {
            build_ebpf(release)?;
            build_userspace(release)?;
        }
        Commands::Run { args: run_args } => {
            build_ebpf(false)?;
            build_userspace(false)?;
            run_sidecar(&run_args)?;
        }
    }

    Ok(())
}

fn build_ebpf(release: bool) -> Result<()> {
    println!("🔧 Building eBPF programs...");

    let mut cmd = Command::new("cargo");
    cmd.current_dir("sidecar-ebpf")
        .env("CARGO_CFG_BPF_TARGET_ARCH", std::env::consts::ARCH)
        .args([
            "+nightly",
            "build",
            "--target=bpfel-unknown-none",
            "-Z",
            "build-std=core",
        ]);

    if release {
        cmd.arg("--release");
    }

    let status = cmd.status()?;
    if !status.success() {
        anyhow::bail!("Failed to build eBPF programs");
    }

    println!("✅ eBPF programs built successfully");
    Ok(())
}

fn build_userspace(release: bool) -> Result<()> {
    println!("🔧 Building userspace loader...");

    let mut cmd = Command::new("cargo");
    cmd.args(["build", "-p", "sidecar"]);

    if release {
        cmd.arg("--release");
    }

    let status = cmd.status()?;
    if !status.success() {
        anyhow::bail!("Failed to build userspace program");
    }

    println!("✅ Userspace program built successfully");
    Ok(())
}

fn run_sidecar(args: &[String]) -> Result<()> {
    println!("🚀 Running sidecar...");
    println!("   (requires root privileges)");

    let mut cmd = Command::new("sudo");
    cmd.arg("./target/debug/sidecar");
    cmd.args(args);

    let status = cmd.status()?;
    if !status.success() {
        anyhow::bail!("Sidecar exited with error");
    }

    Ok(())
}
