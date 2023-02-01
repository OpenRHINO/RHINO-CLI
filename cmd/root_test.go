package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Root Command: rhino
// Here we want to check whether our help message is the same as
// the output of terminal when we type "rhino" in bash/zsh/etc

func TestRootCmd(t *testing.T) {
	// To get our help message, we call `rootCmd.Help`, and redirect
	// its output to our string builder
	// To get the output of terminal, we call `rootCmd.Execute`, and
	// also redirect its output to our string builder
	cmdOutput := new(strings.Builder)

	// redirect output from stdout to our string builder
	rootCmd.SetOut(cmdOutput)
	rootCmd.SetErr(cmdOutput)

	// initialization cobra default commands
	// when we call `rootCmd.Execute`(same result as you type "rhino" in terminal)
	// cobra will automatically init some default commands like `help` and `completion`
	// but if we simply call `rootCmd.Help` without these initialization
	// message in the string builder won't include these commands
	// so this test will not pass
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()
	rootCmd.InitDefaultHelpFlag()

	rootCmd.Help()
	expected := cmdOutput.String()
	cmdOutput.Reset()

	// set command line arguments to an empty string array
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	assert.Equal(t, nil, err, "Execute Failed: Non-nil Error")

	// compare our help message with terminal output
	// helper message: `expeted`
	// terminal output: `actual`
	actual := cmdOutput.String()
	assert.Equal(t, expected, actual)
}
