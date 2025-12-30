.PHONY: all build build-ebpf build-release clean run help

# Default target
all: build

# Build everything (debug)
build:
	cargo xtask build

# Build only eBPF programs
build-ebpf:
	cargo xtask build-ebpf

# Build release version
build-release:
	cargo xtask build --release

# Clean build artifacts
clean:
	cargo clean

# Run the sidecar (requires root)
run: build
	sudo ./target/debug/sidecar

# Run with specific PID
run-pid: build
	@read -p "Enter PID to monitor: " pid; \
	sudo ./target/debug/sidecar --pid $$pid

# Run with debug logging
run-debug: build
	sudo ./target/debug/sidecar --debug

# Install dependencies (Ubuntu/Debian)
deps:
	sudo apt-get update
	sudo apt-get install -y \
		build-essential \
		pkg-config \
		libssl-dev \
		linux-headers-$$(uname -r) \
		clang \
		llvm
	rustup install nightly
	rustup component add rust-src --toolchain nightly
	cargo install bpf-linker

# Check if system supports eBPF
check:
	@echo "Kernel version: $$(uname -r)"
	@echo "BTF support: $$(ls /sys/kernel/btf/vmlinux 2>/dev/null && echo 'YES' || echo 'NO')"
	@echo "BPF filesystem: $$(mount | grep bpf || echo 'Not mounted')"

# Show help
help:
	@echo "eBPF Sidecar Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make build        - Build debug version"
	@echo "  make build-release - Build release version"
	@echo "  make run          - Build and run (requires root)"
	@echo "  make run-pid      - Run with specific PID filter"
	@echo "  make run-debug    - Run with debug logging"
	@echo "  make deps         - Install dependencies"
	@echo "  make check        - Check system eBPF support"
	@echo "  make clean        - Clean build artifacts"
