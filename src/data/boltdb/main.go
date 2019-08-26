package boltdb

import (
	"os"

	"github.com/asdine/storm"
	"github.com/sirupsen/logrus"
)

type DatabaseInstance struct {
	bd  *storm.DB
	log *logrus.Entry
}

func (bdb *DatabaseInstance) Close() error {
	return bdb.bd.Close()
}

func (bdb *DatabaseInstance) Delete(data interface{}) error {
	return bdb.bd.DeleteStruct(GetTagged(data))
}

func (bdb *DatabaseInstance) Find(fieldname string, value interface{}, to interface{}) error {
	taggedTo := GetTagged(to)
	err := bdb.bd.One(fieldname, value, taggedTo)
	GetUntagged(taggedTo, to)
	return err
}
func (bdb *DatabaseInstance) Save(data interface{}) error {
	return bdb.bd.Save(GetTagged(data))
}
func DbOpen(filename string, log *logrus.Entry) (*DatabaseInstance, error) {
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
	bdb := DatabaseInstance{db, log}
	return &bdb, nil
}
