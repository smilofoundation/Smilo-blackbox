package server

import (
	"net/http"
	"testing"
	"Smilo-blackbox/src/server/api"

	"encoding/base64"
	"encoding/json"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/server/config"
	"os"
	"time"
)

func TestMain(m *testing.M) {
	removeIfExists("./blackbox.db")
	removeIfExists("./blackbox.ipc")
	config.LoadConfig("./server_test.conf")

	go StartServer()

	config.WorkDir.Value = ""

	time.Sleep(2000000000)
	retcode := m.Run()
	os.Exit(retcode)
}

func TestUnixSend(t *testing.T) {
	to := make([]string, 1)
	to[0] = "OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="
	sendRequest := api.SendRequest{Payload: base64.StdEncoding.EncodeToString([]byte("1234567890abcdefghijklmnopqrs")), From: "MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=", To: to}
	req, err := json.Marshal(sendRequest)
	if err != nil {
		t.Fail()
	}
	response := doUnixPostJsonRequest(t, "/send", string(req))
	var sendResponse api.SendResponse
	json.Unmarshal([]byte(response), &sendResponse)

	receiveRequest := api.ReceiveRequest{Key: sendResponse.Key, To: sendRequest.To[0]}
	req2, err2 := json.Marshal(receiveRequest)
	require.Empty(t, err2)

	log.Debug("Send Response: " + sendResponse.Key)

	response = doUnixGetJsonRequest(t, "/receive", string(req2))
	var receiveResponse api.ReceiveResponse
	json.Unmarshal([]byte(response), &receiveResponse)

	log.Debug("Receive Response: " + receiveResponse.Payload)
	require.Equal(t, sendRequest.Payload, receiveResponse.Payload)
}

func TestUnixSendRawTransactionGet(t *testing.T) {
	to := make([]string, 1)
	to[0] = "OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="
	payload := "1234567890abcdefghijklmnopqrs"
	encPayload := base64.StdEncoding.EncodeToString([]byte(payload))
	from := "MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk="
	response := doUnixPostRequest(t, "/sendraw", []byte(encPayload), http.Header{"c11n-from": []string{from}, "c11n-to": to})

	key, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		t.Fail()
	}
	urlEncodedKey := base64.URLEncoding.EncodeToString(key)
	log.Debug("Send Response: " + response)
	toBytes, err := base64.StdEncoding.DecodeString(to[0])
	if err != nil {
		t.Fail()
	}
	urlEncodedTo := base64.URLEncoding.EncodeToString(toBytes)
	response = doUnixRequest(t, "/transaction/"+urlEncodedKey+"?to="+urlEncodedTo)
	var receiveResponse api.ReceiveResponse
	json.Unmarshal([]byte(response), &receiveResponse)
	retorno, _ := base64.StdEncoding.DecodeString(receiveResponse.Payload)
	log.Debug("Receive Response: " + receiveResponse.Payload)
	if string(payload) != string(retorno) {
		t.Fail()
	}
}
