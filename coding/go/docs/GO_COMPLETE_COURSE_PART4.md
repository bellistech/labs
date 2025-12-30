# Part V: Capstone - DNS Server

In this final section, we build a complete authoritative DNS server from scratch. This project combines everything we've learned: binary protocols, UDP networking, concurrent programming, and proper Go project structure.

## 21. DNS Protocol Deep Dive

### 21.1 DNS Message Structure

```
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      ID                       |  16 bits
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |  16 bits (flags)
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QDCOUNT                    |  16 bits
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANCOUNT                    |  16 bits
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    NSCOUNT                    |  16 bits
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ARCOUNT                    |  16 bits
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   QUESTION                    |  Variable
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANSWER                     |  Variable
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   AUTHORITY                   |  Variable
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   ADDITIONAL                  |  Variable
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

### 21.2 DNS Record Types

| Type | Value | Description |
|------|-------|-------------|
| A | 1 | IPv4 address |
| NS | 2 | Name server |
| CNAME | 5 | Canonical name (alias) |
| SOA | 6 | Start of authority |
| MX | 15 | Mail exchange |
| TXT | 16 | Text record |
| AAAA | 28 | IPv6 address |

### 21.3 DNS Classes

| Class | Value | Description |
|-------|-------|-------------|
| IN | 1 | Internet (most common) |
| CS | 2 | CSNET (obsolete) |
| CH | 3 | CHAOS |
| HS | 4 | Hesiod |

---

## 22. Parsing DNS Messages

### 22.1 DNS Types and Constants

```go
// File: dns/types.go
package dns

import (
    "fmt"
    "net"
)

// DNS record types
const (
    TypeA     uint16 = 1
    TypeNS    uint16 = 2
    TypeCNAME uint16 = 5
    TypeSOA   uint16 = 6
    TypeMX    uint16 = 15
    TypeTXT   uint16 = 16
    TypeAAAA  uint16 = 28
)

// DNS classes
const (
    ClassIN uint16 = 1  // Internet
)

// DNS response codes
const (
    RcodeNoError        uint8 = 0
    RcodeFormatError    uint8 = 1
    RcodeServerFailure  uint8 = 2
    RcodeNameError      uint8 = 3  // NXDOMAIN
    RcodeNotImplemented uint8 = 4
    RcodeRefused        uint8 = 5
)

// DNS flags
const (
    FlagQR uint16 = 1 << 15 // Query/Response
    FlagAA uint16 = 1 << 10 // Authoritative Answer
    FlagTC uint16 = 1 << 9  // Truncated
    FlagRD uint16 = 1 << 8  // Recursion Desired
    FlagRA uint16 = 1 << 7  // Recursion Available
)

// Header represents a DNS message header
type Header struct {
    ID      uint16
    Flags   uint16
    QDCount uint16 // Question count
    ANCount uint16 // Answer count
    NSCount uint16 // Authority count
    ARCount uint16 // Additional count
}

// Question represents a DNS question
type Question struct {
    Name  string
    Type  uint16
    Class uint16
}

// ResourceRecord represents a DNS resource record
type ResourceRecord struct {
    Name     string
    Type     uint16
    Class    uint16
    TTL      uint32
    RDLength uint16
    RData    []byte
    
    // Parsed data (depending on type)
    Address  net.IP   // For A, AAAA
    Target   string   // For CNAME, NS, MX
    Priority uint16   // For MX
    Text     []string // For TXT
    SOAData  *SOA     // For SOA
}

// SOA represents Start of Authority data
type SOA struct {
    MName   string // Primary nameserver
    RName   string // Admin email (@ replaced with .)
    Serial  uint32
    Refresh uint32
    Retry   uint32
    Expire  uint32
    Minimum uint32
}

// Message represents a complete DNS message
type Message struct {
    Header     Header
    Questions  []Question
    Answers    []ResourceRecord
    Authority  []ResourceRecord
    Additional []ResourceRecord
}

// TypeToString converts record type to string
func TypeToString(t uint16) string {
    switch t {
    case TypeA:
        return "A"
    case TypeAAAA:
        return "AAAA"
    case TypeCNAME:
        return "CNAME"
    case TypeMX:
        return "MX"
    case TypeNS:
        return "NS"
    case TypeTXT:
        return "TXT"
    case TypeSOA:
        return "SOA"
    default:
        return fmt.Sprintf("TYPE%d", t)
    }
}

// StringToType converts string to record type
func StringToType(s string) uint16 {
    switch s {
    case "A":
        return TypeA
    case "AAAA":
        return TypeAAAA
    case "CNAME":
        return TypeCNAME
    case "MX":
        return TypeMX
    case "NS":
        return TypeNS
    case "TXT":
        return TypeTXT
    case "SOA":
        return TypeSOA
    default:
        return 0
    }
}
```

### 22.2 DNS Parser

```go
// File: dns/parser.go
package dns

import (
    "encoding/binary"
    "fmt"
    "net"
    "strings"
)

// Parser handles DNS message parsing
type Parser struct {
    data []byte
    pos  int
}

// NewParser creates a new DNS parser
func NewParser(data []byte) *Parser {
    return &Parser{data: data, pos: 0}
}

// Parse parses a complete DNS message
func (p *Parser) Parse() (*Message, error) {
    msg := &Message{}
    
    // Parse header
    if err := p.parseHeader(&msg.Header); err != nil {
        return nil, fmt.Errorf("header: %w", err)
    }
    
    // Parse questions
    msg.Questions = make([]Question, msg.Header.QDCount)
    for i := 0; i < int(msg.Header.QDCount); i++ {
        if err := p.parseQuestion(&msg.Questions[i]); err != nil {
            return nil, fmt.Errorf("question %d: %w", i, err)
        }
    }
    
    // Parse answers
    msg.Answers = make([]ResourceRecord, msg.Header.ANCount)
    for i := 0; i < int(msg.Header.ANCount); i++ {
        if err := p.parseResourceRecord(&msg.Answers[i]); err != nil {
            return nil, fmt.Errorf("answer %d: %w", i, err)
        }
    }
    
    // Parse authority
    msg.Authority = make([]ResourceRecord, msg.Header.NSCount)
    for i := 0; i < int(msg.Header.NSCount); i++ {
        if err := p.parseResourceRecord(&msg.Authority[i]); err != nil {
            return nil, fmt.Errorf("authority %d: %w", i, err)
        }
    }
    
    // Parse additional
    msg.Additional = make([]ResourceRecord, msg.Header.ARCount)
    for i := 0; i < int(msg.Header.ARCount); i++ {
        if err := p.parseResourceRecord(&msg.Additional[i]); err != nil {
            return nil, fmt.Errorf("additional %d: %w", i, err)
        }
    }
    
    return msg, nil
}

func (p *Parser) parseHeader(h *Header) error {
    if len(p.data) < 12 {
        return fmt.Errorf("header too short")
    }
    
    h.ID = binary.BigEndian.Uint16(p.data[0:2])
    h.Flags = binary.BigEndian.Uint16(p.data[2:4])
    h.QDCount = binary.BigEndian.Uint16(p.data[4:6])
    h.ANCount = binary.BigEndian.Uint16(p.data[6:8])
    h.NSCount = binary.BigEndian.Uint16(p.data[8:10])
    h.ARCount = binary.BigEndian.Uint16(p.data[10:12])
    
    p.pos = 12
    return nil
}

func (p *Parser) parseQuestion(q *Question) error {
    name, err := p.parseName()
    if err != nil {
        return err
    }
    q.Name = name
    
    if p.pos+4 > len(p.data) {
        return fmt.Errorf("question too short")
    }
    
    q.Type = binary.BigEndian.Uint16(p.data[p.pos : p.pos+2])
    q.Class = binary.BigEndian.Uint16(p.data[p.pos+2 : p.pos+4])
    p.pos += 4
    
    return nil
}

func (p *Parser) parseResourceRecord(rr *ResourceRecord) error {
    name, err := p.parseName()
    if err != nil {
        return err
    }
    rr.Name = name
    
    if p.pos+10 > len(p.data) {
        return fmt.Errorf("resource record too short")
    }
    
    rr.Type = binary.BigEndian.Uint16(p.data[p.pos : p.pos+2])
    rr.Class = binary.BigEndian.Uint16(p.data[p.pos+2 : p.pos+4])
    rr.TTL = binary.BigEndian.Uint32(p.data[p.pos+4 : p.pos+8])
    rr.RDLength = binary.BigEndian.Uint16(p.data[p.pos+8 : p.pos+10])
    p.pos += 10
    
    if p.pos+int(rr.RDLength) > len(p.data) {
        return fmt.Errorf("rdata too short")
    }
    
    rr.RData = p.data[p.pos : p.pos+int(rr.RDLength)]
    
    // Parse type-specific data
    switch rr.Type {
    case TypeA:
        if rr.RDLength == 4 {
            rr.Address = net.IP(rr.RData)
        }
    case TypeAAAA:
        if rr.RDLength == 16 {
            rr.Address = net.IP(rr.RData)
        }
    case TypeCNAME, TypeNS:
        savedPos := p.pos
        rr.Target, _ = p.parseName()
        p.pos = savedPos
    case TypeMX:
        if rr.RDLength >= 2 {
            rr.Priority = binary.BigEndian.Uint16(rr.RData[0:2])
            savedPos := p.pos
            p.pos = savedPos + 2  // Skip priority, parse name
            rr.Target, _ = p.parseName()
            p.pos = savedPos
        }
    case TypeTXT:
        rr.Text = p.parseTXT(rr.RData)
    }
    
    p.pos += int(rr.RDLength)
    return nil
}

// parseName handles DNS name compression
func (p *Parser) parseName() (string, error) {
    var labels []string
    visited := make(map[int]bool)
    
    for {
        if p.pos >= len(p.data) {
            return "", fmt.Errorf("name extends past data")
        }
        
        length := int(p.data[p.pos])
        
        // Check for compression pointer (top 2 bits set)
        if length&0xC0 == 0xC0 {
            if p.pos+1 >= len(p.data) {
                return "", fmt.Errorf("invalid compression pointer")
            }
            
            // Get offset
            offset := int(binary.BigEndian.Uint16(p.data[p.pos:p.pos+2]) & 0x3FFF)
            p.pos += 2
            
            // Prevent infinite loops
            if visited[offset] {
                return "", fmt.Errorf("compression loop detected")
            }
            visited[offset] = true
            
            // Save position, jump to offset, parse, restore
            savedPos := p.pos
            p.pos = offset
            rest, err := p.parseName()
            p.pos = savedPos
            if err != nil {
                return "", err
            }
            
            if len(labels) > 0 {
                return strings.Join(labels, ".") + "." + rest, nil
            }
            return rest, nil
        }
        
        // End of name
        if length == 0 {
            p.pos++
            break
        }
        
        // Regular label
        p.pos++
        if p.pos+length > len(p.data) {
            return "", fmt.Errorf("label extends past data")
        }
        
        labels = append(labels, string(p.data[p.pos:p.pos+length]))
        p.pos += length
    }
    
    return strings.Join(labels, "."), nil
}

func (p *Parser) parseTXT(data []byte) []string {
    var texts []string
    pos := 0
    
    for pos < len(data) {
        length := int(data[pos])
        pos++
        
        if pos+length > len(data) {
            break
        }
        
        texts = append(texts, string(data[pos:pos+length]))
        pos += length
    }
    
    return texts
}
```

---

## 23. Building DNS Responses

### 23.1 DNS Message Builder

```go
// File: dns/builder.go
package dns

import (
    "encoding/binary"
    "net"
    "strings"
)

// Builder constructs DNS messages
type Builder struct {
    data []byte
}

// NewBuilder creates a new DNS message builder
func NewBuilder() *Builder {
    return &Builder{
        data: make([]byte, 0, 512),
    }
}

// BuildResponse builds a response message for a query
func (b *Builder) BuildResponse(query *Message, answers []ResourceRecord, authority []ResourceRecord) []byte {
    b.data = b.data[:0]
    
    // Header
    header := Header{
        ID:      query.Header.ID,
        Flags:   FlagQR | FlagAA, // Response + Authoritative
        QDCount: uint16(len(query.Questions)),
        ANCount: uint16(len(answers)),
        NSCount: uint16(len(authority)),
        ARCount: 0,
    }
    
    // Set recursion available if requested
    if query.Header.Flags&FlagRD != 0 {
        header.Flags |= FlagRD
    }
    
    b.writeHeader(&header)
    
    // Questions (echo back)
    for _, q := range query.Questions {
        b.writeQuestion(&q)
    }
    
    // Answers
    for _, rr := range answers {
        b.writeResourceRecord(&rr)
    }
    
    // Authority
    for _, rr := range authority {
        b.writeResourceRecord(&rr)
    }
    
    return b.data
}

// BuildErrorResponse builds an error response
func (b *Builder) BuildErrorResponse(query *Message, rcode uint8) []byte {
    b.data = b.data[:0]
    
    header := Header{
        ID:      query.Header.ID,
        Flags:   FlagQR | FlagAA | uint16(rcode),
        QDCount: uint16(len(query.Questions)),
        ANCount: 0,
        NSCount: 0,
        ARCount: 0,
    }
    
    b.writeHeader(&header)
    
    for _, q := range query.Questions {
        b.writeQuestion(&q)
    }
    
    return b.data
}

func (b *Builder) writeHeader(h *Header) {
    b.writeUint16(h.ID)
    b.writeUint16(h.Flags)
    b.writeUint16(h.QDCount)
    b.writeUint16(h.ANCount)
    b.writeUint16(h.NSCount)
    b.writeUint16(h.ARCount)
}

func (b *Builder) writeQuestion(q *Question) {
    b.writeName(q.Name)
    b.writeUint16(q.Type)
    b.writeUint16(q.Class)
}

func (b *Builder) writeResourceRecord(rr *ResourceRecord) {
    b.writeName(rr.Name)
    b.writeUint16(rr.Type)
    b.writeUint16(rr.Class)
    b.writeUint32(rr.TTL)
    
    // Build RDATA based on type
    rdata := b.buildRData(rr)
    b.writeUint16(uint16(len(rdata)))
    b.data = append(b.data, rdata...)
}

func (b *Builder) buildRData(rr *ResourceRecord) []byte {
    switch rr.Type {
    case TypeA:
        return rr.Address.To4()
    case TypeAAAA:
        return rr.Address.To16()
    case TypeCNAME, TypeNS:
        return b.encodeName(rr.Target)
    case TypeMX:
        data := make([]byte, 2)
        binary.BigEndian.PutUint16(data, rr.Priority)
        data = append(data, b.encodeName(rr.Target)...)
        return data
    case TypeTXT:
        return b.encodeTXT(rr.Text)
    case TypeSOA:
        if rr.SOAData != nil {
            return b.encodeSOA(rr.SOAData)
        }
    }
    return rr.RData
}

func (b *Builder) writeName(name string) {
    b.data = append(b.data, b.encodeName(name)...)
}

func (b *Builder) encodeName(name string) []byte {
    var result []byte
    
    if name == "" || name == "." {
        return []byte{0}
    }
    
    // Remove trailing dot
    name = strings.TrimSuffix(name, ".")
    
    labels := strings.Split(name, ".")
    for _, label := range labels {
        if len(label) > 63 {
            label = label[:63]
        }
        result = append(result, byte(len(label)))
        result = append(result, []byte(label)...)
    }
    result = append(result, 0)
    
    return result
}

func (b *Builder) encodeTXT(texts []string) []byte {
    var result []byte
    for _, text := range texts {
        if len(text) > 255 {
            text = text[:255]
        }
        result = append(result, byte(len(text)))
        result = append(result, []byte(text)...)
    }
    return result
}

func (b *Builder) encodeSOA(soa *SOA) []byte {
    var result []byte
    result = append(result, b.encodeName(soa.MName)...)
    result = append(result, b.encodeName(soa.RName)...)
    
    nums := make([]byte, 20)
    binary.BigEndian.PutUint32(nums[0:4], soa.Serial)
    binary.BigEndian.PutUint32(nums[4:8], soa.Refresh)
    binary.BigEndian.PutUint32(nums[8:12], soa.Retry)
    binary.BigEndian.PutUint32(nums[12:16], soa.Expire)
    binary.BigEndian.PutUint32(nums[16:20], soa.Minimum)
    result = append(result, nums...)
    
    return result
}

func (b *Builder) writeUint16(v uint16) {
    bytes := make([]byte, 2)
    binary.BigEndian.PutUint16(bytes, v)
    b.data = append(b.data, bytes...)
}

func (b *Builder) writeUint32(v uint32) {
    bytes := make([]byte, 4)
    binary.BigEndian.PutUint32(bytes, v)
    b.data = append(b.data, bytes...)
}

// Helper functions to create resource records

// NewARecord creates an A record
func NewARecord(name string, ttl uint32, ip net.IP) ResourceRecord {
    return ResourceRecord{
        Name:    name,
        Type:    TypeA,
        Class:   ClassIN,
        TTL:     ttl,
        Address: ip.To4(),
    }
}

// NewAAAARecord creates an AAAA record
func NewAAAARecord(name string, ttl uint32, ip net.IP) ResourceRecord {
    return ResourceRecord{
        Name:    name,
        Type:    TypeAAAA,
        Class:   ClassIN,
        TTL:     ttl,
        Address: ip.To16(),
    }
}

// NewCNAMERecord creates a CNAME record
func NewCNAMERecord(name string, ttl uint32, target string) ResourceRecord {
    return ResourceRecord{
        Name:   name,
        Type:   TypeCNAME,
        Class:  ClassIN,
        TTL:    ttl,
        Target: target,
    }
}

// NewMXRecord creates an MX record
func NewMXRecord(name string, ttl uint32, priority uint16, target string) ResourceRecord {
    return ResourceRecord{
        Name:     name,
        Type:     TypeMX,
        Class:    ClassIN,
        TTL:      ttl,
        Priority: priority,
        Target:   target,
    }
}

// NewTXTRecord creates a TXT record
func NewTXTRecord(name string, ttl uint32, texts ...string) ResourceRecord {
    return ResourceRecord{
        Name:  name,
        Type:  TypeTXT,
        Class: ClassIN,
        TTL:   ttl,
        Text:  texts,
    }
}

// NewNSRecord creates an NS record
func NewNSRecord(name string, ttl uint32, target string) ResourceRecord {
    return ResourceRecord{
        Name:   name,
        Type:   TypeNS,
        Class:  ClassIN,
        TTL:    ttl,
        Target: target,
    }
}

// NewSOARecord creates an SOA record
func NewSOARecord(name string, ttl uint32, soa *SOA) ResourceRecord {
    return ResourceRecord{
        Name:    name,
        Type:    TypeSOA,
        Class:   ClassIN,
        TTL:     ttl,
        SOAData: soa,
    }
}
```

---

## 24. Zone File Parsing

### 24.1 Zone Data Structure

```go
// File: dns/zone.go
package dns

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "sync"
)

// Zone represents a DNS zone
type Zone struct {
    Name    string
    Records map[string][]ResourceRecord // Keyed by name+type
    SOA     *SOA
    mu      sync.RWMutex
}

// NewZone creates a new zone
func NewZone(name string) *Zone {
    return &Zone{
        Name:    strings.ToLower(name),
        Records: make(map[string][]ResourceRecord),
    }
}

// AddRecord adds a record to the zone
func (z *Zone) AddRecord(rr ResourceRecord) {
    z.mu.Lock()
    defer z.mu.Unlock()
    
    key := z.recordKey(rr.Name, rr.Type)
    z.Records[key] = append(z.Records[key], rr)
    
    if rr.Type == TypeSOA && rr.SOAData != nil {
        z.SOA = rr.SOAData
    }
}

// Lookup finds records matching name and type
func (z *Zone) Lookup(name string, qtype uint16) []ResourceRecord {
    z.mu.RLock()
    defer z.mu.RUnlock()
    
    name = strings.ToLower(name)
    
    // Direct match
    key := z.recordKey(name, qtype)
    if records, ok := z.Records[key]; ok {
        return records
    }
    
    // If looking for A/AAAA, check for CNAME
    if qtype == TypeA || qtype == TypeAAAA {
        cnameKey := z.recordKey(name, TypeCNAME)
        if cnames, ok := z.Records[cnameKey]; ok {
            return cnames
        }
    }
    
    return nil
}

// HasName checks if zone has any records for name
func (z *Zone) HasName(name string) bool {
    z.mu.RLock()
    defer z.mu.RUnlock()
    
    name = strings.ToLower(name)
    
    for key := range z.Records {
        if strings.HasPrefix(key, name+":") {
            return true
        }
    }
    return false
}

// IsAuthoritative checks if this zone is authoritative for the name
func (z *Zone) IsAuthoritative(name string) bool {
    name = strings.ToLower(name)
    zoneName := strings.ToLower(z.Name)
    
    return name == zoneName || strings.HasSuffix(name, "."+zoneName)
}

func (z *Zone) recordKey(name string, qtype uint16) string {
    return strings.ToLower(name) + ":" + strconv.Itoa(int(qtype))
}

// LoadZoneFile loads a zone from BIND-style zone file
func LoadZoneFile(filename string) (*Zone, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var zone *Zone
    var origin string
    var defaultTTL uint32 = 3600
    var currentName string
    
    scanner := bufio.NewScanner(file)
    lineNum := 0
    
    for scanner.Scan() {
        lineNum++
        line := strings.TrimSpace(scanner.Text())
        
        // Skip empty lines and comments
        if line == "" || strings.HasPrefix(line, ";") {
            continue
        }
        
        // Handle directives
        if strings.HasPrefix(line, "$ORIGIN") {
            origin = strings.TrimSpace(strings.TrimPrefix(line, "$ORIGIN"))
            origin = strings.TrimSuffix(origin, ".")
            if zone == nil {
                zone = NewZone(origin)
            }
            continue
        }
        
        if strings.HasPrefix(line, "$TTL") {
            ttlStr := strings.TrimSpace(strings.TrimPrefix(line, "$TTL"))
            ttl, err := parseTTL(ttlStr)
            if err != nil {
                return nil, fmt.Errorf("line %d: invalid TTL: %v", lineNum, err)
            }
            defaultTTL = ttl
            continue
        }
        
        // Parse record
        rr, name, err := parseZoneLine(line, origin, currentName, defaultTTL)
        if err != nil {
            return nil, fmt.Errorf("line %d: %v", lineNum, err)
        }
        
        if name != "" {
            currentName = name
        }
        
        if zone == nil {
            zone = NewZone(origin)
        }
        
        zone.AddRecord(rr)
    }
    
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    
    return zone, nil
}

func parseZoneLine(line, origin, currentName string, defaultTTL uint32) (ResourceRecord, string, error) {
    fields := strings.Fields(line)
    if len(fields) < 3 {
        return ResourceRecord{}, "", fmt.Errorf("too few fields")
    }
    
    var rr ResourceRecord
    var name string
    idx := 0
    
    // First field: name, TTL, class, or type
    field := fields[idx]
    
    // Check if first field is a name
    if !isClassOrType(field) && !isTTL(field) {
        if field == "@" {
            name = origin
        } else if !strings.HasSuffix(field, ".") {
            name = field + "." + origin
        } else {
            name = strings.TrimSuffix(field, ".")
        }
        idx++
    } else {
        name = currentName
    }
    
    rr.Name = name
    rr.TTL = defaultTTL
    rr.Class = ClassIN
    
    // Parse optional TTL
    if idx < len(fields) && isTTL(fields[idx]) {
        ttl, _ := parseTTL(fields[idx])
        rr.TTL = ttl
        idx++
    }
    
    // Parse optional class
    if idx < len(fields) && isClass(fields[idx]) {
        idx++ // Skip class, assume IN
    }
    
    // Parse type
    if idx >= len(fields) {
        return rr, name, fmt.Errorf("missing type")
    }
    
    rr.Type = StringToType(strings.ToUpper(fields[idx]))
    if rr.Type == 0 {
        return rr, name, fmt.Errorf("unknown type: %s", fields[idx])
    }
    idx++
    
    // Parse RDATA
    if idx >= len(fields) {
        return rr, name, fmt.Errorf("missing rdata")
    }
    
    rdata := strings.Join(fields[idx:], " ")
    
    switch rr.Type {
    case TypeA:
        ip := net.ParseIP(rdata)
        if ip == nil || ip.To4() == nil {
            return rr, name, fmt.Errorf("invalid IPv4: %s", rdata)
        }
        rr.Address = ip.To4()
        
    case TypeAAAA:
        ip := net.ParseIP(rdata)
        if ip == nil || ip.To16() == nil || ip.To4() != nil {
            return rr, name, fmt.Errorf("invalid IPv6: %s", rdata)
        }
        rr.Address = ip.To16()
        
    case TypeCNAME, TypeNS:
        target := fields[idx]
        if target == "@" {
            target = origin
        } else if !strings.HasSuffix(target, ".") {
            target = target + "." + origin
        } else {
            target = strings.TrimSuffix(target, ".")
        }
        rr.Target = target
        
    case TypeMX:
        if idx+1 >= len(fields) {
            return rr, name, fmt.Errorf("MX needs priority and target")
        }
        priority, err := strconv.ParseUint(fields[idx], 10, 16)
        if err != nil {
            return rr, name, fmt.Errorf("invalid MX priority: %v", err)
        }
        rr.Priority = uint16(priority)
        
        target := fields[idx+1]
        if !strings.HasSuffix(target, ".") {
            target = target + "." + origin
        } else {
            target = strings.TrimSuffix(target, ".")
        }
        rr.Target = target
        
    case TypeTXT:
        // Handle quoted strings
        text := strings.Trim(rdata, "\"")
        rr.Text = []string{text}
        
    case TypeSOA:
        if len(fields) < idx+7 {
            return rr, name, fmt.Errorf("SOA needs 7 fields")
        }
        soa := &SOA{}
        soa.MName = normalizeSOAName(fields[idx], origin)
        soa.RName = normalizeSOAName(fields[idx+1], origin)
        
        soa.Serial, _ = parseUint32(fields[idx+2])
        soa.Refresh, _ = parseTTL(fields[idx+3])
        soa.Retry, _ = parseTTL(fields[idx+4])
        soa.Expire, _ = parseTTL(fields[idx+5])
        soa.Minimum, _ = parseTTL(fields[idx+6])
        
        rr.SOAData = soa
    }
    
    return rr, name, nil
}

func normalizeSOAName(name, origin string) string {
    if name == "@" {
        return origin
    }
    if !strings.HasSuffix(name, ".") {
        return name + "." + origin
    }
    return strings.TrimSuffix(name, ".")
}

func isClassOrType(s string) bool {
    return isClass(s) || StringToType(strings.ToUpper(s)) != 0
}

func isClass(s string) bool {
    return strings.ToUpper(s) == "IN" || strings.ToUpper(s) == "CH"
}

func isTTL(s string) bool {
    _, err := parseTTL(s)
    return err == nil
}

func parseTTL(s string) (uint32, error) {
    // Handle suffixes: 1h, 1d, 1w
    s = strings.ToLower(s)
    multiplier := uint32(1)
    
    if strings.HasSuffix(s, "w") {
        multiplier = 604800
        s = s[:len(s)-1]
    } else if strings.HasSuffix(s, "d") {
        multiplier = 86400
        s = s[:len(s)-1]
    } else if strings.HasSuffix(s, "h") {
        multiplier = 3600
        s = s[:len(s)-1]
    } else if strings.HasSuffix(s, "m") {
        multiplier = 60
        s = s[:len(s)-1]
    } else if strings.HasSuffix(s, "s") {
        s = s[:len(s)-1]
    }
    
    val, err := strconv.ParseUint(s, 10, 32)
    if err != nil {
        return 0, err
    }
    
    return uint32(val) * multiplier, nil
}

func parseUint32(s string) (uint32, error) {
    val, err := strconv.ParseUint(s, 10, 32)
    return uint32(val), err
}
```

---

## 25. Complete DNS Server

### 25.1 Main Server

```go
// File: cmd/dns-server/main.go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "sync"
    "syscall"
    
    "github.com/bellistech/dns-server/dns"
)

type Server struct {
    zones   map[string]*dns.Zone
    mu      sync.RWMutex
    builder *dns.Builder
    
    udpConn4 *net.UDPConn
    udpConn6 *net.UDPConn
    
    stats struct {
        queries  uint64
        answers  uint64
        nxdomain uint64
        errors   uint64
    }
}

func NewServer() *Server {
    return &Server{
        zones:   make(map[string]*dns.Zone),
        builder: dns.NewBuilder(),
    }
}

func (s *Server) LoadZone(filename string) error {
    zone, err := dns.LoadZoneFile(filename)
    if err != nil {
        return fmt.Errorf("loading %s: %w", filename, err)
    }
    
    s.mu.Lock()
    s.zones[zone.Name] = zone
    s.mu.Unlock()
    
    log.Printf("Loaded zone: %s", zone.Name)
    return nil
}

func (s *Server) Start(ctx context.Context, addr4, addr6 string) error {
    var wg sync.WaitGroup
    
    // Start IPv4 listener
    if addr4 != "" {
        udpAddr4, err := net.ResolveUDPAddr("udp4", addr4)
        if err != nil {
            return fmt.Errorf("resolve IPv4: %w", err)
        }
        
        s.udpConn4, err = net.ListenUDP("udp4", udpAddr4)
        if err != nil {
            return fmt.Errorf("listen IPv4: %w", err)
        }
        
        log.Printf("Listening on IPv4 %s", addr4)
        
        wg.Add(1)
        go func() {
            defer wg.Done()
            s.serveUDP(ctx, s.udpConn4)
        }()
    }
    
    // Start IPv6 listener
    if addr6 != "" {
        udpAddr6, err := net.ResolveUDPAddr("udp6", addr6)
        if err != nil {
            return fmt.Errorf("resolve IPv6: %w", err)
        }
        
        s.udpConn6, err = net.ListenUDP("udp6", udpAddr6)
        if err != nil {
            return fmt.Errorf("listen IPv6: %w", err)
        }
        
        log.Printf("Listening on IPv6 %s", addr6)
        
        wg.Add(1)
        go func() {
            defer wg.Done()
            s.serveUDP(ctx, s.udpConn6)
        }()
    }
    
    wg.Wait()
    return nil
}

func (s *Server) serveUDP(ctx context.Context, conn *net.UDPConn) {
    buffer := make([]byte, 512)
    
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }
        
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            select {
            case <-ctx.Done():
                return
            default:
                log.Printf("Read error: %v", err)
                continue
            }
        }
        
        // Handle in goroutine for concurrency
        go s.handleQuery(conn, clientAddr, buffer[:n])
    }
}

func (s *Server) handleQuery(conn *net.UDPConn, clientAddr *net.UDPAddr, data []byte) {
    s.stats.queries++
    
    // Parse query
    parser := dns.NewParser(data)
    query, err := parser.Parse()
    if err != nil {
        log.Printf("Parse error from %s: %v", clientAddr, err)
        s.stats.errors++
        return
    }
    
    if len(query.Questions) == 0 {
        return
    }
    
    q := query.Questions[0]
    log.Printf("Query from %s: %s %s", clientAddr, q.Name, dns.TypeToString(q.Type))
    
    // Find zone
    zone := s.findZone(q.Name)
    if zone == nil {
        // Not authoritative
        response := s.builder.BuildErrorResponse(query, dns.RcodeRefused)
        conn.WriteToUDP(response, clientAddr)
        return
    }
    
    // Lookup records
    records := zone.Lookup(q.Name, q.Type)
    
    if len(records) == 0 && !zone.HasName(q.Name) {
        // NXDOMAIN
        s.stats.nxdomain++
        response := s.builder.BuildErrorResponse(query, dns.RcodeNameError)
        conn.WriteToUDP(response, clientAddr)
        return
    }
    
    // Build response
    s.stats.answers++
    
    // Get NS records for authority section
    nsRecords := zone.Lookup(zone.Name, dns.TypeNS)
    
    response := s.builder.BuildResponse(query, records, nsRecords)
    conn.WriteToUDP(response, clientAddr)
}

func (s *Server) findZone(name string) *dns.Zone {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // Find most specific zone
    labels := splitLabels(name)
    
    for i := 0; i < len(labels); i++ {
        zoneName := joinLabels(labels[i:])
        if zone, ok := s.zones[zoneName]; ok {
            return zone
        }
    }
    
    return nil
}

func splitLabels(name string) []string {
    name = strings.TrimSuffix(name, ".")
    if name == "" {
        return nil
    }
    return strings.Split(name, ".")
}

func joinLabels(labels []string) string {
    return strings.Join(labels, ".")
}

func (s *Server) Stop() {
    if s.udpConn4 != nil {
        s.udpConn4.Close()
    }
    if s.udpConn6 != nil {
        s.udpConn6.Close()
    }
    
    log.Printf("Statistics: queries=%d, answers=%d, nxdomain=%d, errors=%d",
        s.stats.queries, s.stats.answers, s.stats.nxdomain, s.stats.errors)
}

import "strings"

func main() {
    addr4 := flag.String("4", ":5353", "IPv4 listen address")
    addr6 := flag.String("6", "[::]:5353", "IPv6 listen address")
    zoneFile := flag.String("zone", "", "Zone file to load")
    flag.Parse()
    
    if *zoneFile == "" {
        log.Fatal("Zone file required (-zone)")
    }
    
    server := NewServer()
    
    if err := server.LoadZone(*zoneFile); err != nil {
        log.Fatalf("Failed to load zone: %v", err)
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    
    // Handle shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        log.Println("Shutting down...")
        cancel()
        server.Stop()
    }()
    
    if err := server.Start(ctx, *addr4, *addr6); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

### 25.2 Example Zone File

```
; File: zones/example.com.zone
$ORIGIN example.com.
$TTL 3600

; SOA Record
@   IN  SOA ns1.example.com. admin.example.com. (
            2024010101  ; Serial
            3600        ; Refresh
            1800        ; Retry
            604800      ; Expire
            86400       ; Minimum TTL
        )

; Name Servers
@       IN  NS  ns1.example.com.
@       IN  NS  ns2.example.com.

; A Records (IPv4)
@       IN  A       93.184.216.34
www     IN  A       93.184.216.34
ns1     IN  A       192.0.2.1
ns2     IN  A       192.0.2.2
mail    IN  A       192.0.2.10

; AAAA Records (IPv6)
@       IN  AAAA    2606:2800:220:1:248:1893:25c8:1946
www     IN  AAAA    2606:2800:220:1:248:1893:25c8:1946
ns1     IN  AAAA    2001:db8::1
ns2     IN  AAAA    2001:db8::2

; CNAME Records
ftp     IN  CNAME   www.example.com.
blog    IN  CNAME   www.example.com.

; MX Records
@       IN  MX  10  mail.example.com.
@       IN  MX  20  mail2.example.com.

; TXT Records
@       IN  TXT     "v=spf1 mx -all"
_dmarc  IN  TXT     "v=DMARC1; p=reject"
```

---

## 26. Testing & Deployment

### 26.1 Testing the Server

```bash
# Start the server
go run cmd/dns-server/main.go -zone zones/example.com.zone

# Test with dig (IPv4)
dig @localhost -p 5353 example.com A
dig @localhost -p 5353 www.example.com A
dig @localhost -p 5353 example.com AAAA
dig @localhost -p 5353 example.com MX
dig @localhost -p 5353 example.com TXT
dig @localhost -p 5353 example.com NS
dig @localhost -p 5353 example.com SOA

# Test CNAME
dig @localhost -p 5353 ftp.example.com A

# Test NXDOMAIN
dig @localhost -p 5353 nonexistent.example.com A

# Test IPv6 (if available)
dig @::1 -p 5353 example.com AAAA
```

### 26.2 Unit Tests

```go
// File: dns/parser_test.go
package dns

import (
    "net"
    "testing"
)

func TestParseQuery(t *testing.T) {
    // DNS query for example.com A record
    query := []byte{
        0x12, 0x34, // ID
        0x01, 0x00, // Flags (standard query)
        0x00, 0x01, // Questions: 1
        0x00, 0x00, // Answers: 0
        0x00, 0x00, // Authority: 0
        0x00, 0x00, // Additional: 0
        // Question: example.com A IN
        0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
        0x03, 'c', 'o', 'm',
        0x00,       // End of name
        0x00, 0x01, // Type A
        0x00, 0x01, // Class IN
    }
    
    parser := NewParser(query)
    msg, err := parser.Parse()
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }
    
    if msg.Header.ID != 0x1234 {
        t.Errorf("ID = %x, want 0x1234", msg.Header.ID)
    }
    
    if len(msg.Questions) != 1 {
        t.Fatalf("Questions = %d, want 1", len(msg.Questions))
    }
    
    q := msg.Questions[0]
    if q.Name != "example.com" {
        t.Errorf("Name = %s, want example.com", q.Name)
    }
    if q.Type != TypeA {
        t.Errorf("Type = %d, want %d", q.Type, TypeA)
    }
}

func TestBuildResponse(t *testing.T) {
    query := &Message{
        Header: Header{ID: 0x1234, QDCount: 1},
        Questions: []Question{
            {Name: "example.com", Type: TypeA, Class: ClassIN},
        },
    }
    
    answers := []ResourceRecord{
        NewARecord("example.com", 3600, net.ParseIP("93.184.216.34")),
    }
    
    builder := NewBuilder()
    response := builder.BuildResponse(query, answers, nil)
    
    // Parse response
    parser := NewParser(response)
    msg, err := parser.Parse()
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }
    
    if msg.Header.ID != 0x1234 {
        t.Errorf("Response ID = %x, want 0x1234", msg.Header.ID)
    }
    
    if msg.Header.Flags&FlagQR == 0 {
        t.Error("QR flag not set")
    }
    
    if len(msg.Answers) != 1 {
        t.Fatalf("Answers = %d, want 1", len(msg.Answers))
    }
    
    a := msg.Answers[0]
    expected := net.ParseIP("93.184.216.34").To4()
    if !a.Address.Equal(expected) {
        t.Errorf("Address = %v, want %v", a.Address, expected)
    }
}
```

### 26.3 Project Structure

```
dns-server/
├── go.mod
├── go.sum
├── README.md
├── Makefile
├── cmd/
│   └── dns-server/
│       └── main.go
├── dns/
│   ├── types.go
│   ├── parser.go
│   ├── parser_test.go
│   ├── builder.go
│   ├── builder_test.go
│   ├── zone.go
│   └── zone_test.go
├── zones/
│   └── example.com.zone
└── configs/
    └── server.yaml
```

### 26.4 Makefile

```makefile
.PHONY: build test run clean

BINARY=dns-server
ZONE=zones/example.com.zone

build:
	go build -o bin/$(BINARY) ./cmd/dns-server

test:
	go test -v ./...

run: build
	./bin/$(BINARY) -zone $(ZONE)

clean:
	rm -rf bin/

# Run with specific addresses
run-ipv4: build
	./bin/$(BINARY) -4 :5353 -6 "" -zone $(ZONE)

run-ipv6: build
	./bin/$(BINARY) -4 "" -6 [::]:5353 -zone $(ZONE)

# Integration test with dig
test-dig: run &
	sleep 1
	dig @localhost -p 5353 example.com A
	dig @localhost -p 5353 example.com AAAA
	dig @localhost -p 5353 example.com MX
```

---

## Summary

Congratulations! You've built a complete IPv4/IPv6 dual-stack authoritative DNS server in Go. This project demonstrates:

1. **Binary protocol parsing** - Reading and writing DNS wire format
2. **UDP networking** - Handling connectionless protocols
3. **Concurrent programming** - Goroutines for parallel query handling
4. **IPv6 support** - Dual-stack networking
5. **Configuration** - Zone file parsing
6. **Production patterns** - Graceful shutdown, logging, statistics

### Key Concepts Learned

| Concept | Application |
|---------|-------------|
| Binary encoding | DNS message format |
| Network byte order | Big-endian integers |
| UDP sockets | Connectionless queries |
| Goroutines | Concurrent query handling |
| Mutexes | Thread-safe zone access |
| Context | Graceful shutdown |
| Interfaces | Extensible record types |
| Testing | Table-driven tests |

### Next Steps

1. Add TCP support (for large responses)
2. Implement EDNS0 (extended DNS)
3. Add DNSSEC signing
4. Implement zone transfers (AXFR)
5. Add caching/forwarding
6. Build a CLI management tool
7. Add Prometheus metrics
8. Containerize with Docker

This DNS server provides a solid foundation for understanding network protocols and systems programming in Go.
