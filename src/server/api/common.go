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
	"Smilo-blackbox/src/data/types"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"

	"Smilo-blackbox/src/server/encoding"
	"Smilo-blackbox/src/server/syncpeer"
	"Smilo-blackbox/src/utils"
)

//GetVersion Request path "/version", response plain text version ID
func GetVersion(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(utils.BlackBoxVersion))
	if err != nil {
		log.WithError(err).Error("Could not GetVersion")
	}
}

//Upcheck Request path "/upcheck", response plain text upcheck message.
func Upcheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(utils.UpcheckMessage))
	if err != nil {
		log.WithError(err).Error("Could not Upcheck")
	}
}

//API Request path "/api", response json rest api spec.
func API(w http.ResponseWriter, r *http.Request) {

}

//UnknownRequest will debug unknown reqs
func UnknownRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug("UnknownEndPoint")
}

//RetrieveJSONPayload will retrieve payload based on request
func RetrieveJSONPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) {
	payload := RetrieveAndDecryptPayload(w, r, key, to)
	if payload != nil {
		receiveResp := ReceiveResponse{Payload: base64.StdEncoding.EncodeToString(payload)}
		err := json.NewEncoder(w).Encode(receiveResp)
		if err != nil {
			log.WithError(err).Error("Could not RetrieveJSONPayload, Encode")
		}
		w.Header().Set("Content-Type", "application/json")
	}
}

//RetrieveAndDecryptPayload will retrieve and decrypt the payload
func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte {
	encTrans, err := types.FindEncryptedTransaction(key)
	if err != nil || encTrans == nil {
		message := fmt.Sprintf("Transaction key: %s not found", hex.EncodeToString(key))
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return nil
	}

	encodedPayloadData := encoding.Deserialize(encTrans.EncodedPayload)
	payload := encodedPayloadData.Decode(to)

	if payload == nil {
		message := fmt.Sprintf("Error Encoding Payload on Request: %s", r.URL)
		log.Error(message)
		requestError(w, http.StatusInternalServerError, message)
	}
	return payload
}

//PushTransactionForOtherNodes will push encrypted transaction to other nodes
func PushTransactionForOtherNodes(encryptedTransaction types.EncryptedTransaction, recipient []byte) {
	url, err := syncpeer.GetPeerURL(recipient)
	if err == nil {
		_, err := syncpeer.GetHTTPClient().Post(url+"/push", "application/octet-stream", bytes.NewBuffer([]byte(base64.StdEncoding.EncodeToString(encryptedTransaction.EncodedPayload))))
		if err != nil {
			log.WithError(err).Errorf("Failed to push to %s", base64.StdEncoding.EncodeToString(recipient))
		}
	}
}

// will write status into header and log error if any
func requestError(w http.ResponseWriter, returnCode int, message string) {
	w.WriteHeader(returnCode)
	_, err := fmt.Fprint(w, message)
	if err != nil {
		log.WithError(err).Error("Failed to fmt.Fprintf(w, message)")
	}
}
