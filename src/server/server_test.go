package server

import (
	"testing"
	"net/http"
	"reflect"
	"io/ioutil"
	"os"
	"Smilo-blackbox/src/server/api"
	"github.com/tv42/httpunix"
	"time"
	"net/url"
)

func TestMain(m *testing.M) {
	go StartServer("9000", "")
	time.Sleep(100000000)
	retcode := m.Run();
	os.Exit(retcode)
}

func TestUnixUpcheck(t *testing.T) {
	response := doUnixRequest("/upcheck", t)

	if (!reflect.DeepEqual(response, "I'm up!")) {
		t.Fail()
	}
}

func TestHttpUpcheck(t *testing.T) {
	response := doRequest("http://localhost:9000/upcheck", t)

	if (!reflect.DeepEqual(response, "I'm up!")) {
		t.Fail()
	}
}

func TestUnixVersion(t *testing.T) {
	response := doUnixRequest("/version", t)

	if (!reflect.DeepEqual(response, api.BlackBoxVersion)) {
		t.Fail()
	}
}


func TestHttpVersion(t *testing.T) {
	response := doRequest("http://localhost:9000/version", t)

	if (!reflect.DeepEqual(response, api.BlackBoxVersion)) {
		t.Fail()
	}
}

func TestHttpTransactionGet(t *testing.T) {
	response := doRequest("http://localhost:9000/transaction/1", t)

	if (!reflect.DeepEqual(response, "1")) {
		t.Fail()
	}
}

func TestHttpTransactionDelete(t *testing.T) {
	response := doDeleteRequest("http://localhost:9000/transaction/1", t)

	if (!reflect.DeepEqual(response, "1")) {
		t.Fail()
	}
}

func TestHttpDelete(t *testing.T) {
	params := url.Values{}
    params.Add("Encoded public key", "123456")
	response := doPostRequest("http://localhost:9000/delete", t, params)

	if (!reflect.DeepEqual(response, "123456")) {
		t.Fail()
	}
}

func doUnixRequest(endpoint string, t *testing.T) (string) {
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	u.RegisterLocation("myservice", sockPath)

	var client = http.Client{
		Transport: u,
	}

	response, err := client.Get("http+unix://myservice"+endpoint)
	ret := getResponseData(err, t, response)
	return ret
}

func doDeleteRequest(url string, t *testing.T) (string) {
	client := new(http.Client)
	req, _ := http.NewRequest("DELETE", url, http.NoBody)
	response, err := client.Do(req)
	ret := getResponseData(err, t, response)
	return ret
}

func doPostRequest(_url string, t *testing.T, params url.Values) (string) {
	client := new(http.Client)
	response, err := client.PostForm(_url,params)
	ret := getResponseData(err, t, response)
	return ret
}

func doRequest(url string, t *testing.T) (string) {
	client := new(http.Client)
	response, err := client.Get(url)
	ret := getResponseData(err, t, response)
	return ret
}

func getResponseData(err error, t *testing.T, response *http.Response) string {
	ret := ""
	if (err != nil) {
		t.Fail()
	} else {
		p, error := ioutil.ReadAll(response.Body)
		if (error != nil) {
			t.Fail()
		} else {
			ret = string(p)
		}
	}
	return ret
}