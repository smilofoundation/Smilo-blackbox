package boltdb

import (
	"Smilo-blackbox/src/data/types"
	"time"
)

type EncryptedTransaction struct {
	Hash           []byte `storm:"id"`
	EncodedPayload []byte
	Timestamp      time.Time `storm:"index"`
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
