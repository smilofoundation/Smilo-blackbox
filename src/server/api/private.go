package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/log"
	"Smilo-blackbox/src/server/encoding"
	"Smilo-blackbox/src/data"
	"strings"
)

// It receives headers "c11n-from" and "c11n-to", payload body and returns Status Code 200 and encoded key plain text.
func SendRaw(w http.ResponseWriter, r *http.Request) {

}

// It receives json SendRequest with from, to and payload, returns Status Code 200 and json SendResponse with encoded key.
func Send(w http.ResponseWriter, r *http.Request) {
	var sendReq SendRequest
	err := json.NewDecoder(r.Body).Decode(&sendReq)
	r.Body.Close()
	if err != nil {
		requestError(http.StatusBadRequest, w, fmt.Sprintf("Invalid request: %s, error: %s\n", r.URL, err))
		return
	}

	payload, sender, recipients, msgs := sendReq.Parse()

	if len(msgs) > 0 {
		requestError(http.StatusBadRequest, w, fmt.Sprintf("Invalid request: %s\n %s", r.URL, strings.Join(msgs,"\n")))
		return
	}

	encPayload := encoding.EncodePayloadData(payload, sender, recipients)
	encTrans := data.NewEncryptedTransaction(string(*encPayload.Serialize()))

	encTrans.Save()

	sendResp := SendResponse{Key: base64.StdEncoding.EncodeToString([]byte(encTrans.Hash))}
	json.NewEncoder(w).Encode(sendResp)
	w.Header().Set("Content-Type", "application/json")

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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(params["hash"]))
}

func requestError(returnCode int, w http.ResponseWriter, message string) {
	log.Error(message)
	w.WriteHeader(returnCode)
	fmt.Fprintf(w, message)
}

