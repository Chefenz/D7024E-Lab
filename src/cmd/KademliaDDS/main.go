package main

import (
	"kademlia-app/labCode"
	"os"
	"time"
)

func main() {
	containerName := os.Getenv("CONTAINER_NAME")

	var kademliaNode labCode.Kademlia
	if containerName == "master" {
		kademliaNode = labCode.NewMasterKademliaNode()
		time.Sleep(time.Second * 5)
		go kademliaNode.RoutingTable.UpdateBucketRoutine()
		go kademliaNode.RoutingTable.FindClosestContactsRoutine()
		go kademliaNode.LookupContactRoutine()
		go kademliaNode.Network.Listen("master", 8051)
	} else {
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode = labCode.NewKademliaNode(nodeAddress + ":8051")

		masterNodeId := labCode.NewKademliaID("masterNode")
		masterNodeAddress := "master:8051"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)

		kademliaNode.RoutingTable.AddContact(masterContact)

		go kademliaNode.RoutingTable.UpdateBucketRoutine()
		go kademliaNode.RoutingTable.FindClosestContactsRoutine()
		go kademliaNode.LookupContactRoutine()
		go kademliaNode.Network.Listen(nodeAddress, 8051)

		kademliaNode.LookupContact(&kademliaNode.RoutingTable.Me)

		go kademliaNode.HeartbeatSignal()

	}

	labCode.RunCLI(kademliaNode)

}
