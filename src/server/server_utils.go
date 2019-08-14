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

package server

import (
	"bytes"
	"fmt"
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

// DoUnixPostJSONRequest is used for test real request calls.
func DoUnixPostJSONRequest(t *testing.T, endpoint string, json string) string {
	client := getSocketClient()

	response, err := client.Post("http+unix://myservice"+endpoint, "application/json", bytes.NewBuffer([]byte(json)))
	ret := getResponseData(t, err, response)
	return ret
}

// DoUnixGetJSONRequest is used for test real request calls.
func DoUnixGetJSONRequest(t *testing.T, endpoint string, json string) string {
	client := getSocketClient()
	req, _ := http.NewRequest("GET", "http+unix://myservice"+endpoint, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	ret := getResponseData(t, err, response)
	return ret
}

func getSocketClient() *http.Client {

	socketFile := filepath.Join(config.WorkDir.Value, config.Socket.Value)

	client := GetSocketClient(socketFile)
	return &client
}

// GetSocketClient is used for test real request calls.
func GetSocketClient(socketFile string) http.Client {
	finalPath := utils.BuildFilename(socketFile)
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		log.Error("ERROR: Could not open IPC file, ", " socketFile: ", socketFile, ", ERROR: ", err)
		os.Exit(1)
	}
	u.RegisterLocation("myservice", finalPath)
	var client = http.Client{
		Transport: u,
	}
	return client
}

// DoUnixRequest is used for test real request calls.
func DoUnixRequest(t *testing.T, endpoint string) string {
	client := getSocketClient()

	response, err := client.Get("http+unix://myservice" + endpoint)
	ret := getResponseData(t, err, response)
	return ret
}

// DoUnixPostRequest is used for test real request calls.
func DoUnixPostRequest(t *testing.T, endpoint string, payload []byte, headers http.Header) string {
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
			err = response.Body.Close()
			if err != nil {
				fmt.Println("Could not response.Body.Close()")
			}
		}
	}()
	if err != nil {
		t.Fail()
	} else {
		p, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		} else {
			ret = string(p)
		}
	}
	return ret
}
