# Python Learning Path - Additional Exercises

Supplementary exercises for Phase 2-4, including OOP deep dive, async programming, testing, and more network automation examples.

---

## Phase 2 Extended: Advanced OOP & Patterns

### Day 4 Extended: Complete Network Device Hierarchy

Create `day_04/network_models.py`:

```python
#!/usr/bin/env python3
"""
Day 4 Extended: Complete OOP Network Device Models
Learning: inheritance, composition, abstract classes, protocols
"""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import List, Optional, Dict, Protocol
from datetime import datetime
from enum import Enum
import json


class DeviceStatus(Enum):
    """Device operational status."""
    UP = "up"
    DOWN = "down"
    MAINTENANCE = "maintenance"
    UNKNOWN = "unknown"


class InterfaceStatus(Enum):
    """Interface operational status."""
    UP = "up"
    DOWN = "down"
    ADMIN_DOWN = "admin_down"
    ERROR = "error"


# Protocol for type checking (structural subtyping)
class Configurable(Protocol):
    """Protocol for devices that can be configured."""
    def get_config(self) -> str: ...
    def apply_config(self, config: str) -> bool: ...


@dataclass
class Interface:
    """Network interface data class."""
    name: str
    ip_address: Optional[str] = None
    netmask: str = "255.255.255.0"
    mac_address: Optional[str] = None
    status: InterfaceStatus = InterfaceStatus.DOWN
    speed_mbps: int = 1000
    mtu: int = 1500
    vlan: Optional[int] = None
    description: str = ""
    
    @property
    def is_up(self) -> bool:
        return self.status == InterfaceStatus.UP
    
    def to_dict(self) -> dict:
        return {
            "name": self.name,
            "ip_address": self.ip_address,
            "netmask": self.netmask,
            "mac_address": self.mac_address,
            "status": self.status.value,
            "speed_mbps": self.speed_mbps,
            "mtu": self.mtu,
            "vlan": self.vlan,
            "description": self.description,
        }


class NetworkDevice(ABC):
    """
    Abstract base class for all network devices.
    
    Demonstrates: ABC, abstract methods, class variables, properties
    """
    
    # Class variable - shared across all instances
    device_count = 0
    _registry: Dict[str, 'NetworkDevice'] = {}
    
    def __init__(self, hostname: str, ip_address: str, vendor: str = "generic"):
        self.hostname = hostname
        self._ip_address = ip_address
        self.vendor = vendor
        self._status = DeviceStatus.UNKNOWN
        self.interfaces: Dict[str, Interface] = {}
        self.created_at = datetime.now()
        self.updated_at = datetime.now()
        self._config: List[str] = []
        
        # Track instances
        NetworkDevice.device_count += 1
        NetworkDevice._registry[hostname] = self
    
    @property
    def ip_address(self) -> str:
        """Management IP address (read-only)."""
        return self._ip_address
    
    @property
    def status(self) -> DeviceStatus:
        return self._status
    
    @status.setter
    def status(self, value: DeviceStatus) -> None:
        if not isinstance(value, DeviceStatus):
            raise ValueError(f"Status must be DeviceStatus enum")
        self._status = value
        self.updated_at = datetime.now()
    
    @abstractmethod
    def get_device_type(self) -> str:
        """Return device type string."""
        pass
    
    @abstractmethod
    def show_version(self) -> str:
        """Return version information."""
        pass
    
    def add_interface(self, interface: Interface) -> None:
        """Add an interface to the device."""
        self.interfaces[interface.name] = interface
        self.updated_at = datetime.now()
    
    def get_interface(self, name: str) -> Optional[Interface]:
        """Get interface by name."""
        return self.interfaces.get(name)
    
    def get_up_interfaces(self) -> List[Interface]:
        """Get all interfaces that are up."""
        return [i for i in self.interfaces.values() if i.is_up]
    
    def ping(self) -> bool:
        """Check if device is reachable."""
        return self._status == DeviceStatus.UP
    
    def get_config(self) -> str:
        """Get current configuration."""
        return "\n".join(self._config)
    
    def apply_config(self, config: str) -> bool:
        """Apply configuration lines."""
        self._config.extend(config.strip().split("\n"))
        self.updated_at = datetime.now()
        return True
    
    def to_dict(self) -> dict:
        """Serialize device to dictionary."""
        return {
            "hostname": self.hostname,
            "ip_address": self._ip_address,
            "vendor": self.vendor,
            "device_type": self.get_device_type(),
            "status": self._status.value,
            "interfaces": {
                name: iface.to_dict() 
                for name, iface in self.interfaces.items()
            },
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
        }
    
    def to_json(self) -> str:
        """Serialize to JSON string."""
        return json.dumps(self.to_dict(), indent=2)
    
    @classmethod
    def get_device_count(cls) -> int:
        """Get total number of devices created."""
        return cls.device_count
    
    @classmethod
    def get_device(cls, hostname: str) -> Optional['NetworkDevice']:
        """Get device by hostname from registry."""
        return cls._registry.get(hostname)
    
    @classmethod
    def get_all_devices(cls) -> List['NetworkDevice']:
        """Get all registered devices."""
        return list(cls._registry.values())
    
    @staticmethod
    def is_valid_hostname(hostname: str) -> bool:
        """Validate hostname format."""
        if not hostname or len(hostname) > 253:
            return False
        return all(c.isalnum() or c in "-_." for c in hostname)
    
    def __str__(self) -> str:
        return f"{self.hostname} ({self._ip_address}) - {self._status.value}"
    
    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(hostname='{self.hostname}', ip='{self._ip_address}')"


class Router(NetworkDevice):
    """Router with routing table and BGP support."""
    
    def __init__(self, hostname: str, ip_address: str, vendor: str = "generic"):
        super().__init__(hostname, ip_address, vendor)
        self.routing_table: List[Dict] = []
        self.bgp_neighbors: Dict[str, Dict] = {}
        self.ospf_areas: List[int] = []
    
    def get_device_type(self) -> str:
        return "router"
    
    def show_version(self) -> str:
        return f"{self.vendor} Router - {self.hostname}"
    
    def add_route(self, network: str, next_hop: str, 
                  metric: int = 100, protocol: str = "static") -> None:
        """Add a route to the routing table."""
        self.routing_table.append({
            "network": network,
            "next_hop": next_hop,
            "metric": metric,
            "protocol": protocol,
        })
    
    def add_bgp_neighbor(self, peer_ip: str, remote_as: int, 
                         description: str = "") -> None:
        """Add a BGP neighbor."""
        self.bgp_neighbors[peer_ip] = {
            "remote_as": remote_as,
            "description": description,
            "state": "idle",
        }
    
    def get_routes_to(self, destination: str) -> List[Dict]:
        """Find routes that could reach a destination."""
        # Simplified - in real code, would do prefix matching
        return [r for r in self.routing_table 
                if destination.startswith(r["network"].split("/")[0].rsplit(".", 1)[0])]
    
    def show_ip_route(self) -> str:
        """Display routing table."""
        lines = [f"Routing table for {self.hostname}:", "-" * 60]
        for route in self.routing_table:
            lines.append(
                f"  {route['network']:20} via {route['next_hop']:15} "
                f"[{route['protocol']}/{route['metric']}]"
            )
        return "\n".join(lines)
    
    def show_bgp_summary(self) -> str:
        """Display BGP neighbor summary."""
        lines = [f"BGP Summary for {self.hostname}:", "-" * 60]
        for peer, info in self.bgp_neighbors.items():
            lines.append(f"  {peer:15} AS{info['remote_as']:6} {info['state']}")
        return "\n".join(lines)


class Switch(NetworkDevice):
    """Layer 2/3 Switch with VLAN and MAC table support."""
    
    def __init__(self, hostname: str, ip_address: str, vendor: str = "generic"):
        super().__init__(hostname, ip_address, vendor)
        self.vlans: Dict[int, str] = {1: "default"}
        self.mac_table: Dict[str, Dict] = {}
        self.spanning_tree_enabled = True
    
    def get_device_type(self) -> str:
        return "switch"
    
    def show_version(self) -> str:
        return f"{self.vendor} Switch - {self.hostname}"
    
    def add_vlan(self, vlan_id: int, name: str) -> None:
        """Add a VLAN."""
        if not 1 <= vlan_id <= 4094:
            raise ValueError(f"Invalid VLAN ID: {vlan_id}")
        self.vlans[vlan_id] = name
    
    def delete_vlan(self, vlan_id: int) -> bool:
        """Delete a VLAN."""
        if vlan_id == 1:
            raise ValueError("Cannot delete default VLAN")
        if vlan_id in self.vlans:
            del self.vlans[vlan_id]
            return True
        return False
    
    def learn_mac(self, mac: str, interface: str, vlan: int = 1) -> None:
        """Learn a MAC address on an interface."""
        self.mac_table[mac.upper()] = {
            "interface": interface,
            "vlan": vlan,
            "learned_at": datetime.now().isoformat(),
        }
    
    def lookup_mac(self, mac: str) -> Optional[Dict]:
        """Look up a MAC address."""
        return self.mac_table.get(mac.upper())
    
    def get_macs_on_interface(self, interface: str) -> List[str]:
        """Get all MACs learned on an interface."""
        return [
            mac for mac, info in self.mac_table.items()
            if info["interface"] == interface
        ]
    
    def show_vlan(self) -> str:
        """Display VLAN information."""
        lines = [f"VLAN table for {self.hostname}:", "-" * 40]
        for vlan_id, name in sorted(self.vlans.items()):
            lines.append(f"  VLAN {vlan_id:4}: {name}")
        return "\n".join(lines)
    
    def show_mac_table(self) -> str:
        """Display MAC address table."""
        lines = [f"MAC table for {self.hostname}:", "-" * 60]
        for mac, info in self.mac_table.items():
            lines.append(f"  {mac}  VLAN {info['vlan']:4}  {info['interface']}")
        return "\n".join(lines)


class Firewall(NetworkDevice):
    """Firewall with security rules."""
    
    def __init__(self, hostname: str, ip_address: str, vendor: str = "generic"):
        super().__init__(hostname, ip_address, vendor)
        self.rules: List[Dict] = []
        self.zones: Dict[str, List[str]] = {}
        self.nat_rules: List[Dict] = []
    
    def get_device_type(self) -> str:
        return "firewall"
    
    def show_version(self) -> str:
        return f"{self.vendor} Firewall - {self.hostname}"
    
    def add_zone(self, zone_name: str, interfaces: List[str]) -> None:
        """Add a security zone."""
        self.zones[zone_name] = interfaces
    
    def add_rule(self, name: str, source: str, destination: str,
                 service: str, action: str = "permit") -> None:
        """Add a firewall rule."""
        self.rules.append({
            "name": name,
            "source": source,
            "destination": destination,
            "service": service,
            "action": action,
            "enabled": True,
        })
    
    def add_nat_rule(self, name: str, source: str, 
                     translated: str, nat_type: str = "source") -> None:
        """Add a NAT rule."""
        self.nat_rules.append({
            "name": name,
            "original": source,
            "translated": translated,
            "type": nat_type,
        })
    
    def check_traffic(self, source: str, destination: str, 
                      service: str) -> str:
        """Check if traffic would be permitted."""
        for rule in self.rules:
            if not rule["enabled"]:
                continue
            # Simplified matching
            if (rule["source"] in ["any", source] and
                rule["destination"] in ["any", destination] and
                rule["service"] in ["any", service]):
                return rule["action"]
        return "deny"  # Implicit deny
    
    def show_rules(self) -> str:
        """Display firewall rules."""
        lines = [f"Firewall rules for {self.hostname}:", "-" * 70]
        for i, rule in enumerate(self.rules, 1):
            status = "✓" if rule["enabled"] else "✗"
            lines.append(
                f"  {i:3}. [{status}] {rule['name']:20} "
                f"{rule['source']:15} -> {rule['destination']:15} "
                f"{rule['service']:10} {rule['action'].upper()}"
            )
        return "\n".join(lines)


class LoadBalancer(NetworkDevice):
    """Load balancer with virtual servers and pools."""
    
    def __init__(self, hostname: str, ip_address: str, vendor: str = "generic"):
        super().__init__(hostname, ip_address, vendor)
        self.pools: Dict[str, List[Dict]] = {}
        self.virtual_servers: Dict[str, Dict] = {}
    
    def get_device_type(self) -> str:
        return "load_balancer"
    
    def show_version(self) -> str:
        return f"{self.vendor} Load Balancer - {self.hostname}"
    
    def add_pool(self, name: str) -> None:
        """Create a server pool."""
        self.pools[name] = []
    
    def add_pool_member(self, pool_name: str, ip: str, port: int,
                        weight: int = 1) -> None:
        """Add a member to a pool."""
        if pool_name not in self.pools:
            self.add_pool(pool_name)
        self.pools[pool_name].append({
            "ip": ip,
            "port": port,
            "weight": weight,
            "status": "up",
        })
    
    def add_virtual_server(self, name: str, vip: str, port: int,
                           pool: str, method: str = "round_robin") -> None:
        """Add a virtual server."""
        self.virtual_servers[name] = {
            "vip": vip,
            "port": port,
            "pool": pool,
            "method": method,
            "enabled": True,
        }
    
    def show_pools(self) -> str:
        """Display pool information."""
        lines = [f"Pools for {self.hostname}:", "-" * 50]
        for pool_name, members in self.pools.items():
            lines.append(f"\n  Pool: {pool_name}")
            for m in members:
                lines.append(f"    {m['ip']}:{m['port']} (weight:{m['weight']}) [{m['status']}]")
        return "\n".join(lines)


# Factory function
def create_device(device_type: str, hostname: str, ip_address: str,
                  vendor: str = "generic") -> NetworkDevice:
    """Factory function to create network devices."""
    device_classes = {
        "router": Router,
        "switch": Switch,
        "firewall": Firewall,
        "load_balancer": LoadBalancer,
    }
    
    device_class = device_classes.get(device_type.lower())
    if not device_class:
        raise ValueError(f"Unknown device type: {device_type}")
    
    return device_class(hostname, ip_address, vendor)


def main():
    print("=== Network Device OOP Demo ===\n")
    
    # Create devices using factory
    router = create_device("router", "core-rtr-01", "10.0.0.1", "juniper")
    switch = create_device("switch", "access-sw-01", "10.0.1.1", "arista")
    firewall = create_device("firewall", "edge-fw-01", "10.0.0.254", "paloalto")
    lb = create_device("load_balancer", "prod-lb-01", "10.0.0.100", "f5")
    
    # Set all devices to UP
    for device in [router, switch, firewall, lb]:
        device.status = DeviceStatus.UP
    
    # Configure router
    router.add_interface(Interface("ge-0/0/0", "10.0.1.1", speed_mbps=10000))
    router.add_interface(Interface("ge-0/0/1", "10.0.2.1", speed_mbps=10000))
    router.add_route("0.0.0.0/0", "10.0.0.254", protocol="static")
    router.add_route("192.168.0.0/16", "10.0.1.254", protocol="bgp")
    router.add_bgp_neighbor("10.0.0.2", 65001, "ISP-1")
    router.add_bgp_neighbor("10.0.0.3", 65002, "ISP-2")
    
    # Configure switch
    switch.add_vlan(100, "Servers")
    switch.add_vlan(200, "Users")
    switch.add_vlan(999, "Management")
    switch.learn_mac("aa:bb:cc:dd:ee:01", "Ethernet1", 100)
    switch.learn_mac("aa:bb:cc:dd:ee:02", "Ethernet2", 100)
    switch.learn_mac("aa:bb:cc:dd:ee:03", "Ethernet3", 200)
    
    # Configure firewall
    firewall.add_zone("trust", ["eth1", "eth2"])
    firewall.add_zone("untrust", ["eth0"])
    firewall.add_zone("dmz", ["eth3"])
    firewall.add_rule("allow-web", "trust", "dmz", "http", "permit")
    firewall.add_rule("allow-https", "trust", "dmz", "https", "permit")
    firewall.add_rule("allow-ssh-mgmt", "trust", "any", "ssh", "permit")
    firewall.add_nat_rule("outbound-nat", "10.0.0.0/8", "203.0.113.1", "source")
    
    # Configure load balancer
    lb.add_pool("web-servers")
    lb.add_pool_member("web-servers", "10.0.1.10", 80)
    lb.add_pool_member("web-servers", "10.0.1.11", 80)
    lb.add_pool_member("web-servers", "10.0.1.12", 80)
    lb.add_virtual_server("www", "203.0.113.10", 80, "web-servers")
    
    # Display information
    print(router.show_ip_route())
    print()
    print(router.show_bgp_summary())
    print()
    print(switch.show_vlan())
    print()
    print(switch.show_mac_table())
    print()
    print(firewall.show_rules())
    print()
    print(lb.show_pools())
    
    # Registry demo
    print(f"\n=== Device Registry ===")
    print(f"Total devices: {NetworkDevice.get_device_count()}")
    for device in NetworkDevice.get_all_devices():
        print(f"  {device}")
    
    # Get specific device
    print(f"\nLookup 'core-rtr-01': {NetworkDevice.get_device('core-rtr-01')}")
    
    # JSON serialization
    print(f"\n=== JSON Export ===")
    print(router.to_json()[:500] + "...")


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add day_04/
git commit -m "Day 4 Extended: Complete network device OOP hierarchy"
```

---

## Day 5 Extended: Async Programming

Create `day_05/async_network.py`:

```python
#!/usr/bin/env python3
"""
Day 5 Extended: Asynchronous Programming
Learning: async/await, asyncio, concurrent operations

Install: pip install aiohttp
"""

import asyncio
import time
from typing import List, Dict, Optional
from dataclasses import dataclass
import random


@dataclass
class PingResult:
    """Result of an async ping."""
    host: str
    success: bool
    latency_ms: float
    error: Optional[str] = None


async def async_ping(host: str, timeout: float = 2.0) -> PingResult:
    """
    Simulate an async ping to a host.
    
    In real code, you'd use aioping or subprocess.
    """
    start = time.perf_counter()
    
    # Simulate network latency
    latency = random.uniform(0.01, 0.5)
    await asyncio.sleep(latency)
    
    # Simulate some failures
    if random.random() < 0.1:  # 10% failure rate
        return PingResult(
            host=host,
            success=False,
            latency_ms=0,
            error="Request timed out"
        )
    
    elapsed = (time.perf_counter() - start) * 1000
    return PingResult(host=host, success=True, latency_ms=round(elapsed, 2))


async def fetch_device_info(host: str) -> Dict:
    """Simulate fetching device information."""
    # Simulate API call
    await asyncio.sleep(random.uniform(0.1, 0.3))
    
    return {
        "host": host,
        "hostname": f"device-{host.split('.')[-1]}",
        "uptime": random.randint(1000, 100000),
        "cpu_usage": random.randint(5, 95),
        "memory_usage": random.randint(20, 80),
    }


async def check_port(host: str, port: int, timeout: float = 1.0) -> Dict:
    """Check if a port is open (simulated)."""
    await asyncio.sleep(random.uniform(0.05, 0.2))
    
    # Simulate: common ports more likely to be open
    common_ports = {22, 80, 443, 8080, 8443}
    is_open = port in common_ports or random.random() < 0.3
    
    return {
        "host": host,
        "port": port,
        "open": is_open,
    }


async def scan_host(host: str, ports: List[int]) -> Dict:
    """Scan multiple ports on a host concurrently."""
    tasks = [check_port(host, port) for port in ports]
    results = await asyncio.gather(*tasks)
    
    open_ports = [r["port"] for r in results if r["open"]]
    
    return {
        "host": host,
        "scanned_ports": len(ports),
        "open_ports": open_ports,
    }


async def ping_sweep(network_prefix: str, start: int, end: int) -> List[PingResult]:
    """Ping multiple hosts concurrently."""
    hosts = [f"{network_prefix}.{i}" for i in range(start, end + 1)]
    tasks = [async_ping(host) for host in hosts]
    results = await asyncio.gather(*tasks)
    return results


async def collect_device_metrics(hosts: List[str]) -> List[Dict]:
    """Collect metrics from multiple devices concurrently."""
    tasks = [fetch_device_info(host) for host in hosts]
    results = await asyncio.gather(*tasks, return_exceptions=True)
    
    # Filter out exceptions
    valid_results = [r for r in results if isinstance(r, dict)]
    return valid_results


async def rate_limited_requests(hosts: List[str], 
                                max_concurrent: int = 5) -> List[Dict]:
    """
    Make requests with rate limiting using semaphore.
    """
    semaphore = asyncio.Semaphore(max_concurrent)
    
    async def limited_fetch(host: str) -> Dict:
        async with semaphore:
            return await fetch_device_info(host)
    
    tasks = [limited_fetch(host) for host in hosts]
    return await asyncio.gather(*tasks)


async def with_timeout_example(host: str) -> Optional[Dict]:
    """Demonstrate timeout handling."""
    try:
        result = await asyncio.wait_for(
            fetch_device_info(host),
            timeout=0.1  # Very short timeout
        )
        return result
    except asyncio.TimeoutError:
        print(f"Timeout fetching {host}")
        return None


class AsyncDevicePoller:
    """
    Async device poller with background tasks.
    
    Demonstrates: long-running async tasks, cancellation
    """
    
    def __init__(self, hosts: List[str], interval: float = 5.0):
        self.hosts = hosts
        self.interval = interval
        self.metrics: Dict[str, Dict] = {}
        self._running = False
        self._task: Optional[asyncio.Task] = None
    
    async def poll_once(self) -> None:
        """Poll all devices once."""
        results = await collect_device_metrics(self.hosts)
        for result in results:
            self.metrics[result["host"]] = result
    
    async def _poll_loop(self) -> None:
        """Background polling loop."""
        while self._running:
            await self.poll_once()
            print(f"Polled {len(self.metrics)} devices")
            await asyncio.sleep(self.interval)
    
    async def start(self) -> None:
        """Start background polling."""
        self._running = True
        self._task = asyncio.create_task(self._poll_loop())
    
    async def stop(self) -> None:
        """Stop background polling."""
        self._running = False
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass
    
    def get_metrics(self, host: str) -> Optional[Dict]:
        """Get latest metrics for a host."""
        return self.metrics.get(host)


async def main():
    print("=== Async Network Programming Demo ===\n")
    
    # 1. Sequential vs Concurrent pings
    print("1. Sequential vs Concurrent Pings")
    hosts = [f"10.0.0.{i}" for i in range(1, 6)]
    
    # Sequential
    start = time.perf_counter()
    sequential_results = []
    for host in hosts:
        result = await async_ping(host)
        sequential_results.append(result)
    sequential_time = time.perf_counter() - start
    
    # Concurrent
    start = time.perf_counter()
    concurrent_results = await asyncio.gather(*[async_ping(h) for h in hosts])
    concurrent_time = time.perf_counter() - start
    
    print(f"   Sequential: {sequential_time:.3f}s")
    print(f"   Concurrent: {concurrent_time:.3f}s")
    print(f"   Speedup: {sequential_time/concurrent_time:.1f}x")
    
    # 2. Ping sweep
    print("\n2. Ping Sweep (10.0.0.1-20)")
    results = await ping_sweep("10.0.0", 1, 20)
    alive = [r.host for r in results if r.success]
    print(f"   Alive hosts: {len(alive)}/{len(results)}")
    
    # 3. Port scanning
    print("\n3. Port Scanning")
    ports = [22, 80, 443, 3306, 5432, 8080, 8443]
    scan_result = await scan_host("10.0.0.1", ports)
    print(f"   Host: {scan_result['host']}")
    print(f"   Open ports: {scan_result['open_ports']}")
    
    # 4. Rate-limited requests
    print("\n4. Rate-Limited Requests (max 3 concurrent)")
    many_hosts = [f"10.0.0.{i}" for i in range(1, 11)]
    start = time.perf_counter()
    results = await rate_limited_requests(many_hosts, max_concurrent=3)
    elapsed = time.perf_counter() - start
    print(f"   Fetched {len(results)} devices in {elapsed:.3f}s")
    
    # 5. Background poller (brief demo)
    print("\n5. Background Poller")
    poller = AsyncDevicePoller(hosts[:3], interval=1.0)
    await poller.start()
    await asyncio.sleep(2.5)  # Let it poll a few times
    await poller.stop()
    print(f"   Final metrics count: {len(poller.metrics)}")
    
    # 6. Exception handling
    print("\n6. Gather with Exceptions")
    async def maybe_fail(host: str) -> Dict:
        if random.random() < 0.3:
            raise ConnectionError(f"Failed to connect to {host}")
        return await fetch_device_info(host)
    
    results = await asyncio.gather(
        *[maybe_fail(h) for h in hosts],
        return_exceptions=True
    )
    successes = [r for r in results if isinstance(r, dict)]
    failures = [r for r in results if isinstance(r, Exception)]
    print(f"   Successes: {len(successes)}, Failures: {len(failures)}")


if __name__ == "__main__":
    asyncio.run(main())
```

**Git checkpoint:**
```bash
git add day_05/
git commit -m "Day 5 Extended: Async network programming"
```

---

## Day 6 Extended: Testing

Create `tests/test_network_models.py`:

```python
#!/usr/bin/env python3
"""
Day 6 Extended: Unit Testing
Learning: pytest, fixtures, mocking, parametrization

Install: pip install pytest pytest-cov pytest-asyncio
"""

import pytest
from datetime import datetime
from unittest.mock import Mock, patch, MagicMock
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'day_04'))

from network_models import (
    NetworkDevice, Router, Switch, Firewall,
    Interface, InterfaceStatus, DeviceStatus,
    create_device,
)


# Fixtures - reusable test setup

@pytest.fixture
def sample_router():
    """Create a sample router for testing."""
    router = Router("test-router", "10.0.0.1", "juniper")
    router.status = DeviceStatus.UP
    return router


@pytest.fixture
def sample_switch():
    """Create a sample switch for testing."""
    switch = Switch("test-switch", "10.0.1.1", "arista")
    switch.status = DeviceStatus.UP
    return switch


@pytest.fixture
def sample_interface():
    """Create a sample interface."""
    return Interface(
        name="eth0",
        ip_address="192.168.1.1",
        mac_address="aa:bb:cc:dd:ee:ff",
        status=InterfaceStatus.UP,
        speed_mbps=1000,
    )


# Basic Tests

class TestInterface:
    """Tests for Interface dataclass."""
    
    def test_interface_creation(self, sample_interface):
        """Test basic interface creation."""
        assert sample_interface.name == "eth0"
        assert sample_interface.ip_address == "192.168.1.1"
        assert sample_interface.is_up == True
    
    def test_interface_defaults(self):
        """Test interface default values."""
        iface = Interface(name="eth1")
        assert iface.ip_address is None
        assert iface.status == InterfaceStatus.DOWN
        assert iface.mtu == 1500
    
    def test_interface_to_dict(self, sample_interface):
        """Test interface serialization."""
        data = sample_interface.to_dict()
        assert data["name"] == "eth0"
        assert data["status"] == "up"
        assert "speed_mbps" in data


class TestRouter:
    """Tests for Router class."""
    
    def test_router_creation(self, sample_router):
        """Test router initialization."""
        assert sample_router.hostname == "test-router"
        assert sample_router.get_device_type() == "router"
        assert sample_router.vendor == "juniper"
    
    def test_add_route(self, sample_router):
        """Test adding routes."""
        sample_router.add_route("0.0.0.0/0", "10.0.0.254")
        sample_router.add_route("192.168.0.0/16", "10.0.1.1", protocol="bgp")
        
        assert len(sample_router.routing_table) == 2
        assert sample_router.routing_table[0]["next_hop"] == "10.0.0.254"
    
    def test_add_bgp_neighbor(self, sample_router):
        """Test BGP neighbor configuration."""
        sample_router.add_bgp_neighbor("10.0.0.2", 65001, "ISP")
        
        assert "10.0.0.2" in sample_router.bgp_neighbors
        assert sample_router.bgp_neighbors["10.0.0.2"]["remote_as"] == 65001
    
    def test_show_ip_route_format(self, sample_router):
        """Test routing table output format."""
        sample_router.add_route("0.0.0.0/0", "10.0.0.254")
        output = sample_router.show_ip_route()
        
        assert "test-router" in output
        assert "0.0.0.0/0" in output
        assert "10.0.0.254" in output


class TestSwitch:
    """Tests for Switch class."""
    
    def test_add_vlan(self, sample_switch):
        """Test VLAN creation."""
        sample_switch.add_vlan(100, "Servers")
        
        assert 100 in sample_switch.vlans
        assert sample_switch.vlans[100] == "Servers"
    
    def test_invalid_vlan_id(self, sample_switch):
        """Test invalid VLAN ID raises error."""
        with pytest.raises(ValueError):
            sample_switch.add_vlan(5000, "Invalid")
        
        with pytest.raises(ValueError):
            sample_switch.add_vlan(0, "Invalid")
    
    def test_delete_default_vlan_raises(self, sample_switch):
        """Test that default VLAN cannot be deleted."""
        with pytest.raises(ValueError, match="Cannot delete default VLAN"):
            sample_switch.delete_vlan(1)
    
    def test_mac_learning(self, sample_switch):
        """Test MAC address learning."""
        sample_switch.learn_mac("aa:bb:cc:dd:ee:ff", "Ethernet1", 100)
        
        result = sample_switch.lookup_mac("AA:BB:CC:DD:EE:FF")  # Test case insensitivity
        assert result is not None
        assert result["interface"] == "Ethernet1"
        assert result["vlan"] == 100


# Parametrized Tests

@pytest.mark.parametrize("device_type,expected_class", [
    ("router", Router),
    ("switch", Switch),
    ("firewall", Firewall),
])
def test_device_factory(device_type, expected_class):
    """Test device factory creates correct types."""
    device = create_device(device_type, "test", "10.0.0.1")
    assert isinstance(device, expected_class)


@pytest.mark.parametrize("hostname,expected_valid", [
    ("router-01", True),
    ("switch.core.01", True),
    ("device_name", True),
    ("", False),
    ("a" * 254, False),  # Too long
    ("invalid hostname", False),  # Contains space
])
def test_hostname_validation(hostname, expected_valid):
    """Test hostname validation."""
    assert NetworkDevice.is_valid_hostname(hostname) == expected_valid


@pytest.mark.parametrize("vlan_id,should_raise", [
    (1, False),
    (100, False),
    (4094, False),
    (0, True),
    (4095, True),
    (-1, True),
])
def test_vlan_id_validation(sample_switch, vlan_id, should_raise):
    """Test VLAN ID boundary conditions."""
    if should_raise:
        with pytest.raises(ValueError):
            sample_switch.add_vlan(vlan_id, "Test")
    else:
        sample_switch.add_vlan(vlan_id, "Test")
        assert vlan_id in sample_switch.vlans


# Mocking Tests

class TestWithMocking:
    """Tests demonstrating mocking."""
    
    def test_ping_success(self, sample_router):
        """Test ping when device is up."""
        sample_router._status = DeviceStatus.UP
        assert sample_router.ping() == True
    
    def test_ping_failure(self, sample_router):
        """Test ping when device is down."""
        sample_router._status = DeviceStatus.DOWN
        assert sample_router.ping() == False
    
    @patch('network_models.datetime')
    def test_updated_at_changes(self, mock_datetime, sample_router):
        """Test that updated_at changes on status change."""
        mock_now = datetime(2024, 1, 15, 12, 0, 0)
        mock_datetime.now.return_value = mock_now
        
        sample_router.status = DeviceStatus.MAINTENANCE
        assert sample_router.updated_at == mock_now


# Integration-style Tests

class TestDeviceIntegration:
    """Integration tests for device interactions."""
    
    def test_full_router_configuration(self):
        """Test complete router setup."""
        router = Router("integration-rtr", "10.0.0.1", "cisco")
        router.status = DeviceStatus.UP
        
        # Add interfaces
        router.add_interface(Interface("ge0/0", "10.0.1.1"))
        router.add_interface(Interface("ge0/1", "10.0.2.1"))
        
        # Add routes
        router.add_route("0.0.0.0/0", "10.0.0.254")
        router.add_route("192.168.0.0/16", "10.0.1.254")
        
        # Add BGP
        router.add_bgp_neighbor("10.0.0.2", 65001)
        
        # Verify
        assert len(router.interfaces) == 2
        assert len(router.routing_table) == 2
        assert len(router.bgp_neighbors) == 1
        
        # Test serialization
        data = router.to_dict()
        assert data["hostname"] == "integration-rtr"
        assert "ge0/0" in data["interfaces"]
    
    def test_device_registry(self):
        """Test device registry tracking."""
        # Clear registry first
        NetworkDevice._registry.clear()
        initial_count = NetworkDevice.device_count
        
        r1 = Router("reg-test-1", "10.0.0.1")
        r2 = Router("reg-test-2", "10.0.0.2")
        s1 = Switch("reg-test-3", "10.0.0.3")
        
        assert NetworkDevice.get_device("reg-test-1") == r1
        assert NetworkDevice.get_device("reg-test-2") == r2
        assert len(NetworkDevice.get_all_devices()) >= 3


# Async Tests (requires pytest-asyncio)

@pytest.mark.asyncio
async def test_async_example():
    """Example async test."""
    import asyncio
    
    async def async_operation():
        await asyncio.sleep(0.01)
        return "success"
    
    result = await async_operation()
    assert result == "success"


# Run with: pytest tests/ -v --cov=day_04 --cov-report=html
```

Create `pytest.ini`:

```ini
[pytest]
testpaths = tests
python_files = test_*.py
python_functions = test_*
python_classes = Test*
addopts = -v --tb=short
markers =
    slow: marks tests as slow
    integration: marks tests as integration tests
asyncio_mode = auto
```

**Git checkpoint:**
```bash
git add tests/ pytest.ini
git commit -m "Day 6 Extended: Comprehensive testing with pytest"
```

---

## Running All Tests

```bash
# Install test dependencies
pip install pytest pytest-cov pytest-asyncio

# Run all tests
pytest

# Run with coverage
pytest --cov=. --cov-report=html

# Run specific test file
pytest tests/test_network_models.py

# Run tests matching pattern
pytest -k "test_router"

# Run with verbose output
pytest -v

# Run and stop on first failure
pytest -x

# Run only marked tests
pytest -m "not slow"
```

---

## Summary Commands

```bash
# Complete Phase 2 exercises
cd ~/dev/python/python-learning

# Run all exercises
python day_04/network_models.py
python day_05/async_network.py
pytest tests/ -v

# Check coverage
pytest --cov=day_04 --cov=day_05 --cov-report=term-missing
```

This extended content adds:
1. Complete OOP hierarchy with Router, Switch, Firewall, LoadBalancer
2. Async programming patterns for network operations
3. Comprehensive pytest testing suite
4. Factory pattern and device registry

All following the same structured approach as your Go and Rust projects!
