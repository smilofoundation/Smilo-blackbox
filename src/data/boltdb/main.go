package boltdb

import (
	"os"
	"reflect"
	"sync"
	"time"

	"Smilo-blackbox/src/data/types"

	"github.com/asdine/storm"
	"github.com/sirupsen/logrus"
)

var mutex = sync.Mutex{}

type DatabaseInstance struct {
	bd  *storm.DB
	log *logrus.Entry
}

func (bdb *DatabaseInstance) Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	return bdb.bd.Close()
}

func (bdb *DatabaseInstance) Delete(data interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	return bdb.bd.DeleteStruct(GetTagged(data))
}

func (bdb *DatabaseInstance) Find(fieldname string, value interface{}, to interface{}) error {
	taggedTo := GetTagged(to)
	mutex.Lock()
	locked := true
	defer func() {
		if locked {
			mutex.Unlock()
		}
	}()
	err := bdb.bd.One(fieldname, value, taggedTo)
	mutex.Unlock()
	locked = false
	if err == storm.ErrNotFound {
		return types.ErrNotFound
	}
	GetUntagged(taggedTo, to)
	return err
}

func (bdb *DatabaseInstance) All(instances interface{}) error {
	result := reflect.ValueOf(instances)
	resultItem := reflect.New(reflect.TypeOf(result.Elem().Interface()).Elem()).Elem().Addr().Interface()
	request := GetTaggedArray(instances)
	mutex.Lock()
	locked := true
	defer func() {
		if locked {
			mutex.Unlock()
		}
	}()
	err := bdb.bd.All(request)
	mutex.Unlock()
	locked = false
	if err != nil {
		return err
	}
	result = reflect.ValueOf(
		reflect.MakeSlice(
			reflect.SliceOf(
				reflect.TypeOf(resultItem).Elem()), 0, reflect.ValueOf(request).Elem().Len()).
			Interface())
	for i := 0; i < reflect.ValueOf(request).Elem().Len(); i++ {
		requestItem := reflect.ValueOf(request).Elem().Index(i).Addr().Interface()
		GetUntagged(requestItem, resultItem)
		tmp2 := reflect.ValueOf(resultItem)
		result = reflect.Append(result, tmp2.Elem())
	}
	types.GetUntaggedArrayPtr(result.Interface(), instances)
	return nil
}

func (bdb *DatabaseInstance) AllPeers() (*[]types.Peer, error) {
	var peers []Peer
	mutex.Lock()
	locked := true
	defer func() {
		if locked {
			mutex.Unlock()
		}
	}()
	err := bdb.bd.All(&peers)
	mutex.Unlock()
	locked = false
	if err != nil {
		return nil, err
	}
	allPeers := make([]types.Peer, 0, len(peers))
	for i := range peers {
		tmp := types.Peer{}
		GetUntagged(&peers[i], &tmp)
		allPeers = append(allPeers, tmp)
	}
	return &allPeers, nil
}

func (bdb *DatabaseInstance) Save(data interface{}) error {
	tagged := GetTagged(data)
	mutex.Lock()
	defer mutex.Unlock()
	err := bdb.bd.Save(tagged)
	if err != nil && err == storm.ErrAlreadyExists {
		err = bdb.bd.Update(tagged)
	}
	return err
}

func (bdb *DatabaseInstance) GetNextPeer(postpone time.Duration) (*types.Peer, error) {
	var nextValues []Peer
	mutex.Lock()
	locked := true
	defer func() {
		if locked {
			mutex.Unlock()
		}
	}()
	err := bdb.bd.Range("NextUpdate", time.Unix(0, 0), time.Now(), &nextValues, storm.Limit(10))
	mutex.Unlock()
	locked = false
	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}
	if len(nextValues) == 0 {
		return nil, nil
	}
	var next types.Peer
	GetUntagged(&nextValues[0], &next)
	nextValues[0].NextUpdate = time.Now().Add(postpone)
	mutex.Lock()
	locked = true
	err = bdb.bd.Update(&nextValues[0])
	if err != nil {
		return nil, err
	}
	return &next, nil
}

func DBOpen(filename string, log *logrus.Entry) (*DatabaseInstance, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to start DB file at %s", filename)
		}
	}

	log.Info("Opening DB: ", filename)
	db, err := storm.Open(filename)

	if err != nil {
		defer func() {
			if db != nil {
				_ = db.Close()
			}
			log.WithError(err).Fatal("Could not open DBFile: ", filename, ", error: ", err)
			os.Exit(1)
		}()
	}
	bdb := DatabaseInstance{db, log}
	return &bdb, nil
}
