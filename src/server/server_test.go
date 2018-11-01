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
)

func TestMain(m *testing.M) {
	StartServer("9000")
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

func TestVersion(t *testing.T) {
	response := doRequest("http://localhost:9000/version", t)

	if (!reflect.DeepEqual(response, api.BlackBoxVersion)) {
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