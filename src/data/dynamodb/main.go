package dynamodb

import (
	"time"

	"Smilo-blackbox/src/data/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynDB "github.com/aws/aws-sdk-go/service/dynamodb"
	dynDBAttr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	dynDBExp "github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
			return types.ErrNotFound
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

func (dyndb *DatabaseInstance) AllPeers() (*[]types.Peer, error) {
	peerList := make([]types.Peer, 0, 128)
	input := &dynDB.ScanInput{
		TableName: aws.String(*getTablename(&types.Peer{})),
	}
	out, err := dyndb.db.Scan(input)
	if err != nil {
		return nil, err
	}
	for _, item := range out.Items {
		var peer types.Peer
		err = dynDBAttr.UnmarshalMap(item, &peer)
		if err != nil {
			return nil, err
		}
		peerList = append(peerList, peer)
	}
	return &peerList, nil
}
func (dyndb *DatabaseInstance) GetNextPeer(postpone time.Duration) (*types.Peer, error) {
	cond := dynDBExp.Name("NextUpdate").LessThanEqual(dynDBExp.Value(time.Now()))
	expr, err := dynDBExp.NewBuilder().
		WithFilter(cond).
		Build()
	if err != nil {
		return nil, err
	}
	input := &dynDB.ScanInput{
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(*getTablename(&types.Peer{})),
	}
	out, err := dyndb.db.Scan(input)
	if err != nil {
		return nil, err
	}
	if *out.Count < 1 {
		return nil, nil
	}
	var peer types.Peer
	err = dynDBAttr.UnmarshalMap(out.Items[0], &peer)
	if err == nil {
		peer.NextUpdate = time.Now().Add(postpone)
		err = peer.Save()
		if err == nil {
			return &peer, nil
		}
	}
	return nil, err
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
