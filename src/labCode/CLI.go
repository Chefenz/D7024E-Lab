package labCode

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	SHUTTING_DOWN_NODE_MSG                  = "Exiting - closing down node"
	UNKNOWN_COMMAND_MSG                     = "Unkonwn command - "
	USE_HELP_TO_VIEW_COMMADS_MSG            = "Use command <help> to get a list of availible commands"
	LIST_AVAILIBLE_COMMANDS_MSG             = "The avalible commands are: \n\n\tput\tStore data in the DDS\n\tget\tRetrive data from the DDS\n\texit\tTerminates the node\n\thelp\tList all availible commands\t\n"
	USE_HELP_COMMAND_FOR_DETAILED_USAGE_MSG = "\nUse help <command> for more information about a command"

	INVALID_ARGS_FOR_PUT_COMMAND_MSG = "Invalid args for <put> command\n"
	INVALID_ARGS_FOR_GET_COMMAND_MSG = "Invalid args for <get> command\n"

	PUT_COMMAND_USAGE_MSG = "Usage: \n\t put [argument] \n\nDescription: \n\t Takes a single argument, " +
		"the contents of the file you are uploading, and outputs the hash of the object, if it could be uploaded successfully"
	GET_COMMAND_USAGE_MSG = "Usage: \n\t get [argument] \n\nDescription: \n\t Takes a hash as its only argument, " +
		"and outputs the contents of the object and the node it was retrieved from, if it could be downloaded successfully"
	EXIT_COMMAND_USAGE_MSG = "Usage: \n\t exit \n\nDescription: \n\t Terminates the node"
	HELP_COMMNAD_USAGE_MSG = "Usage \n\t help or help <command> \n\nDescription: \n\t Lists all availible commands " +
		"or more information about a specific command"
)

func StartCLI() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(" ~ ")

		userInput, err := reader.ReadString('\n')

		if err != nil {
			panic(err)
		}
		if userInput == "\n" {
			continue
		}

		parseInput(userInput)
	}

}

func parseInput(userInput string) {
	args := strings.Split(userInput, " ")
	command := strings.Trim(strings.ToLower(args[0]), "\n")

	switch command {
	case "put":
		handlePutCommand(args[1:])
	case "get":
		handleGetCommand(args[1:])
	case "exit":
		printToTerminal(SHUTTING_DOWN_NODE_MSG)
		os.Exit(0)
	case "help":
		handleHelpCommand(args[1:])
	default:
		printToTerminal(UNKNOWN_COMMAND_MSG + "<" + command + "> " + USE_HELP_TO_VIEW_COMMADS_MSG)
	}

}

func handlePutCommand(args []string) {
	if len(args) == 0 {
		printToTerminal(INVALID_ARGS_FOR_PUT_COMMAND_MSG + PUT_COMMAND_USAGE_MSG)
		return
	}

	//TODO - HANDLE THE ACTUAL STORING OF THE DATA
	printToTerminal("TEMPSTR")

}

func handleGetCommand(args []string) {
	if len(args) == 0 || len(args) > 1 {
		printToTerminal(INVALID_ARGS_FOR_GET_COMMAND_MSG + GET_COMMAND_USAGE_MSG)
		return
	}

	//TODO - HANDLE THE ACTUAL GET COMMAND AND VERIFY THAT THE OTHER ARGUMENT IS A 160-BIT SHA1-HASH
	printToTerminal("TEMPSTR")
}

func handleHelpCommand(args []string) {
	if len(args) > 0 {
		command := strings.Trim(strings.ToLower(args[0]), "\n")

		switch command {
		case "put":
			printToTerminal(PUT_COMMAND_USAGE_MSG)
		case "get":
			printToTerminal(GET_COMMAND_USAGE_MSG)
		case "exit":
			printToTerminal(EXIT_COMMAND_USAGE_MSG)
		case "help":
			printToTerminal(HELP_COMMNAD_USAGE_MSG)
		default:
			printToTerminal(LIST_AVAILIBLE_COMMANDS_MSG + USE_HELP_COMMAND_FOR_DETAILED_USAGE_MSG)
		}
	} else {
		printToTerminal(LIST_AVAILIBLE_COMMANDS_MSG + USE_HELP_COMMAND_FOR_DETAILED_USAGE_MSG)
	}
}

func printToTerminal(str string) {
	fmt.Println()
	fmt.Println(str)
	fmt.Println()
}
