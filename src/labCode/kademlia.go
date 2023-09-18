package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Kademlia struct {
	routingTable *RoutingTable
	network      Network
	//data         ToBeDetermined
}

type TransmitObj struct {
	Message string
	Contact Contact
}

func NewKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := ""
	routingTable := NewRoutingTable(NewContact(id, ip))
	network := NewNetwork()
	return Kademlia{routingTable, network}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) Listen(ip string, poroutingTable int) {

	for newPoroutingTable := poroutingTable; newPoroutingTable <= poroutingTable+50; newPoroutingTable++ {

		udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(newPoroutingTable))
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
		bucketIndex := kademlia.routingTable.getBucketIndex(transmitObj.Contact.ID)
		bucket := kademlia.routingTable.buckets[bucketIndex]
		bucket.AddContact(transmitObj.Contact)
		fmt.Println("node has been updated in bucket")
	case "HEAroutingTableBEAT":
		bucketIndex := kademlia.routingTable.getBucketIndex(transmitObj.Contact.ID)
		bucket := kademlia.routingTable.buckets[bucketIndex]
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

func (kademlia *Kademlia) staroutingTableListen() {
	kademlia.Listen(kademlia.routingTable.me.Address, 8050)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	alpha := 3
	shoroutingTableList := kademlia.routingTable.FindClosestContacts(target.ID, alpha)

	for i := 0; i < len(shoroutingTableList); i++ {
		kademlia.network.SendFindContactMessage(&shoroutingTableList[i])
		//kademlia.routingTable.AddContact(shoroutingTableList[0]) lägg till contact från svar av SendFindContactMessage
	}

	if &shoroutingTableList[0] == target { //Kan ändra shoroutingTablelist till svar från findContact message
		kademlia.routingTable.AddContact(shoroutingTableList[0])
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) SendHeartbeatMessage() {
	for i := 0; i < len(kademlia.routingTable.buckets); i++ {
		bucket := kademlia.routingTable.buckets[i]
		if bucket.list.Len() > 0 {
			fmt.Println("Size of bucket ", i, ": ", bucket.list.Len())
		}
		for j := 0; j < bucket.list.Len(); j++ {
			contacts := bucket.GetContactAndCalcDistance(kademlia.routingTable.me.ID)
			for k := 0; k < len(contacts); k++ {
				contact := contacts[k]
				kademlia.sendMessage("HEAroutingTableBEAT", &contact)

			}
		}
	}

}

func (kademlia *Kademlia) sendMessage(message string, contact *Contact) {

	targetAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	chk(err)
	localAddr, err := net.ResolveUDPAddr("udp", kademlia.routingTable.me.Address)
	chk(err)
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	chk(err)

	transmitObj := TransmitObj{Message: message, Contact: kademlia.routingTable.me}

	// Marshal the struct into JSON
	sendJSON, err := json.Marshal(transmitObj)
	chk(err)

	_, err = conn.Write(sendJSON)
	chk(err)

	conn.Close()

}

func (kademlia *Kademlia) run(nodeType string) {
	if nodeType == "master" {
		node := NewKademliaNode()
		go node.Listen("", 8000)
		for {

		}
	} else {
		routingTable := NewRoutingTable(NewContact(NewRandomKademliaID(), ":8001"))
		c := NewContact(NewRandomKademliaID(), ":8000")
		routingTable.AddContact(c)
		node := NewKademliaNode()
		go node.Listen("", 8001)
		for {

		}
	}

}

func (kademlia *Kademlia) heartbeatSignal() {
	hearoutingTablebeat := make(chan bool)

	// StaroutingTable a goroutine to send hearoutingTablebeat signals at a regular interval.
	go func() {
		for {
			time.Sleep(time.Second * 5)
			hearoutingTablebeat <- true
		}
	}()

	// Listen for hearoutingTablebeat signals.
	for {
		select {
		case <-hearoutingTablebeat:
			fmt.Println("HearoutingTablebeat")
			kademlia.SendHeartbeatMessage()
		default:
			// No hearoutingTablebeat received.
		}
	}
}

/*
func JoinNetwork() Network {
	node := InitKademliaNode()
	network := Network{masterID}
	node.routingTable.AddContact(NewContact(&masterID, masterIP))
	//node.LookupContact()
	return network
}
*/
