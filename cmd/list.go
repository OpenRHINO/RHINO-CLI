package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rhinojob "github.com/OpenRHINO/RHINO-Operator/api/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var namespace string
var kubeconfig string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rhino jobs",
	Long:  "\nList all rhino jobs",
	Example: `  rhino list
  rhino list --namespace user_func`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var configPath string
		if len(kubeconfig) == 0 {
			if home := homedir.HomeDir(); home != "" {
				configPath = filepath.Join(home, ".kube", "config")
			} else {
				fmt.Println("Error: kubeconfig file not found, please use --config to specify the absolute path")
				os.Exit(0)
			}
		} else {
			configPath = kubeconfig
		}
		config, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return err
		}

		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			return err
		}

		list, err := listRhinoJob(dynamicClient, namespace)
		if err != nil {
			return err
		}
		if len(list.Items) == 0 {
			fmt.Println("RhinoJob not found in namespace")
			os.Exit(0)
		}
		fmt.Printf("%-20s\t%-15s\t%-5s\n", "Name", "Parallelism", "Status")
		for _, rj := range list.Items {
			fmt.Printf("%-20v\t%-15v\t%-5v\n", rj.Name, *rj.Spec.Parallelism, rj.Status.JobStatus)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace of the rhinojob")
	listCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "kubernetes config path")
}

func listRhinoJob(client dynamic.Interface, namespace string) (*rhinojob.RhinoJobList, error) {
	list, err := client.Resource(RhinoJobGVR).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var rjList rhinojob.RhinoJobList
	if err := json.Unmarshal(data, &rjList); err != nil {
		return nil, err
	}
	return &rjList, nil
}
