# IPv6 Home Lab Guide
## Comprehensive Setup, Security, and Automation

---

## Table of Contents

1. [IPv6 Fundamentals](#ipv6-fundamentals)
2. [Getting IPv6 from Your ISP](#getting-ipv6-from-your-isp)
3. [Network Design](#network-design)
4. [Router Configuration](#router-configuration)
5. [Addressing Strategy](#addressing-strategy)
6. [Security and Firewall Rules](#security-and-firewall-rules)
7. [Split-Horizon DNS](#split-horizon-dns)
8. [Advanced Firewall Hardening](#advanced-firewall-hardening)
9. [Automating Prefix Delegation](#automating-prefix-delegation)
10. [Monitoring and Visibility](#monitoring-and-visibility)
11. [Common Gotchas](#common-gotchas)
12. [Troubleshooting Techniques](#troubleshooting-techniques)

---

## IPv6 Fundamentals

### Address Structure

IPv6 uses 128-bit addresses written in eight groups of four hexadecimal digits separated by colons:

```
2001:0db8:85a3:0000:0000:8a2e:0370:7334
```

**Compression Rules:**

You can compress consecutive zeros with `::` (only once per address):
```
2001:db8:85a3::8a2e:370:7334
```

Leading zeros in each group can be omitted:
```
2001:0db8 → 2001:db8
```

### Address Types

**Global Unicast (`2000::/3`)**: Routable internet addresses, analogous to public IPv4

**Link-Local (`fe80::/10`)**: Auto-configured on every interface, not routable beyond the local link. These always start with `fe80::` and you'll see them constantly in network troubleshooting

**Unique Local (`fc00::/7`)**: Private addressing, similar to RFC 1918 in IPv4. In practice, you'll see `fd00::/8` used

**Multicast (`ff00::/8`)**: Replaces broadcast in IPv4. Common ones:
- `ff02::1` - All nodes on link
- `ff02::2` - All routers on link
- `ff02::1:ff00:0/104` - Solicited-node multicast

### SLAAC and Address Assignment

IPv6 hosts can autoconfigure using SLAAC (Stateless Address Autoconfiguration):

1. Host generates link-local address using `fe80::` + interface ID
2. Performs DAD (Duplicate Address Detection) via Neighbor Solicitation
3. Listens for Router Advertisements to get prefix
4. Combines prefix with interface ID to create global address

The interface ID is typically derived from MAC address using EUI-64 (though privacy extensions randomize this now).

### Neighbor Discovery Protocol

NDP replaces ARP and adds functionality:
- **Neighbor Solicitation/Advertisement**: Address resolution (like ARP)
- **Router Solicitation/Advertisement**: Router discovery and prefix info
- **Redirect**: Route optimization

All NDP uses ICMPv6, which means you cannot blindly block ICMPv6 like many did with ICMPv4.

### Subnet Sizing

Standard practice is `/64` for all subnets regardless of size. This gives you the 64-bit interface ID space for SLAAC and makes life simpler. Point-to-point links sometimes use `/127` to avoid the subnet-router anycast address issue.

A `/48` is typically allocated to sites, giving you 65,536 `/64` subnets.

### Practical Differences from IPv4

- No broadcast - everything is multicast or unicast
- No NAT (in theory) - end-to-end connectivity restored
- No ARP - NDP handles neighbor resolution
- Fragmentation only at source - routers don't fragment, they send "Packet Too Big" ICMPv6
- MTU discovery is mandatory, not optional

---

## Getting IPv6 from Your ISP

### Check Current IPv6 Status

First, verify if you have IPv6 connectivity:

```bash
curl -6 ifconfig.co
# or
curl https://ipv6.icanhazip.com
```

If that works, you've got IPv6. Check what your router received:

```bash
# On your router/gateway
ip -6 addr show
ip -6 route show
```

### Understanding ISP Allocations

Most ISPs give you either:
- **Single /64**: Bare minimum, only one subnet (not ideal for home lab)
- **/56 or /48 delegation**: Multiple subnets via DHCPv6-PD (Prefix Delegation)

Common ISP allocations:
- AT&T Fiber: `/60`
- Comcast/Xfinity: `/60`
- Google Fiber: `/56`
- Verizon FiOS: Varies

If you're getting less than a `/60`, call your ISP and ask for prefix delegation.

---

## Network Design

### Subnet Planning

Assuming you got a `/56` delegation: `2001:db8:1234::/56` (using documentation prefix as example).

Break it down:
```
2001:db8:1234:00::/64 - Main home network (trusted devices)
2001:db8:1234:01::/64 - IoT/untrusted devices
2001:db8:1234:10::/64 - Home lab management network
2001:db8:1234:11::/64 - Home lab VM/container network
2001:db8:1234:12::/64 - Home lab storage network
2001:db8:1234:13::/64 - Home lab DMZ/public services
2001:db8:1234:20::/64 - Guest network
```

You've got 256 `/64`s to work with from a `/56`, so plan accordingly.

### Network Segmentation Benefits

- **Security isolation**: Contain breaches to specific networks
- **Traffic management**: QoS policies per network
- **Simplified firewall rules**: Clear security boundaries
- **Flexibility**: Easy to add new segments

---

## Router Configuration

### Prerequisites

Using a proper router platform:
- pfSense
- OPNsense
- VyOS
- MikroTik RouterOS
- OpenWrt

Not ISP-provided consumer routers (usually limited IPv6 support).

### pfSense/OPNsense Configuration

**Basic Setup:**

1. **WAN Interface**: Configure DHCPv6 client with prefix delegation
   - Services → DHCPv6 Client
   - Request prefix delegation size: 56 (or whatever your ISP provides)
   - Send prefix hint: Yes
   - Prefix delegation size: 56

2. **LAN Interfaces**: Assign delegated prefixes
   - Interfaces → LAN → IPv6 Configuration Type: Track Interface
   - IPv6 Interface: WAN
   - IPv6 Prefix ID: 0 (for first subnet), 1 (for second), etc.

3. **Router Advertisements**: Enable per interface
   - Services → Router Advertisements
   - Router mode: Managed, Assisted, or Unmanaged
   - Router priority: Normal

**Router Advertisement Modes:**

- **Managed**: DHCPv6 assigns addresses (more control, like DHCP for IPv4)
- **Assisted**: SLAAC for addresses, DHCPv6 for DNS/options
- **Unmanaged**: Pure SLAAC (simplest but least control)

For home lab, use **Assisted** on most networks - SLAAC for addressing but DHCPv6 for DNS servers and search domain.

### VyOS Configuration Example

```bash
# WAN interface - get prefix delegation from ISP
set interfaces ethernet eth0 address 'dhcpv6'
set interfaces ethernet eth0 dhcpv6-options prefix-delegation interface eth1 address '0'
set interfaces ethernet eth0 dhcpv6-options prefix-delegation interface eth1 sla-id '0'
set interfaces ethernet eth0 dhcpv6-options prefix-delegation length '56'

# Home network (VLAN 0)
set interfaces ethernet eth1 vif 0 address '2001:db8:1234:0::1/64'
set service router-advert interface eth1.0 prefix ::/64
set service router-advert interface eth1.0 name-server '2606:4700:4700::1111'

# IoT network (VLAN 1)
set interfaces ethernet eth1 vif 1 address '2001:db8:1234:1::1/64'
set service router-advert interface eth1.1 prefix ::/64
set service router-advert interface eth1.1 name-server '2606:4700:4700::1111'

# Lab management (VLAN 10)
set interfaces ethernet eth1 vif 10 address '2001:db8:1234:10::1/64'
set service router-advert interface eth1.10 prefix ::/64
set service router-advert interface eth1.10 name-server '2001:db8:1234:10::1'
```

---

## Addressing Strategy

### Option 1: SLAAC Only (Simplest)

**Pros:**
- Zero configuration
- Automatic address assignment
- Works out of the box

**Cons:**
- No fixed addresses for servers
- Privacy extensions mean addresses constantly change
- Harder to track/monitor specific devices

**Use case**: Guest networks, client devices

### Option 2: SLAAC + Static (Recommended for Home Lab)

**Configuration:**
- General devices use SLAAC
- Servers get manual static addresses in same subnet
- Pick a convention like `::1000` and up for static assignments

**Example server configuration:**

```bash
# On a Linux server
ip -6 addr add 2001:db8:1234:10::1000/64 dev eth0

# Make permanent (netplan example)
cat /etc/netplan/01-netcfg.yaml
```

```yaml
network:
  version: 2
  ethernets:
    eth0:
      dhcp6: no
      addresses:
        - 2001:db8:1234:10::1000/64
      routes:
        - to: ::/0
          via: fe80::1  # Router's link-local
      nameservers:
        addresses:
          - 2606:4700:4700::1111
          - 2606:4700:4700::1001
```

**Rocky Linux (network-scripts):**

```bash
# /etc/sysconfig/network-scripts/ifcfg-eth0
IPV6INIT=yes
IPV6ADDR=2001:db8:1234:10::1000/64
IPV6_DEFAULTGW=fe80::1%eth0
DNS1=2606:4700:4700::1111
DNS2=2606:4700:4700::1001
```

### Option 3: DHCPv6 Stateful (Most Control)

**Pros:**
- Like DHCP for IPv4
- Reservations for servers
- Centralized management

**Cons:**
- More complex to configure
- Some devices don't support DHCPv6 well (looking at you, Android)

**dnsmasq configuration:**

```bash
# /etc/dnsmasq.conf
dhcp-range=2001:db8:1234:10::100,2001:db8:1234:10::1ff,64,1h

# Static assignments
dhcp-host=id:00:01:00:01:2a:3b:4c:5d:6e:7f,2001:db8:1234:10::1000,nas
dhcp-host=id:00:01:00:01:1a:2b:3c:4d:5e:6f,2001:db8:1234:10::1001,prometheus
```

### Addressing Conventions

Recommended scheme within a `/64` subnet:

```
::1           - Router/gateway
::1-::ff      - Infrastructure (DNS, NTP, monitoring)
::100-::1ff   - DHCP/dynamic pool
::1000-::1fff - Static server assignments
::2000-::2fff - Container/VM static assignments
```

---

## Security and Firewall Rules

### Critical Difference from IPv4

You cannot just NAT everything and call it security. IPv6 is end-to-end, so firewall rules matter.

### Basic Firewall Rules (Pseudocode)

```
# Allow established/related
allow state established,related

# Allow ICMPv6 (essential for PMTUD, NDP)
allow icmpv6 type echo-request
allow icmpv6 type echo-reply
allow icmpv6 type destination-unreachable
allow icmpv6 type packet-too-big
allow icmpv6 type time-exceeded
allow icmpv6 type parameter-problem
allow icmpv6 type neighbor-solicitation
allow icmpv6 type neighbor-advertisement
allow icmpv6 type router-solicitation
allow icmpv6 type router-advertisement

# Allow DHCPv6 if using it
allow udp src-port 546 dst-port 547

# Allow outbound initiated by LAN
allow from lan-subnet to any

# Drop everything else inbound from internet
deny from any to lan-subnet
```

### Per-Network Security Policies

**Main Home Network:**
- Allow outbound everything
- Allow inbound to specific services (SSH to admin box, etc.)
- Rate-limit ICMP

**IoT Network:**
- Block IoT devices from initiating to home network
- Allow home network to initiate to IoT
- Block most outbound except DNS, NTP, HTTPS to specific vendors
- No inbound from internet

**Home Lab:**
- Depends on what you're running
- If hosting services, explicit allow rules for them
- Consider separate firewall rules for lab-to-internet vs lab-to-home
- Allow SSH from management network only

### pfSense Firewall Rules Example

**IoT Network Rules:**

```
# Allow IoT to respond to home network
Pass IPv6 from IoT-net to Home-net, state established

# Block IoT initiating to home
Reject IPv6 from IoT-net to Home-net

# Allow IoT to specific internet services
Pass IPv6 from IoT-net to any, port 53 (DNS)
Pass IPv6 from IoT-net to any, port 123 (NTP)
Pass IPv6 from IoT-net to any, port 443 (HTTPS)

# Log and block everything else
Block log IPv6 from IoT-net to any
```

### Privacy Extension Handling

Privacy extensions cause addresses to rotate regularly. On servers, disable them:

```bash
# Linux
sysctl -w net.ipv6.conf.eth0.use_tempaddr=0
echo "net.ipv6.conf.eth0.use_tempaddr=0" >> /etc/sysctl.conf

# Check current addresses
ip -6 addr show eth0
```

You'll see addresses marked:
- `scope global` - SLAAC permanent
- `scope global temporary` - Privacy extension

For client devices, leave privacy extensions enabled for privacy.

---

## Split-Horizon DNS

### Why Split-Horizon DNS?

You want internal devices to resolve lab services to internal IPv6 addresses, while keeping control over what's accessible externally.

### Option 1: dnsmasq (Simplest)

Great for small setups, runs on your router or a dedicated Pi.

```bash
# /etc/dnsmasq.conf

# Listen on lab network interface
interface=eth1
bind-interfaces

# Upstream DNS servers
server=2606:4700:4700::1111
server=2606:4700:4700::1001

# Local domain
domain=lab.home
local=/lab.home/

# Expand simple hostnames
expand-hosts

# Static host entries
host-record=router.lab.home,2001:db8:1234:10::1
host-record=nas.lab.home,2001:db8:1234:10::1000
host-record=k8s-master.lab.home,2001:db8:1234:11::100
host-record=prometheus.lab.home,2001:db8:1234:10::1001
host-record=grafana.lab.home,2001:db8:1234:10::1002

# Wildcard for services
address=/svc.lab.home/2001:db8:1234:11::200

# Enable DNS caching
cache-size=1000

# DHCP integration (if using DHCPv6)
dhcp-range=2001:db8:1234:10::100,2001:db8:1234:10::1ff,64,1h
```

**Configure Router Advertisements to use this DNS:**

```bash
# In pfSense: Services > Router Advertisements > DNS Servers
# Enter: 2001:db8:1234:10::1

# In VyOS:
set service router-advert interface eth1 name-server '2001:db8:1234:10::1'
```

### Option 2: Unbound (More Powerful)

Better for complex setups, supports DNSSEC, DNS-over-TLS.

```bash
# /etc/unbound/unbound.conf.d/lab.conf

server:
    # Listen on IPv6
    interface: 2001:db8:1234:10::1
    interface: ::1
    
    # Access control
    access-control: 2001:db8:1234::/56 allow
    access-control: ::1 allow
    access-control: ::/0 refuse
    
    # Local zone
    local-zone: "lab.home." static
    
    # Local data (AAAA records)
    local-data: "router.lab.home. IN AAAA 2001:db8:1234:10::1"
    local-data: "nas.lab.home. IN AAAA 2001:db8:1234:10::1000"
    local-data: "k8s-master.lab.home. IN AAAA 2001:db8:1234:11::100"
    local-data: "prometheus.lab.home. IN AAAA 2001:db8:1234:10::1001"
    local-data: "grafana.lab.home. IN AAAA 2001:db8:1234:10::1002"
    
    # Reverse DNS (optional)
    local-zone: "0.1.4.3.2.1.8.b.d.0.1.0.0.2.ip6.arpa." static
    local-data-ptr: "2001:db8:1234:10::1 router.lab.home"
    local-data-ptr: "2001:db8:1234:10::1000 nas.lab.home"
    
    # Performance tuning
    num-threads: 2
    msg-cache-size: 8m
    rrset-cache-size: 16m
    
forward-zone:
    name: "."
    forward-addr: 2606:4700:4700::1111@853  # Cloudflare DNS over TLS
    forward-addr: 2606:4700:4700::1001@853
    forward-tls-upstream: yes
```

### Option 3: Split DNS with Public Domain

If you own a real domain and want `*.lab.yourdomain.com` to work internally but resolve differently externally:

```bash
# Unbound with override
server:
    local-zone: "lab.yourdomain.com." transparent

stub-zone:
    name: "lab.yourdomain.com."
    stub-addr: 2001:db8:1234:10::1@5353  # Local authoritative server
```

Then run a separate BIND or NSD instance on port 5353 that authoritatively serves `lab.yourdomain.com` with internal addresses.

### Dynamic DNS Integration

If using stateful DHCPv6, integrate with DNS for automatic updates:

```bash
# dnsmasq handles this automatically with dhcp-range

# For ISC DHCP + BIND9 (more complex)
# In /etc/dhcp/dhcpd6.conf
subnet6 2001:db8:1234:10::/64 {
    range6 2001:db8:1234:10::100 2001:db8:1234:10::1ff;
    option dhcp6.name-servers 2001:db8:1234:10::1;
    ddns-domainname "lab.home.";
    ddns-rev-domainname "ip6.arpa.";
}
```

---

## Advanced Firewall Hardening

### Defense in Depth Strategy

Three layers:
1. **Router edge firewall** (WAN → LAN)
2. **Per-network segmentation rules**
3. **Host-based firewalls** (nftables/iptables)

### Layer 1: Edge Firewall (WAN Interface)

**VyOS Example - Complete WAN Lockdown:**

```bash
# Create firewall ruleset
set firewall ipv6-name WAN_IN default-action drop
set firewall ipv6-name WAN_IN enable-default-log

# Rule 10: Allow established/related
set firewall ipv6-name WAN_IN rule 10 action accept
set firewall ipv6-name WAN_IN rule 10 state established enable
set firewall ipv6-name WAN_IN rule 10 state related enable

# Rule 20: Allow essential ICMPv6
set firewall ipv6-name WAN_IN rule 20 action accept
set firewall ipv6-name WAN_IN rule 20 protocol ipv6-icmp
set firewall ipv6-name WAN_IN rule 20 icmpv6 type destination-unreachable

set firewall ipv6-name WAN_IN rule 21 action accept
set firewall ipv6-name WAN_IN rule 21 protocol ipv6-icmp
set firewall ipv6-name WAN_IN rule 21 icmpv6 type packet-too-big

set firewall ipv6-name WAN_IN rule 22 action accept
set firewall ipv6-name WAN_IN rule 22 protocol ipv6-icmp
set firewall ipv6-name WAN_IN rule 22 icmpv6 type time-exceeded

set firewall ipv6-name WAN_IN rule 23 action accept
set firewall ipv6-name WAN_IN rule 23 protocol ipv6-icmp
set firewall ipv6-name WAN_IN rule 23 icmpv6 type parameter-problem

set firewall ipv6-name WAN_IN rule 24 action accept
set firewall ipv6-name WAN_IN rule 24 protocol ipv6-icmp
set firewall ipv6-name WAN_IN rule 24 icmpv6 type echo-reply

# Rule 30: Rate-limit echo requests (prevent ping floods)
set firewall ipv6-name WAN_IN rule 30 action accept
set firewall ipv6-name WAN_IN rule 30 protocol ipv6-icmp
set firewall ipv6-name WAN_IN rule 30 icmpv6 type echo-request
set firewall ipv6-name WAN_IN rule 30 limit rate 10/minute

# Rule 100: Allow specific services (if hosting)
set firewall ipv6-name WAN_IN rule 100 action accept
set firewall ipv6-name WAN_IN rule 100 destination address 2001:db8:1234:10::1000
set firewall ipv6-name WAN_IN rule 100 destination port 443
set firewall ipv6-name WAN_IN rule 100 protocol tcp

# Rule 9999: Drop and log everything else
set firewall ipv6-name WAN_IN rule 9999 action drop
set firewall ipv6-name WAN_IN rule 9999 log enable

# Apply to WAN interface
set interfaces ethernet eth0 firewall in ipv6-name WAN_IN
```

**pfSense Equivalent:**

Firewall → Rules → WAN (IPv6)
- Add rules in order with logging enabled for drops

### Layer 2: Inter-Network Segmentation

**VyOS Example:**

```bash
# Home network to IoT: allowed
set firewall ipv6-name HOME_TO_IOT default-action accept

# IoT to Home: blocked except established
set firewall ipv6-name IOT_TO_HOME default-action drop
set firewall ipv6-name IOT_TO_HOME rule 10 action accept
set firewall ipv6-name IOT_TO_HOME rule 10 state established enable
set firewall ipv6-name IOT_TO_HOME rule 10 state related enable

# Lab to Home: selective access
set firewall ipv6-name LAB_TO_HOME default-action drop
set firewall ipv6-name LAB_TO_HOME rule 10 action accept
set firewall ipv6-name LAB_TO_HOME rule 10 state established enable
set firewall ipv6-name LAB_TO_HOME rule 10 state related enable

# Allow lab to access DNS/NTP on home network
set firewall ipv6-name LAB_TO_HOME rule 20 action accept
set firewall ipv6-name LAB_TO_HOME rule 20 destination port 53,123
set firewall ipv6-name LAB_TO_HOME rule 20 protocol udp

# Allow lab SSH to specific admin box
set firewall ipv6-name LAB_TO_HOME rule 30 action accept
set firewall ipv6-name LAB_TO_HOME rule 30 destination address 2001:db8:1234:0::10
set firewall ipv6-name LAB_TO_HOME rule 30 destination port 22
set firewall ipv6-name LAB_TO_HOME rule 30 protocol tcp

# Apply to interfaces
set interfaces ethernet eth1 vif 1 firewall out ipv6-name IOT_TO_HOME
set interfaces ethernet eth1 vif 10 firewall out ipv6-name LAB_TO_HOME
```

### Layer 3: Host-Based Firewall (nftables)

**Production-quality nftables configuration:**

```bash
# /etc/nftables.conf

flush ruleset

table ip6 filter {
    chain input {
        type filter hook input priority 0; policy drop;
        
        # Allow loopback
        iif lo accept comment "Allow loopback"
        
        # Allow established/related
        ct state established,related accept comment "Allow established"
        
        # Drop invalid
        ct state invalid drop comment "Drop invalid packets"
        
        # Allow ICMPv6 neighbor discovery
        icmpv6 type { nd-neighbor-solicit, nd-neighbor-advert, nd-router-solicit, nd-router-advert } accept comment "Allow NDP"
        
        # Allow ICMPv6 errors
        icmpv6 type { destination-unreachable, packet-too-big, time-exceeded, parameter-problem } accept comment "Allow ICMP errors"
        
        # Rate-limited ping
        icmpv6 type echo-request limit rate 10/second accept comment "Rate-limited ping"
        
        # SSH from management network only
        ip6 saddr 2001:db8:1234:10::/64 tcp dport 22 ct state new limit rate 5/minute accept comment "SSH from mgmt"
        
        # HTTPS from anywhere
        tcp dport 443 ct state new accept comment "HTTPS"
        
        # Prometheus metrics from monitoring host only
        ip6 saddr 2001:db8:1234:10::1001 tcp dport 9100 ct state new accept comment "Node exporter"
        
        # Log dropped packets (rate limited)
        limit rate 5/minute counter log prefix "nft-drop: " comment "Log drops"
    }
    
    chain forward {
        type filter hook forward priority 0; policy drop;
    }
    
    chain output {
        type filter hook output priority 0; policy accept;
    }
}

# Optional: Rate limiting table
table ip6 ratelimit {
    chain input {
        type filter hook input priority -150; policy accept;
        
        # Limit new connections per source
        ip6 saddr != 2001:db8:1234::/56 ct state new limit rate over 20/minute drop
    }
}
```

**Enable and test:**

```bash
# Test configuration
nft -f /etc/nftables.conf

# Enable on boot
systemctl enable nftables
systemctl start nftables

# View current rules
nft list ruleset

# Monitor drops in real-time
journalctl -kf | grep nft-drop
```

### Application-Level Controls

**SSH Hardening (`/etc/ssh/sshd_config`):**

```bash
# Only listen on specific IPv6 address
ListenAddress 2001:db8:1234:10::1000

# Or prefer IPv6 but allow both
AddressFamily any
ListenAddress ::
ListenAddress 0.0.0.0

# Restrict to specific users/groups
AllowUsers muck
AllowGroups admins

# Key-based auth only
PasswordAuthentication no
PubkeyAuthentication yes
PermitRootLogin no

# Use modern crypto only
KexAlgorithms curve25519-sha256,curve25519-sha256@libssh.org
Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com
```

**Nginx Web Server:**

```nginx
server {
    # Listen on both IPv6 and IPv4
    listen [2001:db8:1234:10::1000]:443 ssl http2;
    listen 443 ssl http2;
    
    server_name nas.lab.home;
    
    # Restrict access by source network
    allow 2001:db8:1234:0::/64;   # Home network
    allow 2001:db8:1234:10::/64;  # Lab network
    deny all;
    
    # SSL configuration
    ssl_certificate /etc/ssl/certs/nas.lab.home.crt;
    ssl_certificate_key /etc/ssl/private/nas.lab.home.key;
    ssl_protocols TLSv1.3;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

**PostgreSQL (`postgresql.conf`):**

```bash
# Listen on specific IPv6 address
listen_addresses = '2001:db8:1234:10::1000,localhost'

# Connection limits
max_connections = 100

# In pg_hba.conf
# TYPE  DATABASE        USER            ADDRESS                 METHOD
hostssl all             all             2001:db8:1234:10::/64   scram-sha-256
hostssl all             all             2001:db8:1234:11::/64   scram-sha-256
```

---

## Automating Prefix Delegation

### The Problem

ISPs can change your delegated prefix at any time:
- DHCP lease renewal
- Service interruption
- Equipment changes
- Random ISP whims

Your prefix changes from `2001:db8:1234::/56` to `2001:db8:5678::/56`, breaking:
- DNS records
- Firewall rules
- Service configurations
- Monitoring systems

### Solution 1: ULA + GUA Dual Addressing

Use Unique Local Addresses (ULA, `fd00::/8`) that never change, plus Global Unicast Addresses (GUA).

**Generate a ULA prefix:**

```bash
# Generate random ULA prefix (do this once, keep forever)
# Format: fd + 40 random bits + subnet + interface ID
# Example: fd12:3456:789a::/48

# Use https://www.unique-local-ipv6.com/ or:
python3 -c "import random; print('fd%02x:%04x:%04x::/48' % (random.randint(0,255), random.randint(0,65535), random.randint(0,65535)))"
```

**Assign both ULA and GUA to servers:**

```bash
# Server configuration with dual addressing
# /etc/netplan/01-netcfg.yaml
network:
  version: 2
  ethernets:
    eth0:
      dhcp6: no
      addresses:
        - fd12:3456:789a:10::1000/64  # ULA - never changes
        - 2001:db8:1234:10::1000/64   # GUA - might change
      routes:
        - to: ::/0
          via: fe80::1
      nameservers:
        addresses:
          - fd12:3456:789a:10::1  # Use ULA for DNS
```

**DNS uses ULA for internal resolution:**

```bash
# dnsmasq configuration
host-record=nas.lab.home,fd12:3456:789a:10::1000
host-record=prometheus.lab.home,fd12:3456:789a:10::1001
host-record=grafana.lab.home,fd12:3456:789a:10::1002
```

**Benefits:**
- Internal services always work regardless of GUA changes
- External access breaks but internal doesn't
- No automation needed for prefix changes

**Considerations:**
- ULA is not routable on internet (by design)
- Need GUA for external services
- Slightly more complex configuration

### Solution 2: Automated Prefix Monitoring and Update

**Python script to monitor and respond to prefix changes:**

```python
#!/usr/bin/env python3
# /usr/local/bin/ipv6-prefix-monitor.py

import subprocess
import json
import time
import sys
import logging
from pathlib import Path

# Configuration
CURRENT_PREFIX_FILE = "/var/lib/ipv6-prefix.json"
LOG_FILE = "/var/log/ipv6-prefix-monitor.log"
NETWORKS = {
    "home": {"subnet_id": "00", "vlan": 0},
    "iot": {"subnet_id": "01", "vlan": 1},
    "lab-mgmt": {"subnet_id": "10", "vlan": 10},
    "lab-vms": {"subnet_id": "11", "vlan": 11},
    "lab-storage": {"subnet_id": "12", "vlan": 12},
}

# Setup logging
logging.basicConfig(
    filename=LOG_FILE,
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def get_delegated_prefix():
    """Get current delegated prefix from DHCPv6-PD"""
    try:
        # Method varies by platform - this is for dhclient
        with open('/var/lib/dhcp/dhclient6.leases', 'r') as f:
            content = f.read()
            
        # Parse for iaprefix
        for line in content.split('\n'):
            if 'iaprefix' in line:
                # Extract prefix from line like:
                # iaprefix 2001:db8:1234::/56 { ... }
                parts = line.strip().split()
                if len(parts) >= 2:
                    prefix = parts[1]
                    logging.info(f"Found delegated prefix: {prefix}")
                    return prefix
                    
    except FileNotFoundError:
        logging.error("DHCPv6 lease file not found")
    except Exception as e:
        logging.error(f"Error reading prefix: {e}")
    
    return None

def load_saved_prefix():
    """Load previously known prefix"""
    try:
        with open(CURRENT_PREFIX_FILE, 'r') as f:
            data = json.load(f)
            return data.get('prefix')
    except FileNotFoundError:
        return None
    except json.JSONDecodeError:
        logging.error("Corrupted prefix file")
        return None

def save_prefix(prefix):
    """Save current prefix to file"""
    data = {
        'prefix': prefix,
        'updated': time.time(),
        'timestamp': time.strftime('%Y-%m-%d %H:%M:%S')
    }
    with open(CURRENT_PREFIX_FILE, 'w') as f:
        json.dump(data, f, indent=2)

def extract_base_prefix(prefix_with_len):
    """Extract base prefix without length"""
    # 2001:db8:1234::/56 -> 2001:db8:1234
    prefix = prefix_with_len.split('/')[0]
    # Remove trailing :: if present
    prefix = prefix.rstrip(':')
    return prefix

def update_network_configs(old_prefix, new_prefix):
    """Update all network configurations with new prefix"""
    logging.info(f"Prefix changed: {old_prefix} -> {new_prefix}")
    
    old_base = extract_base_prefix(old_prefix)
    new_base = extract_base_prefix(new_prefix)
    
    # Update router configurations
    update_router_configs(old_base, new_base)
    
    # Update DNS records
    update_dns_records(old_base, new_base)
    
    # Update firewall rules
    update_firewall_rules(old_base, new_base)
    
    # Send notification
    send_notification(f"IPv6 prefix changed from {old_prefix} to {new_prefix}")

def update_router_configs(old_base, new_base):
    """Update router interface addresses and RAs"""
    logging.info("Updating router configurations")
    
    for name, config in NETWORKS.items():
        subnet_id = config['subnet_id']
        vlan = config['vlan']
        
        old_addr = f"{old_base}:{subnet_id}::1/64"
        new_addr = f"{new_base}:{subnet_id}::1/64"
        
        # VyOS commands (adjust for your router platform)
        commands = [
            f"delete interfaces ethernet eth1 vif {vlan} address '{old_addr}'",
            f"set interfaces ethernet eth1 vif {vlan} address '{new_addr}'",
            f"delete service router-advert interface eth1.{vlan} prefix {old_base}:{subnet_id}::/64",
            f"set service router-advert interface eth1.{vlan} prefix {new_base}:{subnet_id}::/64",
        ]
        
        for cmd in commands:
            try:
                result = subprocess.run(
                    ['vtysh', '-c', 'configure terminal', '-c', cmd],
                    capture_output=True,
                    text=True,
                    timeout=10
                )
                if result.returncode != 0:
                    logging.error(f"Command failed: {cmd} - {result.stderr}")
            except Exception as e:
                logging.error(f"Error running command {cmd}: {e}")
    
    # Commit changes
    try:
        subprocess.run(['vtysh', '-c', 'write memory'], timeout=10)
    except Exception as e:
        logging.error(f"Error committing changes: {e}")

def update_dns_records(old_base, new_base):
    """Update DNS records with new prefix"""
    logging.info("Updating DNS records")
    
    dns_config_file = '/etc/dnsmasq.d/lab-hosts.conf'
    
    try:
        with open(dns_config_file, 'r') as f:
            content = f.read()
        
        # Replace old prefix with new
        updated_content = content.replace(old_base, new_base)
        
        with open(dns_config_file, 'w') as f:
            f.write(updated_content)
        
        # Restart dnsmasq
        subprocess.run(['systemctl', 'restart', 'dnsmasq'], timeout=10)
        logging.info("DNS records updated and dnsmasq restarted")
        
    except Exception as e:
        logging.error(f"Error updating DNS: {e}")

def update_firewall_rules(old_base, new_base):
    """Update firewall rules with new prefix"""
    logging.info("Updating firewall rules")
    
    # Example: Update specific firewall rules
    # This is highly dependent on your firewall platform
    
    try:
        # VyOS example - update a specific rule
        subprocess.run([
            'vtysh', '-c', 'configure terminal',
            '-c', f'set firewall ipv6-name WAN_IN rule 100 destination address {new_base}:10::1000'
        ], timeout=10)
        
        subprocess.run(['vtysh', '-c', 'write memory'], timeout=10)
        
    except Exception as e:
        logging.error(f"Error updating firewall: {e}")

def send_notification(message):
    """Send notification about prefix change"""
    logging.info(f"Notification: {message}")
    
    # Use ntfy.sh, pushover, email, etc.
    try:
        subprocess.run([
            'curl', '-H', 'Priority: high',
            '-d', message,
            'https://ntfy.sh/your-topic-here'
        ], timeout=10)
    except Exception as e:
        logging.error(f"Error sending notification: {e}")

def main():
    """Main monitoring loop"""
    logging.info("Starting IPv6 prefix monitor")
    
    while True:
        try:
            current = get_delegated_prefix()
            
            if current:
                saved = load_saved_prefix()
                
                if saved and saved != current:
                    logging.warning(f"Prefix change detected!")
                    update_network_configs(saved, current)
                
                save_prefix(current)
            else:
                logging.warning("Could not determine current prefix")
            
        except Exception as e:
            logging.error(f"Error in main loop: {e}")
        
        # Check every 5 minutes
        time.sleep(300)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        logging.info("Monitor stopped by user")
        sys.exit(0)
```

**Systemd service file:**

```ini
# /etc/systemd/system/ipv6-prefix-monitor.service
[Unit]
Description=IPv6 Prefix Delegation Monitor
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/bin/python3 /usr/local/bin/ipv6-prefix-monitor.py
Restart=always
RestartSec=10
User=root
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Enable and start:**

```bash
chmod +x /usr/local/bin/ipv6-prefix-monitor.py
systemctl daemon-reload
systemctl enable ipv6-prefix-monitor
systemctl start ipv6-prefix-monitor

# Check status
systemctl status ipv6-prefix-monitor
journalctl -u ipv6-prefix-monitor -f
```

### Solution 3: Dynamic DNS for External Access

Update external DNS when your prefix changes.

```bash
#!/bin/bash
# /usr/local/bin/ddns-update-ipv6.sh

DOMAIN="home.yourdomain.com"
CLOUDFLARE_TOKEN="your-api-token"
ZONE_ID="your-zone-id"
RECORD_ID="your-record-id"

# Get current global IPv6 address
CURRENT_IP=$(curl -s -6 https://ifconfig.co)

if [ -z "$CURRENT_IP" ]; then
    echo "Failed to get IPv6 address"
    exit 1
fi

# Update Cloudflare DNS
RESPONSE=$(curl -s -X PUT "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records/$RECORD_ID" \
     -H "Authorization: Bearer $CLOUDFLARE_TOKEN" \
     -H "Content-Type: application/json" \
     --data "{\"type\":\"AAAA\",\"name\":\"$DOMAIN\",\"content\":\"$CURRENT_IP\",\"ttl\":120,\"proxied\":false}")

if echo "$RESPONSE" | grep -q '"success":true'; then
    echo "Successfully updated $DOMAIN to $CURRENT_IP"
else
    echo "Failed to update DNS: $RESPONSE"
fi
```

**Cron job:**

```bash
# /etc/cron.d/ddns-ipv6
*/5 * * * * root /usr/local/bin/ddns-update-ipv6.sh >> /var/log/ddns-ipv6.log 2>&1
```

### Solution 4: Configuration Management (Ansible)

Use Ansible to template and deploy configurations.

**Inventory variables:**

```yaml
# inventory/group_vars/all.yml
ipv6_prefix: "2001:db8:1234"
ipv6_prefix_len: 56
ipv6_full_prefix: "{{ ipv6_prefix }}::/{{ ipv6_prefix_len }}"

networks:
  home:
    vlan: 0
    subnet_id: "00"
    description: "Main home network"
    router_ip: "{{ ipv6_prefix }}:00::1"
  iot:
    vlan: 1
    subnet_id: "01"
    description: "IoT devices"
    router_ip: "{{ ipv6_prefix }}:01::1"
  lab_mgmt:
    vlan: 10
    subnet_id: "10"
    description: "Lab management"
    router_ip: "{{ ipv6_prefix }}:10::1"
```

**Playbook:**

```yaml
# playbooks/update-ipv6-config.yml
---
- name: Update IPv6 configurations
  hosts: router
  become: yes
  tasks:
    - name: Configure interface IPv6 addresses
      vyos_config:
        lines:
          - set interfaces ethernet eth1 vif {{ item.value.vlan }} address '{{ item.value.router_ip }}/64'
      loop: "{{ networks | dict2items }}"
      
    - name: Configure router advertisements
      vyos_config:
        lines:
          - set service router-advert interface eth1.{{ item.value.vlan }} prefix {{ ipv6_prefix }}:{{ item.value.subnet_id }}::/64
      loop: "{{ networks | dict2items }}"
      
    - name: Save configuration
      vyos_config:
        save: yes

- name: Update DNS server
  hosts: dns_server
  become: yes
  tasks:
    - name: Template dnsmasq configuration
      template:
        src: templates/dnsmasq-hosts.j2
        dest: /etc/dnsmasq.d/lab-hosts.conf
      notify: restart dnsmasq
      
  handlers:
    - name: restart dnsmasq
      service:
        name: dnsmasq
        state: restarted
```

**Template:**

```jinja2
# templates/dnsmasq-hosts.j2
# Auto-generated DNS configuration
# Generated: {{ ansible_date_time.iso8601 }}

{% for network_name, network in networks.items() %}
# {{ network.description }}
{% endfor %}

host-record=router.lab.home,{{ ipv6_prefix }}:10::1
host-record=nas.lab.home,{{ ipv6_prefix }}:10::1000
host-record=prometheus.lab.home,{{ ipv6_prefix }}:10::1001
host-record=grafana.lab.home,{{ ipv6_prefix }}:10::1002
```

**Usage:**

```bash
# When prefix changes, update inventory variable
vim inventory/group_vars/all.yml
# Change: ipv6_prefix: "2001:db8:5678"

# Run playbook
ansible-playbook playbooks/update-ipv6-config.yml
```

---

## Monitoring and Visibility

### Basic Monitoring Script

```bash
#!/bin/bash
# /usr/local/bin/ipv6-status.sh

echo "======================================"
echo "IPv6 Network Status"
echo "======================================"
echo
echo "Current Date: $(date)"
echo

echo "--- Delegated Prefix ---"
ip -6 route show | grep -v fe80 | head -5
echo

echo "--- Interface Addresses ---"
ip -6 addr show | grep "inet6" | grep -v "fe80::" | grep -v "::1/128"
echo

echo "--- Active IPv6 Connections ---"
CONN_COUNT=$(ss -6 -tn 2>/dev/null | wc -l)
echo "Total connections: $((CONN_COUNT - 1))"
echo

echo "--- NDP Table ---"
NDP_COUNT=$(ip -6 neigh show | grep -v FAILED | wc -l)
echo "Neighbor entries: $NDP_COUNT"
ip -6 neigh show | grep -v FAILED | head -10
echo

echo "--- Top IPv6 Talkers ---"
ss -6 -tn 2>/dev/null | awk '{print $5}' | grep -v "^Local" | cut -d: -f1-8 | sort | uniq -c | sort -rn | head -10
echo

echo "--- DNS Resolution Test ---"
dig +short AAAA google.com @2606:4700:4700::1111 | head -1
echo

echo "--- Connectivity Test ---"
if ping6 -c 1 -W 2 2001:4860:4860::8888 >/dev/null 2>&1; then
    echo "✓ IPv6 internet connectivity: OK"
else
    echo "✗ IPv6 internet connectivity: FAILED"
fi
```

### Prometheus Node Exporter

Standard node_exporter already exports IPv6 metrics. Configure it to listen on IPv6:

```bash
# /etc/systemd/system/node_exporter.service
[Unit]
Description=Prometheus Node Exporter
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/node_exporter \
    --web.listen-address="[::]:9100" \
    --collector.netdev.device-include="^(eth|ens|enp|wlan)" \
    --collector.filesystem.fs-types-exclude="^(autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$"

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Custom IPv6 Metrics Exporter

```python
#!/usr/bin/env python3
# /usr/local/bin/ipv6-metrics-exporter.py

from prometheus_client import start_http_server, Gauge, Info
import subprocess
import time
import re

# Define metrics
ndp_entries = Gauge('ipv6_ndp_entries', 'Number of NDP table entries')
active_connections = Gauge('ipv6_active_connections', 'Active IPv6 TCP connections')
prefix_info = Info('ipv6_delegated_prefix', 'Current delegated IPv6 prefix')

def get_ndp_count():
    """Count NDP table entries"""
    try:
        result = subprocess.run(
            ['ip', '-6', 'neigh', 'show'],
            capture_output=True,
            text=True,
            timeout=5
        )
        # Count non-FAILED entries
        entries = [line for line in result.stdout.split('\n') if line and 'FAILED' not in line]
        return len(entries)
    except Exception:
        return 0

def get_connection_count():
    """Count active IPv6 TCP connections"""
    try:
        result = subprocess.run(
            ['ss', '-6', '-tn'],
            capture_output=True,
            text=True,
            timeout=5
        )
        # Count lines minus header
        return len(result.stdout.strip().split('\n')) - 1
    except Exception:
        return 0

def get_delegated_prefix():
    """Get current delegated prefix"""
    try:
        result = subprocess.run(
            ['ip', '-6', 'route', 'show'],
            capture_output=True,
            text=True,
            timeout=5
        )
        # Look for delegated prefix (not link-local)
        for line in result.stdout.split('\n'):
            if 'proto' in line and 'fe80::' not in line and '::1' not in line:
                # Extract prefix
                match = re.search(r'([0-9a-f:]+/\d+)', line)
                if match:
                    return match.group(1)
    except Exception:
        pass
    return "unknown"

def collect_metrics():
    """Collect all metrics"""
    ndp_entries.set(get_ndp_count())
    active_connections.set(get_connection_count())
    prefix = get_delegated_prefix()
    prefix_info.info({'prefix': prefix})

def main():
    # Start HTTP server on IPv6
    start_http_server(9101, addr='::')
    print("IPv6 metrics exporter started on [::]:9101")
    
    while True:
        try:
            collect_metrics()
        except Exception as e:
            print(f"Error collecting metrics: {e}")
        
        time.sleep(15)

if __name__ == '__main__':
    main()
```

**Systemd service:**

```ini
# /etc/systemd/system/ipv6-metrics-exporter.service
[Unit]
Description=IPv6 Metrics Exporter for Prometheus
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/bin/python3 /usr/local/bin/ipv6-metrics-exporter.py
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Prometheus Configuration

```yaml
# /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'node'
    static_configs:
      - targets:
        # Use IPv6 addresses in brackets
        - '[2001:db8:1234:10::1000]:9100'
        - '[2001:db8:1234:10::1001]:9100'
        labels:
          environment: 'homelab'
  
  - job_name: 'ipv6-custom'
    static_configs:
      - targets:
        - '[2001:db8:1234:10::1]:9101'
        - '[2001:db8:1234:10::1000]:9101'
```

### Alerting Rules

```yaml
# /etc/prometheus/rules/ipv6-alerts.yml
groups:
  - name: ipv6_alerts
    interval: 30s
    rules:
      - alert: IPv6PrefixChanged
        expr: changes(ipv6_delegated_prefix[10m]) > 0
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "IPv6 prefix has changed"
          description: "The delegated IPv6 prefix has changed in the last 10 minutes"
      
      - alert: HighNDPEntries
        expr: ipv6_ndp_entries > 200
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Unusually high NDP table entries"
          description: "NDP table has {{ $value }} entries, which is unusually high"
      
      - alert: IPv6ConnectivityLoss
        expr: up{job="node"} == 0
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "IPv6 connectivity lost to {{ $labels.instance }}"
          description: "Cannot scrape metrics from {{ $labels.instance }}"
      
      - alert: NoActiveIPv6Connections
        expr: ipv6_active_connections == 0
        for: 10m
        labels:
          severity: info
        annotations:
          summary: "No active IPv6 connections"
          description: "No IPv6 TCP connections detected for 10 minutes"
```

### Grafana Dashboard

Sample Grafana dashboard JSON snippet:

```json
{
  "dashboard": {
    "title": "IPv6 Home Lab",
    "panels": [
      {
        "title": "NDP Table Entries",
        "type": "graph",
        "targets": [
          {
            "expr": "ipv6_ndp_entries",
            "legendFormat": "NDP Entries"
          }
        ]
      },
      {
        "title": "Active IPv6 Connections",
        "type": "graph",
        "targets": [
          {
            "expr": "ipv6_active_connections",
            "legendFormat": "Connections"
          }
        ]
      },
      {
        "title": "Current Prefix",
        "type": "stat",
        "targets": [
          {
            "expr": "ipv6_delegated_prefix",
            "legendFormat": "Prefix"
          }
        ]
      }
    ]
  }
}
```

### Network Flow Monitoring

Use `nfdump` for NetFlow/IPFIX or `tcpdump` for packet capture:

```bash
# Capture IPv6 traffic for analysis
tcpdump -i eth0 -w /tmp/ipv6-capture.pcap 'ip6' -c 10000

# Analyze with tshark
tshark -r /tmp/ipv6-capture.pcap -q -z conv,ipv6

# Monitor real-time IPv6 traffic
iftop -i eth0 -f 'ip6'
```

---

## Common Gotchas

### Docker IPv6 Support

Docker defaults to IPv4 only. Enable IPv6:

```json
{
  "ipv6": true,
  "fixed-cidr-v6": "2001:db8:1234:ff::/64",
  "experimental": true,
  "ip6tables": true
}
```

Restart Docker:

```bash
systemctl restart docker
```

Test:

```bash
docker run --rm curlimages/curl:latest curl -6 https://ifconfig.co
```

### KVM/libvirt Bridge Configuration

Bridges need IPv6 enabled:

```bash
# Check if IPv6 is disabled on bridge
cat /proc/sys/net/ipv6/conf/virbr0/disable_ipv6  # Should be 0

# If disabled, enable it
echo 0 > /proc/sys/net/ipv6/conf/virbr0/disable_ipv6

# Make permanent
cat >> /etc/sysctl.conf << EOF
net.ipv6.conf.virbr0.disable_ipv6 = 0
EOF
```

### SSH Connection Issues

SSH might prefer IPv4. Force IPv6:

```bash
# From command line
ssh -6 user@hostname

# In ~/.ssh/config
Host server1
    HostName 2001:db8:1234:10::1000
    AddressFamily inet6
    
Host *.lab.home
    AddressFamily inet6
```

### Application Binding

Some applications need explicit IPv6 configuration:

**MySQL/MariaDB:**

```ini
# /etc/mysql/my.cnf
[mysqld]
bind-address = ::
```

**Redis:**

```conf
# /etc/redis/redis.conf
bind :: 0.0.0.0
```

### MTU and PMTUD Issues

IPv6 requires PMTUD to work. If ICMPv6 "Packet Too Big" is filtered, things break silently.

Minimum IPv6 MTU is 1280 bytes.

**Test PMTUD:**

```bash
# Ping with Don't Fragment
ping6 -M do -s 1452 2001:4860:4860::8888

# If this fails but normal ping works, PMTUD is broken
ping6 2001:4860:4860::8888
```

**Fix:**

Ensure ICMPv6 type 2 (Packet Too Big) is allowed in all firewalls.

### Extension Header Filtering

Some firewalls/ACLs drop packets with IPv6 extension headers.

**Test:**

```bash
# This might fail if extension headers are filtered
ping6 -c 5 -s 2000 2001:4860:4860::8888
```

**Solution:**

Configure firewalls to allow extension headers or at least fragment headers.

---

## Troubleshooting Techniques

### Basic Connectivity Tests

```bash
# Ping IPv6
ping6 google.com
ping6 2001:4860:4860::8888

# Ping link-local (requires interface specification)
ping6 -I eth0 fe80::1

# Ping all nodes on link
ping6 -I eth0 ff02::1
```

### Trace Route

```bash
# Trace to destination
traceroute6 2001:4860:4860::8888

# Alternative (doesn't need root)
tracepath6 2001:4860:4860::8888
```

### Check Neighbor Discovery

```bash
# View NDP table
ip -6 neigh show

# Specific interface
ip -6 neigh show dev eth0

# Force refresh
ip -6 neigh flush dev eth0
```

### Routing Table

```bash
# View IPv6 routes
ip -6 route show

# Lookup specific destination
ip -6 route get 2001:4860:4860::8888

# Show only default route
ip -6 route show default
```

### Socket/Connection Inspection

```bash
# All IPv6 TCP sockets
ss -6 -tan

# All IPv6 UDP sockets
ss -6 -uan

# Listening IPv6 TCP sockets
ss -6 -tln

# What's listening on port 443?
ss -6 -tln sport = :443
```

### DNS Testing

```bash
# Query AAAA record
dig AAAA google.com

# Query specific DNS server
dig AAAA google.com @2606:4700:4700::1111

# Reverse DNS lookup
dig -x 2001:4860:4860::8888

# Check what DNS server is being used
cat /etc/resolv.conf
```

### Router Advertisement Monitoring

```bash
# Install rdisc6 (from ndisc6 package)
# Debian/Ubuntu: apt install ndisc6
# RHEL/Rocky: yum install ndisc6

# Listen for router advertisements
rdisc6 eth0

# Expected output shows prefix, MTU, DNS servers
```

### Packet Capture and Analysis

**tcpdump:**

```bash
# Capture all IPv6 traffic
tcpdump -i eth0 -n 'ip6'

# Capture ICMPv6 only
tcpdump -i eth0 -n 'icmp6'

# Capture specific ICMPv6 types
tcpdump -i eth0 -n 'icmp6 and ip6[40] == 135'  # Neighbor Solicitation
tcpdump -i eth0 -n 'icmp6 and ip6[40] == 136'  # Neighbor Advertisement

# Capture BGP over IPv6
tcpdump -i eth0 -n 'ip6 and tcp port 179'

# Save to file for Wireshark
tcpdump -i eth0 -w /tmp/capture.pcap 'ip6'
```

**Wireshark Display Filters:**

```
ipv6.addr == 2001:db8::1
icmpv6.type == 135  # Neighbor Solicitation
icmpv6.type == 136  # Neighbor Advertisement
ipv6.dst == ff02::1  # Multicast to all nodes
tcp.port == 443 && ipv6
```

### Testing Firewall Rules

```bash
# Test from external host
# On external host with IPv6:
nmap -6 2001:db8:1234:10::1000

# Test specific port
nc -6 -zv 2001:db8:1234:10::1000 443

# Test from internal network
# Should work if firewall allows
curl -6 -v https://[2001:db8:1234:10::1000]
```

### Debugging NDP Issues

```bash
# Watch for NDP messages
tcpdump -i eth0 -vv 'icmp6 and (ip6[40] == 135 or ip6[40] == 136)'

# Check if interface is doing DAD (Duplicate Address Detection)
ip -6 addr show eth0 | grep tentative

# Force DAD to complete
ip -6 addr add 2001:db8:1234:10::1000/64 dev eth0
# Wait a few seconds for DAD
ip -6 addr show eth0
```

### Performance Testing

```bash
# iperf3 over IPv6
# On server:
iperf3 -s -B ::

# On client:
iperf3 -c 2001:db8:1234:10::1000

# HTTP download speed test
curl -6 -o /dev/null -w '%{speed_download}\n' https://speed.cloudflare.com/__down?bytes=100000000
```

### Checking SLAAC

```bash
# See if SLAAC addresses are being assigned
ip -6 addr show | grep "scope global"

# Check router advertisements are being sent
tcpdump -i eth1 -vv 'icmp6 and ip6[40] == 134'  # RA type 134

# On client, request new RA
rdisc6 eth0
```

### Log Analysis

```bash
# Check system logs for IPv6 issues
journalctl -k | grep -i ipv6

# Check for NDP issues
journalctl -k | grep -i "neighbour"

# Check for address conflicts
journalctl -k | grep -i "duplicate"

# Firewall logs (if logging enabled)
journalctl -k | grep -i "nft-drop"
```

---

## BGP with IPv6

### Separate Address Families

IPv4 and IPv6 are distinct address families in BGP. You establish one TCP session but negotiate capabilities for both.

**Cisco IOS-XR Example:**

```
router bgp 65001
 neighbor 2001:db8::2 remote-as 65002
 !
 address-family ipv6 unicast
  neighbor 2001:db8::2 activate
  network 2001:db8:1::/48
 exit-address-family
```

**Junos Example:**

```
protocols {
    bgp {
        group external {
            type external;
            neighbor 2001:db8::2 {
                peer-as 65002;
                family inet6 {
                    unicast;
                }
            }
        }
    }
}
```

### Session Establishment

BGP sessions can run over IPv4 or IPv6 transport:

- **IPv6 sessions for IPv6 routes** (most common)
- **IPv4 sessions for both** (easier during migration)
- **Dual sessions** (separate for each AF)

**Link-local BGP sessions:**

```
neighbor fe80::1%eth0 remote-as 65002
```

Zone ID required with link-local addresses.

### Path Selection

Same BGP path selection algorithm as IPv4:
1. Weight (Cisco-specific)
2. Local preference
3. Locally originated
4. AS path length
5. Origin type
6. MED
7. eBGP over iBGP
8. Lowest IGP metric
9. Oldest route
10. Lowest router ID

### Automation Considerations

**TextFSM template differences:**

```
# show bgp ipv6 unicast summary output differs from IPv4
# Neighbor addresses are longer, breaking fixed-width parsing
```

Use structured output (JSON/XML via NETCONF) instead:

```python
from ncclient import manager

with manager.connect(host='router', port=830, username='admin',
                     password='password', hostkey_verify=False) as m:
    result = m.get(filter=('subtree', '''
        <bgp xmlns="http://cisco.com/ns/yang/Cisco-IOS-XR-ipv6-bgp-oper">
          <instances>
            <instance>
              <instance-name>default</instance-name>
              <instance-active>
                <default-vrf>
                  <neighbors/>
                </default-vrf>
              </instance-active>
            </instance>
          </instances>
        </bgp>
    '''))
```

### Prefix Lists and Filters

Different syntax for IPv6:

```
# Cisco
ipv6 prefix-list ALLOW-CUSTOMER seq 10 permit 2001:db8::/32 le 48

# Junos
policy-options {
    prefix-list ipv6-customer {
        2001:db8::/32;
    }
}
```

### Python Automation Example

```python
import ipaddress
from netmiko import ConnectHandler

def validate_ipv6_bgp_peer(device, peer_addr):
    """Validate IPv6 BGP peer is established and receiving routes"""
    try:
        peer = ipaddress.IPv6Address(peer_addr)
        
        connection = ConnectHandler(**device)
        
        # Get BGP summary
        output = connection.send_command(
            f"show bgp ipv6 unicast summary | include {peer_addr}"
        )
        
        if "Established" not in output:
            return {"status": "down", "peer": str(peer)}
            
        # Parse for received routes
        parts = output.split()
        received_routes = parts[-1] if parts else "0"
        
        connection.disconnect()
        
        return {
            "status": "up",
            "peer": str(peer),
            "routes_received": received_routes
        }
        
    except Exception as e:
        return {"status": "error", "peer": peer_addr, "error": str(e)}
```

---

## Summary

This guide covers:

1. **Fundamentals**: Address structure, NDP, SLAAC
2. **ISP Integration**: Getting and using delegated prefixes
3. **Network Design**: Subnetting, segmentation
4. **Configuration**: Router setup for multiple platforms
5. **Addressing**: SLAAC, static, DHCPv6 strategies
6. **Security**: Multi-layer firewall approach
7. **DNS**: Split-horizon internal resolution
8. **Automation**: Handling prefix changes
9. **Monitoring**: Prometheus, Grafana, logging
10. **Troubleshooting**: Comprehensive debugging techniques

**Key Takeaways:**

- IPv6 security requires actual firewall rules, not just NAT
- Use ULA for internal stability, GUA for internet access
- Automate prefix delegation handling
- Monitor NDP table and connection metrics
- Test PMTUD functionality
- Plan for dual-stack operations

**Next Steps:**

1. Document your current IPv6 allocation
2. Implement ULA addressing
3. Set up monitoring before making changes
4. Test firewall rules thoroughly
5. Automate configuration management
6. Plan for prefix changes
