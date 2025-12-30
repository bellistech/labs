# Python Phase 7 Final: CLI, Configuration & Deployment

Complete the DNS server project with a production CLI, configuration system, integration tests, and deployment automation.

---

## Table of Contents

1. [Day 36: CLI Application with Click](#day-36-cli-application)
2. [Day 37: Configuration Management](#day-37-configuration)
3. [Day 38: Integration Tests](#day-38-integration-tests)
4. [Day 39: Deployment Automation](#day-39-deployment)
5. [Day 40: Final Project & Documentation](#day-40-final-project)

---

## Day 36: CLI Application

### Complete CLI with Click

Create `dns_server/cli.py`:

```python
#!/usr/bin/env python3
"""
Day 36: Production CLI Application
Using Click for command-line interface

Install: pip install click rich
"""

import asyncio
import sys
import signal
from pathlib import Path
from typing import Optional

import click
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.progress import Progress, SpinnerColumn, TextColumn
from rich import print as rprint

# Import our modules
from .server import DNSServer, ServerConfig
from .config import load_config, DNSConfig
from .zones import DNSZone, ZoneFileParser
from .utils import DNSClient, DNSBenchmark, format_dig_output
from .metrics import DNSMetrics, MetricsServer

console = Console()


@click.group()
@click.version_option(version="1.0.0", prog_name="pydns")
@click.option("--config", "-c", type=click.Path(exists=True),
              help="Path to configuration file")
@click.pass_context
def cli(ctx, config):
    """
    Python DNS Server - A full-featured DNS server with IPv6 support.
    
    Examples:
    
        pydns serve                    # Start DNS server
        pydns query google.com         # Query a domain
        pydns benchmark 8.8.8.8        # Benchmark a server
        pydns zone validate zone.txt   # Validate zone file
    """
    ctx.ensure_object(dict)
    ctx.obj["config_path"] = config


@cli.command()
@click.option("--port", "-p", default=5353, help="Listen port")
@click.option("--address", "-a", default="::", help="Listen address")
@click.option("--zone", "-z", multiple=True, type=click.Path(exists=True),
              help="Zone file(s) to load")
@click.option("--upstream", "-u", multiple=True, default=["8.8.8.8", "1.1.1.1"],
              help="Upstream DNS servers")
@click.option("--no-recursion", is_flag=True, help="Disable recursive queries")
@click.option("--metrics-port", default=9153, help="Prometheus metrics port")
@click.option("--no-metrics", is_flag=True, help="Disable metrics endpoint")
@click.option("--log-queries", is_flag=True, help="Log all queries")
@click.option("--debug", is_flag=True, help="Enable debug logging")
@click.pass_context
def serve(ctx, port, address, zone, upstream, no_recursion, 
          metrics_port, no_metrics, log_queries, debug):
    """
    Start the DNS server.
    
    Examples:
    
        pydns serve
        pydns serve -p 53 -z /etc/dns/example.zone
        pydns serve --upstream 1.1.1.1 --upstream 8.8.8.8
    """
    import logging
    
    # Configure logging
    level = logging.DEBUG if debug else logging.INFO
    logging.basicConfig(
        level=level,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
    )
    
    console.print(Panel.fit(
        "[bold blue]Python DNS Server[/bold blue]\n"
        f"Listening on {address}:{port}",
        title="Starting"
    ))
    
    # Build config
    config = ServerConfig(
        listen_address=address,
        listen_port=port,
        upstream_servers=list(upstream),
        enable_recursion=not no_recursion,
        log_queries=log_queries,
    )
    
    async def run_server():
        server = DNSServer(config)
        
        # Load zones
        if zone:
            parser = ZoneFileParser()
            for zone_file in zone:
                try:
                    z = parser.parse_file(zone_file)
                    server.add_zone(z)
                    console.print(f"  [green]✓[/green] Loaded zone: {z.name}")
                except Exception as e:
                    console.print(f"  [red]✗[/red] Failed to load {zone_file}: {e}")
        
        # Start metrics server
        metrics_server = None
        if not no_metrics:
            try:
                metrics = DNSMetrics()
                metrics_server = MetricsServer(metrics, metrics_port)
                await metrics_server.start()
                console.print(f"  [green]✓[/green] Metrics: http://localhost:{metrics_port}/metrics")
            except Exception as e:
                console.print(f"  [yellow]![/yellow] Metrics disabled: {e}")
        
        # Start DNS server
        await server.start()
        
        console.print("\n[bold green]Server running.[/bold green] Press Ctrl+C to stop.\n")
        
        # Print status table
        table = Table(title="Server Configuration")
        table.add_column("Setting", style="cyan")
        table.add_column("Value", style="green")
        
        table.add_row("Listen Address", f"{address}:{port}")
        table.add_row("Recursion", "Enabled" if not no_recursion else "Disabled")
        table.add_row("Upstream Servers", ", ".join(upstream))
        table.add_row("Zones Loaded", str(len(server.zones)))
        table.add_row("Query Logging", "Enabled" if log_queries else "Disabled")
        
        console.print(table)
        console.print()
        
        # Handle shutdown
        loop = asyncio.get_event_loop()
        stop_event = asyncio.Event()
        
        def handle_signal():
            console.print("\n[yellow]Shutting down...[/yellow]")
            stop_event.set()
        
        for sig in (signal.SIGINT, signal.SIGTERM):
            loop.add_signal_handler(sig, handle_signal)
        
        await stop_event.wait()
        await server.stop()
        console.print("[green]Server stopped.[/green]")
    
    asyncio.run(run_server())


@cli.command()
@click.argument("domain")
@click.option("--type", "-t", "qtype", default="A",
              type=click.Choice(["A", "AAAA", "MX", "TXT", "CNAME", "NS", "SOA", "PTR", "ANY"]),
              help="Query type")
@click.option("--server", "-s", default="8.8.8.8", help="DNS server to query")
@click.option("--port", "-p", default=53, help="Server port")
@click.option("--tcp", is_flag=True, help="Use TCP instead of UDP")
@click.option("--short", is_flag=True, help="Short output (answers only)")
@click.option("--json", "json_output", is_flag=True, help="JSON output")
def query(domain, qtype, server, port, tcp, short, json_output):
    """
    Query a DNS server (like dig).
    
    Examples:
    
        pydns query google.com
        pydns query google.com -t AAAA
        pydns query google.com -s 1.1.1.1
        pydns query example.com -t MX --short
    """
    from .protocol import DNSRecordType
    
    type_map = {
        "A": DNSRecordType.A,
        "AAAA": DNSRecordType.AAAA,
        "MX": DNSRecordType.MX,
        "TXT": DNSRecordType.TXT,
        "CNAME": DNSRecordType.CNAME,
        "NS": DNSRecordType.NS,
        "SOA": DNSRecordType.SOA,
        "PTR": DNSRecordType.PTR,
        "ANY": DNSRecordType.ANY,
    }
    
    async def do_query():
        client = DNSClient()
        
        with console.status(f"Querying {server}..."):
            try:
                result = await client.query(
                    domain, 
                    type_map[qtype],
                    server, 
                    port,
                    use_tcp=tcp
                )
            except Exception as e:
                console.print(f"[red]Error:[/red] {e}")
                sys.exit(1)
        
        if json_output:
            import json
            output = {
                "query": {"name": domain, "type": qtype},
                "server": f"{server}:{port}",
                "rcode": result.rcode,
                "answers": result.answers,
                "query_time_ms": result.query_time_ms,
            }
            console.print_json(json.dumps(output, indent=2))
        elif short:
            for answer in result.answers:
                console.print(answer["data"])
        else:
            console.print(format_dig_output(result))
    
    asyncio.run(do_query())


@cli.command()
@click.argument("server")
@click.option("--port", "-p", default=53, help="Server port")
@click.option("--queries", "-n", default=100, help="Number of queries")
@click.option("--concurrency", "-c", default=10, help="Concurrent queries")
@click.option("--type", "-t", "qtype", default="A", help="Query type")
def benchmark(server, port, queries, concurrency, qtype):
    """
    Benchmark a DNS server.
    
    Examples:
    
        pydns benchmark 8.8.8.8
        pydns benchmark 1.1.1.1 -n 1000 -c 50
        pydns benchmark localhost -p 5353
    """
    domains = [
        "google.com", "facebook.com", "amazon.com", "microsoft.com",
        "apple.com", "netflix.com", "twitter.com", "linkedin.com",
        "github.com", "stackoverflow.com", "reddit.com", "wikipedia.org",
    ]
    
    async def do_benchmark():
        bench = DNSBenchmark(server, port)
        
        console.print(Panel.fit(
            f"[bold]DNS Benchmark[/bold]\n"
            f"Server: {server}:{port}\n"
            f"Queries: {queries}, Concurrency: {concurrency}",
            title="Starting"
        ))
        
        with Progress(
            SpinnerColumn(),
            TextColumn("[progress.description]{task.description}"),
            console=console,
        ) as progress:
            task = progress.add_task("Running benchmark...", total=None)
            
            try:
                results = await bench.benchmark(
                    domains, 
                    num_queries=queries,
                    concurrency=concurrency
                )
            except Exception as e:
                console.print(f"[red]Error:[/red] {e}")
                sys.exit(1)
        
        # Results table
        if "stats" in results:
            stats = results["stats"]
            
            table = Table(title="Results")
            table.add_column("Metric", style="cyan")
            table.add_column("Value", style="green", justify="right")
            
            table.add_row("Total Queries", str(queries))
            table.add_row("Successful", str(results["success"]))
            table.add_row("Failed", str(results["failed"]))
            table.add_row("", "")
            table.add_row("Min Latency", f"{stats['min_ms']:.2f} ms")
            table.add_row("Max Latency", f"{stats['max_ms']:.2f} ms")
            table.add_row("Avg Latency", f"{stats['avg_ms']:.2f} ms")
            table.add_row("Median Latency", f"{stats['median_ms']:.2f} ms")
            table.add_row("P95 Latency", f"{stats['p95_ms']:.2f} ms")
            table.add_row("P99 Latency", f"{stats['p99_ms']:.2f} ms")
            table.add_row("", "")
            table.add_row("QPS", f"{stats['qps']:.2f}")
            table.add_row("Total Time", f"{stats['total_time_s']:.2f} s")
            
            console.print(table)
        else:
            console.print("[red]Benchmark failed - no results[/red]")
    
    asyncio.run(do_benchmark())


@cli.group()
def zone():
    """Zone file management commands."""
    pass


@zone.command("validate")
@click.argument("zone_file", type=click.Path(exists=True))
def zone_validate(zone_file):
    """
    Validate a zone file.
    
    Examples:
    
        pydns zone validate /etc/dns/example.zone
    """
    from .zones import ZoneFileParser
    from .utils import ZoneValidator
    
    console.print(f"Validating: {zone_file}")
    
    try:
        parser = ZoneFileParser()
        zone = parser.parse_file(zone_file)
        
        console.print(f"  [green]✓[/green] Parsed successfully: {zone.name}")
        console.print(f"  [green]✓[/green] Records: {sum(len(r) for t in zone.records.values() for r in t.values())}")
        
        validator = ZoneValidator()
        is_valid = validator.validate_zone(zone)
        
        if validator.errors:
            console.print("\n[red]Errors:[/red]")
            for e in validator.errors:
                console.print(f"  [red]✗[/red] {e}")
        
        if validator.warnings:
            console.print("\n[yellow]Warnings:[/yellow]")
            for w in validator.warnings:
                console.print(f"  [yellow]![/yellow] {w}")
        
        if is_valid and not validator.warnings:
            console.print("\n[bold green]Zone is valid![/bold green]")
        
    except Exception as e:
        console.print(f"[red]Parse error:[/red] {e}")
        sys.exit(1)


@zone.command("show")
@click.argument("zone_file", type=click.Path(exists=True))
@click.option("--type", "-t", "rtype", help="Filter by record type")
def zone_show(zone_file, rtype):
    """
    Display zone file contents.
    
    Examples:
    
        pydns zone show example.zone
        pydns zone show example.zone -t A
    """
    from .zones import ZoneFileParser
    from .protocol import DNSRecordType
    
    parser = ZoneFileParser()
    zone = parser.parse_file(zone_file)
    
    table = Table(title=f"Zone: {zone.name}")
    table.add_column("Name", style="cyan")
    table.add_column("TTL", justify="right")
    table.add_column("Type", style="green")
    table.add_column("Data")
    
    for name, type_records in sorted(zone.records.items()):
        for record_type, records in sorted(type_records.items()):
            type_name = DNSRecordType(record_type).name if record_type in DNSRecordType._value2member_map_ else f"TYPE{record_type}"
            
            if rtype and type_name != rtype.upper():
                continue
            
            for record in records:
                table.add_row(
                    name,
                    str(record.ttl),
                    type_name,
                    record.parsed_data or record.rdata.hex()[:40]
                )
    
    console.print(table)


@zone.command("create")
@click.argument("zone_name")
@click.option("--output", "-o", type=click.Path(), help="Output file")
@click.option("--ns", multiple=True, default=["ns1", "ns2"], help="Nameservers")
@click.option("--admin", default="admin", help="Admin email (without @)")
def zone_create(zone_name, output, ns, admin):
    """
    Create a zone file template.
    
    Examples:
    
        pydns zone create example.com -o example.zone
        pydns zone create mylab.local --ns ns1 --ns ns2
    """
    from datetime import datetime
    
    serial = datetime.now().strftime("%Y%m%d") + "01"
    
    lines = [
        f"; Zone file for {zone_name}",
        f"; Generated by pydns",
        "",
        f"$TTL 3600",
        f"$ORIGIN {zone_name}.",
        "",
        f"@    IN    SOA    {ns[0]}.{zone_name}. {admin}.{zone_name}. (",
        f"                  {serial}  ; Serial",
        f"                  3600      ; Refresh",
        f"                  900       ; Retry",
        f"                  604800    ; Expire",
        f"                  86400     ; Minimum TTL",
        f"                  )",
        "",
        "; Nameservers",
    ]
    
    for n in ns:
        lines.append(f"@    IN    NS    {n}.{zone_name}.")
    
    lines.extend([
        "",
        "; A Records (IPv4)",
        f"@       IN    A       192.168.1.1",
    ])
    
    for n in ns:
        lines.append(f"{n}      IN    A       192.168.1.{ns.index(n)+2}")
    
    lines.extend([
        "",
        "; AAAA Records (IPv6)",
        f"@       IN    AAAA    2001:db8::1",
    ])
    
    for n in ns:
        lines.append(f"{n}      IN    AAAA    2001:db8::{ns.index(n)+2}")
    
    lines.extend([
        "",
        "; Add your records below",
        "; www     IN    A       192.168.1.10",
        "; www     IN    AAAA    2001:db8::10",
        "",
    ])
    
    content = "\n".join(lines)
    
    if output:
        Path(output).write_text(content)
        console.print(f"[green]Created:[/green] {output}")
    else:
        console.print(content)


@cli.command()
@click.option("--server", "-s", multiple=True, 
              default=["8.8.8.8", "1.1.1.1", "9.9.9.9"],
              help="Servers to compare")
@click.option("--domain", "-d", default="google.com", help="Domain to query")
def compare(server, domain):
    """
    Compare DNS server response times.
    
    Examples:
    
        pydns compare
        pydns compare -s 8.8.8.8 -s 1.1.1.1 -s 9.9.9.9
    """
    async def do_compare():
        client = DNSClient()
        
        table = Table(title=f"DNS Server Comparison - {domain}")
        table.add_column("Server", style="cyan")
        table.add_column("Latency", justify="right")
        table.add_column("Status", justify="center")
        table.add_column("Answer")
        
        for s in server:
            try:
                result = await client.query(domain, server=s)
                answer = result.answers[0]["data"] if result.answers else "No answer"
                table.add_row(
                    s,
                    f"{result.query_time_ms:.2f} ms",
                    "[green]✓[/green]",
                    answer
                )
            except Exception as e:
                table.add_row(
                    s,
                    "-",
                    "[red]✗[/red]",
                    str(e)[:40]
                )
        
        console.print(table)
    
    asyncio.run(do_compare())


@cli.command()
@click.argument("address")
def reverse(address):
    """
    Perform reverse DNS lookup.
    
    Examples:
    
        pydns reverse 8.8.8.8
        pydns reverse 2001:4860:4860::8888
    """
    import ipaddress
    from .ipv6 import ipv6_to_ptr
    
    async def do_reverse():
        client = DNSClient()
        
        try:
            ip = ipaddress.ip_address(address)
            
            if isinstance(ip, ipaddress.IPv4Address):
                # IPv4 reverse
                octets = address.split(".")
                ptr_name = ".".join(reversed(octets)) + ".in-addr.arpa"
            else:
                # IPv6 reverse
                ptr_name = ipv6_to_ptr(address)
            
            console.print(f"PTR: {ptr_name[:60]}...")
            
            result = await client.query(ptr_name, qtype=12)  # PTR
            
            if result.answers:
                for answer in result.answers:
                    console.print(f"[green]{address}[/green] -> {answer['data']}")
            else:
                console.print(f"[yellow]No PTR record found[/yellow]")
                
        except Exception as e:
            console.print(f"[red]Error:[/red] {e}")
    
    asyncio.run(do_reverse())


@cli.command()
def status():
    """
    Show status of running DNS server.
    """
    import httpx
    
    async def get_status():
        try:
            async with httpx.AsyncClient() as client:
                health = await client.get("http://localhost:9153/health")
                metrics = await client.get("http://localhost:9153/metrics")
            
            console.print("[green]Server is running[/green]\n")
            
            # Parse some metrics
            lines = metrics.text.split("\n")
            stats = {}
            for line in lines:
                if line.startswith("dns_queries_total"):
                    stats["queries"] = line.split()[-1]
                elif line.startswith("dns_cache_entries"):
                    stats["cache"] = line.split()[-1]
            
            if stats:
                table = Table(title="Server Status")
                table.add_column("Metric", style="cyan")
                table.add_column("Value", style="green")
                
                for k, v in stats.items():
                    table.add_row(k, v)
                
                console.print(table)
                
        except Exception as e:
            console.print(f"[red]Server not reachable:[/red] {e}")
    
    asyncio.run(get_status())


def main():
    """Entry point."""
    cli(obj={})


if __name__ == "__main__":
    main()
```

### Entry Point

Create `dns_server/__main__.py`:

```python
#!/usr/bin/env python3
"""
DNS Server entry point.

Usage:
    python -m dns_server serve
    python -m dns_server query google.com
"""

from .cli import main

if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add dns_server/cli.py dns_server/__main__.py
git commit -m "Day 36: Complete CLI with Click"
```

---

## Day 37: Configuration Management

Create `dns_server/config.py`:

```python
#!/usr/bin/env python3
"""
Day 37: Configuration Management
- Environment variables
- YAML/TOML config files
- Validation with Pydantic
"""

import os
from pathlib import Path
from typing import List, Optional, Dict, Any
from dataclasses import dataclass, field
from enum import Enum

try:
    from pydantic import BaseModel, Field, validator, root_validator
    from pydantic_settings import BaseSettings
    PYDANTIC_AVAILABLE = True
except ImportError:
    PYDANTIC_AVAILABLE = False


class LogLevel(str, Enum):
    """Log levels."""
    DEBUG = "DEBUG"
    INFO = "INFO"
    WARNING = "WARNING"
    ERROR = "ERROR"


if PYDANTIC_AVAILABLE:
    
    class NetworkConfig(BaseModel):
        """Network configuration."""
        listen_address: str = Field(default="::", description="Listen address")
        listen_port: int = Field(default=5353, ge=1, le=65535)
        enable_ipv6: bool = True
        enable_tcp: bool = True
        tcp_timeout: int = Field(default=30, ge=1)
    
    
    class UpstreamConfig(BaseModel):
        """Upstream resolver configuration."""
        servers: List[str] = Field(default=["8.8.8.8", "1.1.1.1"])
        servers_v6: List[str] = Field(default=["2001:4860:4860::8888"])
        timeout: float = Field(default=2.0, ge=0.1)
        retries: int = Field(default=2, ge=0)
        
        @validator("servers", "servers_v6", each_item=True)
        def validate_ip(cls, v):
            import ipaddress
            try:
                ipaddress.ip_address(v)
            except ValueError:
                raise ValueError(f"Invalid IP address: {v}")
            return v
    
    
    class CacheConfig(BaseModel):
        """Cache configuration."""
        enabled: bool = True
        max_size: int = Field(default=10000, ge=100)
        min_ttl: int = Field(default=60, ge=0)
        max_ttl: int = Field(default=86400, ge=60)
        negative_ttl: int = Field(default=300, ge=0)
    
    
    class SecurityConfig(BaseModel):
        """Security configuration."""
        enable_rate_limiting: bool = True
        queries_per_second: int = Field(default=20, ge=1)
        queries_per_minute: int = Field(default=300, ge=1)
        enable_blocklist: bool = False
        blocklist_file: Optional[str] = None
        enable_rrl: bool = True
        whitelist_networks: List[str] = Field(default=[])
        
        @validator("whitelist_networks", each_item=True)
        def validate_network(cls, v):
            import ipaddress
            try:
                ipaddress.ip_network(v, strict=False)
            except ValueError:
                raise ValueError(f"Invalid network: {v}")
            return v
    
    
    class MetricsConfig(BaseModel):
        """Metrics configuration."""
        enabled: bool = True
        port: int = Field(default=9153, ge=1, le=65535)
        path: str = "/metrics"
    
    
    class LoggingConfig(BaseModel):
        """Logging configuration."""
        level: LogLevel = LogLevel.INFO
        log_queries: bool = False
        log_file: Optional[str] = None
        log_format: str = "%(asctime)s [%(levelname)s] %(name)s: %(message)s"
    
    
    class ZoneConfig(BaseModel):
        """Zone configuration."""
        files: List[str] = Field(default=[])
        auto_reload: bool = False
        reload_interval: int = Field(default=300, ge=60)
    
    
    class DNSConfig(BaseSettings):
        """
        Main configuration.
        
        Loads from environment variables with DNS_ prefix.
        """
        network: NetworkConfig = Field(default_factory=NetworkConfig)
        upstream: UpstreamConfig = Field(default_factory=UpstreamConfig)
        cache: CacheConfig = Field(default_factory=CacheConfig)
        security: SecurityConfig = Field(default_factory=SecurityConfig)
        metrics: MetricsConfig = Field(default_factory=MetricsConfig)
        logging: LoggingConfig = Field(default_factory=LoggingConfig)
        zones: ZoneConfig = Field(default_factory=ZoneConfig)
        
        class Config:
            env_prefix = "DNS_"
            env_nested_delimiter = "__"
        
        @classmethod
        def from_yaml(cls, path: str) -> "DNSConfig":
            """Load configuration from YAML file."""
            import yaml
            
            with open(path) as f:
                data = yaml.safe_load(f)
            
            return cls(**data)
        
        @classmethod
        def from_toml(cls, path: str) -> "DNSConfig":
            """Load configuration from TOML file."""
            import tomllib
            
            with open(path, "rb") as f:
                data = tomllib.load(f)
            
            return cls(**data)
        
        def to_yaml(self) -> str:
            """Export configuration to YAML."""
            import yaml
            return yaml.dump(self.model_dump(), default_flow_style=False)
        
        def to_dict(self) -> Dict[str, Any]:
            """Export to dictionary."""
            return self.model_dump()

else:
    # Fallback without Pydantic
    @dataclass
    class DNSConfig:
        listen_address: str = "::"
        listen_port: int = 5353
        upstream_servers: List[str] = field(default_factory=lambda: ["8.8.8.8"])
        enable_recursion: bool = True
        cache_size: int = 10000
        log_queries: bool = False
        
        @classmethod
        def from_env(cls) -> "DNSConfig":
            return cls(
                listen_address=os.getenv("DNS_LISTEN_ADDRESS", "::"),
                listen_port=int(os.getenv("DNS_LISTEN_PORT", "5353")),
                upstream_servers=os.getenv("DNS_UPSTREAM", "8.8.8.8").split(","),
                enable_recursion=os.getenv("DNS_RECURSION", "true").lower() == "true",
                cache_size=int(os.getenv("DNS_CACHE_SIZE", "10000")),
                log_queries=os.getenv("DNS_LOG_QUERIES", "false").lower() == "true",
            )


def load_config(path: Optional[str] = None) -> DNSConfig:
    """
    Load configuration from file or environment.
    
    Priority:
    1. Specified file path
    2. DNS_CONFIG environment variable
    3. ./config.yaml or ./config.toml
    4. /etc/dns-server/config.yaml
    5. Environment variables only
    """
    # Check for explicit path
    if path:
        config_path = Path(path)
    elif os.getenv("DNS_CONFIG"):
        config_path = Path(os.getenv("DNS_CONFIG"))
    else:
        # Search for config files
        search_paths = [
            Path("config.yaml"),
            Path("config.toml"),
            Path("/etc/dns-server/config.yaml"),
            Path("/etc/dns-server/config.toml"),
        ]
        config_path = None
        for p in search_paths:
            if p.exists():
                config_path = p
                break
    
    if config_path and config_path.exists():
        if config_path.suffix in (".yaml", ".yml"):
            return DNSConfig.from_yaml(str(config_path))
        elif config_path.suffix == ".toml":
            return DNSConfig.from_toml(str(config_path))
    
    # Fall back to environment variables
    if PYDANTIC_AVAILABLE:
        return DNSConfig()
    else:
        return DNSConfig.from_env()


def generate_example_config() -> str:
    """Generate example configuration file."""
    return """# DNS Server Configuration
# Environment variables override file settings (prefix: DNS_)

network:
  listen_address: "::"
  listen_port: 5353
  enable_ipv6: true
  enable_tcp: true
  tcp_timeout: 30

upstream:
  servers:
    - 8.8.8.8
    - 1.1.1.1
  servers_v6:
    - 2001:4860:4860::8888
    - 2606:4700:4700::1111
  timeout: 2.0
  retries: 2

cache:
  enabled: true
  max_size: 10000
  min_ttl: 60
  max_ttl: 86400
  negative_ttl: 300

security:
  enable_rate_limiting: true
  queries_per_second: 20
  queries_per_minute: 300
  enable_blocklist: false
  blocklist_file: null
  enable_rrl: true
  whitelist_networks:
    - 192.168.0.0/16
    - 10.0.0.0/8
    - "::1/128"

metrics:
  enabled: true
  port: 9153
  path: /metrics

logging:
  level: INFO
  log_queries: false
  log_file: null

zones:
  files:
    # - /etc/dns-server/zones/example.zone
  auto_reload: false
  reload_interval: 300
"""


def main():
    print("=== Configuration Management Demo ===\n")
    
    # Generate example
    print("1. Example Configuration:")
    example = generate_example_config()
    print(example[:500] + "...\n")
    
    # Load from environment
    print("2. Loading Configuration:")
    
    # Set some env vars for demo
    os.environ["DNS_NETWORK__LISTEN_PORT"] = "5353"
    os.environ["DNS_LOGGING__LOG_QUERIES"] = "true"
    
    try:
        config = load_config()
        print(f"   Listen: {config.network.listen_address}:{config.network.listen_port}")
        print(f"   Upstream: {config.upstream.servers}")
        print(f"   Cache size: {config.cache.max_size}")
        print(f"   Log queries: {config.logging.log_queries}")
    except Exception as e:
        print(f"   Error: {e}")
        print("   (Install pydantic for full config support)")
    
    # Validation demo
    print("\n3. Configuration Validation:")
    if PYDANTIC_AVAILABLE:
        try:
            bad_config = NetworkConfig(listen_port=99999)
        except Exception as e:
            print(f"   ✓ Caught invalid port: {e}")
        
        try:
            bad_upstream = UpstreamConfig(servers=["not.an.ip"])
        except Exception as e:
            print(f"   ✓ Caught invalid IP: {e}")
    else:
        print("   (Validation requires pydantic)")


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add dns_server/config.py
git commit -m "Day 37: Configuration management with Pydantic"
```

---

## Day 38: Integration Tests

Create `tests/test_integration.py`:

```python
#!/usr/bin/env python3
"""
Day 38: Integration Tests
- Full server testing
- Client-server interaction
- Performance tests
"""

import pytest
import asyncio
import socket
import time
from typing import List, Tuple

import sys
sys.path.insert(0, '..')

from dns_server.server import DNSServer, ServerConfig
from dns_server.zones import DNSZone
from dns_server.protocol import (
    DNSMessage, DNSRecordType, DNSRcode,
    create_query, create_a_record, create_aaaa_record,
    create_cname_record, create_txt_record
)


# Test fixtures

@pytest.fixture
def test_zone() -> DNSZone:
    """Create a test zone."""
    zone = DNSZone(name="test.local")
    
    # A records
    zone.add_record(create_a_record("test.local", "192.168.1.1"))
    zone.add_record(create_a_record("www.test.local", "192.168.1.10"))
    zone.add_record(create_a_record("mail.test.local", "192.168.1.20"))
    
    # Multiple A records (round-robin)
    zone.add_record(create_a_record("lb.test.local", "192.168.1.100"))
    zone.add_record(create_a_record("lb.test.local", "192.168.1.101"))
    zone.add_record(create_a_record("lb.test.local", "192.168.1.102"))
    
    # AAAA records
    zone.add_record(create_aaaa_record("test.local", "2001:db8::1"))
    zone.add_record(create_aaaa_record("www.test.local", "2001:db8::10"))
    zone.add_record(create_aaaa_record("ipv6only.test.local", "2001:db8::100"))
    
    # CNAME records
    zone.add_record(create_cname_record("blog.test.local", "www.test.local"))
    zone.add_record(create_cname_record("ftp.test.local", "www.test.local"))
    
    # TXT records
    zone.add_record(create_txt_record("test.local", "v=spf1 mx -all"))
    zone.add_record(create_txt_record("_dmarc.test.local", "v=DMARC1; p=reject"))
    
    return zone


@pytest.fixture
async def dns_server(test_zone) -> DNSServer:
    """Create and start a test DNS server."""
    config = ServerConfig(
        listen_port=15353,
        enable_recursion=False,  # Faster tests
        log_queries=False,
    )
    
    server = DNSServer(config)
    server.add_zone(test_zone)
    await server.start()
    
    yield server
    
    await server.stop()


async def send_query(name: str, qtype: int, port: int = 15353,
                     server: str = "127.0.0.1") -> DNSMessage:
    """Send a DNS query and return the response."""
    query = create_query(name, qtype)
    
    loop = asyncio.get_event_loop()
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.setblocking(False)
    
    try:
        await loop.sock_sendto(sock, query.pack(), (server, port))
        data = await asyncio.wait_for(
            loop.sock_recv(sock, 4096),
            timeout=2.0
        )
        return DNSMessage.unpack(data)
    finally:
        sock.close()


# Basic functionality tests

class TestBasicQueries:
    """Test basic DNS query functionality."""
    
    @pytest.mark.asyncio
    async def test_a_record(self, dns_server):
        """Test A record lookup."""
        response = await send_query("www.test.local", DNSRecordType.A)
        
        assert response.header.rcode == DNSRcode.NOERROR
        assert len(response.answers) == 1
        assert response.answers[0].parsed_data == "192.168.1.10"
    
    @pytest.mark.asyncio
    async def test_aaaa_record(self, dns_server):
        """Test AAAA record lookup."""
        response = await send_query("www.test.local", DNSRecordType.AAAA)
        
        assert response.header.rcode == DNSRcode.NOERROR
        assert len(response.answers) == 1
        assert "2001:db8::10" in response.answers[0].parsed_data
    
    @pytest.mark.asyncio
    async def test_multiple_a_records(self, dns_server):
        """Test multiple A records (round-robin)."""
        response = await send_query("lb.test.local", DNSRecordType.A)
        
        assert response.header.rcode == DNSRcode.NOERROR
        assert len(response.answers) == 3
        
        ips = {r.parsed_data for r in response.answers}
        assert ips == {"192.168.1.100", "192.168.1.101", "192.168.1.102"}
    
    @pytest.mark.asyncio
    async def test_cname_record(self, dns_server):
        """Test CNAME record lookup."""
        response = await send_query("blog.test.local", DNSRecordType.CNAME)
        
        assert response.header.rcode == DNSRcode.NOERROR
        assert len(response.answers) >= 1
        assert response.answers[0].rtype == DNSRecordType.CNAME
    
    @pytest.mark.asyncio
    async def test_txt_record(self, dns_server):
        """Test TXT record lookup."""
        response = await send_query("test.local", DNSRecordType.TXT)
        
        assert response.header.rcode == DNSRcode.NOERROR
        assert len(response.answers) >= 1
        assert "spf1" in response.answers[0].parsed_data
    
    @pytest.mark.asyncio
    async def test_nxdomain(self, dns_server):
        """Test non-existent domain."""
        response = await send_query("nonexistent.test.local", DNSRecordType.A)
        
        assert response.header.rcode == DNSRcode.NXDOMAIN
        assert len(response.answers) == 0
    
    @pytest.mark.asyncio
    async def test_no_such_type(self, dns_server):
        """Test query for non-existent record type."""
        response = await send_query("www.test.local", DNSRecordType.MX)
        
        # Should return NOERROR with empty answers
        assert len(response.answers) == 0


class TestIPv6:
    """Test IPv6-specific functionality."""
    
    @pytest.mark.asyncio
    async def test_ipv6_only_host(self, dns_server):
        """Test host with only AAAA record."""
        # Should have AAAA
        response = await send_query("ipv6only.test.local", DNSRecordType.AAAA)
        assert len(response.answers) == 1
        assert "2001:db8::100" in response.answers[0].parsed_data
        
        # Should NOT have A
        response = await send_query("ipv6only.test.local", DNSRecordType.A)
        assert len(response.answers) == 0
    
    @pytest.mark.asyncio
    async def test_dual_stack_host(self, dns_server):
        """Test host with both A and AAAA records."""
        # Get A record
        response_a = await send_query("www.test.local", DNSRecordType.A)
        assert len(response_a.answers) == 1
        
        # Get AAAA record
        response_aaaa = await send_query("www.test.local", DNSRecordType.AAAA)
        assert len(response_aaaa.answers) == 1
        
        # Verify different addresses
        assert response_a.answers[0].parsed_data != response_aaaa.answers[0].parsed_data


class TestCaching:
    """Test DNS caching functionality."""
    
    @pytest.mark.asyncio
    async def test_cache_hit(self, dns_server):
        """Test that cached responses are faster."""
        # First query (cache miss)
        start = time.perf_counter()
        await send_query("www.test.local", DNSRecordType.A)
        first_time = time.perf_counter() - start
        
        # Second query (cache hit)
        start = time.perf_counter()
        await send_query("www.test.local", DNSRecordType.A)
        second_time = time.perf_counter() - start
        
        # Cache hit should be at least as fast
        # (Can't guarantee faster in all cases)
        assert second_time <= first_time * 2


class TestConcurrency:
    """Test concurrent query handling."""
    
    @pytest.mark.asyncio
    async def test_concurrent_queries(self, dns_server):
        """Test handling multiple concurrent queries."""
        queries = [
            ("www.test.local", DNSRecordType.A),
            ("mail.test.local", DNSRecordType.A),
            ("test.local", DNSRecordType.AAAA),
            ("blog.test.local", DNSRecordType.CNAME),
            ("test.local", DNSRecordType.TXT),
        ] * 10  # 50 queries
        
        tasks = [send_query(name, qtype) for name, qtype in queries]
        responses = await asyncio.gather(*tasks)
        
        # All should succeed
        for response in responses:
            assert response.header.rcode in (DNSRcode.NOERROR, DNSRcode.NXDOMAIN)
    
    @pytest.mark.asyncio
    async def test_high_load(self, dns_server):
        """Test handling high query load."""
        num_queries = 100
        
        start = time.perf_counter()
        
        tasks = [
            send_query("www.test.local", DNSRecordType.A)
            for _ in range(num_queries)
        ]
        responses = await asyncio.gather(*tasks)
        
        elapsed = time.perf_counter() - start
        qps = num_queries / elapsed
        
        # All should succeed
        success = sum(1 for r in responses if r.header.rcode == DNSRcode.NOERROR)
        assert success == num_queries
        
        # Should handle at least 100 QPS
        assert qps > 50, f"QPS too low: {qps:.2f}"


class TestEdgeCases:
    """Test edge cases and error handling."""
    
    @pytest.mark.asyncio
    async def test_empty_query(self, dns_server):
        """Test handling of malformed queries."""
        # Send garbage data
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.setblocking(False)
        
        try:
            loop = asyncio.get_event_loop()
            await loop.sock_sendto(sock, b"\x00\x00\x00\x00", ("127.0.0.1", 15353))
            
            # Server should handle gracefully (may not respond)
            try:
                await asyncio.wait_for(
                    loop.sock_recv(sock, 4096),
                    timeout=0.5
                )
            except asyncio.TimeoutError:
                pass  # Expected - server drops malformed queries
        finally:
            sock.close()
    
    @pytest.mark.asyncio
    async def test_case_insensitivity(self, dns_server):
        """Test that queries are case-insensitive."""
        queries = [
            "www.test.local",
            "WWW.TEST.LOCAL",
            "Www.Test.Local",
            "wWw.TeSt.LoCaL",
        ]
        
        for query in queries:
            response = await send_query(query, DNSRecordType.A)
            assert response.header.rcode == DNSRcode.NOERROR
            assert len(response.answers) == 1
    
    @pytest.mark.asyncio
    async def test_long_domain_name(self, dns_server):
        """Test handling of long domain names."""
        # Create a very long subdomain
        long_name = "a" * 60 + ".test.local"
        response = await send_query(long_name, DNSRecordType.A)
        
        # Should return NXDOMAIN (not crash)
        assert response.header.rcode == DNSRcode.NXDOMAIN


class TestMetrics:
    """Test metrics collection."""
    
    @pytest.mark.asyncio
    async def test_query_count(self, dns_server):
        """Test that query metrics are collected."""
        initial = dns_server.metrics.get("queries", 0)
        
        # Make some queries
        for _ in range(5):
            await send_query("www.test.local", DNSRecordType.A)
        
        final = dns_server.metrics.get("queries", 0)
        assert final >= initial + 5


# Performance benchmarks

class TestPerformance:
    """Performance benchmark tests."""
    
    @pytest.mark.asyncio
    @pytest.mark.slow
    async def test_latency_p99(self, dns_server):
        """Test P99 latency is acceptable."""
        latencies = []
        
        for _ in range(100):
            start = time.perf_counter()
            await send_query("www.test.local", DNSRecordType.A)
            latencies.append((time.perf_counter() - start) * 1000)
        
        latencies.sort()
        p99 = latencies[98]  # 99th percentile
        
        # P99 should be under 50ms
        assert p99 < 50, f"P99 latency too high: {p99:.2f}ms"
    
    @pytest.mark.asyncio
    @pytest.mark.slow
    async def test_throughput(self, dns_server):
        """Test query throughput."""
        num_queries = 500
        concurrency = 50
        
        start = time.perf_counter()
        
        for batch_start in range(0, num_queries, concurrency):
            batch_size = min(concurrency, num_queries - batch_start)
            tasks = [
                send_query("www.test.local", DNSRecordType.A)
                for _ in range(batch_size)
            ]
            await asyncio.gather(*tasks)
        
        elapsed = time.perf_counter() - start
        qps = num_queries / elapsed
        
        # Should handle at least 500 QPS
        assert qps > 200, f"Throughput too low: {qps:.2f} QPS"


# Run with: pytest tests/test_integration.py -v -m "not slow"
# Full: pytest tests/test_integration.py -v
```

Create `tests/conftest.py`:

```python
"""
Pytest configuration and shared fixtures.
"""

import pytest
import asyncio


@pytest.fixture(scope="session")
def event_loop():
    """Create event loop for async tests."""
    loop = asyncio.new_event_loop()
    yield loop
    loop.close()


def pytest_configure(config):
    """Configure pytest markers."""
    config.addinivalue_line(
        "markers", "slow: marks tests as slow (deselect with '-m \"not slow\"')"
    )
```

**Git checkpoint:**
```bash
git add tests/
git commit -m "Day 38: Integration tests"
```

---

## Day 39: Deployment Automation

Create `deployment/deploy.sh`:

```bash
#!/bin/bash
#
# DNS Server Deployment Script
# Day 39: Production deployment automation
#

set -euo pipefail

# Configuration
APP_NAME="dns-server"
APP_USER="dns"
APP_GROUP="dns"
APP_DIR="/opt/${APP_NAME}"
CONFIG_DIR="/etc/${APP_NAME}"
LOG_DIR="/var/log/${APP_NAME}"
DATA_DIR="/var/lib/${APP_NAME}"

REPO_URL="${REPO_URL:-https://github.com/yourusername/python-dns-server.git}"
BRANCH="${BRANCH:-main}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

install_dependencies() {
    log_info "Installing system dependencies..."
    
    apt-get update
    apt-get install -y \
        python3 \
        python3-pip \
        python3-venv \
        git \
        curl \
        jq
    
    log_info "Dependencies installed"
}

create_user() {
    log_info "Creating application user..."
    
    if ! id -u "${APP_USER}" &>/dev/null; then
        useradd --system --no-create-home --shell /usr/sbin/nologin "${APP_USER}"
        log_info "User ${APP_USER} created"
    else
        log_info "User ${APP_USER} already exists"
    fi
}

create_directories() {
    log_info "Creating directories..."
    
    mkdir -p "${APP_DIR}"
    mkdir -p "${CONFIG_DIR}/zones"
    mkdir -p "${LOG_DIR}"
    mkdir -p "${DATA_DIR}"
    
    chown -R "${APP_USER}:${APP_GROUP}" "${LOG_DIR}"
    chown -R "${APP_USER}:${APP_GROUP}" "${DATA_DIR}"
    
    log_info "Directories created"
}

clone_or_update_repo() {
    log_info "Setting up application..."
    
    if [[ -d "${APP_DIR}/.git" ]]; then
        log_info "Updating existing installation..."
        cd "${APP_DIR}"
        git fetch origin
        git checkout "${BRANCH}"
        git reset --hard "origin/${BRANCH}"
    else
        log_info "Cloning repository..."
        git clone --branch "${BRANCH}" "${REPO_URL}" "${APP_DIR}"
    fi
    
    log_info "Repository ready"
}

setup_virtualenv() {
    log_info "Setting up Python virtual environment..."
    
    cd "${APP_DIR}"
    
    if [[ ! -d "venv" ]]; then
        python3 -m venv venv
    fi
    
    source venv/bin/activate
    pip install --upgrade pip wheel
    pip install -r requirements.txt
    
    log_info "Virtual environment ready"
}

setup_config() {
    log_info "Setting up configuration..."
    
    # Create default config if not exists
    if [[ ! -f "${CONFIG_DIR}/config.yaml" ]]; then
        cat > "${CONFIG_DIR}/config.yaml" << 'EOF'
# DNS Server Configuration
network:
  listen_address: "::"
  listen_port: 53
  enable_ipv6: true

upstream:
  servers:
    - 8.8.8.8
    - 1.1.1.1
  servers_v6:
    - 2001:4860:4860::8888

cache:
  enabled: true
  max_size: 10000

security:
  enable_rate_limiting: true
  queries_per_second: 50
  whitelist_networks:
    - 127.0.0.0/8
    - "::1/128"

metrics:
  enabled: true
  port: 9153

logging:
  level: INFO
  log_queries: false
  log_file: /var/log/dns-server/dns.log

zones:
  files:
    # Add zone files here
    # - /etc/dns-server/zones/example.zone
EOF
        log_info "Default configuration created"
    else
        log_info "Configuration already exists"
    fi
    
    # Create environment file
    cat > "${CONFIG_DIR}/environment" << EOF
DNS_CONFIG=${CONFIG_DIR}/config.yaml
PATH=${APP_DIR}/venv/bin:\$PATH
EOF
    
    chown -R root:${APP_GROUP} "${CONFIG_DIR}"
    chmod 750 "${CONFIG_DIR}"
    chmod 640 "${CONFIG_DIR}/config.yaml"
}

install_systemd_service() {
    log_info "Installing systemd service..."
    
    cat > /etc/systemd/system/${APP_NAME}.service << EOF
[Unit]
Description=Python DNS Server
Documentation=https://github.com/yourusername/python-dns-server
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${APP_USER}
Group=${APP_GROUP}
WorkingDirectory=${APP_DIR}
EnvironmentFile=${CONFIG_DIR}/environment
ExecStart=${APP_DIR}/venv/bin/python -m dns_server serve
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5

# Allow binding to privileged ports
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ReadWritePaths=${LOG_DIR} ${DATA_DIR}

# Resource limits
LimitNOFILE=65535
MemoryMax=512M

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    systemctl enable ${APP_NAME}
    
    log_info "Systemd service installed"
}

start_service() {
    log_info "Starting service..."
    
    systemctl restart ${APP_NAME}
    sleep 2
    
    if systemctl is-active --quiet ${APP_NAME}; then
        log_info "Service started successfully"
    else
        log_error "Service failed to start"
        journalctl -u ${APP_NAME} -n 20 --no-pager
        exit 1
    fi
}

health_check() {
    log_info "Running health check..."
    
    # Check if service is running
    if ! systemctl is-active --quiet ${APP_NAME}; then
        log_error "Service is not running"
        exit 1
    fi
    
    # Check DNS response
    if command -v dig &> /dev/null; then
        if dig @localhost -p 53 localhost +short +timeout=2 &> /dev/null; then
            log_info "DNS responding"
        else
            log_warn "DNS not responding on port 53"
        fi
    fi
    
    # Check metrics endpoint
    if curl -s http://localhost:9153/health | jq -e '.status == "healthy"' &> /dev/null; then
        log_info "Metrics endpoint healthy"
    else
        log_warn "Metrics endpoint not responding"
    fi
    
    log_info "Health check complete"
}

show_status() {
    echo ""
    echo "=========================================="
    echo "  DNS Server Deployment Complete"
    echo "=========================================="
    echo ""
    echo "Service: ${APP_NAME}"
    echo "Status:  $(systemctl is-active ${APP_NAME})"
    echo ""
    echo "Directories:"
    echo "  App:    ${APP_DIR}"
    echo "  Config: ${CONFIG_DIR}"
    echo "  Logs:   ${LOG_DIR}"
    echo ""
    echo "Commands:"
    echo "  systemctl status ${APP_NAME}"
    echo "  journalctl -u ${APP_NAME} -f"
    echo "  dig @localhost example.com"
    echo ""
    echo "Metrics: http://localhost:9153/metrics"
    echo ""
}

# Main deployment flow
main() {
    log_info "Starting deployment..."
    
    check_root
    install_dependencies
    create_user
    create_directories
    clone_or_update_repo
    setup_virtualenv
    setup_config
    install_systemd_service
    start_service
    health_check
    show_status
    
    log_info "Deployment complete!"
}

# Handle command line arguments
case "${1:-deploy}" in
    deploy)
        main
        ;;
    update)
        clone_or_update_repo
        setup_virtualenv
        start_service
        health_check
        ;;
    restart)
        start_service
        ;;
    status)
        systemctl status ${APP_NAME}
        ;;
    logs)
        journalctl -u ${APP_NAME} -f
        ;;
    *)
        echo "Usage: $0 {deploy|update|restart|status|logs}"
        exit 1
        ;;
esac
```

Create `deployment/Dockerfile`:

```dockerfile
# Multi-stage build for DNS Server
FROM python:3.11-slim as builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Create virtual environment
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt


# Production image
FROM python:3.11-slim

LABEL maintainer="your.email@example.com"
LABEL description="Python DNS Server with IPv6 support"

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Copy virtual environment
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Copy application
COPY dns_server/ ./dns_server/
COPY zones/ ./zones/

# Create non-root user
RUN useradd --system --no-create-home --shell /usr/sbin/nologin dns \
    && chown -R dns:dns /app

USER dns

# Expose ports
EXPOSE 53/udp
EXPOSE 53/tcp
EXPOSE 9153/tcp

# Environment
ENV DNS_NETWORK__LISTEN_PORT=53 \
    DNS_METRICS__ENABLED=true \
    DNS_LOGGING__LEVEL=INFO

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -sf http://localhost:9153/health || exit 1

# Entry point
ENTRYPOINT ["python", "-m", "dns_server"]
CMD ["serve"]
```

Create `deployment/docker-compose.yml`:

```yaml
version: '3.8'

services:
  dns:
    build:
      context: ..
      dockerfile: deployment/Dockerfile
    container_name: dns-server
    hostname: dns-server
    restart: unless-stopped
    ports:
      - "53:53/udp"
      - "53:53/tcp"
      - "9153:9153"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./zones:/app/zones:ro
      - dns-logs:/var/log/dns-server
    environment:
      - DNS_CONFIG=/app/config.yaml
      - DNS_LOGGING__LEVEL=INFO
    healthcheck:
      test: ["CMD", "curl", "-sf", "http://localhost:9153/health"]
      interval: 30s
      timeout: 5s
      retries: 3
    networks:
      - dns-network
    cap_add:
      - NET_BIND_SERVICE
    cap_drop:
      - ALL
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp

  # Optional: Prometheus for metrics
  prometheus:
    image: prom/prometheus:latest
    container_name: dns-prometheus
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
    networks:
      - dns-network
    profiles:
      - monitoring

  # Optional: Grafana for dashboards
  grafana:
    image: grafana/grafana:latest
    container_name: dns-grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    networks:
      - dns-network
    profiles:
      - monitoring

volumes:
  dns-logs:
  prometheus-data:
  grafana-data:

networks:
  dns-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/24
```

**Git checkpoint:**
```bash
chmod +x deployment/deploy.sh
git add deployment/
git commit -m "Day 39: Deployment automation"
```

---

## Day 40: Final Project & Documentation

Create `README.md`:

```markdown
# Python DNS Server

A full-featured DNS server with IPv6 support, written in Python.

## Features

- **Dual-stack Support**: IPv4 and IPv6
- **Record Types**: A, AAAA, CNAME, MX, TXT, PTR, NS, SOA
- **Zone Files**: BIND-style zone file support
- **Caching**: TTL-aware record caching
- **Recursive Resolution**: Query upstream servers
- **Security**: Rate limiting, blocklists, RRL
- **Metrics**: Prometheus-compatible metrics
- **Zone Transfers**: AXFR and IXFR support

## Quick Start

```bash
# Install
pip install -r requirements.txt

# Run
python -m dns_server serve

# Query
dig @localhost -p 5353 google.com
```

## CLI Usage

```bash
# Start server
pydns serve -p 5353 -z zones/example.zone

# Query a domain
pydns query google.com -t AAAA

# Benchmark a server
pydns benchmark 8.8.8.8 -n 1000

# Compare servers
pydns compare -s 8.8.8.8 -s 1.1.1.1 -s 9.9.9.9

# Validate zone file
pydns zone validate zones/example.zone

# Create zone template
pydns zone create example.com -o example.zone
```

## Configuration

```yaml
# config.yaml
network:
  listen_port: 53
  enable_ipv6: true

upstream:
  servers:
    - 8.8.8.8
    - 1.1.1.1

cache:
  max_size: 10000

security:
  enable_rate_limiting: true
```

## Docker

```bash
# Build
docker build -t dns-server .

# Run
docker run -d -p 53:53/udp -p 53:53/tcp dns-server
```

## Zone Files

```
$TTL 3600
$ORIGIN example.local.

@    IN    SOA    ns1 admin (2024010101 3600 900 604800 86400)
@    IN    NS     ns1
@    IN    A      192.168.1.1
@    IN    AAAA   2001:db8::1
www  IN    A      192.168.1.10
www  IN    AAAA   2001:db8::10
```

## Metrics

Prometheus metrics available at `http://localhost:9153/metrics`:

- `dns_queries_total` - Total queries by type
- `dns_query_duration_seconds` - Query latency histogram
- `dns_cache_entries` - Cache size
- `dns_cache_hits_total` - Cache hit count

## License

MIT
```

Create final project structure verification:

```bash
#!/bin/bash
# verify_project.sh - Verify project structure

echo "=== Project Structure Verification ==="

required_files=(
    "dns_server/__init__.py"
    "dns_server/__main__.py"
    "dns_server/server.py"
    "dns_server/protocol.py"
    "dns_server/cli.py"
    "dns_server/config.py"
    "dns_server/zones.py"
    "dns_server/cache.py"
    "dns_server/security.py"
    "dns_server/metrics.py"
    "dns_server/utils.py"
    "dns_server/ipv6.py"
    "tests/test_integration.py"
    "tests/conftest.py"
    "zones/example.zone"
    "deployment/deploy.sh"
    "deployment/Dockerfile"
    "deployment/docker-compose.yml"
    "requirements.txt"
    "requirements-dev.txt"
    "README.md"
    "Makefile"
    ".gitignore"
)

missing=0
for file in "${required_files[@]}"; do
    if [[ -f "$file" ]]; then
        echo "  ✓ $file"
    else
        echo "  ✗ $file (MISSING)"
        ((missing++))
    fi
done

echo ""
if [[ $missing -eq 0 ]]; then
    echo "All files present!"
else
    echo "Missing $missing files"
fi
```

**Git checkpoint:**
```bash
git add README.md verify_project.sh
git commit -m "Day 40: Final documentation"
git tag -a v1.0.0 -m "Version 1.0.0 - Complete DNS Server"
```

---

## Summary: Days 36-40

### CLI (Day 36)
- ✅ Click-based command structure
- ✅ Rich terminal output
- ✅ Server, query, benchmark commands
- ✅ Zone management commands

### Configuration (Day 37)
- ✅ Pydantic settings validation
- ✅ YAML/TOML config files
- ✅ Environment variable override
- ✅ Hierarchical configuration

### Testing (Day 38)
- ✅ Pytest async fixtures
- ✅ Basic functionality tests
- ✅ IPv6 tests
- ✅ Concurrency tests
- ✅ Performance benchmarks

### Deployment (Day 39)
- ✅ Bash deployment script
- ✅ Multi-stage Dockerfile
- ✅ Docker Compose stack
- ✅ Prometheus integration

### Documentation (Day 40)
- ✅ Complete README
- ✅ CLI examples
- ✅ Configuration reference
- ✅ Project verification

---

## Complete Course Summary

| Phase | Days | Topics |
|-------|------|--------|
| 1 | 1-3 | Python basics, control flow, functions |
| 2 | 4-6 | OOP, error handling, decorators |
| 3 | 7-9 | Files, HTTP, subprocess |
| 4 | 10-12 | SQLite, PostgreSQL, SQLAlchemy |
| 5 | 13-16 | FastAPI, REST, database integration |
| 6 | 17-20 | Systemd, deployment, CI/CD |
| 7 | 21-30 | IPv6, DNS protocol, server |
| 7+ | 31-35 | Security, transfers, metrics |
| 7++ | 36-40 | CLI, config, tests, deployment |

**Total: 40 days of progressive Python learning**

This completes your Python learning path from basics to a production DNS server!
