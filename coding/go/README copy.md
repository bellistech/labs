# Go Complete Course: From Zero to Production Systems

A comprehensive Go programming course that takes you from absolute beginner to building production network services. The course culminates in building a fully functional IPv4/IPv6 dual-stack authoritative DNS server.

## Course Overview

**Target Audience:** Script writers and programmers familiar with Python, Ruby, Bash, or similar languages wanting to build production backend systems.

**Duration:** Self-paced, approximately 20-30 hours of material

**Final Project:** A complete authoritative DNS server supporting A, AAAA, CNAME, MX, TXT, NS records with IPv4/IPv6 dual-stack support.

## Course Structure

### Part I: Go Fundamentals (Week 1)
1. Introduction & Setup
2. Variables, Types & Constants
3. Control Flow
4. Functions
5. Data Structures (Arrays, Slices, Maps)

### Part II: Intermediate Go (Week 2)
6. Structs & Methods
7. Interfaces
8. Error Handling
9. Packages & Modules
10. Testing

### Part III: Concurrency (Week 3)
11. Goroutines
12. Channels
13. Select & Timeouts
14. Sync Primitives (Mutex, RWMutex, Once, Pool)
15. Context

### Part IV: Systems Programming (Week 4)
16. File I/O
17. Binary Data & Encoding
18. Network Programming Basics
19. TCP Servers & Clients
20. UDP Servers & Clients

### Part V: Capstone - DNS Server (Week 4 continued)
21. DNS Protocol Deep Dive
22. Parsing DNS Messages
23. Building DNS Responses
24. Zone File Parsing
25. Complete DNS Server
26. Testing & Deployment

### Part VI: Traceroute Clone (Week 5) ⭐ NEW!
27. How Internet Routing Works
28. TTL (Time To Live) Explained
29. ICMP Protocol Deep Dive
30. Raw Sockets in Go
31. Building a Complete Traceroute
32. Reverse DNS Lookups

## Contents

```
go-course/
├── README.md                           # This file
├── docs/
│   ├── GO_COMPLETE_COURSE_PART1.md    # Fundamentals
│   ├── GO_COMPLETE_COURSE_PART2.md    # Intermediate
│   ├── GO_COMPLETE_COURSE_PART3.md    # Concurrency & Systems
│   ├── GO_COMPLETE_COURSE_PART4.md    # DNS Capstone
│   ├── GO_COMPLETE_COURSE_PART5_TRACEROUTE.md  # Traceroute (ELI5 style!)
│   └── GO_CHEAT_SHEET.md              # Quick reference
├── dns-server/                         # Working DNS server project
│   ├── README.md
│   ├── Makefile
│   ├── go.mod
│   ├── cmd/dns-server/main.go
│   ├── dns/
│   │   ├── types.go
│   │   ├── parser.go + parser_test.go
│   │   ├── builder.go
│   │   └── zone.go + zone_test.go
│   └── zones/example.com.zone
├── traceroute/                         # ⭐ NEW! Traceroute clone project
│   ├── README.md
│   ├── Makefile
│   ├── go.mod
│   └── main.go                        # ~500 lines, heavily commented
└── examples/
    ├── basics/
    ├── concurrency/
    └── networking/
```

## Quick Start

### Prerequisites

```bash
# Install Go (1.21 or later)
# macOS
brew install go

# Ubuntu/Debian
sudo apt-get install golang-go

# Verify
go version
```

### Run the DNS Server

```bash
cd dns-server

# Build
go build -o dns-server ./cmd/dns-server

# Run (uses port 5353)
./dns-server -zone zones/example.com.zone

# In another terminal, test with dig
dig @localhost -p 5353 example.com A
dig @localhost -p 5353 example.com AAAA
dig @localhost -p 5353 example.com MX
dig @localhost -p 5353 ftp.example.com A
```

### Run the Traceroute Clone (Week 5) ⭐ NEW!

```bash
cd traceroute

# Build
go build -o traceroute main.go

# Run (requires sudo for raw sockets!)
sudo ./traceroute google.com
sudo ./traceroute 8.8.8.8
sudo ./traceroute amazon.com
```

## What You'll Learn

| Topic | Skills |
|-------|--------|
| Go Basics | Variables, types, functions, control flow |
| Data Structures | Slices, maps, structs |
| OOP in Go | Methods, interfaces, composition |
| Error Handling | Custom errors, wrapping, best practices |
| Concurrency | Goroutines, channels, mutexes |
| Context | Cancellation, timeouts, values |
| Networking | TCP/UDP servers, IPv4/IPv6 |
| Binary Protocols | Parsing/building wire formats |
| Testing | Unit tests, benchmarks, table-driven tests |

## Key Concepts by Section

### Part I: Fundamentals
- Go's type system and zero values
- Slices vs arrays
- Map operations and iteration
- String handling and Unicode

### Part II: Intermediate
- Value vs pointer receivers
- Interface satisfaction
- Error wrapping (Go 1.13+)
- Table-driven tests

### Part III: Concurrency
- Goroutine lifecycle
- Channel patterns (generator, pipeline, fan-out/fan-in)
- select for multiplexing
- Context for cancellation

### Part IV: Systems Programming
- Binary encoding (big-endian, little-endian)
- UDP connectionless protocols
- IPv4/IPv6 dual-stack
- Graceful shutdown

### Part V: DNS Server
- DNS wire format (RFC 1035)
- Name compression
- Zone file parsing
- Authoritative responses

## DNS Server Features

The capstone DNS server includes:

- ✅ IPv4 and IPv6 dual-stack support
- ✅ A, AAAA, CNAME, MX, NS, TXT record types
- ✅ BIND-style zone file parsing
- ✅ Concurrent query handling
- ✅ NXDOMAIN responses
- ✅ Statistics tracking
- ✅ Graceful shutdown

## Recommended Learning Path

1. **Week 1**: Parts 1-2 (Fundamentals + Intermediate)
   - Complete all code examples
   - Write unit tests for exercises

2. **Week 2**: Part 3 (Concurrency)
   - Focus on channel patterns
   - Build a simple concurrent program

3. **Week 3**: Part 3 continued (Systems Programming)
   - Build TCP echo server
   - Build UDP client/server

4. **Week 4**: Parts 4-5 (DNS Server Capstone)
   - Follow along building the DNS server
   - Add your own record types
   - Deploy and test

5. **Week 5**: Part 6 (Traceroute Clone) ⭐ NEW!
   - Learn how internet routing works
   - Understand ICMP and raw sockets
   - Build a real network diagnostic tool
   - Extend with IPv6 or geographic lookups

## Next Projects

After completing this course, try:

1. **HTTP Server** - Build a REST API with net/http
2. **gRPC Service** - Protocol Buffers and gRPC
3. **Database App** - PostgreSQL with database/sql
4. **CLI Tool** - Build with cobra/viper
5. **Metrics Collector** - Like Prometheus node_exporter

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [RFC 1035 - DNS](https://datatracker.ietf.org/doc/html/rfc1035)

## License

MIT License - See LICENSE file
