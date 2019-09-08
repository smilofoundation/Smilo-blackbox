package types

import "time"

var DBI DatabaseInstance

type DatabaseInstance interface {
	Close() error
	Delete(data interface{}) error
	Find(fieldname string, value interface{}, to interface{}) error
	AllPeers() (*[]Peer, error)
	Save(data interface{}) error
	GetNextPeer(pospone time.Duration) (*Peer, error)
}
