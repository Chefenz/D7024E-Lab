package labCode

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	// Constants for messages
	ShuttingDownNodeMsg               = "Exiting - closing down node"
	UnknownCommandMsg                 = "Unknown command - "
	UseHelpToViewCommandsMsg          = "Use command <help> to get a list of available commands"
	ListAvailableCommandsMsg          = "The available commands are: \n\n\tput\tStore data in the DDS\n\tget\tRetrieve data from the DDS\n\texit\tTerminates the node\n\thelp\tList all available commands\t\n"
	UseHelpCommandForDetailedUsageMsg = "\nUse help <command> for more information about a command"
	InvalidArgsForPutCommandMsg       = "Invalid args for <put> command\n"
	InvalidArgsForGetCommandMsg       = "Invalid args for <get> command\n"
	NotAValidHashMSG                  = "Argument was not a 160 bit hash\n"

	PutCommandUsageMsg = "Usage: \n\t put [argument] \n\nDescription: \n\t Takes a single argument, " +
		"the contents of the file you are uploading, and outputs the hash of the object, if it could be uploaded successfully"
	GetCommandUsageMsg = "Usage: \n\t get [argument] \n\nDescription: \n\t Takes a hash as its only argument, " +
		"and outputs the contents of the object and the node it was retrieved from, if it could be downloaded successfully"
	ExitCommandUsageMsg = "Usage: \n\t exit \n\nDescription: \n\t Terminates the node"
	HelpCommandUsageMsg = "Usage \n\t help or help <command> \n\nDescription: \n\t Lists all available commands " +
		"or more information about a specific command"
)

type CLI struct {
	KademliaNode   Kademlia
	CLINetworkChan *chan string
}

// Creates and returns a new CLI struct
func NewCli(kademliaNode Kademlia, CLIChan *chan string) CLI {
	return CLI{KademliaNode: kademliaNode, CLINetworkChan: CLIChan}
}

// Runs the CLI interface
func (cli *CLI) RunCLI() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(" ~ ")

		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unexpected Error: ", err)
			os.Exit(1)
		}
		if userInput == "\n" {
			continue
		}

		response, err := cli.handleUserInput(userInput)
		if err != nil {
			cli.printToTerminal(err.Error())
			continue
		}
		if response == ShuttingDownNodeMsg {
			cli.printToTerminal(ShuttingDownNodeMsg)
			os.Exit(0)
		}

		cli.printToTerminal(response)

	}

}

// Handles the userInput
func (cli *CLI) handleUserInput(userInp string) (string, error) {
	command, args, err := cli.parseInput(userInp)
	if err != nil {
		return "", err
	}

	switch command {
	case "put":
		response, err := cli.handlePutCommand(args)
		return response, err
	case "get":
		response, err := cli.handleGetCommand(args)
		return response, err
	case "exit":
		response := ShuttingDownNodeMsg
		return response, nil
	case "help":
		response := cli.handleHelpCommand(args)
		return response, nil
	default:
		response := UnknownCommandMsg + "<" + command + "> " + UseHelpToViewCommandsMsg
		return response, nil

	}
}

// Parses the userinput
func (cli *CLI) parseInput(userInput string) (string, []string, error) {
	if userInput == "" {
		return "", nil, errors.New("Error: empty string to parseInput")
	}

	args := strings.Split(userInput, " ")
	command := strings.Trim(strings.ToLower(args[0]), "\n")

	return command, args[1:], nil
}

// Handles the put command by telling the kademlia node to do a store RPC
func (cli *CLI) handlePutCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New(InvalidArgsForPutCommandMsg + PutCommandUsageMsg)

	}

	dataStr := strings.Join(args, " ")
	dataByte := []byte(dataStr)
	cli.KademliaNode.Store(dataByte)

	result := ""
	select {
	case resp := <-*cli.CLINetworkChan:
		result = "Store successful! Hash: " + resp
	case <-time.After(rpcTimeout + time.Millisecond*90):
		result = "Timeout: could not store"
	}

	return result, nil

}

// Handles the get command by telling the kademlia node to do a LookUpData RPC
func (cli *CLI) handleGetCommand(args []string) (string, error) {
	if len(args) == 0 || len(args) > 1 {
		return "", errors.New(InvalidArgsForGetCommandMsg + GetCommandUsageMsg)
	}

	dataIDStr := strings.Join(args, " ")
	if len(dataIDStr) != 40 {
		fmt.Println(len(dataIDStr))
		return "", errors.New(NotAValidHashMSG)
	}

	cli.KademliaNode.LookupData(dataIDStr)

	result := ""
	select {
	case resp := <-*cli.CLINetworkChan:
		result = "The retrived data: " + resp
	case <-time.After(rpcTimeout + time.Millisecond*90):
		result = "Timeout: No data received"
	}

	return result, nil
}

// Handles the help command
func (cli *CLI) handleHelpCommand(args []string) string {
	var rtnMsg string

	if len(args) > 0 {
		command := strings.Trim(strings.ToLower(args[0]), "\n")

		switch command {
		case "put":
			rtnMsg = PutCommandUsageMsg
		case "get":
			rtnMsg = GetCommandUsageMsg
		case "exit":
			rtnMsg = ExitCommandUsageMsg
		case "help":
			rtnMsg = HelpCommandUsageMsg
		default:
			rtnMsg = ListAvailableCommandsMsg + UseHelpCommandForDetailedUsageMsg
		}
	} else {
		rtnMsg = ListAvailableCommandsMsg + UseHelpCommandForDetailedUsageMsg
	}

	return rtnMsg
}

// Prints output to the terminal
func (cli *CLI) printToTerminal(str string) {
	fmt.Println()
	fmt.Println(str)
	fmt.Println()
}
