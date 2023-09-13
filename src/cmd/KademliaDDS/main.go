package main

import (
	"fmt"
	"os"
)

func main() {
	// Get the container name from the environment variable.
	containerName := os.Getenv("CONTAINER_ID")

	// Print the container name.
	fmt.Println(containerName)
}
