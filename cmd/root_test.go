/*
 * Copyright 2023 RHINO Team
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cmd_test

import (
	"testing"

	"github.com/OpenRHINO/RHINO-CLI/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	// Test if rootCmd has the correct Use, Short, and Long values
	assert.Equal(t, "rhino", rootCmd.Use)
	assert.Equal(t, "\nRHINO-CLI - Manage your OpenRHINO functions and jobs", rootCmd.Short)

	// Test if rootCmd has the correct subcommands
	expectedSubcommands := []string{"create", "build", "delete", "run", "list", "docker-run", "version", "logs"}
	actualSubcommands := getSubcommandNames(rootCmd)

	assert.Equal(t, len(expectedSubcommands), len(actualSubcommands), "Number of subcommands should be equal")

	for _, expected := range expectedSubcommands {
		assert.Contains(t, actualSubcommands, expected, "Expected subcommand %s is not present", expected)
	}
}

func getSubcommandNames(cmd *cobra.Command) []string {
	subcommandNames := []string{}
	for _, subCmd := range cmd.Commands() {
		subcommandNames = append(subcommandNames, subCmd.Name())
	}
	return subcommandNames
}
