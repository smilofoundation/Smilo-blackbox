package crypt

import (
	"os"
	"testing"

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
