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
	routingTable RoutingTable
	network      Network
	//data         ToBeDetermined
}

func (kademlia *Kademlia) startListen() {
	kademlia.network.Listen(kademlia.routingTable.me.Address, 8050)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	alpha := 3
	shortList := kademlia.routingTable.FindClosestContacts(target.ID, alpha)

	for i := 0; i < len(shortList); i++ {
		kademlia.network.SendFindContactMessage(&shortList[i])
		//kademlia.routingTable.AddContact(shortList[0]) lägg till contact från svar av SendFindContactMessage
	}

	if &shortList[0] == target { //Kan ändra shortlist till svar från findContact message
		kademlia.routingTable.AddContact(shortList[0])
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
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

func InitKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := ""
	rt := NewRoutingTable(NewContact(id, ip))
	network := NewNetwork()
	Listen(rt.me.Address, 8050)
	return Kademlia{*rt, network}
}
