package labCode

import (
	"fmt"
	"reflect"
	"time"
)

const (
	alpha = 3 //The alpha value in the kademlia

	chkDataDecayinter = 25 * time.Second //The interval in which the node will check for data decay

	foundTargetDecayinter = 10 * time.Second

	dataDecayTime = 2 * time.Minute // How long the data will be stored in the node before it will be regarded as decayed.

)

type Kademlia struct {
	RoutingTable      *RoutingTable
	Network           Network
	DataStorage       map[string]DataStorageObject
	DataReadChan      *chan ReadOperation        //For sending read requests to the data storage
	DataWriteChan     *chan WriteOperation       //For sending write requests to the data storage
	FindConValueChan  *chan FindContCloseToValOp //For looking up contacts close to a target value
	dataManagerTicker *time.Ticker               //Periodically tells the node to check for decayed data
	StopChan          *chan string
}

type DataStorageObject struct {
	Data []byte
	Time time.Time
}

type TransmitObj struct {
	Message        string
	Sender         Contact
	Data           interface{}
	RPC_created_at time.Time
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

type ReturnFindValuePayload struct {
	Data      string
	Shortlist []Contact
	TargetKey *KademliaID
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

// Data needed to continue finding new nodes
type LookupContOp struct {
	Contact        *Contact
	RPC_created_at time.Time
}

//Struct for giving a request to look up nodes closest to a target

type FindContCloseToValOp struct {
	TargetID *KademliaID
	Resp     chan []Contact
}

func NewKademliaNode(ip string) (Kademlia, *chan string) {
	id := NewRandomKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan LookupContOp)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	CLIChan := make(chan string)
	findContCloseToValChan := make(chan FindContCloseToValOp)
	dataManagerTicker := time.NewTicker(chkDataDecayinter)
	stopChan := make(chan string)
	routingTable := NewRoutingTable(NewContact(id, ip), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan, &findContCloseToValChan)
	return Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan, &findContCloseToValChan, dataManagerTicker, &stopChan}, &CLIChan
}

func NewMasterKademliaNode() (Kademlia, *chan string) {
	id := NewMasterKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan LookupContOp)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	CLIChan := make(chan string)
	findContCloseToValChan := make(chan FindContCloseToValOp)
	dataManagerTicker := time.NewTicker(chkDataDecayinter)
	stopChan := make(chan string)
	routingTable := NewRoutingTable(NewContact(id, "master"+":8051"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan, &findContCloseToValChan)
	return Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan, &findContCloseToValChan, dataManagerTicker, &stopChan}, &CLIChan
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) Ping(contact *Contact) {
	transmitObj := TransmitObj{Message: "PING", Data: kademlia.RoutingTable.Me, RPC_created_at: time.Now()}
	kademlia.Network.sendMessage(&transmitObj, contact)
}

func (kademlia *Kademlia) LookupContact(target *Contact, RPC_created_at time.Time) {
	shortlist := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)

	for i := 0; i < len(shortlist); i++ {

		findContactPayload := FindContactPayload{Sender: kademlia.RoutingTable.Me, Target: *target}

		var transmitObj TransmitObj

		if RPC_created_at.IsZero() {
			transmitObj = TransmitObj{Message: "FIND_CONTACT", Data: findContactPayload, RPC_created_at: RPC_created_at}
		} else {
			transmitObj = TransmitObj{Message: "FIND_CONTACT", Data: findContactPayload, RPC_created_at: time.Now()}
		}

		kademlia.Network.sendMessage(&transmitObj, &shortlist[i])

	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	dataKademliaID := NewKademliaID(hash)

	closestContactsToTargetLst := kademlia.RoutingTable.FindClosestContacts(dataKademliaID, alpha)

	findValuePayload := FindValuePayload{Key: dataKademliaID}
	transmitObj := TransmitObj{Message: "FIND_VALUE", Sender: kademlia.RoutingTable.Me, Data: findValuePayload, RPC_created_at: time.Now()}
	for i := 0; i < len(closestContactsToTargetLst); i++ {
		kademlia.Network.sendMessage(&transmitObj, &closestContactsToTargetLst[i])
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	strData := string(data)
	newDataId := NewKademliaDataID(strData)

	closestContactsLst := kademlia.RoutingTable.FindClosestContacts(newDataId, alpha)

	storePayload := StorePayload{Key: newDataId, Data: strData}
	transmitObj := TransmitObj{Message: "STORE", Sender: kademlia.RoutingTable.Me, Data: storePayload, RPC_created_at: time.Now()}

	if len(closestContactsLst) == 0 {
		*kademlia.Network.CLIChan <- ""
	}

	for i := 0; i < len(closestContactsLst); i++ {
		kademlia.Network.sendMessage(&transmitObj, &closestContactsLst[i])
	}

}

func (kademlia *Kademlia) SendHeartbeatMessages() {
	for i := 0; i < len(kademlia.RoutingTable.buckets); i++ {
		bucket := kademlia.RoutingTable.buckets[i]
		if bucket.list.Len() > 0 {
			//fmt.Println("Size of bucket ", i, ": ", bucket.list.Len())
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
			time.Sleep(time.Second * 60)
			heartbeat <- true
		}
	}()

	// Listen for heartbeat signals.
	for {
		select {
		case <-heartbeat:
			fmt.Println("Heartbeat")
			//fmt.Println("Foundtarget: ", kademlia.Network.FoundTarget)
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
			lookupContOp := <-*kademlia.Network.LookupChan

			kademlia.LookupContact(lookupContOp.Contact, lookupContOp.RPC_created_at)
		}
	}
}

func (kademlia *Kademlia) LookupCloseContactsToValueRoutine() {
	for {
		findContactNearValueStruct := <-*kademlia.FindConValueChan
		targetID := findContactNearValueStruct.TargetID

		closestContactsLst := kademlia.RoutingTable.FindClosestContacts(targetID, alpha)

		findContactNearValueStruct.Resp <- closestContactsLst

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
		case <-Kademlia.dataManagerTicker.C:
			for key, value := range Kademlia.DataStorage {
				insertedAt := value.Time
				durationSinceInsert := time.Since(insertedAt)

				//Delete all stored objects that has been stored for more than the set decay time
				if durationSinceInsert > dataDecayTime {
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

func (kademlia *Kademlia) StopAllRoutines() {
	/*for i := 0; i < 4; i++ {
		*kademlia.StopChan <- "Stop"
	}*/
	close(*kademlia.StopChan)

}
