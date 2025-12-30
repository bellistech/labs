// TCP Echo Server - A simple example of TCP server programming in Go
//
// This server accepts connections and echoes back whatever the client sends,
// prefixed with "Echo: ". Each client is handled in its own goroutine.
//
// Usage:
//   go run echo_server.go
//
// Test with netcat:
//   nc localhost 8080
//   (type a message and press Enter)
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Listen on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("Echo server listening on :8080")
	log.Println("Test with: nc localhost 8080")

	// Accept connections forever
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)

	// Welcome message
	fmt.Fprintf(conn, "Welcome to Echo Server! Type 'quit' to exit.\n")

	// Read lines from client
	reader := bufio.NewReader(conn)

	for {
		// Read until newline
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v", clientAddr, err)
			return
		}

		// Trim whitespace
		message := strings.TrimSpace(line)
		log.Printf("[%s] Received: %s", clientAddr, message)

		// Check for quit command
		if strings.ToLower(message) == "quit" {
			fmt.Fprintf(conn, "Goodbye!\n")
			return
		}

		// Echo back
		response := fmt.Sprintf("Echo: %s\n", message)
		conn.Write([]byte(response))
	}
}
