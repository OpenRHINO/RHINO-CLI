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
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	rhinojob "github.com/OpenRHINO/RHINO-Operator/api/v1alpha2"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type LogsOptions struct {
	rhinojobName string
	kubeconfig   string
	namespace    string
	launcher     bool
	worker       int
	follow       bool
}

func NewLogsCommand() *cobra.Command {
	logsOpts := &LogsOptions{}
	logsCmd := &cobra.Command{
		Use:   "logs [name]",
		Short: "Print logs for a RHINO job",
		Long:  "\nPrint pod logs for specified RHINO job by name",
		Example: `  rhino logs user_job --namespace user_space
  rhino logs user_job -w 0 -f`,
		Args: logsOpts.argsCheck,
		RunE: logsOpts.runLogs,
	}

	logsCmd.Flags().StringVarP(&logsOpts.namespace, "namespace", "n", "", "namespace of the RHINO job")
	logsCmd.Flags().StringVar(&logsOpts.kubeconfig, "kubeconfig", "", "path to the kubeconfig file")
	logsCmd.Flags().BoolVarP(&logsOpts.launcher, "launcher", "l", false, "get the log of the launcher pod")
	logsCmd.Flags().IntVarP(&logsOpts.worker, "worker", "w", -1, "get the log of w_th worker pod(0 <= w < worker_num)")
	logsCmd.Flags().BoolVarP(&logsOpts.follow, "follow", "f", false, "continuously track the latest updates to the log output")
	return logsCmd
}

func (l *LogsOptions) argsCheck(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("job name cannot be empty")
	}
	l.rhinojobName = args[0]

	if l.worker < -1 {
		return fmt.Errorf("worker pod number cannot be negative")
	}

	if l.worker >= 0 && l.launcher {
		return fmt.Errorf("cannot specify both launcher and worker pod")
	}

	var err error
	l.kubeconfig, err = getKubeconfigPath(l.kubeconfig)
	if err != nil {
		return fmt.Errorf("%v, please set the kubeconfig path by --kubeconfig", err)
	}

	return nil
}

func (l *LogsOptions) runLogs(cmd *cobra.Command, args []string) error {
	dynamicClient, currentNamespace, err := buildFromKubeconfig(l.kubeconfig)
	if err != nil {
		return err
	}
	if l.namespace == "" {
		l.namespace = *currentNamespace
	}

	config, err := clientcmd.BuildConfigFromFlags("", l.kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// 获取 rhinojob
	rhinoJobUnstructured, err := dynamicClient.Resource(RhinoJobGVR).Namespace(l.namespace).Get(context.TODO(), l.rhinojobName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	rhinoJobBytes, err := rhinoJobUnstructured.MarshalJSON()
	if err != nil {
		return err
	}

	var rj rhinojob.RhinoJob
	err = json.Unmarshal(rhinoJobBytes, &rj)
	if err != nil {
		return err
	}

	if l.launcher || l.worker == -1 {
		// 获取 launcher pod 的日志
		err = l.getPodLogs(clientset, rj.Status.LauncherPodNames, 0)
		if err != nil {
			return err
		}
	} else if l.worker >= 0 {
		// 获取指定编号的 worker pod 的日志
		if l.worker > len(rj.Status.WorkerPodNames)-1 {
			return fmt.Errorf("worker pod index out of range [0, %d]", len(rj.Status.WorkerPodNames)-1)
		}
		err = l.getPodLogs(clientset, rj.Status.WorkerPodNames, l.worker)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LogsOptions) getPodLogs(clientset *kubernetes.Clientset, podNames []string, index int) error {
	podLogOpts := corev1.PodLogOptions{
		Follow: l.follow,
	}

	request := clientset.CoreV1().Pods(l.namespace).GetLogs(podNames[index], &podLogOpts)
	podLogs, err := request.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	if l.follow {
		// If follow is set, then continuously output the logs to the console
		_, err = io.Copy(os.Stdout, podLogs)
		if err != nil {
			return fmt.Errorf("failed to output logs continuously: %w", err)
		}
	} else {
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return fmt.Errorf("failed to retrieve logs: %w", err)
		}
		str := buf.String()
		fmt.Printf("Logs for Pod %s: \n%s\n", podNames[index], str)
	}
	return nil
}
