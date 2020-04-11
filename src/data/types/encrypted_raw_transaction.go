package types

import "time"

// EncryptedRawTransaction holds hash and payload
type EncryptedRawTransaction struct {
	Hash           []byte `key:"true"`
	EncodedPayload []byte
	Sender         []byte
	Timestamp      time.Time
}

// NewEncryptedRawTransaction will create a new encrypted transaction based on the provided payload
func NewEncryptedRawTransaction(encodedPayload []byte, sender []byte) *EncryptedRawTransaction {
	trans := EncryptedRawTransaction{
		Hash:           calculateHash(encodedPayload),
		EncodedPayload: encodedPayload,
		Sender:         sender,
		Timestamp:      time.Now(),
	}
	return &trans
}

// FindEncryptedRawTransaction will find a encrypted transaction for a hash
func FindEncryptedRawTransaction(hash []byte) (*EncryptedRawTransaction, error) {
	var t EncryptedRawTransaction
	err := DBI.Find("Hash", hash, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

//Save saves into db
func (et *EncryptedRawTransaction) Save() error {
	return DBI.Save(et)
}

//Delete delete it on the db
func (et *EncryptedRawTransaction) Delete() error {
	return DBI.Delete(et)
}
