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

package types

import (
	"time"

	"golang.org/x/crypto/sha3"
)

// EncryptedTransaction holds hash and payload
type EncryptedTransaction struct {
	Hash           []byte `key:"true"`
	EncodedPayload []byte
	Timestamp      time.Time
}

// NewEncryptedTransaction will create a new encrypted transaction based on the provided payload
func NewEncryptedTransaction(encodedPayload []byte) *EncryptedTransaction {
	trans := EncryptedTransaction{
		Hash:           calculateHash(encodedPayload),
		EncodedPayload: encodedPayload,
		Timestamp:      time.Now(),
	}
	return &trans
}

func calculateHash(encodedPayload []byte) []byte {
	tmp := sha3.Sum512(encodedPayload)
	return tmp[:]
}

// CreateEncryptedTransaction will encrypt the transaction
func CreateEncryptedTransaction(hash []byte, encodedPayload []byte, timestamp time.Time) *EncryptedTransaction {
	trans := EncryptedTransaction{
		Hash:           hash,
		EncodedPayload: encodedPayload,
		Timestamp:      timestamp,
	}
	return &trans
}

// FindEncryptedTransaction will find a encrypted transaction for a hash
func FindEncryptedTransaction(hash []byte) (*EncryptedTransaction, error) {
	var t EncryptedTransaction
	err := DBI.Find("Hash", hash, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

//Save saves into db
func (et *EncryptedTransaction) Save() error {
	return DBI.Save(et)
}

//Delete delete it on the db
func (et *EncryptedTransaction) Delete() error {
	return DBI.Delete(et)
}
