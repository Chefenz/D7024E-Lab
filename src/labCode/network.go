package labCode

import (
	"fmt"
	"net"
	"os"
)

type Network struct {
}

func NewNetwork() Network {
	return Network{}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// Establish connection over tcp on port 8000

	targetAddr := ":8000"
	conn, err := net.Dial("tcp", targetAddr)

	if err != nil {
		// Handle the error, log it, or return an error message
		fmt.Printf("Failed to connect: %v", err)
		return
	}
	fmt.Println("Connected to: ", targetAddr)

	data := []byte("PING")
	_, writeErr := conn.Write(data)

	if writeErr != nil {
		// Handle the write error
		fmt.Printf("Failed to write data: %v", writeErr)
	}
	fmt.Println("Sent message ", string(data))

	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	fmt.Println("Waiting for response...")

	// Wait for a response
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		os.Exit(1)
	}

	fmt.Println("Response recieved")

	// Received data from the target node, convert to string
	response := string(buffer[:n])

	// Check if the response is "PONG"
	if response == "PONG" {
		fmt.Println("Received PONG response from", targetAddr)
	} else {
		fmt.Println("Received an unexpected response:", response)
	}

	// Close the connection when done
	conn.Close()
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
