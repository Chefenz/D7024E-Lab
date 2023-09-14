package main

import (
	"fmt"
	"os"
)

func main() {
	containerName := os.Getenv("CONTAINER_NAME")
	fmt.Printf("Container Name: %s\n", containerName)
}
