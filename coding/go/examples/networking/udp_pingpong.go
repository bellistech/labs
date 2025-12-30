// UDP Ping-Pong - Example of UDP server and client in Go
//
// This demonstrates connectionless UDP communication. The server
// responds to "ping" messages with "pong" and echoes other messages.
//
// Usage:
//   # Run server
//   go run udp_pingpong.go server
//
//   # Run client (in another terminal)
//   go run udp_pingpong.go client
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run udp_pingpong.go [server|client]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		runServer()
	case "client":
		runClient()
	default:
		fmt.Println("Unknown command. Use 'server' or 'client'")
		os.Exit(1)
	}
}

func runServer() {
	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		log.Fatalf("ResolveUDPAddr: %v", err)
	}

	// Create UDP connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("ListenUDP: %v", err)
	}
	defer conn.Close()

	log.Println("UDP server listening on :9999")

	buffer := make([]byte, 1024)

	for {
		// Read from UDP (blocking)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("ReadFromUDP error: %v", err)
			continue
		}

		message := string(buffer[:n])
		log.Printf("Received from %s: %s", clientAddr, message)

		// Respond based on message
		var response string
		switch message {
		case "ping":
			response = "pong"
		case "time":
			response = time.Now().Format(time.RFC3339)
		default:
			response = fmt.Sprintf("echo: %s", message)
		}

		// Send response
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			log.Printf("WriteToUDP error: %v", err)
		}
	}
}

func runClient() {
	// Resolve server address
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:9999")
	if err != nil {
		log.Fatalf("ResolveUDPAddr: %v", err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("DialUDP: %v", err)
	}
	defer conn.Close()

	messages := []string{"ping", "time", "hello world", "ping"}

	for _, msg := range messages {
		// Send message
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Write error: %v", err)
			continue
		}
		log.Printf("Sent: %s", msg)

		// Set read deadline
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		// Read response
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		log.Printf("Received: %s", string(buffer[:n]))
		fmt.Println()

		time.Sleep(500 * time.Millisecond)
	}
}
