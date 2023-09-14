package main

import (
	"fmt"
	"os"
)

func main() {
	//Use KEY = 'HOSTNAME' aswell later to retrive the address of the node
	containerName := os.Getenv("CONTAINER_NAME")
	fmt.Printf("Container Name: %s\n", containerName)
}
