package syncpeer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/data/types"
	"Smilo-blackbox/src/utils"
)

func TestMain(m *testing.M) {
	data.SetEngine("boltdb")
	filename := utils.BuildFilename("blackbox_sync_test.db")
	_ = os.Remove(filename)
	data.SetFilename(filename)
	data.Start()
	pubKey, _ := base64.StdEncoding.DecodeString("rYxIwmdlrqetxTYolgXBq+qVBQCT29IYyWq9JIGgNWU=")
	privKey, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAI=")
	crypt.PutKeyPair(crypt.KeyPair{PublicKey: pubKey, PrivateKey: privKey})
	crypt.ReadRandom = func(b []byte) (int, error) {
		for i := 0; i < len(b); i++ {
			b[i] = 0
		}
		return len(b), nil
	}
	SetHostURL("http://localhost:9001")
	errorCode := m.Run()
	_ = os.Remove(filename)
	os.Exit(errorCode)
}

func TestStartSync(t *testing.T) {
	cli := GetHTTPClient()
	cli.PostResponseFunction = func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
		switch url {
		case "http://localhost:9002/partyinfo":
			return &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBuffer(readAll(t, "./9002_partyinfo.json"))),
				StatusCode: http.StatusOK,
			}, nil
		case "http://localhost:9003/partyinfo":
			return &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBuffer(readAll(t, "./9003_partyinfo.json"))),
				StatusCode: http.StatusOK,
			}, nil
		case "http://localhost:9004/partyinfo":
			return &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBuffer(readAll(t, "./9004_partyinfo.json"))),
				StatusCode: http.StatusOK,
			}, nil
		case "http://localhost:9002":
			return &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBuffer(readAll(t, "./9002.json"))),
				StatusCode: http.StatusOK,
			}, nil
		}
		return nil, fmt.Errorf("Unknown test request %s", url)
	}
	StartSync()
	peer, err := types.FindNextUpdatablePeer(0)
	if err != nil {
		t.Fail()
	}
	require.Nil(t, peer)
	PeerAdd("http://localhost:9002")
	peer, err = types.FindNextUpdatablePeer(10 * time.Second)
	if err != nil {
		t.Fail()
	}
	require.NotNil(t, peer)
	peer, err = types.FindNextUpdatablePeer(0)
	if err != nil {
		t.Fail()
	}
	require.Nil(t, peer)

	time.Sleep(20 * time.Second)
	pubkey, _ := base64.StdEncoding.DecodeString("/TOE4TKtAqVsePRVR+5AA43HkAK5DSntkOCO7nYq5xU=")
	url, err := GetPeerURL(pubkey)
	require.Nil(t, err)
	require.Equal(t, "http://localhost:9002", url)
	peers, _ := types.GetAllPeers()
	require.Equal(t, len(*peers), 3)
	time.Sleep(6 * time.Second)
}

func readAll(t *testing.T, filename string) []byte {
	ret, err := utils.ReadAllFile(filename, log)
	if err != nil {
		t.Fatalf("Error reading json test return: %s", err)
	}
	return ret
}
