// Graceful Shutdown - Production-ready server patterns
//
// This example demonstrates how to build a server that:
// - Handles OS signals (SIGINT, SIGTERM)
// - Completes in-flight requests before shutting down
// - Sets shutdown timeouts
// - Properly cleans up resources
//
// This pattern is essential for:
// - Kubernetes deployments (pod termination)
// - Systemd services
// - Any production server
//
// Usage:
//   go run graceful_shutdown.go
//   (Press Ctrl+C to trigger graceful shutdown)
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Server represents our production-ready server
type Server struct {
	listener    net.Listener
	connections map[net.Conn]struct{}
	connMu      sync.Mutex
	wg          sync.WaitGroup
	
	// Metrics
	totalConns   uint64
	activeConns  int64
	totalQueries uint64
	
	// Shutdown coordination
	shutdownCh chan struct{}
	isShutdown atomic.Bool
}

func NewServer(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	
	return &Server{
		listener:    listener,
		connections: make(map[net.Conn]struct{}),
		shutdownCh:  make(chan struct{}),
	}, nil
}

func (s *Server) Start(ctx context.Context) {
	log.Printf("Server listening on %s", s.listener.Addr())
	
	for {
		// Check if we should stop accepting
		select {
		case <-ctx.Done():
			return
		default:
		}
		
		// Set accept deadline so we can check context periodically
		s.listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))
		
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				continue // Timeout, check context and retry
			}
			if s.isShutdown.Load() {
				return
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		
		// Track connection
		s.connMu.Lock()
		s.connections[conn] = struct{}{}
		s.connMu.Unlock()
		
		atomic.AddUint64(&s.totalConns, 1)
		atomic.AddInt64(&s.activeConns, 1)
		
		// Handle connection
		s.wg.Add(1)
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		conn.Close()
		
		s.connMu.Lock()
		delete(s.connections, conn)
		s.connMu.Unlock()
		
		atomic.AddInt64(&s.activeConns, -1)
		s.wg.Done()
	}()
	
	clientAddr := conn.RemoteAddr().String()
	log.Printf("[%s] Connected", clientAddr)
	
	buf := make([]byte, 1024)
	
	for {
		select {
		case <-ctx.Done():
			// Server shutting down - inform client
			conn.Write([]byte("Server shutting down, goodbye!\n"))
			log.Printf("[%s] Disconnected (server shutdown)", clientAddr)
			return
		default:
		}
		
		// Set read deadline for responsiveness
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		
		n, err := conn.Read(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				continue
			}
			log.Printf("[%s] Disconnected: %v", clientAddr, err)
			return
		}
		
		// Simulate some work
		atomic.AddUint64(&s.totalQueries, 1)
		workDuration := time.Duration(50+rand.Intn(200)) * time.Millisecond
		time.Sleep(workDuration)
		
		// Send response
		response := fmt.Sprintf("Processed: %s", string(buf[:n]))
		conn.Write([]byte(response))
	}
}

func (s *Server) Shutdown(timeout time.Duration) error {
	log.Println("Starting graceful shutdown...")
	s.isShutdown.Store(true)
	
	// Stop accepting new connections
	s.listener.Close()
	
	// Signal all handlers to stop
	close(s.shutdownCh)
	
	// Wait for existing connections with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		log.Println("All connections closed gracefully")
		return nil
	case <-time.After(timeout):
		// Force close remaining connections
		s.connMu.Lock()
		for conn := range s.connections {
			conn.Close()
		}
		s.connMu.Unlock()
		return fmt.Errorf("shutdown timeout, %d connections force-closed", len(s.connections))
	}
}

func (s *Server) Stats() {
	log.Printf("Stats: total_connections=%d, active=%d, queries=%d",
		atomic.LoadUint64(&s.totalConns),
		atomic.LoadInt64(&s.activeConns),
		atomic.LoadUint64(&s.totalQueries))
}

func main() {
	// Create server
	server, err := NewServer(":8080")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	
	// Handle OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	// Start server in background
	go server.Start(ctx)
	
	// Print usage
	log.Println("Server ready. Test with: nc localhost 8080")
	log.Println("Press Ctrl+C to initiate graceful shutdown")
	
	// Periodic stats
	statsTicker := time.NewTicker(5 * time.Second)
	defer statsTicker.Stop()
	
	// Wait for signal
	for {
		select {
		case sig := <-sigCh:
			log.Printf("Received signal: %v", sig)
			
			// Cancel context to stop accepting
			cancel()
			
			// Graceful shutdown with 10 second timeout
			if err := server.Shutdown(10 * time.Second); err != nil {
				log.Printf("Shutdown warning: %v", err)
			}
			
			server.Stats()
			log.Println("Server stopped")
			return
			
		case <-statsTicker.C:
			server.Stats()
		}
	}
}
