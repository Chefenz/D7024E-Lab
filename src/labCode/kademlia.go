package labCode

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Kademlia struct {
	rt *RoutingTable
}

func InitKademliaNode(rt *RoutingTable) Kademlia {
	return Kademlia{rt}
}

func (kademlia *Kademlia) Listen(ip string, port int) {
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

					var recievedObj struct {
						message string
						contact Contact
					}
					fmt.Println(len(buf))
					fmt.Println(len(tmp))
					buf = append(buf, tmp[:n]...)

					error := json.Unmarshal(buf[:n], &recievedObj)
					if error != nil {
						fmt.Println("failed to deserialize")
						fmt.Println(error)
						return
					}
					//buf = append(buf, tmp[:n]...)

					fmt.Println("Handling RPC: ", recievedObj)
					switch recievedObj.message {
					case "PING":
						fmt.Println("Returning message: ", "PONG")
						c.Write([]byte("PONG"))
					case "HEARTBEAT":
						bucketIndex := kademlia.rt.getBucketIndex(recievedObj.contact.ID)
						bucket := kademlia.rt.buckets[bucketIndex]
						bucket.AddContact(*&recievedObj.contact)
						fmt.Println("node with id: %v has been update in bucket", recievedObj.contact.ID)
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

func (kademlia *Kademlia) Ping(contact *Contact) {
	message, conn := kademlia.sendMessage("PING", contact)

	kademlia.recieveMessage(message, contact, conn)

	// Close the connection when done
	conn.Close()
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) SendHeartbeatMessage() {
	for i := 0; i < len(kademlia.rt.buckets); i++ {
		bucket := kademlia.rt.buckets[i]
		if bucket.list.Len() > 0 {
			fmt.Println("Size of bucket ", i, ": ", bucket.list.Len())
		}
		for j := 0; j < bucket.list.Len(); j++ {
			contacts := bucket.GetContactAndCalcDistance(kademlia.rt.me.ID)
			for k := 0; k < len(contacts); k++ {
				contact := contacts[k]
				_, conn := kademlia.sendMessage("HEARTBEAT", &contact)
				// Close the connection when done
				conn.Close()

			}
		}
	}

}

func (kademlia *Kademlia) sendMessage(message string, contact *Contact) (string, net.Conn) {
	// Establish connection over tcp on port 8000
	targetAddr := contact.Address
	conn, err := net.Dial("tcp", targetAddr)

	if err != nil {
		// Handle the error, log it, or return an error message
		fmt.Printf("Failed to connect: %v", err)
		return "", conn
	}

	fmt.Println("Connected to: ", targetAddr)

	sendObj := struct {
		Message string   `json:"message"`
		Contact *Contact `json:"contact"`
	}{
		Message: message,
		Contact: contact,
	}

	// Marshal the struct into JSON
	sendJSON, err := json.Marshal(sendObj)

	buf := make([]byte, 0, 4096) // big buffer
	fmt.Println(json.Unmarshal(buf, sendJSON))
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v", err)
		return "", conn
	}

	//aBytes := []byte(fmt.Sprintf("%v", message))
	//bBytes := []byte(fmt.Sprintf("%v", *contact))

	//data := make([]byte, len(aBytes)+len(bBytes))
	//copy(data, aBytes)
	//copy(data[len(aBytes):], bBytes)
	_, writeErr := conn.Write(sendJSON)

	if writeErr != nil {
		// Handle the write error
		fmt.Printf("Failed to write data: %v", writeErr)
	}

	return message, conn

}

func (kademlia *Kademlia) recieveMessage(recievedMessage string, contact *Contact, conn net.Conn) {
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
		fmt.Println("Received PONG response from", contact.Address)
		bucketIndex := kademlia.rt.getBucketIndex(contact.ID)
		bucket := kademlia.rt.buckets[bucketIndex]
		bucket.AddContact(*contact)

	} else {
		fmt.Println("Received an unexpected response:", response)
	}

	// Close the connection when done
	conn.Close()
}
