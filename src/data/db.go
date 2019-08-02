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
	"Smilo-blackbox/src/data/boltdb"
	"Smilo-blackbox/src/data/dynamodb"
	"Smilo-blackbox/src/data/redis"
	"Smilo-blackbox/src/data/types"
)


var dbFile string
var dbEngine = ""


// SetFilename set filename
func SetFilename(filename string) {
	dbFile = filename
}

func SetEngine(engine string) {
	dbEngine = engine
}
// Start will start the db
func Start() {
	var err error
	switch dbEngine {
		case "boltdb":
			types.DBI, err = boltdb.BoltDBOpen(dbFile, log)
		case "dynamodb":
			types.DBI, err = dynamodb.DynamoDBOpen(dbFile, log)
		case "redis":
		    types.DBI, err = redis.RedisOpen(dbFile, log)
		default:
		    panic("Unknown Database Engine")
	}
	if err != nil {
        log.Fatal("Unable to connect to database.")
	}

}
