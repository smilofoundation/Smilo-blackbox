package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/tv42/httpunix"

	"path/filepath"
	"Smilo-blackbox/src/server/config"
	"Smilo-blackbox/src/utils"
)

func doUnixPostJsonRequest(t *testing.T, endpoint string, json string) string {
	client := getSocketClient()

	response, err := client.Post("http+unix://myservice"+endpoint, "application/json", bytes.NewBuffer([]byte(json)))
	ret := getResponseData(t, err, response)
	return ret
}

func doUnixGetJsonRequest(t *testing.T, endpoint string, json string) string {
	client := getSocketClient()
	req, _ := http.NewRequest("GET", "http+unix://myservice"+endpoint, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	ret := getResponseData(t, err, response)
	return ret
}

func getSocketClient() *http.Client {
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}

	socketFile := filepath.Join(config.WorkDir.Value, config.Socket.Value)

	finalPath := utils.BuildFilename(socketFile)

	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		log.Error("ERROR: Could not open IPC file, ", " socketFile: ", socketFile, ", ERROR: ", err)
		os.Exit(1)
	}

	u.RegisterLocation("myservice", finalPath)
	var client = http.Client{
		Transport: u,
	}
	return &client
}

func doUnixRequest(t *testing.T, endpoint string) string {
	client := getSocketClient()

	response, err := client.Get("http+unix://myservice" + endpoint)
	ret := getResponseData(t, err, response)
	return ret
}

func doUnixPostRequest(t *testing.T, endpoint string, payload []byte, headers http.Header) string {
	client := getSocketClient()

	req, _ := http.NewRequest("POST", "http+unix://myservice"+endpoint, bytes.NewBuffer(payload))
	req.Header = headers
	req.Header.Set("Content-Type", "application/octet-stream")
	response, err := client.Do(req)

	ret := getResponseData(t, err, response)
	return ret
}

func getResponseData(t *testing.T, err error, response *http.Response) string {
	ret := ""
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
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

