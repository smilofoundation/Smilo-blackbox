package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

// It receives headers "c11n-from" and "c11n-to", payload body and returns Status Code 200 and encoded key plain text.
func SendRaw(w http.ResponseWriter, r *http.Request) {

}

// It receives json SendRequest with from, to and payload, returns Status Code 200 and json SendResponse with encoded key.
func Send(w http.ResponseWriter, r *http.Request) {

}
// Deprecated API
// It receives a ReceiveRequest json with an encoded key (hash) and to values, returns decrypted payload
func Receive(w http.ResponseWriter, r *http.Request) {

}

// it receives a GET request with a hash on path and query var "to" with encoded hash and to, returns decrypted payload
func TransactionGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	r.ParseForm()
	to := r.Form.Get("to")
	fmt.Println(to)
	w.WriteHeader(200)
	w.Write([]byte(params["hash"]))
}

