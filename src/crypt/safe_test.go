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
	"os"
	"testing"

	"encoding/base64"
	"fmt"

	"github.com/stretchr/testify/require"
	"github.com/twystd/tweetnacl-go"
)

var EMPTY = []byte("")

func TestMain(m *testing.M) {
	retcode := m.Run()
	os.Exit(retcode)
}

func TestComputeSharedKey(t *testing.T) {
	keyPair1, _ := tweetnacl.CryptoBoxKeyPair()
	keyPair2, _ := tweetnacl.CryptoBoxKeyPair()
	sharedKey1 := ComputeSharedKey(keyPair1.SecretKey, keyPair2.PublicKey)
	sharedKey2 := ComputeSharedKey(keyPair2.SecretKey, keyPair1.PublicKey)
	require.Equal(t, sharedKey1, sharedKey2)
}

func TestComputeSharedKey2(t *testing.T) {
	keyPair, _ := tweetnacl.CryptoBoxKeyPair()
	sharedKey := ComputeSharedKey(keyPair.SecretKey, keyPair.PublicKey)
	require.NotEqual(t, sharedKey, EMPTY)
}

func TestEncryptDecryptPayload(t *testing.T) {
	keyPair1, _ := tweetnacl.CryptoBoxKeyPair()
	keyPair2, _ := tweetnacl.CryptoBoxKeyPair()
	sharedKey := ComputeSharedKey(keyPair1.SecretKey, keyPair2.PublicKey)
	payload := []byte("12345678901234567890123456789012345678901234567890123456789012345678901234567890")
	encryptedMessage := EncryptPayload(sharedKey, payload, nil)
	payload2 := DecryptPayload(sharedKey, encryptedMessage, nil)
	require.Equal(t, payload, payload2, "Return (Nounce Zero): "+string(payload2))

	t.Log("Encrypted Message (Nounce Zero): " + string(encryptedMessage))
	var nounce = make([]byte, 24, 24)
	for i := 0; i < 24; i++ {
		nounce[i] = byte(i)
	}
	encryptedMessage = EncryptPayload(sharedKey, payload, nounce)
	payload2 = DecryptPayload(sharedKey, encryptedMessage, nounce)
	require.Equal(t, payload, payload2, "Return: (Nounce Non Zero): "+string(payload2))

	t.Log("Encrypted Message (Nounce Zero): " + string(encryptedMessage))
}

func TestPublicKeyFromPrivateKey(t *testing.T) {
	keyPair, _ := tweetnacl.CryptoBoxKeyPair()
	public, _ := ComputePublicKey(keyPair.SecretKey)
	require.Equal(t, keyPair.PublicKey, public, "Different Public Keys!")
}

func TestGenerateTestKeys(t *testing.T) {
	t.SkipNow()
	for i := byte(1); i < 10; i++ {
		privateKey := make([]byte, 32)
		privateKey[31] = i
		publicKey, _ := ComputePublicKey(privateKey)
		WritePrivateKeyFile(base64.StdEncoding.EncodeToString(privateKey), "../../../keys/testkey"+fmt.Sprint(i)+".key")
		WritePublicKeyFile(base64.StdEncoding.EncodeToString(publicKey), "../../../keys/testkey"+fmt.Sprint(i)+".pub")
	}
}
