package labCode

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

type Network struct {
}

func Listen(ip string, port int) {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	portString := strconv.Itoa(port)
	l, err := net.Listen("tcp", ip+":"+portString)
	if err != nil {
		log.Fatal(err)
	}
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
				//fmt.Println("got", n, "bytes.")
				buf = append(buf, tmp[:n]...)

			}
			fmt.Println("total size:", len(buf))
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// Establish connection over tcp on port 8000
	conn, err := net.Dial("tcp", ":8000")

	if err != nil {
		// Handle the error, log it, or return an error message
		fmt.Printf("Failed to connect: %v", err)
		return
	}

	data := []byte("ping")
	_, writeErr := conn.Write(data)

	if writeErr != nil {
		// Handle the write error
		fmt.Printf("Failed to write data: %v", writeErr)
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
