package test

import (
	"testing"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"Smilo-blackbox/src/server/api"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"Smilo-blackbox/src/server"
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
}

func TestIntegrationSendAll(t *testing.T) {
	waitNodesUp([]int{int(9001),int(9002),int(9003),int(9004),int(9005)})
	time.Sleep(1 * time.Minute)
	to := make([]string, 4)
	to[0] = testServers[1].PublicKey
	to[1] = testServers[2].PublicKey
	to[2] = testServers[3].PublicKey
	to[3] = testServers[4].PublicKey
	sendResponse := sendTestPayload(t, testServers[0], to)

	for i:=1; i<5; i++ {
		receiveResponse := receiveTestPayload(t, testServers[i], sendResponse.Key)
		if receiveResponse.Payload != TEST_PAYLOAD {
			require.Equal(t, TEST_PAYLOAD, receiveResponse.Payload,"Payload not received on Server "+fmt.Sprint(i))
		}
	}
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

func sendTestPayload(t *testing.T, targetServer TestServer, to []string) (api.SendResponse) {
	sendRequest := api.SendRequest{Payload: TEST_PAYLOAD, From: targetServer.PublicKey, To: to}
	req, err := json.Marshal(sendRequest)
	if err != nil {
		t.Fail()
	}
	response := doSendRequest(t, targetServer, string(req))
	var sendResponse api.SendResponse
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
     	return false
	 }
	 return true
}

func DoPostJsonRequest(t *testing.T, _url string, json string) string {
	client := new(http.Client)
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
	client := new(http.Client)
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