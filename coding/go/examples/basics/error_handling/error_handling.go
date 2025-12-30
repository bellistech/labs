// Error Handling in Go - Patterns and best practices
//
// Go doesn't have exceptions. Instead, functions return error values
// that must be checked explicitly. This example covers:
//
// - Basic error checking
// - Creating custom errors
// - Error wrapping (Go 1.13+)
// - errors.Is and errors.As
// - Sentinel errors
//
// Usage:
//   go run error_handling.go
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Sentinel errors - predefined errors for specific conditions
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("resource already exists")
)

// Custom error type with additional context
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Another custom error type
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s on %s: %v",
		e.Operation, e.Table, e.Err)
}

// Unwrap allows errors.Is and errors.As to see the wrapped error
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

func main() {
	fmt.Println("=== Error Handling in Go ===")
	fmt.Println()

	// Basic error checking
	fmt.Println("1. Basic Error Checking")
	fmt.Println("-----------------------")

	result, err := divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	result, err = divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println()
	fmt.Println("2. Sentinel Errors")
	fmt.Println("------------------")

	user, err := findUser("alice")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("User not found in database")
	} else if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Printf("Found user: %s\n", user)
	}

	fmt.Println()
	fmt.Println("3. Custom Error Types")
	fmt.Println("--------------------")

	err = validateEmail("not-an-email")
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)

		// Use errors.As to check for specific error type
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  Field: %s\n", valErr.Field)
			fmt.Printf("  Message: %s\n", valErr.Message)
		}
	}

	fmt.Println()
	fmt.Println("4. Error Wrapping")
	fmt.Println("-----------------")

	err = processFile("nonexistent.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Check if the underlying error is a path error
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			fmt.Printf("  Path: %s\n", pathErr.Path)
			fmt.Printf("  Op: %s\n", pathErr.Op)
		}

		// Check if it wraps os.ErrNotExist
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("  -> File does not exist!")
		}
	}

	fmt.Println()
	fmt.Println("5. Error Chain")
	fmt.Println("--------------")

	err = performDatabaseOperation()
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)

		// Check if it's a DatabaseError
		var dbErr *DatabaseError
		if errors.As(err, &dbErr) {
			fmt.Printf("  Operation: %s\n", dbErr.Operation)
			fmt.Printf("  Table: %s\n", dbErr.Table)
		}

		// Check if the root cause is ErrNotFound
		if errors.Is(err, ErrNotFound) {
			fmt.Println("  -> Root cause: resource not found")
		}
	}

	fmt.Println()
	fmt.Println("6. Multiple Return Values Pattern")
	fmt.Println("----------------------------------")

	// Go convention: return (result, error)
	data, err := fetchData("https://api.example.com/data")
	if err != nil {
		fmt.Printf("Fetch failed: %v\n", err)
		// Handle error (retry, log, return, etc.)
	} else {
		fmt.Printf("Fetched %d bytes\n", len(data))
	}
}

// Basic error creation
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Using sentinel errors
func findUser(username string) (string, error) {
	// Simulated database lookup
	users := map[string]string{
		"bob":   "Bob Smith",
		"carol": "Carol Jones",
	}

	if name, ok := users[username]; ok {
		return name, nil
	}
	return "", ErrNotFound
}

// Custom error type
func validateEmail(email string) error {
	// Simple validation
	if len(email) < 5 || email[0] == '@' {
		return &ValidationError{
			Field:   "email",
			Message: "must be at least 5 characters and not start with @",
		}
	}
	return nil
}

// Error wrapping with fmt.Errorf and %w
func processFile(filename string) error {
	_, err := os.Open(filename)
	if err != nil {
		// Wrap the original error with context
		return fmt.Errorf("failed to process file %s: %w", filename, err)
	}
	return nil
}

// Nested error wrapping
func performDatabaseOperation() error {
	// Simulate a low-level error
	lowLevelErr := ErrNotFound

	// Wrap in a database error
	dbErr := &DatabaseError{
		Operation: "SELECT",
		Table:     "users",
		Err:       lowLevelErr,
	}

	// Add more context at a higher level
	return fmt.Errorf("user service failed: %w", dbErr)
}

// Simulated data fetch
func fetchData(url string) ([]byte, error) {
	// Simulated - would normally make HTTP request
	if url == "" {
		return nil, errors.New("empty URL")
	}
	return []byte("simulated data"), nil
}
