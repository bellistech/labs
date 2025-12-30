// Interfaces in Go - Implicit satisfaction and polymorphism
//
// Go's interfaces are satisfied implicitly - there's no "implements"
// keyword. If a type has all the methods an interface requires,
// it automatically satisfies that interface.
//
// This example demonstrates:
// - Interface definition
// - Implicit satisfaction
// - Polymorphism
// - Interface composition
// - Empty interface (any)
// - Type assertions and type switches
//
// Usage:
//   go run interfaces.go
package main

import (
	"fmt"
	"math"
)

// Shape interface - any type with these methods is a Shape
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Stringer interface (like fmt.Stringer)
type Stringer interface {
	String() string
}

// Combined interface - composition
type PrintableShape interface {
	Shape
	Stringer
}

// Rectangle implements Shape and Stringer
type Rectangle struct {
	Width  float64
	Height float64
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

// Circle implements Shape and Stringer
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
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Triangle implements Shape
type Triangle struct {
	A, B, C float64 // Side lengths
}

func (t Triangle) Area() float64 {
	// Heron's formula
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func main() {
	fmt.Println("=== Interfaces Demo ===")
	fmt.Println()

	// Create shapes
	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 7},
		Triangle{A: 3, B: 4, C: 5},
	}

	// Polymorphism - same interface, different implementations
	fmt.Println("All shapes (polymorphism):")
	for _, s := range shapes {
		printShapeInfo(s)
	}

	fmt.Println()
	fmt.Println("=== Type Assertions ===")

	// Type assertion: interface -> concrete type
	var s Shape = Rectangle{Width: 3, Height: 4}

	// Safe type assertion (ok idiom)
	if rect, ok := s.(Rectangle); ok {
		fmt.Printf("It's a rectangle: %v\n", rect)
		fmt.Printf("Width: %.2f, Height: %.2f\n", rect.Width, rect.Height)
	}

	// This would panic if s is not a Circle
	// circle := s.(Circle) // panic!

	fmt.Println()
	fmt.Println("=== Type Switch ===")

	// Type switch - handle different types
	for _, shape := range shapes {
		describeShape(shape)
	}

	fmt.Println()
	fmt.Println("=== Empty Interface (any) ===")

	// any (alias for interface{}) can hold any type
	var anything any

	anything = 42
	fmt.Printf("int: %v (type: %T)\n", anything, anything)

	anything = "hello"
	fmt.Printf("string: %v (type: %T)\n", anything, anything)

	anything = Circle{Radius: 5}
	fmt.Printf("Circle: %v (type: %T)\n", anything, anything)

	// Use type assertion to get back to concrete type
	if c, ok := anything.(Circle); ok {
		fmt.Printf("Got circle back, area: %.2f\n", c.Area())
	}

	fmt.Println()
	fmt.Println("=== Interface Composition ===")

	// PrintableShape requires both Shape and Stringer
	var ps PrintableShape = Rectangle{Width: 8, Height: 6}
	fmt.Printf("PrintableShape: %s, Area: %.2f\n", ps.String(), ps.Area())

	// Triangle doesn't implement Stringer, so it's not a PrintableShape
	// var ps2 PrintableShape = Triangle{3, 4, 5}  // Won't compile!
}

// printShapeInfo accepts any Shape
func printShapeInfo(s Shape) {
	fmt.Printf("  Area: %8.2f, Perimeter: %8.2f", s.Area(), s.Perimeter())

	// Check if it also implements Stringer
	if str, ok := s.(Stringer); ok {
		fmt.Printf(" - %s", str.String())
	} else {
		fmt.Printf(" - %T", s)
	}
	fmt.Println()
}

// describeShape uses type switch
func describeShape(s Shape) {
	switch v := s.(type) {
	case Rectangle:
		fmt.Printf("Rectangle: %.2f x %.2f\n", v.Width, v.Height)
	case Circle:
		fmt.Printf("Circle: radius %.2f\n", v.Radius)
	case Triangle:
		fmt.Printf("Triangle: sides %.2f, %.2f, %.2f\n", v.A, v.B, v.C)
	default:
		fmt.Printf("Unknown shape: %T\n", v)
	}
}
