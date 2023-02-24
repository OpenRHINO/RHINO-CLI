package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var image string
var file string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build MPI function/project",
	Long:  "\nBuild MPI function/project into a docker image",
	Example: `  Build single mpi func:
    rhino build ./hello.cpp --image foo/hello:v1.0
  Build mpi proj(located at root of the folder with Makefile)
    rhino build -i bar/mpibench:v2.1
  Build mpi proj(provide relative path and parameters)
    rhino build -f ./config/Makefile -i bar/mpibench:v2.1 -- make -j all arch=Linux`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(image) == 0 {
			return fmt.Errorf("please provide the image name")
		} else if len(args) == 0 || args[0] == "make" {
			fmt.Println("No source file is specified, use current folder to build")
			path = "./"
		} else {
			path = args[0]
		}
		if err := builder(args, image, path, file); err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(0)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&image, "image", "i", "", "full image form: [registry]/[namespace]/[name]:[tag]")
	buildCmd.Flags().StringVarP(&file, "file", "f", "", "relative path of the makefile")
}
		
func builder(args []string, image string, path string, file string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	var execArgs []string
	var execCommand string
	var buildCommand []string = []string{"make"}
	var makefilePath string
	funcName := getFuncName(image)

	if f.IsDir() {
		if len(file) == 0 {
			makefilePath = filepath.Join(path, "/Makefile")
		} else {
			makefilePath = file
		}
		_, err = os.Stat(makefilePath)
		if err != nil {
			return err
		}
		fmt.Println("Makefile path:", makefilePath)

		if len(args) > 0 {
			buildCommand = args
		}
		fmt.Println("Build command:", buildCommand)		

		_, err = os.Stat("proj.dockerfile") 
		if err != nil {
			fmt.Println("Template not found. Pulling now...")
			templateURL := "https://github.com/OpenRHINO/templates/raw/main/proj.zip"
			if err := downloadTemplate(templateURL, "./"); err != nil {
				fmt.Println("Error:", err.Error())
				os.Exit(0)
			}
		}

		execCommand = "docker"
		execArgs = []string{
			"build", "-t", image,
			"--build-arg", "func_name=" + funcName,
			"--build-arg", "file=" + makefilePath,
			"--build-arg", "make_args=" + strings.Join(buildCommand[1:], " "),
			"-f", "./proj.dockerfile", ".",
		}
	} else {
		_, err = os.Stat("func.dockerfile") 
		if err != nil {
			return fmt.Errorf("build template not found. Please use 'rhino create' first")
		}
		suffix := filepath.Ext(path)
		var compile string
		if suffix == ".c" {
			compile = "mpicc"
		} else if suffix == ".cpp" {
			compile = "mpic++"
		} else {
			return fmt.Errorf("only supports programs written in c or cpp")
		}

		execCommand = "docker"
		execArgs = []string{
			"build", "-t", image,
			"--build-arg", "func_name=" + funcName,
			"--build-arg", "file=" + path,
			"--build-arg", "compile=" + compile,
			"-f", "./func.dockerfile", ".",
		}
	}
	
	cmd := exec.Command(execCommand, execArgs...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		cmdOutput := scanner.Text()
		fmt.Println(cmdOutput)
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func execute(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

func getFuncName(image string) string {
	nameTag := strings.Split(image, "/")
	funcName := strings.Split(nameTag[len(nameTag)-1], ":")[0]
	return funcName
}
