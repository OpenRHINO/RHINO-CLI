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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testFuncRunNamespace = "rhino-test"

func TestRunSingleJob(t *testing.T) {
	// change work directory to ${workspaceFolder}
	cwd, err := os.Getwd()
	assert.Equal(t, nil, err, "test run failed: %s", errorMessage(err))
	if strings.HasSuffix(cwd, "cmd") {
		os.Chdir("..")
	}
	rootCmd := NewRootCommand()
	// use `rhino build` to build template
	os.Chdir("templates/func")
	testFuncImageName := "test-run-func-cpp:v1"
	rootCmd.SetArgs([]string{"build", "--image", testFuncImageName})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "preparatory work build failed: %s", errorMessage(err))

	// test run command
	execShellCmd("kubectl", []string{"create", "namespace", testFuncRunNamespace})
	rootCmd.SetArgs([]string{"run", testFuncImageName, "--namespace", testFuncRunNamespace})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "test run failed: %s", errorMessage(err))

	// use `kubectl get rhinojob` to check whether rhinojob has been created
	fmt.Println("Wait 60s and check job status")
	time.Sleep(60 * time.Second) //60s is long enough for this simple test job to complete
	cmdOutput, err := execShellCmd("kubectl", []string{"get", "rhinojob", "--namespace", testFuncRunNamespace})
	assert.Equal(t, nil, err, "test run failed: %s", errorMessage(err))
	assert.Equal(t, true, strings.Contains(cmdOutput, "Completed"), "rhinojob failed to start")

	// use defer to ensure that these commands are executed even if the test fails
	defer execShellCmd("sh", []string{"-c", "docker rmi -f $(docker images | grep none | grep second | awk '{print $3}')"})
	defer execShellCmd("docker", []string{"rmi", testFuncImageName})
	defer execShellCmd("kubectl", []string{"delete", "namespace", testFuncRunNamespace, "--force", "--grace-period=0"})
}
