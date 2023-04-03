package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	rhinojob "github.com/OpenRHINO/RHINO-Operator/api/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/homedir"
)

type RunOptions struct {
	parallel   int
	timeToLive int
	dataPath   string
	dataServer string
	funcName   string

	kubeconfig string
	namespace  string
}

func NewRunCommand() *cobra.Command {
	runOpts := &RunOptions{}
	runCmd := &cobra.Command{
		Use:   "run [image]",
		Short: "Submit and run a RHINO job",
		Long:  "\nSubmit an MPI function/project and run it as a RHINO job",
		Example: `  rhino run hello:v1.0 --namespace user_space
  rhino run foo/matmul:v2.1 --np 4 -- arg1 arg2 
  rhino run mpi/testbench -n 32 -t 800 --server 10.0.0.7 --dir /mnt -- --in=/data/file --out=/data/out`,
		RunE: runOpts.run,
	}

	runCmd.Flags().StringVar(&runOpts.dataServer, "server", "", "IP address of an NFS server")
	runCmd.Flags().StringVar(&runOpts.dataPath, "dir", "", "a directory in the NFS server, to store data and shared with all the MPI processes")
	runCmd.MarkFlagsRequiredTogether("server", "dir")
	runCmd.Flags().IntVar(&runOpts.parallel, "np", 1, "the number of MPI processes")
	runCmd.Flags().IntVarP(&runOpts.timeToLive, "ttl", "t", 600, "Time To Live (seconds). The RHINO job will be deleted after this time, whether it is completed or not.")
	runCmd.Flags().StringVarP(&runOpts.namespace, "namespace", "n", "", "the namespace of the RHINO job")
	runCmd.Flags().StringVar(&runOpts.kubeconfig, "kubeconfig", "", "the path of the kubeconfig file")

	return runCmd
}

func (r *RunOptions) run(cmd *cobra.Command, args []string) error {
	// Check the arguments
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}
	r.funcName = getFuncName(args[0])
	if r.parallel < 1 {
		return fmt.Errorf("the number of MPI processes (--np) must be greater than 0")
	}
	if r.timeToLive < 0 {
		return fmt.Errorf("the time to live (--ttl) must be greater than or equal to 0")
	}
	if r.kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			r.kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			fmt.Println("Error: kubeconfig file not found, please use --config to specify the absolute path")
			os.Exit(0)
		}
	}

	dynamicClient, currentNamespace, err := buildFromKubeconfig(r.kubeconfig)
	if err != nil {
		return err
	}
	if r.namespace == "" {
		r.namespace = *currentNamespace
	}

	// Create a RHINO job
	_, err = r.runRhinoJob(dynamicClient, args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println("RhinoJob", r.funcName, "created")
	return nil
}


func (r *RunOptions) printYAML(args []string) (yamlFile string) {
	yamlFile = `apiVersion: openrhino.org/v1alpha1
kind: RhinoJob
metadata:
  labels:
    app.kubernetes.io/name: rhinojob 
    app.kubernetes.io/instance: rhinojob-`
	yamlFile += r.funcName + `
    app.kubernetes.io/part-of: rhino-operator
    app.kubernetes.io/managed-by: kustomizex
    app.kubernetes.io/created-by: rhino-operator
  name: rhinojob-`
	yamlFile += r.funcName + `
spec:
  image: "`
	yamlFile += args[0] + `"
  ttl: `
	yamlFile += strconv.Itoa(r.timeToLive) + `
  parallelism: `
	yamlFile += strconv.Itoa(r.parallel) + ` 
  appExec: "./mpi-func"`
	if len(args) > 1 {
		yamlFile += `
  appArgs: [`
		for i := 1; i < len(args); i++ {
			yamlFile += `"` + args[i] + `", `
		}
		yamlFile += `]`
	}
	if len(r.dataServer) != 0 {
		yamlFile += `
  dataServer: "` + r.dataServer + `"`
		yamlFile += `
  dataPath: "` + r.dataPath + `"`
	}
	return yamlFile
}

func (r *RunOptions) runRhinoJob(client dynamic.Interface, args []string) (*rhinojob.RhinoJobList, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	rhinojobYAML := r.printYAML(args)
	if _, _, err := decoder.Decode([]byte(rhinojobYAML), nil, obj); err != nil {
		return nil, err
	}
	createdRhinoJob, err := client.Resource(RhinoJobGVR).Namespace(r.namespace).Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}
	data, err := createdRhinoJob.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var rj rhinojob.RhinoJobList
	if err := json.Unmarshal(data, &rj); err != nil {
		return nil, err
	}
	return &rj, nil
}