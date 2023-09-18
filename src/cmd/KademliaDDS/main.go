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
	node := labCode.NewKademliaNode(":8051")
	c := labCode.NewContact(labCode.NewRandomKademliaID(), ":8050")
	c.CalcDistance(node.RoutingTable.Me.ID)
	node.RoutingTable.Me.CalcDistance(c.ID)
	fmt.Println(node.RoutingTable.Me)
	fmt.Println(c)
	node.RoutingTable.AddContact(c)
	go node.Listen("", 8051)
	node.Ping(&c)
	for {

	}

}

func startMasterTypeNode() {
	node := labCode.NewKademliaNode(":8050")
	go node.Listen("", 8050)
	for {

	}
}
