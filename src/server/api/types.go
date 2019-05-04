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
	"fmt"
)

//SendRequest will marshal/unmarshal payload from and to
type SendRequest struct {
	// Payload is the transaction payload data we wish to store.
	Payload string `json:"payload"`
	// From is the sender node identification.
	From string `json:"from"`
	// To is a list of the recipient nodes that should be privy to this transaction payload.
	To []string `json:"to"`
}

//SendResponse marshal/unmarshal a key
type SendResponse struct {
	// Key is the key that can be used to retrieve the submitted transaction.
	Key string `json:"key"`
}

//ReceiveRequest marshal/unmarshal key and to
type ReceiveRequest struct {
	Key string `json:"key"`
	To  string `json:"to"`
}

//ReceiveResponse will marshal/unmarshal payload
type ReceiveResponse struct {
	Payload string `json:"payload"`
}

//DeleteRequest marshal/unmarshal key
type DeleteRequest struct {
	Key string `json:"key"`
}

//ResendRequest will marshal/unmarshal type, pub and pk
type ResendRequest struct {
	// Type is the resend request type. It should be either "all" or "individual" depending on if
	// you want to request an individual transaction, or all transactions associated with a node.
	Type      string `json:"type"`
	PublicKey string `json:"publicKey"`
	Key       string `json:"key,omitempty"`
}

//PeerURL will marshal/unmarshal url
type PeerURL struct {
	URL string `json:"url"`
}

//Parse will process send parsing
func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string) {
	msgs := make([]string, 0, len(e.To)+2)
	payload, err := base64.StdEncoding.DecodeString(e.Payload)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode payload: %s, error: %s", e.Payload, err))
	}
	sender, err := base64.StdEncoding.DecodeString(e.From)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode sender: %s, error: %s", e.From, err))
	}
	recipients := make([][]byte, len(e.To))
	for i, value := range e.To {
		recipient, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Unable to decode recipient: %s, error: %s", value, err))
		} else {
			recipients[i] = recipient
		}
	}
	return payload, sender, recipients, msgs
}

//Parse will process receiving parsing
func (e *ReceiveRequest) Parse() ([]byte, []byte, []string) {
	msgs := make([]string, 0, len(e.To)+2)
	key, err := base64.StdEncoding.DecodeString(e.Key)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode Key: %s, error: %s", e.Key, err))
	}
	to, err := base64.StdEncoding.DecodeString(e.To)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Unable to decode To: %s, error: %s", e.To, err))
	}

	return key, to, msgs
}
