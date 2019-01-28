package sync

import (
	"errors"
	"time"
	"Smilo-blackbox/src/crypt"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"net/http"
	"bytes"
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
	retPublicKeys, err := GetPublicKeysFromOtherNode(url, crypt.GetPublicKeys()[0])
	if err != nil {
		return make([][]byte,0,1)
	}
	return retPublicKeys
}

func PeerAdd(url string) {
	for _, peer := range peerList {
		if peer.url == url {
			return
		}
	}
	peerList = append(peerList, Peer{ url: url, publicKeys:make([][]byte, 0, 128)})
}

func GetPeerURL(publicKey []byte) (string, error) {
	peer := publicKeysHashMap[string(publicKey)]
	if peer != nil {
		return peer.url, nil
	}
	return "", errors.New("Unknow Public Key Peer")
}

func GetPublicKeysFromOtherNode(url string, publicKey []byte) ([][]byte, error) {
	reqJson := PartyInfoRequest{SenderKey:base64.StdEncoding.EncodeToString(publicKey)}
	reqStr, err := json.Marshal(reqJson)
	privateKey := crypt.GetPrivateKey(publicKey)
	var retPubKeys = make([][]byte, 0, 128)
	if err != nil {
		return nil, err
	}
	response, err := new(http.Client).Post(url + "/partyinfo","application/json", bytes.NewBuffer(reqStr))
	if err != nil {
		return nil, err
	}
    if response.StatusCode != http.StatusOK {
    	return nil, errors.New(response.Status)
	}
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	var responseJson PartyInfoResponse
	p, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return nil, err
	}
	err = json.Unmarshal(p, &responseJson)
	if err != nil {
		return nil, err
	}
	for _, provenKey := range responseJson.PublicKeys {
		remotePublicKey, _ := base64.StdEncoding.DecodeString(provenKey.Key)
		remoteProof, _ := base64.StdEncoding.DecodeString(provenKey.Proof)
		sharedKey := crypt.ComputeSharedKey(privateKey, remotePublicKey)
		if string(crypt.DecryptPayload(sharedKey, remoteProof, nil)) != "" {
			retPubKeys = append(retPubKeys, remotePublicKey)
		}
	}
	return retPubKeys, nil
}
