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
	StopChan     *chan string
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
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	stopChan := make(chan string)
	routingTable := NewRoutingTable(NewContact(id, ip), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan)
	return Kademlia{routingTable, network, make(map[KademliaID][]byte), &stopChan}
}

func NewMasterKademliaNode() Kademlia {
	id := NewKademliaID("masterNode")
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	stopChan := make(chan string)
	routingTable := NewRoutingTable(NewContact(id, "master"+":8051"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan)
	return Kademlia{routingTable, network, make(map[KademliaID][]byte), &stopChan}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) Ping(contact *Contact) {
	transmitObj := TransmitObj{Message: "PING", Data: kademlia.RoutingTable.Me}
	kademlia.Network.sendMessage(&transmitObj, contact)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	shortlist := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)

	for i := 0; i < len(shortlist); i++ {

		findContactPayload := FindContactPayload{Sender: kademlia.RoutingTable.Me, Target: *target}

		transmitObj := TransmitObj{Message: "FIND_CONTACT", Data: findContactPayload}

		kademlia.Network.sendMessage(&transmitObj, &shortlist[i])

	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) SendHeartbeatMessages() {
	for i := 0; i < len(kademlia.RoutingTable.buckets); i++ {
		bucket := kademlia.RoutingTable.buckets[i]
		if bucket.list.Len() > 0 {
			fmt.Println("Size of bucket ", i, ": ", bucket.list.Len())
		}
		for j := 0; j < bucket.list.Len(); j++ {
			contacts := bucket.GetContactAndCalcDistance(kademlia.RoutingTable.Me.ID)
			for k := 0; k < len(contacts); k++ {
				contact := contacts[k]
				kademlia.Ping(&contact)
			}
		}
	}
}

func (kademlia *Kademlia) HeartbeatSignal(stopChan <-chan string) {
	heartbeat := make(chan bool)

	// Start a goroutine to send heartbeat signals at a regular interval.
	go func() {
		for {
			time.Sleep(time.Second * 30)
			heartbeat <- true
		}
	}()

	// Listen for heartbeat signals.
	for {
		select {
		case <-heartbeat:
			fmt.Println("Heartbeat")
			kademlia.SendHeartbeatMessages()
		case <-stopChan:

		default:
			// No heartbeat received.
		}
	}
}

func (kademlia *Kademlia) LookupContactRoutine(stopChan <-chan string) {
	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping look up contact routine...")
			return
		default:
			target := <-*kademlia.Network.LookupChan

			kademlia.LookupContact(&target)
		}
	}
}

func (kademlia *Kademlia) StopAllRoutines() {
	/*for i := 0; i < 4; i++ {
		*kademlia.StopChan <- "Stop"
	}*/
	close(*kademlia.StopChan)

}
