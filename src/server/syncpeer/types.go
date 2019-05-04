package syncpeer

import (
	sync2 "sync"
	"time"
)

//PartyInfoRequest used to marshal/unmarshal json
type PartyInfoRequest struct {
	SenderURL   string `json:"url"`
	SenderKey   string `json:"key"`
	SenderNonce string `json:"nonce"`
}

//PartyInfoResponse used to marshal/unmarshal json
type PartyInfoResponse struct {
	PublicKeys []ProvenPublicKey `json:"publicKeys"`
	PeerURLs   []string          `json:"peers"`
}

//ProvenPublicKey used to marshal/unmarshal json
type ProvenPublicKey struct {
	Key   string `json:"key"`
	Proof string `json:"proof"`
}

//Peer used to marshal/unmarshal json
type Peer struct {
	url         string
	publicKeys  [][]byte
	skipcycles  int
	failures    int
	lastFailure time.Time
	tries       int
}

//SafePublicKeyMap used to marshal/unmarshal json
type SafePublicKeyMap struct {
	sync2.RWMutex
	internal map[string]*Peer
}

//NewSafePublicKeyMap create new key
func NewSafePublicKeyMap() *SafePublicKeyMap {
	return &SafePublicKeyMap{
		internal: make(map[string]*Peer),
	}
}

//Get will get internal key
func (spm *SafePublicKeyMap) Get(key string) (value *Peer) {
	spm.RLock()
	defer spm.RUnlock()
	result, _ := spm.internal[key]
	return result
}

//Delete will delete internal key
func (spm *SafePublicKeyMap) Delete(key string) {
	spm.Lock()
	defer spm.Unlock()
	delete(spm.internal, key)
}

//Store will store key
func (spm *SafePublicKeyMap) Store(key string, value *Peer) {
	spm.Lock()
	defer spm.Unlock()
	spm.internal[key] = value
}
