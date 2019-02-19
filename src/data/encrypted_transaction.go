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
