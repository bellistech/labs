# Lab 5: Building the Clos Network Fabric

## What You'll Learn

- Create a 3-layer Clos network topology in LXD
- Configure FRR (Free Range Routing) for BGP
- Set up VXLAN tunnels between layers
- Understand automatic route discovery
- Build the foundation for multi-datacenter Kubernetes

## Estimated Time

- Understanding: 30 minutes
- Setup: 60 minutes  
- Testing & verification: 45 minutes
- **Total: ~2.5 hours**

## Prerequisites

- Completed Labs 1-4
- LXD installed and initialized
- Comfortable with NixOS configuration
- Understanding of networking basics (from Foundation 04)

---

## Part 1: Understanding Our Target Network

### What We're Building

```
                    CORE LAYER
                  ┌─────────┐
                  │ Core-1  │
                  │ Core-2  │
                  └────┬────┘
                       │
            ┌──────────┼──────────┐
            │          │          │
      DISTRIBUTION LAYER
      ┌─────────┐ ┌─────────┐ ┌─────────┐
      │ Dist-1  │ │ Dist-2  │ │ Dist-3  │
      └────┬────┘ └────┬────┘ └────┬────┘
           │           │           │
        ACCESS LAYER
      ┌─────────┬──────────┬──────────┐
      │ Pod-1   │ Pod-2    │ Pod-3    │
      │ Pod-4   │ Pod-5    │ Pod-6    │
      └─────────┴──────────┴──────────┘
```

### The Job of Each Layer

**Core Layer (2 containers)**
- Route traffic between distribution routers
- Handle inter-datacenter traffic
- Announce routes to rest of network

**Distribution Layer (3 containers)**
- Connect core to access layer
- Discover services (EVPN)
- Load balance traffic
- Redistribute routes

**Access Layer (6 containers)**
- Kubernetes pods/worker nodes
- VXLAN tunnel endpoints
- Local networking

---

## Part 2: Create LXD Containers for Network Fabric

### Step 1: Create Core Routers

```bash
# Create core layer (top of network)
lxc launch images:nixos/unstable core-1
lxc launch images:nixos/unstable core-2

# Verify they're running
lxc list | grep core
```

### Step 2: Create Distribution Layer

```bash
# Create distribution routers (middle layer)
lxc launch images:nixos/unstable dist-1
lxc launch images:nixos/unstable dist-2
lxc launch images:nixos/unstable dist-3

# Verify
lxc list | grep dist
```

### Step 3: Create Access Layer (Pods)

```bash
# Create Kubernetes pods (bottom layer)
lxc launch images:nixos/unstable pod-1
lxc launch images:nixos/unstable pod-2
lxc launch images:nixos/unstable pod-3
lxc launch images:nixos/unstable pod-4
lxc launch images:nixos/unstable pod-5
lxc launch images:nixos/unstable pod-6

# Verify all 12 containers
lxc list
```

### Step 4: Configure LXD Networking

```bash
# Create custom network for fabric
lxc network create fabric-net ipv4.address=10.100.0.1/24 ipv4.dhcp=false

# Attach containers to network (one example)
lxc network attach fabric-net core-1
lxc network attach fabric-net core-2
lxc network attach fabric-net dist-1
lxc network attach fabric-net dist-2
lxc network attach fabric-net dist-3
lxc network attach fabric-net pod-1
# ... attach all pods

# Verify
lxc network show fabric-net
```

---

## Part 3: Configure FRR on Core Routers

### Step 1: Create NixOS Configuration for Core-1

```bash
lxc shell core-1
```

Edit `/etc/nixos/configuration.nix`:

```nix
{ config, pkgs, ... }:

{
  # FRR routing daemon
  services.frr.enable = true;
  services.frr.bgp = {
    enable = true;
    
    # This router's identity
    asn = 65000;  # Autonomous System Number
    routerId = "10.100.1.1";
  };
  
  # Networking
  networking.hostname = "core-1";
  networking.interfaces.eth0.ipv4.addresses = [
    {
      address = "10.100.1.1";
      prefixLength = 24;
    }
  ];
  
  # SSH for administration
  services.openssh.enable = true;
  
  system.stateVersion = "23.11";
}
EOF
```

Rebuild:
```bash
sudo nixos-rebuild switch
```

### Step 2: Configure FRR (BGP Daemon)

The FRR configuration controls how this router talks to others.

```bash
# Edit FRR configuration
sudo nano /etc/frr/frr.conf
```

Add this configuration:

```
! Core-1 FRR Configuration
!
hostname core-1
!
router bgp 65000
  bgp router-id 10.100.1.1
  
  ! Core-2 is in same AS
  neighbor 10.100.2.1 remote-as 65000
  neighbor 10.100.2.1 description "core-2"
  
  ! Distribution routers in lower AS
  neighbor 10.100.10.1 remote-as 65001
  neighbor 10.100.10.1 description "dist-1"
  neighbor 10.100.11.1 remote-as 65002
  neighbor 10.100.11.1 description "dist-2"
  neighbor 10.100.12.1 remote-as 65003
  neighbor 10.100.12.1 description "dist-3"
  
  ! Address family
  address-family ipv4 unicast
    network 10.100.1.0/24
    redistribute connected
  exit-address-family
  
  address-family l2vpn evpn
    neighbor 10.100.2.1 activate
    neighbor 10.100.10.1 activate
    neighbor 10.100.11.1 activate
    neighbor 10.100.12.1 activate
  exit-address-family
!
```

**What this means (ELI5):**
```
"I'm router core-1 with ID 10.100.1.1"
"I belong to AS 65000 (my group of routers)"
"Here are my neighbors:"
  "core-2 (my equal, same AS)"
  "dist-1, dist-2, dist-3 (routers below me)"
"When we talk, use standard BGP"
"Also enable EVPN for service discovery"
"Share my connected networks with neighbors"
```

### Step 3: Start FRR Services

```bash
sudo systemctl restart frr

# Check if it started
sudo systemctl status frr

# View BGP status
sudo vtysh -c "show ip bgp summary"
```

### Step 4: Repeat for Core-2

Repeat Steps 1-3 for core-2:
- Change IP to 10.100.2.1
- Change router-id to 10.100.2.1
- Neighbor to 10.100.1.1 (core-1)

---

## Part 4: Configure Distribution Routers

### Dist-1 Configuration

```bash
lxc shell dist-1
```

Edit `/etc/nixos/configuration.nix`:

```nix
{ config, pkgs, ... }:

{
  # FRR for routing
  services.frr.enable = true;
  services.frr.bgp = {
    enable = true;
    asn = 65001;
    routerId = "10.100.10.1";
  };
  
  # VXLAN support
  boot.kernel.sysctl."net.ipv4.ip_forward" = 1;
  boot.kernel.sysctl."net.bridge.bridge-nf-call-iptables" = 1;
  
  networking.hostname = "dist-1";
  networking.interfaces.eth0.ipv4.addresses = [
    {
      address = "10.100.10.1";
      prefixLength = 24;
    }
  ];
  
  # Enable VXLAN
  networking.interfaces.vxlan100 = {
    virtual = true;
    virtualOwner = "root";
  };
  
  services.openssh.enable = true;
  system.stateVersion = "23.11";
}
```

Rebuild:
```bash
sudo nixos-rebuild switch
```

FRR Configuration for dist-1:

```bash
sudo nano /etc/frr/frr.conf
```

```
hostname dist-1
!
router bgp 65001
  bgp router-id 10.100.10.1
  
  ! Connect to core routers
  neighbor 10.100.1.1 remote-as 65000
  neighbor 10.100.1.1 description "core-1"
  neighbor 10.100.2.1 remote-as 65000
  neighbor 10.100.2.1 description "core-2"
  
  ! Connect to access pods
  neighbor 10.100.20.1 remote-as 65010
  neighbor 10.100.20.1 description "pod-1"
  neighbor 10.100.30.1 remote-as 65020
  neighbor 10.100.30.1 description "pod-4"
  
  address-family ipv4 unicast
    network 10.100.10.0/24
    redistribute connected
  exit-address-family
  
  address-family l2vpn evpn
    neighbor 10.100.1.1 activate
    neighbor 10.100.2.1 activate
    neighbor 10.100.20.1 activate
    neighbor 10.100.30.1 activate
  exit-address-family
!
```

Repeat for dist-2 and dist-3 with appropriate IP addresses.

---

## Part 5: Configure VXLAN Tunnels

### Create VXLAN Tunnel (Core-1 to Dist-1)

On core-1:

```bash
# Create VXLAN interface
sudo ip link add vxlan100 type vxlan id 100 remote 10.100.10.1 local 10.100.1.1 dstport 4789

# Enable it
sudo ip link set vxlan100 up

# Add it to bridge (for Kubernetes)
sudo brctl addif br-vxlan vxlan100

# Verify
ip addr show vxlan100
```

**What this does (ELI5):**
```
"Create a virtual tunnel interface"
  "ID 100 = tunnel name"
  "Remote 10.100.10.1 = other end of tunnel"
  "Local 10.100.1.1 = this end"
  "Port 4789 = standard VXLAN port"
"Turn it on"
"Connect it to bridge so traffic can use it"
```

### Automate with NixOS

Instead of manual commands, add to NixOS config:

```nix
# In configuration.nix
networking.interfaces.vxlan100 = {
  virtual = true;
  virtualOwner = "root";
};

systemd.services.vxlan-setup = {
  description = "Setup VXLAN tunnel";
  after = [ "network-online.target" ];
  wantedBy = [ "multi-user.target" ];
  serviceConfig = {
    Type = "oneshot";
    ExecStart = ''
      ${pkgs.iproute2}/bin/ip link add vxlan100 type vxlan id 100 \
        remote 10.100.10.1 local 10.100.1.1 dstport 4789
      ${pkgs.iproute2}/bin/ip link set vxlan100 up
    '';
  };
};
```

---

## Part 6: Testing the Network

### Step 1: Check BGP Neighbors

```bash
lxc shell core-1
sudo vtysh

# Inside vtysh
show ip bgp neighbors

# Should show:
# Neighbor        V    AS MsgRcvd MsgSent   TblVer  InQ OutQ  Up/Down State/PfxRcd
# 10.100.2.1      4 65000       5       5        2    0    0 00:02:30        0
# 10.100.10.1     4 65001       3       3        2    0    0 00:01:15        0
```

**What this means:**
```
"Show me all neighbors I'm talking to"
"V = IP version (4 = IPv4)"
"AS = Their Autonomous System number"
"State/PfxRcd = Are they connected and how many prefixes learned"
```

### Step 2: Check Routes Learned

```bash
# Still in vtysh
show ip bgp summary

show ip route bgp

# Should show:
# B   10.100.10.0/24 [200/0] via 10.100.10.1, 00:01:20
# B   10.100.11.0/24 [200/0] via 10.100.11.1, 00:01:20
# B   10.100.12.0/24 [200/0] via 10.100.12.1, 00:01:20
```

### Step 3: Ping Test Between Layers

```bash
# From core-1
ping 10.100.10.1  # Should reach dist-1

# From dist-1
ping 10.100.1.1   # Should reach core-1

# From pod-1
ping 10.100.10.1  # Should reach dist-1
```

### Step 4: View EVPN Routes (Advanced)

```bash
# In vtysh on dist-1
show bgp l2vpn evpn route

# Shows which services are where
```

---

## Part 7: Scaling Observations

### What You Should See

```
✓ BGP established between routers
✓ Routes automatically discovered
✓ VXLAN tunnels passing traffic
✓ Pods can reach distribution layer
✓ Distribution layer can reach core
✓ Traffic follows optimal paths

This is automatic - no manual routing!
```

### If Something Doesn't Work

```
Problem: "Neighbor not connecting"
Solution: 
  1. Check IPs are reachable: ping <neighbor-ip>
  2. Check FRR started: systemctl status frr
  3. Check config syntax: vtysh -c "show run"
  4. Check firewall: sudo iptables -L

Problem: "Routes not learned"
Solution:
  1. Check neighbor state is "Established"
  2. Check address-family is enabled
  3. Check redistribute is working
  4. Wait 30 seconds for BGP to converge

Problem: "VXLAN tunnel down"
Solution:
  1. Check remote end is reachable
  2. Check VID (id) is correct
  3. Check port 4789 is open
  4. Verify MAC forwarding
```

---

## Verification Checklist

- [ ] All 12 LXD containers running
- [ ] Core layer routing established
- [ ] Distribution layer connected
- [ ] BGP neighbors showing as "Established"
- [ ] Routes showing in routing table
- [ ] VXLAN tunnels created
- [ ] Ping test between layers successful
- [ ] FRR daemon running on all routers
- [ ] No errors in FRR logs
- [ ] Traffic flowing automatically

---

## What You've Built

✅ **3-layer network topology**
- Core routers handling inter-datacenter traffic
- Distribution routers managing access
- Access layer supporting pods

✅ **Automatic routing**
- BGP discovers routes between layers
- No manual route configuration
- Self-healing on failures

✅ **Service discovery**
- EVPN learns what services are where
- Applications find each other automatically

✅ **Foundation for Kubernetes**
- Network ready for pod deployment
- VXLAN tunnels for pod communication
- All containers can talk

---

## Real-World Comparison

```
What You Built:          Real Data Center:
────────────────         ─────────────────
Core-1, Core-2    →      Core routers (expensive!)
Dist-1/2/3        →      Leaf/spine architecture
Pod-1 to 6        →      Physical servers
VXLAN tunnels     →      Layer 2 fabric
BGP/EVPN          →      Overlay network
FRR               →      Production router OS
```

**You just built a local mockup of modern cloud infrastructure!**

---

## Next: Lab 6

Now that networking is working:
- Deploy Kubernetes across the pods
- Scale applications
- Watch automatic routing handle traffic
- Deploy across multiple datacenters

[Next: Lab 6 - Kubernetes on Clos Network](../lab-06-kubernetes-deployment/README.md)

---

## Reference

- [FRR Documentation](https://docs.frrouting.org/)
- [VXLAN RFC](https://tools.ietf.org/html/rfc7348)
- [BGP RFC](https://tools.ietf.org/html/rfc4271)
- [EVPN RFC](https://tools.ietf.org/html/rfc7432)
- [Clos Network Paper](https://en.wikipedia.org/wiki/Clos_network)
