package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"Smilo-blackbox/src/server/api"

	"github.com/tv42/httpunix"
	"encoding/base64"
	"encoding/json"
	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/config"
)

func TestMain(m *testing.M) {
	removeIfExists("./blackbox.db")
	removeIfExists("./blackbox.sock")
	config.LoadConfig("./server_test.conf")
	data.Start("./blackbox.db")
	go StartServer(9000, "")
	time.Sleep(500000000)
	retcode := m.Run()
	os.Exit(retcode)
}

func removeIfExists(file string) {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		os.Remove(file)
	}
}

func TestUnixUpcheck(t *testing.T) {
	response := doUnixRequest("/upcheck", t)

	if !reflect.DeepEqual(response, "I'm up!") {
		t.Fail()
	}
}

func TestHttpUpcheck(t *testing.T) {
	response := doRequest("http://localhost:9000/upcheck", t)

	if !reflect.DeepEqual(response, "I'm up!") {
		t.Fail()
	}
}

func TestUnixVersion(t *testing.T) {
	response := doUnixRequest("/version", t)

	if !reflect.DeepEqual(response, api.BlackBoxVersion) {
		t.Fail()
	}
}

func TestHttpVersion(t *testing.T) {
	response := doRequest("http://localhost:9000/version", t)

	if !reflect.DeepEqual(response, api.BlackBoxVersion) {
		t.Fail()
	}
}

func TestHttpTransactionDelete(t *testing.T) {
	_, status := doDeleteRequest("http://localhost:9000/transaction/1", t)

	if status != 204 {
		t.Fail()
	}
}

func TestHttpDelete(t *testing.T) {
	//params := url.Values{}
	//params.Add("Encoded public key", "123456")
	//tmp := &api.DeleteRequest{ Key: "123456" }
	//params, _ := json.Marshal(tmp)
	response := doPostJsonRequest("http://localhost:9000/delete", t, "{\"key\": \"123456\" }")

	if !reflect.DeepEqual(response, "123456") {
		t.Fail()
	}
}

func TestUnixSend(t *testing.T) {
	to := make([]string,1)
	to[0] = "OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="
	sendRequest := api.SendRequest{ Payload: base64.StdEncoding.EncodeToString([]byte("1234567890abcdefghijklmnopqrs")), From: "MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=", To: to}
	req, err := json.Marshal(sendRequest)
	if err != nil {
		t.Fail()
	}
	response := doUnixPostJsonRequest("/send", t, string(req))
	var sendResponse api.SendResponse
	json.Unmarshal([]byte(response),&sendResponse)

    receiveRequest := api.ReceiveRequest{ Key: sendResponse.Key, To: sendRequest.To[0]}
	req2, err2 := json.Marshal(receiveRequest)
	if err2 != nil {
		t.Fail()
	}
	log.Debug("Send Response: " + sendResponse.Key)

	response = doUnixGetJsonRequest("/receive", t, string(req2))
	var receiveResponse api.ReceiveResponse
	json.Unmarshal([]byte(response),&receiveResponse)

	log.Debug("Receive Response: " + receiveResponse.Payload)
	if sendRequest.Payload != receiveResponse.Payload {
		t.Fail()
	}
}

func TestUnixSendRawTransactionGet(t *testing.T) {
	to := make([]string,1)
	to[0] = "OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="
	payload := "1234567890abcdefghijklmnopqrs"
	encPayload := base64.StdEncoding.EncodeToString([]byte(payload))
	from := "MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk="
	response := doUnixPostRequest("/sendraw", t, []byte(encPayload), http.Header{"c11n-from" : []string{from}, "c11n-to" : to})

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
	response = doUnixRequest("/transaction/" + urlEncodedKey + "?to="+urlEncodedTo, t)
	var receiveResponse api.ReceiveResponse
	json.Unmarshal([]byte(response),&receiveResponse)
    retorno, _ := base64.StdEncoding.DecodeString(receiveResponse.Payload)
	log.Debug("Receive Response: " + receiveResponse.Payload)
	if string(payload) != string(retorno) {
		t.Fail()
	}
}

func doUnixPostJsonRequest(endpoint string, t *testing.T, json string) string {
	client := getSocketClient()

	response, err := client.Post("http+unix://myservice" + endpoint, "application/json", bytes.NewBuffer([]byte(json)))
	ret := getResponseData(err, t, response)
	return ret
}

func doUnixGetJsonRequest(endpoint string, t *testing.T, json string) string {
	client := getSocketClient()
	req, _ := http.NewRequest("GET","http+unix://myservice" + endpoint, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	ret := getResponseData(err, t, response)
	return ret
}

func getSocketClient() *http.Client {
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	u.RegisterLocation("myservice", sockPath)
	var client = http.Client{
		Transport: u,
	}
	return &client
}

func doUnixRequest(endpoint string, t *testing.T) string {
	client := getSocketClient()

	response, err := client.Get("http+unix://myservice" + endpoint)
	ret := getResponseData(err, t, response)
	return ret
}

func doUnixPostRequest(endpoint string, t *testing.T, payload []byte, headers http.Header) string {
	client := getSocketClient()

	req, _ := http.NewRequest("POST","http+unix://myservice" + endpoint, bytes.NewBuffer(payload))
	req.Header = headers
	req.Header.Set("Content-Type", "application/octet-stream")
	response, err := client.Do(req)

	ret := getResponseData(err, t, response)
	return ret
}

func doDeleteRequest(url string, t *testing.T) (string, int) {
	client := new(http.Client)
	req, _ := http.NewRequest("DELETE", url, http.NoBody)
	response, err := client.Do(req)
	ret := getResponseData(err, t, response)
	return ret, response.StatusCode
}

func doPostRequest(_url string, t *testing.T, params url.Values) string {
	client := new(http.Client)
	response, err := client.PostForm(_url, params)
	ret := getResponseData(err, t, response)
	return ret
}

func doPostJsonRequest(_url string, t *testing.T, json string) string {
	client := new(http.Client)
	response, err := client.Post(_url, "application/json", bytes.NewBuffer([]byte(json)))
	ret := getResponseData(err, t, response)
	return ret
}

func doRequest(url string, t *testing.T) string {
	client := new(http.Client)
	response, err := client.Get(url)
	ret := getResponseData(err, t, response)
	return ret
}

func getResponseData(err error, t *testing.T, response *http.Response) string {
	ret := ""
	defer response.Body.Close()
	if err != nil {
		t.Fail()
	} else {
		p, error := ioutil.ReadAll(response.Body)
		if error != nil {
			t.Fail()
		} else {
			ret = string(p)
		}
	}
	return ret
}
