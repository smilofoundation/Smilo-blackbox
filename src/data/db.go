package data

import (
	"os"

	"github.com/asdine/storm"

	"path"

	"Smilo-blackbox/src/server/config"
)

var db *storm.DB

func init() {
	Start()
}

func Start() {
	var err error
	DBFile := config.GetString(config.DBFileStr)

	dbPATH := path.Join(config.WorkDir, DBFile)
	db, err = storm.Open(dbPATH)

	if err != nil {
		log.Fatal("Could not open DBFile: ", dbPATH, ", error: ", err)
		os.Exit(1)
	}
}
