package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

const alpha = 3

type Kademlia struct {
	RoutingTable *RoutingTable
	network      Network
	data         map[KademliaID][]byte
}

type TransmitObj struct {
	Message string
	Data    interface{}
}

type FindContactPayload struct {
	Sender Contact
	Target Contact
}

type ReturnFindContactPayload struct {
	Shortlist []Contact
	Target    Contact
}

func NewKademliaNode(ip string) Kademlia {
	id := NewRandomKademliaID()
	routingTable := NewRoutingTable(NewContact(id, ip))
	network := NewNetwork()
	return Kademlia{routingTable, network, make(map[KademliaID][]byte)}
}

func NewMasterKademliaNode() Kademlia {
	id := NewKademliaID("masterNode")
	routingTable := NewRoutingTable(NewContact(id, "master"+":8050"))
	network := NewNetwork()
	return Kademlia{routingTable, network, map[KademliaID][]byte{}}
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

	switch transmitObj.Message {
	case "PING":
		fmt.Println("Returning message: ", "PONG")

		Contact := DataToContact(transmitObj)
		transmitObj := TransmitObj{Message: "PONG", Data: kademlia.RoutingTable.Me}
		kademlia.sendMessage(&transmitObj, Contact)
	case "PONG":

		Contact := DataToContact(transmitObj)

		fmt.Println("Received PONG response from", Contact.Address)
		bucketIndex := kademlia.RoutingTable.getBucketIndex(Contact.ID)
		bucket := kademlia.RoutingTable.buckets[bucketIndex]
		bucket.AddContact(*Contact)
		fmt.Println("node has been updated in bucket")
	case "HEARTBEAT":

		Contact := DataToContact(transmitObj)

		bucketIndex := kademlia.RoutingTable.getBucketIndex(Contact.ID)
		bucket := kademlia.RoutingTable.buckets[bucketIndex]
		bucket.AddContact(*Contact)
		fmt.Println("node has been updated in bucket")
	case "FIND_CONTACT":

		findContactPayloadMap, ok := transmitObj.Data.(map[string]interface{})

		if ok != true {
			fmt.Println("Data is not a Map")
		}

		//fmt.Println(findContactPayloadMap)

		var findContactPayload *FindContactPayload
		err := mapstructure.Decode(findContactPayloadMap, &findContactPayload)
		chk(err)

		//fmt.Println(findContactPayload)

		kademlia.RoutingTable.AddContact(findContactPayload.Sender)
		closestContacts := kademlia.RoutingTable.FindClosestContacts(findContactPayload.Target.ID, alpha)

		returnFindContactPayload := ReturnFindContactPayload{Shortlist: closestContacts, Target: findContactPayload.Target}

		transmitObj := TransmitObj{Message: "RETURN_FIND_CONTACT", Data: returnFindContactPayload}

		kademlia.sendMessage(&transmitObj, &findContactPayload.Sender)

	case "RETURN_FIND_CONTACT":

		returnFindContactPayloadMap, ok := transmitObj.Data.(map[string]interface{})

		if ok != true {
			fmt.Println("Data is not a Map")
		}

		var returnFindContactPayload *ReturnFindContactPayload
		err := mapstructure.Decode(returnFindContactPayloadMap, &returnFindContactPayload)
		chk(err)

		shortlist := returnFindContactPayload.Shortlist
		target := returnFindContactPayload.Target
		foundTarget := false

		for i := 0; i < len(shortlist); i++ {
			kademlia.RoutingTable.AddContact(shortlist[0])
			if *shortlist[i].ID == *target.ID {
				foundTarget = true
				fmt.Println("Found The Target Node :)")
			}
		}
		if foundTarget == false {
			fmt.Println("Did Not Find The Target Node Will Try Again")
			kademlia.LookupContact(&target)

		}

	case "FIND_DATA":
		fmt.Println("This should handle finddata")
	case "STORE":
		fmt.Println("This should handle store")
	}

}

func DataToContact(obj TransmitObj) (contact *Contact) {
	contactMap, ok := obj.Data.(map[string]interface{})

	if ok != true {
		fmt.Println("Data is not a Map")
	}

	err := mapstructure.Decode(contactMap, &contact)
	chk(err)

	return contact
}

func (kademlia *Kademlia) Ping(contact *Contact) {

	//fmt.Println("sending ping to addr:", contact.Address)
	transmitObj := TransmitObj{Message: "PING", Data: kademlia.RoutingTable.Me}
	kademlia.sendMessage(&transmitObj, contact)
}

func (kademlia *Kademlia) startListen() {
	kademlia.Listen(kademlia.RoutingTable.Me.Address, 8050)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	shortList := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)

	for i := 0; i < len(shortList); i++ {

		findContactPayload := FindContactPayload{Sender: kademlia.RoutingTable.Me, Target: *target}

		transmitObj := TransmitObj{Message: "FIND_CONTACT", Data: findContactPayload}

		kademlia.sendMessage(&transmitObj, &shortList[i])
		//kademlia.routingTable.AddContact(shortList[0]) lägg till contact från svar av SendFindContactMessage
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) SendHeartbeatMessage() {
	for i := 0; i < len(kademlia.RoutingTable.buckets); i++ {
		bucket := kademlia.RoutingTable.buckets[i]
		if bucket.list.Len() > 0 {
			fmt.Println("Size of bucket ", i, ": ", bucket.list.Len())
		}
		for j := 0; j < bucket.list.Len(); j++ {
			contacts := bucket.GetContactAndCalcDistance(kademlia.RoutingTable.Me.ID)
			for k := 0; k < len(contacts); k++ {
				contact := contacts[k]
				transmitObj := TransmitObj{Message: "HEARTBEAT", Data: kademlia.RoutingTable.Me}
				kademlia.sendMessage(&transmitObj, &contact)

			}
		}
	}

}

func (kademlia *Kademlia) sendMessage(transmitObj *TransmitObj, contact *Contact) {

	targetAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	chk(err)

	fmt.Println("Target Address: ", targetAddr)
	fmt.Println(contact)

	localAddr, err := net.ResolveUDPAddr("udp", kademlia.RoutingTable.Me.Address)
	chk(err)
	fmt.Println("First error place")
	fmt.Println("Target Address: ", targetAddr)
	fmt.Println(contact)
	fmt.Println(localAddr)
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	chk(err)
	fmt.Println("second error place")
	fmt.Println("Target Address: ", targetAddr)
	fmt.Println(contact)
	fmt.Println(localAddr)

	// Marshal the struct into JSON
	sendJSON, err := json.Marshal(transmitObj)
	chk(err)

	_, err = conn.Write(sendJSON)
	chk(err)

	conn.Close()

}

func (kademlia *Kademlia) Run(nodeType string) {
	if nodeType == "master" {
		node := NewKademliaNode(":8050")
		go node.Listen("", 8050)
		for {

		}
	} else {
		node := NewKademliaNode(":8051")
		go node.Listen("", 8051)
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

/*
func (kademlia *Kademlia) JoinNetwork() {
	node := NewKademliaNode()
	network := Network{}

	masterNodeId := NewKademliaID("masterNode")
	masterNodeAddress := "master"
	masterContact := NewContact(masterNodeId, masterNodeAddress)

	node.RoutingTable.AddContact(NewContact(&masterID, masterIP))
	//node.LookupContact()
	kademlia.network = network
}

*/
