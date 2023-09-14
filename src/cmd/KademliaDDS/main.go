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

	if connType == "master" {
		startMasterTypeNode()
	} else if connType == "other" {
		startOtherTypeNode()
	} else {
		fmt.Printf("please provide the arg (master or other)")
	}
}

func startOtherTypeNode() {
	rt := labCode.NewRoutingTable(labCode.NewContact(labCode.NewRandomKademliaID(), ":8001"))
	/* for i := 1; i <= 2; i++ {
		address := fmt.Sprintf(":80%02d", i)
		fmt.Printf(address)
		c := labCode.NewContact(labCode.NewRandomKademliaID(), address)
		rt.AddContact(c)
	} */

	c := labCode.NewContact(labCode.NewRandomKademliaID(), ":8000")
	rt.AddContact(c)

	node := labCode.InitKademliaNode(rt)
	go node.Listen("", 8001)
	heartbeatSignal(&node)
	for {

	}
}

func startMasterTypeNode() {
	rt := labCode.NewRoutingTable(labCode.NewContact(labCode.NewRandomKademliaID(), ":8000"))
	node := labCode.InitKademliaNode(rt)
	go node.Listen("", 8000)
	for {

	}
}

func heartbeatSignal(node *labCode.Kademlia) {
	heartbeat := make(chan bool)

	// Start a goroutine to send heartbeat signals at a regular interval.
	go func() {
		for {
			time.Sleep(time.Second * 15)
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
