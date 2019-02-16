package syncpeer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	sync2 "sync"
	"time"

	"Smilo-blackbox/src/crypt"
)

var (
	peerChannel         = make(chan *Peer, 1024)
	peerList            []*Peer
	publicKeysHashMap   = NewSafePublicKeyMap()
	keepRunning         = true
	mutex               sync2.RWMutex
	timeBetweenCycles   = 13 * time.Second
	timeBetweenRequests = 2 * time.Second
)

func StartSync() {
	go sync()
}

func StopSync() {
	keepRunning = false
}

func SetTimeBetweenRequests(seconds int) {
	timeBetweenRequests = time.Duration(seconds) * time.Second
}

func SetTimeBetweenCycles(seconds int) {
	timeBetweenCycles = time.Duration(seconds) * time.Second
}

func sync() {
	keepRunning = true
	for keepRunning {
		time.Sleep(timeBetweenCycles)
		updateAllPeers()
		updatePeersList()
	}
}

func updatePeersList() {
	mutex.Lock()
	defer mutex.Unlock()
	for j := len(peerList) - 1; j >= 0; j-- {
		if peerList[j].tries > 3 {
			peer := peerList[j]
			if j < len(peerList)-1 {
				peerList = append(peerList[0:j], peerList[j+1:]...)
			} else {
				peerList = peerList[0:j]
			}
			for _, pubKeys := range peer.publicKeys {
				if publicKeysHashMap.Get(string(pubKeys)) == peer {
					publicKeysHashMap.Delete(string(pubKeys))
				}
			}
		}
	}
	for {
		select {
		case p := <-peerChannel:
			alreadyExists := false
			for _, peer := range peerList {
				if peer.url == p.url {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				peerList = append(peerList, p)
			}
		default:
			return
		}
	}
}

func updateAllPeers() {
	mutex.RLock()
	defer mutex.RUnlock()
	for i, peer := range peerList {
		if peer.failures > 10 {
			if time.Since(peer.lastFailure) > (15 * time.Minute) {
				peer.failures = 0
				peer.tries++
			}
		} else {
			time.Sleep(timeBetweenRequests)
			updatePeer(i)
		}
	}
}

func updatePeer(i int) {
	publicKeys, err := queryPeer(peerList[i].url)
	if err != nil {
		peerList[i].failures++
		peerList[i].lastFailure = time.Now()
	} else {
		peerList[i].failures = 0
		peerList[i].tries = 0
		peerList[i].publicKeys = publicKeys
	}
	for j := range publicKeys {
		publicKeysHashMap.Store(string(publicKeys[j]), peerList[i])
	}
}

func queryPeer(url string) ([][]byte, error) {
	pubKeys := crypt.GetPublicKeys()
	if len(pubKeys) == 0 {
		panic("Could find valid public keys, please provide a valid pub key and check your config file")
	}
	retPublicKeys, urls, err := GetPublicKeysFromOtherNode(url, pubKeys[0])
	peerAddAll(urls...)
	if err != nil {
		return make([][]byte, 0, 1), err
	}
	return retPublicKeys, nil
}

func peerAddAll(urls ...string) {
	for _, url := range urls {
		PeerAdd(url)
	}
}

func PeerAdd(url string) {
	peerChannel <- &Peer{url: url, publicKeys: make([][]byte, 0, 128), failures: 0, tries: 0}
}

func GetPeers() []string {
	mutex.RLock()
	defer mutex.RUnlock()
	urls := make([]string, 0, len(peerList))
	for _, peer := range peerList {
		urls = append(urls, peer.url)
	}
	return urls
}

func GetPeerURL(publicKey []byte) (string, error) {
	peer := publicKeysHashMap.Get(string(publicKey))
	if peer != nil {
		return peer.url, nil
	}
	return "", errors.New("unknown Public Key Peer")
}

func GetPublicKeysFromOtherNode(url string, publicKey []byte) ([][]byte, []string, error) {
	reqJson := PartyInfoRequest{SenderKey: base64.StdEncoding.EncodeToString(publicKey)}
	reqStr, err := json.Marshal(reqJson)
	privateKey := crypt.GetPrivateKey(publicKey)
	var retPubKeys = make([][]byte, 0, 128)
	if err != nil {
		return nil, nil, err
	}
	response, err := new(http.Client).Post(url+"/partyinfo", "application/json", bytes.NewBuffer(reqStr))
	if err != nil {
		return nil, nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, nil, errors.New(response.Status)
	}
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	var responseJson PartyInfoResponse
	p, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(p, &responseJson)
	if err != nil {
		return nil, nil, err
	}
	for _, provenKey := range responseJson.PublicKeys {
		remotePublicKey, _ := base64.StdEncoding.DecodeString(provenKey.Key)
		remoteProof, _ := base64.StdEncoding.DecodeString(provenKey.Proof)
		sharedKey := crypt.ComputeSharedKey(privateKey, remotePublicKey)
		if string(crypt.DecryptPayload(sharedKey, remoteProof, nil)) != "" {
			retPubKeys = append(retPubKeys, remotePublicKey)
		}
	}
	return retPubKeys, responseJson.PeerURLs, nil
}
