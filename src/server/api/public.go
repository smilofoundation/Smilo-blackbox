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
// GetPartyInfo It receives a POST request with a json containing url and key, returns local publicKeys and a proof that private key is known.
func GetPartyInfo(w http.ResponseWriter, r *http.Request) {
	var jsonReq syncpeer.PartyInfoRequest
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error (%s) decoding json.", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	key, err := base64.StdEncoding.DecodeString(jsonReq.SenderKey)
	publicKeys := crypt.GetPublicKeys()
	responseJson := syncpeer.PartyInfoResponse{PublicKeys: make([]syncpeer.ProvenPublicKey, 0, len(publicKeys)), PeerURLs: syncpeer.GetPeers()}
	for _, pubkey := range publicKeys {
		sharedKey := crypt.ComputeSharedKey(crypt.GetPrivateKey(pubkey), key)
		randomPayload, _ := crypt.NewRandomKey()
		responseJson.PublicKeys = append(responseJson.PublicKeys, syncpeer.ProvenPublicKey{Key: base64.StdEncoding.EncodeToString(pubkey), Proof: base64.StdEncoding.EncodeToString(crypt.EncryptPayload(sharedKey, randomPayload, nil))})
	}
	json.NewEncoder(w).Encode(responseJson)
	w.Header().Set("Content-Type", "application/json")
	syncpeer.PeerAdd(jsonReq.SenderURL)
}

// Push It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.
func Push(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			message := fmt.Sprintf("Cannot deserialize payload.")
			log.Error(message)
			requestError(w, http.StatusInternalServerError, message)
		}
	}()
	encPayload, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if encPayload == nil {
		message := fmt.Sprintf("Invalid request: %s, missing payload.", r.URL)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	payload, err := base64.StdEncoding.DecodeString(string(encPayload))
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error decoding payload: (%s), %s", r.URL, encPayload, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	encoding.Deserialize(payload)
	encTrans := data.NewEncryptedTransaction(payload)

	if encTrans == nil {
		message := fmt.Sprintf("Cannot save transaction.")
		log.Error(message)
		requestError(w, http.StatusInternalServerError, message)
		return
	}

	encTrans.Save()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(base64.StdEncoding.EncodeToString(encTrans.Hash)))
}

// ReceiveRaw Receive a GET request with header params c11n-key and c11n-to, return unencrypted payload
func ReceiveRaw(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("c11n-key")
	to := r.Header.Get("c11n-to")

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
		message := fmt.Sprintf("Invalid request: %s, c11n-key header (%s) is not a valid key.", r.URL, key)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	public, err := base64.StdEncoding.DecodeString(to)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, c11n-to header (%s) is not a valid key.", r.URL, to)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}

	payload := RetrieveAndDecryptPayload(w, r, hash, public)
	if payload != nil {
		log.Info("Found transaction! ", base64.StdEncoding.EncodeToString(payload))
		w.Write([]byte(base64.StdEncoding.EncodeToString(payload)))
	} else {
		log.WithField("key", key).WithField("hash",hash).WithField("public", public).
			Error("Could not find valid data for the request.")
	}

}

// Resend It receives a POST request with a json ResendRequest containing type (INDIVIDUAL, ALL), publicKey and key(for individual requests),
// it returns encoded payload for INDIVIDUAL or it does one push request for each payload and returns empty for type ALL.
func Resend(w http.ResponseWriter, r *http.Request) {
	var jsonReq ResendRequest
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error (%s) decoding json.", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	if strings.ToUpper(jsonReq.Type) == "INDIVIDUAL" {
		key, err := base64.StdEncoding.DecodeString(jsonReq.Key)
		if err != nil {
			message := fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.", r.URL, jsonReq.Key)
			log.Error(message)
			requestError(w, http.StatusBadRequest, message)
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
			message := fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.", r.URL, jsonReq.Type)
			log.Error(message)
			requestError(w, http.StatusBadRequest, message)
		}
	}
}

// Delete Deprecated API
// It receives a POST request with a json containing a DeleteRequest with key and returns Status 200 if succeed, 404 otherwise.
func Delete(w http.ResponseWriter, r *http.Request) {
	var jsonReq DeleteRequest
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &jsonReq)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error (%s) decoding json.", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	key, err := base64.StdEncoding.DecodeString(jsonReq.Key)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, Key (%s) is not a valid BASE64 key.", r.URL, jsonReq.Key)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	encTrans, err := data.FindEncryptedTransaction(key)
	if encTrans == nil {
		message := fmt.Sprintf("Transaction key: %s not found", jsonReq.Key)
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return
	}
	encTrans.Delete()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete successful"))
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

//TODO
// ConfigPeersGet Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 404 if not found.
func ConfigPeersGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	publicKey, err := base64.URLEncoding.DecodeString(params["index"])
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, Public Key (%s) is not a valid BASE64 key.", r.URL, params["index"])
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	url, err := syncpeer.GetPeerURL(publicKey)
	if err != nil {
		message := fmt.Sprintf("Public key: %s not found", params["index"])
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return
	}
	jsonResponse := PeerUrl{Url: url}
	out, _ := json.Marshal(jsonResponse)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

//TODO
// Metrics Receive a GET request and return Status Code 200 and server internal status information in plain text.
func Metrics(w http.ResponseWriter, r *http.Request) {

}
