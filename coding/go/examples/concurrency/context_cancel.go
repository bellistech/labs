// Context Cancellation - Managing goroutine lifecycles
//
// Context is Go's standard way to:
// - Cancel operations across goroutine boundaries
// - Set deadlines and timeouts
// - Pass request-scoped values
//
// This example demonstrates graceful shutdown of multiple goroutines
// using context cancellation.
//
// Usage:
//   go run context_cancel.go
//   (Press Ctrl+C to trigger shutdown)
package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	fmt.Println("=== Context Cancellation Demo ===")
	fmt.Println("Press Ctrl+C to trigger graceful shutdown")
	fmt.Println()

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %v\n", sig)
		fmt.Println("Cancelling all workers...")
		cancel()
	}()

	// Start workers
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers stopped. Goodbye!")
}

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d started\n", id)

	ticker := time.NewTicker(time.Duration(500+rand.Intn(500)) * time.Millisecond)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ctx.Done():
			// Context cancelled - clean up and exit
			fmt.Printf("Worker %d stopping (processed %d items): %v\n",
				id, count, ctx.Err())
			return

		case <-ticker.C:
			// Do some work
			count++
			fmt.Printf("Worker %d: tick #%d\n", id, count)
		}
	}
}
