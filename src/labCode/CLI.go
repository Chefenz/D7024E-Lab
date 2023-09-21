package labCode

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
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

func NewCli(kademliaNode Kademlia, CLIChan *chan string) CLI {
	return CLI{KademliaNode: kademliaNode, CLINetworkChan: CLIChan}
}

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

		command, args, err := cli.parseInput(userInput)
		if err != nil {
			fmt.Println("Unexpected Error: ", err)
			os.Exit(1)
		}

		switch command {
		case "put":
			response, err := cli.handlePutCommand(args)
			if err != nil {
				cli.printToTerminal(err.Error())
				break
			}
			cli.printToTerminal(response)
		case "get":
			response, err := cli.handleGetCommand(args)
			if err != nil {
				cli.printToTerminal(err.Error())
				break
			}
			cli.printToTerminal(response)
		case "exit":
			cli.printToTerminal(ShuttingDownNodeMsg)
			os.Exit(0)
		case "help":
			response := cli.handleHelpCommand(args)
			cli.printToTerminal(response)
		default:
			cli.printToTerminal(UnknownCommandMsg + "<" + command + "> " + UseHelpToViewCommandsMsg)
		}

	}

}

func (cli *CLI) parseInput(userInput string) (string, []string, error) {
	if userInput == "" {
		return "", nil, errors.New("empty string")
	}
	args := strings.Split(userInput, " ")
	command := strings.Trim(strings.ToLower(args[0]), "\n")

	return command, args[1:], nil
}

func (cli *CLI) handlePutCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New(InvalidArgsForPutCommandMsg + PutCommandUsageMsg)

	}

	dataStr := strings.Join(args, " ")
	dataByte := []byte(dataStr)
	fmt.Println(dataByte)
	fmt.Println("Before Store")
	cli.KademliaNode.Store(dataByte)
	result := <-*cli.CLINetworkChan
	fmt.Println("After Store")

	if result == "" {
		return "", errors.New("Could not save")
	}

	return "Store successful! Hash: " + result, nil

}

func (cli *CLI) handleGetCommand(args []string) (string, error) {
	if len(args) == 0 || len(args) > 1 {
		return "", errors.New(InvalidArgsForGetCommandMsg + GetCommandUsageMsg)

	}

	//TODO - HANDLE THE ACTUAL GET COMMAND AND VERIFY THAT THE OTHER ARGUMENT IS A 160-BIT SHA1-HASH
	return "TEMPSTR", nil
}

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

func (cli *CLI) printToTerminal(str string) {
	fmt.Println()
	fmt.Println(str)
	fmt.Println()
}
