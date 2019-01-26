package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tv42/httpunix"

	"Smilo-blackbox/src/server/config"
	"path/filepath"
	"path"
	"strings"
	"fmt"
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


	currentDir, _ := os.Getwd()
	var workDir string

	isServer := strings.HasSuffix(currentDir, "/server")
	isData := strings.HasSuffix(currentDir, "/data")
	isRoot := strings.HasSuffix(currentDir, "/Smilo-blackbox")
	if isServer {
		workDir = "../../"
		fmt.Println("Contains /server")
	} else if isData  {
		workDir = "../../"
		fmt.Println("Contains /data")
	} else if isRoot {
		fmt.Println("is root dir")
		workDir = ""
	}

	socketFile := filepath.Join(config.WorkDir.Value, config.Socket.Value)

	finalPath := path.Join(workDir, socketFile)

	finalPath = path.Join(currentDir, finalPath)

	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		log.Error("ERROR: Could not open IPC file, ", " socketFile: ", socketFile, ", ERROR: ",err)
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

func doDeleteRequest(t *testing.T, url string) (string, int) {
	client := new(http.Client)
	req, _ := http.NewRequest("DELETE", url, http.NoBody)
	response, err := client.Do(req)
	ret := getResponseData(t, err, response)
	require.NotEmpty(t, ret)
	return ret, response.StatusCode
}

func doPostRequest(t *testing.T, _url string, params url.Values) string {
	client := new(http.Client)
	response, err := client.PostForm(_url, params)
	ret := getResponseData(t, err, response)
	return ret
}

func doPostJsonRequest(t *testing.T, _url string, json string) string {
	client := new(http.Client)
	response, err := client.Post(_url, "application/json", bytes.NewBuffer([]byte(json)))
	ret := getResponseData(t, err, response)
	return ret
}

func doRequest(t *testing.T, url string) string {
	client := new(http.Client)
	response, err := client.Get(url)
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

func removeIfExists(file string) {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		os.Remove(file)
	}
}
