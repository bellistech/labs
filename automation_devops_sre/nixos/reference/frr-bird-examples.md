# Reference: FRR/BIRD Configuration Examples

## FRR (Free Range Routing) - Complete Examples

### Core Router Configuration (Core-1)

```
! Core-1 BGP Configuration for Clos Network
!

hostname core-1
!
! Basic system config
!
router bgp 65000
  bgp router-id 10.100.1.1
  bgp log-neighbor-changes
  !
  ! iBGP (Internal BGP) - same AS
  neighbor 10.100.2.1 remote-as 65000
  neighbor 10.100.2.1 description "Core-2"
  neighbor 10.100.2.1 update-source 10.100.1.1
  !
  ! eBGP (External BGP) - to distribution layer
  neighbor 10.100.10.1 remote-as 65001
  neighbor 10.100.10.1 description "Distribution-1"
  neighbor 10.100.11.1 remote-as 65002
  neighbor 10.100.11.1 description "Distribution-2"
  neighbor 10.100.12.1 remote-as 65003
  neighbor 10.100.12.1 description "Distribution-3"
  !
  ! IPv4 routes
  address-family ipv4 unicast
    ! Advertise this router's network
    network 10.100.1.0/24
    !
    ! Learn from all neighbors
    neighbor 10.100.2.1 activate
    neighbor 10.100.10.1 activate
    neighbor 10.100.11.1 activate
    neighbor 10.100.12.1 activate
    !
    ! Redistribute directly connected networks
    redistribute connected
  exit-address-family
  !
  ! EVPN (Ethernet VPN) for service discovery
  address-family l2vpn evpn
    neighbor 10.100.2.1 activate
    neighbor 10.100.10.1 activate
    neighbor 10.100.11.1 activate
    neighbor 10.100.12.1 activate
  exit-address-family
!

! Interface configuration
interface eth0
  description "Uplink to network"
  ip address 10.100.1.1/24
  !
  ! BFD for fast failure detection
  ip bfd

! Logging
log stdout
log monitor
!
debug bgp events
debug bgp zebra
!
```

### Distribution Router (Dist-1)

```
hostname dist-1
!
router bgp 65001
  bgp router-id 10.100.10.1
  bgp log-neighbor-changes
  !
  ! Connect to cores
  neighbor 10.100.1.1 remote-as 65000
  neighbor 10.100.1.1 description "Core-1"
  neighbor 10.100.2.1 remote-as 65000
  neighbor 10.100.2.1 description "Core-2"
  !
  ! Connect to pods
  neighbor 10.100.20.1 remote-as 65010
  neighbor 10.100.20.1 description "Pod-1"
  neighbor 10.100.30.1 remote-as 65020
  neighbor 10.100.30.1 description "Pod-4"
  !
  address-family ipv4 unicast
    network 10.100.10.0/24
    neighbor 10.100.1.1 activate
    neighbor 10.100.2.1 activate
    neighbor 10.100.20.1 activate
    neighbor 10.100.30.1 activate
    redistribute connected
  exit-address-family
  !
  address-family l2vpn evpn
    neighbor 10.100.1.1 activate
    neighbor 10.100.2.1 activate
    neighbor 10.100.20.1 activate
    neighbor 10.100.30.1 activate
  exit-address-family
!

interface eth0
  description "To network"
  ip address 10.100.10.1/24
  ip bfd
!
```

---

## BIRD (Alternative Routing Daemon) - Complete Examples

### BIRD Configuration for Core Router

```
# BIRD Configuration - Core-1

router id 10.100.1.1;

# Logging
log stdout all;
log syslog { debug, trace, info, remote, warning, error, auth, fatal };

# Device protocol - learn local IPs
protocol device {
}

# Direct connections
protocol direct {
  interface "*";
}

# BGP Template for common settings
template bgp common_bgp {
  local as 65000;
  
  # Timers (fast convergence)
  connect delay time 5;
  connect retry time 10;
  hold time 30;
  keepalive time 10;
  
  # Features
  graceful restart on;
  bfd on;
  
  # Add all routes to master table
  add paths on;
  add paths limit 2;
}

# BGP Peer: Core-2 (iBGP)
protocol bgp CORE2 from common_bgp {
  neighbor 10.100.2.1 as 65000;
  description "Core-2";
  
  # Next hop self required for iBGP
  next hop self;
}

# BGP Peer: Distribution-1
protocol bgp DIST1 from common_bgp {
  neighbor 10.100.10.1 as 65001;
  description "Distribution-1";
  local as 65000;
  multihop 1;
}

# BGP Peer: Distribution-2
protocol bgp DIST2 from common_bgp {
  neighbor 10.100.11.1 as 65002;
  description "Distribution-2";
  local as 65000;
}

# BGP Peer: Distribution-3
protocol bgp DIST3 from common_bgp {
  neighbor 10.100.12.1 as 65003;
  description "Distribution-3";
  local as 65000;
}

# Static route for loopback
protocol static {
  route 10.100.1.0/24 via 10.100.1.1;
}

# Kernel protocol - sync routes to kernel
protocol kernel {
  kernel table 254;
  learn on;
  scan time 10;
}
```

### BIRD Distribution Router

```
router id 10.100.10.1;

log stdout all;

protocol device {
}

protocol direct {
  interface "*";
}

template bgp clos_bgp {
  local as 65001;
  connect delay time 5;
  connect retry time 10;
  hold time 30;
  keepalive time 10;
  graceful restart on;
  bfd on;
  add paths on;
}

# Upstream connections
protocol bgp CORE1 from clos_bgp {
  neighbor 10.100.1.1 as 65000;
  description "Upstream-Core-1";
  local as 65001;
}

protocol bgp CORE2 from clos_bgp {
  neighbor 10.100.2.1 as 65000;
  description "Upstream-Core-2";
  local as 65001;
}

# Downstream to pods
protocol bgp POD1 from clos_bgp {
  neighbor 10.100.20.1 as 65010;
  description "Pod-1-Leaf";
  local as 65001;
}

protocol bgp POD4 from clos_bgp {
  neighbor 10.100.30.1 as 65020;
  description "Pod-4-Leaf";
  local as 65001;
}

protocol static {
  route 10.100.10.0/24 via 10.100.10.1;
}

protocol kernel {
  kernel table 254;
  learn on;
  scan time 10;
}
```

---

## VXLAN Tunnel Configuration

### Create VXLAN Interface (with BUM handling)

```bash
#!/bin/bash
# Create VXLAN tunnel between two distribution routers

# Parameters
LOCAL_IP="10.100.10.1"
REMOTE_IP="10.100.11.1"
VXLAN_ID="100"
VXLAN_PORT="4789"
INTERFACE_NAME="vxlan100"

# Create VXLAN interface
ip link add $INTERFACE_NAME type vxlan \
  id $VXLAN_ID \
  remote $REMOTE_IP \
  local $LOCAL_IP \
  dstport $VXLAN_PORT \
  dev eth0

# Bring up
ip link set up $INTERFACE_NAME

# Add to bridge (for switching)
brctl addif br-vxlan $INTERFACE_NAME

# Configure IP (if needed)
ip addr add 192.168.100.1/24 dev $INTERFACE_NAME

# Enable ARP suppression (reduce flooding)
bridge link set dev $INTERFACE_NAME neigh_suppress on

echo "VXLAN tunnel created: $INTERFACE_NAME"
ip addr show $INTERFACE_NAME
```

### NixOS VXLAN Systemd Service

```nix
systemd.services.vxlan-setup = {
  description = "Setup VXLAN tunnels for Clos network";
  
  after = [ "network-online.target" ];
  wants = [ "network-online.target" ];
  wantedBy = [ "multi-user.target" ];
  
  serviceConfig = {
    Type = "oneshot";
    RemainAfterExit = true;
    
    ExecStart = ''
      # Create tunnel
      ${pkgs.iproute2}/bin/ip link add vxlan100 type vxlan \
        id 100 remote 10.100.11.1 local 10.100.10.1 \
        dstport 4789 dev eth0
      
      # Bring up
      ${pkgs.iproute2}/bin/ip link set up vxlan100
      
      # Add to bridge
      ${pkgs.bridge-utils}/bin/brctl addif br-vxlan vxlan100 || true
      
      # Enable ARP suppression
      ${pkgs.bridge-utils}/bin/bridge link set dev vxlan100 neigh_suppress on
    '';
    
    ExecStop = ''
      ${pkgs.iproute2}/bin/ip link del vxlan100 || true
    '';
  };
};
```

---

## Monitoring Commands

### Check BGP Status

```bash
# FRR
vtysh -c "show ip bgp summary"
vtysh -c "show ip bgp neighbors"
vtysh -c "show ip bgp route-map"
vtysh -c "show ip route bgp"

# BIRD
birdc show route protocol bgp
birdc show protocol
birdc show status

# Kernel routing table
ip route show
ip route show table all
```

### Monitor VXLAN Traffic

```bash
# Show VXLAN interface stats
ip -s link show vxlan100

# Capture VXLAN packets
tcpdump -i eth0 -n "udp port 4789"

# Watch tunnel data
watch -n 1 'ip -s link show vxlan100'

# Show VXLAN FDB (forwarding database)
bridge fdb show
```

### Kubernetes Network Debugging

```bash
# Show Cilium status
kubectl -n kube-system get pods | grep cilium

# Check service IPs
kubectl get svc -A

# Trace pod traffic
kubectl exec -it pod-name -- traceroute destination-ip

# Monitor network policies
kubectl get networkpolicies -A

# Show pod IPs and nodes
kubectl get pods -A -o wide
```

---

## Common Configuration Patterns

### BGP Route Filtering

```
! FRR - Accept only specific prefixes
router bgp 65000
  !
  ! Create route map
  route-map ACCEPT-PODS permit 10
    match ip address prefix-list PODS
  route-map ACCEPT-PODS deny 20
  !
  ! Apply to neighbor
  neighbor 10.100.10.1 route-map ACCEPT-PODS in
  !
  ! Define prefix list
  ip prefix-list PODS seq 10 permit 192.168.0.0/16
  ip prefix-list PODS seq 20 permit 10.32.0.0/16
!
```

### BIRD Conditional Route Filtering

```
# BIRD - Accept routes only from specific AS
filter import_filter {
  if bgp_path.last ~ [ 65001, 65002, 65003 ] then accept;
  else reject;
}

protocol bgp DIST1 {
  import filter import_filter;
}
```

### BGP Community Tags

```
! FRR - Mark routes with community
route-map ADD-COMMUNITY permit 10
  set community 65000:100
  set community 65000:DC-A additive

route-map TAG-PODS permit 10
  match ip address prefix-list PODS
  set community 65000:PODS additive
```

---

## Troubleshooting Patterns

### BGP Neighbor Not Connecting

```bash
# Check neighbor reachability
ping neighbor-ip

# Check BGP status
show ip bgp neighbors neighbor-ip

# Check if neighbor sees us
ssh neighbor-router show ip bgp neighbors our-ip

# Check password (if using MD5)
vtysh -c "show run router bgp"

# Check AS numbers match expectations
show ip bgp neighbors | grep "remote AS"
```

### Routes Not Being Learned

```bash
# Check if neighbor established
show ip bgp neighbors | grep "Established"

# Check address family activated
show ip bgp neighbors | grep "afi"

# Check if prefix is being announced
show ip bgp summary | grep "PfxRcd"

# Manually check routes
show ip bgp all

# Check redistribute settings
show run router bgp
```

### VXLAN Not Passing Traffic

```bash
# Check interface up
ip addr show vxlan100

# Check remote reachability
ping remote-ip

# Check BUM handling
bridge fdb show

# Monitor VXLAN stats
watch -n 1 'ip -s link show vxlan100'

# Check MTU (should be 1550 for VXLAN overlay)
ip link show | grep mtu
```

---

## Performance Tuning

### BGP Timers for Fast Convergence

```
! Aggressive timers for test/dev
neighbor NEIGHBOR timers 3 10
neighbor NEIGHBOR timers connect 5

! Production timers (more stable)
neighbor NEIGHBOR timers 60 180
neighbor NEIGHBOR timers connect 30
```

### VXLAN MTU Optimization

```bash
# Check path MTU
ip link set mtu 9000 dev eth0  # Jumbo frames
ip link set mtu 1550 dev vxlan100  # Account for VXLAN header

# Test with ping
ping -M do -s 1472 destination
```

### BGP Memory Optimization

```
! Limit route cache
bgp peertype external route-map RM-FILTER in

! Soft-reconfiguration for route refreshes
neighbor NEIGHBOR soft-reconfiguration inbound
```

---

## Reference

- [FRR Configuration Reference](https://docs.frrouting.org/)
- [BIRD User Guide](https://bird.network.cz/)
- [VXLAN RFC 7348](https://tools.ietf.org/html/rfc7348)
- [BGP RFC 4271](https://tools.ietf.org/html/rfc4271)
- [EVPN RFC 7432](https://tools.ietf.org/html/rfc7432)
