package boltdb

import (
	"github.com/asdine/storm"
	"github.com/sirupsen/logrus"
	"os"
)

type BoltdbDatabaseInstance struct {
     bd *storm.DB
     log *logrus.Entry
}

func (bdb *BoltdbDatabaseInstance) Close() error {
	return bdb.bd.Close()
}

func (bdb *BoltdbDatabaseInstance) Delete(data interface{}) error {
	return bdb.bd.DeleteStruct(GetTagged(data))
}

func (bdb *BoltdbDatabaseInstance) Find(fieldname string, value interface{}, to interface{}) error {
	to_ := GetTagged(to)
    err := bdb.bd.One(fieldname, value, to_)
	GetUntagged(to_, to)
    return err
}
func (bdb *BoltdbDatabaseInstance) Save(data interface{}) error {
    return bdb.bd.Save(GetTagged(data))
}
func BoltDBOpen(filename string, log *logrus.Entry) (*BoltdbDatabaseInstance,error) {
	_, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to start DB file at %s", filename)
	}

	log.Info("Opening DB: ", filename)
	db, err := storm.Open(filename)

	if err != nil {
		defer func() {
			err = db.Close()
			log.WithError(err).Fatal("Could not open DBFile: ", filename, ", error: ", err)
			os.Exit(1)
		}()
	}
	bdb := BoltdbDatabaseInstance{db, log}
	return &bdb, nil
}
