package main

import (
	"kademlia-app/labCode"
	"os"
)

func main() {
	//containerName := os.Getenv("CONTAINER_NAME")
	containerName := "master"

	var kademliaNode labCode.Kademlia
	if containerName == "master" {
		kademliaNode = labCode.NewMasterKademliaNode()
	} else {
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode = labCode.NewKademliaNode(nodeAddress)

		masterNodeId := labCode.NewKademliaID("masterNodeTest")
		masterNodeAddress := "master"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)

		kademliaNode.AddContact(masterContact)

		//TODO: Add find contact on myself
	}

	labCode.RunCLI(kademliaNode)

}
