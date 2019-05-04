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

package encoding

import (
	"bytes"
	"encoding/binary"

	"github.com/sirupsen/logrus"

	"io/ioutil"

	"Smilo-blackbox/src/crypt"
)

//EncodedPayloadData holds the encoded payload data (sender,nonce,cipher,recipientNonce and list)
type EncodedPayloadData struct {
	Sender         []byte
	Nonce          []byte
	Cipher         []byte
	RecipientNonce []byte
	RecipientList  [][]byte
}

//Serialize obj
func (e *EncodedPayloadData) Serialize() *[]byte {
	var buffer = bytes.NewBuffer([]byte{})
	serializeBytes(e.Sender, buffer)
	serializeBytes(e.Cipher, buffer)
	serializeBytes(e.Nonce, buffer)
	serializeArray(e.RecipientList, buffer)
	serializeBytes(e.RecipientNonce, buffer)
	ret, _ := ioutil.ReadAll(buffer)
	return &ret
}

//Deserialize obj
func Deserialize(encodedPayload []byte) *EncodedPayloadData {
	e := EncodedPayloadData{}
	buffer := bytes.NewBuffer(encodedPayload)
	e.Sender = deserializeBytes(buffer)
	e.Cipher = deserializeBytes(buffer)
	e.Nonce = deserializeBytes(buffer)
	e.RecipientList = deserializeArray(buffer)
	e.RecipientNonce = deserializeBytes(buffer)
	return &e
}

//EncodePayloadData encode payload data
func EncodePayloadData(payload []byte, sender []byte, recipients [][]byte) (*EncodedPayloadData, error) {
	masterkey, err := crypt.NewRandomKey()
	if err != nil {
		return nil, err
	}
	recipientsNonce, err := crypt.NewRandomNonce()
	if err != nil {
		return nil, err
	}
	nonce, err := crypt.NewRandomNonce()
	if err != nil {
		return nil, err
	}
	cipher := crypt.EncryptPayload(masterkey, payload, nonce)
	senderPrivate := crypt.GetPrivateKey(sender)
	recipientsEncryptedKey := make([][]byte, len(recipients))
	for i := 0; i < len(recipients); i++ {
		sharedKey := crypt.ComputeSharedKey(senderPrivate, recipients[i])
		recipientsEncryptedKey[i] = crypt.EncryptPayload(sharedKey, masterkey, recipientsNonce)
	}
	e := EncodedPayloadData{
		Sender:         sender,
		Cipher:         cipher,
		Nonce:          nonce,
		RecipientNonce: recipientsNonce,
		RecipientList:  recipientsEncryptedKey,
	}
	return &e, nil
}

//Decode will decode
func (e *EncodedPayloadData) Decode(to []byte) []byte {
	var publicKey = to
	privateKey := crypt.GetPrivateKey(e.Sender)
	if privateKey == nil {
		privateKey = crypt.GetPrivateKey(to)
		publicKey = e.Sender
	}
	sharedKey := crypt.ComputeSharedKey(privateKey, publicKey)

	var masterKey []byte
	for _, recipient := range e.RecipientList {
		masterKey = crypt.DecryptPayload(sharedKey, recipient, e.RecipientNonce)
		if len(masterKey) > 0 {
			break
		}
	}
	if len(masterKey) > 0 {
		return crypt.DecryptPayload(masterKey, e.Cipher, e.Nonce)
	}
	return nil
}

func serializeBytes(data []byte, buffer *bytes.Buffer) {
	tmp := make([]byte, 8)
	size := len(data)
	buffer.Grow(size + 8)
	binary.BigEndian.PutUint64(tmp, uint64(size))
	buffer.Write(tmp)
	buffer.Write(data)
}

func serializeArray(data [][]byte, buffer *bytes.Buffer) {
	tmp := make([]byte, 8)
	buffer.Grow(8)
	binary.BigEndian.PutUint64(tmp, uint64(len(data)))
	buffer.Write(tmp)
	for _, i := range data {
		serializeBytes(i, buffer)
	}
}

func deserializeBytes(buffer *bytes.Buffer) []byte {
	var sizeB = make([]byte, 8)
	_, err := buffer.Read(sizeB)
	if err != nil {
		logrus.WithError(err).Error("Could not buffer.Read")
	}
	size := binary.BigEndian.Uint64(sizeB)
	data := buffer.Next(int(size))
	return data
}

func deserializeArray(buffer *bytes.Buffer) [][]byte {
	sizeB := buffer.Next(8)
	size := binary.BigEndian.Uint64(sizeB)
	data := make([][]byte, size)
	for i := uint64(0); i < size; i++ {
		data[i] = deserializeBytes(buffer)
	}
	return data
}
