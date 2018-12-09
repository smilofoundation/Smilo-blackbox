package crypt

import (
	"bytes"
	"os"
	"testing"

	"github.com/twystd/tweetnacl-go"
	"github.com/ethereum/go-ethereum/log"
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
	if !bytes.Equal(sharedKey1, sharedKey2) {
		log.Error("Shared Key 1: " + string(sharedKey1))
		log.Error("Shared Key 2: " + string(sharedKey2))
		t.Fail()
	}
}

func TestComputeSharedKey2(t *testing.T) {
	keyPair, _ := tweetnacl.CryptoBoxKeyPair()
	sharedKey := ComputeSharedKey(keyPair.SecretKey, keyPair.PublicKey)
	if bytes.Equal(sharedKey, EMPTY) {
		t.Fail()
	}
}

func TestEncryptDecryptPayload(t *testing.T) {
	keyPair1, _ := tweetnacl.CryptoBoxKeyPair()
	keyPair2, _ := tweetnacl.CryptoBoxKeyPair()
	sharedKey := ComputeSharedKey(keyPair1.SecretKey, keyPair2.PublicKey)
	payload := []byte("12345678901234567890123456789012345678901234567890123456789012345678901234567890")
	encryptedMessage := EncryptPayload(sharedKey, payload, nil)
	payload2 := DecryptPayload(sharedKey, encryptedMessage, nil)
	if !bytes.Equal(payload, payload2) {
		log.Error("Return (Nounce Zero): " + string(payload2))
		t.Fail()
	}
	log.Info("Encrypted Message (Nounce Zero): " + string(encryptedMessage))
	var nounce = make([]byte, 24, 24)
	for i := 0; i < 24; i++ {
		nounce[i] = byte(i)
	}
	encryptedMessage = EncryptPayload(sharedKey, payload, nounce)
	payload2 = DecryptPayload(sharedKey, encryptedMessage, nounce)
	if !bytes.Equal(payload, payload2) {
		log.Error("Return: (Nounce Non Zero): " + string(payload2))
		t.Fail()
	}
	log.Info("Encrypted Message (Nounce Zero): " + string(encryptedMessage))
}
