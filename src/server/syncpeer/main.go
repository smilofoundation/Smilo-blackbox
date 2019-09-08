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
	"time"

	"crypto/tls"
	"crypto/x509"

	"github.com/sirupsen/logrus"

	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/data/types"
)

var (
	//peerChannel         = make(chan *Peer, 1024)
	//peerList            []*Peer
	keepRunning         = true
	//mutex               sync2.RWMutex
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
		updatePeer()
	}
}

func updatePeer() {
	peer, err := types.FindNextUpdatablePeer(2*client.Timeout)
	if err != nil {
		//panic
	}
	if peer.Failures > 10 {
		if time.Since(peer.LastFailure) > (15 * time.Minute) {
			peer.Failures = 0
			peer.Tries++
		}
	} else {
		if peer.SkipCycles > 0 {
			peer.SkipCycles--
		} else {
			updateFromRemotePeerData(peer)
		}
	}
	if peer.Tries > 3 {
		for _, pubKeys := range peer.PublicKeys {
			pkurl, err := types.FindPublicKeyUrl(pubKeys)
			if err != nil {
				//panic
			}
			if pkurl.URL == peer.URL {
				pkurl.Delete()
			}
		}
        peer.Delete()
	}
	peer.Save()
}

func updateFromRemotePeerData(peer *types.Peer) {
	publicKeys, err := queryPeer(peer.URL)
	if err != nil {
		peer.Failures++
		peer.LastFailure = time.Now()
		log.Errorf("Unable to query the peer: %s, Error: %s", peer.URL, err)
	} else {
		peer.Failures = 0
		peer.Tries = 0
		peer.PublicKeys = publicKeys
		peer.SkipCycles = 10 * len(publicKeys)
	}
	for j := range publicKeys {
		pkURL := types.NewPublicKeyUrl(publicKeys[j], peer.URL)
		pkURL.Save()
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
	err := types.UpdateNewPeers(urls)
	if err != nil {
		//panic
	}
}

func PeerAdd(url string) {
	peerAddAll(url)
}
//GetPeers get peers
func GetPeers() []string {
	peerList, err := types.GetAllPeers()
	if err != nil {
		// panic
	}
	urls := make([]string, 0, len(*peerList))
	for _, peer := range *peerList {
		urls = append(urls, peer.URL)
	}
	return urls
}

//GetPeerURL get url
func GetPeerURL(publicKey []byte) (string, error) {
	peer, err := types.FindPublicKeyUrl(publicKey)
	if err != nil {
		//panic
	}
	if peer != nil {
		return peer.URL, nil
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
