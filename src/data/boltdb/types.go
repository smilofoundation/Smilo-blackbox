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

type PublicKeyURL struct {
	PublicKey []byte `storm:"id"`
	URL       string `storm:"index"`
}

type Peer struct {
	URL         string `storm:"id"`
	PublicKeys  [][]byte
	Failures    int
	LastFailure time.Time
	Tries       int
	NextUpdate  time.Time `storm:"index"`
}

func GetTagged(dat interface{}) interface{} {
	switch dat.(type) { //nolint
	case *types.EncryptedTransaction:
		dat2 := EncryptedTransaction(*dat.(*types.EncryptedTransaction))
		return &dat2
	case *types.EncryptedRawTransaction:
		dat2 := EncryptedRawTransaction(*dat.(*types.EncryptedRawTransaction))
		return &dat2
	case *types.PublicKeyURL:
		dat2 := PublicKeyURL(*dat.(*types.PublicKeyURL))
		return &dat2
	case *types.Peer:
		dat2 := Peer(*dat.(*types.Peer))
		return &dat2
	default:
		return dat
	}
}

func GetUntagged(dat interface{}, gen interface{}) {
	switch dat.(type) { //nolint
	case *EncryptedTransaction:
		dat2 := types.EncryptedTransaction(*dat.(*EncryptedTransaction))
		*gen.(*types.EncryptedTransaction) = dat2
	case *EncryptedRawTransaction:
		dat2 := types.EncryptedRawTransaction(*dat.(*EncryptedRawTransaction))
		*gen.(*types.EncryptedRawTransaction) = dat2
	case *PublicKeyURL:
		dat2 := types.PublicKeyURL(*dat.(*PublicKeyURL))
		*gen.(*types.PublicKeyURL) = dat2
	case *Peer:
		dat2 := types.Peer(*dat.(*Peer))
		*gen.(*types.Peer) = dat2
	}
}

func GetTaggedArray(dat interface{}) interface{} {
	switch dat.(type) { //nolint
	case *[]types.EncryptedTransaction:
		return &[]EncryptedTransaction{}
	case *[]types.EncryptedRawTransaction:
		return &[]EncryptedRawTransaction{}
	case *[]types.PublicKeyURL:
		return &[]PublicKeyURL{}
	case *[]types.Peer:
		return &[]Peer{}
	default:
		return dat
	}
}
