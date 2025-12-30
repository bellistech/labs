# Part III: Concurrency (Continued)

## 15. Context (Continued)

### 15.2 Context Propagation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    // Add request ID
    ctx = context.WithValue(ctx, "requestID", "req-12345")
    
    // Start the request chain
    result, err := handleRequest(ctx)
    if err != nil {
        fmt.Printf("Request failed: %v\n", err)
        return
    }
    fmt.Printf("Result: %s\n", result)
}

func handleRequest(ctx context.Context) (string, error) {
    requestID := ctx.Value("requestID").(string)
    fmt.Printf("[%s] Handling request\n", requestID)
    
    // Fetch data (propagate context)
    data, err := fetchData(ctx)
    if err != nil {
        return "", fmt.Errorf("fetch failed: %w", err)
    }
    
    // Process data
    result, err := processData(ctx, data)
    if err != nil {
        return "", fmt.Errorf("process failed: %w", err)
    }
    
    return result, nil
}

func fetchData(ctx context.Context) (string, error) {
    requestID := ctx.Value("requestID").(string)
    fmt.Printf("[%s] Fetching data\n", requestID)
    
    // Simulate slow fetch
    select {
    case <-time.After(200 * time.Millisecond):
        return "raw data", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func processData(ctx context.Context, data string) (string, error) {
    requestID := ctx.Value("requestID").(string)
    fmt.Printf("[%s] Processing: %s\n", requestID, data)
    
    // Simulate processing
    select {
    case <-time.After(100 * time.Millisecond):
        return "processed: " + data, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

### 15.3 Context Best Practices

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

// Rule 1: Context should be the first parameter
func GoodFunction(ctx context.Context, name string) error {
    return nil
}

// Rule 2: Don't store context in structs (usually)
type BadService struct {
    ctx context.Context  // Don't do this
}

type GoodService struct {
    // No context field
}

func (s *GoodService) DoWork(ctx context.Context) error {
    // Accept context as parameter
    return nil
}

// Rule 3: Use context values sparingly (request-scoped data only)
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
)

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}

// Rule 4: Check context.Done() in long operations
func LongOperation(ctx context.Context, items []int) error {
    for i, item := range items {
        // Check cancellation periodically
        select {
        case <-ctx.Done():
            return fmt.Errorf("cancelled at item %d: %w", i, ctx.Err())
        default:
        }
        
        // Process item
        processItem(item)
    }
    return nil
}

func processItem(item int) {
    time.Sleep(10 * time.Millisecond)
}

// Rule 5: Always call cancel function
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
    // Request already has context
    ctx := r.Context()
    
    // Add timeout for this specific operation
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()  // Always defer cancel!
    
    result, err := doExpensiveOperation(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "Result: %s", result)
}

func doExpensiveOperation(ctx context.Context) (string, error) {
    select {
    case <-time.After(1 * time.Second):
        return "done", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    // Example usage
    ctx := context.Background()
    ctx = WithRequestID(ctx, "req-abc123")
    
    fmt.Printf("Request ID: %s\n", GetRequestID(ctx))
    
    // With timeout
    ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()
    
    items := make([]int, 100)
    if err := LongOperation(ctx, items); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Completed successfully")
    }
}
```

---

# Part IV: Systems Programming

## 16. File I/O

### 16.1 Reading Files

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

func main() {
    // Read entire file at once (small files only)
    data, err := os.ReadFile("example.txt")
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
    } else {
        fmt.Printf("Content: %s\n", data)
    }
    
    // Open file for reading
    file, err := os.Open("example.txt")
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    
    // Read with buffer
    buf := make([]byte, 1024)
    for {
        n, err := file.Read(buf)
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
    }
    
    // Buffered reading (efficient for large files)
    file2, _ := os.Open("example.txt")
    defer file2.Close()
    
    reader := bufio.NewReader(file2)
    
    // Read line by line
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            if len(line) > 0 {
                fmt.Printf("Last line: %s\n", line)
            }
            break
        }
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            break
        }
        fmt.Printf("Line: %s", line)
    }
    
    // Scanner (convenient for line reading)
    file3, _ := os.Open("example.txt")
    defer file3.Close()
    
    scanner := bufio.NewScanner(file3)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        fmt.Printf("%d: %s\n", lineNum, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Printf("Scanner error: %v\n", err)
    }
}
```

### 16.2 Writing Files

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // Write entire file at once
    data := []byte("Hello, World!\n")
    err := os.WriteFile("output.txt", data, 0644)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // Create/truncate file
    file, err := os.Create("output2.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer file.Close()
    
    // Write string
    n, err := file.WriteString("Line 1\n")
    fmt.Printf("Wrote %d bytes\n", n)
    
    // Write bytes
    file.Write([]byte("Line 2\n"))
    
    // Buffered writing (efficient)
    writer := bufio.NewWriter(file)
    writer.WriteString("Buffered line 1\n")
    writer.WriteString("Buffered line 2\n")
    writer.Flush()  // Don't forget to flush!
    
    // Append to file
    appendFile, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer appendFile.Close()
    
    appendFile.WriteString("Appended line\n")
    
    // Open with multiple flags
    // os.O_RDONLY - read only
    // os.O_WRONLY - write only
    // os.O_RDWR   - read/write
    // os.O_CREATE - create if not exists
    // os.O_APPEND - append to file
    // os.O_TRUNC  - truncate file
    // os.O_EXCL   - error if file exists (with O_CREATE)
    
    f, err := os.OpenFile("data.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer f.Close()
    
    f.WriteString("Created with flags\n")
}
```

### 16.3 File Operations

```go
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func main() {
    // File info
    info, err := os.Stat("example.txt")
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("File does not exist")
        } else {
            fmt.Printf("Error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Name: %s\n", info.Name())
    fmt.Printf("Size: %d bytes\n", info.Size())
    fmt.Printf("Mode: %s\n", info.Mode())
    fmt.Printf("ModTime: %s\n", info.ModTime())
    fmt.Printf("IsDir: %t\n", info.IsDir())
    
    // Check if file exists
    if _, err := os.Stat("somefile.txt"); os.IsNotExist(err) {
        fmt.Println("File does not exist")
    }
    
    // Rename/move file
    err = os.Rename("old.txt", "new.txt")
    if err != nil {
        fmt.Printf("Rename error: %v\n", err)
    }
    
    // Copy file
    src, _ := os.Open("source.txt")
    defer src.Close()
    dst, _ := os.Create("destination.txt")
    defer dst.Close()
    
    bytesCopied, err := io.Copy(dst, src)
    fmt.Printf("Copied %d bytes\n", bytesCopied)
    
    // Delete file
    err = os.Remove("temp.txt")
    if err != nil {
        fmt.Printf("Remove error: %v\n", err)
    }
    
    // Create directory
    err = os.Mkdir("newdir", 0755)
    err = os.MkdirAll("path/to/nested/dir", 0755)  // Create all parents
    
    // Remove directory
    os.Remove("emptydir")           // Only if empty
    os.RemoveAll("path/to/nested")  // Remove all contents
    
    // List directory
    entries, _ := os.ReadDir(".")
    for _, entry := range entries {
        info, _ := entry.Info()
        fmt.Printf("%s\t%d\t%s\n", info.Mode(), info.Size(), entry.Name())
    }
    
    // Walk directory tree
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("%s (%d bytes)\n", path, info.Size())
        return nil
    })
    
    // Temporary file
    tmpFile, err := os.CreateTemp("", "prefix-*.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer os.Remove(tmpFile.Name())
    defer tmpFile.Close()
    
    fmt.Printf("Temp file: %s\n", tmpFile.Name())
    tmpFile.WriteString("temporary content")
    
    // Temporary directory
    tmpDir, _ := os.MkdirTemp("", "myapp-*")
    defer os.RemoveAll(tmpDir)
    fmt.Printf("Temp dir: %s\n", tmpDir)
}
```

---

## 17. Binary Data & Encoding

### 17.1 Binary Reading/Writing

```go
package main

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

func main() {
    // Write binary data
    buf := new(bytes.Buffer)
    
    // Write uint32 (big endian)
    var num uint32 = 0x12345678
    binary.Write(buf, binary.BigEndian, num)
    fmt.Printf("BigEndian: %x\n", buf.Bytes())  // [12 34 56 78]
    
    // Write uint32 (little endian)
    buf.Reset()
    binary.Write(buf, binary.LittleEndian, num)
    fmt.Printf("LittleEndian: %x\n", buf.Bytes())  // [78 56 34 12]
    
    // Write struct
    type Header struct {
        Magic   uint32
        Version uint16
        Flags   uint16
        Length  uint32
    }
    
    buf.Reset()
    header := Header{
        Magic:   0x44454641,  // "AFED" (little endian)
        Version: 1,
        Flags:   0x0003,
        Length:  1024,
    }
    binary.Write(buf, binary.BigEndian, header)
    fmt.Printf("Header bytes: %x\n", buf.Bytes())
    
    // Read binary data
    reader := bytes.NewReader(buf.Bytes())
    var readHeader Header
    binary.Read(reader, binary.BigEndian, &readHeader)
    fmt.Printf("Read header: %+v\n", readHeader)
    
    // Manual byte manipulation
    data := []byte{0x00, 0x00, 0x01, 0x00}  // 256 in big endian
    value := binary.BigEndian.Uint32(data)
    fmt.Printf("Value: %d\n", value)
    
    // Put bytes manually
    out := make([]byte, 4)
    binary.BigEndian.PutUint32(out, 65535)
    fmt.Printf("65535 as bytes: %x\n", out)
    
    // Working with bits
    var flags uint8 = 0
    flags |= (1 << 0)  // Set bit 0
    flags |= (1 << 2)  // Set bit 2
    fmt.Printf("Flags: %08b\n", flags)  // 00000101
    
    // Check bit
    if flags&(1<<2) != 0 {
        fmt.Println("Bit 2 is set")
    }
    
    // Clear bit
    flags &^= (1 << 0)  // Clear bit 0
    fmt.Printf("After clear: %08b\n", flags)  // 00000100
}
```

### 17.2 Encoding (JSON, XML, etc.)

```go
package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
)

type Person struct {
    Name    string   `json:"name" xml:"name"`
    Age     int      `json:"age" xml:"age"`
    Email   string   `json:"email,omitempty" xml:"email,omitempty"`
    Tags    []string `json:"tags" xml:"tag"`
    private string   // Not exported, won't be encoded
}

func main() {
    // JSON encoding
    person := Person{
        Name:  "Alice",
        Age:   30,
        Email: "alice@example.com",
        Tags:  []string{"developer", "gopher"},
    }
    
    // Marshal (struct to JSON)
    jsonData, err := json.Marshal(person)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("JSON: %s\n", jsonData)
    
    // Pretty print
    jsonPretty, _ := json.MarshalIndent(person, "", "  ")
    fmt.Printf("Pretty JSON:\n%s\n", jsonPretty)
    
    // Unmarshal (JSON to struct)
    jsonStr := `{"name":"Bob","age":25,"tags":["admin"]}`
    var person2 Person
    json.Unmarshal([]byte(jsonStr), &person2)
    fmt.Printf("Unmarshaled: %+v\n", person2)
    
    // Dynamic JSON (interface{})
    var dynamic map[string]interface{}
    json.Unmarshal([]byte(jsonStr), &dynamic)
    fmt.Printf("Dynamic: %v\n", dynamic)
    fmt.Printf("Name: %s\n", dynamic["name"].(string))
    
    // JSON streaming (encoder/decoder)
    // For large files or network streams
    
    // XML encoding
    xmlData, _ := xml.MarshalIndent(person, "", "  ")
    fmt.Printf("\nXML:\n%s\n", xmlData)
    
    // XML with root element
    type People struct {
        XMLName xml.Name `xml:"people"`
        Persons []Person `xml:"person"`
    }
    
    people := People{
        Persons: []Person{person, person2},
    }
    xmlData2, _ := xml.MarshalIndent(people, "", "  ")
    fmt.Printf("\nPeople XML:\n%s\n", xmlData2)
}
```

### 17.3 Base64 and Hex Encoding

```go
package main

import (
    "encoding/base64"
    "encoding/hex"
    "fmt"
)

func main() {
    original := []byte("Hello, World!")
    
    // Base64 Standard encoding
    encoded := base64.StdEncoding.EncodeToString(original)
    fmt.Printf("Base64: %s\n", encoded)
    
    decoded, _ := base64.StdEncoding.DecodeString(encoded)
    fmt.Printf("Decoded: %s\n", decoded)
    
    // Base64 URL-safe encoding
    urlSafe := base64.URLEncoding.EncodeToString(original)
    fmt.Printf("URL-safe Base64: %s\n", urlSafe)
    
    // Hex encoding
    hexStr := hex.EncodeToString(original)
    fmt.Printf("Hex: %s\n", hexStr)
    
    hexDecoded, _ := hex.DecodeString(hexStr)
    fmt.Printf("Hex decoded: %s\n", hexDecoded)
    
    // Dump (readable hex)
    dump := hex.Dump(original)
    fmt.Printf("Hex dump:\n%s", dump)
}
```

---

## 18. Network Programming Basics

### 18.1 DNS Lookups

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Lookup IP addresses
    ips, err := net.LookupIP("google.com")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("IPs for google.com:")
    for _, ip := range ips {
        if ip.To4() != nil {
            fmt.Printf("  IPv4: %s\n", ip)
        } else {
            fmt.Printf("  IPv6: %s\n", ip)
        }
    }
    
    // Lookup hostname from IP
    names, _ := net.LookupAddr("8.8.8.8")
    fmt.Printf("\nHostnames for 8.8.8.8: %v\n", names)
    
    // Lookup MX records
    mxRecords, _ := net.LookupMX("google.com")
    fmt.Println("\nMX records:")
    for _, mx := range mxRecords {
        fmt.Printf("  %s (priority %d)\n", mx.Host, mx.Pref)
    }
    
    // Lookup TXT records
    txtRecords, _ := net.LookupTXT("google.com")
    fmt.Println("\nTXT records:")
    for _, txt := range txtRecords {
        fmt.Printf("  %s\n", txt)
    }
    
    // Lookup NS records
    nsRecords, _ := net.LookupNS("google.com")
    fmt.Println("\nNS records:")
    for _, ns := range nsRecords {
        fmt.Printf("  %s\n", ns.Host)
    }
    
    // Lookup CNAME
    cname, _ := net.LookupCNAME("www.google.com")
    fmt.Printf("\nCNAME: %s\n", cname)
    
    // Lookup SRV records
    _, srvRecords, _ := net.LookupSRV("", "", "_http._tcp.google.com")
    fmt.Println("\nSRV records:")
    for _, srv := range srvRecords {
        fmt.Printf("  %s:%d (priority %d, weight %d)\n",
            srv.Target, srv.Port, srv.Priority, srv.Weight)
    }
}
```

### 18.2 IP Address Manipulation

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Parse IP address
    ip := net.ParseIP("192.168.1.1")
    fmt.Printf("Parsed IPv4: %s\n", ip)
    
    ip6 := net.ParseIP("2001:db8::1")
    fmt.Printf("Parsed IPv6: %s\n", ip6)
    
    // Check IP version
    if ip.To4() != nil {
        fmt.Println("Is IPv4")
    }
    if ip6.To4() == nil && ip6.To16() != nil {
        fmt.Println("Is IPv6")
    }
    
    // Parse CIDR notation
    ip, network, err := net.ParseCIDR("192.168.1.0/24")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("IP: %s, Network: %s\n", ip, network)
    fmt.Printf("Network IP: %s\n", network.IP)
    fmt.Printf("Mask: %s\n", network.Mask)
    
    // Check if IP is in network
    testIP := net.ParseIP("192.168.1.100")
    fmt.Printf("%s in %s: %t\n", testIP, network, network.Contains(testIP))
    
    testIP2 := net.ParseIP("192.168.2.1")
    fmt.Printf("%s in %s: %t\n", testIP2, network, network.Contains(testIP2))
    
    // IPv4 to bytes
    ipv4 := net.ParseIP("192.168.1.1").To4()
    fmt.Printf("IPv4 bytes: %v\n", []byte(ipv4))
    
    // IPv6 to bytes
    ipv6Bytes := net.ParseIP("::1").To16()
    fmt.Printf("IPv6 bytes: %v\n", []byte(ipv6Bytes))
    
    // Create IP from bytes
    newIP := net.IPv4(10, 0, 0, 1)
    fmt.Printf("Created IPv4: %s\n", newIP)
    
    // Special addresses
    fmt.Printf("\nLoopback IPv4: %s\n", net.IPv4(127, 0, 0, 1))
    fmt.Printf("Loopback IPv6: %s\n", net.IPv6loopback)
    fmt.Printf("Unspecified IPv4: %s\n", net.IPv4zero)
    fmt.Printf("Unspecified IPv6: %s\n", net.IPv6unspecified)
    
    // Check special addresses
    loopback := net.ParseIP("127.0.0.1")
    fmt.Printf("Is loopback: %t\n", loopback.IsLoopback())
    
    private := net.ParseIP("192.168.1.1")
    fmt.Printf("Is private: %t\n", private.IsPrivate())
    
    multicast := net.ParseIP("224.0.0.1")
    fmt.Printf("Is multicast: %t\n", multicast.IsMulticast())
}
```

---

## 19. TCP Servers & Clients

### 19.1 Simple TCP Server

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

func main() {
    // Listen on TCP port
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Server listening on :8080")
    
    for {
        // Accept connection
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Accept error: %v\n", err)
            continue
        }
        
        // Handle in goroutine
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    clientAddr := conn.RemoteAddr().String()
    fmt.Printf("Client connected: %s\n", clientAddr)
    
    reader := bufio.NewReader(conn)
    
    for {
        // Read until newline
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Printf("Client %s disconnected\n", clientAddr)
            return
        }
        
        message = strings.TrimSpace(message)
        fmt.Printf("Received from %s: %s\n", clientAddr, message)
        
        // Echo back with modification
        response := fmt.Sprintf("Server received: %s\n", message)
        conn.Write([]byte(response))
        
        // Exit on "quit"
        if message == "quit" {
            return
        }
    }
}
```

### 19.2 TCP Client

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    // Connect to server
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Printf("Connection error: %v\n", err)
        return
    }
    defer conn.Close()
    
    fmt.Println("Connected to server")
    
    // Read from stdin, send to server
    stdinReader := bufio.NewReader(os.Stdin)
    serverReader := bufio.NewReader(conn)
    
    for {
        fmt.Print("> ")
        input, _ := stdinReader.ReadString('\n')
        input = strings.TrimSpace(input)
        
        // Send to server
        fmt.Fprintf(conn, "%s\n", input)
        
        // Read response
        response, err := serverReader.ReadString('\n')
        if err != nil {
            fmt.Printf("Server disconnected: %v\n", err)
            return
        }
        
        fmt.Printf("Server: %s", response)
        
        if input == "quit" {
            return
        }
    }
}
```

### 19.3 Concurrent TCP Server with Context

```go
package main

import (
    "bufio"
    "context"
    "fmt"
    "net"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type Server struct {
    listener net.Listener
    clients  map[net.Conn]struct{}
    mu       sync.Mutex
    wg       sync.WaitGroup
}

func NewServer(addr string) (*Server, error) {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }
    
    return &Server{
        listener: listener,
        clients:  make(map[net.Conn]struct{}),
    }, nil
}

func (s *Server) Start(ctx context.Context) {
    fmt.Printf("Server listening on %s\n", s.listener.Addr())
    
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            select {
            case <-ctx.Done():
                return
            default:
                fmt.Printf("Accept error: %v\n", err)
                continue
            }
        }
        
        s.mu.Lock()
        s.clients[conn] = struct{}{}
        s.mu.Unlock()
        
        s.wg.Add(1)
        go s.handleClient(ctx, conn)
    }
}

func (s *Server) handleClient(ctx context.Context, conn net.Conn) {
    defer s.wg.Done()
    defer func() {
        conn.Close()
        s.mu.Lock()
        delete(s.clients, conn)
        s.mu.Unlock()
    }()
    
    clientAddr := conn.RemoteAddr().String()
    fmt.Printf("Client connected: %s\n", clientAddr)
    
    reader := bufio.NewReader(conn)
    
    for {
        select {
        case <-ctx.Done():
            conn.Write([]byte("Server shutting down\n"))
            return
        default:
        }
        
        // Set read deadline for non-blocking check
        conn.SetReadDeadline(time.Now().Add(1 * time.Second))
        
        message, err := reader.ReadString('\n')
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                continue  // Timeout, check context and retry
            }
            fmt.Printf("Client %s disconnected\n", clientAddr)
            return
        }
        
        fmt.Printf("[%s] %s", clientAddr, message)
        conn.Write([]byte("OK\n"))
    }
}

func (s *Server) Shutdown() {
    fmt.Println("Shutting down server...")
    s.listener.Close()
    
    // Close all client connections
    s.mu.Lock()
    for conn := range s.clients {
        conn.Close()
    }
    s.mu.Unlock()
    
    // Wait for all handlers to finish
    s.wg.Wait()
    fmt.Println("Server stopped")
}

func main() {
    server, err := NewServer(":8080")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    
    // Handle signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        cancel()
        server.Shutdown()
    }()
    
    server.Start(ctx)
}
```

---

## 20. UDP Servers & Clients

### 20.1 Simple UDP Server

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Listen on UDP port
    addr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer conn.Close()
    
    fmt.Println("UDP server listening on :8080")
    
    buffer := make([]byte, 1024)
    
    for {
        // Read from UDP
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Read error: %v\n", err)
            continue
        }
        
        message := string(buffer[:n])
        fmt.Printf("Received from %s: %s\n", clientAddr, message)
        
        // Send response
        response := fmt.Sprintf("Echo: %s", message)
        conn.WriteToUDP([]byte(response), clientAddr)
    }
}
```

### 20.2 UDP Client

```go
package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
    // Resolve server address
    serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Create UDP connection
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer conn.Close()
    
    // Send message
    message := []byte("Hello, UDP Server!")
    _, err = conn.Write(message)
    if err != nil {
        fmt.Printf("Write error: %v\n", err)
        return
    }
    fmt.Printf("Sent: %s\n", message)
    
    // Set read timeout
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    
    // Read response
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("Read error: %v\n", err)
        return
    }
    
    fmt.Printf("Received: %s\n", string(buffer[:n]))
}
```

### 20.3 UDP with IPv4 and IPv6

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Listen on IPv4 only
    addr4, _ := net.ResolveUDPAddr("udp4", ":8080")
    conn4, err := net.ListenUDP("udp4", addr4)
    if err != nil {
        fmt.Printf("IPv4 error: %v\n", err)
    } else {
        defer conn4.Close()
        fmt.Println("Listening on IPv4 :8080")
    }
    
    // Listen on IPv6 only
    addr6, _ := net.ResolveUDPAddr("udp6", "[::]:8081")
    conn6, err := net.ListenUDP("udp6", addr6)
    if err != nil {
        fmt.Printf("IPv6 error: %v\n", err)
    } else {
        defer conn6.Close()
        fmt.Println("Listening on IPv6 [::]:8081")
    }
    
    // Listen on both (dual-stack)
    // Note: behavior depends on OS
    addrDual, _ := net.ResolveUDPAddr("udp", ":8082")
    connDual, err := net.ListenUDP("udp", addrDual)
    if err != nil {
        fmt.Printf("Dual-stack error: %v\n", err)
    } else {
        defer connDual.Close()
        fmt.Println("Listening on dual-stack :8082")
    }
    
    // Keep running
    select {}
}
```

---

*Continue to Part V: Capstone - DNS Server...*
