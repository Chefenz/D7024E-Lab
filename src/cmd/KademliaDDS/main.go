package main

import (
	"fmt"
	"os"
)

func main() {
	containerName := os.Getenv("HOSTNAME")
	fmt.Printf("Container Name: %s\n", containerName)
}
