package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"Smilo-blackbox/src/data"

	"encoding/base64"
	"strings"

	"Smilo-blackbox/src/server/encoding"

	"github.com/gorilla/mux"
	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/server/syncpeer"
)

//TODO
// It receives a POST request with a json containing url and key, returns local publicKeys and a proof that private key is known.
func GetPartyInfo(w http.ResponseWriter, r *http.Request) {
	var jsonReq syncpeer.PartyInfoRequest
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error (%s) decoding json.\n", r.URL, err))
		return
	}
    key, err := base64.StdEncoding.DecodeString(jsonReq.SenderKey)
    publicKeys := crypt.GetPublicKeys()
    responseJson := syncpeer.PartyInfoResponse{ PublicKeys: make([]syncpeer.ProvenPublicKey,0,len(publicKeys)), PeerURLs: syncpeer.GetPeers()}
    for _,pubkey := range publicKeys {
    	sharedKey := crypt.ComputeSharedKey(crypt.GetPrivateKey(pubkey), key)
    	randomPayload, _ := crypt.NewRandomKey()
        responseJson.PublicKeys = append(responseJson.PublicKeys, syncpeer.ProvenPublicKey{ Key: base64.StdEncoding.EncodeToString(pubkey), Proof: base64.StdEncoding.EncodeToString(crypt.EncryptPayload(sharedKey, randomPayload, nil))})
	}
	json.NewEncoder(w).Encode(responseJson)
	w.Header().Set("Content-Type", "application/json")
    syncpeer.PeerAdd(jsonReq.SenderURL)
}

// It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.
func Push(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			requestError(w, http.StatusInternalServerError, fmt.Sprintf("Cannot deserialize payload."))
		}
	}()
	encPayload, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if encPayload == nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, missing payload.\n", r.URL))
		return
	}

	payload, err := base64.StdEncoding.DecodeString(string(encPayload))
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error decoding payload: (%s), %s\n", r.URL, encPayload, err))
		return
	}

	encoding.Deserialize(payload)
	encTrans := data.NewEncryptedTransaction(payload)

	if encTrans == nil {
		requestError(w, http.StatusInternalServerError, fmt.Sprintf("Cannot save transaction."))
		return
	}

	encTrans.Save()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(base64.StdEncoding.EncodeToString(encTrans.Hash)))
}

// Receive a GET request with header params c11n-key and c11n-to, return unencrypted payload
func ReceiveRaw(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("c11n-key")
	to := r.Header.Get("c11n-to")

	if key == "" || to == "" {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, invalid headers.\n", r.URL))
		return
	}
	hash, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, c11n-key header (%s) is not a valid key.\n", r.URL, key))
		return
	}
	public, err := base64.StdEncoding.DecodeString(to)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, c11n-to header (%s) is not a valid key.\n", r.URL, to))
		return
	}

	payload := RetrieveAndDecryptPayload(w, r, hash, public)
	if payload != nil {
		w.Write([]byte(base64.StdEncoding.EncodeToString(payload)))
	}

}

// It receives a POST request with a json ResendRequest containing type (INDIVIDUAL, ALL), publicKey and key(for individual requests),
// it returns encoded payload for INDIVIDUAL or it does one push request for each payload and returns empty for type ALL.
func Resend(w http.ResponseWriter, r *http.Request) {
	var jsonReq ResendRequest
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error (%s) decoding json.\n", r.URL, err))
		return
	}
	if strings.ToUpper(jsonReq.Type) == "INDIVIDUAL" {
		key, err := base64.StdEncoding.DecodeString(jsonReq.Key)
		if err != nil {
			requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.\n", r.URL, jsonReq.Key))
			return
		}
		encTrans, err := data.FindEncryptedTransaction(key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(base64.StdEncoding.EncodeToString(encTrans.Encoded_Payload)))
	} else {
		if strings.ToUpper(jsonReq.Type) == "ALL" {
			//TODO Implement loop of push requests
			w.WriteHeader(http.StatusNoContent)
		} else {
			requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.\n", r.URL, jsonReq.Type))
		}
	}
}

// Deprecated API
// It receives a POST request with a json containing a DeleteRequest with key and returns Status 200 if succeed, 404 otherwise.
func Delete(w http.ResponseWriter, r *http.Request) {
	var jsonReq DeleteRequest
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error (%s) decoding json.\n", r.URL, err))
		return
	}
	key, err := base64.StdEncoding.DecodeString(jsonReq.Key)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.\n", r.URL, jsonReq.Key))
		return
	}
	encTrans, err := data.FindEncryptedTransaction(key)
	if encTrans == nil {
		requestError(w, http.StatusNotFound, fmt.Sprintf("Transaction key: %s not found\n", jsonReq.Key))
		return
	}
	encTrans.Delete()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete successful"))
}

// It receives a DELETE request with a key on path string and returns 204 if succeed, 404 otherwise.
func TransactionDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key, err := base64.URLEncoding.DecodeString(params["key"])
	if err != nil || params["key"] == "" {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.\n", r.URL, params["key"]))
		return
	}
	encTrans, err := data.FindEncryptedTransaction(key)
	if encTrans == nil {
		requestError(w, http.StatusNotFound, fmt.Sprintf("Transaction key: %s not found\n", params["key"]))
		return
	}
	encTrans.Delete()
	w.WriteHeader(http.StatusNoContent)
}

//TODO
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

//TODO
// Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 500 otherwise
func ConfigPeersGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["index"])
	jsonResponse := data.Peer{}
	out, _ := json.Marshal(jsonResponse)
	w.WriteHeader(200)
	w.Write(out)
}

//TODO
// Receive a GET request and return Status Code 200 and server internal status information in plain text.
func Metrics(w http.ResponseWriter, r *http.Request) {

}
