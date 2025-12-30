# Part II: Intermediate Go

## 6. Structs & Methods

### 6.1 Defining Structs

```go
package main

import (
    "fmt"
    "time"
)

// Basic struct
type Person struct {
    FirstName string
    LastName  string
    Age       int
    Email     string
}

// Struct with embedded struct
type Employee struct {
    Person            // Embedded (anonymous field)
    EmployeeID string
    Department string
    HireDate   time.Time
}

// Struct with tags (used by encoding/json, databases, etc.)
type User struct {
    ID        int    `json:"id" db:"user_id"`
    Username  string `json:"username" db:"username"`
    Email     string `json:"email,omitempty" db:"email"`
    Password  string `json:"-"` // Ignored in JSON
}

func main() {
    // Create struct instances
    p1 := Person{
        FirstName: "Alice",
        LastName:  "Smith",
        Age:       30,
        Email:     "alice@example.com",
    }
    fmt.Printf("Person: %+v\n", p1)
    
    // Partial initialization (others get zero values)
    p2 := Person{FirstName: "Bob"}
    fmt.Printf("Partial: %+v\n", p2)
    
    // Positional initialization (not recommended)
    p3 := Person{"Charlie", "Brown", 25, "charlie@example.com"}
    fmt.Printf("Positional: %+v\n", p3)
    
    // Access and modify fields
    p1.Age = 31
    fmt.Printf("Updated age: %d\n", p1.Age)
    
    // Pointer to struct
    p4 := &Person{FirstName: "Dave"}
    p4.Age = 40  // Go automatically dereferences
    fmt.Printf("Pointer: %+v\n", *p4)
    
    // Anonymous struct (inline, one-off use)
    config := struct {
        Host string
        Port int
    }{
        Host: "localhost",
        Port: 8080,
    }
    fmt.Printf("Config: %+v\n", config)
    
    // Embedded struct
    emp := Employee{
        Person: Person{
            FirstName: "Eve",
            LastName:  "Johnson",
            Age:       28,
        },
        EmployeeID: "E12345",
        Department: "Engineering",
        HireDate:   time.Now(),
    }
    
    // Access embedded fields directly
    fmt.Printf("Employee name: %s %s\n", emp.FirstName, emp.LastName)
    fmt.Printf("Full employee: %+v\n", emp)
    
    // Compare structs (only if all fields are comparable)
    p5 := Person{FirstName: "Alice", LastName: "Smith", Age: 30, Email: "alice@example.com"}
    fmt.Printf("p1 == p5: %t\n", p1 == p5)
}
```

### 6.2 Methods

```go
package main

import (
    "fmt"
    "math"
)

type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver - works on copy
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Pointer receiver - can modify the struct
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// Pointer receiver - convention: if one method uses pointer, all should
func (r *Rectangle) String() string {
    return fmt.Sprintf("Rectangle(%.2f x %.2f)", r.Width, r.Height)
}

// Methods on any type (must be in same package)
type MyFloat float64

func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

// Circle for comparison
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    
    // Call methods
    fmt.Printf("Area: %.2f\n", rect.Area())
    fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
    fmt.Printf("String: %s\n", rect.String())
    
    // Modify with pointer receiver
    rect.Scale(2)
    fmt.Printf("After scale: %s\n", rect.String())
    
    // Methods can be called on pointer or value
    // Go automatically converts
    rectPtr := &Rectangle{Width: 3, Height: 4}
    fmt.Printf("Pointer area: %.2f\n", rectPtr.Area())
    
    // Custom type method
    f := MyFloat(-42.5)
    fmt.Printf("Abs: %.2f\n", f.Abs())
    
    // Both shapes have Area() - but no interface yet
    circle := Circle{Radius: 5}
    fmt.Printf("Circle area: %.2f\n", circle.Area())
}
```

### 6.3 Constructors and Factory Functions

```go
package main

import (
    "fmt"
    "time"
)

type Server struct {
    Host     string
    Port     int
    Timeout  time.Duration
    MaxConns int
    running  bool  // private field (lowercase)
}

// Constructor function (Go convention: New<Type>)
func NewServer(host string, port int) *Server {
    return &Server{
        Host:     host,
        Port:     port,
        Timeout:  30 * time.Second,  // Default
        MaxConns: 100,               // Default
        running:  false,
    }
}

// Constructor with options pattern
type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) {
        s.Timeout = d
    }
}

func WithMaxConns(n int) ServerOption {
    return func(s *Server) {
        s.MaxConns = n
    }
}

func NewServerWithOptions(host string, port int, opts ...ServerOption) *Server {
    s := &Server{
        Host:     host,
        Port:     port,
        Timeout:  30 * time.Second,
        MaxConns: 100,
    }
    
    for _, opt := range opts {
        opt(s)
    }
    
    return s
}

func (s *Server) Start() error {
    s.running = true
    fmt.Printf("Server starting on %s:%d\n", s.Host, s.Port)
    return nil
}

func (s *Server) IsRunning() bool {
    return s.running
}

func main() {
    // Simple constructor
    s1 := NewServer("localhost", 8080)
    fmt.Printf("Server 1: %+v\n", s1)
    
    // Constructor with options
    s2 := NewServerWithOptions("0.0.0.0", 443,
        WithTimeout(60*time.Second),
        WithMaxConns(1000),
    )
    fmt.Printf("Server 2: %+v\n", s2)
    
    s2.Start()
    fmt.Printf("Is running: %t\n", s2.IsRunning())
}
```

---

## 7. Interfaces

### 7.1 Defining and Implementing Interfaces

```go
package main

import (
    "fmt"
    "math"
)

// Interface definition
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Stringer interface (from fmt package)
type Stringer interface {
    String() string
}

// Rectangle implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
    return fmt.Sprintf("Rectangle(%.2f x %.2f)", r.Width, r.Height)
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
    return fmt.Sprintf("Circle(r=%.2f)", c.Radius)
}

// Triangle implements Shape
type Triangle struct {
    A, B, C float64  // Side lengths
}

func (t Triangle) Area() float64 {
    // Heron's formula
    s := (t.A + t.B + t.C) / 2
    return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
    return t.A + t.B + t.C
}

// Function accepting interface
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// Function accepting Stringer
func PrintString(s fmt.Stringer) {
    fmt.Println(s.String())
}

func main() {
    // Concrete types
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}
    triangle := Triangle{A: 3, B: 4, C: 5}
    
    // Use through interface
    shapes := []Shape{rect, circle, triangle}
    
    fmt.Println("All shapes:")
    for _, s := range shapes {
        PrintShapeInfo(s)
    }
    
    // Only Rectangle and Circle implement Stringer
    fmt.Println("\nStringers:")
    PrintString(rect)
    PrintString(circle)
    // PrintString(triangle)  // Would not compile
}
```

### 7.2 Interface Composition and Empty Interface

```go
package main

import "fmt"

// Small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Composed interfaces
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// Empty interface - can hold any value
func PrintAny(v interface{}) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}

// Type alias for empty interface (Go 1.18+)
func PrintAnyNew(v any) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}

// Type assertion
func ProcessValue(v interface{}) {
    // Type assertion with ok check
    if str, ok := v.(string); ok {
        fmt.Printf("String of length %d: %s\n", len(str), str)
        return
    }
    
    if num, ok := v.(int); ok {
        fmt.Printf("Integer doubled: %d\n", num*2)
        return
    }
    
    fmt.Printf("Unknown type: %T\n", v)
}

// Type switch
func Describe(v interface{}) {
    switch val := v.(type) {
    case nil:
        fmt.Println("nil value")
    case int:
        fmt.Printf("Integer: %d\n", val)
    case float64:
        fmt.Printf("Float: %.2f\n", val)
    case string:
        fmt.Printf("String: %s\n", val)
    case bool:
        fmt.Printf("Boolean: %t\n", val)
    case []int:
        fmt.Printf("Int slice: %v\n", val)
    default:
        fmt.Printf("Unknown: %T = %v\n", val, val)
    }
}

func main() {
    // Empty interface can hold anything
    var anything interface{}
    anything = 42
    PrintAny(anything)
    
    anything = "hello"
    PrintAny(anything)
    
    anything = []int{1, 2, 3}
    PrintAny(anything)
    
    // Slice of any type
    mixed := []interface{}{1, "two", 3.0, true, nil}
    fmt.Println("\nMixed slice:")
    for _, v := range mixed {
        Describe(v)
    }
    
    // Type assertions
    fmt.Println("\nType assertions:")
    ProcessValue("hello")
    ProcessValue(42)
    ProcessValue(3.14)
    
    // Dangerous type assertion (panics if wrong)
    // str := anything.(string)  // Panics if not string
    
    // Safe type assertion
    if str, ok := anything.(string); ok {
        fmt.Printf("It's a string: %s\n", str)
    } else {
        fmt.Println("Not a string")
    }
}
```

### 7.3 Common Interfaces

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "sort"
    "strings"
)

// Implement io.Reader
type RepeatReader struct {
    char  byte
    count int
    pos   int
}

func NewRepeatReader(char byte, count int) *RepeatReader {
    return &RepeatReader{char: char, count: count}
}

func (r *RepeatReader) Read(p []byte) (n int, err error) {
    if r.pos >= r.count {
        return 0, io.EOF
    }
    
    for n = 0; n < len(p) && r.pos < r.count; n++ {
        p[n] = r.char
        r.pos++
    }
    return n, nil
}

// Implement sort.Interface
type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Implement error interface
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func main() {
    // Using our Reader implementation
    reader := NewRepeatReader('A', 10)
    buf := make([]byte, 4)
    
    fmt.Println("Reading from RepeatReader:")
    for {
        n, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
    }
    
    // Standard library readers
    stringReader := strings.NewReader("Hello, World!")
    data, _ := io.ReadAll(stringReader)
    fmt.Printf("\nFrom string reader: %s\n", data)
    
    // Using io.Copy
    var output bytes.Buffer
    stringReader2 := strings.NewReader("Copy this!")
    io.Copy(&output, stringReader2)
    fmt.Printf("Copied: %s\n", output.String())
    
    // Sorting with sort.Interface
    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
    }
    
    fmt.Println("\nBefore sort:", people)
    sort.Sort(ByAge(people))
    fmt.Println("After sort:", people)
    
    // Using error interface
    err := ValidationError{Field: "email", Message: "invalid format"}
    fmt.Printf("\nError: %v\n", err)
}
```

---

## 8. Error Handling

### 8.1 Error Basics

```go
package main

import (
    "errors"
    "fmt"
    "os"
    "strconv"
)

// Simple error creation
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Formatted error
func validateAge(age int) error {
    if age < 0 {
        return fmt.Errorf("invalid age: %d (must be non-negative)", age)
    }
    if age > 150 {
        return fmt.Errorf("invalid age: %d (too old)", age)
    }
    return nil
}

func main() {
    // Basic error handling
    result, err := divide(10, 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %.2f\n", result)
    }
    
    // Check error immediately
    if err := validateAge(-5); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    }
    
    // Real-world: file operations
    file, err := os.Open("nonexistent.txt")
    if err != nil {
        fmt.Printf("File error: %v\n", err)
    } else {
        defer file.Close()
        // Process file...
    }
    
    // Real-world: parsing
    if num, err := strconv.Atoi("abc"); err != nil {
        fmt.Printf("Parse error: %v\n", err)
    } else {
        fmt.Printf("Parsed: %d\n", num)
    }
}
```

### 8.2 Custom Error Types

```go
package main

import (
    "fmt"
    "time"
)

// Custom error type
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: field '%s' with value '%v': %s",
        e.Field, e.Value, e.Message)
}

// Network error with retry info
type NetworkError struct {
    Op        string
    URL       string
    Err       error
    Retryable bool
    RetryAt   time.Time
}

func (e *NetworkError) Error() string {
    return fmt.Sprintf("%s %s: %v", e.Op, e.URL, e.Err)
}

func (e *NetworkError) Unwrap() error {
    return e.Err
}

// Temporary error interface (common pattern)
type TemporaryError interface {
    Temporary() bool
}

func (e *NetworkError) Temporary() bool {
    return e.Retryable
}

// Function returning custom error
func validateUser(name string, age int) error {
    if name == "" {
        return &ValidationError{
            Field:   "name",
            Value:   name,
            Message: "cannot be empty",
        }
    }
    if age < 0 || age > 150 {
        return &ValidationError{
            Field:   "age",
            Value:   age,
            Message: "must be between 0 and 150",
        }
    }
    return nil
}

func main() {
    // Use custom error
    err := validateUser("", 25)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        
        // Type assertion to get details
        if valErr, ok := err.(*ValidationError); ok {
            fmt.Printf("  Field: %s\n", valErr.Field)
            fmt.Printf("  Value: %v\n", valErr.Value)
        }
    }
    
    // Network error with retry
    netErr := &NetworkError{
        Op:        "GET",
        URL:       "https://api.example.com/data",
        Err:       fmt.Errorf("connection timeout"),
        Retryable: true,
        RetryAt:   time.Now().Add(5 * time.Second),
    }
    
    fmt.Printf("\nNetwork error: %v\n", netErr)
    
    // Check if retryable
    if tempErr, ok := interface{}(netErr).(TemporaryError); ok && tempErr.Temporary() {
        fmt.Printf("Will retry at: %v\n", netErr.RetryAt)
    }
}
```

### 8.3 Error Wrapping and Unwrapping (Go 1.13+)

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

// Sentinel errors
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInternal     = errors.New("internal error")
)

func findUser(id int) error {
    if id == 0 {
        // Wrap with context
        return fmt.Errorf("findUser: %w", ErrNotFound)
    }
    if id < 0 {
        return fmt.Errorf("findUser: invalid id %d: %w", id, ErrUnauthorized)
    }
    return nil
}

func processUser(id int) error {
    if err := findUser(id); err != nil {
        // Wrap again with more context
        return fmt.Errorf("processUser: %w", err)
    }
    return nil
}

func main() {
    // Error wrapping chain
    err := processUser(0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        
        // errors.Is - check if error chain contains target
        if errors.Is(err, ErrNotFound) {
            fmt.Println("  -> User not found (errors.Is)")
        }
        if errors.Is(err, ErrUnauthorized) {
            fmt.Println("  -> Unauthorized")
        }
        
        // errors.Unwrap - get wrapped error
        unwrapped := errors.Unwrap(err)
        fmt.Printf("  Unwrapped: %v\n", unwrapped)
        
        // Keep unwrapping
        for unwrapped != nil {
            fmt.Printf("    -> %v\n", unwrapped)
            unwrapped = errors.Unwrap(unwrapped)
        }
    }
    
    // errors.As - extract specific error type
    err2 := processFile("nonexistent.txt")
    if err2 != nil {
        var pathErr *os.PathError
        if errors.As(err2, &pathErr) {
            fmt.Printf("\nPath error: %v\n", pathErr)
            fmt.Printf("  Op: %s\n", pathErr.Op)
            fmt.Printf("  Path: %s\n", pathErr.Path)
        }
    }
}

func processFile(path string) error {
    _, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("processFile: %w", err)
    }
    return nil
}
```

---

## 9. Packages & Modules

### 9.1 Package Basics

```go
// File: mathutil/mathutil.go
package mathutil

import "math"

// Exported (public) - uppercase first letter
func Add(a, b float64) float64 {
    return a + b
}

func Multiply(a, b float64) float64 {
    return a * b
}

// Exported constant
const Pi = math.Pi

// unexported (private) - lowercase first letter
func helper() {
    // Only accessible within this package
}

// Exported type
type Calculator struct {
    Precision int
}

// unexported field
type config struct {
    maxValue float64
}
```

```go
// File: mathutil/advanced.go
package mathutil  // Same package, different file

import "math"

func Sqrt(x float64) float64 {
    return math.Sqrt(x)
}

func Power(base, exp float64) float64 {
    return math.Pow(base, exp)
}
```

```go
// File: main.go
package main

import (
    "fmt"
    "github.com/bellistech/myproject/mathutil"
)

func main() {
    result := mathutil.Add(3, 5)
    fmt.Printf("3 + 5 = %.2f\n", result)
    
    fmt.Printf("Pi = %.4f\n", mathutil.Pi)
    
    calc := mathutil.Calculator{Precision: 2}
    fmt.Printf("Calculator: %+v\n", calc)
}
```

### 9.2 init() Functions

```go
package main

import "fmt"

// Package-level variables (initialized before init)
var startTime = getCurrentTime()

// init() runs automatically before main()
// Can have multiple init() in a file and package
func init() {
    fmt.Println("First init()")
}

func init() {
    fmt.Println("Second init()")
}

func main() {
    fmt.Println("main()")
}

func getCurrentTime() string {
    fmt.Println("Initializing startTime")
    return "2024-01-01"
}

// Output:
// Initializing startTime
// First init()
// Second init()
// main()
```

### 9.3 Go Modules in Detail

```bash
# Create new module
go mod init github.com/bellistech/myproject

# Add dependency
go get github.com/gorilla/mux@v1.8.0

# Add specific version
go get github.com/lib/pq@v1.10.9

# Update all dependencies
go get -u ./...

# Tidy (remove unused, add missing)
go mod tidy

# Vendor dependencies
go mod vendor

# List dependencies
go list -m all

# Why is a dependency needed?
go mod why github.com/gorilla/mux

# Download dependencies
go mod download
```

**go.mod example:**
```
module github.com/bellistech/myproject

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
    google.golang.org/grpc v1.59.0
)

require (
    // Indirect dependencies (auto-managed)
    golang.org/x/net v0.17.0 // indirect
)
```

---

## 10. Testing

### 10.1 Basic Tests

```go
// File: mathutil/mathutil.go
package mathutil

func Add(a, b int) int {
    return a + b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

```go
// File: mathutil/mathutil_test.go
package mathutil

import "testing"

// Test function must start with Test
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

// Table-driven tests (Go idiom)
func TestAddTable(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"mixed", -2, 3, 1},
        {"zero", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

// Test with error
func TestDivide(t *testing.T) {
    // Normal case
    result, err := Divide(10, 2)
    if err != nil {
        t.Fatalf("Divide(10, 2) returned error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10, 2) = %d; want 5", result)
    }
    
    // Error case
    _, err = Divide(10, 0)
    if err == nil {
        t.Error("Divide(10, 0) should return error")
    }
}
```

### 10.2 Benchmarks

```go
// File: mathutil/mathutil_test.go
package mathutil

import "testing"

func BenchmarkAdd(b *testing.B) {
    // b.N is set by testing framework
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}

func BenchmarkAddParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Add(100, 200)
        }
    })
}
```

```bash
# Run benchmarks
go test -bench=. -benchmem

# Output:
# BenchmarkAdd-8           1000000000     0.318 ns/op    0 B/op    0 allocs/op
# BenchmarkAddParallel-8   1000000000     0.098 ns/op    0 B/op    0 allocs/op
```

### 10.3 Example Tests (Documentation)

```go
// File: mathutil/example_test.go
package mathutil_test  // Note: _test suffix for external test package

import (
    "fmt"
    "github.com/bellistech/myproject/mathutil"
)

func ExampleAdd() {
    result := mathutil.Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

func ExampleDivide() {
    result, err := mathutil.Divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result)
    // Output: 5
}
```

### 10.4 Test Helpers and Setup

```go
package mathutil

import (
    "os"
    "testing"
)

// TestMain for setup/teardown
func TestMain(m *testing.M) {
    // Setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Teardown
    teardown()
    
    os.Exit(code)
}

func setup() {
    // Initialize test database, etc.
}

func teardown() {
    // Cleanup
}

// Test helper function
func assertEqual(t *testing.T, got, want int) {
    t.Helper()  // Marks this as helper (better error locations)
    if got != want {
        t.Errorf("got %d; want %d", got, want)
    }
}

func TestWithHelper(t *testing.T) {
    assertEqual(t, Add(2, 3), 5)
    assertEqual(t, Add(-1, 1), 0)
}
```

```bash
# Run tests
go test ./...

# Verbose output
go test -v ./...

# Run specific test
go test -run TestAdd ./...

# Coverage
go test -cover ./...

# Coverage with HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Race detector
go test -race ./...
```

---

*Continue to Part III: Concurrency...*
