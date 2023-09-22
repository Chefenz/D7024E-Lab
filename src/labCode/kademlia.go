package labCode

import (
	"fmt"
	"reflect"
	"time"
)

const alpha = 3
const dataDecayInterval = 7

type Kademlia struct {
	RoutingTable      *RoutingTable
	Network           Network
	DataStorage       map[string]DataStorageObject
	DataReadChan      *chan ReadOperation  //For sending read requests to the data storage
	DataWriteChan     *chan WriteOperation //For sending write requests to the data storage
	dataManagerTicker *time.Ticker
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

type FindValuePayload struct {
	Key *KademliaID
}

type ReturnFindValueDataPayload struct {
	Data string
}

type ReturnFindValueContactsPayload struct {
	TargetKey *KademliaID
	Data      []Contact
}

type StorePayload struct {
	Key  *KademliaID
	Data string
}

type ReturnStorePayload struct {
	Key *KademliaID
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

func NewKademliaNode(ip string) (Kademlia, *chan string) {
	id := NewRandomKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	CLIChan := make(chan string)
	dataManagerTicker := time.NewTicker(dataDecayInterval * time.Second)
	routingTable := NewRoutingTable(NewContact(id, ip), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan)
	return Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan, dataManagerTicker}, &CLIChan
}

func NewMasterKademliaNode() (Kademlia, *chan string) {
	id := NewMasterKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	CLIChan := make(chan string)
	dataManagerTicker := time.NewTicker(dataDecayInterval * time.Second)
	routingTable := NewRoutingTable(NewContact(id, "master"+":8051"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan)
	return Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan, dataManagerTicker}, &CLIChan
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
	dataKademliaID := NewKademliaID(hash)

	closestContactsLst := kademlia.RoutingTable.FindClosestContacts(dataKademliaID, alpha)

	findValuePayload := FindValuePayload{Key: dataKademliaID}
	transmitObj := TransmitObj{Message: "FIND_VALUE", Sender: kademlia.RoutingTable.Me, Data: findValuePayload}
	for i := 0; i < len(closestContactsLst); i++ {
		kademlia.Network.sendMessage(&transmitObj, &closestContactsLst[i])
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	strData := string(data)
	newDataId := NewKademliaDataID(strData)

	closestContactsLst := kademlia.RoutingTable.FindClosestContacts(newDataId, alpha)

	storePayload := StorePayload{Key: newDataId, Data: strData}
	transmitObj := TransmitObj{Message: "STORE", Sender: kademlia.RoutingTable.Me, Data: storePayload}

	if len(closestContactsLst) == 0 {
		*kademlia.Network.CLIChan <- ""
	}

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
			fmt.Println("In read case")
			dataStorageObject := Kademlia.DataStorage[read.Key]
			fmt.Println("Key used:", read.Key)
			fmt.Println("The datastorage object", dataStorageObject)
			if reflect.TypeOf(dataStorageObject) == nil {
				read.Resp <- nil
			} else {
				read.Resp <- dataStorageObject.Data
			}
		case write := <-*Kademlia.DataWriteChan:
			fmt.Println("in write case")
			key := write.Key
			data := write.Data

			newDataStorageObject := DataStorageObject{Data: data, Time: time.Now()}
			Kademlia.DataStorage[key] = newDataStorageObject

			fmt.Println("Map after store:", Kademlia.DataStorage)
			fmt.Println("lenght of map", len(Kademlia.DataStorage))

			write.Resp <- true

		case <-Kademlia.dataManagerTicker.C:
			for key, value := range Kademlia.DataStorage {
				insertedAt := value.Time
				durationSinceInsert := time.Since(insertedAt)

				//Delete all stored objects that has been stored for more than 1 minute
				if durationSinceInsert > time.Minute {
					delete(Kademlia.DataStorage, key)
					fmt.Println("DATA OBJECT DELETED BECAUSE OF DECAY")
					fmt.Println("lenght of map after deletion", len(Kademlia.DataStorage))

				}

			}

		default:
			//No write or read request has been issued
		}
	}
}
