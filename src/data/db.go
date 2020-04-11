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

const BOLTDBENGINE = "boltdb"
const DYNAMODBENGINE = "dynamodb"
const REDISENGINE = "redis"

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
	case BOLTDBENGINE:
		types.DBI, err = boltdb.DbOpen(dbFile, log)
	case DYNAMODBENGINE:
		types.DBI, err = dynamodb.DbOpen(dbFile, log)
	case REDISENGINE:
		types.DBI, err = redis.DBOpen(dbFile, log)
	default:
		panic("Unknown Database Engine")
	}
	if err != nil {
		log.Fatal("Unable to connect to database.")
	}

}

func Migrate(fromEngine string, fromFile string, toEngine string, toFile string) error {
	SetEngine(fromEngine)
	SetFilename(fromFile)
	Start()
	var encryptedTransactions []types.EncryptedTransaction
	var encryptedRawTransactions []types.EncryptedRawTransaction
	var peers []types.Peer

	err := types.GetAll(&encryptedTransactions)
	if err != nil {
		log.WithError(err).Fatal("Unable to get all transactions")
	}

	err = types.GetAll(&encryptedRawTransactions)
	if err != nil {
		log.WithError(err).Fatal("Unable to get all raw transactions")
	}

	err = types.GetAll(&peers)
	if err != nil {
		log.WithError(err).Fatal("Unable to get all peers")
	}

	SetEngine(toEngine)
	SetFilename(toFile)
	Start()

	for _, item := range encryptedTransactions {
		err = item.Save()
		if err != nil {
			log.WithError(err).Fatal("Unable to save all transactions")
		}
	}

	for _, item := range encryptedRawTransactions {
		err = item.Save()
		if err != nil {
			log.WithError(err).Fatal("Unable to save all raw transactions")
		}
	}

	for _, item := range peers {
		err = item.Save()
		if err != nil {
			log.WithError(err).Fatal("Unable to save all peers")
		}
		for _, pk := range item.PublicKeys {
			pkURL := types.NewPublicKeyURL(pk, item.URL)
			err := pkURL.Save()
			if err != nil {
				log.WithError(err).Panic("Could not save peer public key.")
			}
		}
	}
	return nil
}
