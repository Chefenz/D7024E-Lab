package main

import (
	"fmt"
	"kademlia-app/labCode"
	"os"
)

func main() {
	connType := os.Args[1]

	fmt.Println("Running node as:", connType)

	if connType == "master" {
		startMasterTypeNode()
	} else if connType == "other" {
		startOtherTypeNode()
	} else {
		fmt.Printf("please provide the arg (master or other)")
	}
}

func startOtherTypeNode() {
	n := labCode.Network{}
	c := labCode.NewContact(labCode.NewRandomKademliaID(), "")
	n.SendPingMessage(&c)
}

func startMasterTypeNode() {
	labCode.Listen("", 8000)
}
