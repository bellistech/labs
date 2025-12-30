# Python Phase 7 Extended: DNS Security & Advanced Features

Building on the DNS server foundation with security features, zone transfers, metrics, and operational tooling.

---

## Table of Contents

1. [Day 31: DNS Security - Rate Limiting & Validation](#day-31-dns-security)
2. [Day 32: Zone Transfers (AXFR/IXFR)](#day-32-zone-transfers)
3. [Day 33: Prometheus Metrics & Monitoring](#day-33-metrics)
4. [Day 34: DNS Utilities & Troubleshooting Tools](#day-34-utilities)
5. [Day 35: Complete Project Structure](#day-35-project-structure)

---

## Day 31: DNS Security

### Rate Limiting & Query Validation

Create `dns_server/security.py`:

```python
#!/usr/bin/env python3
"""
Day 31: DNS Security Features
- Rate limiting per client IP
- Query validation
- Response rate limiting (RRL)
- Blocklists
"""

import time
import ipaddress
import re
import logging
from typing import Dict, List, Set, Optional, Tuple
from dataclasses import dataclass, field
from collections import defaultdict
from enum import Enum
import hashlib

logger = logging.getLogger("dns.security")


class ThreatLevel(Enum):
    """Threat classification levels."""
    NONE = 0
    LOW = 1
    MEDIUM = 2
    HIGH = 3
    CRITICAL = 4


@dataclass
class RateLimitConfig:
    """Rate limiting configuration."""
    queries_per_second: int = 20
    queries_per_minute: int = 300
    burst_allowance: int = 50
    block_duration_seconds: int = 60
    whitelist_networks: List[str] = field(default_factory=list)


@dataclass
class ClientState:
    """Per-client rate limiting state."""
    query_times: List[float] = field(default_factory=list)
    blocked_until: float = 0
    total_queries: int = 0
    blocked_queries: int = 0
    
    def is_blocked(self) -> bool:
        return time.time() < self.blocked_until
    
    def cleanup_old_queries(self, window_seconds: int = 60):
        """Remove query times older than window."""
        cutoff = time.time() - window_seconds
        self.query_times = [t for t in self.query_times if t > cutoff]


class RateLimiter:
    """
    Token bucket rate limiter with per-IP tracking.
    
    Implements both queries-per-second and queries-per-minute limits.
    """
    
    def __init__(self, config: RateLimitConfig = None):
        self.config = config or RateLimitConfig()
        self.clients: Dict[str, ClientState] = defaultdict(ClientState)
        self.whitelist: Set[ipaddress.IPv4Network | ipaddress.IPv6Network] = set()
        
        # Parse whitelist networks
        for net in self.config.whitelist_networks:
            try:
                self.whitelist.add(ipaddress.ip_network(net, strict=False))
            except ValueError:
                logger.warning(f"Invalid whitelist network: {net}")
    
    def is_whitelisted(self, client_ip: str) -> bool:
        """Check if client IP is whitelisted."""
        try:
            ip = ipaddress.ip_address(client_ip)
            return any(ip in net for net in self.whitelist)
        except ValueError:
            return False
    
    def check_rate_limit(self, client_ip: str) -> Tuple[bool, str]:
        """
        Check if client is within rate limits.
        
        Returns:
            (allowed, reason) tuple
        """
        if self.is_whitelisted(client_ip):
            return True, "whitelisted"
        
        state = self.clients[client_ip]
        now = time.time()
        
        # Check if blocked
        if state.is_blocked():
            state.blocked_queries += 1
            return False, f"blocked until {state.blocked_until - now:.0f}s"
        
        # Cleanup old entries
        state.cleanup_old_queries(60)
        
        # Check per-minute limit
        if len(state.query_times) >= self.config.queries_per_minute:
            state.blocked_until = now + self.config.block_duration_seconds
            logger.warning(f"Rate limit exceeded for {client_ip}, blocking")
            return False, "per-minute limit exceeded"
        
        # Check per-second limit (last 1 second)
        recent = [t for t in state.query_times if t > now - 1]
        if len(recent) >= self.config.queries_per_second:
            # Allow burst but track it
            if len(recent) >= self.config.queries_per_second + self.config.burst_allowance:
                state.blocked_until = now + self.config.block_duration_seconds
                return False, "burst limit exceeded"
        
        # Record query
        state.query_times.append(now)
        state.total_queries += 1
        
        return True, "ok"
    
    def get_stats(self) -> Dict:
        """Get rate limiter statistics."""
        now = time.time()
        active = sum(1 for s in self.clients.values() if s.query_times)
        blocked = sum(1 for s in self.clients.values() if s.is_blocked())
        
        return {
            "active_clients": active,
            "blocked_clients": blocked,
            "total_tracked": len(self.clients),
            "whitelist_size": len(self.whitelist),
        }
    
    def cleanup(self, max_age_seconds: int = 3600):
        """Remove stale client entries."""
        cutoff = time.time() - max_age_seconds
        stale = [
            ip for ip, state in self.clients.items()
            if not state.query_times or max(state.query_times) < cutoff
        ]
        for ip in stale:
            del self.clients[ip]
        
        if stale:
            logger.debug(f"Cleaned up {len(stale)} stale rate limit entries")


class QueryValidator:
    """
    Validate DNS queries for security issues.
    """
    
    # Suspicious patterns
    SUSPICIOUS_PATTERNS = [
        r'\.{2,}',           # Multiple consecutive dots
        r'^-',               # Label starting with hyphen
        r'-$',               # Label ending with hyphen
        r'[^\w\-\.]',        # Invalid characters
    ]
    
    # Known malicious TLDs (example - real list would be much larger)
    SUSPICIOUS_TLDS = {'tk', 'ml', 'ga', 'cf', 'gq'}
    
    # Maximum label and name lengths per RFC
    MAX_LABEL_LENGTH = 63
    MAX_NAME_LENGTH = 253
    
    def __init__(self):
        self.patterns = [re.compile(p) for p in self.SUSPICIOUS_PATTERNS]
        self.blocklist: Set[str] = set()
        self.blocklist_patterns: List[re.Pattern] = []
    
    def add_blocklist(self, domains: List[str]):
        """Add domains to blocklist."""
        for domain in domains:
            domain = domain.lower().strip()
            if domain.startswith('*.'):
                # Wildcard pattern
                pattern = re.compile(
                    r'(^|\.)'  + re.escape(domain[2:]) + r'$',
                    re.IGNORECASE
                )
                self.blocklist_patterns.append(pattern)
            else:
                self.blocklist.add(domain)
    
    def load_blocklist_file(self, filepath: str):
        """Load blocklist from file (one domain per line)."""
        try:
            with open(filepath) as f:
                domains = [
                    line.strip() for line in f
                    if line.strip() and not line.startswith('#')
                ]
            self.add_blocklist(domains)
            logger.info(f"Loaded {len(domains)} blocklist entries")
        except Exception as e:
            logger.error(f"Failed to load blocklist: {e}")
    
    def validate_query(self, name: str, qtype: int, 
                       client_ip: str) -> Tuple[bool, str, ThreatLevel]:
        """
        Validate a DNS query.
        
        Returns:
            (valid, reason, threat_level) tuple
        """
        name = name.lower().rstrip('.')
        
        # Check name length
        if len(name) > self.MAX_NAME_LENGTH:
            return False, "name too long", ThreatLevel.MEDIUM
        
        # Check label lengths
        for label in name.split('.'):
            if len(label) > self.MAX_LABEL_LENGTH:
                return False, "label too long", ThreatLevel.MEDIUM
            if len(label) == 0:
                return False, "empty label", ThreatLevel.LOW
        
        # Check for suspicious patterns
        for pattern in self.patterns:
            if pattern.search(name):
                return False, f"suspicious pattern", ThreatLevel.MEDIUM
        
        # Check blocklist
        if name in self.blocklist:
            return False, "blocklisted", ThreatLevel.HIGH
        
        for pattern in self.blocklist_patterns:
            if pattern.search(name):
                return False, "blocklist pattern match", ThreatLevel.HIGH
        
        # Check TLD
        parts = name.split('.')
        if parts and parts[-1] in self.SUSPICIOUS_TLDS:
            # Don't block, just flag
            return True, "suspicious TLD", ThreatLevel.LOW
        
        # Check for potential DNS tunneling (very long subdomains)
        if any(len(label) > 50 for label in parts[:-2]):
            return True, "possible tunneling", ThreatLevel.MEDIUM
        
        # Check for excessive subdomain depth
        if len(parts) > 10:
            return True, "excessive depth", ThreatLevel.LOW
        
        return True, "ok", ThreatLevel.NONE
    
    def calculate_entropy(self, s: str) -> float:
        """Calculate Shannon entropy of a string."""
        from collections import Counter
        import math
        
        if not s:
            return 0.0
        
        counts = Counter(s)
        length = len(s)
        
        return -sum(
            (count / length) * math.log2(count / length)
            for count in counts.values()
        )
    
    def detect_tunneling(self, name: str) -> Tuple[bool, float]:
        """
        Detect potential DNS tunneling.
        
        DNS tunneling often uses high-entropy subdomain labels
        to encode data.
        """
        labels = name.lower().rstrip('.').split('.')
        
        # Check subdomain labels (not TLD or main domain)
        if len(labels) < 3:
            return False, 0.0
        
        subdomains = labels[:-2]
        
        # Calculate average entropy
        if not subdomains:
            return False, 0.0
        
        avg_entropy = sum(
            self.calculate_entropy(label) for label in subdomains
        ) / len(subdomains)
        
        # High entropy (> 3.5) suggests encoded data
        # Normal domains typically have entropy < 3.0
        is_suspicious = avg_entropy > 3.5 and any(len(l) > 20 for l in subdomains)
        
        return is_suspicious, avg_entropy


class ResponseRateLimiter:
    """
    Response Rate Limiting (RRL) to mitigate DNS amplification attacks.
    
    Limits identical responses to the same client prefix.
    """
    
    def __init__(self, 
                 responses_per_second: int = 5,
                 window_seconds: int = 1,
                 slip_ratio: int = 2,
                 ipv4_prefix: int = 24,
                 ipv6_prefix: int = 56):
        self.rps = responses_per_second
        self.window = window_seconds
        self.slip_ratio = slip_ratio  # 1 in N gets truncated response
        self.ipv4_prefix = ipv4_prefix
        self.ipv6_prefix = ipv6_prefix
        
        # Track responses: (client_prefix, response_hash) -> [timestamps]
        self.responses: Dict[Tuple[str, str], List[float]] = defaultdict(list)
        self.slip_counter: Dict[Tuple[str, str], int] = defaultdict(int)
    
    def get_client_prefix(self, client_ip: str) -> str:
        """Get client network prefix for grouping."""
        try:
            ip = ipaddress.ip_address(client_ip)
            if isinstance(ip, ipaddress.IPv4Address):
                network = ipaddress.IPv4Network(
                    f"{client_ip}/{self.ipv4_prefix}", strict=False
                )
            else:
                network = ipaddress.IPv6Network(
                    f"{client_ip}/{self.ipv6_prefix}", strict=False
                )
            return str(network.network_address)
        except ValueError:
            return client_ip
    
    def hash_response(self, qname: str, qtype: int, rcode: int) -> str:
        """Create hash of response characteristics."""
        key = f"{qname.lower()}:{qtype}:{rcode}"
        return hashlib.md5(key.encode()).hexdigest()[:16]
    
    def check_response(self, client_ip: str, qname: str,
                       qtype: int, rcode: int) -> Tuple[bool, bool]:
        """
        Check if response should be rate limited.
        
        Returns:
            (allow, truncate) tuple
            - allow=True, truncate=False: Send full response
            - allow=True, truncate=True: Send truncated response (slip)
            - allow=False, truncate=False: Drop response
        """
        prefix = self.get_client_prefix(client_ip)
        resp_hash = self.hash_response(qname, qtype, rcode)
        key = (prefix, resp_hash)
        
        now = time.time()
        
        # Cleanup old entries
        cutoff = now - self.window
        self.responses[key] = [t for t in self.responses[key] if t > cutoff]
        
        # Check rate
        if len(self.responses[key]) >= self.rps:
            # Rate exceeded - check slip
            self.slip_counter[key] += 1
            if self.slip_counter[key] % self.slip_ratio == 0:
                # Slip: send truncated response
                return True, True
            else:
                # Drop
                return False, False
        
        # Record response
        self.responses[key].append(now)
        return True, False
    
    def cleanup(self):
        """Remove stale entries."""
        cutoff = time.time() - self.window * 2
        
        stale_keys = [
            key for key, times in self.responses.items()
            if not times or max(times) < cutoff
        ]
        
        for key in stale_keys:
            del self.responses[key]
            self.slip_counter.pop(key, None)


@dataclass
class SecurityConfig:
    """Complete security configuration."""
    enable_rate_limiting: bool = True
    enable_query_validation: bool = True
    enable_rrl: bool = True
    blocklist_file: Optional[str] = None
    rate_limit: RateLimitConfig = field(default_factory=RateLimitConfig)
    log_blocked_queries: bool = True
    log_suspicious_queries: bool = True


class DNSSecurityManager:
    """
    Unified security manager for DNS server.
    """
    
    def __init__(self, config: SecurityConfig = None):
        self.config = config or SecurityConfig()
        
        self.rate_limiter = RateLimiter(self.config.rate_limit)
        self.validator = QueryValidator()
        self.rrl = ResponseRateLimiter()
        
        # Load blocklist if configured
        if self.config.blocklist_file:
            self.validator.load_blocklist_file(self.config.blocklist_file)
        
        # Statistics
        self.stats = defaultdict(int)
    
    def check_query(self, client_ip: str, qname: str, 
                    qtype: int) -> Tuple[bool, str]:
        """
        Perform all security checks on incoming query.
        
        Returns:
            (allowed, reason) tuple
        """
        # Rate limiting
        if self.config.enable_rate_limiting:
            allowed, reason = self.rate_limiter.check_rate_limit(client_ip)
            if not allowed:
                self.stats["rate_limited"] += 1
                if self.config.log_blocked_queries:
                    logger.warning(f"Rate limited {client_ip}: {reason}")
                return False, f"rate_limit: {reason}"
        
        # Query validation
        if self.config.enable_query_validation:
            valid, reason, threat = self.validator.validate_query(
                qname, qtype, client_ip
            )
            
            if not valid:
                self.stats["validation_blocked"] += 1
                if self.config.log_blocked_queries:
                    logger.warning(
                        f"Blocked query from {client_ip}: {qname} - {reason}"
                    )
                return False, f"validation: {reason}"
            
            if threat != ThreatLevel.NONE:
                self.stats[f"threat_{threat.name.lower()}"] += 1
                if self.config.log_suspicious_queries:
                    logger.info(
                        f"Suspicious query from {client_ip}: {qname} - "
                        f"{reason} (threat: {threat.name})"
                    )
        
        self.stats["allowed"] += 1
        return True, "ok"
    
    def check_response(self, client_ip: str, qname: str,
                       qtype: int, rcode: int) -> Tuple[bool, bool]:
        """
        Check response for RRL.
        
        Returns:
            (allow, truncate) tuple
        """
        if not self.config.enable_rrl:
            return True, False
        
        allow, truncate = self.rrl.check_response(
            client_ip, qname, qtype, rcode
        )
        
        if not allow:
            self.stats["rrl_dropped"] += 1
        elif truncate:
            self.stats["rrl_slipped"] += 1
        
        return allow, truncate
    
    def get_stats(self) -> Dict:
        """Get security statistics."""
        return {
            **dict(self.stats),
            "rate_limiter": self.rate_limiter.get_stats(),
        }
    
    def periodic_cleanup(self):
        """Periodic cleanup of state."""
        self.rate_limiter.cleanup()
        self.rrl.cleanup()


def main():
    print("=== DNS Security Demo ===\n")
    
    # Rate limiter demo
    print("1. Rate Limiter:")
    config = RateLimitConfig(
        queries_per_second=5,
        queries_per_minute=20,
        whitelist_networks=["192.168.0.0/16", "::1/128"]
    )
    limiter = RateLimiter(config)
    
    client = "10.0.0.100"
    for i in range(25):
        allowed, reason = limiter.check_rate_limit(client)
        if i < 5 or i >= 20:
            print(f"   Query {i+1}: {allowed} ({reason})")
    
    print(f"   Stats: {limiter.get_stats()}")
    
    # Query validator demo
    print("\n2. Query Validator:")
    validator = QueryValidator()
    validator.add_blocklist(["malware.example.com", "*.evil.com"])
    
    test_queries = [
        ("www.google.com", 1),
        ("malware.example.com", 1),
        ("test.evil.com", 1),
        ("a" * 70 + ".com", 1),  # Label too long
        ("normal.suspicious.tk", 1),  # Suspicious TLD
        ("aGVsbG8gd29ybGQgdGhpcyBpcyBhIHRlc3Q.tunnel.example.com", 1),  # Tunneling
    ]
    
    for name, qtype in test_queries:
        valid, reason, threat = validator.validate_query(name, qtype, "10.0.0.1")
        status = "✓" if valid else "✗"
        print(f"   {status} {name[:40]:40} - {reason} ({threat.name})")
    
    # Tunneling detection
    print("\n3. DNS Tunneling Detection:")
    tunneling_tests = [
        "www.google.com",
        "aGVsbG8gd29ybGQgdGhpcyBpcyBlbmNvZGVkIGRhdGE.tunnel.example.com",
        "dGhpcyBpcyBhIHZlcnkgbG9uZyBlbmNvZGVkIHN0cmluZw.data.evil.com",
    ]
    
    for name in tunneling_tests:
        is_tunnel, entropy = validator.detect_tunneling(name)
        status = "⚠️ TUNNEL" if is_tunnel else "  normal"
        print(f"   {status} entropy={entropy:.2f} - {name[:50]}")
    
    # Security manager demo
    print("\n4. Security Manager:")
    security = DNSSecurityManager()
    
    # Simulate queries
    for i in range(10):
        allowed, reason = security.check_query("10.0.0.1", "test.example.com", 1)
    
    allowed, reason = security.check_query("10.0.0.1", "malware.example.com", 1)
    print(f"   Malware domain: allowed={allowed} ({reason})")
    
    print(f"   Stats: {security.get_stats()}")


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add dns_server/security.py
git commit -m "Day 31: DNS security features"
```

---

## Day 32: Zone Transfers

Create `dns_server/zone_transfer.py`:

```python
#!/usr/bin/env python3
"""
Day 32: DNS Zone Transfers (AXFR/IXFR)
- AXFR: Full zone transfer
- IXFR: Incremental zone transfer
- TCP transport for large transfers
"""

import asyncio
import socket
import struct
import logging
from typing import List, Dict, Optional, Tuple, AsyncIterator
from dataclasses import dataclass, field
from datetime import datetime

import sys
sys.path.insert(0, '..')
from day_23.dns_protocol import (
    DNSMessage, DNSHeader, DNSQuestion, DNSRecord,
    DNSRecordType, DNSClass, DNSRcode,
    encode_name, decode_name
)

logger = logging.getLogger("dns.transfer")


@dataclass
class ZoneSerial:
    """Zone serial number with change tracking."""
    serial: int
    last_updated: datetime = field(default_factory=datetime.now)
    
    @classmethod
    def generate(cls) -> 'ZoneSerial':
        """Generate serial in YYYYMMDDnn format."""
        now = datetime.now()
        base = int(now.strftime("%Y%m%d")) * 100
        return cls(serial=base + 1)
    
    def increment(self) -> int:
        """Increment serial number."""
        self.serial += 1
        self.last_updated = datetime.now()
        return self.serial


@dataclass
class SOARecord:
    """Start of Authority record data."""
    mname: str      # Primary nameserver
    rname: str      # Admin email (@ replaced with .)
    serial: int     # Zone serial number
    refresh: int    # Refresh interval (seconds)
    retry: int      # Retry interval (seconds)
    expire: int     # Expire time (seconds)
    minimum: int    # Minimum TTL
    
    def pack(self) -> bytes:
        """Pack SOA record data."""
        data = encode_name(self.mname)
        data += encode_name(self.rname)
        data += struct.pack("!IIIII",
            self.serial, self.refresh, self.retry,
            self.expire, self.minimum
        )
        return data
    
    @classmethod
    def unpack(cls, data: bytes, offset: int = 0) -> Tuple['SOARecord', int]:
        """Unpack SOA record data."""
        mname, offset = decode_name(data, offset)
        rname, offset = decode_name(data, offset)
        serial, refresh, retry, expire, minimum = struct.unpack(
            "!IIIII", data[offset:offset+20]
        )
        return cls(
            mname=mname, rname=rname, serial=serial,
            refresh=refresh, retry=retry, expire=expire, minimum=minimum
        ), offset + 20


@dataclass
class ZoneChange:
    """Represents a zone change for IXFR."""
    serial: int
    timestamp: datetime
    additions: List[DNSRecord] = field(default_factory=list)
    deletions: List[DNSRecord] = field(default_factory=list)


class ZoneTransferServer:
    """
    Zone transfer server supporting AXFR and IXFR.
    
    Runs on TCP port 53 (zone transfers require TCP).
    """
    
    def __init__(self, zones: Dict[str, 'DNSZone'],
                 allowed_clients: List[str] = None,
                 listen_port: int = 5353):
        self.zones = zones
        self.allowed_clients = set(allowed_clients or [])
        self.listen_port = listen_port
        self.server = None
        
        # Track zone changes for IXFR
        self.zone_changes: Dict[str, List[ZoneChange]] = {}
    
    def is_allowed(self, client_ip: str) -> bool:
        """Check if client is allowed to transfer."""
        if not self.allowed_clients:
            return True  # Allow all if no restrictions
        return client_ip in self.allowed_clients
    
    async def start(self):
        """Start the zone transfer server."""
        self.server = await asyncio.start_server(
            self._handle_client,
            '0.0.0.0',
            self.listen_port
        )
        logger.info(f"Zone transfer server listening on port {self.listen_port}")
    
    async def stop(self):
        """Stop the server."""
        if self.server:
            self.server.close()
            await self.server.wait_closed()
    
    async def _handle_client(self, reader: asyncio.StreamReader,
                             writer: asyncio.StreamWriter):
        """Handle incoming TCP connection."""
        addr = writer.get_extra_info('peername')
        client_ip = addr[0] if addr else "unknown"
        
        logger.info(f"Zone transfer connection from {client_ip}")
        
        if not self.is_allowed(client_ip):
            logger.warning(f"Unauthorized transfer attempt from {client_ip}")
            writer.close()
            await writer.wait_closed()
            return
        
        try:
            while True:
                # Read length-prefixed DNS message (TCP format)
                length_data = await reader.readexactly(2)
                length = struct.unpack("!H", length_data)[0]
                
                if length == 0:
                    break
                
                message_data = await reader.readexactly(length)
                query = DNSMessage.unpack(message_data)
                
                # Handle the query
                await self._handle_query(query, writer, client_ip)
                
        except asyncio.IncompleteReadError:
            pass
        except Exception as e:
            logger.error(f"Error handling transfer: {e}")
        finally:
            writer.close()
            await writer.wait_closed()
    
    async def _handle_query(self, query: DNSMessage,
                            writer: asyncio.StreamWriter,
                            client_ip: str):
        """Handle zone transfer query."""
        if not query.questions:
            return
        
        q = query.questions[0]
        zone_name = q.name.lower().rstrip('.')
        
        logger.info(f"Transfer request: {zone_name} type={q.qtype} from {client_ip}")
        
        if zone_name not in self.zones:
            # Send NXDOMAIN
            response = self._create_error(query, DNSRcode.NXDOMAIN)
            await self._send_message(writer, response)
            return
        
        zone = self.zones[zone_name]
        
        if q.qtype == DNSRecordType.SOA:
            # SOA query - return SOA record
            await self._send_soa(query, zone, writer)
        
        elif q.qtype == 252:  # AXFR
            await self._send_axfr(query, zone, writer)
        
        elif q.qtype == 251:  # IXFR
            await self._send_ixfr(query, zone, writer)
        
        else:
            response = self._create_error(query, DNSRcode.NOTIMP)
            await self._send_message(writer, response)
    
    async def _send_axfr(self, query: DNSMessage, zone: 'DNSZone',
                         writer: asyncio.StreamWriter):
        """
        Send full zone transfer (AXFR).
        
        Format:
        1. SOA record
        2. All other records
        3. SOA record (again, marks end of transfer)
        """
        logger.info(f"Sending AXFR for {zone.name}")
        
        # Get SOA record
        soa_records = zone.get_records(zone.name, DNSRecordType.SOA)
        if not soa_records:
            response = self._create_error(query, DNSRcode.SERVFAIL)
            await self._send_message(writer, response)
            return
        
        soa = soa_records[0]
        records_sent = 0
        
        # Send SOA first
        response = self._create_response(query, [soa])
        await self._send_message(writer, response)
        records_sent += 1
        
        # Send all other records
        batch = []
        batch_size = 50  # Records per message
        
        for name, type_records in zone.records.items():
            for rtype, records in type_records.items():
                if rtype == DNSRecordType.SOA:
                    continue  # Skip SOA, we handle it separately
                
                for record in records:
                    batch.append(record)
                    
                    if len(batch) >= batch_size:
                        response = self._create_response(query, batch)
                        await self._send_message(writer, response)
                        records_sent += len(batch)
                        batch = []
        
        # Send remaining records
        if batch:
            response = self._create_response(query, batch)
            await self._send_message(writer, response)
            records_sent += len(batch)
        
        # Send SOA again to mark end
        response = self._create_response(query, [soa])
        await self._send_message(writer, response)
        records_sent += 1
        
        logger.info(f"AXFR complete: {records_sent} records sent")
    
    async def _send_ixfr(self, query: DNSMessage, zone: 'DNSZone',
                         writer: asyncio.StreamWriter):
        """
        Send incremental zone transfer (IXFR).
        
        If client's serial is current, send just SOA.
        Otherwise, send changes or fall back to AXFR.
        """
        # Get client's serial from query authority section
        client_serial = 0
        if query.authority:
            for rr in query.authority:
                if rr.rtype == DNSRecordType.SOA:
                    soa, _ = SOARecord.unpack(rr.rdata)
                    client_serial = soa.serial
                    break
        
        zone_name = zone.name.lower()
        soa_records = zone.get_records(zone_name, DNSRecordType.SOA)
        
        if not soa_records:
            response = self._create_error(query, DNSRcode.SERVFAIL)
            await self._send_message(writer, response)
            return
        
        current_soa = soa_records[0]
        # Parse current SOA to get serial
        soa_data, _ = SOARecord.unpack(current_soa.rdata)
        current_serial = soa_data.serial
        
        logger.info(f"IXFR request: client={client_serial}, current={current_serial}")
        
        if client_serial == current_serial:
            # Client is up to date - send just SOA
            response = self._create_response(query, [current_soa])
            await self._send_message(writer, response)
            return
        
        # Check if we have incremental changes
        changes = self.zone_changes.get(zone_name, [])
        applicable = [c for c in changes if c.serial > client_serial]
        
        if not applicable:
            # No incremental data, fall back to AXFR
            logger.info("No IXFR data available, falling back to AXFR")
            await self._send_axfr(query, zone, writer)
            return
        
        # Send IXFR
        # Format: SOA (new), [SOA (old), deletions, SOA (new), additions]*, SOA (new)
        
        response = self._create_response(query, [current_soa])
        await self._send_message(writer, response)
        
        for change in applicable:
            # Send old SOA, deletions, new SOA, additions
            # (Simplified - real IXFR is more complex)
            all_records = change.deletions + change.additions
            if all_records:
                response = self._create_response(query, all_records)
                await self._send_message(writer, response)
        
        # Final SOA
        response = self._create_response(query, [current_soa])
        await self._send_message(writer, response)
        
        logger.info(f"IXFR complete: {len(applicable)} changes sent")
    
    async def _send_soa(self, query: DNSMessage, zone: 'DNSZone',
                        writer: asyncio.StreamWriter):
        """Send SOA record response."""
        soa_records = zone.get_records(zone.name, DNSRecordType.SOA)
        if soa_records:
            response = self._create_response(query, soa_records)
        else:
            response = self._create_error(query, DNSRcode.NXDOMAIN)
        await self._send_message(writer, response)
    
    def _create_response(self, query: DNSMessage,
                         answers: List[DNSRecord]) -> DNSMessage:
        """Create response message."""
        return DNSMessage(
            header=DNSHeader(
                id=query.header.id,
                qr=1,
                aa=1,
                rcode=DNSRcode.NOERROR,
            ),
            questions=query.questions,
            answers=answers,
        )
    
    def _create_error(self, query: DNSMessage, rcode: int) -> DNSMessage:
        """Create error response."""
        return DNSMessage(
            header=DNSHeader(
                id=query.header.id,
                qr=1,
                rcode=rcode,
            ),
            questions=query.questions,
        )
    
    async def _send_message(self, writer: asyncio.StreamWriter,
                            message: DNSMessage):
        """Send length-prefixed DNS message over TCP."""
        data = message.pack()
        length = struct.pack("!H", len(data))
        writer.write(length + data)
        await writer.drain()


class ZoneTransferClient:
    """
    Zone transfer client for pulling zones from primary servers.
    """
    
    def __init__(self, timeout: float = 30.0):
        self.timeout = timeout
    
    async def axfr(self, zone_name: str, server: str,
                   port: int = 53) -> List[DNSRecord]:
        """
        Perform AXFR (full zone transfer).
        
        Returns list of all records in the zone.
        """
        logger.info(f"Starting AXFR for {zone_name} from {server}")
        
        records = []
        soa_count = 0
        
        async for record in self._transfer(zone_name, 252, server, port):
            records.append(record)
            
            if record.rtype == DNSRecordType.SOA:
                soa_count += 1
                if soa_count == 2:
                    # Second SOA marks end of transfer
                    break
        
        logger.info(f"AXFR complete: {len(records)} records")
        return records
    
    async def ixfr(self, zone_name: str, server: str,
                   current_serial: int, port: int = 53) -> List[DNSRecord]:
        """
        Perform IXFR (incremental zone transfer).
        """
        logger.info(f"Starting IXFR for {zone_name} from {server}")
        
        records = []
        
        async for record in self._transfer(zone_name, 251, server, port,
                                           serial=current_serial):
            records.append(record)
        
        logger.info(f"IXFR complete: {len(records)} records")
        return records
    
    async def _transfer(self, zone_name: str, qtype: int,
                        server: str, port: int,
                        serial: int = None) -> AsyncIterator[DNSRecord]:
        """
        Perform zone transfer and yield records.
        """
        # Create query
        query = DNSMessage(
            header=DNSHeader(id=12345, rd=0),
            questions=[DNSQuestion(name=zone_name, qtype=qtype)]
        )
        
        # For IXFR, include current SOA in authority section
        if serial is not None:
            soa_rdata = SOARecord(
                mname="", rname="", serial=serial,
                refresh=0, retry=0, expire=0, minimum=0
            ).pack()
            query.authority.append(DNSRecord(
                name=zone_name,
                rtype=DNSRecordType.SOA,
                rdata=soa_rdata
            ))
            query.header.nscount = 1
        
        # Connect
        reader, writer = await asyncio.wait_for(
            asyncio.open_connection(server, port),
            timeout=self.timeout
        )
        
        try:
            # Send query
            data = query.pack()
            writer.write(struct.pack("!H", len(data)) + data)
            await writer.drain()
            
            # Receive responses
            while True:
                try:
                    length_data = await asyncio.wait_for(
                        reader.readexactly(2),
                        timeout=self.timeout
                    )
                except asyncio.IncompleteReadError:
                    break
                
                length = struct.unpack("!H", length_data)[0]
                if length == 0:
                    break
                
                message_data = await reader.readexactly(length)
                response = DNSMessage.unpack(message_data)
                
                if response.header.rcode != DNSRcode.NOERROR:
                    raise Exception(f"Transfer failed: {response.header.rcode}")
                
                for record in response.answers:
                    yield record
                
        finally:
            writer.close()
            await writer.wait_closed()


async def main():
    print("=== Zone Transfer Demo ===\n")
    
    # SOA record handling
    print("1. SOA Record:")
    soa = SOARecord(
        mname="ns1.example.com",
        rname="admin.example.com",
        serial=2024010101,
        refresh=3600,
        retry=900,
        expire=604800,
        minimum=86400
    )
    packed = soa.pack()
    unpacked, _ = SOARecord.unpack(packed)
    print(f"   Serial: {unpacked.serial}")
    print(f"   Primary NS: {unpacked.mname}")
    print(f"   Admin: {unpacked.rname}")
    
    # Zone serial
    print("\n2. Zone Serial Management:")
    serial = ZoneSerial.generate()
    print(f"   Generated: {serial.serial}")
    serial.increment()
    print(f"   Incremented: {serial.serial}")
    
    print("\n3. Zone Transfer (requires server):")
    print("   AXFR: dig @server AXFR example.com")
    print("   IXFR: dig @server IXFR example.com")
    
    # Client demo (would need a real server)
    print("\n4. Transfer Client Usage:")
    print("""
    client = ZoneTransferClient()
    
    # Full transfer
    records = await client.axfr("example.com", "ns1.example.com")
    
    # Incremental transfer
    records = await client.ixfr("example.com", "ns1.example.com", 
                                current_serial=2024010100)
    """)


if __name__ == "__main__":
    asyncio.run(main())
```

**Git checkpoint:**
```bash
git add dns_server/zone_transfer.py
git commit -m "Day 32: Zone transfer support (AXFR/IXFR)"
```

---

## Day 33: Prometheus Metrics

Create `dns_server/metrics.py`:

```python
#!/usr/bin/env python3
"""
Day 33: Prometheus Metrics & Monitoring
- Query metrics
- Cache statistics
- Performance tracking
- Health endpoints
"""

import time
import asyncio
from typing import Dict, List, Optional
from dataclasses import dataclass, field
from collections import defaultdict
from contextlib import contextmanager
import threading


@dataclass
class MetricValue:
    """Single metric value with labels."""
    value: float
    labels: Dict[str, str] = field(default_factory=dict)
    timestamp: float = field(default_factory=time.time)


class Counter:
    """Prometheus-style counter metric."""
    
    def __init__(self, name: str, description: str, labels: List[str] = None):
        self.name = name
        self.description = description
        self.label_names = labels or []
        self._values: Dict[tuple, float] = defaultdict(float)
        self._lock = threading.Lock()
    
    def inc(self, value: float = 1, **labels):
        """Increment counter."""
        key = tuple(labels.get(l, "") for l in self.label_names)
        with self._lock:
            self._values[key] += value
    
    def get(self, **labels) -> float:
        """Get counter value."""
        key = tuple(labels.get(l, "") for l in self.label_names)
        return self._values.get(key, 0)
    
    def collect(self) -> List[MetricValue]:
        """Collect all metric values."""
        result = []
        with self._lock:
            for key, value in self._values.items():
                labels = dict(zip(self.label_names, key))
                result.append(MetricValue(value=value, labels=labels))
        return result


class Gauge:
    """Prometheus-style gauge metric."""
    
    def __init__(self, name: str, description: str, labels: List[str] = None):
        self.name = name
        self.description = description
        self.label_names = labels or []
        self._values: Dict[tuple, float] = {}
        self._lock = threading.Lock()
    
    def set(self, value: float, **labels):
        """Set gauge value."""
        key = tuple(labels.get(l, "") for l in self.label_names)
        with self._lock:
            self._values[key] = value
    
    def inc(self, value: float = 1, **labels):
        """Increment gauge."""
        key = tuple(labels.get(l, "") for l in self.label_names)
        with self._lock:
            self._values[key] = self._values.get(key, 0) + value
    
    def dec(self, value: float = 1, **labels):
        """Decrement gauge."""
        self.inc(-value, **labels)
    
    def get(self, **labels) -> float:
        """Get gauge value."""
        key = tuple(labels.get(l, "") for l in self.label_names)
        return self._values.get(key, 0)
    
    def collect(self) -> List[MetricValue]:
        """Collect all metric values."""
        result = []
        with self._lock:
            for key, value in self._values.items():
                labels = dict(zip(self.label_names, key))
                result.append(MetricValue(value=value, labels=labels))
        return result


class Histogram:
    """Prometheus-style histogram metric."""
    
    DEFAULT_BUCKETS = (
        0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0
    )
    
    def __init__(self, name: str, description: str,
                 labels: List[str] = None,
                 buckets: tuple = None):
        self.name = name
        self.description = description
        self.label_names = labels or []
        self.buckets = buckets or self.DEFAULT_BUCKETS
        
        self._counts: Dict[tuple, Dict[float, int]] = defaultdict(
            lambda: {b: 0 for b in self.buckets}
        )
        self._sums: Dict[tuple, float] = defaultdict(float)
        self._totals: Dict[tuple, int] = defaultdict(int)
        self._lock = threading.Lock()
    
    def observe(self, value: float, **labels):
        """Observe a value."""
        key = tuple(labels.get(l, "") for l in self.label_names)
        with self._lock:
            self._sums[key] += value
            self._totals[key] += 1
            for bucket in self.buckets:
                if value <= bucket:
                    self._counts[key][bucket] += 1
    
    @contextmanager
    def time(self, **labels):
        """Context manager to time operations."""
        start = time.perf_counter()
        try:
            yield
        finally:
            self.observe(time.perf_counter() - start, **labels)
    
    def collect(self) -> List[MetricValue]:
        """Collect all metric values."""
        result = []
        with self._lock:
            for key in self._counts:
                labels = dict(zip(self.label_names, key))
                
                # Bucket values
                for bucket, count in self._counts[key].items():
                    bucket_labels = {**labels, "le": str(bucket)}
                    result.append(MetricValue(
                        value=count,
                        labels=bucket_labels
                    ))
                
                # +Inf bucket
                inf_labels = {**labels, "le": "+Inf"}
                result.append(MetricValue(
                    value=self._totals[key],
                    labels=inf_labels
                ))
                
                # Sum and count
                result.append(MetricValue(
                    value=self._sums[key],
                    labels={**labels, "_type": "sum"}
                ))
                result.append(MetricValue(
                    value=self._totals[key],
                    labels={**labels, "_type": "count"}
                ))
        
        return result


class DNSMetrics:
    """
    DNS server metrics collection.
    """
    
    def __init__(self):
        # Query metrics
        self.queries_total = Counter(
            "dns_queries_total",
            "Total DNS queries received",
            ["type", "class", "zone"]
        )
        
        self.responses_total = Counter(
            "dns_responses_total",
            "Total DNS responses sent",
            ["rcode", "type"]
        )
        
        # Latency
        self.query_duration = Histogram(
            "dns_query_duration_seconds",
            "DNS query processing duration",
            ["type"],
            buckets=(0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5)
        )
        
        # Cache metrics
        self.cache_size = Gauge(
            "dns_cache_entries",
            "Number of entries in DNS cache"
        )
        
        self.cache_hits = Counter(
            "dns_cache_hits_total",
            "Cache hit count"
        )
        
        self.cache_misses = Counter(
            "dns_cache_misses_total",
            "Cache miss count"
        )
        
        # Zone metrics
        self.zone_records = Gauge(
            "dns_zone_records",
            "Number of records in zone",
            ["zone"]
        )
        
        # Security metrics
        self.blocked_queries = Counter(
            "dns_blocked_queries_total",
            "Blocked queries",
            ["reason"]
        )
        
        self.rate_limited = Counter(
            "dns_rate_limited_total",
            "Rate limited queries"
        )
        
        # Connection metrics
        self.active_connections = Gauge(
            "dns_active_connections",
            "Active TCP connections"
        )
        
        # Upstream metrics
        self.upstream_queries = Counter(
            "dns_upstream_queries_total",
            "Queries sent to upstream servers",
            ["server"]
        )
        
        self.upstream_failures = Counter(
            "dns_upstream_failures_total",
            "Failed upstream queries",
            ["server", "reason"]
        )
        
        self.upstream_latency = Histogram(
            "dns_upstream_latency_seconds",
            "Upstream query latency",
            ["server"]
        )
    
    def record_query(self, qtype: str, qclass: str, zone: str):
        """Record incoming query."""
        self.queries_total.inc(type=qtype, **{"class": qclass}, zone=zone)
    
    def record_response(self, rcode: str, qtype: str, duration: float):
        """Record outgoing response."""
        self.responses_total.inc(rcode=rcode, type=qtype)
        self.query_duration.observe(duration, type=qtype)
    
    def record_cache_hit(self):
        """Record cache hit."""
        self.cache_hits.inc()
    
    def record_cache_miss(self):
        """Record cache miss."""
        self.cache_misses.inc()
    
    def update_cache_size(self, size: int):
        """Update cache size gauge."""
        self.cache_size.set(size)
    
    def record_blocked(self, reason: str):
        """Record blocked query."""
        self.blocked_queries.inc(reason=reason)
    
    def record_upstream(self, server: str, duration: float, success: bool):
        """Record upstream query."""
        self.upstream_queries.inc(server=server)
        self.upstream_latency.observe(duration, server=server)
        if not success:
            self.upstream_failures.inc(server=server, reason="timeout")
    
    def export_prometheus(self) -> str:
        """Export metrics in Prometheus text format."""
        lines = []
        
        def format_metric(name: str, description: str,
                         metric_type: str, values: List[MetricValue]):
            lines.append(f"# HELP {name} {description}")
            lines.append(f"# TYPE {name} {metric_type}")
            
            for mv in values:
                if mv.labels:
                    label_str = ",".join(
                        f'{k}="{v}"' for k, v in mv.labels.items()
                        if not k.startswith("_")
                    )
                    lines.append(f"{name}{{{label_str}}} {mv.value}")
                else:
                    lines.append(f"{name} {mv.value}")
        
        # Export all metrics
        format_metric(
            self.queries_total.name,
            self.queries_total.description,
            "counter",
            self.queries_total.collect()
        )
        
        format_metric(
            self.responses_total.name,
            self.responses_total.description,
            "counter",
            self.responses_total.collect()
        )
        
        format_metric(
            self.cache_size.name,
            self.cache_size.description,
            "gauge",
            self.cache_size.collect()
        )
        
        format_metric(
            self.cache_hits.name,
            self.cache_hits.description,
            "counter",
            self.cache_hits.collect()
        )
        
        format_metric(
            self.cache_misses.name,
            self.cache_misses.description,
            "counter",
            self.cache_misses.collect()
        )
        
        format_metric(
            self.blocked_queries.name,
            self.blocked_queries.description,
            "counter",
            self.blocked_queries.collect()
        )
        
        # Histogram (special format)
        for mv in self.query_duration.collect():
            if "_type" in mv.labels:
                suffix = "_" + mv.labels["_type"]
                label_str = ",".join(
                    f'{k}="{v}"' for k, v in mv.labels.items()
                    if not k.startswith("_")
                )
                if label_str:
                    lines.append(f"{self.query_duration.name}{suffix}{{{label_str}}} {mv.value}")
                else:
                    lines.append(f"{self.query_duration.name}{suffix} {mv.value}")
            elif "le" in mv.labels:
                label_str = ",".join(f'{k}="{v}"' for k, v in mv.labels.items())
                lines.append(f"{self.query_duration.name}_bucket{{{label_str}}} {mv.value}")
        
        return "\n".join(lines)


class MetricsServer:
    """
    HTTP server for Prometheus metrics endpoint.
    """
    
    def __init__(self, metrics: DNSMetrics, port: int = 9153):
        self.metrics = metrics
        self.port = port
        self.server = None
    
    async def start(self):
        """Start metrics HTTP server."""
        from aiohttp import web
        
        app = web.Application()
        app.router.add_get("/metrics", self._handle_metrics)
        app.router.add_get("/health", self._handle_health)
        app.router.add_get("/ready", self._handle_ready)
        
        runner = web.AppRunner(app)
        await runner.setup()
        
        self.server = web.TCPSite(runner, "0.0.0.0", self.port)
        await self.server.start()
        
        print(f"Metrics server listening on :{self.port}")
    
    async def _handle_metrics(self, request):
        """Handle /metrics endpoint."""
        from aiohttp import web
        
        content = self.metrics.export_prometheus()
        return web.Response(
            text=content,
            content_type="text/plain; charset=utf-8"
        )
    
    async def _handle_health(self, request):
        """Handle /health endpoint."""
        from aiohttp import web
        return web.json_response({"status": "healthy"})
    
    async def _handle_ready(self, request):
        """Handle /ready endpoint."""
        from aiohttp import web
        return web.json_response({"status": "ready"})


def main():
    print("=== DNS Metrics Demo ===\n")
    
    metrics = DNSMetrics()
    
    # Simulate some queries
    print("1. Recording Queries:")
    for _ in range(100):
        metrics.record_query("A", "IN", "example.com")
    for _ in range(50):
        metrics.record_query("AAAA", "IN", "example.com")
    for _ in range(20):
        metrics.record_query("MX", "IN", "example.com")
    
    print(f"   Total A queries: {metrics.queries_total.get(type='A')}")
    print(f"   Total AAAA queries: {metrics.queries_total.get(type='AAAA')}")
    
    # Simulate responses
    print("\n2. Recording Responses:")
    import random
    for _ in range(150):
        duration = random.uniform(0.0001, 0.01)
        rcode = random.choice(["NOERROR"] * 9 + ["NXDOMAIN"])
        metrics.record_response(rcode, "A", duration)
    
    # Cache metrics
    print("\n3. Cache Metrics:")
    for _ in range(80):
        metrics.record_cache_hit()
    for _ in range(20):
        metrics.record_cache_miss()
    metrics.update_cache_size(1500)
    
    print(f"   Cache hits: {metrics.cache_hits.get()}")
    print(f"   Cache misses: {metrics.cache_misses.get()}")
    print(f"   Cache size: {metrics.cache_size.get()}")
    
    # Export Prometheus format
    print("\n4. Prometheus Export (sample):")
    export = metrics.export_prometheus()
    # Show first 20 lines
    for line in export.split("\n")[:20]:
        print(f"   {line}")
    print("   ...")
    
    print(f"\n   Total export: {len(export)} bytes, {len(export.split(chr(10)))} lines")


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add dns_server/metrics.py
git commit -m "Day 33: Prometheus metrics"
```

---

## Day 34: DNS Utilities

Create `dns_server/utils.py`:

```python
#!/usr/bin/env python3
"""
Day 34: DNS Utilities & Troubleshooting Tools
- dig-like query tool
- DNS debugging
- Zone validation
- Performance testing
"""

import asyncio
import socket
import time
import random
from typing import List, Dict, Optional, Tuple
from dataclasses import dataclass
import ipaddress

import sys
sys.path.insert(0, '..')
from day_23.dns_protocol import (
    DNSMessage, DNSHeader, DNSQuestion, DNSRecord,
    DNSRecordType, DNSRcode, create_query
)


@dataclass
class QueryResult:
    """DNS query result."""
    query_name: str
    query_type: str
    server: str
    rcode: str
    answers: List[Dict]
    authority: List[Dict]
    additional: List[Dict]
    query_time_ms: float
    message_size: int
    flags: Dict[str, bool]


class DNSClient:
    """
    DNS client for querying servers.
    Similar to dig command.
    """
    
    def __init__(self, timeout: float = 5.0):
        self.timeout = timeout
    
    async def query(self, name: str, qtype: int = DNSRecordType.A,
                    server: str = "8.8.8.8", port: int = 53,
                    use_tcp: bool = False) -> QueryResult:
        """
        Send DNS query and return result.
        """
        query = create_query(name, qtype)
        start = time.perf_counter()
        
        if use_tcp:
            response_data = await self._query_tcp(query, server, port)
        else:
            response_data = await self._query_udp(query, server, port)
        
        elapsed = (time.perf_counter() - start) * 1000
        response = DNSMessage.unpack(response_data)
        
        return QueryResult(
            query_name=name,
            query_type=DNSRecordType(qtype).name,
            server=f"{server}:{port}",
            rcode=DNSRcode(response.header.rcode).name,
            answers=[self._format_record(r) for r in response.answers],
            authority=[self._format_record(r) for r in response.authority],
            additional=[self._format_record(r) for r in response.additional],
            query_time_ms=round(elapsed, 2),
            message_size=len(response_data),
            flags={
                "qr": bool(response.header.qr),
                "aa": bool(response.header.aa),
                "tc": bool(response.header.tc),
                "rd": bool(response.header.rd),
                "ra": bool(response.header.ra),
            }
        )
    
    async def _query_udp(self, query: DNSMessage, server: str, port: int) -> bytes:
        """Send UDP query."""
        loop = asyncio.get_event_loop()
        
        try:
            ipaddress.IPv6Address(server)
            family = socket.AF_INET6
        except:
            family = socket.AF_INET
        
        sock = socket.socket(family, socket.SOCK_DGRAM)
        sock.setblocking(False)
        
        try:
            await loop.sock_sendto(sock, query.pack(), (server, port))
            data = await asyncio.wait_for(
                loop.sock_recv(sock, 4096),
                timeout=self.timeout
            )
            return data
        finally:
            sock.close()
    
    async def _query_tcp(self, query: DNSMessage, server: str, port: int) -> bytes:
        """Send TCP query."""
        import struct
        
        reader, writer = await asyncio.wait_for(
            asyncio.open_connection(server, port),
            timeout=self.timeout
        )
        
        try:
            data = query.pack()
            writer.write(struct.pack("!H", len(data)) + data)
            await writer.drain()
            
            length_data = await reader.readexactly(2)
            length = struct.unpack("!H", length_data)[0]
            return await reader.readexactly(length)
        finally:
            writer.close()
            await writer.wait_closed()
    
    def _format_record(self, record: DNSRecord) -> Dict:
        """Format record for display."""
        return {
            "name": record.name,
            "type": DNSRecordType(record.rtype).name if record.rtype in DNSRecordType._value2member_map_ else f"TYPE{record.rtype}",
            "ttl": record.ttl,
            "data": record.parsed_data or record.rdata.hex(),
        }


class DNSBenchmark:
    """
    DNS performance benchmarking tool.
    """
    
    def __init__(self, server: str = "8.8.8.8", port: int = 53):
        self.server = server
        self.port = port
        self.client = DNSClient(timeout=2.0)
    
    async def benchmark(self, domains: List[str],
                        qtype: int = DNSRecordType.A,
                        num_queries: int = 100,
                        concurrency: int = 10) -> Dict:
        """
        Benchmark DNS server performance.
        """
        print(f"Benchmarking {self.server}:{self.port}")
        print(f"Queries: {num_queries}, Concurrency: {concurrency}")
        
        results = {
            "success": 0,
            "failed": 0,
            "latencies": [],
            "errors": [],
        }
        
        semaphore = asyncio.Semaphore(concurrency)
        
        async def query_one(domain: str) -> Optional[float]:
            async with semaphore:
                try:
                    result = await self.client.query(
                        domain, qtype, self.server, self.port
                    )
                    return result.query_time_ms
                except Exception as e:
                    results["errors"].append(str(e))
                    return None
        
        start = time.perf_counter()
        
        # Generate queries
        queries = [
            random.choice(domains)
            for _ in range(num_queries)
        ]
        
        # Run queries
        tasks = [query_one(d) for d in queries]
        latencies = await asyncio.gather(*tasks)
        
        total_time = time.perf_counter() - start
        
        # Calculate stats
        valid_latencies = [l for l in latencies if l is not None]
        results["success"] = len(valid_latencies)
        results["failed"] = num_queries - len(valid_latencies)
        results["latencies"] = valid_latencies
        
        if valid_latencies:
            results["stats"] = {
                "min_ms": min(valid_latencies),
                "max_ms": max(valid_latencies),
                "avg_ms": sum(valid_latencies) / len(valid_latencies),
                "median_ms": sorted(valid_latencies)[len(valid_latencies)//2],
                "p95_ms": sorted(valid_latencies)[int(len(valid_latencies)*0.95)],
                "p99_ms": sorted(valid_latencies)[int(len(valid_latencies)*0.99)],
                "qps": results["success"] / total_time,
                "total_time_s": total_time,
            }
        
        return results
    
    def print_results(self, results: Dict):
        """Print benchmark results."""
        print(f"\n{'='*50}")
        print("Benchmark Results")
        print(f"{'='*50}")
        print(f"Success: {results['success']}")
        print(f"Failed:  {results['failed']}")
        
        if "stats" in results:
            stats = results["stats"]
            print(f"\nLatency:")
            print(f"  Min:    {stats['min_ms']:.2f} ms")
            print(f"  Max:    {stats['max_ms']:.2f} ms")
            print(f"  Avg:    {stats['avg_ms']:.2f} ms")
            print(f"  Median: {stats['median_ms']:.2f} ms")
            print(f"  P95:    {stats['p95_ms']:.2f} ms")
            print(f"  P99:    {stats['p99_ms']:.2f} ms")
            print(f"\nThroughput:")
            print(f"  QPS:    {stats['qps']:.2f}")
            print(f"  Time:   {stats['total_time_s']:.2f} s")


class ZoneValidator:
    """
    Validate DNS zone configuration.
    """
    
    def __init__(self):
        self.errors: List[str] = []
        self.warnings: List[str] = []
    
    def validate_zone(self, zone: 'DNSZone') -> bool:
        """
        Validate a DNS zone.
        """
        self.errors = []
        self.warnings = []
        
        # Check SOA record exists
        soa_records = zone.get_records(zone.name, DNSRecordType.SOA)
        if not soa_records:
            self.errors.append(f"Missing SOA record for {zone.name}")
        elif len(soa_records) > 1:
            self.errors.append(f"Multiple SOA records for {zone.name}")
        
        # Check NS records exist
        ns_records = zone.get_records(zone.name, DNSRecordType.NS)
        if not ns_records:
            self.errors.append(f"Missing NS records for {zone.name}")
        elif len(ns_records) < 2:
            self.warnings.append(f"Only {len(ns_records)} NS record(s) - recommend at least 2")
        
        # Check for orphan CNAMEs
        for name, type_records in zone.records.items():
            if DNSRecordType.CNAME in type_records:
                if len(type_records) > 1:
                    self.errors.append(
                        f"CNAME at {name} coexists with other records"
                    )
                
                # Check CNAME target exists
                cname = type_records[DNSRecordType.CNAME][0]
                target = cname.parsed_data
                if target and not zone.get_records(target, DNSRecordType.A):
                    if not zone.get_records(target, DNSRecordType.AAAA):
                        self.warnings.append(
                            f"CNAME {name} -> {target} target has no A/AAAA"
                        )
        
        # Check MX targets have A/AAAA records
        mx_records = zone.get_records(zone.name, DNSRecordType.MX)
        for mx in mx_records:
            # MX parsed_data is "priority exchange"
            if mx.parsed_data:
                parts = mx.parsed_data.split()
                if len(parts) >= 2:
                    exchange = parts[1]
                    if not zone.get_records(exchange, DNSRecordType.A):
                        self.warnings.append(
                            f"MX {exchange} has no A record"
                        )
        
        return len(self.errors) == 0
    
    def print_report(self, zone_name: str):
        """Print validation report."""
        print(f"\nZone Validation Report: {zone_name}")
        print("=" * 50)
        
        if self.errors:
            print(f"\nErrors ({len(self.errors)}):")
            for e in self.errors:
                print(f"  ✗ {e}")
        
        if self.warnings:
            print(f"\nWarnings ({len(self.warnings)}):")
            for w in self.warnings:
                print(f"  ⚠ {w}")
        
        if not self.errors and not self.warnings:
            print("  ✓ Zone is valid")
        
        print()


def format_dig_output(result: QueryResult) -> str:
    """Format query result like dig output."""
    lines = []
    
    lines.append(f"; <<>> Python DNS Client <<>> {result.query_name} {result.query_type}")
    lines.append(f";; Got answer:")
    lines.append(f";; ->>HEADER<<- opcode: QUERY, status: {result.rcode}")
    
    flags = ", ".join(k for k, v in result.flags.items() if v)
    lines.append(f";; flags: {flags}; QUERY: 1, ANSWER: {len(result.answers)}, "
                 f"AUTHORITY: {len(result.authority)}, ADDITIONAL: {len(result.additional)}")
    
    if result.answers:
        lines.append("\n;; ANSWER SECTION:")
        for a in result.answers:
            lines.append(f"{a['name']:30} {a['ttl']:6} IN {a['type']:6} {a['data']}")
    
    if result.authority:
        lines.append("\n;; AUTHORITY SECTION:")
        for a in result.authority:
            lines.append(f"{a['name']:30} {a['ttl']:6} IN {a['type']:6} {a['data']}")
    
    lines.append(f"\n;; Query time: {result.query_time_ms} msec")
    lines.append(f";; SERVER: {result.server}")
    lines.append(f";; MSG SIZE  rcvd: {result.message_size}")
    
    return "\n".join(lines)


async def main():
    print("=== DNS Utilities Demo ===\n")
    
    client = DNSClient()
    
    # Basic queries
    print("1. Basic Queries (like dig):")
    
    queries = [
        ("google.com", DNSRecordType.A),
        ("google.com", DNSRecordType.AAAA),
        ("google.com", DNSRecordType.MX),
    ]
    
    for name, qtype in queries:
        try:
            result = await client.query(name, qtype)
            print(f"\n{format_dig_output(result)}")
        except Exception as e:
            print(f"   Error: {e}")
    
    # Benchmark
    print("\n\n2. Performance Benchmark:")
    benchmark = DNSBenchmark("8.8.8.8")
    
    domains = [
        "google.com", "facebook.com", "amazon.com",
        "microsoft.com", "apple.com", "netflix.com",
    ]
    
    try:
        results = await benchmark.benchmark(domains, num_queries=50, concurrency=5)
        benchmark.print_results(results)
    except Exception as e:
        print(f"   Benchmark error: {e}")
    
    # Compare servers
    print("\n3. Server Comparison:")
    servers = [
        ("8.8.8.8", "Google"),
        ("1.1.1.1", "Cloudflare"),
        ("9.9.9.9", "Quad9"),
    ]
    
    for server, name in servers:
        try:
            result = await client.query("example.com", DNSRecordType.A, server)
            print(f"   {name:12} ({server}): {result.query_time_ms:.2f} ms")
        except Exception as e:
            print(f"   {name:12} ({server}): Error - {e}")


if __name__ == "__main__":
    asyncio.run(main())
```

**Git checkpoint:**
```bash
git add dns_server/utils.py
git commit -m "Day 34: DNS utilities and tools"
```

---

## Day 35: Complete Project Structure

```
python-dns-server/
├── dns_server/
│   ├── __init__.py
│   ├── __main__.py           # Entry point
│   ├── server.py             # Main DNS server
│   ├── protocol.py           # DNS protocol
│   ├── cache.py              # Caching
│   ├── zones.py              # Zone management
│   ├── parser.py             # Zone file parser
│   ├── resolver.py           # Recursive resolver
│   ├── security.py           # Rate limiting, validation
│   ├── zone_transfer.py      # AXFR/IXFR
│   ├── metrics.py            # Prometheus metrics
│   ├── utils.py              # CLI tools
│   ├── ipv6.py               # IPv6 utilities
│   └── config.py             # Configuration
├── tests/
│   ├── __init__.py
│   ├── conftest.py
│   ├── test_protocol.py
│   ├── test_server.py
│   ├── test_security.py
│   ├── test_cache.py
│   └── test_zones.py
├── zones/
│   ├── example.local.zone
│   └── reverse.zone
├── deployment/
│   ├── dns-server.service
│   ├── config.env
│   ├── deploy.sh
│   └── Dockerfile
├── docs/
│   ├── README.md
│   ├── CONFIGURATION.md
│   ├── SECURITY.md
│   └── API.md
├── .github/
│   └── workflows/
│       └── ci.yml
├── .gitignore
├── docker-compose.yml
├── Jenkinsfile
├── Makefile
├── requirements.txt
├── requirements-dev.txt
└── pyproject.toml
```

### requirements.txt

```
# Core
aiohttp>=3.8.0

# Optional: DoH client
httpx>=0.24.0
```

### requirements-dev.txt

```
# Testing
pytest>=7.0.0
pytest-asyncio>=0.21.0
pytest-cov>=4.0.0

# Code quality
black>=23.0.0
isort>=5.12.0
flake8>=6.0.0
mypy>=1.0.0

# Docs
mkdocs>=1.4.0
mkdocs-material>=9.0.0
```

---

## Summary: Extended Phase 7

### Security (Day 31)
- ✅ Per-IP rate limiting
- ✅ Query validation
- ✅ Blocklist support
- ✅ DNS tunneling detection
- ✅ Response Rate Limiting (RRL)

### Zone Transfers (Day 32)
- ✅ AXFR (full transfer)
- ✅ IXFR (incremental)
- ✅ TCP transport
- ✅ Transfer client

### Metrics (Day 33)
- ✅ Prometheus format
- ✅ Query/response counters
- ✅ Latency histograms
- ✅ Cache statistics
- ✅ Health endpoints

### Utilities (Day 34)
- ✅ dig-like client
- ✅ Performance benchmarking
- ✅ Zone validation
- ✅ Server comparison

This completes the extended DNS server course!
