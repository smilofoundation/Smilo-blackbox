package data

import (
	"time"

	"golang.org/x/crypto/sha3"
)

type Encrypted_Transaction struct {
	Hash            []byte `storm:"id"`
	Encoded_Payload []byte
	Timestamp       time.Time `storm:"index"`
}

func NewEncryptedTransaction(encoded_payload []byte) *Encrypted_Transaction {
	trans := Encrypted_Transaction{
		Hash:            calculateHash(encoded_payload),
		Encoded_Payload: encoded_payload,
		Timestamp:       time.Now(),
	}
	return &trans
}

func calculateHash(encoded_payload []byte) []byte {
	tmp := sha3.Sum512(encoded_payload)
	return tmp[:]
}

func CreateEncryptedTransaction(hash []byte, encoded_payload []byte, timestamp time.Time) *Encrypted_Transaction {
	trans := Encrypted_Transaction{
		Hash:            hash,
		Encoded_Payload: encoded_payload,
		Timestamp:       timestamp,
	}
	return &trans
}

func FindEncryptedTransaction(hash []byte) (*Encrypted_Transaction, error) {
	var t Encrypted_Transaction
	err := db.One("Hash", hash, &t)
	if err != nil {
		log.Error("Unable to find transaction.")
		return nil, err
	}
	return &t, nil
}

func (et *Encrypted_Transaction) Save() error {
	return db.Save(et)
}

func (et *Encrypted_Transaction) Delete() error {
	return db.DeleteStruct(et)
}
