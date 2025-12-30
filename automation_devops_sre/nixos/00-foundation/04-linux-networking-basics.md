# Foundation 04: Linux Networking Basics - ELI5

## The Big Picture

Imagine your computer is a city. Networks are like streets connecting buildings.

```
Traditional Single Computer:
  Your Computer = One Building
  
Computer Network:
  Your Computer = One Building in a City
  Network = Streets connecting all buildings
  Router = Traffic controller at intersections
  Firewall = Security guard at city entrance
```

---

## Core Concepts (ELI5 Style)

### IP Address = House Address

```
Your house address: 192.168.1.100

Just like:
  - 123 Main Street, Anytown, USA
  
Packets are like letters:
  - They have a FROM address (your house)
  - They have a TO address (destination house)
  - Mail carrier delivers them (router)
```

### MAC Address = Your House Description

```
MAC Address = Physical building appearance
IP Address = Your house number

Example:
  MAC: aa:bb:cc:dd:ee:ff (what your network card looks like)
  IP: 192.168.1.100 (where your house is on the street)
```

### Subnet = A Neighborhood

```
192.168.1.0/24 = One neighborhood
  192.168.1.1 = Street address for the neighborhood
  192.168.1.100 = Your house
  192.168.1.255 = The end of the street

/24 means "first 24 bits are the neighborhood, last 8 bits are the house number"
```

---

## VXLAN = Virtual Tunnel System

### The Analogy

```
VXLAN = Invisible tunnel through the city

Without VXLAN:
  Computer A → Computer B (Direct street)
  
With VXLAN:
  Computer A → Tunnel entrance
    ↓ (tunnel travels through unknown streets)
    ↓
  Tunnel exit → Computer B

Why use it?
  • Network is too complex to understand
  • Want to group computers together "virtually"
  • Computers are on different physical networks
  • Want to create isolated virtual networks
```

### Real Example

```
Physical Network (what router sees):
  Building 1 (IP: 10.0.0.0/24)
  Building 2 (IP: 10.0.1.0/24)
  
Virtual Network Inside VXLAN Tunnel:
  Apartment 1A (thinks it's 192.168.100.0/24)
  Apartment 1B (thinks it's 192.168.100.0/24)
  
The tunnel "pretends" they're on same network
even though they're in different buildings
```

---

## BGP = Post Office for Routers

### How it works (ELI5)

```
Regular routing:
  "To get to House 5, turn left at Oak Street"
  
BGP routing:
  Post offices announce to each other:
  "I deliver mail to Houses 1-10 on Main Street"
  "I deliver mail to Houses 11-20 on Oak Street"
  
Other post offices remember:
  "Oh! If I need to reach House 15, send to Oak Street post office"
```

### In Network Terms

```
Router A announces:
  "I manage network 10.0.0.0/24"
  
Router B announces:
  "I manage network 10.0.1.0/24"
  
Router C learns:
  "To reach 10.0.0.0/24, ask Router A"
  "To reach 10.0.1.0/24, ask Router B"
```

---

## EVPN = BGP for Layer 2 Networks

### The Analogy

```
Regular BGP = Post office routes
  "I deliver to houses on Main Street"
  
EVPN = Post office knows WHO lives there too
  "I deliver to houses on Main Street"
  "AND John Smith lives at 5 Main Street"
  "AND Jane Doe lives at 10 Main Street"
```

### Why it matters

```
Without EVPN:
  Router just knows "packets to 10.0.0.0/24 go here"
  
With EVPN:
  Router knows:
  "Packets to 10.0.0.0/24 go here"
  "MAC address aa:bb:cc:dd:ee:ff is at location X"
  "Broadcast traffic should go to Y"
```

---

## Clos Network = Perfect Staircase

### The Topology

```
Regular Network (star):
  All routers connect to one center router
  Problem: Center router becomes bottleneck
  
Clos Network (staircase):
  
  Level 3: 2 Core Routers (top of stairs)
           /              \
          /                \
  Level 2: Distribution Layer (middle of stairs)
       /    |    |    \
      /     |    |     \
  Level 1: Access Layer (bottom of stairs)
    Pod1  Pod2  Pod3  Pod4
    
  
  Every pod connects to multiple distribution routers
  Every distribution router connects to multiple core routers
  Result: No bottlenecks, full bandwidth available
```

### Why Clos is Special

```
Traditional:
  Problem: All traffic goes through center
  Result: Slow when busy
  
Clos (like modern data centers):
  Multiple paths available
  Traffic spreads out
  Result: Predictable, non-blocking bandwidth
```

---

## FRR = Routing Daemon (The Traffic Controller)

### What it does

```
Without FRR:
  Manually configure: "Route this IP to that gateway"
  Problem: Complex, error-prone, can't adapt
  
With FRR:
  FRR runs continuously
  Listens to router announcements (BGP, OSPF, etc)
  Automatically updates routes
  Adapts if network fails
  Result: Self-healing network
```

### In Practice

```
Router A (running FRR):
  • Listens to other routers with BGP
  • "Router B says it handles 10.0.1.0/24"
  • "Router C says it handles 10.0.2.0/24"
  • Automatically adds routes
  • If Router B stops responding:
    "Router B died, remove its routes"
    "Wait... Router D just announced it handles 10.0.1.0/24 now"
    "Update routes automatically"
```

---

## BIRD = Alternative Traffic Controller

### Same job, different style

```
FRR = Configuration-heavy routing daemon
      "Tell me exactly what to do"
      Good for: Complex networks, specific control
      
BIRD = Configuration-light routing daemon
       "Figure it out based on these rules"
       Good for: Simple setups, flexible policies
```

### Both do the same core job:

```
1. Listen for router announcements (BGP, OSPF)
2. Build routing table
3. Update kernel routes
4. Adapt when network changes
5. Announce own routes to other routers
```

---

## Kubernetes Pods in This Network

### The Analogy

```
Without network awareness:
  Pod 1 in Datacenter A
  Pod 2 in Datacenter B
  
  They can't talk smoothly
  Traffic bounces around
  Slow and unpredictable
  
With VXLAN + BGP + EVPN:
  Pod 1 in Datacenter A
  Pod 2 in Datacenter B
  
  Network tunnels connect them
  Routes automatically discovered
  All pods appear on same network
  Fast and predictable
```

### What We're Building

```
6 Kubernetes Pods arranged in 2 Virtual Data Centers:

Datacenter 1:
  Pod 1 (2 nodes)
  Pod 2 (2 nodes)
  Pod 3 (2 nodes)
  
Datacenter 2:
  Pod 4 (2 nodes)
  Pod 5 (2 nodes)
  Pod 6 (2 nodes)

Connected by:
  VXLAN tunnels (virtual cables)
  BGP routing (automatic path discovery)
  EVPN (service discovery)
  Clos topology (perfect connectivity)
  
Result: Looks like one giant Kubernetes cluster
even though it's spread across 2 datacenters!
```

---

## The 3-Layer Clos We're Building

### Access Layer (Bottom)
```
6 NixOS Containers (1 per pod)
Each runs:
  • Kubernetes control plane OR worker nodes
  • FRR/BIRD for routing
  • VXLAN tunnel endpoints
```

### Distribution Layer (Middle)
```
4 NixOS Containers (routing fabric)
Each runs:
  • FRR with BGP
  • EVPN service discovery
  • Route redistribution
```

### Core Layer (Top)
```
2 NixOS Containers (central routing)
Each runs:
  • FRR with BGP
  • Full mesh with distribution layer
  • Route summarization
```

---

## Network Diagram (ASCII Art)

```
┌─────────┐      ┌─────────┐
│ Core-1  │──────│ Core-2  │
│ BGP/FRR │      │ BGP/FRR │
└────┬────┘      └────┬────┘
     │                │
  ┌──┴─────────────┬──┴──┐
  │   Full Mesh    │     │
  │                │     │
┌─┴────┐    ┌──────┴──┐ ┌┴─────┐   ┌──────┐
│Dist-1 │    │Dist-2   │ │Dist-3│   │Dist-4│
│EVPN   │────│  BGP    ├─┤VXLAN │───│ FRR  │
└──┬─┬──┘    └────┬────┘ └┬─────┘   └──────┘
   │ └──────┬─────┘       │
   │        │             │
┌──┴─┐   ┌──┴─┐   ┌──────┴──┐   ┌──────┐
│Pod1│   │Pod2│   │Pod3    │   │Pod4 │
│K8s │───│K8s │───│K8s     ├───│K8s  │
└────┘   └────┘   └────────┘   └──────┘

┌──────┐   ┌──────┐
│Pod5 │   │Pod6  │
│K8s  │───│K8s   │
└──────┘   └──────┘

(All connected via VXLAN tunnels and BGP routing)
```

---

## Why This Design

✅ **Scalable** - Add pods easily
✅ **Redundant** - Multiple paths everywhere
✅ **Automatic** - BGP discovers routes
✅ **Realistic** - Matches real data center design
✅ **Learning-Friendly** - Learn real networking concepts
✅ **Local** - Works on single LXD host
✅ **Reproducible** - All in NixOS config

---

## What Happens in Practice

### Scenario: Pod 1 (DC A) talks to Pod 5 (DC B)

```
Step 1: Pod 1 sends packet to Pod 5
  Source: 192.168.100.10 (Pod 1)
  Dest: 192.168.150.10 (Pod 5)
  
Step 2: Pod 1's K8s node checks routing table
  "To reach 192.168.150.0/24, send to distribution layer"
  
Step 3: Packet reaches distribution router
  "Oh, 192.168.150.0/24 is in different datacenter"
  "FRR learned from BGP: use VXLAN tunnel to Dist-4"
  
Step 4: VXLAN tunnel
  Original packet wrapped inside:
  Source: 10.0.2.50 (Dist-1)
  Dest: 10.0.3.50 (Dist-4)
  Inside: 192.168.100.10 → 192.168.150.10
  
Step 5: Packet travels through core routers
  Core-1 → Core-2 (BGP says best path)
  
Step 6: Packet arrives at Dist-4
  VXLAN unwraps inner packet
  "Ah! 192.168.150.10 is local here"
  
Step 7: Packet delivered to Pod 5
  Pod 5 receives packet normally
  Doesn't know it traveled through VXLAN!
  
Step 8: Reply follows same path back
```

---

## Network Services Discovery (EVPN)

### How pods find each other

```
Without EVPN:
  Pod 1 knows Pod 5's IP (manual config)
  Pod 1 sends ARP: "Who has 192.168.150.10?"
  Network floods query everywhere
  Inefficient!
  
With EVPN:
  Pod 5 starts, announces:
  "I'm 192.168.150.10, reach me via Dist-4"
  (EVPN gossip through BGP)
  
  Pod 1 wants to reach Pod 5:
  Checks local cache: "Pod 5 is via Dist-4"
  Sends packet directly
  No flooding, instant routing!
```

---

## Next: Let's Build It

In the upcoming labs, we'll:
1. Create NixOS containers for each layer
2. Configure FRR with BGP and EVPN
3. Set up VXLAN tunnels
4. Deploy Kubernetes across the cluster
5. Scale to multiple datacenters
6. Watch automatic routing work

**All using NixOS declarative configuration!**

---

## Key Takeaways

✅ VXLAN = Virtual tunnels connecting distant networks
✅ BGP = Routers automatically discovering routes
✅ EVPN = Service discovery through BGP
✅ Clos = Perfect topology for data centers
✅ FRR/BIRD = Routing daemons that do the work
✅ Together = Local mockup of real cloud infrastructure

**Ready to build it?** → Lab 5: Setting up the network fabric
