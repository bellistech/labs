# Go Cheat Sheet

A quick reference for Go syntax and common patterns.

## Basics

```go
// Package declaration (every file needs one)
package main

// Imports
import (
    "fmt"
    "strings"
)

// Main function (entry point)
func main() {
    fmt.Println("Hello")
}
```

## Variables

```go
// Declaration with type
var name string = "Alice"
var age int = 30

// Type inference
var city = "NYC"           // string inferred
count := 42                // short declaration (only in functions)

// Multiple variables
var x, y int = 1, 2
a, b := "hello", "world"

// Constants
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)
```

## Types

```go
// Basic types
bool                       // true, false
string                     // "hello"
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
complex64, complex128
byte                       // alias for uint8
rune                       // alias for int32 (Unicode)

// Zero values
var i int      // 0
var f float64  // 0.0
var b bool     // false
var s string   // ""
var p *int     // nil
```

## Control Flow

```go
// If statement
if x > 0 {
    fmt.Println("positive")
} else if x < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}

// If with initialization
if err := doSomething(); err != nil {
    return err
}

// Switch
switch day {
case "Mon", "Tue", "Wed", "Thu", "Fri":
    fmt.Println("weekday")
case "Sat", "Sun":
    fmt.Println("weekend")
default:
    fmt.Println("unknown")
}

// Type switch
switch v := x.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
default:
    fmt.Println("unknown type")
}

// For loop (only loop in Go)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style
for x < 100 {
    x *= 2
}

// Infinite loop
for {
    // break to exit
}

// Range over slice/map
for i, v := range slice {
    fmt.Println(i, v)
}

for key, value := range myMap {
    fmt.Println(key, value)
}
```

## Functions

```go
// Basic function
func add(a, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // naked return
}

// Variadic function
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// Function as value
fn := func(x int) int { return x * 2 }

// Closure
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// Defer (runs when function exits)
func example() {
    defer fmt.Println("cleanup")  // runs last
    fmt.Println("work")
}
```

## Data Structures

```go
// Array (fixed size)
var arr [5]int
arr := [3]int{1, 2, 3}
arr := [...]int{1, 2, 3}  // size inferred

// Slice (dynamic)
slice := []int{1, 2, 3}
slice := make([]int, 5)      // len=5, cap=5
slice := make([]int, 0, 10)  // len=0, cap=10

// Slice operations
slice = append(slice, 4, 5)
sub := slice[1:3]            // elements 1, 2
len(slice)                   // length
cap(slice)                   // capacity

// Map
m := make(map[string]int)
m := map[string]int{"a": 1, "b": 2}

m["key"] = 42                // set
value := m["key"]            // get
value, ok := m["key"]        // check existence
delete(m, "key")             // delete

// Struct
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 30}
p := Person{"Alice", 30}     // positional
p.Name = "Bob"               // access/modify
```

## Methods & Interfaces

```go
// Method (value receiver)
func (p Person) Greet() string {
    return "Hello, " + p.Name
}

// Method (pointer receiver - can modify)
func (p *Person) Birthday() {
    p.Age++
}

// Interface
type Speaker interface {
    Speak() string
}

// Implement by having the methods
func (p Person) Speak() string {
    return p.Name + " says hi"
}

// Use interface
func announce(s Speaker) {
    fmt.Println(s.Speak())
}

// Empty interface (any type)
var x interface{}  // or: var x any
x = 42
x = "hello"

// Type assertion
s := x.(string)           // panics if wrong type
s, ok := x.(string)       // safe version
```

## Error Handling

```go
// Return error
func doThing() error {
    if problem {
        return errors.New("something went wrong")
    }
    return nil
}

// Check error
result, err := doThing()
if err != nil {
    return fmt.Errorf("doThing failed: %w", err)  // wrap
}

// Custom error type
type MyError struct {
    Code    int
    Message string
}

func (e *MyError) Error() string {
    return e.Message
}

// Check error type
if errors.Is(err, ErrNotFound) { }
var myErr *MyError
if errors.As(err, &myErr) { }
```

## Concurrency

```go
// Goroutine
go doSomething()
go func() {
    fmt.Println("anonymous goroutine")
}()

// Channel
ch := make(chan int)        // unbuffered
ch := make(chan int, 10)    // buffered

ch <- 42                    // send
value := <-ch               // receive
close(ch)                   // close

// Range over channel
for v := range ch {
    fmt.Println(v)
}

// Select
select {
case v := <-ch1:
    fmt.Println("from ch1:", v)
case ch2 <- x:
    fmt.Println("sent to ch2")
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("no communication ready")
}

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()

// Mutex
var mu sync.Mutex
mu.Lock()
// critical section
mu.Unlock()

// Context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return result
}
```

## Common Packages

```go
// Strings
strings.Contains(s, "sub")
strings.Split(s, ",")
strings.Join(parts, "-")
strings.ToLower(s)
strings.TrimSpace(s)

// fmt
fmt.Println("hello")
fmt.Printf("name: %s, age: %d\n", name, age)
fmt.Sprintf("formatted %s", s)

// strconv
strconv.Atoi("42")           // string -> int
strconv.Itoa(42)             // int -> string
strconv.ParseFloat("3.14", 64)

// time
time.Now()
time.Sleep(1 * time.Second)
t.Format("2006-01-02 15:04:05")  // reference time!
time.Parse("2006-01-02", "2024-01-15")

// encoding/json
json.Marshal(obj)            // struct -> []byte
json.Unmarshal(data, &obj)   // []byte -> struct

// os
os.Open("file.txt")
os.Create("file.txt")
os.ReadFile("file.txt")
os.WriteFile("file.txt", data, 0644)
os.Args                      // command line args
os.Getenv("HOME")

// net
net.Listen("tcp", ":8080")
net.Dial("tcp", "localhost:8080")
net.ParseIP("192.168.1.1")

// net/http
http.Get("https://example.com")
http.ListenAndServe(":8080", handler)
```

## Project Structure

```
myproject/
├── go.mod              # module definition
├── go.sum              # dependency checksums
├── main.go             # or cmd/myapp/main.go
├── internal/           # private packages
│   └── server/
│       └── server.go
├── pkg/                # public packages
│   └── utils/
│       └── utils.go
└── api/                # API definitions (protobuf, etc.)
```

## Common Commands

```bash
go mod init mymodule     # create new module
go mod tidy              # clean up dependencies
go build                 # compile
go run main.go           # compile and run
go test ./...            # run all tests
go test -v               # verbose
go test -cover           # coverage
go fmt ./...             # format code
go vet ./...             # static analysis
go get package@version   # add dependency
```

## Testing

```go
// file: math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

// Table-driven test
func TestAddTable(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }
    
    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("Add(%d, %d) = %d; want %d",
                tt.a, tt.b, got, tt.want)
        }
    }
}

// Benchmark
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}
```
