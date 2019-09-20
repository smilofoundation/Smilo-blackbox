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

package crypt

import (
	"crypto/rand"

	"sync"

	"github.com/twystd/tweetnacl-go"
)

var emptyReturn = []byte("")

var emptyNounce = []byte("\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000")

var computedKeys = make(map[string][]byte)

var mutex = sync.RWMutex{}

var ReadRandom = func(b []byte) (n int, err error) {
	return rand.Read(b)
}

// KeyPair holds PrivateKey and PublicKey
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

var keys = make(map[string]KeyPair)
var pairs = make([]KeyPair, 0, 128)

// PutKeyPair will put a pair into pairs var
func PutKeyPair(pair KeyPair) {
	keys[string(pair.PublicKey)] = pair
	pairs = append(pairs, pair)
}

// GetPublicKeys will get all pubs in memory
func GetPublicKeys() [][]byte {
	publicKeys := make([][]byte, 0, len(pairs))
	for _, pair := range pairs {
		publicKeys = append(publicKeys, pair.PublicKey)
	}
	return publicKeys
}

// GetPrivateKey will get the pk for a pub
func GetPrivateKey(publicKey []byte) []byte {
	return keys[string(publicKey)].PrivateKey
}

// NewRandomKey generate new key
func NewRandomKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := ReadRandom(b)
	return b, err
}

// NewRandomNonce generate new nonce
func NewRandomNonce() ([]byte, error) {
	b := make([]byte, 24)
	_, err := ReadRandom(b)
	return b, err
}

// ComputeSharedKey compute a shareKey based on two keys
func ComputeSharedKey(senderKey []byte, publicKey []byte) []byte {
	var ret []byte
	mutex.RLock()
	ret, ok := computedKeys[string(senderKey)+string(publicKey)]
	mutex.RUnlock()
	if !ok {
		var err error
		ret, err = tweetnacl.CryptoBoxBeforeNM(publicKey, senderKey)
		if err != nil {
			ret = emptyReturn
		} else {
			mutex.Lock()
			defer mutex.Unlock()
			computedKeys[string(senderKey)+string(publicKey)] = ret
		}
	}
	return ret
}

// EncryptPayload will encrypt payload based on a key and nonce
func EncryptPayload(sharedKey []byte, payload []byte, nonce []byte) []byte {
	var ret []byte
	if nonce == nil {
		nonce = emptyNounce
	}
	ret, err := tweetnacl.CryptoSecretBox(payload, nonce, sharedKey)
	if err != nil {
		ret = emptyReturn
	}
	return ret
}

// DecryptPayload will decrypt a payload based on a key and nonce
func DecryptPayload(sharedKey []byte, encryptedPayload []byte, nonce []byte) []byte {
	var ret []byte
	if nonce == nil {
		nonce = emptyNounce
	}
	ret, err := tweetnacl.CryptoSecretBoxOpen(encryptedPayload, nonce, sharedKey)
	if err != nil {
		ret = emptyReturn
	}
	return ret
}
