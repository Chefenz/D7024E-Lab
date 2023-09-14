package labCode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for testing the parseInput function
func TestParseInputCommand(t *testing.T) {
	asserts := assert.New(t)

	mockUserInp := "userEnteredCommand"

	command, args, err := parseInput(mockUserInp)

	asserts.Equal("userenteredcommand", command)
	asserts.Equal([]string{}, args)
	asserts.Nil(err)
}

func TestParseInputCommandAndArgs(t *testing.T) {
	asserts := assert.New(t)

	mockUserInp := "userEnteredCommand userEnteredFlagOne userEnteredFlagTwo"

	command, args, err := parseInput(mockUserInp)

	asserts.Equal("userenteredcommand", command)
	asserts.Equal([]string{"userEnteredFlagOne", "userEnteredFlagTwo"}, args)
	asserts.Nil(err)
}

func TestParseInputEmptyUserInput(t *testing.T) {
	asserts := assert.New(t)

	mockEmptyUserInp := ""

	command, args, err := parseInput(mockEmptyUserInp)

	asserts.Equal("", command)
	asserts.Nil(args)
	asserts.Error(err)
}

// Tests for testing the handleHelpCommand
func TestHandleHelpCommandNoArgs(t *testing.T) {
	asserts := assert.New(t)

	emptyArgsSlice := []string{}

	response := handleHelpCommand(emptyArgsSlice)

	asserts.Equal(ListAvailableCommandsMsg+UseHelpCommandForDetailedUsageMsg, response)

}

func TestHandleHelpCommandPut(t *testing.T) {
	asserts := assert.New(t)

	helpCommandArg := []string{"put"}

	response := handleHelpCommand(helpCommandArg)

	asserts.Equal(PutCommandUsageMsg, response)

}

func TestHandleHelpCommandGet(t *testing.T) {
	asserts := assert.New(t)

	helpCommandArg := []string{"get"}

	response := handleHelpCommand(helpCommandArg)

	asserts.Equal(GetCommandUsageMsg, response)

}

func TestHandleHelpCommandExit(t *testing.T) {
	asserts := assert.New(t)

	helpCommandArg := []string{"exit"}

	response := handleHelpCommand(helpCommandArg)

	asserts.Equal(ExitCommandUsageMsg, response)

}

func TestHandleHelpCommandHelp(t *testing.T) {
	asserts := assert.New(t)

	helpCommandArg := []string{"help"}

	response := handleHelpCommand(helpCommandArg)

	asserts.Equal(HelpCommandUsageMsg, response)

}

func TestHandleHelpCommandUnknownCommand(t *testing.T) {
	asserts := assert.New(t)

	helpCommandArg := []string{"unknownCommand"}

	response := handleHelpCommand(helpCommandArg)

	asserts.Equal(ListAvailableCommandsMsg+UseHelpCommandForDetailedUsageMsg, response)
}

// Tests for handleGetCommand
func TestHandleGetCommandNoArgs(t *testing.T) {
	asserts := assert.New(t)

	getCommandArgs := []string{}

	response, err := handleGetCommand(getCommandArgs)

	asserts.Equal("", response)
	asserts.Error(err)
}

func TestHandleGetCommandTooManyArgs(t *testing.T) {
	asserts := assert.New(t)

	getCommandArgs := []string{"args1", "args2"}

	response, err := handleGetCommand(getCommandArgs)

	asserts.Equal("", response)
	asserts.Error(err)
}

func TestHandleGetCommandValueExists(t *testing.T) {} //TODO

func TestHandleGetCommandValueDoesNotExist(t *testing.T) {} //TODO

// Tests for handlePutCommand
func TestHandlePutCommandNoArgs(t *testing.T) {
	asserts := assert.New(t)

	putCommandArgs := []string{}

	response, err := handlePutCommand(putCommandArgs)

	asserts.Equal("", response)
	asserts.Error(err)
}

func TestHandlePutCommandSuccessfulStorage(t *testing.T) {} //TODO

func TestHandlePutCommandUnsuccessfulStorage(t *testing.T) {} //TODO
