package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Whitelist = []string{"logs/", "bin/", "themes/", "config.toml", "sqlite/", ".cupignore"}

const REPO_NAME = "Contagio"

const (
	INFO_MSG  = "\x1b[42m\x1b[44m[INFO]\x1b[0m"
	ERROR_MSG = "\x1b[42m\x1b[41m[INFO]\x1b[0m"
	INPUT_MSG = "\x1b[42m\x1b[100m[INPUT]\x1b[0m"
)

type ResponseVersion struct {
	Tag_name    string `json:"tag_name"`
	Tarball_url string `json:"tarball_url"`
}

func parseConfig() {
	config, err := os.ReadFile(".cupignore")
	if err != nil {
		fmt.Printf("%s \".cupignore\" not found. Skipping...\n", INFO_MSG)
		return
	}

	for _, i := range strings.Split(string(config), "\n") {
		if len(i) == 0 || strings.HasPrefix(i, "//") || i == " " || i == "" {
			continue
		}

		Whitelist = append(Whitelist, i)
	}

}

func backup(source, destination, backupName string) {
	filepath.Walk(source, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, _ := filepath.Rel(source, sourcePath)
		destinationPath := filepath.Join(destination, relativePath)

		if info.IsDir() {
			os.MkdirAll(destinationPath, os.ModePerm)
		} else {
			if filepath.Base(sourcePath) == backupName {
				return nil
			}

			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			destFile, err := os.Create(destinationPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, sourceFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func getCurrentVersion() string {
	version, err := os.ReadFile(".VERSION")
	if err != nil || string(version) == "" {
		fmt.Printf("%s \".VERSION\" not found\n", ERROR_MSG)
		return ""
	}

	return string(version)
}

func downloadVersion(old_version, new_version string) {
	var response []byte
	var response_struct ResponseVersion

	if new_version == "" {
		r, err := http.Get("https://api.github.com/repos/TryZeroOne/" + REPO_NAME + "/releases/latest")
		if err != nil {
			fmt.Printf("%s HTTP get error: "+err.Error(), ERROR_MSG)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("%s ReadAll error: "+err.Error(), ERROR_MSG)
			return
		}
		response = body

	} else {
		r, err := http.Get("https://api.github.com/repos/TryZeroOne/" + REPO_NAME + "/releases/tags/" + new_version)
		if err != nil {
			fmt.Printf("%s HTTP get error: "+err.Error(), ERROR_MSG)
			return
		}

		if r.StatusCode == 404 {
			fmt.Printf("%s Version ( "+new_version+" ) not found. See https://github.com/TryZeroOne/"+REPO_NAME+"/releases\n", ERROR_MSG)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("%s ReadAll error: "+err.Error(), ERROR_MSG)
			return
		}

		response = body
	}

	err := json.Unmarshal([]byte(response), &response_struct)
	if err != nil {
		fmt.Printf("%s Unmarshal error: "+err.Error(), ERROR_MSG)
	}

	if response_struct.Tag_name == old_version {
		fmt.Printf("%s The version is already installed\n", ERROR_MSG)
		return
	}

	new_version = response_struct.Tag_name

	for {
		var input string
		fmt.Printf("%s Version found ( "+new_version+" )! Install? [y/n]:\n", INPUT_MSG)
		fmt.Scanln(&input)

		if strings.ToLower(input) != "y" && strings.ToLower(input) != "n" {
			continue
		}

		if strings.ToLower(input) != "y" {
			os.Exit(0)
		}

		break
	}

	backup_name := "backup_" + old_version
	src, _ := os.Getwd()
	dst := "./" + backup_name

	os.RemoveAll(dst)
	backup(src, dst, backup_name)
	fmt.Printf("%s Backup has been created with name \"%s\"\n", INFO_MSG, backup_name)
	fmt.Printf("%s Downloading version %s ...\n", INFO_MSG, new_version)
	os.RemoveAll("./tmp")
	os.Mkdir("./tmp", os.ModePerm)

	r, err := http.Get(response_struct.Tarball_url)
	if err != nil {
		fmt.Printf("%s HTTP get error: "+err.Error(), ERROR_MSG)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("%s ReadAll error: "+err.Error(), ERROR_MSG)
		return
	}

	os.WriteFile("./tmp/res.tar.gz", body, os.ModePerm)

	exec.Command("tar", "-zxf", "tmp/res.tar.gz", "-C", "tmp/").Run()

	dir, err := os.ReadDir("./tmp")
	if err != nil {
		fmt.Printf("%s ReadDir error: "+err.Error(), ERROR_MSG)
		return
	}

	var new_dir string

	for _, i := range dir {
		if strings.HasPrefix(i.Name(), "TryZeroOne") {
			new_dir = i.Name()
			for _, i := range Whitelist {
				os.RemoveAll("./tmp/" + new_dir + "/" + i)
			}
		}
	}

	cur_dir, err := os.ReadDir("./")
	if err != nil {
		fmt.Printf("%s ReadDir error: "+err.Error(), ERROR_MSG)
		return
	}

	for _, x := range cur_dir {
		if strings.HasPrefix(x.Name(), "tmp") || strings.HasPrefix(x.Name(), "backup") {
			continue
		}

		status := func() bool {
			for _, i := range Whitelist {
				if strings.HasPrefix(i, x.Name()) {
					return true
				}
			}
			return false
		}

		if status() {
			continue
		}

		os.RemoveAll(x.Name())
	}

	dir, err = os.ReadDir("./tmp/" + new_dir)
	if err != nil {
		fmt.Printf("%s ReadDir error: "+err.Error(), ERROR_MSG)
		return
	}
	for _, i := range dir {
		fmt.Println("./tmp/" + new_dir + "/" + i.Name() + " ./")
		exec.Command("cp", "-r", "./tmp/"+new_dir+"/"+i.Name(), "./").Run()
	}

	os.WriteFile(".VERSION", []byte(new_version), os.ModePerm)

}

func main() {
	var new_version string

	if len(os.Args) > 2 {
		if os.Args[1] == "--version" {
			if len(os.Args[2]) < 1 {
				fmt.Printf("%s Invalid args\n", ERROR_MSG)
				return
			}
			new_version = os.Args[2]
		}
	}

	parseConfig()
	old_version := getCurrentVersion()
	if old_version == "" {
		return
	}

	downloadVersion(old_version, new_version)

	fmt.Printf("%s New version installed. Enjoy!\n", INFO_MSG)
}
