package server

import (
	"testing"
	"net/http"
	"reflect"
	"io/ioutil"
)

func TestRequest(t *testing.T) {
	go StartServer("9000")
	client := new(http.Client)
	response, err := client.Get("http://localhost:9000/upcheck")
	if (err != nil) {
		t.Fail()
	} else {
		p, error := ioutil.ReadAll(response.Body)
		if (error != nil) {
			t.Fail()
		} else {
			if (!reflect.DeepEqual(p,[]byte("I'm up!"))) {
				t.Fail()
			}
		}
	}
}