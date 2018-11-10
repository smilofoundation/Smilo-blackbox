package crypt

import "github.com/twystd/tweetnacl-go"

var empty_return = []byte("")

var empty_nounce = []byte("\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000")

func ComputeSharedKey(senderKey []byte, publicKey []byte) ([]byte) {
	var ret []byte
	ret, err := tweetnacl.CryptoBoxBeforeNM(publicKey, senderKey)
	if err != nil {
		ret = empty_return
	}
	return ret;
}

func EncryptPayload(sharedKey []byte, payload []byte, nounce []byte) ([]byte) {
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

func DecryptPayload(sharedKey []byte, encrypted_payload []byte, nounce []byte) ([]byte) {
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