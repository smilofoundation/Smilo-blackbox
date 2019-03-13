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

package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/encoding"

	"io/ioutil"

	"github.com/gorilla/mux"

	"encoding/hex"

	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/server/syncpeer"
	"Smilo-blackbox/src/utils"
)

// SendRaw It receives headers "bb0x-from" and "bb0x-to", payload body and returns Status Code 200 and encoded key plain text.
func SendRaw(w http.ResponseWriter, r *http.Request) {
	var fromEncoded []byte
	var err error

	from := r.Header.Get(utils.HeaderFrom)
	to := r.Header.Get(utils.HeaderTo)

	if to == "" {
		message := fmt.Sprintf("Invalid request: %s, invalid headers. to:%s", r.URL, to)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	if from != "" {
		fromEncoded, err = base64.StdEncoding.DecodeString(from)
		if err != nil {
			message := fmt.Sprintf("Invalid request: %s, bb0x-from header (%s) is not a valid key.", r.URL, from)
			log.Error(message)
			requestError(w, http.StatusBadRequest, message)
			return
		}
	} else {
		//use default
		defaultPubKey := base64.StdEncoding.EncodeToString(crypt.GetPublicKeys()[0])
		fromEncoded, err = base64.StdEncoding.DecodeString(defaultPubKey)
		log.WithField("defaultPubKey", defaultPubKey).Info("Request from NOT filled, will use default PubKey")
		if err != nil {
			message := fmt.Sprintf("Invalid request: %s, bb0x-from header (%s) is not a valid key.", r.URL, from)
			log.Error(message)
			requestError(w, http.StatusBadRequest, message)
			return
		}
	}

	recipients, errors := splitToString(to)
	if len(errors) > 0 {
		message := fmt.Sprintf("Invalid request: %s, %s.", r.URL, strings.Join(errors, ", "))
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	encPayload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || encPayload == nil {
		message := fmt.Sprintf("Invalid request: %s, missing payload, err: %s", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(encPayload)))
	n, err := base64.StdEncoding.Decode(dbuf, encPayload)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error decoding payload: (%s), %s", r.URL, encPayload, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	payload := dbuf[:n]
	if len(payload) == 0 {
		message := fmt.Sprintf("Invalid request: %s, len of payload after decode is zero: (%s), %s", r.URL, encPayload, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
	}

	encTrans := createNewEncodedTransaction(w, r, payload, fromEncoded, recipients)

	if encTrans != nil {
		w.Write([]byte(base64.StdEncoding.EncodeToString(encTrans.Hash)))
		w.Header().Set("Content-Type", "text/plain")
	}
}

// SendRaw It receives headers "bb0x-from" and "bb0x-to", payload body and returns Status Code 200 and encoded key plain text.
func SendSignedTx(w http.ResponseWriter, r *http.Request) {
	var err error

	to := r.Header.Get(utils.HeaderTo)

	if to == "" {
		message := fmt.Sprintf("Invalid request: %s, invalid headers. to:%s", r.URL, to)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	recipients, errors := splitToString(to)
	if len(errors) > 0 {
		message := fmt.Sprintf("Invalid request: %s, %s.", r.URL, strings.Join(errors, ", "))
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	encodedHash, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || encodedHash == nil {
		message := fmt.Sprintf("Invalid request: %s, missing transaction hash, err: %s", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	key, err := base64.StdEncoding.DecodeString(string(encodedHash))
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, hash value (%s) is not a valid key.", r.URL, encodedHash)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	encTrans, err := data.FindEncryptedTransaction(key)
	if err != nil || encTrans == nil {
		message := fmt.Sprintf("Transaction key: %s not found", hex.EncodeToString(key))
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return
	}

	pushToAllRecipients(recipients, encTrans)
	w.Header().Set("Content-Type", "text/plain")

}

func splitToString(to string) ([][]byte, []string) {
	encodedRecipients := strings.Split(to, ",")
	var errors []string
	var recipients = make([][]byte, len(encodedRecipients))
	for i := 0; i < len(encodedRecipients); i++ {
		decodedValue, err := base64.StdEncoding.DecodeString(encodedRecipients[i])
		if err != nil {
			errors = append(errors, fmt.Sprintf("bb0x-to header (%s) is not a valid key", encodedRecipients[i]))
		}
		recipients[i] = decodedValue
	}
	return recipients, errors
}

func createNewEncodedTransaction(w http.ResponseWriter, r *http.Request, payload []byte, fromEncoded []byte, recipients [][]byte) *data.Encrypted_Transaction {
	encPayload, err := encoding.EncodePayloadData(payload, fromEncoded, recipients)
	if err != nil {
		message := fmt.Sprintf("Error Encoding Payload on Request: url: %s, err: %s", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusInternalServerError, message)
		return nil
	}
	encTrans := data.NewEncryptedTransaction(*encPayload.Serialize())
	encTrans.Save()
	pushToAllRecipients(recipients, encTrans)
	return encTrans
}

func pushToAllRecipients(recipients [][]byte, encTrans *data.Encrypted_Transaction) {
	for _, recipient := range recipients {
		PushTransactionForOtherNodes(*encTrans, recipient)
	}
}

// ReceiveRaw Receive a GET request with header params bb0x-key and bb0x-to, return unencrypted payload
func ReceiveRaw(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(utils.HeaderKey)
	to := r.Header.Get(utils.HeaderTo)

	if key == "" {
		message := fmt.Sprintf("Invalid request: %s, invalid headers. key: %s", r.URL, key)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	if to == "" {
		//use default
		defaultPubKey := base64.StdEncoding.EncodeToString(crypt.GetPublicKeys()[0])
		to = defaultPubKey
		log.WithField("defaultPubKey", defaultPubKey).Info("Request to NOT filled, will use default PubKey")
	}

	hash, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, bb0x-key header (%s) is not a valid key.", r.URL, key)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	public, err := base64.StdEncoding.DecodeString(to)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, bb0x-to header (%s) is not a valid key.", r.URL, to)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	payload := RetrieveAndDecryptPayload(w, r, hash, public)
	if payload != nil {
		log.Info("Found transaction! ", base64.StdEncoding.EncodeToString(payload))
		w.Write([]byte(base64.StdEncoding.EncodeToString(payload)))
	} else {
		log.WithField("key", key).WithField("hash", hash).WithField("public", public).
			Error("Could not find valid data for the request.")
	}

}

// TransactionGet it receives a GET request with a hash on path and query var "to" with encoded hash and to, returns decrypted payload
func TransactionGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	r.ParseForm()
	encodedTo := r.Form.Get("to")
	hash := params["hash"]
	if hash == "" || encodedTo == "" {
		message := fmt.Sprintf("Invalid request: %s, invalid query.", r.URL)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	var errors []string
	key, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		errors = append(errors, "Invalid hash value.")
	}
	to, err2 := base64.URLEncoding.DecodeString(encodedTo)
	if err2 != nil {
		errors = append(errors, "Invalid reciepient(to) value.")
	}

	if len(errors) > 0 {
		message := fmt.Sprintf("Invalid request: %s %s", r.URL, strings.Join(errors, "\n"))
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	RetrieveJsonPayload(w, r, key, to)

}

// ConfigPeersPut It receives a PUT request with a json containing a Peer url and returns Status Code 200.
func ConfigPeersPut(w http.ResponseWriter, r *http.Request) {
	jsonReq := PeerUrl{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error (%s) decoding json.", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	syncpeer.PeerAdd(jsonReq.Url)
	w.WriteHeader(http.StatusNoContent)
}

// ConfigPeersGet Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 404 if not found.
func ConfigPeersGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	publicKey, err := base64.URLEncoding.DecodeString(params["publickey"])
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, Public Key (%s) is not a valid BASE64 key.", r.URL, params["publickey"])
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	url, err := syncpeer.GetPeerURL(publicKey)
	if err != nil {
		message := fmt.Sprintf("Public key: %s not found", params["publickey"])
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return
	}
	jsonResponse := PeerUrl{Url: url}
	out, _ := json.Marshal(jsonResponse)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// TransactionDelete It receives a DELETE request with a key on path string and returns 204 if succeed, 404 otherwise.
func TransactionDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key, err := base64.URLEncoding.DecodeString(params["key"])
	if err != nil || params["key"] == "" {
		message := fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.", r.URL, params["key"])
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	encTrans, err := data.FindEncryptedTransaction(key)
	if encTrans == nil {
		message := fmt.Sprintf("Transaction key: %s not found", params["key"])
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return
	}
	encTrans.Delete()
	w.WriteHeader(http.StatusNoContent)
}

//TODO
// Metrics Receive a GET request and return Status Code 200 and server internal status information in plain text.
func Metrics(w http.ResponseWriter, r *http.Request) {

}
