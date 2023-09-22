package labCode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoutingTable(t *testing.T) {
	testContact1 := NewContact(NewKademliaID("TEST1"), "TEST_ADDRES1")

	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	routingTable := NewRoutingTable(testContact1, &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)

	var testBuckets [IDLength * 8]*bucket
	for i := 0; i < IDLength*8; i++ {
		testBuckets[i] = newBucket()
	}

	routingTable.buckets = testBuckets

	testRoutingTable := RoutingTable{Me: testContact1, buckets: testBuckets, BucketChan: &bucketChan, BucketWaitChan: &bucketWaitChan, FindChan: &findChan, ReturnFindChan: &returnFindChan}
	testTable := &testRoutingTable
	assert.Equal(t, routingTable, testTable)
}

func TestAddContactRoutingtable(t *testing.T) {
	testMeContact := NewContact(NewKademliaID("TEST1"), "TEST_ADDRES1")

	bucketChan := make(chan Contact, 1)
	bucketWaitChan := make(chan bool)
	findChan := make(chan Contact)
	returnFindChan := make(chan []Contact)
	routingTable := NewRoutingTable(testMeContact, &bucketChan, &bucketWaitChan, &findChan, &returnFindChan)

	addedContact := NewContact(NewKademliaID("ADDED"), "ADDED_ADDRES")

	routingTable.AddContact(addedContact)

	closestContact := routingTable.FindClosestContacts(addedContact.ID, 1)
	addedContact.Distance = closestContact[0].Distance
	//testMeContact.Distance = nil

	assert.Equal(t, closestContact[0], addedContact)

}

/*
func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}
*/
