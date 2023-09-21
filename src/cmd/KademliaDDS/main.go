package main

import (
	"kademlia-app/labCode"
	"os"
)

func main() {
	containerName := os.Getenv("CONTAINER_NAME")

	var kademliaNode labCode.Kademlia
	var CLINetworkChan *chan string
	if containerName == "master" {
		kademliaNode, CLINetworkChan = labCode.NewMasterKademliaNode()
		go kademliaNode.Network.Listen("master", 8051)
		go kademliaNode.RoutingTable.UpdateBucketRoutine()
		go kademliaNode.RoutingTable.FindClosestContactsRoutine()
		go kademliaNode.LookupContactRoutine()

		go kademliaNode.HeartbeatSignal()
		go kademliaNode.DataStorageManager()

	} else {
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode, CLINetworkChan = labCode.NewKademliaNode(nodeAddress + ":8051")

		masterNodeId := labCode.NewMasterKademliaID()
		masterNodeAddress := "master:8051"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)

		kademliaNode.RoutingTable.AddContact(masterContact)

		go kademliaNode.Network.Listen(nodeAddress, 8051)
		go kademliaNode.RoutingTable.UpdateBucketRoutine()
		go kademliaNode.RoutingTable.FindClosestContactsRoutine()
		go kademliaNode.LookupContactRoutine()

		kademliaNode.LookupContact(&kademliaNode.RoutingTable.Me)

		go kademliaNode.HeartbeatSignal()
		go kademliaNode.DataStorageManager()

	}

	CLI := labCode.NewCli(kademliaNode, CLINetworkChan)
	CLI.RunCLI()

}
