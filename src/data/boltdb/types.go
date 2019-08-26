package boltdb

import (
	"time"

	"Smilo-blackbox/src/data/types"
)

type EncryptedTransaction struct {
	Hash           []byte `storm:"id"`
	EncodedPayload []byte
	Timestamp      time.Time `storm:"index"`
}

type EncryptedRawTransaction struct {
	Hash           []byte `storm:"id"`
	EncodedPayload []byte
	Sender         []byte
	Timestamp      time.Time `storm:"index"`
}

func GetTagged(dat interface{}) interface{} {
	switch dat.(type) {
	case *types.EncryptedTransaction:
		dat2 := EncryptedTransaction(*dat.(*types.EncryptedTransaction))
		return &dat2
	case *types.EncryptedRawTransaction:
		dat2 := EncryptedRawTransaction(*dat.(*types.EncryptedRawTransaction))
		return &dat2
	default:
		return dat
	}
}

func GetUntagged(dat interface{}, gen interface{}) {
	switch dat.(type) {
	case *EncryptedTransaction:
		dat2 := types.EncryptedTransaction(*dat.(*EncryptedTransaction))
		*gen.(*types.EncryptedTransaction) = dat2
	case *EncryptedRawTransaction:
		dat2 := types.EncryptedRawTransaction(*dat.(*EncryptedRawTransaction))
		*gen.(*types.EncryptedRawTransaction) = dat2
	}
}
