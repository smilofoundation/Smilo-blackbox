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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"Smilo-blackbox/src/data"

	"encoding/base64"
	"strings"

	"Smilo-blackbox/src/server/encoding"

	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/server/syncpeer"
)

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
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error (%s) decoding sender public key.", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	nonce, err := base64.StdEncoding.DecodeString(jsonReq.SenderNonce)
	if err != nil {
		message := fmt.Sprintf("Invalid request: %s, error (%s) decoding sender nonce.", r.URL, err)
		log.Error(message)
		requestError(w, http.StatusBadRequest, message)
		return
	}
	publicKeys := crypt.GetPublicKeys()
	responseJson := syncpeer.PartyInfoResponse{PublicKeys: make([]syncpeer.ProvenPublicKey, 0, len(publicKeys)), PeerURLs: syncpeer.GetPeers()}
	for _, pubkey := range publicKeys {
		sharedKey := crypt.ComputeSharedKey(crypt.GetPrivateKey(pubkey), key)
		randomPayload, _ := crypt.NewRandomKey()
		responseJson.PublicKeys = append(responseJson.PublicKeys, syncpeer.ProvenPublicKey{Key: base64.StdEncoding.EncodeToString(pubkey), Proof: base64.StdEncoding.EncodeToString(crypt.EncryptPayload(sharedKey, randomPayload, nonce))})
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
