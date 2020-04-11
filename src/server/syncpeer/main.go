package syncpeer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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
	keepRunning         = true
	timeBetweenCycles   = 13 * time.Second
	timeBetweenRequests = 2 * time.Second
	timeBetweenUpdates  = 15 * time.Minute
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
	client = SetupHTTPClientWrapper()
)

func getRequestTimeout() (t time.Duration) {
	value := os.Getenv("REQUEST_TIMEOUT")
	_, err := strconv.Atoi(value)
	if err == nil {
		t, err = time.ParseDuration(value + "s")
		if err == nil {
			return t
		}
	}
	t = 30 * time.Second
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
		client = SetupHTTPClientWrapper()
	}
	return ok
}

func SetupHTTPClientWrapper() *HTTPClientWrapper {
	clientWrapper := HTTPClientWrapper{
		http.Client{
			Transport: tr,
			Timeout:   getRequestTimeout(),
		},
		nil,
		nil,
	}
	clientWrapper.RequestResponseFunction = func(req *http.Request) (response *http.Response, e error) {
		return clientWrapper.Client.Do(req)
	}
	clientWrapper.PostResponseFunction = func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
		return clientWrapper.Client.Post(url, contentType, body)
	}
	return &clientWrapper
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
		time.Sleep(timeBetweenRequests)
		updatePeer()
	}
}

func updatePeer() {
	peer, err := types.FindNextUpdatablePeer(timeBetweenCycles)
	if err != nil {
		log.Panicf("Unable to get next the peer, Error: %s", err)
	}
	if peer == nil {
		time.Sleep(timeBetweenCycles)
		return
	}
	log.Debugf("Starting update of: %s", peer.URL)
	updateFromRemotePeerData(peer)
}

func updateFromRemotePeerData(peer *types.Peer) {
	publicKeys, err := queryPeer(peer.URL)
	if err != nil {
		peer.Failures++
		peer.LastFailure = time.Now()
		peer.NextUpdate = time.Now().Add(timeBetweenCycles)
		log.Errorf("Unable to query the peer: %s, Error: %s", peer.URL, err)
		if peer.Failures > 9 {
			peer.Failures = 0
			peer.Tries++
			peer.NextUpdate = time.Now().Add(timeBetweenUpdates)
		}
		if peer.Tries > 3 {
			for _, pubKeys := range peer.PublicKeys {
				pkurl, err := types.FindPublicKeyURL(pubKeys)
				if err != nil {
					log.Panicf("Unable to get public key, Error: %s", err)
				}
				if pkurl.URL == peer.URL {
					err := pkurl.Delete()
					if err != nil {
						log.WithError(err).Panic("Could not delete remote peer public key.")
					}
				}
			}
			err = peer.Delete()
			if err != nil {
				log.WithError(err).Panic("Could not delete peer.")
			}
			return
		}
	} else {
		peer.Failures = 0
		peer.Tries = 0
		peer.PublicKeys = publicKeys
		peer.NextUpdate = time.Now().Add(timeBetweenUpdates)
	}
	SavePublicKeys(publicKeys, peer)
	err = peer.Save()
	if err != nil {
		log.WithError(err).Panic("Could not save peer.")
	}
}

func SavePublicKeys(publicKeys [][]byte, peer *types.Peer) {
	for j := range publicKeys {
		pkURL := types.NewPublicKeyURL(publicKeys[j], peer.URL)
		err := pkURL.Save()
		if err != nil {
			log.WithError(err).Panic("Could not save peer public key.")
		}
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
	err := types.UpdateNewPeers(urls, hostURL)
	if err != nil {
		log.WithError(err).Panic("Unable to insert peers to database.")
	}
}

func PeerAdd(url string) {
	peerAddAll(url)
}

//GetPeers get peers
func GetPeers() []string {
	peerList, err := types.GetAllPeers()
	if err != nil {
		log.WithError(err).Panic("Unable to retrieve peer list from database.")
	}
	urls := make([]string, 0, len(*peerList))
	for _, peer := range *peerList {
		urls = append(urls, peer.URL)
	}
	return urls
}

//GetPeerURL get url
func GetPeerURL(publicKey []byte) (string, error) {

	peer, err := types.FindPublicKeyURL(publicKey)
	if err != nil && err != types.ErrNotFound {
		log.WithError(err).Panic("Unable to query database")
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
		log.Info(base64.StdEncoding.EncodeToString(crypt.EncryptPayload(sharedKey, nonce, nonce)))
		if string(crypt.DecryptPayload(sharedKey, remoteProof, nonce)) != "" {
			retPubKeys = append(retPubKeys, remotePublicKey)
			log.Debugf("Public Key accepted: %s", provenKey.Key)
		} else {
			log.Debugf("Public Key ignored: %s", provenKey.Key)
		}
	}
	return retPubKeys, responseJSON.PeerURLs, nil
}

//GetHTTPClient get http client
func GetHTTPClient() *HTTPClientWrapper {
	return client
}
