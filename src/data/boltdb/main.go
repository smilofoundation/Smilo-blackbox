package boltdb

import (
	"Smilo-blackbox/src/data/types"
	"os"
	"time"

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

func (bdb *DatabaseInstance) AllPeers () (*[]types.Peer, error) {
	var peers []Peer
	err := bdb.bd.All(&peers)
	if err != nil {
		return nil, err
	}
	allPeers := make([]types.Peer, 0, len(peers))
	for _,peer := range peers {
		tmp := &types.Peer{}
		GetUntagged(peer, tmp)
		allPeers = append(allPeers, *tmp)
	}
	return &allPeers, nil
}
func (bdb *DatabaseInstance) Save(data interface{}) error {
	return bdb.bd.Save(GetTagged(data))
}

func (bdb *DatabaseInstance) GetNextPeer(postpone time.Duration) (*types.Peer, error) {
	var nextValues []Peer
	err := bdb.bd.Range("NextUpdate",time.Unix(0,0), time.Now(), &nextValues, storm.Limit(1))
    if err != nil {
    	return nil, err
	}
	if len(nextValues) == 0 {
		return nil, nil
	}
	var next *types.Peer
	GetUntagged(nextValues[0], next)
	nextValues[0].NextUpdate = time.Now().Add(postpone)
	err = bdb.bd.Save(nextValues[0])
	if err != nil {
		return nil, err
	}
	return next, nil
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
