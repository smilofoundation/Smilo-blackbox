package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"Smilo-blackbox/src/data/types"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	utils2 "Smilo-blackbox/src/utils"
)

var (
	peerName, peerKey = utils2.GetMetadata(&types.Peer{})
)

const peerIndexName = "NextPeerIndex"

type DatabaseInstance struct {
	bd  *redis.Client
	log *logrus.Entry
}

func (rds *DatabaseInstance) Close() error {
	return rds.bd.Close()
}

func (rds *DatabaseInstance) Delete(data interface{}) error {
	name, key := utils2.GetMetadata(data)
	value := utils2.GetField(data, key)
	ret := rds.bd.Del(GetKey(name, value))
	err := ret.Err()
	if err == nil && name == peerName {
		ret := rds.bd.ZRem(peerIndexName, value)
		err = ret.Err()
	}
	return err
}

func (rds *DatabaseInstance) Find(fieldname string, value interface{}, to interface{}) error {
	name, key := utils2.GetMetadata(to)
	if key == fieldname {
		ret := rds.bd.Get(GetKey(name, value))
		str, err := ret.Result()
		if err != nil {
			return err
		}
		data := GetTagged(to)
		err = json.Unmarshal([]byte(str), &data)
		if err != nil {
			return err
		}
		GetUntagged(data, to)
		return nil
	}
	return fmt.Errorf("wrong key field %s, expected %s", fieldname, key)
}

func (rds *DatabaseInstance) Save(data interface{}) error {
	name, key := utils2.GetMetadata(data)
	value := utils2.GetField(data, key)
	tagged := GetTagged(data)
	bytesValue, err := json.Marshal(tagged)
	if err != nil {
		return err
	}
    keyValue := GetKey(name, value)
	ret := rds.bd.Set(keyValue, bytesValue, -1)
	err = ret.Err()
	if err == nil && name == peerName {
		score := float64(tagged.(Peer).NextUpdate.Unix())
		ret := rds.bd.ZAdd(peerIndexName, redis.Z{Score:score,Member:keyValue})
		err = ret.Err()
	}
	return err
}

func (rds *DatabaseInstance) AllPeers() (*[]types.Peer, error) {
	var cursor uint64
	keys, cursor, err := rds.bd.Scan(cursor, GetKey(peerName, "*"), 128).Result()
	if err != nil {
		return nil, err
	}
	allPeers := make([]types.Peer, 0, len(keys))
	for _, key := range keys {
		var peer types.Peer
		err := rds.Find(peerKey, GetKeyValue(peerName, key),&peer)
		if err != nil {
			return nil, err
		}
		allPeers = append(allPeers, peer)
	}
	return &allPeers, nil
}

func (rds *DatabaseInstance) GetNextPeer(postpone time.Duration) (*types.Peer, error) {
	ret := rds.bd.ZRange(peerIndexName,0,0)
	err := ret.Err()
	var peer types.Peer
	if err == nil {
		list, err := ret.Result()
		if err == nil && len(list) > 0 {
			err := rds.Find(peerKey, list[0] , &peer)
			if err == nil {
				if peer.NextUpdate.Before(time.Now()) {
					peer.NextUpdate.Add(postpone)
					err = peer.Save()
					if err == nil {
						return &peer, nil
					}
				} else {
					return nil, nil
				}
			}
		}
	}
	return nil, err
}

func DBOpen(filename string, log *logrus.Entry) (*DatabaseInstance, error) {
	byteValue, err := utils2.ReadAllFile(filename, log)
	if err == nil {
		var options redis.Options
		err = json.Unmarshal(byteValue, &options)
		if err == nil {
			db := redis.NewClient(&options)
			_, err = db.Ping().Result()
			if err == nil {
				return &DatabaseInstance{db, log}, nil
			}
			log.WithError(err).Fatal("Unable to connect to redis error: ", err)
		} else {
			log.WithError(err).Fatal("Unable to parse redis config file: ", filename, ", error: ", err)
		}
	} else {
		log.WithError(err).Fatal("Could not open redis config file: ", filename, ", error: ", err)
	}
	return nil, err
}
