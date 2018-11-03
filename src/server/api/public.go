package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"Smilo-blackbox/src/server/data"
)

// It receives a POST request with a binary encoded PartyInfo, updates it and returns updated PartyInfo encoded.
func GetPartyInfo(w http.ResponseWriter, r *http.Request) {

}

// It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.
func Push(w http.ResponseWriter, r *http.Request) {

}

// Receive a GET request with header params c11n-key and c11n-to, return unencrypted payload
func ReceiveRaw(w http.ResponseWriter, r *http.Request) {

}

// It receives a POST request with a json ResendRequest containing type (INDIVIDUAL, ALL), publicKey and key(for individual requests),
// it returns encoded payload for INDIVIDUAL or it does one push request for each payload and returns empty for type ALL.
func Resend(w http.ResponseWriter, r *http.Request) {

}

// Deprecated API
// It receives a POST request with a json containing a DeleteRequest with key and returns Status 200 if succeed, 404 otherwise.
func Delete(w http.ResponseWriter, r *http.Request) {
	jsonReq := DeleteRequest{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &jsonReq)
	fmt.Println(err)
	w.WriteHeader(200)
    w.Write([]byte(jsonReq.Key))
}

// It receives a DELETE request with a key on path string and returns 204 if succeed, 404 otherwise.
func TransactionDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.WriteHeader(204)
	w.Write([]byte(params["key"]))
}

// It receives a PUT request with a json containing a Peer and returns Status Code 200 and the new peer URL.
func ConfigPeersPut(w http.ResponseWriter, r *http.Request) {
    jsonReq := data.Peer{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &jsonReq)
	fmt.Println(err)
	newId := string("123456")
	w.WriteHeader(200)
	w.Write([]byte(mux.CurrentRoute(r).GetName() + "/" + newId))
}

// Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 500 otherwise
func ConfigPeersGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["index"])
	jsonResponse := data.Peer{}
	out, _ := json.Marshal(jsonResponse)
	w.WriteHeader(200)
	w.Write(out)
}

// Receive a GET request and return Status Code 200 and server internal status information in plain text.
func Metrics(w http.ResponseWriter, r *http.Request) {

}
