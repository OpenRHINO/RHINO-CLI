package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

type versionOptions struct {
	clientVersion string
	kubernetesVersion  string
	rhinoVersion string
}


func NewVersionCommand() *cobra.Command {
	versionOpts := &versionOptions{}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of CLI, K8s, and the RHINO Operator",
		Long:  "\nShow the version of CLI, K8s, and the RHINO Operator",
		Example: `rhino version`,
		RunE: versionOpts.version,
	}
	versionCmd.Flags().StringVarP(&versionOpts.clientVersion, "client", "c", "", "the client version")
	versionCmd.Flags().StringVarP(&versionOpts.kubernetesVersion, "kubernetes", "k", "", "the kubernetes version")
	versionCmd.Flags().StringVarP(&versionOpts.rhinoVersion, "rhino", "r", "", "the rhino version")

	return versionCmd
}

func (v *versionOptions) version(cmd *cobra.Command, args []string) error {
	// Check the arguments
	if len(args) != 0 {
		cmd.Help()
		os.Exit(0)
	}

	v.clientVersion = "0.2.0"
	
	

}



