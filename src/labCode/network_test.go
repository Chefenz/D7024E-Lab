package labCode

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func networkSetup() (network Network) {
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
	routingTable := NewRoutingTable(NewContact(id, "master"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	return NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan, &findContCloseToValChan)
}

func TestNewNetwork(t *testing.T) {

	network := networkSetup()

	asserts := assert.New(t)

	asserts.NotNil(network)
}

func TestListen(t *testing.T) {

	listenContact := NewContact(NewRandomKademliaID(), ":8051")

	kademliaNode, _ := NewKademliaNode("master", "master:8051")

	testData := TransmitObj{Message: "ELSE", Data: "Nothing"}

	asserts := assert.New(t)

	go kademliaNode.Network.Listen("", 8051, *kademliaNode.StopChan)

	asserts.NotPanics(func() { kademliaNode.Network.sendMessage(&testData, &listenContact) }, "The code did panic")

	kademliaNode.StopAllRoutines()
}

func TestHandlePing(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)

	testData := TransmitObj{Message: "PING", Data: senderContact, RPC_created_at: time.Now()}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handlePing(transmitObject)

	//Check that buckets have been updated as expected

	contacts := kademliaNode.RoutingTable.FindClosestContacts(senderContact.ID, 1)
	fmt.Println(contacts)

	asserts := assert.New(t)

	asserts.Equal(senderContact.ID, contacts[0].ID)

	kademliaNode.StopAllRoutines()

}

func TestHandlePong(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)

	testData := TransmitObj{Message: "PING", Data: senderContact, RPC_created_at: time.Now()}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handlePong(transmitObject)

	//Check that buckets have been updated as expected

	contacts := kademliaNode.RoutingTable.FindClosestContacts(senderContact.ID, 1)

	asserts := assert.New(t)

	asserts.Equal(senderContact.ID, contacts[0].ID)

	kademliaNode.StopAllRoutines()

}

func TestHandleFindContactSameNode(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)
	go kademliaNode.RoutingTable.FindClosestContactsRoutine(*kademliaNode.StopChan)
	go kademliaNode.LookupContactRoutine(*kademliaNode.StopChan)

	testData := TransmitObj{
		Message: "FIND_CONTACT",
		Data: FindContactPayload{
			Sender: senderContact,
			Target: senderContact,
		},
		RPC_created_at: time.Now(),
	}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handleFindContact(transmitObject)

	//Check that buckets have been updated as expected

	contacts := kademliaNode.RoutingTable.FindClosestContacts(senderContact.ID, 1)

	asserts := assert.New(t)

	asserts.Equal(senderContact.ID, contacts[0].ID)

	kademliaNode.StopAllRoutines()

}

func TestHandleReturnFindContact(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)
	go kademliaNode.RoutingTable.FindClosestContactsRoutine(*kademliaNode.StopChan)
	go kademliaNode.LookupContactRoutine(*kademliaNode.StopChan)

	shortlist := []Contact{
		senderContact,
	}

	testData := TransmitObj{
		Message: "RETURN_FIND_CONTACT",
		Data: ReturnFindContactPayload{
			Shortlist: shortlist,
			Target:    senderContact,
		},
		RPC_created_at: time.Now(),
	}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handleReturnFindContact(transmitObject)

	//Check that buckets have been updated as expected

	asserts := assert.New(t)

	asserts.Equal(kademliaNode.Network.FoundTarget, true)

	kademliaNode.StopAllRoutines()

}

func TestHandleFindValueNoStoreValue(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.LookupCloseContactsToValueRoutine()
	go kademliaNode.DataStorageManager()
	testData := TransmitObj{
		Message: "FIND_VALUE",
		Data: FindValuePayload{
			Key: senderContact.ID,
		},
		RPC_created_at: time.Now(),
	}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handleFindValue(transmitObject)

	//Check that value can not be found since it has not been stored yet

	asserts := assert.New(t)

	asserts.Equal(kademliaNode.Network.FoundValue, false)

	kademliaNode.StopAllRoutines()

}

func TestHandleFindValueWithStoreValue(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.LookupCloseContactsToValueRoutine()
	go kademliaNode.DataStorageManager()
	testData := TransmitObj{
		Message: "FIND_VALUE",
		Data: FindValuePayload{
			Key: senderContact.ID,
		},
		RPC_created_at: time.Now(),
	}

	// First store value in datastorage

	key := senderContact.ID.String()
	data := []byte("test")

	newDataStorageObject := DataStorageObject{Data: data, Time: time.Now()}
	kademliaNode.DataStorage[key] = newDataStorageObject

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	transmitObject, _ = kademliaNode.Network.handleFindValue(transmitObject)

	//Check that no contacts are returned in shortlist since correct value has been found.

	asserts := assert.New(t)

	jsonData, err = json.Marshal(transmitObject)
	chk(err)
	transmitObject, err = kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	returnFindValuePayload := decodeTransmitObj(transmitObject, "ReturnFindValuePayload").(*ReturnFindValuePayload)

	asserts.Equal(len(returnFindValuePayload.Shortlist), 0)

	kademliaNode.StopAllRoutines()

}

func TestHandleReturnValue(t *testing.T) {
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	cliChan := make(chan string)
	kademliaNode.Network.CLIChan = &cliChan
	go func() {
		text := <-cliChan
		fmt.Println(text)
	}()
	go kademliaNode.LookupCloseContactsToValueRoutine()
	go kademliaNode.DataStorageManager()

	testData := TransmitObj{
		Message: "RETURN_FIND_VALUE",
		Sender:  kademliaNode.Network.Me,
		Data: ReturnFindValuePayload{
			Data:      "test",
			Shortlist: nil,
			TargetKey: nil,
		},
		RPC_created_at: time.Now(),
	}

	// First store value in datastorage

	key := kademliaNode.Network.Me.ID.String()
	data := []byte("test")

	newDataStorageObject := DataStorageObject{Data: data, Time: time.Now()}
	kademliaNode.DataStorage[key] = newDataStorageObject

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handleReturnFindValue(transmitObject)

	//Check that value has been found.

	asserts := assert.New(t)

	asserts.Equal(kademliaNode.Network.FoundValue, true)

	kademliaNode.StopAllRoutines()

}

func TestHandleStore(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.LookupCloseContactsToValueRoutine()
	go kademliaNode.DataStorageManager()
	testData := TransmitObj{
		Message: "STORE",
		Data: StorePayload{
			Key:  senderContact.ID,
			Data: "test",
		},
		RPC_created_at: time.Now(),
	}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)
	kademliaNode.Network.handleStore(transmitObject)

	//Check that no contacts are returned in shortlist since correct value has been found.

	asserts := assert.New(t)

	key := senderContact.ID.String()

	dataStorageObject := kademliaNode.DataStorage[key]

	asserts.Equal(string(dataStorageObject.Data), "test")

	kademliaNode.StopAllRoutines()

}

func TestHandleReturnStore(t *testing.T) {
	senderContact := NewContact(NewRandomKademliaID(), ":8053")
	kademliaNode, _ := NewKademliaNode("master", "master:8051")
	go kademliaNode.LookupCloseContactsToValueRoutine()
	go kademliaNode.DataStorageManager()
	testData := TransmitObj{
		Message: "RETURN_STORE",
		Sender:  senderContact,
		Data: ReturnStorePayload{
			Key: senderContact.ID,
		},
		RPC_created_at: time.Now(),
	}

	jsonData, err := json.Marshal(testData)
	chk(err)
	transmitObject, err := kademliaNode.Network.unmarshalTransmitObj(jsonData)
	chk(err)

	//Since it prints to a channel for CLI so do we only check if the code does not panic.

	asserts := assert.New(t)

	asserts.NotPanics(func() { kademliaNode.Network.handleReturnStore(transmitObject) }, "The code did panic")

	kademliaNode.StopAllRoutines()

}
