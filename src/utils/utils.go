package utils

import (
	"flag"
	"os"
	"path"
	"strings"
)

func BuildFilename(filename string) string {
	currentDir, _ := os.Getwd()
	var workDir string
	var newDBFile string
	if flag.Lookup("test.v") != nil {
		isServer := strings.HasSuffix(currentDir, "/server")
		isData := strings.HasSuffix(currentDir, "/data")
		isConfig := strings.HasSuffix(currentDir, "/config")
		isRoot := strings.HasSuffix(currentDir, "/Smilo-blackbox")
		if isServer {
			workDir = "../../"
		} else if isData {
			workDir = "../../"
		} else if isRoot {
			workDir = ""
		} else if isConfig {
			workDir = "../../../"
		}
	} else {
		workDir = ""
	}
	newDBFile = path.Join(currentDir, workDir)
	newDBFile = path.Join(newDBFile, filename)
	return newDBFile
}
