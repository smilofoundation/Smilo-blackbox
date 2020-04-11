package types

import (
	"errors"
	"time"
)

var DBI DatabaseInstance

type DatabaseInstance interface {
	Close() error
	Delete(data interface{}) error
	Find(fieldname string, value interface{}, to interface{}) error
	All(instances interface{}) error
	AllPeers() (*[]Peer, error)
	Save(data interface{}) error
	GetNextPeer(pospone time.Duration) (*Peer, error)
}

var ErrNotFound = errors.New("not found")

func GetAll(instances interface{}) error {
	return DBI.All(instances)
}

func GetUntaggedArrayPtr(dat interface{}, gen interface{}) {
	switch dat.(type) {
	case []EncryptedTransaction:
		*gen.(*[]EncryptedTransaction) = dat.([]EncryptedTransaction)
	case []EncryptedRawTransaction:
		*gen.(*[]EncryptedRawTransaction) = dat.([]EncryptedRawTransaction)
	case []PublicKeyURL:
		*gen.(*[]PublicKeyURL) = dat.([]PublicKeyURL)
	case []Peer:
		*gen.(*[]Peer) = dat.([]Peer)
	}
}
