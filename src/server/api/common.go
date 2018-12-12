package api

import (
	"net/http"
	"Smilo-blackbox/src/data"
	"fmt"
	"encoding/hex"
	"Smilo-blackbox/src/server/encoding"
	"encoding/base64"
	"encoding/json"
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
func RetrieveJsonPayload(key []byte, w http.ResponseWriter, to []byte, r *http.Request) {
	payload := RetrieveAndDecryptPayload(key, w, to, r)
	if payload != nil {
		receiveResp := ReceiveResponse{Payload: base64.StdEncoding.EncodeToString(payload)}
		json.NewEncoder(w).Encode(receiveResp)
		w.Header().Set("Content-Type", "application/json")
	}
}

func RetrieveAndDecryptPayload(key []byte, w http.ResponseWriter, to []byte, r *http.Request) []byte {
	encTrans := data.FindEncryptedTransaction(key)
	if encTrans == nil {
		requestError(http.StatusNotFound, w, fmt.Sprintf("Transaction key: %s not found\n", hex.EncodeToString(key)))
		return nil
	}
	encodedPayloadData := encoding.Deserialize([]byte(encTrans.Encoded_Payload))
	payload := encodedPayloadData.Decode(to)
	if payload == nil {
		requestError(http.StatusInternalServerError, w, fmt.Sprintf("Error Encoding Payload on Request: %s\n", r.URL))
	}
	return payload
}