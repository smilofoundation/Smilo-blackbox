package data

import (
	"github.com/asdine/storm"
	"os"
	"github.com/golang/glog"
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
    	glog.Fatal(err)
    	os.Exit(1)
	}
}


