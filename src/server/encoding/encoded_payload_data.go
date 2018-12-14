package encoding

import (
	"bytes"
	"encoding/binary"

	"io/ioutil"

	"Smilo-blackbox/src/crypt"
)

type Encoded_Payload_Data struct {
	Sender         []byte
	Nonce          []byte
	Cipher         []byte
	RecipientNonce []byte
	RecipientList  [][]byte
}

func (e *Encoded_Payload_Data) Serialize() *[]byte {
	var buffer = bytes.NewBuffer([]byte{})
	serializeBytes(e.Sender, buffer)
	serializeBytes(e.Cipher, buffer)
	serializeBytes(e.Nonce, buffer)
	serializeArray(e.RecipientList, buffer)
	serializeBytes(e.RecipientNonce, buffer)
	ret, _ := ioutil.ReadAll(buffer)
	return &ret
}

func Deserialize(encodedPayload []byte) *Encoded_Payload_Data {
	e := Encoded_Payload_Data{}
	buffer := bytes.NewBuffer(encodedPayload)
	e.Sender = deserializeBytes(buffer)
	e.Cipher = deserializeBytes(buffer)
	e.Nonce = deserializeBytes(buffer)
	e.RecipientList = deserializeArray(buffer)
	e.RecipientNonce = deserializeBytes(buffer)
	return &e
}

func EncodePayloadData(payload []byte, sender []byte, recipients [][]byte) (*Encoded_Payload_Data, error) {
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
	e := Encoded_Payload_Data{
		Sender:         sender,
		Cipher:         cipher,
		Nonce:          nonce,
		RecipientNonce: recipientsNonce,
		RecipientList:  recipientsEncryptedKey,
	}
	return &e, nil
}

func (e *Encoded_Payload_Data) Decode(to []byte) []byte {
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
	} else {
		return nil
	}
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
	buffer.Read(sizeB)
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
