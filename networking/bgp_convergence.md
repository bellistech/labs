# Faster BGP Convergence — Personal Notes (ELI5 → Deep Dive)

---

## What “BGP convergence” means (ELI5 version)

Imagine the internet like a city full of roads.

- **BGP** is the map that says which roads exist.
- **Convergence** is how fast everyone updates their map after a road breaks.
- Slow convergence = cars keep driving into a dead end.
- Fast convergence = detour signs flip instantly and traffic keeps moving.

The trick:
**Don’t wait for the map to redraw if you already know a detour.**

Convergence has four stages:

1. **Detect** the failure  
2. **Fail over** traffic immediately  
3. **Recalculate and propagate routes**  
4. **Stay stable while things churn**  

Fast networks optimize the first two so the last two hurt less.

---

## 1) Failure Detection — “How fast do I notice the road is broken?”

### Bidirectional Forwarding Detection (BFD)

**ELI5:**  
Routers tap each other on the shoulder many times per second.  
If the tapping stops, the link is dead.

- Sub-second detection
- Independent of BGP
- Best option in almost all cases

### FRR — BFD with BGP
```bash
bfd
 peer 192.0.2.2
  interval 300 min_rx 300 multiplier 3
```

```bash
router bgp 65001
 neighbor 192.0.2.2 remote-as 65002
 neighbor 192.0.2.2 bfd
```

### BIRD — BFD with BGP
```bird
protocol bfd {
  interface "*";
  interval 300 ms;
  multiplier 3;
}
```

```bird
protocol bgp upstream {
  neighbor 192.0.2.2 as 65002;
  bfd yes;
}
```

---

### Fast External Fallover

**ELI5:**  
“If the cable is unplugged, stop talking immediately.”

- Only works for direct physical link-down events
- Does not detect upstream blackholes

FRR example:
```bash
router bgp 65001
 neighbor 192.0.2.2 fall-over
```

---

### Reduced BGP Timers (last resort)

**ELI5:**  
Check on your friend every 3 seconds instead of every 30.

FRR example:
```bash
router bgp 65001
 neighbor 192.0.2.2 timers 3 9
```

Use only if BFD is unavailable.

---

## 2) Failover — “How fast does traffic move?”

### BGP Prefix Independent Convergence (PIC)

**ELI5:**  
The detour is already programmed into the router before anything breaks.

- Immediate data-plane failover
- Prefix count does not matter
- Requires ECMP or alternate next-hops

FRR note:
PIC is automatic when multiple next-hops exist.

---

## 3) Control Plane Reaction — “How fast does the map update?”

### Next-Hop Tracking (NHT)

**ELI5:**  
Instead of checking every minute, BGP gets notified instantly when a next-hop disappears.

FRR example:
```bash
router bgp 65001
 bgp nexthop trigger delay 0
```

BIRD (automatic with kernel + IGP):
```bird
protocol kernel {
  learn;
}
```

---

### Fast IGP Convergence

**ELI5:**  
BGP can only be as fast as the IGP underneath it.

FRR OSPF tuning:
```bash
router ospf
 timers throttle spf 50 200 5000
 timers throttle lsa 50 200 5000
```

---

### Route Reflector Design

- RRs are convergence choke points
- Always deploy redundant RRs
- Watch CPU and memory closely

---

### Add-Path (iBGP)

- Advertise multiple paths for same prefix
- Reduces path hunting
- Increases memory usage

---

## 4) Propagation and Stability

### MRAI

**ELI5:**  
“Wait a moment before yelling again so you don’t spam everyone.”

FRR example:
```bash
router bgp 65001
 neighbor 192.0.2.2 advertisement-interval 0
```

Use sparingly.

---

### Graceful Restart

**ELI5:**  
“Pretend I didn’t disappear while rebooting.”

FRR:
```bash
router bgp 65001
 bgp graceful-restart
```

Can cause blackholing if misused.

---

## Stability & Protection

- Control Plane Policing (CoPP)
- QoS for routing protocols
- Max-prefix limits
- Be careful with route flap dampening

---

## Scaling Notes

- Convergence is CPU and FIB-install bound
- Fewer routes = faster recovery
- Design beats timer tuning

---

## Condensed Cheat Sheet

- **Best detection:** BFD  
- **Best failover:** PIC  
- **Fast iBGP:** NHT + fast IGP  
- **Hidden bottleneck:** Route Reflectors  
- **Last resort:** Timer tuning  

---
