package syncpeer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	sync2 "sync"
	"time"

	"crypto/tls"
	"crypto/x509"

	"github.com/sirupsen/logrus"

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
	hostURL             string
	log                 = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "syncpeer",
	})
	tr = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{RootCAs: getOrCreateCertPool()},
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   getRequestTimeout(),
	}
)

func getRequestTimeout() (t time.Duration) {
	value := os.Getenv("REQUEST_TIMEOUT")
	vint, err := strconv.Atoi(value)
	if err != nil {
		t = 30 * time.Second
	} else {
		t = time.Duration(vint) * time.Second
	}
	return t
}
func getOrCreateCertPool() *x509.CertPool {
	rootCAs, err := x509.SystemCertPool()
	if rootCAs == nil || err != nil {
		rootCAs = x509.NewCertPool()
	}
	return rootCAs
}

//AppendCertificate append cert
func AppendCertificate(cert []byte) bool {
	ok := tr.TLSClientConfig.RootCAs.AppendCertsFromPEM(cert)
	if !ok {
		log.Error("Unable to append additional Root CA certificate.")
	} else {
		client = &http.Client{Transport: tr}
	}
	return ok
}

//StartSync start sync
func StartSync() {
	go sync()
}

//func StopSync() {
//	keepRunning = false
//}
//
//func SetTimeBetweenRequests(seconds int) {
//	timeBetweenRequests = time.Duration(seconds) * time.Second
//}
//
//func SetTimeBetweenCycles(seconds int) {
//	timeBetweenCycles = time.Duration(seconds) * time.Second
//}

//SetHostURL set host
func SetHostURL(url string) {
	hostURL = url
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
			if peer.skipcycles > 0 {
				peer.skipcycles--
				continue
			}
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
		log.Errorf("Unable to query the peer: %s, Error: %s", peerList[i].url, err)
	} else {
		peerList[i].failures = 0
		peerList[i].tries = 0
		peerList[i].publicKeys = publicKeys
		peerList[i].skipcycles = 10 * len(publicKeys)
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

//PeerAdd add peer url
func PeerAdd(url string) {
	if url != hostURL {
		peerChannel <- &Peer{url: url, publicKeys: make([][]byte, 0, 128), failures: 0, tries: 0, skipcycles: 0}
	}
}

//GetPeers get peers
func GetPeers() []string {
	mutex.RLock()
	defer mutex.RUnlock()
	urls := make([]string, 0, len(peerList))
	for _, peer := range peerList {
		urls = append(urls, peer.url)
	}
	return urls
}

//GetPeerURL get url
func GetPeerURL(publicKey []byte) (string, error) {
	peer := publicKeysHashMap.Get(string(publicKey))
	if peer != nil {
		return peer.url, nil
	}
	return "", errors.New("unknown Public Key Peer")
}

//GetPublicKeysFromOtherNode get pub from other nodes
func GetPublicKeysFromOtherNode(url string, publicKey []byte) ([][]byte, []string, error) {
	nonce, err := crypt.NewRandomNonce()
	if err != nil {
		return nil, nil, err
	}
	reqJSON := PartyInfoRequest{SenderURL: hostURL, SenderNonce: base64.StdEncoding.EncodeToString(nonce), SenderKey: base64.StdEncoding.EncodeToString(publicKey)}
	reqStr, err := json.Marshal(reqJSON)
	privateKey := crypt.GetPrivateKey(publicKey)
	var retPubKeys = make([][]byte, 0, 128)
	if err != nil {
		return nil, nil, err
	}
	cli := GetHTTPClient()
	response, err := cli.Post(url+"/partyinfo", "application/json", bytes.NewBuffer(reqStr)) //nolint:bodyclose
	defer func() {
		if response != nil && response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				log.WithError(err).Error("Could not response.Body.Close()")
			}
		}
	}()
	if err != nil {
		return nil, nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, nil, errors.New(response.Status)
	}

	var responseJSON PartyInfoResponse
	p, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(p, &responseJSON)
	if err != nil {
		return nil, nil, err
	}
	for _, provenKey := range responseJSON.PublicKeys {
		remotePublicKey, err := base64.StdEncoding.DecodeString(provenKey.Key)
		if err != nil {
			continue
		}
		remoteProof, err := base64.StdEncoding.DecodeString(provenKey.Proof)
		if err != nil {
			continue
		}
		sharedKey := crypt.ComputeSharedKey(privateKey, remotePublicKey)
		if string(crypt.DecryptPayload(sharedKey, remoteProof, nonce)) != "" {
			retPubKeys = append(retPubKeys, remotePublicKey)
		}
	}
	return retPubKeys, responseJSON.PeerURLs, nil
}

//GetHTTPClient get http client
func GetHTTPClient() *http.Client {
	return client
}
