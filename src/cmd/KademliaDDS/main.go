package main

import (
	"fmt"
	"kademlia-app/labCode"
	"os"
)

func main() {
	containerName := os.Getenv("CONTAINER_NAME")

	var kademliaNode labCode.Kademlia
	if containerName == "master" {
		fmt.Println("Start of Masternode")
		kademliaNode = labCode.NewMasterKademliaNode()
		go kademliaNode.Listen("master", 8050)
		fmt.Println("End of node")
	} else {
		fmt.Println("Start of node")
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode = labCode.NewKademliaNode(nodeAddress + ":8050")
		fmt.Println("After Creation")

		masterNodeId := labCode.NewKademliaID("masterNode")
		masterNodeAddress := "master"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)
		fmt.Println("After master contact creation")

		kademliaNode.RoutingTable.AddContact(masterContact)
		fmt.Println("After addcontact Master")

		go kademliaNode.Listen(nodeAddress, 8050)
		fmt.Println("After Listen")

		kademliaNode.LookupContact(&kademliaNode.RoutingTable.Me)
		fmt.Println("After Lookup Contact")
	}

	labCode.RunCLI(kademliaNode)

}
