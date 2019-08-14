// Copyright 2019 The Smilo-blackbox Authors
// This file is part of the Smilo-blackbox library.
//
// The Smilo-blackbox library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Smilo-blackbox library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Smilo-blackbox library. If not, see <http://www.gnu.org/licenses/>.

package test

import (
	"Smilo-blackbox/src/server"
	"Smilo-blackbox/src/server/api"
	"Smilo-blackbox/src/server/syncpeer"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)
type TestServer struct {
	Port int
	Client http.Client
	PublicKey string
}
var (
	testServers = make([]TestServer,5)
	TEST_PAYLOAD = base64.StdEncoding.EncodeToString([]byte("1234567890abcdefghijklmnopqrs"))
)
func init() {
	testServers[0].Port = 9001
	testServers[0].Client = server.GetSocketClient("./blackbox1.ipc")
	testServers[0].PublicKey = "/TOE4TKtAqVsePRVR+5AA43HkAK5DSntkOCO7nYq5xU="
	testServers[1].Port = 9002
	testServers[1].Client = server.GetSocketClient("./blackbox2.ipc")
	testServers[1].PublicKey = "rYxIwmdlrqetxTYolgXBq+qVBQCT29IYyWq9JIGgNWU="
	testServers[2].Port = 9003
	testServers[2].Client = server.GetSocketClient("./blackbox3.ipc")
	testServers[2].PublicKey = "mVL7flODxSLJVN6U8uRiDT4IzZ5ySK0jIH+e9VyQQUQ="
	testServers[3].Port = 9004
	testServers[3].Client = server.GetSocketClient("./blackbox4.ipc")
	testServers[3].PublicKey = "7BPSOhfa8XR1DfetZ6hqHU6r7I9RdjZgUoHB2xjGZkk="
	testServers[4].Port = 9005
	testServers[4].Client = server.GetSocketClient("./blackbox5.ipc")
	testServers[4].PublicKey = "PSe+1pnRmrR910zyTVL6ngJOFXLPu8CBW+hjFI0+dhw="
	certData, _ := ioutil.ReadFile("./rootCA.crt")
	syncpeer.AppendCertificate(certData)
}

func TestMain(m *testing.M) {
	waitNodesUp([]int{int(9001),int(9002),int(9003),int(9004),int(9005)})
	time.Sleep(30 * time.Second)
	m.Run()
}

func TestIntegrationSendAll(t *testing.T) {
	to := make([]string, 4)
	to[0] = testServers[1].PublicKey
	to[1] = testServers[2].PublicKey
	to[2] = testServers[3].PublicKey
	to[3] = testServers[4].PublicKey
	sendResponse := sendTestPayload(t, testServers[0], to)

	for i:=1; i<5; i++ {
		receiveResponse := receiveTestPayload(t, testServers[i], sendResponse.Key)
		require.Equal(t, TEST_PAYLOAD, receiveResponse.Payload,"Payload not received on Server "+fmt.Sprint(i))
	}
}

func TestIntegrationSendFew(t *testing.T) {
	to := make([]string, 2)
	to[0] = testServers[1].PublicKey
	to[1] = testServers[2].PublicKey
	sendResponse := sendTestPayload(t, testServers[0], to)

	for i:=1; i<5; i++ {
		receiveResponse := receiveTestPayload(t, testServers[i], sendResponse.Key)
		if i<3 {
			require.Equal(t, TEST_PAYLOAD, receiveResponse.Payload,"Payload not received on Server "+fmt.Sprint(i))
		} else {
			require.Equal(t, "", receiveResponse.Payload,"Payload received by mistake on Server "+fmt.Sprint(i))
		}
	}

}

func TestPeerURLPropagation(t *testing.T) {
	response := getPartyInfo(t, testServers[3])
	require.NotEmpty(t, response.PeerURLs)
	require.Equal(t, 5, len(response.PeerURLs))
}

func getPartyInfo(t *testing.T, targetServer TestServer) syncpeer.PartyInfoResponse {
	partyInfoRequest := syncpeer.PartyInfoRequest{SenderURL:"", SenderNonce: base64.StdEncoding.EncodeToString(make([]byte,24)), SenderKey:targetServer.PublicKey}
	req, err := json.Marshal(partyInfoRequest)
	require.Empty(t,err)
	response := DoPostJSONRequest(t, "http://localhost:"+fmt.Sprint(targetServer.Port)+"/partyinfo", string(req))
	var partyInfoResponse syncpeer.PartyInfoResponse
	json.Unmarshal([]byte(response), &partyInfoResponse)
	return partyInfoResponse
}

func receiveTestPayload(t *testing.T, targetServer TestServer, key string) api.ReceiveResponse {
	receiveRequest := api.ReceiveRequest{Key: key, To: targetServer.PublicKey}
	req2, err2 := json.Marshal(receiveRequest)
	require.Empty(t, err2)
	response := doReceiveRequest(t, targetServer, string(req2))
	var receiveResponse api.ReceiveResponse
	json.Unmarshal([]byte(response), &receiveResponse)
	return receiveResponse
}

func sendTestPayload(t *testing.T, targetServer TestServer, to []string) (api.KeyJson) {
	sendRequest := api.SendRequest{Payload: TEST_PAYLOAD, From: targetServer.PublicKey, To: to}
	req, err := json.Marshal(sendRequest)
	if err != nil {
		t.Fail()
	}
	response := doSendRequest(t, targetServer, string(req))
	var sendResponse api.KeyJson
	json.Unmarshal([]byte(response), &sendResponse)
	return sendResponse
}

func doReceiveRequest(t *testing.T, targetServer TestServer, json string) (string) {
	req, _ := http.NewRequest("GET", "http+unix://myservice/receive", bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")
	response, _ := targetServer.Client.Do(req)
	ret, _ := getResponseData(response)
	return ret
}

func doSendRequest(t *testing.T, targetServer TestServer, json string) (string) {
	response, _ := targetServer.Client.Post("http+unix://myservice/send", "application/json", bytes.NewBuffer([]byte(json)))
	ret, _ := getResponseData(response)
	return ret
}

func waitNodesUp(ports []int) {
	for {
		allUp := true
		for _, port := range ports {
			if !getUpcheck(port) {
				allUp = false
				log.Info("Node %i still down.", port)
				break
			}
		}
		if allUp {
			break
		}
	}
}

func getUpcheck(port int) bool {
     ret := DoRequest("http://localhost:"+fmt.Sprint(port)+"/upcheck")
     if ret == "" {
     	ret = DoRequest("https://localhost:"+fmt.Sprint(port)+"/upcheck")
        if ret == "" {
			return false
		}
	 }
	 return true
}

func DoPostJSONRequest(t *testing.T, _url string, json string) string {
	client := syncpeer.GetHTTPClient()
	response, err := client.Post(_url, "application/json", bytes.NewBuffer([]byte(json)))
	if err != nil {
		return ""
	}
	ret, err := getResponseData(response)
	if err != nil {
		return ""
	}
	return ret
}

func DoRequest(url string) string {
	client := syncpeer.GetHTTPClient()
	response, err := client.Get(url)
	if err != nil {
		return ""
	}
	ret, err := getResponseData(response)
	if err != nil {
		return ""
	}
	return ret
}

func getResponseData(response *http.Response) (string, error) {
	ret := ""
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
	p, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return "", error
	}
	ret = string(p)

	return ret, nil
}
