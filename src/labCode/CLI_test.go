package labCode

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests for testing the parseInput function
func TestParseInputCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockUserInp := "userEnteredCommand"
	command, args, err := CLI.parseInput(mockUserInp)

	asserts.Equal("userenteredcommand", command)
	asserts.Equal([]string{}, args)
	asserts.Nil(err)

}

func TestPareInputEmptyStringErr(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockUserInp := ""
	command, args, err := CLI.parseInput(mockUserInp)

	asserts.Equal("", command)
	asserts.Nil(args)
	asserts.Equal("Error: empty string to parseInput", err.Error())
}

func TestParseInputCommandAndArgs(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockUserInp := "userEnteredCommand userEnteredFlagOne userEnteredFlagTwo"
	command, args, err := CLI.parseInput(mockUserInp)

	asserts.Equal("userenteredcommand", command)
	asserts.Equal([]string{"userEnteredFlagOne", "userEnteredFlagTwo"}, args)
	asserts.Nil(err)
}

// Tests for testing the handleHelpCommand
func TestHandleHelpCommandNoArgs(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	emptyArgsSlice := []string{}
	response := CLI.handleHelpCommand(emptyArgsSlice)

	asserts.Equal(ListAvailableCommandsMsg+UseHelpCommandForDetailedUsageMsg, response)

}

func TestHandleHelpCommandPut(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	helpCommandArg := []string{"put"}
	response := CLI.handleHelpCommand(helpCommandArg)

	asserts.Equal(PutCommandUsageMsg, response)

}

func TestHandleHelpCommandGet(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	helpCommandArg := []string{"get"}
	response := CLI.handleHelpCommand(helpCommandArg)

	asserts.Equal(GetCommandUsageMsg, response)

}

func TestHandleHelpCommandExit(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	helpCommandArg := []string{"exit"}
	response := CLI.handleHelpCommand(helpCommandArg)

	asserts.Equal(ExitCommandUsageMsg, response)

}

func TestHandleHelpCommandHelp(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	helpCommandArg := []string{"help"}
	response := CLI.handleHelpCommand(helpCommandArg)

	asserts.Equal(HelpCommandUsageMsg, response)

}

func TestHandleHelpCommandUnknownCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	helpCommandArg := []string{"unknownCommand"}
	response := CLI.handleHelpCommand(helpCommandArg)

	asserts.Equal(ListAvailableCommandsMsg+UseHelpCommandForDetailedUsageMsg, response)
}

// Tests for handleGetCommand
func TestHandleGetCommandNoArgs(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	emptyGetCommandArgs := []string{}
	response, err := CLI.handleGetCommand(emptyGetCommandArgs)

	asserts.Equal("", response)
	asserts.Equal(InvalidArgsForGetCommandMsg+GetCommandUsageMsg, err.Error())
}

func TestHandleGetCommandTooManyArgs(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	GetCommandArgs := []string{"argumentOne", "argumentTwo"}
	response, err := CLI.handleGetCommand(GetCommandArgs)

	asserts.Equal("", response)
	asserts.Equal(InvalidArgsForGetCommandMsg+GetCommandUsageMsg, err.Error())
}

func TestHandleGetCommandInvalidHashInput(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	GetCommandArgs := []string{"sggfyagsajfkasy"}
	response, err := CLI.handleGetCommand(GetCommandArgs)

	asserts.Equal("", response)
	asserts.Equal(NotAValidHashMSG, err.Error())

}

func TestHandleGetCommandValueExists(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedKademliaNode, _ := NewKademliaNode("", "")
	mockedData := "mockedData"
	hashOfMockedData := NewKademliaDataID(mockedData)

	go waitToFireResultsGet(testClINetworkChan, mockedKademliaNode, mockedData)

	getCommandValidHashArgs := []string{hashOfMockedData.String()}
	results, err := CLI.handleGetCommand(getCommandValidHashArgs)

	asserts.Nil(err)
	asserts.Equal("The retrived data: "+mockedData+" "+mockedKademliaNode.Network.Me.ID.String(), results)

}

func waitToFireResultsGet(testClINetworkChan *chan string, kademliaNode Kademlia, mockedData string) {
	time.Sleep(1 * time.Second)
	*testClINetworkChan <- mockedData + " " + kademliaNode.Network.Me.ID.String()
}

func TestHandleGetCommandValueDoesNotExist(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedData := "mockedData"
	hashOfMockedData := NewKademliaDataID(mockedData)

	getCommandValidHashArgs := []string{hashOfMockedData.String()}
	results, err := CLI.handleGetCommand(getCommandValidHashArgs)

	asserts.Nil(err)
	asserts.Equal("Timeout: No data received", results)

}

// Tests for handlePutCommand
func TestHandlePutCommandNoArgs(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	emptyPutCommandArgs := []string{}
	response, err := CLI.handlePutCommand(emptyPutCommandArgs)

	asserts.Equal("", response)
	asserts.Equal(InvalidArgsForPutCommandMsg+PutCommandUsageMsg, err.Error())

}

func TestHandlePutCommandValidStore(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockData := "test"
	mockDataHash := NewKademliaDataID(mockData)

	go waitToFireResultsPut(*mockDataHash, testClINetworkChan)

	putCommandArgs := []string{mockData}
	results, err := CLI.handlePutCommand(putCommandArgs)

	asserts.Nil(err)
	asserts.Equal("Store successful! Hash: "+mockDataHash.String(), results)

}

func waitToFireResultsPut(mockDataHash KademliaID, testCLINetworkChan *chan string) {
	time.Sleep(1 * time.Second)
	*testCLINetworkChan <- mockDataHash.String()

}

func TestHandlePutCommandTimeout(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockData := "test"

	putCommandArgs := []string{mockData}
	results, err := CLI.handlePutCommand(putCommandArgs)

	asserts.Nil(err)
	asserts.Equal("Timeout: could not store", results)

}

// Tests for handleUserInput
func TestHandleUserInpPutCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockData := "test"
	mockDataHash := NewKademliaDataID(mockData)

	go waitToFireResultsPut(*mockDataHash, testClINetworkChan)

	mockedUserInp := "put" + " " + mockData
	results, err := CLI.handleUserInput(mockedUserInp)

	asserts.Nil(err)
	asserts.Equal("Store successful! Hash: "+mockDataHash.String(), results)

}

func TestHandleUserInpGetCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedKademliaNode, _ := NewKademliaNode("", "")
	mockedData := "mockedData"
	hashOfMockedData := NewKademliaDataID(mockedData)

	go waitToFireResultsGet(testClINetworkChan, mockedKademliaNode, mockedData)

	mockedUserInp := "get" + " " + hashOfMockedData.String()
	results, err := CLI.handleUserInput(mockedUserInp)

	asserts.Nil(err)
	asserts.Equal("The retrived data: "+mockedData+" "+mockedKademliaNode.Network.Me.ID.String(), results)
}

func TestHandleUserInpExitCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedUserInp := "exit"
	results, err := CLI.handleUserInput(mockedUserInp)

	asserts.Nil(err)
	asserts.Equal(ShuttingDownNodeMsg, results)
}

func TestHandleUserInpHelpCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedUserInp := "help"
	results, err := CLI.handleUserInput(mockedUserInp)

	asserts.Nil(err)
	asserts.Equal(ListAvailableCommandsMsg+UseHelpCommandForDetailedUsageMsg, results)

}

func TestHandleUserInpOtherCommand(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedUserInp := "nonExistingCommand"
	parsedMockeduserInp, _, _ := CLI.parseInput(mockedUserInp)
	results, err := CLI.handleUserInput(mockedUserInp)

	asserts.Nil(err)
	asserts.Equal(UnknownCommandMsg+"<"+parsedMockeduserInp+"> "+UseHelpToViewCommandsMsg, results)

}

func TestHandleUserInpError(t *testing.T) {
	asserts := assert.New(t)

	testKademliaNode, testClINetworkChan := NewKademliaNode("", "")
	CLI := NewCli(testKademliaNode, testClINetworkChan)

	mockedUserInp := ""
	results, err := CLI.handleUserInput(mockedUserInp)

	asserts.Equal("", results)
	asserts.Equal("Error: empty string to parseInput", err.Error())

}
