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

	for newPort := port; newPort <= port+50; newPort++ {
		portString := strconv.Itoa(newPort)
		targetAddr := ip + ":" + portString
		l, err := net.Listen("tcp", targetAddr)
		if err != nil {
			continue
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
				fmt.Println("Connection from", conn.RemoteAddr())
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
					buf = append(buf, tmp[:n]...)

					fmt.Println("Handling RPC: ", string(buf))
					switch string(buf) {
					case "PING":
						fmt.Println("Returning message: ", "PONG")
						c.Write([]byte("PONG"))
					case "HEARTBEAT":

					case "FINDCONTACT":
						fmt.Println("This should handle lookup")
					case "FINDDATA":
						fmt.Println("This should handle finddata")
					case "STORE":
						fmt.Println("This should handle store")
					}
				}
			}(conn)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
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
