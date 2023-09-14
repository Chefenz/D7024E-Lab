package main

import (
	"fmt"
	"kademlia-app/labCode"
	"os"
)

func main() {
	connType := os.Args[1]

	fmt.Println("Running node as:", connType)

	//newKademliaID := labCode.NewRandomKademliaID()
	//me := labCode.NewContact(newKademliaID, "localhost") //Localhost for now
	//routingTable := labCode.NewRoutingTable(me)
	//fmt.Println(routingTable)

	if connType == "master" {
		go startMasterTypeNode()
	} else if connType == "other" {
		go startOtherTypeNode()
		println("test")
	} else {
		fmt.Printf("please provide the arg (master or other)")
	}

	//Start the command line user face
	labCode.RunCLI()
}

func startOtherTypeNode() {
	n := labCode.Network{}
	c := labCode.NewContact(labCode.NewRandomKademliaID(), "")
	n.SendPingMessage(&c)
}

func startMasterTypeNode() {
	labCode.Listen("", 8000)
}
