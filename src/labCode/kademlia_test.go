package labCode

/*
// Test for sending ping and recieving pong as reponse
func testPing(t *testing.T) {

	masterContact := NewContact(NewRandomKademliaID(), ":8000")
	masterNode := NewKademliaNode("")
	go masterNode.run("master")

	otherContact := NewContact(NewRandomKademliaID(), ":8001")
	otherRt := NewRoutingTable(otherContact)
	otherRt.AddContact(masterContact)
	otherNode := NewKademliaNode("")
	go otherNode.run("other")

	asserts := assert.New(t)

	asserts.Panics(func() { otherNode.Ping(&masterContact) }, "The code did not panic")
}

func testHeartbeat(t *testing.T) {
	masterContact := NewContact(NewRandomKademliaID(), ":8000")
	masterNode := NewKademliaNode("")
	go masterNode.run("master")

	otherContact := NewContact(NewRandomKademliaID(), ":8001")
	otherRt := NewRoutingTable(otherContact)
	otherRt.AddContact(masterContact)
	otherNode := NewKademliaNode("")
	go otherNode.run("other")

	otherNode.heartbeatSignal()

	asserts := assert.New(t)

	asserts.Panics(func() { otherNode.heartbeatSignal() }, "The code did not panic")

}

*/
