package labCode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContact(t *testing.T) {
	testID := NewKademliaID("TEST")
	createdContact := NewContact(testID, "TEST_ADDRES")
	testContact := Contact{ID: testID, Address: "TEST_ADDRES", Distance: nil}

	assert.Equal(t, createdContact, testContact, "Should be equal")

}

func TestContactCalcDistance(t *testing.T) {
	testID1 := NewKademliaID("TEST1")
	testID2 := NewKademliaID("TEST2")
	createdContact1 := NewContact(testID1, "TEST_ADDRES1")
	createdContact2 := NewContact(testID2, "TEST_ADDRES2")

	createdContact1.CalcDistance(createdContact2.ID)
	createdContact2.CalcDistance(createdContact1.ID)

	assert.Equal(t, *createdContact1.Distance, *createdContact2.Distance, "Should be equal")

}

func TestLess(t *testing.T) {
	testContact1 := NewContact(NewKademliaID("TEST1"), "TEST_ADDRES1")
	testContact2 := NewContact(NewKademliaID("TEST2"), "TEST_ADDRES2")
	testContact3 := NewContact(NewKademliaID("TEST3"), "TEST_ADDRES1")

	var distance1 KademliaID
	myIDValue1 := [IDLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance1 = myIDValue1
	testContact1.Distance = &distance1

	var distance2 KademliaID
	myIDValue2 := [IDLength]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance2 = myIDValue2
	testContact2.Distance = &distance2

	var distance3 KademliaID
	myIDValue3 := [IDLength]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance3 = myIDValue3
	testContact3.Distance = &distance3

	trueBool1 := testContact1.Less(&testContact2)
	trueBool2 := testContact1.Less(&testContact3)
	trueBool3 := testContact2.Less(&testContact3)
	falseBool1 := testContact2.Less(&testContact1)
	falseBool2 := testContact3.Less(&testContact1)
	falseBool3 := testContact3.Less(&testContact2)
	falseBool4 := testContact3.Less(&testContact3)

	assert.True(t, trueBool1)
	assert.True(t, trueBool2)
	assert.True(t, trueBool3)
	assert.False(t, falseBool1)
	assert.False(t, falseBool2)
	assert.False(t, falseBool3)
	assert.False(t, falseBool4)

}

func TestString(t *testing.T) {
	var distance1 KademliaID
	myIDValue1 := [IDLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance1 = myIDValue1

	testContact1 := NewContact(&distance1, "TEST_ADDRES1")

	stringContact := testContact1.String()

	testContact := ("contact(\"0000000000000000000000000000000000000000\", \"TEST_ADDRES1\")")

	assert.Equal(t, testContact, stringContact)
}

func TestContactCandidates(t *testing.T) {
	// Create some test contacts
	testContact1 := NewContact(NewKademliaID("TEST1"), "TEST_ADDRES1")
	testContact2 := NewContact(NewKademliaID("TEST2"), "TEST_ADDRES2")
	testContact3 := NewContact(NewKademliaID("TEST3"), "TEST_ADDRES1")

	var distance1 KademliaID
	myIDValue1 := [IDLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance1 = myIDValue1
	testContact1.Distance = &distance1

	var distance2 KademliaID
	myIDValue2 := [IDLength]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance2 = myIDValue2
	testContact2.Distance = &distance2

	var distance3 KademliaID
	myIDValue3 := [IDLength]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	distance3 = myIDValue3
	testContact3.Distance = &distance3

	// Create ContactCandidates and append contacts
	testCandidates := ContactCandidates{}
	testCandidates.Append([]Contact{testContact3, testContact1, testContact2})

	// Ensure the length is correct
	assert.Equal(t, 3, testCandidates.Len())

	// Test sorting
	testCandidates.Sort()
	sortedContacts := testCandidates.GetContacts(3)

	// Check if the contacts are sorted by distance
	trueBool1 := sortedContacts[0].Less(&sortedContacts[1])
	trueBool2 := sortedContacts[0].Less(&sortedContacts[2])
	trueBool3 := sortedContacts[1].Less(&sortedContacts[2])
	falseBool1 := sortedContacts[2].Less(&sortedContacts[1])
	falseBool2 := sortedContacts[2].Less(&sortedContacts[0])
	falseBool3 := sortedContacts[1].Less(&sortedContacts[0])
	falseBool4 := sortedContacts[0].Less(&sortedContacts[0])

	assert.True(t, trueBool1)
	assert.True(t, trueBool2)
	assert.True(t, trueBool3)
	assert.False(t, falseBool1)
	assert.False(t, falseBool2)
	assert.False(t, falseBool3)
	assert.False(t, falseBool4)

	// Check if GetContacts returns the correct number of contacts
	getContacts := testCandidates.GetContacts(2)
	getContactsStruct := ContactCandidates{contacts: getContacts}
	assert.Equal(t, 2, getContactsStruct.Len())
}

/*

func TestContactString(t *testing.T) {
	kadID := NewKademliaID("test_id")
	contact := NewContact(kadID, "test_address")
	expectedString := `contact("test_id", "test_address")`

	if contact.String() != expectedString {
		t.Errorf("String() expected to return %s, but got %s", expectedString, contact.String())
	}
}
*/
