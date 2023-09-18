package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Kademlia struct {
	rt *RoutingTable
}

type TransmitObj struct {
	Message string
	Contact Contact
}

func InitKademliaNode(rt *RoutingTable) Kademlia {
	return Kademlia{rt}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) Listen(ip string, port int) {

	for newPort := port; newPort <= port+50; newPort++ {

		udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(newPort))
		chk(err)
		conn, err := net.ListenUDP("udp", udpAddr)
		chk(err)

		fmt.Println("Listening to: ", udpAddr)

		defer conn.Close()

		buffer := make([]byte, 4096)

		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from UDP connection:", err)
				continue
			}
			if len(buffer) > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])

				go kademlia.handleRPC(data, conn)
			}
		}
	}
}

func (kademlia *Kademlia) handleRPC(data []byte, conn *net.UDPConn) {

	var transmitObj TransmitObj

	err := json.Unmarshal(data, &transmitObj)
	chk(err)

	fmt.Println("Handling RPC: ", transmitObj.Message)
	targetAddr, err := net.ResolveUDPAddr("udp", transmitObj.Contact.Address)
	chk(err)

	switch transmitObj.Message {
	case "PING":
		fmt.Println("Returning message: ", "PONG")
		fmt.Println("Sending message to: ", targetAddr)

		kademlia.sendMessage("PONG", &transmitObj.Contact)
	case "PONG":
		fmt.Println("Received PONG response from", transmitObj.Contact.Address)
		bucketIndex := kademlia.rt.getBucketIndex(transmitObj.Contact.ID)
		bucket := kademlia.rt.buckets[bucketIndex]
		bucket.AddContact(transmitObj.Contact)
		fmt.Println("node has been updated in bucket")
	case "HEARTBEAT":
		bucketIndex := kademlia.rt.getBucketIndex(transmitObj.Contact.ID)
		bucket := kademlia.rt.buckets[bucketIndex]
		bucket.AddContact(transmitObj.Contact)
		fmt.Println("node has been updated in bucket")
	case "FINDCONTACT":
		fmt.Println("This should handle lookup")
	case "FINDDATA":
		fmt.Println("This should handle finddata")
	case "STORE":
		fmt.Println("This should handle store")
	}

}

func (kademlia *Kademlia) Ping(contact *Contact) {

	fmt.Println("sending ping to addr:", contact.Address)
	kademlia.sendMessage("PING", contact)
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
				kademlia.sendMessage("HEARTBEAT", &contact)

			}
		}
	}

}

func (kademlia *Kademlia) sendMessage(message string, contact *Contact) {

	targetAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	chk(err)
	localAddr, err := net.ResolveUDPAddr("udp", kademlia.rt.me.Address)
	chk(err)
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	chk(err)

	transmitObj := TransmitObj{Message: message, Contact: kademlia.rt.me}

	// Marshal the struct into JSON
	sendJSON, err := json.Marshal(transmitObj)
	chk(err)

	_, err = conn.Write(sendJSON)
	chk(err)

	conn.Close()

}

func (kademlia *Kademlia) run(nodeType string) {
	if nodeType == "master" {
		rt := NewRoutingTable(NewContact(NewRandomKademliaID(), ":8000"))
		node := InitKademliaNode(rt)
		go node.Listen("", 8000)
		for {

		}
	} else {
		rt := NewRoutingTable(NewContact(NewRandomKademliaID(), ":8001"))
		c := NewContact(NewRandomKademliaID(), ":8000")
		rt.AddContact(c)
		node := InitKademliaNode(rt)
		go node.Listen("", 8001)
		for {

		}
	}

}

func (kademlia *Kademlia) heartbeatSignal() {
	heartbeat := make(chan bool)

	// Start a goroutine to send heartbeat signals at a regular interval.
	go func() {
		for {
			time.Sleep(time.Second * 5)
			heartbeat <- true
		}
	}()

	// Listen for heartbeat signals.
	for {
		select {
		case <-heartbeat:
			fmt.Println("Heartbeat")
			kademlia.SendHeartbeatMessage()
		default:
			// No heartbeat received.
		}
	}
}
