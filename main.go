package main

import (
	"C"
	"os"

	"github.com/OpenRHINO/RHINO-CLI/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
