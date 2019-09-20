package redis

import (
	"fmt"
	"time"

	"Smilo-blackbox/src/data/types"
)

type EncryptedTransaction struct {
	Hash           []byte    `json:"hash"`
	EncodedPayload []byte    `json:"encodedpayload"`
	Timestamp      time.Time `json:"timestamp"`
}

type EncryptedRawTransaction struct {
	Hash           []byte    `json:"hash"`
	EncodedPayload []byte    `json:"encodedpayload"`
	Sender         []byte    `json:"sender"`
	Timestamp      time.Time `json:"timestamp"`
}

type PublicKeyURL struct {
	PublicKey []byte `json:"publickey"`
	URL       string `json:"url"`
}

type Peer struct {
	URL         string    `json:"url"`
	PublicKeys  [][]byte  `json:"publickeys"`
	Failures    int       `json:"failures"`
	LastFailure time.Time `json:"lastfailure"`
	Tries       int       `json:"tries"`
	NextUpdate  time.Time `json:"nextupdate"`
}

func GetKey(name string, value interface{}) string {
	return fmt.Sprintf("%s:%v", name, value)
}

func GetTagged(dat interface{}) interface{} {
	switch dat.(type) {
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
	switch dat.(type) {
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
