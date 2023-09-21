package labCode

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

const alpha = 3

type Kademlia struct {
	RoutingTable  *RoutingTable
	Network       Network
	DataStorage   map[string]DataStorageObject
	DataReadChan  *chan ReadOperation  //For sending read requests to the data storage
	DataWriteChan *chan WriteOperation //For sending write requests to the data storage
}

type DataStorageObject struct {
	Data []byte
	Time time.Time
}

type TransmitObj struct {
	Message string
	Sender  Contact
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

type StorePayload struct {
	Key  *KademliaID
	Wg   *sync.WaitGroup
	Data []byte
}

type ReturnStorePayload struct {
	Key *KademliaID
	Wg  *sync.WaitGroup
}

// Structs for read and write operations (move these to an appropriate place later)
type ReadOperation struct {
	Key  string
	Resp chan []byte
}

type WriteOperation struct {
	Key  string
	Data []byte
	Resp chan bool
}

func NewKademliaNode(ip string) Kademlia {
	id := NewRandomKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	routingTable := NewRoutingTable(NewContact(id, ip), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan)
	return Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan}
}

func NewMasterKademliaNode() Kademlia {
	id := NewMasterKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	routingTable := NewRoutingTable(NewContact(id, "master"+":8051"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan)
	return Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan}
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

		kademlia.Network.sendMessage(&transmitObj, &shortlist[i])

	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	fmt.Println("In store")
	strData := string(data)
	newDataId := NewKademliaDataID(strData)

	closestContactsLst := kademlia.RoutingTable.FindClosestContacts(newDataId, alpha)
	var wg sync.WaitGroup
	wg.Add(len(closestContactsLst))

	storePayload := StorePayload{Key: newDataId, Wg: &wg, Data: data}
	fmt.Println("The StorePayload:", storePayload)
	fmt.Println("The data in the Data field:", data)
	fmt.Println("The type of the data:", reflect.TypeOf(data))
	transmitObj := TransmitObj{Message: "STORE", Sender: kademlia.RoutingTable.Me, Data: storePayload}

	for i := 0; i < len(closestContactsLst); i++ {
		kademlia.Network.sendMessage(&transmitObj, &closestContactsLst[i])
	}
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
			time.Sleep(time.Second * 30)
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
	for {
		target := <-*kademlia.Network.LookupChan

		kademlia.LookupContact(&target)
	}
}

func (Kademlia *Kademlia) DataStorageManager() {
	for {
		select {
		case read := <-*Kademlia.DataReadChan:
			dataStorageObject := Kademlia.DataStorage[read.Key]
			if reflect.TypeOf(dataStorageObject) == nil {
				read.Resp <- nil
			} else {
				read.Resp <- dataStorageObject.Data
			}
		case write := <-*Kademlia.DataWriteChan:
			key := write.Key
			data := write.Data

			newDataStorageObject := DataStorageObject{Data: data, Time: time.Now()}
			Kademlia.DataStorage[key] = newDataStorageObject

			write.Resp <- true
		default:
			//No write or read request has been issued
			continue
		}
	}
}
