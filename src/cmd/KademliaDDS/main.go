package main

import (
	"kademlia-app/labCode"
)

func main() {

	var c = labCode.NewContact(labCode.NewRandomKademliaID(), "localhost")
	var network = labCode.Network{}
	network.SendPingMessage(&c)

}
