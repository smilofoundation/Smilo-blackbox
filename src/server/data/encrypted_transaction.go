package data

import (
	"time"
	"github.com/golang/glog"
	"golang.org/x/crypto/sha3"
)

type Encrypted_Transaction struct {
	Hash            string    `storm:"id"`
	Encoded_Payload string
	Timestamp       time.Time `storm:"index"`
}

func NewEncryptedTransaction(encoded_payload string) (*Encrypted_Transaction) {
	trans := Encrypted_Transaction{
		Hash:            calculateHash(encoded_payload),
		Encoded_Payload: encoded_payload,
		Timestamp:       time.Now(),
	}
	return &trans;
}

func calculateHash(encoded_payload string) string {
	tmp := sha3.Sum512([]byte(encoded_payload))
	return string(tmp[:])
}

func CreateEncryptedTransaction(hash string, encoded_payload string, timestamp time.Time) (*Encrypted_Transaction) {
	trans := Encrypted_Transaction{
		Hash:            hash,
		Encoded_Payload: encoded_payload,
		Timestamp:       timestamp,
	}
	return &trans;
}

func FindEncryptedTransaction(hash string) (*Encrypted_Transaction) {
	var t Encrypted_Transaction
	err := db.One("Hash", hash, &t)
	if err != nil {
		glog.Error("Unable to find transaction.")
		return nil
	}
	return &t
}

func (et *Encrypted_Transaction) Save() {
	db.Save(et)
}

func (et *Encrypted_Transaction) Delete() {
	db.DeleteStruct(et)
}
