package api

import (
	"encoding/base64"
	"fmt"
)

type SendRequest struct {
	// Payload is the transaction payload data we wish to store.
	Payload string `json:"payload"`
	// From is the sender node identification.
	From string `json:"from"`
	// To is a list of the recipient nodes that should be privy to this transaction payload.
	To []string `json:"to"`
}

type SendResponse struct {
	// Key is the key that can be used to retrieve the submitted transaction.
	Key string `json:"key"`
}

type ReceiveRequest struct {
	Key string `json:"key"`
	To  string `json:"to"`
}

type ReceiveResponse struct {
	Payload string `json:"payload"`
}

type DeleteRequest struct {
	Key string `json:"key"`
}

type ResendRequest struct {
	// Type is the resend request type. It should be either "all" or "individual" depending on if
	// you want to request an individual transaction, or all transactions associated with a node.
	Type      string `json:"type"`
	PublicKey string `json:"publicKey"`
	Key       string `json:"key,omitempty"`
}

type PeerUrl struct {
	Url string `json:"url"`
}

func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string) {
	msgs := make([]string, 0, len(e.To)+2)
	payload, err := base64.StdEncoding.DecodeString(e.Payload)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode payload: %s, error: %s\n", e.Payload, err))
	}
	sender, err := base64.StdEncoding.DecodeString(e.From)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode sender: %s, error: %s\n", e.From, err))
	}
	recipients := make([][]byte, len(e.To))
	for i, value := range e.To {
		recipient, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Unable to decode recipient: %s, error: %s\n", value, err))
		} else {
			recipients[i] = recipient
		}
	}
	return payload, sender, recipients, msgs
}

func (e *ReceiveRequest) Parse() ([]byte, []byte, []string) {
	msgs := make([]string, 0, len(e.To)+2)
	key, err := base64.StdEncoding.DecodeString(e.Key)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode Key: %s, error: %s\n", e.Key, err))
	}
	to, err := base64.StdEncoding.DecodeString(e.To)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode To: %s, error: %s\n", e.To, err))
	}

	return key, to, msgs
}
