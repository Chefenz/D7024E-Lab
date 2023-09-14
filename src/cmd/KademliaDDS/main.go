package main

import (
	"fmt"
	"kademlia-app/labCode"
	//"os"
)

func main() {

	fmt.Println("Network Join Test")
	masterID := startMasterTypeNode()
	labCode.JoinNetwork(masterID, "MasterNode")

	/*
		connType := os.Args[1]

		fmt.Println("Running node as:", connType)

		if connType == "master" {
			startMasterTypeNode()
		} else if connType == "other" {
			startOtherTypeNode()
		} else {
			fmt.Printf("please provide the arg (master or other)")
		}
	*/
}

func startOtherTypeNode() {
	n := labCode.Network{}
	c := labCode.NewContact(labCode.NewRandomKademliaID(), "")
	labCode.NewRoutingTable(c)
	//labCode.AddContact()

	//labCode.lookupContact(self)
	n.SendPingMessage(&c)
}

func startMasterTypeNode() *labCode.KademliaID {
	labCode.Listen("MasterNode", 8000)
	masterID := labCode.NewRandomKademliaID()
	masterContact := labCode.NewContact(masterID, "MasterNode")
	//masterContact := labCode.NewContact(labCode.NewRandomKademliaID(), "MasterNode")
	labCode.NewRoutingTable(masterContact)
	return (masterID)
}
