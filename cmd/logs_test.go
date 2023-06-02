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
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogsSingleJob(t *testing.T) {
	// change work directory to ${workspaceFolder}
	cwd, err := os.Getwd()
	assert.Equal(t, nil, err, "test logs failed: %s", errorMessage(err))
	if strings.HasSuffix(cwd, "cmd") {
		os.Chdir("..")
	}
	rootCmd := NewRootCommand()

	// use `rhino build` to build template
	os.Chdir("templates/func")
	testFuncName := "test-logs-func-cpp"
	testFuncImageName := "test-logs-func-cpp:v1"
	rootCmd.SetArgs([]string{"build", "--image", testFuncImageName})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "preparatory work build failed: %s", errorMessage(err))

	// test run command
	testFuncRunNamespace := "test-logs-namespace"
	execShellCmd("kubectl", []string{"create", "namespace", testFuncRunNamespace})
	rootCmd.SetArgs([]string{"run", testFuncImageName, "--namespace", testFuncRunNamespace})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "preparatory work run failed: %s", errorMessage(err))
	fmt.Println("Wait 10s and check job status")
	time.Sleep(10 * time.Second)

	// test logs command for launcher pod
	rescueStdout := os.Stdout
	r, w, err := os.Pipe()
	assert.Equal(t, nil, err, "test logs failed: %s", errorMessage(err))

	os.Stdout = w
	rootCmd.SetArgs([]string{"logs", testFuncName, "--namespace", testFuncRunNamespace})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "test launcher pod logs failed: %s", errorMessage(err))

	// switch back
	w.Close()
	os.Stdout = rescueStdout

	// copy the output of `rhino logs` from read end of the pipe to a string builder
	buf := new(strings.Builder)
	io.Copy(buf, r)
	r.Close()

	cmdOutput := buf.String()
	cmdOutputLines := strings.Split(cmdOutput, "\n")

	var foundLauncherPodLog bool
	testLauncherPodLogPrefix := "Logs for Pod"
	for _, line := range cmdOutputLines {
		if strings.HasPrefix(line, testLauncherPodLogPrefix) {
			foundLauncherPodLog = true
			break
		}
	}
	assert.Equal(t, true, foundLauncherPodLog, "test logs failed: logs output does not contain logs of the launcher pod created in this test")

	// test logs command for worker pod
	rescueStdout = os.Stdout
	r, w, err = os.Pipe()
	assert.Equal(t, nil, err, "test logs failed: %s", errorMessage(err))

	os.Stdout = w
	rootCmd.SetArgs([]string{"logs", testFuncName, "-w", "0", "--namespace", testFuncRunNamespace})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "test worker pod logs failed: %s", errorMessage(err))

	// switch back
	w.Close()
	os.Stdout = rescueStdout

	// copy the output of `rhino logs` from read end of the pipe to a string builder
	buf = new(strings.Builder)
	io.Copy(buf, r)
	r.Close()

	cmdOutput = buf.String()
	cmdOutputLines = strings.Split(cmdOutput, "\n")

	var foundWorkerPodLog bool
	testWorkerPodLogPrefix := "Logs for Pod"
	for _, line := range cmdOutputLines {
		if strings.HasPrefix(line, testWorkerPodLogPrefix) {
			foundWorkerPodLog = true
			break
		}
	}
	assert.Equal(t, true, foundWorkerPodLog, "test logs failed: logs output does not contain logs of the worker pod created in this test")

	// delete test namespace and rhinojob created just now
	execShellCmd("kubectl", []string{"delete", "namespace", testFuncRunNamespace, "--force", "--grace-period=0"})

	// delete the image built just now
	execShellCmd("docker", []string{"rmi", testFuncImageName})
	execShellCmd("sh", []string{"-c", "docker rmi -f $(docker images | grep none | grep second | awk '{print $3}')"})
}
