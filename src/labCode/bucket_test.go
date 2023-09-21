package labCode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBucket(t *testing.T) {
	testbucket := newBucket()

	assert.NotNil(t, testbucket.list, "List should not be nil")

}

func TestAddContact(t *testing.T) {
	testBucket := newBucket()

	// Add a contact to the bucket
	contact1 := Contact{ID: NewKademliaID("TEST"), Address: "TEST_ADDRES"}
	testBucket.AddContact(contact1)

	// Check that the bucket contains the added contact
	assert.Equal(t, contact1, testBucket.list.Front().Value, "Contact should be at the front of the bucket")

	contact2 := Contact{ID: NewKademliaID("TEST2"), Address: "TEST_ADDRES2"}
	testBucket.AddContact(contact2)
	assert.Equal(t, contact2, testBucket.list.Front().Value, "Contact should be at the front of the bucket")

	testBucket.AddContact(contact1)
	assert.Equal(t, contact1, testBucket.list.Front().Value, "Contact should be at the front of the bucket")
}

func TestAddContactBucketFull(t *testing.T) {
	// Create a new bucket instance with a maximum bucket size of 1
	testBucket := newBucket()

	contact1 := Contact{ID: NewRandomKademliaID(), Address: "TEST_ADDRES"}
	testBucket.AddContact(contact1)

	assert.Equal(t, 1, testBucket.list.Len())

	for i := 0; i < 25; i++ {
		loopcontact := Contact{ID: NewRandomKademliaID(), Address: "LOOP_ADDRES"}
		testBucket.AddContact(loopcontact)
	}

	testBucket.AddContact(contact1)

	// Check that the second contact is at the front of the bucket,
	// as it should have replaced the first contact due to bucket size limitation
	assert.Equal(t, contact1, testBucket.list.Front().Value, "First contact should be at the front of the bucket")

	bucketLenght := testBucket.Len()
	assert.Equal(t, bucketLenght, 20)
}

func TestGetContactAndCalcDistance(t *testing.T) {
	testBucket := newBucket()

	// Add some sample contacts to the bucket
	contact1 := Contact{ID: NewKademliaID("TEST"), Address: "TEST_ADDRES"}
	testBucket.AddContact(contact1)
	contact2 := Contact{ID: NewKademliaID("TEST2"), Address: "TEST_ADDRES2"}
	testBucket.AddContact(contact2)

	// Create a target KademliaID
	target := NewKademliaID("TARGET")

	// Call the GetContactAndCalcDistance method
	contacts := testBucket.GetContactAndCalcDistance(target)

	// Check that the contacts were retrieved and their distances were calculated
	assert.Len(t, contacts, 2, "Expected 2 contacts in the result")
}
