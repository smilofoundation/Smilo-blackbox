package dynamodb

import (
	"Smilo-blackbox/src/data/types"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	dynDB "github.com/aws/aws-sdk-go/service/dynamodb"
	dynDBAttr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type DatabaseInstance struct {
	db  *dynDB.DynamoDB
	log *logrus.Entry
}

func (dyndb *DatabaseInstance) Close() error {
	dyndb.db = nil
	return nil
}

func (dyndb *DatabaseInstance) Delete(data interface{}) error {
	item, err := GetDeleteItemInput(data)
	if err != nil {
		return err
	}
	_, err = dyndb.db.DeleteItem(item)
	return err
}

func (dyndb *DatabaseInstance) Find(fieldname string, value interface{}, to interface{}) error {
	input, err := GetItemInput(fieldname, value, to)
	if err != nil {
		dyndb.log.Error("Unable to convert to a Dynamo DB type")
		return err
	}
	ret, err := dyndb.db.GetItem(input)
	if err == nil {
		if ret.Item == nil {
			return errors.New("not found")
		}
		err = dynDBAttr.UnmarshalMap(ret.Item, to)
	}
	return err
}
func (dyndb *DatabaseInstance) Save(data interface{}) error {
	item, err := GetPutItemInput(data)
	if err != nil {
		return err
	}
	_, err = dyndb.db.PutItem(item)
	return err
}

func (dyndb *DatabaseInstance) AllPeers () (*[]types.Peer, error) {
	return nil, nil
}
func (dyndb *DatabaseInstance) GetNextPeer(postpone time.Duration) (*types.Peer, error) {
    //TODO: Implement NextPeer for DynamoDB
	return nil, nil
}

func DbOpen(filename string, log *logrus.Entry) (*DatabaseInstance, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db := dynDB.New(sess)
	log.Info("Opening DB: ", filename)

	bdb := DatabaseInstance{db, log}
	return &bdb, nil
}
