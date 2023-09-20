package labCode

import (
	"fmt"
	"time"
)

const alpha = 3

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      Network
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
	bucketChan := make(chan Contact, 1)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	routingTable := NewRoutingTable(NewContact(id, ip), &bucketChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &lookupChan, &findChan, &returnFindChan)
	fmt.Println(bucketChan)
	fmt.Println(lookupChan)
	fmt.Println(findChan)
	fmt.Println(returnFindChan)
	fmt.Println(routingTable)
	fmt.Println(network)
	return Kademlia{routingTable, network, make(map[KademliaID][]byte)}
}

func NewMasterKademliaNode() Kademlia {
	id := NewKademliaID("masterNode")
	bucketChan := make(chan Contact, 1)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	routingTable := NewRoutingTable(NewContact(id, "master"+":8051"), &bucketChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &lookupChan, &findChan, &returnFindChan)
	return Kademlia{routingTable, network, make(map[KademliaID][]byte)}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) Ping(contact *Contact) {
	//fmt.Println("sending ping to addr:", contact.Address)
	transmitObj := TransmitObj{Message: "PING", Data: kademlia.RoutingTable.Me}
	kademlia.Network.sendMessage(&transmitObj, contact)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	shortlist := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)

	for i := 0; i < len(shortlist); i++ {

		findContactPayload := FindContactPayload{Sender: kademlia.RoutingTable.Me, Target: *target}

		transmitObj := TransmitObj{Message: "FIND_CONTACT", Data: findContactPayload}

		if *shortlist[i].ID == *kademlia.RoutingTable.Me.ID {
			fmt.Println("in found myself in shortlist lookup contact")
		}
		fmt.Println("Outside if in lookup contact")
		kademlia.Network.sendMessage(&transmitObj, &shortlist[i])

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
				kademlia.Network.sendMessage(&transmitObj, &contact)

			}
		}
	}

}

func (kademlia *Kademlia) HeartbeatSignal() {
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

func (kademlia *Kademlia) LookupContactRoutine() {
	target := <-*kademlia.Network.LookupChan

	kademlia.LookupContact(&target)
}
