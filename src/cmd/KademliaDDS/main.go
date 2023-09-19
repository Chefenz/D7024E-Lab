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
		go kademliaNode.Listen("master", 8051)
	} else {
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode = labCode.NewKademliaNode(nodeAddress + ":8051")

		masterNodeId := labCode.NewKademliaID("masterNode")
		masterNodeAddress := "master:8051"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)

		kademliaNode.RoutingTable.AddContact(masterContact)

		go kademliaNode.Listen(nodeAddress, 8051)

		kademliaNode.LookupContact(&kademliaNode.RoutingTable.Me)

		go kademliaNode.SendHeartbeatMessage()

	}

	labCode.RunCLI(kademliaNode)

}
