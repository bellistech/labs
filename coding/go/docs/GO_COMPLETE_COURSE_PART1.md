# Go Crash Course: From Zero to Production Systems

A comprehensive guide from absolute beginner to building production network services in Go. This course culminates in building a fully functional IPv6-capable authoritative DNS server.

**Target Audience:** Script writers and programmers familiar with Python, Ruby, Bash, or similar languages wanting to build production backend systems.

**Final Project:** A complete authoritative DNS server supporting A, AAAA, CNAME, MX, TXT, NS, and SOA records with IPv4/IPv6 dual-stack support.

---

## Table of Contents

### Part I: Go Fundamentals
1. [Introduction & Setup](#1-introduction--setup)
2. [Variables, Types & Constants](#2-variables-types--constants)
3. [Control Flow](#3-control-flow)
4. [Functions](#4-functions)
5. [Data Structures](#5-data-structures)

### Part II: Intermediate Go
6. [Structs & Methods](#6-structs--methods)
7. [Interfaces](#7-interfaces)
8. [Error Handling](#8-error-handling)
9. [Packages & Modules](#9-packages--modules)
10. [Testing](#10-testing)

### Part III: Concurrency
11. [Goroutines](#11-goroutines)
12. [Channels](#12-channels)
13. [Select & Timeouts](#13-select--timeouts)
14. [Sync Primitives](#14-sync-primitives)
15. [Context](#15-context)

### Part IV: Systems Programming
16. [File I/O](#16-file-io)
17. [Binary Data & Encoding](#17-binary-data--encoding)
18. [Network Programming Basics](#18-network-programming-basics)
19. [TCP Servers & Clients](#19-tcp-servers--clients)
20. [UDP Servers & Clients](#20-udp-servers--clients)

### Part V: Capstone - DNS Server
21. [DNS Protocol Deep Dive](#21-dns-protocol-deep-dive)
22. [Parsing DNS Messages](#22-parsing-dns-messages)
23. [Building DNS Responses](#23-building-dns-responses)
24. [Zone File Parsing](#24-zone-file-parsing)
25. [Complete DNS Server](#25-complete-dns-server)
26. [Testing & Deployment](#26-testing--deployment)

---

# Part I: Go Fundamentals

## 1. Introduction & Setup

### 1.1 Why Go?

Go (Golang) was designed at Google for building reliable, efficient software. Key characteristics:

| Feature | Benefit |
|---------|---------|
| **Compiled** | Single static binary, no runtime dependencies |
| **Strongly typed** | Catches bugs at compile time |
| **Garbage collected** | No manual memory management |
| **Built-in concurrency** | Goroutines and channels |
| **Fast compilation** | Seconds, not minutes |
| **Simple syntax** | ~25 keywords, easy to read |

### 1.2 Installation

```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang-go

# Or download from https://go.dev/dl/

# Verify
go version
# Output: go version go1.21.0 linux/amd64
```

### 1.3 Hello World

Create a file `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Run it:
```bash
go run main.go
# Output: Hello, World!

# Or compile to binary
go build -o hello main.go
./hello
```

### 1.4 Understanding the Code

```go
package main        // Every file belongs to a package
                    // 'main' package = executable program

import "fmt"        // Import the format package (printing)
                    // Standard library, no installation needed

func main() {       // Entry point - must be in package main
    fmt.Println("Hello, World!")  // Print with newline
}
```

### 1.5 Go Modules (Project Setup)

```bash
# Create a new project
mkdir myproject
cd myproject
go mod init github.com/bellistech/myproject

# This creates go.mod - like package.json or requirements.txt
cat go.mod
# module github.com/bellistech/myproject
# go 1.21
```

### 1.6 Project Structure Convention

```
myproject/
├── go.mod              # Module definition
├── go.sum              # Dependency checksums (auto-generated)
├── main.go             # Entry point for executables
├── cmd/                # Multiple executables
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── internal/           # Private packages (can't be imported externally)
│   └── parser/
│       └── parser.go
├── pkg/                # Public packages (can be imported)
│   └── dns/
│       └── dns.go
└── README.md
```

---

## 2. Variables, Types & Constants

### 2.1 Variable Declaration

```go
package main

import "fmt"

func main() {
    // Explicit declaration with type
    var name string = "Alice"
    var age int = 30
    var active bool = true
    
    // Type inference (compiler determines type)
    var city = "New York"  // string
    var count = 42         // int
    
    // Short declaration (most common, only inside functions)
    country := "USA"       // string
    score := 95.5          // float64
    
    // Multiple declarations
    var x, y, z int = 1, 2, 3
    a, b, c := "hello", 42, true
    
    // Zero values (uninitialized variables get default)
    var unsetString string  // ""
    var unsetInt int        // 0
    var unsetBool bool      // false
    var unsetFloat float64  // 0.0
    
    fmt.Println(name, age, active, city, count, country, score)
    fmt.Println(x, y, z, a, b, c)
    fmt.Printf("Zero values: '%s', %d, %t, %f\n", 
        unsetString, unsetInt, unsetBool, unsetFloat)
}
```

### 2.2 Basic Types

```go
package main

import "fmt"

func main() {
    // Integers (signed)
    var i8 int8 = 127              // -128 to 127
    var i16 int16 = 32767          // -32768 to 32767
    var i32 int32 = 2147483647     // -2^31 to 2^31-1
    var i64 int64 = 9223372036854775807
    var i int = 42                 // Platform dependent (32 or 64 bit)
    
    // Integers (unsigned)
    var u8 uint8 = 255             // 0 to 255 (alias: byte)
    var u16 uint16 = 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615
    
    // Floats
    var f32 float32 = 3.14
    var f64 float64 = 3.141592653589793  // Default for float literals
    
    // Complex numbers
    var c64 complex64 = 1 + 2i
    var c128 complex128 = 1 + 2i
    
    // Boolean
    var isReady bool = true
    
    // String (immutable UTF-8)
    var message string = "Hello, 世界"
    
    // Rune (Unicode code point, alias for int32)
    var letter rune = 'A'        // 65
    var chinese rune = '世'      // 19990
    
    // Byte (alias for uint8)
    var b byte = 'A'             // 65
    
    fmt.Printf("int8: %d, int64: %d\n", i8, i64)
    fmt.Printf("uint8: %d, uint64: %d\n", u8, u64)
    fmt.Printf("float32: %f, float64: %.15f\n", f32, f64)
    fmt.Printf("complex: %v\n", c128)
    fmt.Printf("string: %s, length: %d bytes\n", message, len(message))
    fmt.Printf("rune: %c (%d), chinese: %c (%d)\n", letter, letter, chinese, chinese)
}
```

### 2.3 Type Conversions

Go requires explicit type conversions (no implicit casting):

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // Numeric conversions
    var i int = 42
    var f float64 = float64(i)      // int to float64
    var u uint = uint(i)            // int to uint
    
    // String conversions
    var s string = strconv.Itoa(i)  // int to string: "42"
    var i2, _ = strconv.Atoi("123") // string to int: 123
    
    // Float to string
    var fs string = strconv.FormatFloat(3.14, 'f', 2, 64)  // "3.14"
    var f2, _ = strconv.ParseFloat("3.14", 64)             // 3.14
    
    // Byte slice to string and back
    var bytes []byte = []byte("hello")    // string to []byte
    var str string = string(bytes)        // []byte to string
    
    fmt.Printf("int: %d, float: %f, uint: %d\n", i, f, u)
    fmt.Printf("string: %s, back to int: %d\n", s, i2)
    fmt.Printf("float string: %s, back to float: %f\n", fs, f2)
    fmt.Printf("bytes: %v, string: %s\n", bytes, str)
}
```

### 2.4 Constants

```go
package main

import "fmt"

// Package-level constants
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)

// Typed constants
const MaxSize int = 1024

// iota - auto-incrementing constant generator
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
    Wednesday      // 3
    Thursday       // 4
    Friday         // 5
    Saturday       // 6
)

// iota with expressions
const (
    _  = iota             // 0 (ignore first)
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20 = 1048576
    GB                    // 1 << 30 = 1073741824
    TB                    // 1 << 40
)

// Bit flags with iota
const (
    FlagRead  = 1 << iota  // 1
    FlagWrite              // 2
    FlagExec               // 4
)

func main() {
    fmt.Printf("Pi: %f\n", Pi)
    fmt.Printf("Tuesday: %d\n", Tuesday)
    fmt.Printf("KB: %d, MB: %d, GB: %d\n", KB, MB, GB)
    fmt.Printf("Read|Write: %d\n", FlagRead|FlagWrite)
}
```

---

## 3. Control Flow

### 3.1 If Statements

```go
package main

import "fmt"

func main() {
    x := 10
    
    // Basic if
    if x > 5 {
        fmt.Println("x is greater than 5")
    }
    
    // If-else
    if x > 15 {
        fmt.Println("x is greater than 15")
    } else {
        fmt.Println("x is not greater than 15")
    }
    
    // If-else if-else
    if x < 0 {
        fmt.Println("negative")
    } else if x < 10 {
        fmt.Println("single digit")
    } else {
        fmt.Println("double digit or more")
    }
    
    // If with initialization statement (scope limited to if block)
    if y := x * 2; y > 15 {
        fmt.Printf("y (%d) is greater than 15\n", y)
    }
    // y is not accessible here
    
    // Common pattern: error checking
    if err := doSomething(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}

func doSomething() error {
    return nil
}
```

### 3.2 Switch Statements

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // Basic switch
    day := time.Now().Weekday()
    switch day {
    case time.Saturday, time.Sunday:
        fmt.Println("Weekend!")
    case time.Friday:
        fmt.Println("TGIF!")
    default:
        fmt.Println("Weekday")
    }
    
    // Switch with initialization
    switch os := runtime.GOOS; os {
    case "darwin":
        fmt.Println("macOS")
    case "linux":
        fmt.Println("Linux")
    default:
        fmt.Printf("Other: %s\n", os)
    }
    
    // Switch without expression (like if-else chain)
    hour := time.Now().Hour()
    switch {
    case hour < 12:
        fmt.Println("Morning")
    case hour < 17:
        fmt.Println("Afternoon")
    default:
        fmt.Println("Evening")
    }
    
    // Type switch
    var i interface{} = "hello"
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
    
    // Fallthrough (explicit, unlike C)
    n := 1
    switch n {
    case 1:
        fmt.Println("one")
        fallthrough
    case 2:
        fmt.Println("two (fallthrough)")
    case 3:
        fmt.Println("three")
    }
}
```

### 3.3 For Loops

Go has only one loop construct: `for`. It can do everything:

```go
package main

import "fmt"

func main() {
    // Traditional for loop
    for i := 0; i < 5; i++ {
        fmt.Printf("i = %d\n", i)
    }
    
    // While-style loop
    j := 0
    for j < 5 {
        fmt.Printf("j = %d\n", j)
        j++
    }
    
    // Infinite loop
    k := 0
    for {
        if k >= 3 {
            break
        }
        fmt.Printf("k = %d\n", k)
        k++
    }
    
    // Range over slice
    fruits := []string{"apple", "banana", "cherry"}
    for index, fruit := range fruits {
        fmt.Printf("%d: %s\n", index, fruit)
    }
    
    // Range - ignore index
    for _, fruit := range fruits {
        fmt.Printf("Fruit: %s\n", fruit)
    }
    
    // Range - only index
    for index := range fruits {
        fmt.Printf("Index: %d\n", index)
    }
    
    // Range over map
    ages := map[string]int{"Alice": 30, "Bob": 25}
    for name, age := range ages {
        fmt.Printf("%s is %d\n", name, age)
    }
    
    // Range over string (iterates runes, not bytes)
    for i, r := range "Hello, 世界" {
        fmt.Printf("Index %d: %c (U+%04X)\n", i, r, r)
    }
    
    // Continue and break
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            continue  // Skip even numbers
        }
        if i > 7 {
            break     // Stop at 7
        }
        fmt.Printf("Odd: %d\n", i)
    }
    
    // Labeled break (break outer loop)
    outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                break outer
            }
            fmt.Printf("(%d, %d)\n", i, j)
        }
    }
}
```

---

## 4. Functions

### 4.1 Basic Functions

```go
package main

import "fmt"

// No parameters, no return
func sayHello() {
    fmt.Println("Hello!")
}

// With parameters
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// With return value
func add(a, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Named return values
func rectangle(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return  // Naked return - returns named values
}

// Variadic function (variable number of arguments)
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

func main() {
    sayHello()
    greet("Alice")
    
    result := add(3, 5)
    fmt.Printf("3 + 5 = %d\n", result)
    
    quotient, err := divide(10, 3)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("10 / 3 = %f\n", quotient)
    }
    
    area, perimeter := rectangle(5, 3)
    fmt.Printf("Area: %f, Perimeter: %f\n", area, perimeter)
    
    fmt.Printf("Sum: %d\n", sum(1, 2, 3, 4, 5))
    
    // Spreading a slice
    nums := []int{10, 20, 30}
    fmt.Printf("Sum of slice: %d\n", sum(nums...))
}
```

### 4.2 First-Class Functions

```go
package main

import (
    "fmt"
    "sort"
)

// Function as parameter
func apply(f func(int) int, value int) int {
    return f(value)
}

// Function returning function (closure)
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

// Function as field in struct
type Operation struct {
    Name string
    Fn   func(int, int) int
}

func main() {
    // Function variable
    double := func(x int) int {
        return x * 2
    }
    fmt.Printf("Double 5: %d\n", double(5))
    
    // Pass function as argument
    result := apply(double, 10)
    fmt.Printf("Apply double to 10: %d\n", result)
    
    // Anonymous function (immediately invoked)
    func(msg string) {
        fmt.Println(msg)
    }("Hello from anonymous!")
    
    // Closure
    triple := makeMultiplier(3)
    fmt.Printf("Triple 7: %d\n", triple(7))
    
    // Closure capturing variable
    counter := 0
    increment := func() int {
        counter++
        return counter
    }
    fmt.Println(increment(), increment(), increment())  // 1 2 3
    
    // Function in struct
    ops := []Operation{
        {"add", func(a, b int) int { return a + b }},
        {"sub", func(a, b int) int { return a - b }},
        {"mul", func(a, b int) int { return a * b }},
    }
    for _, op := range ops {
        fmt.Printf("%s(10, 3) = %d\n", op.Name, op.Fn(10, 3))
    }
    
    // Real-world example: custom sort
    names := []string{"Charlie", "Alice", "Bob"}
    sort.Slice(names, func(i, j int) bool {
        return names[i] < names[j]
    })
    fmt.Printf("Sorted: %v\n", names)
}
```

### 4.3 Defer, Panic, and Recover

```go
package main

import "fmt"

func main() {
    // Defer - executes when function returns (LIFO order)
    fmt.Println("=== Defer ===")
    defer fmt.Println("1st defer")
    defer fmt.Println("2nd defer")
    defer fmt.Println("3rd defer")
    fmt.Println("Main body")
    // Output: Main body, 3rd defer, 2nd defer, 1st defer
    
    // Common pattern: cleanup
    fmt.Println("\n=== Cleanup Pattern ===")
    processFile()
    
    // Panic and recover
    fmt.Println("\n=== Panic & Recover ===")
    safeDivide(10, 2)
    safeDivide(10, 0)  // Would panic, but we recover
    fmt.Println("Program continues!")
}

func processFile() {
    fmt.Println("Opening file...")
    defer fmt.Println("Closing file (deferred)")
    
    fmt.Println("Processing file...")
    // Even if panic happens here, defer runs
}

func safeDivide(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
        }
    }()
    
    if b == 0 {
        panic("division by zero!")
    }
    fmt.Printf("%d / %d = %d\n", a, b, a/b)
}
```

---

## 5. Data Structures

### 5.1 Arrays

```go
package main

import "fmt"

func main() {
    // Arrays have fixed size (part of the type)
    var a [5]int                    // [0, 0, 0, 0, 0]
    a[0] = 10
    a[4] = 50
    fmt.Printf("Array a: %v\n", a)
    
    // Initialize with values
    b := [5]int{1, 2, 3, 4, 5}
    fmt.Printf("Array b: %v\n", b)
    
    // Compiler counts elements
    c := [...]string{"apple", "banana", "cherry"}
    fmt.Printf("Array c: %v (len=%d)\n", c, len(c))
    
    // Specific indices
    d := [5]int{0: 100, 4: 400}  // [100, 0, 0, 0, 400]
    fmt.Printf("Array d: %v\n", d)
    
    // Arrays are values (copied on assignment)
    e := b
    e[0] = 999
    fmt.Printf("b[0]: %d, e[0]: %d (different!)\n", b[0], e[0])
    
    // Array length and capacity
    fmt.Printf("len(b): %d, cap(b): %d\n", len(b), cap(b))
}
```

### 5.2 Slices

Slices are the workhorses of Go - dynamic, flexible views into arrays:

```go
package main

import "fmt"

func main() {
    // Create slice with make
    s1 := make([]int, 5)        // len=5, cap=5, all zeros
    s2 := make([]int, 3, 10)    // len=3, cap=10
    fmt.Printf("s1: %v (len=%d, cap=%d)\n", s1, len(s1), cap(s1))
    fmt.Printf("s2: %v (len=%d, cap=%d)\n", s2, len(s2), cap(s2))
    
    // Slice literal
    fruits := []string{"apple", "banana", "cherry"}
    fmt.Printf("fruits: %v\n", fruits)
    
    // Append (may allocate new underlying array)
    fruits = append(fruits, "date")
    fruits = append(fruits, "elderberry", "fig")
    fmt.Printf("fruits after append: %v\n", fruits)
    
    // Slicing (creates view, not copy)
    slice1 := fruits[1:4]    // ["banana", "cherry", "date"]
    slice2 := fruits[:3]     // ["apple", "banana", "cherry"]
    slice3 := fruits[2:]     // ["cherry", "date", "elderberry", "fig"]
    fmt.Printf("slice1: %v\n", slice1)
    fmt.Printf("slice2: %v\n", slice2)
    fmt.Printf("slice3: %v\n", slice3)
    
    // Slices share underlying array!
    slice1[0] = "BANANA"
    fmt.Printf("fruits after modifying slice1: %v\n", fruits)
    
    // Copy slice (creates independent copy)
    copyFruits := make([]string, len(fruits))
    copy(copyFruits, fruits)
    copyFruits[0] = "APPLE"
    fmt.Printf("fruits: %v\n", fruits)
    fmt.Printf("copyFruits: %v\n", copyFruits)
    
    // Nil slice vs empty slice
    var nilSlice []int
    emptySlice := []int{}
    fmt.Printf("nil slice: %v (nil=%t)\n", nilSlice, nilSlice == nil)
    fmt.Printf("empty slice: %v (nil=%t)\n", emptySlice, emptySlice == nil)
    // Both have len=0, but nil slice is nil
    
    // Multi-dimensional slice
    matrix := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    fmt.Printf("matrix[1][2] = %d\n", matrix[1][2])  // 6
    
    // Remove element (by slicing and appending)
    nums := []int{1, 2, 3, 4, 5}
    i := 2  // Remove index 2
    nums = append(nums[:i], nums[i+1:]...)
    fmt.Printf("After removing index 2: %v\n", nums)  // [1, 2, 4, 5]
    
    // Insert element
    nums = []int{1, 2, 4, 5}
    i = 2
    nums = append(nums[:i], append([]int{3}, nums[i:]...)...)
    fmt.Printf("After inserting 3 at index 2: %v\n", nums)  // [1, 2, 3, 4, 5]
}
```

### 5.3 Maps

```go
package main

import "fmt"

func main() {
    // Create map with make
    ages := make(map[string]int)
    ages["Alice"] = 30
    ages["Bob"] = 25
    ages["Charlie"] = 35
    fmt.Printf("ages: %v\n", ages)
    
    // Map literal
    scores := map[string]int{
        "Alice":   95,
        "Bob":     87,
        "Charlie": 92,
    }
    fmt.Printf("scores: %v\n", scores)
    
    // Access value
    aliceAge := ages["Alice"]
    fmt.Printf("Alice's age: %d\n", aliceAge)
    
    // Check if key exists
    age, exists := ages["Dave"]
    if exists {
        fmt.Printf("Dave's age: %d\n", age)
    } else {
        fmt.Println("Dave not found")
    }
    
    // Shorthand check
    if age, ok := ages["Bob"]; ok {
        fmt.Printf("Bob's age: %d\n", age)
    }
    
    // Delete key
    delete(ages, "Charlie")
    fmt.Printf("After delete: %v\n", ages)
    
    // Iterate (order is randomized!)
    fmt.Println("Iterating:")
    for name, score := range scores {
        fmt.Printf("  %s: %d\n", name, score)
    }
    
    // Nil map (can read, but cannot write)
    var nilMap map[string]int
    fmt.Printf("nil map['x']: %d\n", nilMap["x"])  // Returns zero value
    // nilMap["x"] = 1  // PANIC: assignment to nil map
    
    // Map of slices
    graph := make(map[string][]string)
    graph["A"] = []string{"B", "C"}
    graph["B"] = []string{"A", "D"}
    fmt.Printf("graph: %v\n", graph)
    
    // Set using map (Go doesn't have sets)
    set := make(map[string]struct{})
    set["apple"] = struct{}{}
    set["banana"] = struct{}{}
    
    if _, exists := set["apple"]; exists {
        fmt.Println("apple is in set")
    }
    
    // Count occurrences
    words := []string{"apple", "banana", "apple", "cherry", "apple"}
    wordCount := make(map[string]int)
    for _, word := range words {
        wordCount[word]++
    }
    fmt.Printf("Word counts: %v\n", wordCount)
}
```

### 5.4 Strings and Runes

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    s := "Hello, 世界!"
    
    // String length (bytes vs runes)
    fmt.Printf("String: %s\n", s)
    fmt.Printf("Byte length: %d\n", len(s))           // 14 bytes
    fmt.Printf("Rune count: %d\n", utf8.RuneCountInString(s))  // 10 runes
    
    // Iterate by byte (don't do this for Unicode)
    fmt.Println("\nBy byte:")
    for i := 0; i < len(s); i++ {
        fmt.Printf("%d: %x\n", i, s[i])
    }
    
    // Iterate by rune (correct for Unicode)
    fmt.Println("\nBy rune:")
    for i, r := range s {
        fmt.Printf("%d: %c (U+%04X)\n", i, r, r)
    }
    
    // Convert to rune slice for manipulation
    runes := []rune(s)
    runes[7] = '世'  // Modify
    fmt.Printf("Modified: %s\n", string(runes))
    
    // String builder (efficient concatenation)
    var builder strings.Builder
    for i := 0; i < 5; i++ {
        builder.WriteString("Go")
        builder.WriteString(" ")
    }
    result := builder.String()
    fmt.Printf("Built string: %s\n", result)
    
    // Common string operations
    fmt.Printf("Contains '世': %t\n", strings.Contains(s, "世"))
    fmt.Printf("HasPrefix 'Hello': %t\n", strings.HasPrefix(s, "Hello"))
    fmt.Printf("Index of ',': %d\n", strings.Index(s, ","))
    fmt.Printf("ToUpper: %s\n", strings.ToUpper(s))
    fmt.Printf("Replace: %s\n", strings.Replace(s, "世界", "World", 1))
    fmt.Printf("Split: %v\n", strings.Split("a,b,c", ","))
    fmt.Printf("Join: %s\n", strings.Join([]string{"a", "b", "c"}, "-"))
    fmt.Printf("TrimSpace: '%s'\n", strings.TrimSpace("  hello  "))
}
```

---

*Continue to Part II: Intermediate Go...*
