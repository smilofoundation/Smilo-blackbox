package data

import (
	"bytes"
	"encoding/binary"
)

type Encoded_Payload_Data struct {
	Sender []byte
	Nonce []byte
	Cipher []byte
	RecipientNonce []byte
	RecipientList [][]byte
}

func (e *Encoded_Payload_Data) Serialize() (*[]byte){
    buffer := bytes.NewBuffer([]byte(""))
    encodeBytes(e.Sender, buffer)
	encodeBytes(e.Cipher, buffer)
	encodeBytes(e.Nonce, buffer)
	encodeArray(e.RecipientList, buffer)
	encodeBytes(e.RecipientNonce, buffer)
    ret := buffer.Bytes()
    return &ret
}

func Deserialize(encoded_payload []byte) (*Encoded_Payload_Data) {
    e := Encoded_Payload_Data{}
    buffer := bytes.NewBuffer(encoded_payload)
    e.Sender = *decodeBytes(buffer)
    e.Cipher = *decodeBytes(buffer)
    e.Nonce = *decodeBytes(buffer)
    e.RecipientList = *decodeArray(buffer)
    e.RecipientNonce = *decodeBytes(buffer)
    return &e
}

func encodeBytes(data []byte, buffer *bytes.Buffer) {
	tmp := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp, uint64(len(data)))
	buffer.Write(tmp)
	buffer.Write(data)
}

func encodeArray(data [][]byte, buffer *bytes.Buffer) {
	tmp := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp, uint64(len(data)))
	buffer.Write(tmp)
	for _, i := range data {
		encodeBytes(i, buffer)
	}
}

func decodeBytes(buffer *bytes.Buffer) (*[]byte) {
	sizeB := buffer.Next(8)
	size := binary.BigEndian.Uint64(sizeB)
	data := buffer.Next(int(size))
	return &data
}

func decodeArray(buffer *bytes.Buffer) (*[][]byte) {
	sizeB := buffer.Next(8)
	size := binary.BigEndian.Uint64(sizeB)
	data := make([][]byte, size)
	for i:=uint64(0); i<size; i++ {
		data[i] = *decodeBytes(buffer)
	}
	return &data
}
