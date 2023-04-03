package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


func NewVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Short:   "Show the version of CLI, K8s, and the RHINO Operator",
		Long:    "\nShow the version of CLI, K8s, and the RHINO Operator",
		Example: `rhino version`,
		RunE:    version,
	}

	return versionCmd
}

func (v *versionOptions) version(cmd *cobra.Command, args []string) error {
	// Check the arguments
	if len(args) != 0 {
		cmd.Help()
		os.Exit(0)
		fmt.Print("Length of arguments is not zero.")
	}

	v.clientVersion = "0.2.0"

	return nil

}
