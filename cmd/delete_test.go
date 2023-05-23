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

func TestDeleteSingleJob(t *testing.T) {
	// change work directory to ${workspaceFolder}
	cwd, err := os.Getwd()
	assert.Equal(t, nil, err, "test delete failed: %s", errorMessage(err))
	if strings.HasSuffix(cwd, "cmd") {
		os.Chdir("..")
	}
	rootCmd := NewRootCommand()
	// use `rhino build` to build template
	os.Chdir("templates/func")
	testFuncName := "test-delete-func-cpp"
	testFuncImageName := "test-delete-func-cpp:v1"
	rootCmd.SetArgs([]string{"build", "--image", testFuncImageName})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "preparatory work build failed: %s", errorMessage(err))

	// test run command
	execShellCmd("kubectl", []string{"create", "namespace", testFuncRunNamespace})
	rootCmd.SetArgs([]string{"run", testFuncImageName, "--namespace", testFuncRunNamespace})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "preparatory work run failed: %s", errorMessage(err))

	// test delete
	fmt.Println("Wait 10s and check job status")
	time.Sleep(10 * time.Second)
	testRhinoJobName := testFuncName
	rootCmd.SetArgs([]string{"delete", testRhinoJobName, "--namespace", testFuncRunNamespace})
	err = rootCmd.Execute()
	assert.Equal(t, nil, err, "test delete failed: %s", errorMessage(err))

	// check if the rhinojob created just now is deleted successfully
	actualCmdOutput, err := execShellCmd("kubectl", []string{"get", "rhinojob", "-n", testFuncRunNamespace})
	assert.Equal(t, nil, err, "test delete failed: %s", errorMessage(err))

	expetedCmdOutput := "No resources found in rhino-test namespace.\n"
	assert.Equal(t, expetedCmdOutput, actualCmdOutput, "test delete failed:\n"+
		"expected kubectl output: %s\nactual kubectl output: %s\n",
		expetedCmdOutput, actualCmdOutput)

	// delete the image built just now
	defer execShellCmd("sh", []string{"-c", "docker rmi -f $(docker images | grep none | grep second | awk '{print $3}')"})
	defer execShellCmd("docker", []string{"rmi", testFuncImageName})

	// delete test namespace and rhinojob created just now
	defer execShellCmd("kubectl", []string{"delete", "namespace", testFuncRunNamespace, "--force", "--grace-period=0"})
}
