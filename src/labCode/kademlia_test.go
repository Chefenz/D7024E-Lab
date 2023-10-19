package labCode

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewKademliaNode(t *testing.T) {

	id := NewRandomKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan LookupContOp)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	stopChan := make(chan string)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	CLIChan := make(chan string)
	findContCloseToValChan := make(chan FindContCloseToValOp)
	dataManagerTicker := time.NewTicker(chkDataDecayinter)
	ToggleNonCLIPrintouts := false
	routingTable := NewRoutingTable(NewContact(id, ":8051"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan, &findContCloseToValChan, &ToggleNonCLIPrintouts)
	kademliaNode, _ := Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan, &findContCloseToValChan, dataManagerTicker, &stopChan, &sync.Mutex{}}, &CLIChan

	kademliaNode2, _ := NewKademliaNode(":8051")

	asserts := assert.New(t)

	asserts.NotEqual(kademliaNode, kademliaNode2)

}

func TestNewMasterKademliaNode(t *testing.T) {

	id := NewKademliaDataID("masterNode")
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan LookupContOp)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	stopChan := make(chan string)
	dataReadChan := make(chan ReadOperation)
	dataWriteChan := make(chan WriteOperation)
	CLIChan := make(chan string)
	findContCloseToValChan := make(chan FindContCloseToValOp)
	dataManagerTicker := time.NewTicker(chkDataDecayinter)
	ToggleNonCLIPrintouts := false
	routingTable := NewRoutingTable(NewContact(id, "master"+":8051"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	network := NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan, &dataReadChan, &dataWriteChan, &CLIChan, &findContCloseToValChan, &ToggleNonCLIPrintouts)
	kademliaNode, _ := Kademlia{routingTable, network, make(map[string]DataStorageObject), &dataReadChan, &dataWriteChan, &findContCloseToValChan, dataManagerTicker, &stopChan, &sync.Mutex{}}, &CLIChan

	kademliaNode2, _ := NewMasterKademliaNode()

	asserts := assert.New(t)

	asserts.NotEqual(kademliaNode, kademliaNode2)

}

// Test for sending ping and it not panicing
func TestPing(t *testing.T) {

	listenContact := NewContact(NewRandomKademliaID(), ":8053")

	masterNode, _ := NewMasterKademliaNode()
	otherNode, _ := NewKademliaNode("other")

	asserts := assert.New(t)

	go masterNode.Network.Listen("", 8053, *masterNode.StopChan)

	asserts.NotPanics(func() { otherNode.Ping(&listenContact) }, "The code did panic")

	masterNode.StopAllRoutines()
}

func TestLookupContact(t *testing.T) {

	targetContact := NewContact(NewRandomKademliaID(), ":8055")

	masterNode, _ := NewMasterKademliaNode()
	otherNode, _ := NewKademliaNode("other")

	asserts := assert.New(t)

	go masterNode.Network.Listen("", 8054, *masterNode.StopChan)

	asserts.NotPanics(func() { otherNode.LookupContact(&targetContact, time.Now()) }, "The code did panic")

	masterNode.StopAllRoutines()
}

/*
func TestLookupData(t *testing.T) {

	masterNode, _ := NewMasterKademliaNode()
	otherNode, _ := NewKademliaNode("other")

	asserts := assert.New(t)

	go masterNode.Network.Listen("", 8055, *masterNode.StopChan)

	asserts.NotPanics(func() { otherNode.LookupData("insert data here") }, "The code did panic")

	masterNode.StopAllRoutines()

}

func TestStore(t *testing.T) {

	masterNode, _ := NewMasterKademliaNode()
	otherNode, _ := NewKademliaNode("localhost")

	asserts := assert.New(t)

	go masterNode.Network.Listen("", 8056, *masterNode.StopChan)

	data := []byte("Hello")
	asserts.NotPanics(func() { otherNode.Store(data) }, "The code did panic")

	masterNode.StopAllRoutines()

}
*/

func TestHeartbeat(t *testing.T) {

	masterNode, _ := NewMasterKademliaNode()
	otherNode, _ := NewKademliaNode("other")

	asserts := assert.New(t)

	go masterNode.Network.Listen("", 8057, *masterNode.StopChan)

	asserts.NotPanics(func() { otherNode.SendHeartbeatMessages() }, "The code did panic")

	masterNode.StopAllRoutines()

}
