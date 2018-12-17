package data

import (
	"os"

	"github.com/asdine/storm"

	"Smilo-blackbox/src/server/config"
)

var db *storm.DB

func init() {
	Start()
}

func Start() {
	var err error
	db, err = storm.Open(config.DBFile.Value)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
