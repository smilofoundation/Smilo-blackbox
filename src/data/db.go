package data

import (
	"os"

	"github.com/asdine/storm"

)

var db *storm.DB

var dbFile string

func SetFilename(filename string) {
	dbFile = filename
}

func Start() {
	var err error

	log.Info("Opening DB: ", dbFile)
	db, err = storm.Open(dbFile)

	if err != nil {
		defer db.Close()
		log.Fatal("Could not open DBFile: ", dbFile, ", error: ", err)
		os.Exit(1)
	}
}

