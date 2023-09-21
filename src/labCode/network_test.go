package labCode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func networkSetup() (network Network) {
	id := NewRandomKademliaID()
	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	lookupChan := make(chan Contact)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	routingTable := NewRoutingTable(NewContact(id, "master"), &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)
	return NewNetwork(routingTable.Me, &bucketChan, &bucketWaitChan, &lookupChan, &findChan, &returnFindChan)
}

func TestNewNetwork(t *testing.T) {

	network := networkSetup()

	asserts := assert.New(t)

	asserts.NotNil(network)
}

func TestListen(t *testing.T) {

	listenContact := NewContact(NewRandomKademliaID(), "localhost:8051")

	kademliaNode := NewMasterKademliaNode()

	testData := TransmitObj{Message: "ELSE", Data: "Nothing"}

	asserts := assert.New(t)

	go kademliaNode.Network.Listen("localhost", 8051, *kademliaNode.StopChan)

	asserts.NotPanics(func() { kademliaNode.Network.sendMessage(&testData, &listenContact) }, "The code did panic")

	kademliaNode.StopAllRoutines()
}

func TestHandleRPC(t *testing.T) {

	listenContact := NewContact(NewRandomKademliaID(), "localhost:8052")
	senderContact := NewContact(NewRandomKademliaID(), "localhost:8053")
	targetContact := NewContact(NewRandomKademliaID(), "localhost:8053")

	kademliaNode := NewMasterKademliaNode()

	go kademliaNode.Network.Listen("localhost", 8052, *kademliaNode.StopChan)
	go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)
	go kademliaNode.RoutingTable.FindClosestContactsRoutine(*kademliaNode.StopChan)
	go kademliaNode.LookupContactRoutine(*kademliaNode.StopChan)

	shortlist := []Contact{
		targetContact,
	}

	testData := []TransmitObj{
		{
			Message: "PING",
			Data:    senderContact,
		},
		{
			Message: "PONG",
			Data:    senderContact,
		},
		{
			Message: "HEARTBEAT",
			Data:    senderContact,
		},
		{
			Message: "FIND_CONTACT",
			Data: FindContactPayload{
				Sender: senderContact,
				Target: targetContact,
			},
		},
		{
			Message: "RETURN_FIND_CONTACT",
			Data: ReturnFindContactPayload{
				Shortlist: shortlist,
				Target:    targetContact,
			},
		},

		{
			Message: "FIND_DATA",
			Data:    senderContact,
		},
		{
			Message: "STORE",
			Data:    senderContact,
		},
	}

	asserts := assert.New(t)

	for i := 0; i < len(testData); i++ {

		asserts.NotPanics(func() { kademliaNode.Network.sendMessage(&testData[i], &listenContact) }, "The code did panic")
	}
	kademliaNode.StopAllRoutines()

}
