package main

import (
	"fmt"
	"kademlia-app/labCode"
	"os"
	"time"
)

func main() {
	connType := os.Args[1]

	fmt.Println("Running node as:", connType)

	//newKademliaID := labCode.NewRandomKademliaID()
	//me := labCode.NewContact(newKademliaID, "localhost") //Localhost for now
	//routingTable := labCode.NewRoutingTable(me)
	//fmt.Println(routingTable)

	if connType == "master" {
		go startMasterTypeNode()
	} else if connType == "other" {
		go startOtherTypeNode()
		println("test")
	} else {
		fmt.Printf("please provide the arg (master or other)")
	}

	//Start the command line user face
	//labCode.RunCLI()
}

func startOtherTypeNode() {
	node := labCode.NewKademliaNode()
	go node.Listen("", 8001)
	heartbeatSignal(&node)
	time.Sleep(time.Second)
	//node.Ping(&c)
	for {

	}
}

func startMasterTypeNode() {
	node := labCode.NewKademliaNode()
	go node.Listen("", 8000)
	for {

	}
}

func heartbeatSignal(node *labCode.Kademlia) {
	heartbeat := make(chan bool)

	// Start a goroutine to send heartbeat signals at a regular interval.
	go func() {
		for {
			time.Sleep(time.Second * 5)
			heartbeat <- true
		}
	}()

	// Listen for heartbeat signals.
	for {
		select {
		case <-heartbeat:
			fmt.Println("Heartbeat")
			node.SendHeartbeatMessage()
		default:
			// No heartbeat received.
		}
	}
}
