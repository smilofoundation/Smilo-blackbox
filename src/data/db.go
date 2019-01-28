package data

import (
	"os"

	"github.com/asdine/storm"

	"path"
	"strings"

)

var db *storm.DB

var dbFile string

func SetFilename(filename string) {
	dbFile = filename
}

func Start() {
	var err error
	currentDir, _ := os.Getwd()
	var workDir string
	var newDBFile string

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

	newDBFile = path.Join(currentDir, workDir)

	newDBFile = path.Join(newDBFile, dbFile)

	log.Info("Opening DB: ", newDBFile)
	db, err = storm.Open(newDBFile)

	if err != nil {
		defer db.Close()
		log.Fatal("Could not open DBFile: ", dbFile, ", error: ", err)
		os.Exit(1)
	}
}
