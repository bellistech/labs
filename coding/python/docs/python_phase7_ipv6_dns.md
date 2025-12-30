# Python Learning Path - Phase 7: IPv6 & DNS Server

Building a functional DNS server with full IPv6 support. This phase covers IPv6 addressing, DNS protocol internals, and building a production-ready DNS server from scratch.

**Duration:** Days 21-30
**Prerequisites:** Phases 1-6 completed
**Final Project:** Production DNS server with dual-stack support

---

## Table of Contents

1. [Day 21-22: IPv6 Fundamentals](#day-21-22-ipv6-fundamentals)
2. [Day 23-24: DNS Protocol Deep Dive](#day-23-24-dns-protocol-deep-dive)
3. [Day 25-27: Building the DNS Server](#day-25-27-building-the-dns-server)
4. [Day 28-29: Advanced Features](#day-28-29-advanced-features)
5. [Day 30: Production Deployment](#day-30-production-deployment)

---

## Day 21-22: IPv6 Fundamentals

### Exercise 21.1: IPv6 Address Handling

Create `day_21/ipv6_basics.py`:

```python
#!/usr/bin/env python3
"""
Day 21: IPv6 Fundamentals
Learning: IPv6 addressing, subnetting, address types
"""

import ipaddress
from typing import List, Tuple, Optional, Dict
from dataclasses import dataclass
from enum import Enum


class IPv6AddressType(Enum):
    """IPv6 address types."""
    GLOBAL_UNICAST = "Global Unicast (2000::/3)"
    LINK_LOCAL = "Link-Local (fe80::/10)"
    UNIQUE_LOCAL = "Unique Local (fc00::/7)"
    MULTICAST = "Multicast (ff00::/8)"
    LOOPBACK = "Loopback (::1)"
    UNSPECIFIED = "Unspecified (::)"
    IPV4_MAPPED = "IPv4-Mapped (::ffff:0:0/96)"
    DOCUMENTATION = "Documentation (2001:db8::/32)"
    OTHER = "Other"


@dataclass
class IPv6Info:
    """Detailed IPv6 address information."""
    address: str
    compressed: str
    exploded: str
    address_type: IPv6AddressType
    is_private: bool
    is_global: bool
    is_multicast: bool
    is_link_local: bool
    network: Optional[str] = None
    prefix_length: Optional[int] = None
    interface_id: Optional[str] = None


def classify_ipv6(addr: str) -> IPv6AddressType:
    """Classify an IPv6 address by type."""
    try:
        ip = ipaddress.IPv6Address(addr)
    except ipaddress.AddressValueError:
        return IPv6AddressType.OTHER
    
    if ip.is_loopback:
        return IPv6AddressType.LOOPBACK
    elif ip == ipaddress.IPv6Address("::"):
        return IPv6AddressType.UNSPECIFIED
    elif ip.is_link_local:
        return IPv6AddressType.LINK_LOCAL
    elif ip.is_multicast:
        return IPv6AddressType.MULTICAST
    elif ip.is_private:
        if int(ip) >> 120 in (0xfc, 0xfd):
            return IPv6AddressType.UNIQUE_LOCAL
        return IPv6AddressType.LINK_LOCAL
    elif ip.is_global:
        doc_net = ipaddress.IPv6Network("2001:db8::/32")
        if ip in doc_net:
            return IPv6AddressType.DOCUMENTATION
        return IPv6AddressType.GLOBAL_UNICAST
    elif ip.ipv4_mapped:
        return IPv6AddressType.IPV4_MAPPED
    
    return IPv6AddressType.OTHER


def analyze_ipv6(addr: str) -> IPv6Info:
    """Get detailed information about an IPv6 address."""
    if "/" in addr:
        network = ipaddress.IPv6Network(addr, strict=False)
        ip = network.network_address
        prefix_len = network.prefixlen
        net_str = str(network)
    else:
        ip = ipaddress.IPv6Address(addr)
        prefix_len = None
        net_str = None
    
    # Extract interface ID (last 64 bits)
    interface_id = format(int(ip) & ((1 << 64) - 1), '016x')
    interface_id = ':'.join(interface_id[i:i+4] for i in range(0, 16, 4))
    
    return IPv6Info(
        address=addr,
        compressed=ip.compressed,
        exploded=ip.exploded,
        address_type=classify_ipv6(str(ip)),
        is_private=ip.is_private,
        is_global=ip.is_global,
        is_multicast=ip.is_multicast,
        is_link_local=ip.is_link_local,
        network=net_str,
        prefix_length=prefix_len,
        interface_id=interface_id,
    )


def ipv6_to_ptr(addr: str) -> str:
    """Convert IPv6 address to PTR record format (ip6.arpa)."""
    ip = ipaddress.IPv6Address(addr)
    hex_str = ip.exploded.replace(":", "")
    reversed_hex = ".".join(reversed(hex_str))
    return f"{reversed_hex}.ip6.arpa"


def ptr_to_ipv6(ptr: str) -> str:
    """Convert PTR record back to IPv6 address."""
    ptr = ptr.replace(".ip6.arpa", "")
    hex_digits = ptr.split(".")
    hex_digits.reverse()
    hex_str = "".join(hex_digits)
    addr_str = ":".join(hex_str[i:i+4] for i in range(0, 32, 4))
    return ipaddress.IPv6Address(addr_str).compressed


def generate_eui64(mac: str, prefix: str = "fe80::/64") -> str:
    """
    Generate EUI-64 interface identifier from MAC address.
    Used for SLAAC (Stateless Address Autoconfiguration).
    """
    mac = mac.replace(":", "").replace("-", "").lower()
    if len(mac) != 12:
        raise ValueError("Invalid MAC address")
    
    oui = mac[:6]
    nic = mac[6:]
    eui64 = oui + "fffe" + nic
    
    # Flip the 7th bit (Universal/Local bit)
    first_byte = int(eui64[:2], 16) ^ 0x02
    eui64 = format(first_byte, '02x') + eui64[2:]
    
    interface_id = ":".join(eui64[i:i+4] for i in range(0, 16, 4))
    
    network = ipaddress.IPv6Network(prefix, strict=False)
    prefix_int = int(network.network_address) >> 64 << 64
    interface_int = int(eui64, 16)
    full_addr = ipaddress.IPv6Address(prefix_int | interface_int)
    
    return full_addr.compressed


def subnet_ipv6(network: str, new_prefix: int) -> List[str]:
    """Subnet an IPv6 network into smaller networks."""
    net = ipaddress.IPv6Network(network, strict=False)
    
    if new_prefix <= net.prefixlen:
        raise ValueError(f"New prefix must be larger than {net.prefixlen}")
    
    subnets = list(net.subnets(new_prefix=new_prefix))
    return [str(s) for s in subnets]


def summarize_networks(networks: List[str]) -> List[str]:
    """Summarize/aggregate a list of IPv6 networks."""
    network_objs = [ipaddress.IPv6Network(n, strict=False) for n in networks]
    collapsed = ipaddress.collapse_addresses(network_objs)
    return [str(n) for n in collapsed]


def ipv4_to_ipv6_mapped(ipv4: str) -> str:
    """Convert IPv4 to IPv4-mapped IPv6 address."""
    v4 = ipaddress.IPv4Address(ipv4)
    v6_int = 0xffff << 32 | int(v4)
    return ipaddress.IPv6Address(v6_int).compressed


class IPv6Calculator:
    """IPv6 subnet calculator."""
    
    def __init__(self, network: str):
        self.network = ipaddress.IPv6Network(network, strict=False)
    
    @property
    def total_addresses(self) -> int:
        return self.network.num_addresses
    
    @property
    def first_address(self) -> str:
        return str(self.network.network_address)
    
    @property
    def last_address(self) -> str:
        return str(self.network.broadcast_address)
    
    def get_nth_address(self, n: int) -> str:
        if n >= self.total_addresses:
            raise ValueError(f"Index {n} exceeds network size")
        return str(self.network.network_address + n)
    
    def contains(self, addr: str) -> bool:
        try:
            ip = ipaddress.IPv6Address(addr)
            return ip in self.network
        except:
            return False
    
    def summary(self) -> Dict:
        return {
            "network": str(self.network),
            "prefix_length": self.network.prefixlen,
            "first_address": self.first_address,
            "last_address": self.last_address,
            "total_addresses": self.total_addresses,
        }


def main():
    print("=== IPv6 Fundamentals ===\n")
    
    # Address analysis
    print("1. Address Analysis:")
    addresses = [
        "2001:db8:85a3::8a2e:370:7334",
        "fe80::1",
        "::1",
        "ff02::1",
        "2607:f8b0:4004:800::200e",
    ]
    
    for addr in addresses:
        info = analyze_ipv6(addr)
        print(f"   {addr}")
        print(f"      Type: {info.address_type.value}")
        print(f"      Global: {info.is_global}")
    
    # PTR conversion
    print("\n2. PTR Records:")
    addr = "2001:db8::1"
    ptr = ipv6_to_ptr(addr)
    print(f"   {addr} -> {ptr[:50]}...")
    
    # EUI-64
    print("\n3. EUI-64 Generation:")
    mac = "00:1a:2b:3c:4d:5e"
    link_local = generate_eui64(mac)
    print(f"   MAC: {mac}")
    print(f"   Link-Local: {link_local}")
    
    # Subnetting
    print("\n4. Subnetting:")
    network = "2001:db8::/48"
    subnets = subnet_ipv6(network, 64)[:4]
    print(f"   {network} -> /64 subnets:")
    for s in subnets:
        print(f"      {s}")
    
    # Calculator
    print("\n5. Calculator:")
    calc = IPv6Calculator("2001:db8::/32")
    print(f"   Network: {calc.network}")
    print(f"   Addresses: {calc.total_addresses:,}")


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add day_21/
git commit -m "Day 21: IPv6 fundamentals"
```

---

## Day 23-24: DNS Protocol Deep Dive

### Exercise 23.1: DNS Protocol Implementation

Create `day_23/dns_protocol.py`:

```python
#!/usr/bin/env python3
"""
Day 23: DNS Protocol Implementation
Learning: DNS message format, record types, wire format
Reference: RFC 1035, RFC 3596 (AAAA records)
"""

import struct
import random
from enum import IntEnum
from dataclasses import dataclass, field
from typing import List, Optional, Tuple
import ipaddress


class DNSRecordType(IntEnum):
    """DNS record types."""
    A = 1
    NS = 2
    CNAME = 5
    SOA = 6
    PTR = 12
    MX = 15
    TXT = 16
    AAAA = 28
    SRV = 33
    ANY = 255


class DNSClass(IntEnum):
    """DNS classes."""
    IN = 1
    ANY = 255


class DNSRcode(IntEnum):
    """DNS response codes."""
    NOERROR = 0
    FORMERR = 1
    SERVFAIL = 2
    NXDOMAIN = 3
    NOTIMP = 4
    REFUSED = 5


@dataclass
class DNSHeader:
    """DNS message header (12 bytes)."""
    id: int = 0
    qr: int = 0           # 0=query, 1=response
    opcode: int = 0
    aa: int = 0           # Authoritative Answer
    tc: int = 0           # Truncation
    rd: int = 1           # Recursion Desired
    ra: int = 0           # Recursion Available
    z: int = 0
    rcode: int = 0
    qdcount: int = 0
    ancount: int = 0
    nscount: int = 0
    arcount: int = 0
    
    def pack(self) -> bytes:
        flags = (
            (self.qr << 15) | (self.opcode << 11) |
            (self.aa << 10) | (self.tc << 9) |
            (self.rd << 8) | (self.ra << 7) |
            (self.z << 4) | self.rcode
        )
        return struct.pack(
            "!HHHHHH", self.id, flags,
            self.qdcount, self.ancount, self.nscount, self.arcount
        )
    
    @classmethod
    def unpack(cls, data: bytes) -> 'DNSHeader':
        id_, flags, qd, an, ns, ar = struct.unpack("!HHHHHH", data[:12])
        return cls(
            id=id_,
            qr=(flags >> 15) & 1,
            opcode=(flags >> 11) & 0xf,
            aa=(flags >> 10) & 1,
            tc=(flags >> 9) & 1,
            rd=(flags >> 8) & 1,
            ra=(flags >> 7) & 1,
            z=(flags >> 4) & 7,
            rcode=flags & 0xf,
            qdcount=qd, ancount=an, nscount=ns, arcount=ar
        )


@dataclass
class DNSQuestion:
    """DNS question section."""
    name: str
    qtype: int = DNSRecordType.A
    qclass: int = DNSClass.IN
    
    def pack(self) -> bytes:
        return encode_name(self.name) + struct.pack("!HH", self.qtype, self.qclass)


@dataclass
class DNSRecord:
    """DNS resource record."""
    name: str
    rtype: int
    rclass: int = DNSClass.IN
    ttl: int = 300
    rdata: bytes = b""
    parsed_data: Optional[str] = None
    
    def pack(self) -> bytes:
        return (
            encode_name(self.name) +
            struct.pack("!HHIH", self.rtype, self.rclass, self.ttl, len(self.rdata)) +
            self.rdata
        )


@dataclass
class DNSMessage:
    """Complete DNS message."""
    header: DNSHeader = field(default_factory=DNSHeader)
    questions: List[DNSQuestion] = field(default_factory=list)
    answers: List[DNSRecord] = field(default_factory=list)
    authority: List[DNSRecord] = field(default_factory=list)
    additional: List[DNSRecord] = field(default_factory=list)
    
    def pack(self) -> bytes:
        self.header.qdcount = len(self.questions)
        self.header.ancount = len(self.answers)
        self.header.nscount = len(self.authority)
        self.header.arcount = len(self.additional)
        
        data = self.header.pack()
        for q in self.questions:
            data += q.pack()
        for rr in self.answers + self.authority + self.additional:
            data += rr.pack()
        return data
    
    @classmethod
    def unpack(cls, data: bytes) -> 'DNSMessage':
        header = DNSHeader.unpack(data[:12])
        offset = 12
        
        questions = []
        for _ in range(header.qdcount):
            name, offset = decode_name(data, offset)
            qtype, qclass = struct.unpack("!HH", data[offset:offset+4])
            offset += 4
            questions.append(DNSQuestion(name=name, qtype=qtype, qclass=qclass))
        
        def parse_records(count):
            nonlocal offset
            records = []
            for _ in range(count):
                name, offset = decode_name(data, offset)
                rtype, rclass, ttl, rdlen = struct.unpack("!HHIH", data[offset:offset+10])
                offset += 10
                rdata = data[offset:offset+rdlen]
                offset += rdlen
                parsed = parse_rdata(rtype, rdata, data)
                records.append(DNSRecord(
                    name=name, rtype=rtype, rclass=rclass,
                    ttl=ttl, rdata=rdata, parsed_data=parsed
                ))
            return records
        
        return cls(
            header=header,
            questions=questions,
            answers=parse_records(header.ancount),
            authority=parse_records(header.nscount),
            additional=parse_records(header.arcount),
        )


def encode_name(name: str) -> bytes:
    """Encode domain name to DNS wire format."""
    if not name or name == ".":
        return b"\x00"
    
    result = b""
    for label in name.rstrip(".").split("."):
        result += bytes([len(label)]) + label.encode("ascii")
    return result + b"\x00"


def decode_name(data: bytes, offset: int) -> Tuple[str, int]:
    """Decode domain name from DNS wire format with compression."""
    labels = []
    original_offset = offset
    jumped = False
    
    while True:
        if offset >= len(data):
            break
        
        length = data[offset]
        
        # Compression pointer
        if (length & 0xc0) == 0xc0:
            if not jumped:
                original_offset = offset + 2
            pointer = struct.unpack("!H", data[offset:offset+2])[0] & 0x3fff
            offset = pointer
            jumped = True
            continue
        
        offset += 1
        if length == 0:
            break
        
        labels.append(data[offset:offset+length].decode("ascii"))
        offset += length
    
    return ".".join(labels) if labels else ".", original_offset if not jumped else original_offset


def parse_rdata(rtype: int, rdata: bytes, full_msg: bytes = None) -> str:
    """Parse rdata based on record type."""
    try:
        if rtype == DNSRecordType.A:
            return str(ipaddress.IPv4Address(rdata))
        elif rtype == DNSRecordType.AAAA:
            return str(ipaddress.IPv6Address(rdata))
        elif rtype == DNSRecordType.TXT:
            texts = []
            offset = 0
            while offset < len(rdata):
                length = rdata[offset]
                offset += 1
                texts.append(rdata[offset:offset+length].decode("utf-8", errors="replace"))
                offset += length
            return " ".join(texts)
        else:
            return rdata.hex()
    except Exception as e:
        return f"(error: {e})"


# Record creation helpers

def create_a_record(name: str, ip: str, ttl: int = 300) -> DNSRecord:
    """Create an A record."""
    rdata = ipaddress.IPv4Address(ip).packed
    return DNSRecord(name=name, rtype=DNSRecordType.A, ttl=ttl, rdata=rdata, parsed_data=ip)


def create_aaaa_record(name: str, ip: str, ttl: int = 300) -> DNSRecord:
    """Create an AAAA record."""
    rdata = ipaddress.IPv6Address(ip).packed
    return DNSRecord(name=name, rtype=DNSRecordType.AAAA, ttl=ttl, rdata=rdata, parsed_data=ip)


def create_cname_record(name: str, target: str, ttl: int = 300) -> DNSRecord:
    """Create a CNAME record."""
    rdata = encode_name(target)
    return DNSRecord(name=name, rtype=DNSRecordType.CNAME, ttl=ttl, rdata=rdata, parsed_data=target)


def create_txt_record(name: str, text: str, ttl: int = 300) -> DNSRecord:
    """Create a TXT record."""
    encoded = text.encode("utf-8")
    rdata = bytes([len(encoded)]) + encoded
    return DNSRecord(name=name, rtype=DNSRecordType.TXT, ttl=ttl, rdata=rdata, parsed_data=text)


def create_ptr_record(name: str, target: str, ttl: int = 300) -> DNSRecord:
    """Create a PTR record."""
    rdata = encode_name(target)
    return DNSRecord(name=name, rtype=DNSRecordType.PTR, ttl=ttl, rdata=rdata, parsed_data=target)


def create_query(name: str, qtype: int = DNSRecordType.A) -> DNSMessage:
    """Create a DNS query message."""
    return DNSMessage(
        header=DNSHeader(id=random.randint(0, 65535), rd=1),
        questions=[DNSQuestion(name=name, qtype=qtype)]
    )


def create_response(query: DNSMessage, answers: List[DNSRecord],
                    rcode: int = DNSRcode.NOERROR) -> DNSMessage:
    """Create a DNS response message."""
    return DNSMessage(
        header=DNSHeader(
            id=query.header.id, qr=1, rd=query.header.rd, ra=1, rcode=rcode
        ),
        questions=query.questions,
        answers=answers,
    )


def main():
    print("=== DNS Protocol Demo ===\n")
    
    # Encode/decode domain names
    print("1. Domain Name Encoding:")
    name = "www.example.com"
    encoded = encode_name(name)
    decoded, _ = decode_name(encoded, 0)
    print(f"   Original: {name}")
    print(f"   Encoded: {encoded.hex()}")
    print(f"   Decoded: {decoded}")
    
    # Create query
    print("\n2. DNS Query:")
    query = create_query("google.com", DNSRecordType.A)
    query_bytes = query.pack()
    print(f"   Query ID: {query.header.id}")
    print(f"   Size: {len(query_bytes)} bytes")
    
    # Create records
    print("\n3. DNS Records:")
    records = [
        create_a_record("example.com", "93.184.216.34"),
        create_aaaa_record("example.com", "2606:2800:220:1:248:1893:25c8:1946"),
        create_txt_record("example.com", "v=spf1 -all"),
    ]
    for r in records:
        print(f"   {r.name} {DNSRecordType(r.rtype).name} {r.parsed_data}")


if __name__ == "__main__":
    main()
```

**Git checkpoint:**
```bash
git add day_23/
git commit -m "Day 23: DNS protocol implementation"
```

---

## Day 25-27: Building the DNS Server

### Exercise 25.1: DNS Server

Create `day_25/dns_server.py`:

```python
#!/usr/bin/env python3
"""
Day 25-27: Production DNS Server
Features: Dual-stack, zones, caching, recursive resolution
"""

import asyncio
import socket
import logging
import time
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, field
from collections import defaultdict
import ipaddress

import sys
sys.path.insert(0, '..')
from day_23.dns_protocol import (
    DNSMessage, DNSHeader, DNSQuestion, DNSRecord,
    DNSRecordType, DNSClass, DNSRcode,
    create_a_record, create_aaaa_record, create_cname_record,
    create_txt_record, create_response, create_query,
    encode_name
)

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger("dns")


@dataclass
class CachedRecord:
    """Cached DNS record with expiration."""
    record: DNSRecord
    expires_at: float
    
    @property
    def is_expired(self) -> bool:
        return time.time() > self.expires_at
    
    @property
    def remaining_ttl(self) -> int:
        return max(0, int(self.expires_at - time.time()))


@dataclass
class DNSZone:
    """DNS zone with records."""
    name: str
    records: Dict[str, Dict[int, List[DNSRecord]]] = field(default_factory=dict)
    
    def add_record(self, record: DNSRecord) -> None:
        name = record.name.lower().rstrip(".")
        if name not in self.records:
            self.records[name] = {}
        if record.rtype not in self.records[name]:
            self.records[name][record.rtype] = []
        self.records[name][record.rtype].append(record)
    
    def get_records(self, name: str, rtype: int) -> List[DNSRecord]:
        name = name.lower().rstrip(".")
        if name in self.records:
            if rtype == DNSRecordType.ANY:
                return [r for recs in self.records[name].values() for r in recs]
            return self.records[name].get(rtype, [])
        return []


class DNSCache:
    """DNS record cache with TTL."""
    
    def __init__(self, max_size: int = 10000):
        self.cache: Dict[Tuple[str, int], List[CachedRecord]] = {}
        self.max_size = max_size
    
    def get(self, name: str, rtype: int) -> List[DNSRecord]:
        key = (name.lower(), rtype)
        if key not in self.cache:
            return []
        
        valid = []
        for cached in self.cache[key]:
            if not cached.is_expired:
                record = cached.record
                record.ttl = cached.remaining_ttl
                valid.append(record)
        
        if not valid:
            del self.cache[key]
        return valid
    
    def put(self, records: List[DNSRecord]) -> None:
        for record in records:
            key = (record.name.lower(), record.rtype)
            expires_at = time.time() + record.ttl
            cached = CachedRecord(record=record, expires_at=expires_at)
            
            if key not in self.cache:
                self.cache[key] = []
            self.cache[key].append(cached)


@dataclass
class ServerConfig:
    """DNS server configuration."""
    listen_address: str = "::"
    listen_port: int = 5353
    upstream_servers: List[str] = field(default_factory=lambda: ["8.8.8.8", "1.1.1.1"])
    upstream_v6: List[str] = field(default_factory=lambda: ["2001:4860:4860::8888"])
    enable_recursion: bool = True
    log_queries: bool = True


class DNSServerProtocol(asyncio.DatagramProtocol):
    """DNS server UDP protocol."""
    
    def __init__(self, server: 'DNSServer'):
        self.server = server
        self.transport = None
    
    def connection_made(self, transport):
        self.transport = transport
    
    def datagram_received(self, data: bytes, addr: Tuple):
        asyncio.create_task(self._handle(data, addr))
    
    async def _handle(self, data: bytes, addr: Tuple):
        try:
            query = DNSMessage.unpack(data)
            response = await self.server.handle_query(query, addr)
            self.transport.sendto(response.pack(), addr)
        except Exception as e:
            logger.error(f"Error from {addr}: {e}")


class DNSServer:
    """
    DNS server with IPv6 support.
    
    Features:
    - Dual-stack listening
    - Zone support
    - Record caching
    - Recursive resolution
    """
    
    def __init__(self, config: ServerConfig = None):
        self.config = config or ServerConfig()
        self.zones: Dict[str, DNSZone] = {}
        self.cache = DNSCache()
        self.metrics = defaultdict(int)
        self._transports = []
    
    def add_zone(self, zone: DNSZone) -> None:
        self.zones[zone.name.lower()] = zone
        logger.info(f"Added zone: {zone.name}")
    
    async def start(self) -> None:
        loop = asyncio.get_event_loop()
        
        # IPv6 socket
        try:
            t6, _ = await loop.create_datagram_endpoint(
                lambda: DNSServerProtocol(self),
                local_addr=(self.config.listen_address, self.config.listen_port),
                family=socket.AF_INET6,
            )
            self._transports.append(t6)
            logger.info(f"Listening on [::]:{self.config.listen_port} (IPv6)")
        except Exception as e:
            logger.warning(f"IPv6 bind failed: {e}")
        
        # IPv4 socket
        try:
            t4, _ = await loop.create_datagram_endpoint(
                lambda: DNSServerProtocol(self),
                local_addr=("0.0.0.0", self.config.listen_port),
                family=socket.AF_INET,
            )
            self._transports.append(t4)
            logger.info(f"Listening on 0.0.0.0:{self.config.listen_port} (IPv4)")
        except Exception as e:
            logger.warning(f"IPv4 bind failed: {e}")
    
    async def stop(self) -> None:
        for t in self._transports:
            t.close()
        self._transports.clear()
        logger.info("Server stopped")
    
    async def handle_query(self, query: DNSMessage, addr: Tuple) -> DNSMessage:
        self.metrics["queries"] += 1
        
        if not query.questions:
            return self._error_response(query, DNSRcode.FORMERR)
        
        q = query.questions[0]
        
        if self.config.log_queries:
            logger.info(f"Query from {addr[0]}: {q.name} {DNSRecordType(q.qtype).name}")
        
        # Check local zones
        answers = self._lookup_local(q.name, q.qtype)
        if answers:
            self.metrics["local"] += 1
            return create_response(query, answers)
        
        # Check cache
        cached = self.cache.get(q.name, q.qtype)
        if cached:
            self.metrics["cached"] += 1
            return create_response(query, cached)
        
        # Recursive resolution
        if self.config.enable_recursion and query.header.rd:
            try:
                answers = await self._resolve(q.name, q.qtype)
                if answers:
                    self.metrics["recursive"] += 1
                    self.cache.put(answers)
                    return create_response(query, answers)
            except Exception as e:
                logger.error(f"Recursive failed: {e}")
        
        self.metrics["nxdomain"] += 1
        return self._error_response(query, DNSRcode.NXDOMAIN)
    
    def _lookup_local(self, name: str, rtype: int) -> List[DNSRecord]:
        name_lower = name.lower().rstrip(".")
        for zone_name, zone in self.zones.items():
            if name_lower == zone_name or name_lower.endswith("." + zone_name):
                records = zone.get_records(name_lower, rtype)
                if records:
                    return records
                
                # Handle CNAME
                if rtype != DNSRecordType.CNAME:
                    cnames = zone.get_records(name_lower, DNSRecordType.CNAME)
                    if cnames:
                        target = cnames[0].parsed_data
                        target_records = zone.get_records(target, rtype)
                        return cnames + target_records
        return []
    
    async def _resolve(self, name: str, rtype: int) -> List[DNSRecord]:
        query = create_query(name, rtype)
        
        for server in self.config.upstream_servers:
            try:
                response = await self._upstream_query(query, server)
                if response and response.answers:
                    return response.answers
            except:
                continue
        
        for server in self.config.upstream_v6:
            try:
                response = await self._upstream_query(query, server)
                if response and response.answers:
                    return response.answers
            except:
                continue
        
        return []
    
    async def _upstream_query(self, query: DNSMessage, server: str) -> Optional[DNSMessage]:
        loop = asyncio.get_event_loop()
        
        try:
            ipaddress.IPv6Address(server)
            family = socket.AF_INET6
        except:
            family = socket.AF_INET
        
        sock = socket.socket(family, socket.SOCK_DGRAM)
        sock.setblocking(False)
        
        try:
            await loop.sock_sendto(sock, query.pack(), (server, 53))
            data = await asyncio.wait_for(loop.sock_recv(sock, 4096), timeout=2.0)
            return DNSMessage.unpack(data)
        except asyncio.TimeoutError:
            return None
        finally:
            sock.close()
    
    def _error_response(self, query: DNSMessage, rcode: int) -> DNSMessage:
        return create_response(query, [], rcode=rcode)


def create_example_zone() -> DNSZone:
    """Create example zone for testing."""
    zone = DNSZone(name="example.local")
    
    # A records
    zone.add_record(create_a_record("example.local", "192.168.1.1"))
    zone.add_record(create_a_record("www.example.local", "192.168.1.10"))
    zone.add_record(create_a_record("mail.example.local", "192.168.1.20"))
    zone.add_record(create_a_record("ns1.example.local", "192.168.1.2"))
    
    # AAAA records
    zone.add_record(create_aaaa_record("example.local", "2001:db8::1"))
    zone.add_record(create_aaaa_record("www.example.local", "2001:db8::10"))
    zone.add_record(create_aaaa_record("mail.example.local", "2001:db8::20"))
    zone.add_record(create_aaaa_record("ipv6only.example.local", "2001:db8::100"))
    
    # CNAME
    zone.add_record(create_cname_record("blog.example.local", "www.example.local"))
    zone.add_record(create_cname_record("ftp.example.local", "www.example.local"))
    
    # TXT
    zone.add_record(create_txt_record("example.local", "v=spf1 mx -all"))
    zone.add_record(create_txt_record("_dmarc.example.local", "v=DMARC1; p=reject"))
    
    return zone


async def main():
    print("=== DNS Server ===\n")
    
    config = ServerConfig(listen_port=5353, log_queries=True)
    server = DNSServer(config)
    server.add_zone(create_example_zone())
    
    print(f"Starting on port {config.listen_port}")
    print(f"Zones: {list(server.zones.keys())}")
    print("\nTest with:")
    print(f"  dig @localhost -p {config.listen_port} www.example.local A")
    print(f"  dig @localhost -p {config.listen_port} www.example.local AAAA")
    print(f"  dig @localhost -p {config.listen_port} google.com A  # recursive")
    print("\nPress Ctrl+C to stop\n")
    
    await server.start()
    
    try:
        while True:
            await asyncio.sleep(60)
            logger.info(f"Metrics: {dict(server.metrics)}")
    except KeyboardInterrupt:
        print("\nStopping...")
    finally:
        await server.stop()


if __name__ == "__main__":
    asyncio.run(main())
```

**Git checkpoint:**
```bash
git add day_25/
git commit -m "Day 25-27: DNS server with IPv6 support"
```

---

## Day 28-29: Advanced Features

### DNS-over-HTTPS Client

Create `day_28/doh_client.py`:

```python
#!/usr/bin/env python3
"""
Day 28: DNS-over-HTTPS (DoH) Client
Modern encrypted DNS resolution
"""

import asyncio
import base64
from typing import List, Dict, Optional
from dataclasses import dataclass

import sys
sys.path.insert(0, '..')
from day_23.dns_protocol import DNSMessage, DNSRecordType, create_query


DOH_PROVIDERS = {
    "cloudflare": "https://cloudflare-dns.com/dns-query",
    "google": "https://dns.google/dns-query",
    "quad9": "https://dns.quad9.net/dns-query",
}


@dataclass
class DoHResponse:
    """DoH response."""
    status: int
    answers: List[Dict]
    query_time_ms: float = 0


class DoHClient:
    """DNS-over-HTTPS client."""
    
    def __init__(self, provider: str = "cloudflare"):
        self.endpoint = DOH_PROVIDERS.get(provider, provider)
    
    async def query_json(self, name: str, rtype: str = "A") -> DoHResponse:
        """Query using JSON API."""
        try:
            import aiohttp
        except ImportError:
            raise ImportError("pip install aiohttp")
        
        import time
        start = time.perf_counter()
        
        async with aiohttp.ClientSession() as session:
            async with session.get(
                self.endpoint,
                params={"name": name, "type": rtype},
                headers={"Accept": "application/dns-json"},
            ) as resp:
                data = await resp.json()
        
        elapsed = (time.perf_counter() - start) * 1000
        
        return DoHResponse(
            status=data.get("Status", 0),
            answers=data.get("Answer", []),
            query_time_ms=round(elapsed, 2),
        )
    
    async def query_wire(self, name: str, rtype: int = DNSRecordType.A) -> DNSMessage:
        """Query using wire format (RFC 8484)."""
        try:
            import aiohttp
        except ImportError:
            raise ImportError("pip install aiohttp")
        
        query = create_query(name, rtype)
        
        async with aiohttp.ClientSession() as session:
            async with session.post(
                self.endpoint,
                data=query.pack(),
                headers={
                    "Content-Type": "application/dns-message",
                    "Accept": "application/dns-message",
                },
            ) as resp:
                data = await resp.read()
        
        return DNSMessage.unpack(data)


class DoHResolver:
    """DoH resolver with fallback."""
    
    def __init__(self, providers: List[str] = None):
        self.providers = providers or ["cloudflare", "google"]
        self.clients = [DoHClient(p) for p in self.providers]
    
    async def resolve(self, name: str, rtype: str = "A") -> List[str]:
        for client in self.clients:
            try:
                response = await client.query_json(name, rtype)
                if response.status == 0 and response.answers:
                    return [a.get("data") for a in response.answers]
            except:
                continue
        return []
    
    async def resolve_both(self, name: str) -> Dict[str, List[str]]:
        """Resolve both A and AAAA."""
        a, aaaa = await asyncio.gather(
            self.resolve(name, "A"),
            self.resolve(name, "AAAA")
        )
        return {"A": a, "AAAA": aaaa}


async def main():
    print("=== DoH Client Demo ===\n")
    
    client = DoHClient("cloudflare")
    
    print("1. JSON API Query:")
    try:
        resp = await client.query_json("google.com", "A")
        print(f"   Status: {resp.status}, Time: {resp.query_time_ms}ms")
        for a in resp.answers:
            print(f"   {a.get('name')} -> {a.get('data')}")
    except Exception as e:
        print(f"   Error: {e}")
    
    print("\n2. AAAA Query:")
    try:
        resp = await client.query_json("google.com", "AAAA")
        for a in resp.answers:
            print(f"   {a.get('name')} -> {a.get('data')}")
    except Exception as e:
        print(f"   Error: {e}")
    
    print("\n3. Dual-stack Resolution:")
    resolver = DoHResolver()
    for domain in ["github.com", "cloudflare.com"]:
        results = await resolver.resolve_both(domain)
        print(f"   {domain}:")
        print(f"      A: {results['A']}")
        print(f"      AAAA: {results['AAAA']}")


if __name__ == "__main__":
    asyncio.run(main())
```

**Git checkpoint:**
```bash
git add day_28/
git commit -m "Day 28: DNS-over-HTTPS client"
```

---

## Day 30: Production Deployment

### Systemd Service

Create `deployment/dns-server.service`:

```ini
[Unit]
Description=Python DNS Server
After=network.target

[Service]
Type=simple
User=dns
Group=dns
WorkingDirectory=/opt/dns-server
Environment="PATH=/opt/dns-server/venv/bin"
EnvironmentFile=/etc/dns-server/config.env
ExecStart=/opt/dns-server/venv/bin/python -m dns_server
Restart=always
RestartSec=5

# Allow privileged port binding
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Security
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### Configuration

Create `deployment/config.env`:

```bash
# Network
DNS_LISTEN_PORT=53
DNS_ENABLE_IPV6=true

# Upstream
DNS_UPSTREAM_SERVERS=8.8.8.8,1.1.1.1
DNS_UPSTREAM_V6=2001:4860:4860::8888

# Features
DNS_ENABLE_RECURSION=true
DNS_CACHE_SIZE=10000

# Logging
DNS_LOG_QUERIES=true
DNS_LOG_LEVEL=INFO
```

### Zone File Example

Create `zones/example.local.zone`:

```
$TTL 3600
$ORIGIN example.local.

@       IN  SOA   ns1 admin (2024010101 3600 900 604800 86400)
@       IN  NS    ns1
@       IN  NS    ns2

; IPv4
@       IN  A     192.168.1.1
ns1     IN  A     192.168.1.2
ns2     IN  A     192.168.1.3
www     IN  A     192.168.1.10
mail    IN  A     192.168.1.20

; IPv6
@       IN  AAAA  2001:db8::1
ns1     IN  AAAA  2001:db8::2
ns2     IN  AAAA  2001:db8::3
www     IN  AAAA  2001:db8::10
mail    IN  AAAA  2001:db8::20
ipv6    IN  AAAA  2001:db8::100

; Aliases
blog    IN  CNAME www
ftp     IN  CNAME www

; Mail
@       IN  MX    10 mail

; TXT
@       IN  TXT   "v=spf1 mx -all"
```

### Dockerfile

```dockerfile
FROM python:3.11-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY dns_server/ ./dns_server/
COPY zones/ ./zones/

RUN useradd --system dns && chown -R dns:dns /app
USER dns

EXPOSE 53/udp 53/tcp

HEALTHCHECK --interval=30s --timeout=5s \
    CMD python -c "import socket; s=socket.socket(socket.AF_INET,socket.SOCK_DGRAM); s.settimeout(1); s.sendto(b'',('127.0.0.1',53))"

CMD ["python", "-m", "dns_server"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  dns:
    build: .
    ports:
      - "53:53/udp"
      - "53:53/tcp"
    volumes:
      - ./zones:/app/zones:ro
    environment:
      - DNS_LOG_QUERIES=true
      - DNS_ENABLE_RECURSION=true
    restart: unless-stopped
```

---

## Testing

### Manual Tests

```bash
# A record
dig @localhost -p 5353 www.example.local A

# AAAA record
dig @localhost -p 5353 www.example.local AAAA

# Both
dig @localhost -p 5353 www.example.local ANY

# Recursive (external)
dig @localhost -p 5353 google.com A

# IPv6 transport
dig @::1 -p 5353 www.example.local AAAA

# TXT record
dig @localhost -p 5353 example.local TXT
```

### pytest Tests

Create `tests/test_dns.py`:

```python
import pytest
import asyncio
import socket
from day_23.dns_protocol import DNSMessage, DNSRecordType, create_query
from day_25.dns_server import DNSServer, ServerConfig, create_example_zone


@pytest.fixture
async def server():
    config = ServerConfig(listen_port=15353, enable_recursion=False)
    srv = DNSServer(config)
    srv.add_zone(create_example_zone())
    await srv.start()
    yield srv
    await srv.stop()


async def query(name: str, rtype: int, port: int = 15353) -> DNSMessage:
    q = create_query(name, rtype)
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.setblocking(False)
    loop = asyncio.get_event_loop()
    
    try:
        await loop.sock_sendto(sock, q.pack(), ("127.0.0.1", port))
        data = await asyncio.wait_for(loop.sock_recv(sock, 4096), timeout=2)
        return DNSMessage.unpack(data)
    finally:
        sock.close()


@pytest.mark.asyncio
async def test_a_record(server):
    resp = await query("www.example.local", DNSRecordType.A)
    assert resp.header.rcode == 0
    assert len(resp.answers) > 0
    assert resp.answers[0].parsed_data == "192.168.1.10"


@pytest.mark.asyncio
async def test_aaaa_record(server):
    resp = await query("www.example.local", DNSRecordType.AAAA)
    assert resp.header.rcode == 0
    assert "2001:db8::10" in resp.answers[0].parsed_data


@pytest.mark.asyncio
async def test_nxdomain(server):
    resp = await query("nonexistent.example.local", DNSRecordType.A)
    assert resp.header.rcode == 3  # NXDOMAIN
```

Run: `pytest tests/test_dns.py -v`

---

## Summary: Phase 7 Learning Outcomes

### IPv6 Skills
- ✅ Address types and classification
- ✅ EUI-64 interface ID generation  
- ✅ Subnetting and summarization
- ✅ PTR records (ip6.arpa)
- ✅ Dual-stack networking

### DNS Protocol
- ✅ Message format (RFC 1035)
- ✅ Record types (A, AAAA, CNAME, TXT, PTR, MX)
- ✅ Wire format encoding/decoding
- ✅ Name compression

### Server Development
- ✅ Async UDP server
- ✅ Dual-stack sockets
- ✅ Zone management
- ✅ Caching with TTL
- ✅ Recursive resolution

### Modern DNS
- ✅ DNS-over-HTTPS (DoH)
- ✅ Multiple providers

### Production
- ✅ Systemd service
- ✅ Docker deployment
- ✅ Testing

---

## Next Steps

1. Add DNS-over-TLS (DoT) - port 853
2. Add DNSSEC validation
3. Add rate limiting
4. Add Prometheus metrics
5. Add zone transfers (AXFR)

This completes Phase 7 of your Python learning path!
