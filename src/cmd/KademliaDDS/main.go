package main

import (
	"fmt"
	"os"
)

func main() {
	containerName := os.Args[1]
	fmt.Println("The name of the container is:", containerName)
}
