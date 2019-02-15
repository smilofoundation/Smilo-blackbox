package syncpeer

import (
	sync2 "sync"
	"time"
)

type PartyInfoRequest struct {
	SenderURL string `json:url`
	SenderKey string `json:"key"`
}

type PartyInfoResponse struct {
	PublicKeys []ProvenPublicKey `json:"publicKeys"`
	PeerURLs   []string          `json:peers`
}

type ProvenPublicKey struct {
	Key   string `json:"key"`
	Proof string `json:"proof"`
}

type Peer struct {
	url         string
	publicKeys  [][]byte
	failures    int
	lastFailure time.Time
	tries       int
}

type SafePublicKeyMap struct {
	sync2.RWMutex
	internal map[string]*Peer
}

func NewSafePublicKeyMap() *SafePublicKeyMap {
	return &SafePublicKeyMap{
		internal: make(map[string]*Peer),
	}
}

func (spm *SafePublicKeyMap) Get(key string) (value *Peer) {
	spm.RLock()
	defer spm.RUnlock()
	result, _ := spm.internal[key]
	return result
}

func (spm *SafePublicKeyMap) Delete(key string) {
	spm.Lock()
	defer spm.Unlock()
	delete(spm.internal, key)
}

func (spm *SafePublicKeyMap) Store(key string, value *Peer) {
	spm.Lock()
	defer spm.Unlock()
	spm.internal[key] = value
}
