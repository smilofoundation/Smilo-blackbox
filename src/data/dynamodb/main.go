package dynamodb

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	dynDB "github.com/aws/aws-sdk-go/service/dynamodb"
	dynDBAttr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type DynamodbDatabaseInstance struct {
	db *dynDB.DynamoDB
	log *logrus.Entry
}

func (dyndb *DynamodbDatabaseInstance) Close() error {
	dyndb.db=nil
	return nil
}

func (dyndb *DynamodbDatabaseInstance) Delete(data interface{}) error {
	item, err := GetDeleteItemInput(data)
	if err != nil {
		return err
	}
	_,err = dyndb.db.DeleteItem(item)
	return err
}

func (dyndb *DynamodbDatabaseInstance) Find(fieldname string, value interface{}, to interface{}) error {
	input, err := GetItemInput(fieldname, value, to)
	if err != nil {
		dyndb.log.Error("Unable to convert to a Dynamo DB type")
		return err
	}
	ret, err := dyndb.db.GetItem(input)
	if err == nil {
		if ret.Item == nil {
			return errors.New("Not found")
		}
		err = dynDBAttr.UnmarshalMap(ret.Item, to)
	}
	return err
}
func (dyndb *DynamodbDatabaseInstance) Save(data interface{}) error {
	item, err := GetPutItemInput(data)
	if err != nil {
		return err
	}
	_,err = dyndb.db.PutItem(item)
	return err
}
func DynamoDBOpen(filename string, log *logrus.Entry) (*DynamodbDatabaseInstance,error) {

	//sess, err := session.NewSession()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//if err != nil {
	//	return nil, err
	//}
	db := dynDB.New(sess)
	log.Info("Opening DB: ", filename)
//	db, err := storm.Open(filename)

	//if err != nil {
	//	defer func() {
	//		err = db.Close()
	//		log.WithError(err).Fatal("Could not open DBFile: ", filename, ", error: ", err)
	//		os.Exit(1)
	//	}()
//	}
	bdb := DynamodbDatabaseInstance{db, log}
	return &bdb, nil
}
