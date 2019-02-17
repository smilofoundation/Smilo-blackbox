package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"

	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/encoding"
	"Smilo-blackbox/src/server/syncpeer"
)

const BlackBoxVersion = "Smilo Black Box 0.1.0"
const UpcheckMessage = "I'm up!"

const HeaderFrom = "c11n-from"
const HeaderTo = "c11n-to"
const HeaderKey = "c11n-key"

// Request path "/version", response plain text version ID
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(BlackBoxVersion))
}

// Request path "/upcheck", response plain text upcheck message.
func Upcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(UpcheckMessage))
}

// Request path "/api", response json rest api spec.
func Api(w http.ResponseWriter, r *http.Request) {

}

func UnknownRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug("UnknownEndPoint")
}

func RetrieveJsonPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) {
	payload := RetrieveAndDecryptPayload(w, r, key, to)
	if payload != nil {
		receiveResp := ReceiveResponse{Payload: base64.StdEncoding.EncodeToString(payload)}
		json.NewEncoder(w).Encode(receiveResp)
		w.Header().Set("Content-Type", "application/json")
	}
}

func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte {
	encTrans, err := data.FindEncryptedTransaction(key)
	if err != nil || encTrans == nil {
		message := fmt.Sprintf("Transaction key: %s not found", hex.EncodeToString(key))
		log.Error(message)
		requestError(w, http.StatusNotFound, message)
		return nil
	}

	encodedPayloadData := encoding.Deserialize([]byte(encTrans.Encoded_Payload))
	payload := encodedPayloadData.Decode(to)

	if payload == nil {
		message := fmt.Sprintf("Error Encoding Payload on Request: %s", r.URL)
		log.Error(message)
		requestError(w, http.StatusInternalServerError, message)
	}
	return payload
}

func PushTransactionForOtherNodes(encryptedTransaction data.Encrypted_Transaction, recipient []byte) {
	url, err := syncpeer.GetPeerURL(recipient)
	if err == nil {
		_, err := new(http.Client).Post(url+"/push", "application/octet-stream", bytes.NewBuffer([]byte(base64.StdEncoding.EncodeToString(encryptedTransaction.Encoded_Payload))))
		if err != nil {
			log.WithError(err).Errorf("Failed to push to %s", base64.StdEncoding.EncodeToString(recipient))
		}
	}
}

func requestError(w http.ResponseWriter, returnCode int, message string) {
	w.WriteHeader(returnCode)
	fmt.Fprintf(w, message)
}
