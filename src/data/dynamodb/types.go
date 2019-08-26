package dynamodb

import (
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	dynDB "github.com/aws/aws-sdk-go/service/dynamodb"
	dynDBAttr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	utils2 "Smilo-blackbox/src/utils"
)

func GetDeleteItemInput(data interface{}) (*dynDB.DeleteItemInput, error) {
	keys, err := dynDBAttr.MarshalMap(data)
	if err != nil {
		return nil, err
	}
	t := utils2.GetType(data)
	for key := range keys {

		field, ok := t.FieldByName(key)
		if ok && field.Tag.Get("key") == "true" {
			continue
		}
		delete(keys, key)
	}
	input := &dynDB.DeleteItemInput{
		Key:       keys,
		TableName: getTablename(data),
	}
	return input, nil
}

func GetPutItemInput(data interface{}) (*dynDB.PutItemInput, error) {
	keys, err := dynDBAttr.MarshalMap(data)
	if err != nil {
		return nil, err
	}
	input := &dynDB.PutItemInput{
		Item:                   keys,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              getTablename(data),
	}
	return input, nil
}

func GetItemInput(fieldname string, value interface{}, to interface{}) (*dynDB.GetItemInput, error) {
	av, err := dynDBAttr.Marshal(value)
	if err != nil {
		return nil, err
	}
	keys := map[string]*dynDB.AttributeValue{
		fieldname: av,
	}
	input := &dynDB.GetItemInput{
		Key:       keys,
		TableName: getTablename(to),
	}
	return input, err
}

func getTablename(obj interface{}) *string {
	s := reflect.TypeOf(obj).String()
	return aws.String(strings.TrimPrefix(s, "*types."))
}
