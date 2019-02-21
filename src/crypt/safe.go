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

	"github.com/twystd/tweetnacl-go"
	"sync"
)

var empty_return = []byte("")

var empty_nounce = []byte("\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000")

var computedKeys = make(map[string][]byte)

var mutex = sync.RWMutex{}

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

var keys = make(map[string]KeyPair)
var pairs = make([]KeyPair, 0, 128)

func PutKeyPair(pair KeyPair) {
	keys[string(pair.PublicKey)] = pair
	pairs = append(pairs, pair)
}

func GetPublicKeys() [][]byte {
	publicKeys := make([][]byte, 0, len(pairs))
	for _, pair := range pairs {
		publicKeys = append(publicKeys, pair.PublicKey)
	}
	return publicKeys
}

func GetPrivateKey(publicKey []byte) []byte {
	return keys[string(publicKey)].PrivateKey
}

func NewRandomKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return b, err
}

func NewRandomNonce() ([]byte, error) {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	return b, err
}

func ComputeSharedKey(senderKey []byte, publicKey []byte) []byte {
	var ret []byte
	mutex.RLock()
    ret, ok := computedKeys[string(senderKey)+string(publicKey)]
    mutex.RUnlock()
    if !ok {
    	var err error
		ret, err = tweetnacl.CryptoBoxBeforeNM(publicKey, senderKey)
		if err != nil {
			ret = empty_return
		} else {
			mutex.Lock()
			defer mutex.Unlock()
			computedKeys[string(senderKey)+string(publicKey)] = ret
		}
	}
	return ret
}

func EncryptPayload(sharedKey []byte, payload []byte, nounce []byte) []byte {
	var ret []byte
	if nounce == nil {
		nounce = empty_nounce
	}
	ret, err := tweetnacl.CryptoSecretBox(payload, nounce, sharedKey)
	if err != nil {
		ret = empty_return
	}
	return ret
}

func DecryptPayload(sharedKey []byte, encrypted_payload []byte, nounce []byte) []byte {
	var ret []byte
	if nounce == nil {
		nounce = empty_nounce
	}
	ret, err := tweetnacl.CryptoSecretBoxOpen(encrypted_payload, nounce, sharedKey)
	if err != nil {
		ret = empty_return
	}
	return ret
}
