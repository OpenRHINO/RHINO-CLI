package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var imageName string
var execArgs []string
var parallel int
var execTime int
var dataPath string
var dataServer string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Submit and run rhino job",
	// TODO: edit Usage(run build ...)
	Long: "\nSubmit and run rhino job",
	Example: `  rhino run hello:v1.0
  rhino run foo/matmul:v2.1 --np 4 -- arg1 arg2 
  rhino run mpi/testbench -n 32 -t 800 --server 10.0.0.7 --dir /mnt -- --in=/data/file --out=/data/out`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		fmt.Println(printYAML(args))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&dataServer, "server", "", "NFS server ip")
	runCmd.Flags().StringVar(&dataPath, "dir", "", "shared directory in NFS server")
	runCmd.MarkFlagsRequiredTogether("server", "dir")
	runCmd.Flags().IntVarP(&parallel, "np", "n", 1, "mpi processes")
	runCmd.Flags().IntVarP(&execTime, "ttl", "t", 600, "estimated execution time(s)")
}

func printYAML(args []string) (yamlFile string) {
	funcName := getFuncName(args[0])
	fmt.Println(args)
	yamlFile = `apiVersion: openrhino.org/v1alpha1
kind: RhinoJob
metadata:
  labels:
    app.kubernetes.io/name: rhinojob 
    app.kubernetes.io/instance: rhinojob-`
	yamlFile += funcName + `
    app.kubernetes.io/part-of: rhino-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: rhino-operator
  name: rhinojob-`
	yamlFile += funcName + `
spec:
  image: "`
	yamlFile += args[0] + `"
  ttl: `
	yamlFile += strconv.Itoa(execTime) + `
  parallelism: `
	yamlFile += strconv.Itoa(parallel) + ` 
  appExec: "./`
	yamlFile += funcName
	if len(args) > 1 {
		yamlFile += `"
  appArgs: [`
		for i := 1; i < len(args); i++ {
			yamlFile += `"` + args[i] + `", `
		}
		yamlFile += `]`
	}
	if len(dataServer) != 0 {
		yamlFile += `
  dataServer: "` + dataServer + `"`
		yamlFile += `
  dataPath: "` + dataPath + `"`
	}
	return yamlFile
}
