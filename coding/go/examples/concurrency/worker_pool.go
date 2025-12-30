// Worker Pool - A common Go concurrency pattern
//
// This example demonstrates how to process jobs concurrently using
// a fixed number of worker goroutines. This pattern is useful when:
// - You have many jobs to process
// - Each job is independent
// - You want to limit concurrent operations
//
// Usage:
//   go run worker_pool.go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents work to be done
type Job struct {
	ID      int
	Payload string
}

// Result represents the output of a job
type Result struct {
	JobID    int
	Output   string
	Duration time.Duration
}

func main() {
	// Configuration
	numWorkers := 3
	numJobs := 10

	// Create channels
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{
			ID:      j,
			Payload: fmt.Sprintf("data-%d", j),
		}
	}
	close(jobs) // No more jobs

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("Results:")
	fmt.Println("--------")
	for result := range results {
		fmt.Printf("Job %d: %s (took %v)\n",
			result.JobID, result.Output, result.Duration)
	}
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job.ID)
		start := time.Now()

		// Simulate work
		output := processJob(job)

		duration := time.Since(start)
		fmt.Printf("Worker %d finished job %d\n", id, job.ID)

		results <- Result{
			JobID:    job.ID,
			Output:   output,
			Duration: duration,
		}
	}
}

func processJob(job Job) string {
	// Simulate variable processing time
	sleepTime := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(sleepTime)

	return fmt.Sprintf("processed(%s)", job.Payload)
}
