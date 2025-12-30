// Binary Protocol Parsing - Working with network wire formats
//
// This example demonstrates how to parse and build binary protocols
// like those used in DNS, TCP/IP headers, and other network services.
//
// Key concepts:
// - Big-endian (network byte order) vs little-endian
// - Bit manipulation for flags
// - Using encoding/binary package
// - Struct packing/unpacking
//
// Usage:
//   go run binary_protocol.go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Simulated protocol header (like a simplified DNS or custom protocol)
// 
// Wire format (16 bytes total):
//   0                   1                   2                   3
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |           Message ID          |             Flags             |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                          Sequence                            |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                          Timestamp                           |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                       Payload Length                         |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// Header represents our protocol header
type Header struct {
	MessageID     uint16
	Flags         uint16
	Sequence      uint32
	Timestamp     uint32
	PayloadLength uint32
}

// Flag bit positions
const (
	FlagRequest   uint16 = 1 << 15 // Bit 15: Request (0) / Response (1)
	FlagError     uint16 = 1 << 14 // Bit 14: Error flag
	FlagEncrypted uint16 = 1 << 13 // Bit 13: Payload encrypted
	FlagCompressed uint16 = 1 << 12 // Bit 12: Payload compressed
	// Bits 0-11: Reserved or protocol-specific
)

func main() {
	fmt.Println("=== Binary Protocol Parsing Demo ===")
	fmt.Println()

	// Create a header
	original := Header{
		MessageID:     0x1234,
		Flags:         FlagRequest | FlagEncrypted,
		Sequence:      42,
		Timestamp:     1700000000,
		PayloadLength: 256,
	}

	fmt.Println("Original header:")
	printHeader(&original)

	// Serialize to bytes (network byte order = big-endian)
	data := serializeHeader(&original)
	fmt.Printf("\nSerialized (%d bytes):\n", len(data))
	hexDump(data)

	// Parse back from bytes
	parsed, err := parseHeader(data)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Println("\nParsed header:")
	printHeader(parsed)

	fmt.Println()
	fmt.Println("=== Manual Byte Manipulation ===")
	fmt.Println()

	// Show manual parsing
	manualParseDemo(data)

	fmt.Println()
	fmt.Println("=== Bit Manipulation for Flags ===")
	fmt.Println()

	flagsDemo()
}

// serializeHeader converts Header to bytes (big-endian)
func serializeHeader(h *Header) []byte {
	buf := new(bytes.Buffer)
	
	// Write each field in network byte order (big-endian)
	binary.Write(buf, binary.BigEndian, h.MessageID)
	binary.Write(buf, binary.BigEndian, h.Flags)
	binary.Write(buf, binary.BigEndian, h.Sequence)
	binary.Write(buf, binary.BigEndian, h.Timestamp)
	binary.Write(buf, binary.BigEndian, h.PayloadLength)
	
	return buf.Bytes()
}

// parseHeader converts bytes back to Header
func parseHeader(data []byte) (*Header, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf("header too short: %d bytes", len(data))
	}

	h := &Header{}
	reader := bytes.NewReader(data)

	binary.Read(reader, binary.BigEndian, &h.MessageID)
	binary.Read(reader, binary.BigEndian, &h.Flags)
	binary.Read(reader, binary.BigEndian, &h.Sequence)
	binary.Read(reader, binary.BigEndian, &h.Timestamp)
	binary.Read(reader, binary.BigEndian, &h.PayloadLength)

	return h, nil
}

// Manual parsing without encoding/binary.Read
func manualParseDemo(data []byte) {
	// Parse uint16 (2 bytes, big-endian)
	messageID := binary.BigEndian.Uint16(data[0:2])
	fmt.Printf("MessageID: 0x%04X (manual parse)\n", messageID)

	// Parse uint32 (4 bytes, big-endian)
	sequence := binary.BigEndian.Uint32(data[4:8])
	fmt.Printf("Sequence: %d (manual parse)\n", sequence)

	// Manual byte-by-byte (educational)
	msgID := uint16(data[0])<<8 | uint16(data[1])
	fmt.Printf("MessageID: 0x%04X (byte-by-byte)\n", msgID)
}

func flagsDemo() {
	var flags uint16 = 0

	fmt.Println("Starting flags: 0b" + fmt.Sprintf("%016b", flags))

	// Set flags using OR
	flags |= FlagRequest
	fmt.Printf("After setting Request:   0b%016b\n", flags)

	flags |= FlagEncrypted
	fmt.Printf("After setting Encrypted: 0b%016b\n", flags)

	flags |= FlagCompressed
	fmt.Printf("After setting Compressed: 0b%016b\n", flags)

	// Check flags using AND
	fmt.Println()
	fmt.Printf("Is Request?    %v\n", flags&FlagRequest != 0)
	fmt.Printf("Is Error?      %v\n", flags&FlagError != 0)
	fmt.Printf("Is Encrypted?  %v\n", flags&FlagEncrypted != 0)
	fmt.Printf("Is Compressed? %v\n", flags&FlagCompressed != 0)

	// Clear flag using AND NOT
	flags &^= FlagCompressed
	fmt.Printf("\nAfter clearing Compressed: 0b%016b\n", flags)
	fmt.Printf("Is Compressed? %v\n", flags&FlagCompressed != 0)

	// Toggle flag using XOR
	flags ^= FlagError
	fmt.Printf("\nAfter toggling Error: 0b%016b\n", flags)
	fmt.Printf("Is Error? %v\n", flags&FlagError != 0)
}

func printHeader(h *Header) {
	fmt.Printf("  MessageID:     0x%04X\n", h.MessageID)
	fmt.Printf("  Flags:         0b%016b\n", h.Flags)
	fmt.Printf("    - Request:   %v\n", h.Flags&FlagRequest != 0)
	fmt.Printf("    - Error:     %v\n", h.Flags&FlagError != 0)
	fmt.Printf("    - Encrypted: %v\n", h.Flags&FlagEncrypted != 0)
	fmt.Printf("    - Compressed:%v\n", h.Flags&FlagCompressed != 0)
	fmt.Printf("  Sequence:      %d\n", h.Sequence)
	fmt.Printf("  Timestamp:     %d\n", h.Timestamp)
	fmt.Printf("  PayloadLength: %d\n", h.PayloadLength)
}

func hexDump(data []byte) {
	for i, b := range data {
		if i > 0 && i%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
}
