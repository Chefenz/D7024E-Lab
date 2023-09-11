package main

import (
	"fmt"
	"kademlia-app/labCode"
	"os"
)

func main() {
	connType := os.Getenv("TYPE")

	if connType == "master" {
		startOtherTypeNode()
	} else if connType == "other" {
		startMasterTypeNode()
	} else {
		fmt.Printf("please provide the TYPE arg (TYPE=master or TYPE=other)")
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
