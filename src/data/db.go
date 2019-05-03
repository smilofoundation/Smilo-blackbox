// Copyright 2019 The Smilo-blackbox Authors
// This file is part of the Smilo-blackbox library.
//
// The Smilo-blackbox library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Smilo-blackbox library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Smilo-blackbox library. If not, see <http://www.gnu.org/licenses/>.

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
	_, err := os.Create(dbFile)
	if err != nil {
		log.Fatalf("Failed to start DB file at %s", dbFile)
	}

	log.Info("Opening DB: ", dbFile)
	db, err = storm.Open(dbFile)

	if err != nil {
		defer func(){
			err = db.Close()
			log.WithError(err).Fatal("Could not open DBFile: ", dbFile, ", error: ", err)
			os.Exit(1)
		}()
	}
}
