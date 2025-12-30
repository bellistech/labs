# Network Automation Lab Exercises & GitHub Project Templates

## Quick Start GitHub Projects

### Project 1: BGP Network Monitor (Go)
```go
// main.go - BGP Session Monitor in Go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type BGPSession struct {
    PeerIP      string    `json:"peer_ip"`
    ASN         int       `json:"asn"`
    State       string    `json:"state"`
    Uptime      int64     `json:"uptime"`
    PrefixCount int       `json:"prefix_count"`
    LastChecked time.Time `json:"last_checked"`
}

type BGPMonitor struct {
    sessions map[string]*BGPSession
    mu       sync.RWMutex
    metrics  *BGPMetrics
}

type BGPMetrics struct {
    sessionsUp   prometheus.Gauge
    sessionsDown prometheus.Gauge
    prefixCount  prometheus.Gauge
}

// TODO: Implement these functions
// - NewBGPMonitor()
// - CheckBGPSession()
// - GetSessionStatus()
// - UpdateMetrics()
// - ServeHTTP handlers

func main() {
    // Initialize monitor
    // Start monitoring goroutines
    // Setup HTTP API
    // Setup Prometheus metrics
    // Run server
}
```

### Project 2: EVPN/VXLAN Configuration Generator (Python)
```python
#!/usr/bin/env python3
"""
EVPN/VXLAN Configuration Generator
Generates vendor-specific configs for spine-leaf EVPN fabric
"""

import yaml
import jinja2
from typing import Dict, List
import click
import ipaddress

class EVPNFabric:
    def __init__(self, topology_file: str):
        self.topology = self._load_topology(topology_file)
        self.vteps = []
        self.spines = []
        self.leafs = []
        
    def _load_topology(self, file_path: str) -> Dict:
        """Load topology from YAML file"""
        with open(file_path, 'r') as f:
            return yaml.safe_load(f)
    
    def generate_configs(self, vendor: str = 'juniper') -> Dict[str, str]:
        """Generate device configurations"""
        configs = {}
        
        # Generate spine configs
        for spine in self.topology['spines']:
            configs[spine['hostname']] = self._generate_spine_config(spine, vendor)
        
        # Generate leaf configs
        for leaf in self.topology['leafs']:
            configs[leaf['hostname']] = self._generate_leaf_config(leaf, vendor)
        
        return configs
    
    def _generate_spine_config(self, spine: Dict, vendor: str) -> str:
        """Generate spine switch configuration"""
        # Template loading and rendering
        pass
    
    def _generate_leaf_config(self, leaf: Dict, vendor: str) -> str:
        """Generate leaf switch configuration"""
        # Template loading and rendering
        pass
    
    def validate_config(self, config: str, vendor: str) -> bool:
        """Validate generated configuration"""
        # Add validation logic
        pass

@click.command()
@click.option('--topology', '-t', required=True, help='Topology YAML file')
@click.option('--vendor', '-v', default='juniper', help='Vendor (juniper/arista/cisco)')
@click.option('--output', '-o', default='configs/', help='Output directory')
def main(topology, vendor, output):
    """Generate EVPN/VXLAN configurations"""
    fabric = EVPNFabric(topology)
    configs = fabric.generate_configs(vendor)
    
    # Save configurations
    for device, config in configs.items():
        with open(f"{output}/{device}.conf", 'w') as f:
            f.write(config)
    
    click.echo(f"Generated {len(configs)} configurations")

if __name__ == "__main__":
    main()
```

### Project 3: Network Automation CI/CD Pipeline (Jenkins)
```groovy
// Jenkinsfile - Network Automation Pipeline
pipeline {
    agent any
    
    environment {
        ANSIBLE_HOST_KEY_CHECKING = 'False'
        GIT_REPO = 'https://github.com/yourusername/network-automation'
        SLACK_CHANNEL = '#network-ops'
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: "${GIT_REPO}"
            }
        }
        
        stage('Syntax Validation') {
            parallel {
                stage('Validate YAML') {
                    steps {
                        sh 'yamllint inventory/*.yml'
                    }
                }
                stage('Validate Jinja2') {
                    steps {
                        sh 'python scripts/validate_templates.py'
                    }
                }
                stage('Ansible Syntax') {
                    steps {
                        sh 'ansible-playbook --syntax-check playbooks/*.yml'
                    }
                }
            }
        }
        
        stage('Unit Tests') {
            steps {
                sh 'pytest tests/ -v --junit-xml=results.xml'
            }
        }
        
        stage('Config Generation') {
            steps {
                sh 'python scripts/generate_configs.py --env staging'
            }
        }
        
        stage('Pre-Deployment Validation') {
            steps {
                sh 'python scripts/validate_configs.py'
            }
        }
        
        stage('Deploy to Staging') {
            steps {
                sh 'ansible-playbook -i inventory/staging.yml playbooks/deploy.yml --check'
            }
        }
        
        stage('Integration Tests') {
            steps {
                sh 'robot tests/integration/'
            }
        }
        
        stage('Approval') {
            when {
                branch 'main'
            }
            steps {
                input message: 'Deploy to Production?', ok: 'Deploy'
            }
        }
        
        stage('Production Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh 'ansible-playbook -i inventory/production.yml playbooks/deploy.yml'
            }
        }
        
        stage('Post-Deployment Verification') {
            steps {
                sh 'python scripts/verify_deployment.py'
            }
        }
    }
    
    post {
        always {
            junit 'results.xml'
            archiveArtifacts artifacts: 'configs/**/*.conf', allowEmptyArchive: true
        }
        success {
            slackSend channel: "${SLACK_CHANNEL}", 
                     color: 'good', 
                     message: "Network deployment successful: ${env.JOB_NAME} - ${env.BUILD_NUMBER}"
        }
        failure {
            slackSend channel: "${SLACK_CHANNEL}", 
                     color: 'danger', 
                     message: "Network deployment failed: ${env.JOB_NAME} - ${env.BUILD_NUMBER}"
        }
    }
}
```

### Project 4: Kubernetes Network Operator with Cilium
```yaml
# cilium-network-policy.yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: bgp-peering-policy
  namespace: network-automation
spec:
  endpointSelector:
    matchLabels:
      app: bgp-speaker
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: route-reflector
    toPorts:
    - ports:
      - port: "179"
        protocol: TCP
  egress:
  - toEndpoints:
    - matchLabels:
        app: route-reflector
    toPorts:
    - ports:
      - port: "179"
        protocol: TCP
  - toFQDNs:
    - matchPattern: "*.monitoring.local"
    toPorts:
    - ports:
      - port: "9090"
        protocol: TCP
---
# bgp-config.yaml
apiVersion: cilium.io/v2alpha1
kind: CiliumBGPPeeringPolicy
metadata:
  name: datacenter-bgp
spec:
  nodeSelector:
    matchLabels:
      bgp: enabled
  virtualRouters:
  - localASN: 65001
    exportPodCIDR: true
    neighbors:
    - peerAddress: 10.0.0.1/32
      peerASN: 65000
      connectRetryTimeSeconds: 30
      holdTimeSeconds: 90
      keepAliveTimeSeconds: 30
    - peerAddress: 10.0.0.2/32
      peerASN: 65000
      connectRetryTimeSeconds: 30
      holdTimeSeconds: 90
      keepAliveTimeSeconds: 30
```

### Project 5: Ansible Network Automation Playbooks
```yaml
# playbooks/configure_evpn.yml
---
- name: Configure EVPN/VXLAN Fabric
  hosts: network_devices
  gather_facts: no
  vars:
    vtep_loopback: "{{ hostvars[inventory_hostname]['vtep_ip'] }}"
    
  tasks:
    - name: Configure Underlay BGP
      block:
        - name: Configure BGP for spines
          junos_bgp:
            config:
              as_number: "{{ bgp_asn }}"
              router_id: "{{ router_id }}"
              neighbors: "{{ bgp_neighbors }}"
          when: device_role == "spine"
          
        - name: Configure BGP for leafs
          junos_bgp:
            config:
              as_number: "{{ bgp_asn }}"
              router_id: "{{ router_id }}"
              neighbors: "{{ bgp_neighbors }}"
          when: device_role == "leaf"
      tags: underlay
      
    - name: Configure EVPN Overlay
      junos_config:
        lines:
          - set protocols evpn encapsulation vxlan
          - set protocols evpn extended-vni-list all
          - set switch-options vtep-source-interface lo0.0
          - set switch-options route-distinguisher {{ router_id }}:1
          - set switch-options vrf-target target:{{ bgp_asn }}:1
      when: device_role == "leaf"
      tags: overlay
      
    - name: Configure VXLAN
      junos_vlans:
        config:
          - name: "{{ item.name }}"
            vlan_id: "{{ item.vlan_id }}"
            vxlan:
              vni: "{{ item.vni }}"
      loop: "{{ vxlan_vlans }}"
      when: device_role == "leaf"
      tags: vxlan
      
    - name: Verify EVPN Status
      junos_command:
        commands:
          - show evpn database
          - show bgp summary
          - show ethernet-switching vxlan-tunnel-end-point remote
      register: evpn_output
      tags: verify
      
    - name: Save configuration
      junos_config:
        commit: yes
        comment: "EVPN configuration deployed by Ansible"
      tags: commit
```

## Lab Environment Setup Scripts

### Quick Lab Setup (bash)
```bash
#!/bin/bash
# setup_lab.sh - Quick Network Lab Environment Setup

set -e

echo "Setting up Network Automation Lab Environment..."

# Install required packages
install_dependencies() {
    echo "Installing dependencies..."
    sudo apt-get update
    sudo apt-get install -y \
        docker.io \
        docker-compose \
        python3-pip \
        golang-go \
        ansible \
        git \
        jq \
        yamllint \
        tree
    
    # Python packages
    pip3 install --user \
        netmiko \
        napalm \
        junos-eznc \
        pyyaml \
        jinja2 \
        pytest \
        flask \
        fastapi \
        uvicorn
}

# Setup containerized network devices
setup_network_containers() {
    echo "Setting up network device containers..."
    
    # Create docker-compose.yml for cEOS
    cat > docker-compose.yml <<EOF
version: '3.8'

services:
  spine1:
    image: ceosimage:latest
    container_name: spine1
    hostname: spine1
    environment:
      CEOS: 1
      container: docker
      ETBA: 1
      SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT: 1
      INTFTYPE: eth
    networks:
      mgmt:
        ipv4_address: 172.20.0.11
    volumes:
      - ./configs/spine1:/mnt/flash
    privileged: true
    
  spine2:
    image: ceosimage:latest
    container_name: spine2
    hostname: spine2
    environment:
      CEOS: 1
      container: docker
      ETBA: 1
      SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT: 1
      INTFTYPE: eth
    networks:
      mgmt:
        ipv4_address: 172.20.0.12
    volumes:
      - ./configs/spine2:/mnt/flash
    privileged: true
    
  leaf1:
    image: ceosimage:latest
    container_name: leaf1
    hostname: leaf1
    environment:
      CEOS: 1
      container: docker
      ETBA: 1
      SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT: 1
      INTFTYPE: eth
    networks:
      mgmt:
        ipv4_address: 172.20.0.21
    volumes:
      - ./configs/leaf1:/mnt/flash
    privileged: true
    
  leaf2:
    image: ceosimage:latest
    container_name: leaf2
    hostname: leaf2
    environment:
      CEOS: 1
      container: docker
      ETBA: 1
      SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT: 1
      INTFTYPE: eth
    networks:
      mgmt:
        ipv4_address: 172.20.0.22
    volumes:
      - ./configs/leaf2:/mnt/flash
    privileged: true

networks:
  mgmt:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/24
EOF
}

# Setup Jenkins
setup_jenkins() {
    echo "Setting up Jenkins..."
    docker run -d \
        --name jenkins \
        -p 8080:8080 \
        -p 50000:50000 \
        -v jenkins_home:/var/jenkins_home \
        jenkins/jenkins:lts
    
    echo "Jenkins initial password:"
    sleep 10
    docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
}

# Setup Kubernetes with Kind
setup_kubernetes() {
    echo "Setting up Kubernetes..."
    
    # Install kind
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.17.0/kind-linux-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
    
    # Create cluster
    kind create cluster --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
networking:
  disableDefaultCNI: true
EOF
    
    # Install Cilium
    curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz
    sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
    rm cilium-linux-amd64.tar.gz
    
    cilium install
    cilium status --wait
}

# Setup monitoring stack
setup_monitoring() {
    echo "Setting up Prometheus and Grafana..."
    
    # Create monitoring docker-compose
    cat > monitoring-compose.yml <<EOF
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
    networks:
      - monitoring
    
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - monitoring

volumes:
  prometheus_data:
  grafana_data:

networks:
  monitoring:
    driver: bridge
EOF
    
    docker-compose -f monitoring-compose.yml up -d
}

# Main execution
main() {
    install_dependencies
    setup_network_containers
    setup_jenkins
    setup_kubernetes
    setup_monitoring
    
    echo "Lab environment setup complete!"
    echo "Access points:"
    echo "  Jenkins: http://localhost:8080"
    echo "  Prometheus: http://localhost:9090"
    echo "  Grafana: http://localhost:3000 (admin/admin)"
    echo ""
    echo "Network devices:"
    echo "  Spine1: 172.20.0.11"
    echo "  Spine2: 172.20.0.12"
    echo "  Leaf1: 172.20.0.21"
    echo "  Leaf2: 172.20.0.22"
}

main
```

## Daily Practice Problems

### Day 1 Practice: Python Network Automation
```python
# Challenge: Build a multi-threaded network device configuration backup tool
import concurrent.futures
import datetime
import os
from netmiko import ConnectHandler
from typing import Dict, List
import yaml

def backup_device(device: Dict) -> Dict:
    """Backup configuration from a single device"""
    try:
        connection = ConnectHandler(**device)
        config = connection.send_command("show running-config")
        connection.disconnect()
        
        # Save to file
        timestamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"backups/{device['host']}_{timestamp}.conf"
        
        os.makedirs("backups", exist_ok=True)
        with open(filename, 'w') as f:
            f.write(config)
        
        return {'device': device['host'], 'status': 'success', 'file': filename}
    except Exception as e:
        return {'device': device['host'], 'status': 'failed', 'error': str(e)}

def main():
    # Load devices from inventory
    with open('inventory.yml', 'r') as f:
        devices = yaml.safe_load(f)['devices']
    
    # Parallel execution
    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
        results = list(executor.map(backup_device, devices))
    
    # Report results
    for result in results:
        print(f"{result['device']}: {result['status']}")

if __name__ == "__main__":
    main()
```

### Day 2 Practice: Go BGP Parser
```go
// Challenge: Parse BGP table and find best paths
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type BGPRoute struct {
    Network    string
    NextHop    string
    Metric     int
    LocalPref  int
    Weight     int
    Path       []int
    Origin     string
    Valid      bool
    Best       bool
}

func parseBGPTable(filename string) ([]BGPRoute, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var routes []BGPRoute
    scanner := bufio.NewScanner(file)
    
    // Regex patterns for parsing
    routePattern := regexp.MustCompile(`^[\*>]\s+(\S+)\s+(\S+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.+)`)
    
    for scanner.Scan() {
        line := scanner.Text()
        matches := routePattern.FindStringSubmatch(line)
        
        if len(matches) > 0 {
            route := BGPRoute{
                Network:   matches[1],
                NextHop:   matches[2],
                Metric:    parseIntOrZero(matches[3]),
                LocalPref: parseIntOrZero(matches[4]),
                Weight:    parseIntOrZero(matches[5]),
                Path:      parseASPath(matches[6]),
                Valid:     strings.Contains(line, "*"),
                Best:      strings.Contains(line, ">"),
            }
            routes = append(routes, route)
        }
    }
    
    return routes, scanner.Err()
}

func parseIntOrZero(s string) int {
    val, _ := strconv.Atoi(s)
    return val
}

func parseASPath(path string) []int {
    var asPath []int
    parts := strings.Fields(path)
    for _, part := range parts {
        if as, err := strconv.Atoi(part); err == nil {
            asPath = append(asPath, as)
        }
    }
    return asPath
}

func findBestPath(routes []BGPRoute) *BGPRoute {
    // Implement BGP best path selection algorithm
    // 1. Highest Weight
    // 2. Highest Local Preference
    // 3. Locally originated
    // 4. Shortest AS Path
    // 5. Lowest Origin Type
    // 6. Lowest MED
    // 7. eBGP over iBGP
    // 8. Lowest IGP metric
    // 9. Oldest route
    // 10. Lowest Router ID
    // 11. Lowest peer IP
    
    // Simplified version
    var best *BGPRoute
    for i := range routes {
        if !routes[i].Valid {
            continue
        }
        
        if best == nil {
            best = &routes[i]
            continue
        }
        
        // Compare based on BGP attributes
        if routes[i].Weight > best.Weight {
            best = &routes[i]
        } else if routes[i].Weight == best.Weight {
            if routes[i].LocalPref > best.LocalPref {
                best = &routes[i]
            } else if routes[i].LocalPref == best.LocalPref {
                if len(routes[i].Path) < len(best.Path) {
                    best = &routes[i]
                }
            }
        }
    }
    
    return best
}

func main() {
    routes, err := parseBGPTable("bgp_table.txt")
    if err != nil {
        fmt.Printf("Error parsing BGP table: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Parsed %d routes\n", len(routes))
    
    // Group by network
    networkMap := make(map[string][]BGPRoute)
    for _, route := range routes {
        networkMap[route.Network] = append(networkMap[route.Network], route)
    }
    
    // Find best path for each network
    for network, routes := range networkMap {
        best := findBestPath(routes)
        if best != nil {
            fmt.Printf("Best path for %s: via %s, AS Path: %v\n", 
                network, best.NextHop, best.Path)
        }
    }
}
```

## Interview Questions Bank

### Technical Questions
1. **BGP**: Explain BGP path selection algorithm and how to influence it
2. **EVPN/VXLAN**: How does EVPN Type-2 route work? What's in a Type-5 route?
3. **Python**: How would you handle concurrent API calls to 100 network devices?
4. **Go**: Explain goroutines vs threads. How do channels work?
5. **Kubernetes**: How does Cilium implement network policies using eBPF?
6. **CI/CD**: Design a zero-downtime network upgrade pipeline
7. **Monitoring**: How would you detect and alert on BGP flapping?
8. **Ansible**: How to handle idempotency in network automation?
9. **Terraform**: How does Terraform handle state in network automation?
10. **Linux**: Explain network namespaces and their use in containerization

### Scenario-Based Questions
1. **Troubleshooting**: "BGP session is established but no routes are being received"
2. **Design**: "Design automation for 50-location WAN with different vendors"
3. **Incident**: "Production EVPN fabric has a split-brain condition"
4. **Scaling**: "Current automation takes 2 hours for 1000 devices, how to improve?"
5. **Security**: "Implement zero-trust networking in existing infrastructure"

### Behavioral Questions (STAR Format)
1. Tell me about a complex network issue you automated
2. Describe a time when automation caused an outage
3. How did you handle conflicting priorities in on-call situations?
4. Example of learning a new technology quickly
5. Time when you improved team processes

## Resources & References

### Documentation Links
- BGP: https://www.rfc-editor.org/rfc/rfc4271
- EVPN: https://www.rfc-editor.org/rfc/rfc7432
- VXLAN: https://www.rfc-editor.org/rfc/rfc7348
- Cilium: https://docs.cilium.io/
- Ansible Network: https://docs.ansible.com/ansible/latest/network/
- Go Networking: https://pkg.go.dev/net

### GitHub Repositories to Study
- https://github.com/networktocode/ntc-templates
- https://github.com/napalm-automation/napalm
- https://github.com/aristanetworks/goarista
- https://github.com/cilium/cilium

### Books (Quick Reference)
- "Network Programmability and Automation" - Jason Edelman
- "BGP Design and Implementation" - Randy Zhang
- "Cloud Native DevOps with Kubernetes" - John Arundel

Remember: Focus on demonstrating practical skills through your GitHub portfolio. The interviewer wants to see you can build real solutions, not just memorize concepts.
