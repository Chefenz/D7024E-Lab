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
		go kademliaNode.Network.Listen("master", 8051, *kademliaNode.StopChan)
		go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)
		go kademliaNode.RoutingTable.FindClosestContactsRoutine(*kademliaNode.StopChan)
		go kademliaNode.LookupContactRoutine(*kademliaNode.StopChan)

		go kademliaNode.HeartbeatSignal(*kademliaNode.StopChan)
		go kademliaNode.LookupCloseContactsToValueRoutine()
		go kademliaNode.DataStorageManager()

	} else {
		nodeAddress := os.Getenv("HOSTNAME")
		kademliaNode, CLINetworkChan = labCode.NewKademliaNode(nodeAddress + ":8051")

		masterNodeId := labCode.NewMasterKademliaID()
		masterNodeAddress := "master:8051"
		masterContact := labCode.NewContact(masterNodeId, masterNodeAddress)

		kademliaNode.RoutingTable.AddContact(masterContact)

		go kademliaNode.Network.Listen(nodeAddress, 8051, *kademliaNode.StopChan)
		go kademliaNode.RoutingTable.UpdateBucketRoutine(*kademliaNode.StopChan)
		go kademliaNode.RoutingTable.FindClosestContactsRoutine(*kademliaNode.StopChan)
		go kademliaNode.LookupContactRoutine(*kademliaNode.StopChan)

		kademliaNode.LookupContact(&kademliaNode.RoutingTable.Me)

		go kademliaNode.HeartbeatSignal(*kademliaNode.StopChan)
		go kademliaNode.LookupCloseContactsToValueRoutine()
		go kademliaNode.DataStorageManager()

	}

	CLI := labCode.NewCli(kademliaNode, CLINetworkChan)
	CLI.RunCLI()

}
