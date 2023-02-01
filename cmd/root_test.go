package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	cmdOutput := new(strings.Builder)
	rootCmd.SetOut(cmdOutput)
	rootCmd.SetErr(cmdOutput)

	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()
	rootCmd.InitDefaultHelpFlag()
	rootCmd.Help()

	expected := cmdOutput.String()
	cmdOutput.Reset()
	
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	assert.Equal(t, nil, err, "Execute Failed: Non-nil Error")

	actual := cmdOutput.String()
	assert.Equal(t, expected, actual)
}
