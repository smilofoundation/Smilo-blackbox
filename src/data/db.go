package data

import (
	"os"

	"github.com/asdine/storm"
)

var databaseFile string
var db *storm.DB

func Start(_databaseFile string) {
	if _databaseFile != "" {
		databaseFile = _databaseFile
	} else {
		databaseFile = "blackbox.db"
	}
	var err error
	db, err = storm.Open(databaseFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
