// Pipeline Pattern - Chaining stages of processing
//
// A pipeline is a series of stages connected by channels. Each stage
// is a group of goroutines running the same function. In each stage:
// - Goroutines receive values from upstream via inbound channels
// - Perform some function on that data
// - Send values downstream via outbound channels
//
// This pattern is great for:
// - Data transformation pipelines
// - ETL (Extract, Transform, Load) operations
// - Stream processing
//
// Usage:
//   go run pipeline.go
package main

import (
	"fmt"
	"strings"
	"sync"
)

func main() {
	fmt.Println("=== Pipeline Example ===")
	fmt.Println()

	// Create input data
	input := []string{
		"  hello world  ",
		"  GO IS AWESOME  ",
		"  Pipeline Pattern  ",
		"  Concurrent Programming  ",
	}

	// Build pipeline: input -> trim -> lowercase -> addPrefix -> output
	//
	//  [input] --> [trim] --> [lowercase] --> [addPrefix] --> [output]
	//

	// Stage 1: Generate values from slice
	source := generate(input)

	// Stage 2: Trim whitespace
	trimmed := trim(source)

	// Stage 3: Convert to lowercase
	lowered := lowercase(trimmed)

	// Stage 4: Add prefix
	prefixed := addPrefix(lowered, ">> ")

	// Consume the pipeline
	fmt.Println("Pipeline output:")
	for result := range prefixed {
		fmt.Println(result)
	}

	fmt.Println()
	fmt.Println("=== Fan-Out / Fan-In Example ===")
	fmt.Println()

	// Fan-out: Multiple goroutines read from the same channel
	// Fan-in: Multiple channels merged into one

	// Generate numbers
	numbers := generateNumbers(1, 10)

	// Fan out to 3 workers that square numbers
	workers := 3
	channels := make([]<-chan int, workers)
	for i := 0; i < workers; i++ {
		channels[i] = square(numbers)
	}

	// Fan in (merge results)
	merged := fanIn(channels...)

	// Consume
	fmt.Println("Squared numbers (order may vary):")
	for n := range merged {
		fmt.Printf("%d ", n)
	}
	fmt.Println()
}

// generate creates a channel and sends strings from a slice
func generate(values []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, v := range values {
			out <- v
		}
	}()
	return out
}

// trim removes leading/trailing whitespace
func trim(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- strings.TrimSpace(s)
		}
	}()
	return out
}

// lowercase converts strings to lowercase
func lowercase(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- strings.ToLower(s)
		}
	}()
	return out
}

// addPrefix adds a prefix to each string
func addPrefix(in <-chan string, prefix string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- prefix + s
		}
	}()
	return out
}

// generateNumbers creates a channel of numbers
func generateNumbers(start, count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i < start+count; i++ {
			out <- i
		}
	}()
	return out
}

// square reads from in, squares each number, sends to out
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// fanIn merges multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}

	// Close output when all inputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
