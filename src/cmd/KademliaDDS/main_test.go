package main

import (
	"net"
	"testing"
)

func TestDefaultPing(t *testing.T) {
	conn, err := net.Dial("tcp", ":2000")

	if err != nil {
		// Handle the error, log it, or return an error message
		t.Fatalf("Failed to connect: %v", err)
		return
	}

	data := []byte("Hello world")
	_, writeErr := conn.Write(data)

	if writeErr != nil {
		// Handle the write error
		t.Fatalf("Failed to write data: %v", writeErr)
	}

	// Close the connection when done
	conn.Close()
}
