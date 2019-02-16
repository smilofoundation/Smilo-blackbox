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

	"Smilo-blackbox/src/crypt"
)

// It receives headers "c11n-from" and "c11n-to", payload body and returns Status Code 200 and encoded key plain text.
func SendRaw(w http.ResponseWriter, r *http.Request) {
	var fromEncoded []byte
	var err error

	from := r.Header.Get("c11n-from")
	to := r.Header.Get("c11n-to")

	if to == "" {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, invalid headers. to:%s", r.URL, to))
		return
	}

	if from != "" {
		fromEncoded, err = base64.StdEncoding.DecodeString(from)
		if err != nil {
			requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, c11n-from header (%s) is not a valid key.\n", r.URL, from))
			return
		}
	} else {
		//use default
		defaultPubKey := base64.StdEncoding.EncodeToString(crypt.GetPublicKeys()[0])
		fromEncoded, err = base64.StdEncoding.DecodeString(defaultPubKey)
		log.WithField("defaultPubKey", defaultPubKey).Info("Request from NOT filled, will use default PubKey")
		if err != nil {
			requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, c11n-from header (%s) is not a valid key.\n", r.URL, from))
			return
		}
	}

	encodedRecipients := strings.Split(to, ",")
	var errors []string
	var recipients = make([][]byte, len(encodedRecipients))
	for i := 0; i < len(encodedRecipients); i++ {
		decodedValue, err := base64.StdEncoding.DecodeString(encodedRecipients[i])
		if err != nil {
			errors = append(errors, fmt.Sprintf("c11n-to header (%s) is not a valid key", encodedRecipients[i]))
		}
		recipients[i] = decodedValue
	}
	if len(errors) > 0 {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, %s.", r.URL, strings.Join(errors, ", ")))
		return
	}
	encPayload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || encPayload == nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, missing payload, err: %s", r.URL, err))
		return
	}

	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(encPayload)))
	n, err := base64.StdEncoding.Decode(dbuf, encPayload)
	if err != nil {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, error decoding payload: (%s), %s", r.URL, encPayload, err))
		return
	}

	payload := dbuf[:n]
	if len(payload) == 0 {
		requestError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %s, len of payload after decode is zero: (%s), %s", r.URL, encPayload, err))
	}

	encTrans := createNewEncodedTransaction(w, r, payload, fromEncoded, recipients)

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

func createNewEncodedTransaction(w http.ResponseWriter, r *http.Request, payload []byte, fromEncoded []byte, recipients [][]byte) *data.Encrypted_Transaction {
	encPayload, err := encoding.EncodePayloadData(payload, fromEncoded, recipients)
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
