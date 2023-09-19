package main

import (
	"kademlia-app/labCode"
	"os"
)

func main() {
	containerName := os.Getenv("CONTAINER_NAME")

	var kademliaNode labCode.Kademlia
	if containerName == "master" {
		kademliaNode = labCode.NewMasterKademliaNode()
		go kademliaNode.Listen("master", 8050)
	} else {
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode = labCode.NewKademliaNode(nodeAddress + ":8050")

		masterNodeId := labCode.NewKademliaID("masterNode")
		masterNodeAddress := "master"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)

		kademliaNode.RoutingTable.AddContact(masterContact)

		go kademliaNode.Listen(nodeAddress, 8050)

		kademliaNode.LookupContact(&kademliaNode.RoutingTable.Me)
	}

	labCode.RunCLI(kademliaNode)

}
