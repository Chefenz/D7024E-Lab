package main

import (
	"fmt"
	"os"
)

func main() {
	// Get the container name from the environment variable.
	containerName := os.Getenv("NAMES")

	// Print the container name.
	fmt.Println(containerName)
}
