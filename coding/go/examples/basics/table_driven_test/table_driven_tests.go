// Table-Driven Tests - The Go testing idiom
//
// Table-driven tests are the standard way to write tests in Go.
// Benefits:
// - Easy to add new test cases
// - Clear test case documentation
// - DRY (Don't Repeat Yourself)
// - Easy to see what's being tested
//
// Run tests:
//   go test -v
//   go test -cover
//   go test -bench=.
package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode"
)

// ============================================================
// Functions to test
// ============================================================

// Add adds two integers
func Add(a, b int) int {
	return a + b
}

// Divide divides a by b, returns error if b is zero
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// IsPalindrome checks if a string is a palindrome (ignoring case and spaces)
func IsPalindrome(s string) bool {
	// Clean string: lowercase, remove non-letters
	var clean []rune
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) {
			clean = append(clean, r)
		}
	}
	
	// Check palindrome
	for i := 0; i < len(clean)/2; i++ {
		if clean[i] != clean[len(clean)-1-i] {
			return false
		}
	}
	return true
}

// FizzBuzz returns fizz, buzz, fizzbuzz, or the number
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprintf("%d", n)
	}
}

// ============================================================
// Tests
// ============================================================

func TestAdd(t *testing.T) {
	// Table of test cases
	tests := []struct {
		name     string // Test case name
		a, b     int    // Inputs
		expected int    // Expected output
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -2, 3, 1},
		{"with zero", 5, 0, 5},
		{"both zero", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
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

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		wantError bool
	}{
		{"simple division", 10, 2, 5, false},
		{"division with decimals", 7, 2, 3.5, false},
		{"division by zero", 10, 0, 0, true},
		{"zero dividend", 0, 5, 0, false},
		{"negative numbers", -10, 2, -5, false},
		{"both negative", -10, -2, 5, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("Divide(%v, %v) expected error, got nil",
						tt.a, tt.b)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Divide(%v, %v) unexpected error: %v",
					tt.a, tt.b, err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v; want %v",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"racecar", true},
		{"hello", false},
		{"A man a plan a canal Panama", true},
		{"Was it a car or a cat I saw", true},
		{"", true},
		{"a", true},
		{"ab", false},
		{"Madam", true},
		{"Never odd or even", true},
	}
	
	for _, tt := range tests {
		// Use input as test name (truncated if too long)
		name := tt.input
		if len(name) > 20 {
			name = name[:20] + "..."
		}
		if name == "" {
			name = "empty string"
		}
		
		t.Run(name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{1, "1"},
		{2, "2"},
		{3, "Fizz"},
		{4, "4"},
		{5, "Buzz"},
		{6, "Fizz"},
		{10, "Buzz"},
		{15, "FizzBuzz"},
		{30, "FizzBuzz"},
		{45, "FizzBuzz"},
		{7, "7"},
	}
	
	for _, tt := range tests {
		t.Run(fmt.Sprintf("n=%d", tt.n), func(t *testing.T) {
			result := FizzBuzz(tt.n)
			if result != tt.expected {
				t.Errorf("FizzBuzz(%d) = %q; want %q",
					tt.n, result, tt.expected)
			}
		})
	}
}

// ============================================================
// Benchmarks
// ============================================================

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	// Benchmark with different input sizes
	inputs := []string{
		"a",
		"racecar",
		"A man a plan a canal Panama",
	}
	
	for _, input := range inputs {
		b.Run(fmt.Sprintf("len=%d", len(input)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsPalindrome(input)
			}
		})
	}
}

// ============================================================
// Example tests (appear in documentation)
// ============================================================

func ExampleAdd() {
	fmt.Println(Add(2, 3))
	// Output: 5
}

func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("racecar"))
	fmt.Println(IsPalindrome("hello"))
	// Output:
	// true
	// false
}

func ExampleFizzBuzz() {
	for i := 1; i <= 15; i++ {
		fmt.Println(FizzBuzz(i))
	}
	// Output:
	// 1
	// 2
	// Fizz
	// 4
	// Buzz
	// Fizz
	// 7
	// 8
	// Fizz
	// Buzz
	// 11
	// Fizz
	// 13
	// 14
	// FizzBuzz
}

// Main function for standalone execution
func main() {
	fmt.Println("This file is meant to be run with 'go test'")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  go test -v                    # Run all tests verbosely")
	fmt.Println("  go test -run TestAdd          # Run specific test")
	fmt.Println("  go test -cover                # Show coverage")
	fmt.Println("  go test -bench=.              # Run benchmarks")
	fmt.Println("  go test -bench=. -benchmem    # Benchmarks with memory")
}
