package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	utils2 "Smilo-blackbox/src/utils"
)

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
	return ret.Err()
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

	ret := rds.bd.Set(GetKey(name, value), bytesValue, -1)
	return ret.Err()
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
