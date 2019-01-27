package sync

import (
	"errors"
	"time"
)

var peerList []Peer

var publicKeysHashMap map[string]*Peer

var keepRunning = true

type Peer struct {
	url string
	publicKeys [][]byte
}

func StartSync() {
    go sync()
}

func StopSync() {
	keepRunning = false
}

func sync() {
	keepRunning = true
	for keepRunning {
		time.Sleep(15*time.Second)
		for i, _ := range peerList {
			updatePeer(i)
		}
	}
}

func updatePeer(i int) {
	publicKeys := queryPeer(peerList[i].url)
	peerList[i].publicKeys = publicKeys
	for j, _ := range publicKeys {
		publicKeysHashMap[string(publicKeys[j])] = &peerList[i]
	}
}

func queryPeer(url string) ([][]byte) {
	// TODO: query PeerInfo
	return make([][]byte,0,1)
}

func PeerAdd(url string) {
	peerList = append(peerList, Peer{ url: url, publicKeys:make([][]byte, 0, 128)})
}

func GetPeerURL(publicKey []byte) (string, error) {
	peer := publicKeysHashMap[string(publicKey)]
	if peer != nil {
		return peer.url, nil
	}
	return "", errors.New("Unknow Public Key Peer")
}
