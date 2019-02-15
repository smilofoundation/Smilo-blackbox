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
)

// It receives headers "c11n-from" and "c11n-to", payload body and returns Status Code 200 and encoded key plain text.
func SendRaw(w http.ResponseWriter, r *http.Request) {
	from := r.Header.Get("c11n-from")
	to := r.Header.Get("c11n-to")

	if from == "" || to == "" {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, invalid headers.\n", r.URL))
		return
	}
	sender, err := base64.StdEncoding.DecodeString(from)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, c11n-from header (%s) is not a valid key.\n", r.URL, from))
		return
	}
	encodedRecipients := strings.Split(to, ",")
	var error []string
	var recipients = make([][]byte, len(encodedRecipients))
	for i := 0; i < len(encodedRecipients); i++ {
		decodedValue, err := base64.StdEncoding.DecodeString(encodedRecipients[i])
		if err != nil {
			error = append(error, fmt.Sprintf("c11n-to header (%s) is not a valid key", encodedRecipients[i]))
		}
		recipients[i] = decodedValue
	}
	if len(error) > 0 {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, %s.", r.URL, strings.Join(error, ", ")))
		return
	}
	encPayload, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if encPayload == nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, missing payload.\n", r.URL))
		return
	}

	payload, err := base64.StdEncoding.DecodeString(string(encPayload))
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error decoding payload: (%s), %s\n", r.URL, encPayload, err))
		return
	}

	encTrans := createNewEncodedTransaction(w, r, payload, sender, recipients)

	if encTrans != nil {
		w.Write([]byte(base64.StdEncoding.EncodeToString(encTrans.Hash)))
		w.Header().Set("Content-Type", "text/plain")
	}
}

// It receives json SendRequest with from, to and payload, returns Status Code 200 and json SendResponse with encoded key.
func Send(w http.ResponseWriter, r *http.Request) {
	var sendReq SendRequest
	err := json.NewDecoder(r.Body).Decode(&sendReq)
	r.Body.Close()
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error: %s\n", r.URL, err))
		return
	}

	payload, sender, recipients, msgs := sendReq.Parse()

	if len(msgs) > 0 {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s\n %s", r.URL, strings.Join(msgs, "\n")))
		return
	}

	encTrans := createNewEncodedTransaction(w, r, payload, sender, recipients)

	if encTrans != nil {
		sendResp := SendResponse{Key: base64.StdEncoding.EncodeToString(encTrans.Hash)}
		json.NewEncoder(w).Encode(sendResp)
		w.Header().Set("Content-Type", "application/json")
	}
}

func createNewEncodedTransaction(w http.ResponseWriter, r *http.Request, payload []byte, sender []byte, recipients [][]byte) *data.Encrypted_Transaction {
	encPayload, err := encoding.EncodePayloadData(payload, sender, recipients)
	if err != nil {
		requestError(w, http.StatusInternalServerError, fmt.Sprintf("Error Encoding Payload on Request: %s\n %s\n", r.URL, err))
		return nil
	}
	encTrans := data.NewEncryptedTransaction(*encPayload.Serialize())
	encTrans.Save()
	for _, recipient := range recipients {
		PushTransactionForOtherNodes(*encTrans, recipient)
	}
	return encTrans
}

// Deprecated API
// It receives a ReceiveRequest json with an encoded key (hash) and to values, returns decrypted payload
func Receive(w http.ResponseWriter, r *http.Request) {
	var receiveReq ReceiveRequest
	err := json.NewDecoder(r.Body).Decode(&receiveReq)
	r.Body.Close()
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error: %s\n", r.URL, err))
		return
	}

	key, to, msgs := receiveReq.Parse()

	if len(msgs) > 0 {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s\n %s", r.URL, strings.Join(msgs, "\n")))
		return
	}

	RetrieveJsonPayload(w, r, key, to)

}

// it receives a GET request with a hash on path and query var "to" with encoded hash and to, returns decrypted payload
func TransactionGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	r.ParseForm()
	encodedTo := r.Form.Get("to")
	hash := params["hash"]
	if hash == "" || encodedTo == "" {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, invalid query.\n", r.URL))
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
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s\n %s", r.URL, strings.Join(errors, "\n")))
		return
	}

	RetrieveJsonPayload(w, r, key, to)

}

func requestError(w http.ResponseWriter, returnCode int, message string) {
	log.Error(message)
	w.WriteHeader(returnCode)
	fmt.Fprintf(w, message)
}
