package data

import (
	"os"

	"github.com/asdine/storm"

	"Smilo-blackbox/src/server/config"
	"fmt"
	"strings"
	"path"
)

var db *storm.DB

func init() {
	Start()
}

func Start() {
	var err error
	currentDir, _ := os.Getwd()
	var workDir string
	var newDBFile string
	var dbFile = config.DBFile.Value

	isServer := strings.HasSuffix(currentDir, "/server")
	isData := strings.HasSuffix(currentDir, "/data")
	isRoot := strings.HasSuffix(currentDir, "/Smilo-blackbox")
	if isServer {
		workDir = "../../"
		fmt.Println("db, Contains /server")
	} else if isData {
		workDir = "../../"
		fmt.Println("db, Contains /data")
	} else if isRoot {
		fmt.Println("db, is root dir")
		workDir = ""
	}

	newDBFile = path.Join(currentDir, workDir)

	newDBFile = path.Join(newDBFile, dbFile)

	fmt.Println("Opening DB: ", newDBFile)
	db, err = storm.Open(newDBFile)

	if err != nil {
		defer db.Close()
		log.Fatal("Could not open DBFile: ", dbFile, ", error: ", err)
		os.Exit(1)
	}
}
