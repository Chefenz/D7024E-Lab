package labCode

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Network struct {
	masterID KademliaID
}

func Listen(ip string, port int) {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	portString := strconv.Itoa(port)
	targetAddr := ip + ":" + portString

	l, err := net.Listen("tcp", targetAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening to: ", targetAddr)
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			fmt.Println("Connection established")
			// Echo all incoming data.
			buf := make([]byte, 0, 4096) // big buffer
			tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
			for {
				n, err := conn.Read(tmp)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error:", err)
					}
					break
				}
				fmt.Println("got", n, "bytes.")
				buf = append(buf, tmp[:n]...)

				fmt.Println("Recieved message: ", string(buf))
				fmt.Println("Returning message: ", "PONG")
				c.Write([]byte("PONG"))
				fmt.Println("Message sent")

			}

			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
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
