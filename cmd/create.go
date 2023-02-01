package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var language string
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new mpi function/project",
	Long:  "\nCreate a new mpi function/project",
	Example: `  C function:   rhino create func_name --lang c
  C++ function: rhino create func_name -l cpp
  MPI project:  rhino create proj_name --lang dockerfile`,
	Args: argsCheck,
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&language, "lang", "l", "", "language template to use")
}

func argsCheck(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && len(language) == 0 {
		cmd.Help()
		os.Exit(0)
	} else if len(language) == 0 {
		return fmt.Errorf("language cannot be empty")
	} else if len(args) == 0 {
		return fmt.Errorf("function or project name cannot be empty")
	}
	validName := regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	matchString := validName.MatchString(args[0])
	if !matchString {
		fmt.Println("Error: name can only contain a~z, 0~9 and _")
		os.Exit(0)
	}
	if language != "c" && language != "cpp" && language != "dockerfile" {
		return fmt.Errorf("only support c, cpp and dockerfile")
	}
	return nil
}

func runCreate(cmd *cobra.Command, args []string) error {
	dirName := args[0]
	if _, err := os.Stat(dirName); err == nil {
		fmt.Println("Error: folder", dirName, "already exists")
		os.Exit(0)
	}
	if err := os.Mkdir(dirName, 0700); err != nil {
		fmt.Println("Error: folder", dirName, "could not create")
		os.Exit(0)
	}

	// download the template from github, if fail, delete folder
	templateURL := "https://github.com/OpenRHINO/templates/raw/main/func.zip"
	if err := downloadTemplate(templateURL, dirName); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(0)
	}

	return nil
}

func downloadTemplate(templateURL, dstDir string) error {
	// 下载模板文件包，并暂存为 template.zip (模板文件包应当为zip格式，但可能不叫这个名字)
	resp, err := http.Get(templateURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	zipfile, err := os.Create("template.zip")
	if err != nil {
		return err
	}
	defer zipfile.Close()
	_, err = io.Copy(zipfile, resp.Body)
	if err != nil {
		return err
	}

	// 解压模板文件包
	zr, err := zip.OpenReader("template.zip")
	if err != nil {
		return err
	}
	for _, file := range zr.File {
		path := filepath.Join(dstDir, file.Name)
		// 如果是目录，则创建目录，并跳过当前循环，继续处理下一个
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		// 解压文件到目标文件夹
		fr, err := file.Open()
		if err != nil {
			return err
		}
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		fw.Close()
		fr.Close()
	}
	zr.Close()

	// 删除缓存的模板文件包
	if err := os.Remove("template.zip"); err != nil {
		return err
	}

	return nil
}
