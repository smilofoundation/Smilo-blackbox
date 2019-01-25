package data

import (
	"github.com/asdine/storm"

	"Smilo-blackbox/src/server/config"
	"os"
	"path"
)

var db *storm.DB

func Start(dbFile string) {

	var err error
	if dbFile == "" {
		dbFile = config.GetString(config.DBFileStr)

		dbFile = path.Join(config.WorkDir, dbFile)
	}

	db, err = storm.Open(dbFile)
	if err != nil {
		defer db.Close()
		log.Fatal("Could not open DBFile: ", dbFile, ", error: ", err)
		os.Exit(1)
	}
}
