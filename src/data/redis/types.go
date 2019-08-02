package redis

import (
	"Smilo-blackbox/src/data/types"
	"fmt"
	"time"
)

type EncryptedTransaction struct {
	Hash           []byte `json:"hash"`
	EncodedPayload []byte `json:"encodedpayload"`
	Timestamp      time.Time `json:"timestamp"`
}

func GetKey(name string, value interface{}) string {
	return fmt.Sprintf("%s:%v", name, value)
}

func GetTagged(dat interface{}) interface{} {
	switch dat.(type) {
	case *types.EncryptedTransaction :
		dat2 := EncryptedTransaction(*dat.(*types.EncryptedTransaction))
		return &dat2
	default:
		return dat
	}
}

func GetUntagged(dat interface {}, gen interface{}) {
	switch dat.(type) {
	case *EncryptedTransaction:
		dat2 := types.EncryptedTransaction(*dat.(*EncryptedTransaction))
		*gen.(*types.EncryptedTransaction) = dat2
	default:
		gen = dat
	}
}
